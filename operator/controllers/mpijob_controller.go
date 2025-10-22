package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	trainingv2 "github.com/kubeflow/training-operator/pkg/apis/kubeflow.org/v2"
)

const (
	// FinalizerName is the finalizer name for MPIJob
	FinalizerName = "mpijob.kubeflow.org/finalizer"

	// RequeueDelay is the delay for requeuing reconciliation
	RequeueDelay = 10 * time.Second

	// ControllerName is the name of the controller
	ControllerName = "mpijob-controller"
)

// MPIJobReconciler reconciles a MPIJob object
type MPIJobReconciler struct {
	client.Client
	Log          logr.Logger
	Scheme       *runtime.Scheme
	StatusTracker *MPIJobStatusTracker
}

// +kubebuilder:rbac:groups=kubeflow.org,resources=mpijobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kubeflow.org,resources=mpijobs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=kubeflow.org,resources=mpijobs/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop
func (r *MPIJobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("mpijob", req.NamespacedName)

	// Fetch the MPIJob instance
	mpijob := &trainingv2.MPIJob{}
	err := r.Get(ctx, req.NamespacedName, mpijob)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("MPIJob resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get MPIJob")
		return ctrl.Result{}, err
	}

	// Check if the MPIJob is being deleted
	if !mpijob.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, mpijob, log)
	}

	// Add finalizer if not present
	if !controllerutil.ContainsFinalizer(mpijob, FinalizerName) {
		controllerutil.AddFinalizer(mpijob, FinalizerName)
		if err := r.Update(ctx, mpijob); err != nil {
			log.Error(err, "Failed to add finalizer")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// Initialize status if needed
	if mpijob.Status.Conditions == nil {
		mpijob.Status.Conditions = []metav1.Condition{}
	}
	if mpijob.Status.ReplicaStatuses == nil {
		mpijob.Status.ReplicaStatuses = make(map[string]*trainingv2.ReplicaStatus)
	}

	// Reconcile the MPIJob
	result, err := r.reconcileMPIJob(ctx, mpijob, log)
	if err != nil {
		log.Error(err, "Failed to reconcile MPIJob")
		r.StatusTracker.UpdateCondition(mpijob, trainingv2.JobFailed, "ReconcileError", err.Error())
		if statusErr := r.Status().Update(ctx, mpijob); statusErr != nil {
			log.Error(statusErr, "Failed to update status")
		}
		return result, err
	}

	// Update status
	if err := r.Status().Update(ctx, mpijob); err != nil {
		log.Error(err, "Failed to update MPIJob status")
		return ctrl.Result{}, err
	}

	return result, nil
}

// reconcileMPIJob performs the main reconciliation logic
func (r *MPIJobReconciler) reconcileMPIJob(ctx context.Context, mpijob *trainingv2.MPIJob, log logr.Logger) (ctrl.Result, error) {
	// 1. Create launcher pod if not exists
	launcherCreated, err := r.reconcileLauncher(ctx, mpijob, log)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to reconcile launcher: %w", err)
	}

	// 2. Create worker pods if not exists
	workersCreated, err := r.reconcileWorkers(ctx, mpijob, log)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to reconcile workers: %w", err)
	}

	// 3. Create headless service for worker discovery
	serviceCreated, err := r.reconcileService(ctx, mpijob, log)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to reconcile service: %w", err)
	}

	// 4. Create ConfigMap for MPI hostfile if needed
	configMapCreated, err := r.reconcileConfigMap(ctx, mpijob, log)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to reconcile configmap: %w", err)
	}

	// 5. Update job status based on pod states
	if err := r.StatusTracker.UpdateJobStatus(ctx, mpijob, r.Client); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to update job status: %w", err)
	}

	// 6. Handle job completion
	if r.isJobFinished(mpijob) {
		log.Info("MPIJob finished", "status", r.getJobPhase(mpijob))
		return r.handleJobCompletion(ctx, mpijob, log)
	}

	// 7. If we just created resources, requeue to check status
	if launcherCreated || workersCreated || serviceCreated || configMapCreated {
		log.Info("Resources created, requeuing to check status")
		return ctrl.Result{RequeueAfter: RequeueDelay}, nil
	}

	// 8. Requeue to monitor job progress
	return ctrl.Result{RequeueAfter: RequeueDelay}, nil
}

// reconcileLauncher creates or updates the launcher pod
func (r *MPIJobReconciler) reconcileLauncher(ctx context.Context, mpijob *trainingv2.MPIJob, log logr.Logger) (bool, error) {
	launcherSpec, ok := mpijob.Spec.MPIReplicaSpecs["Launcher"]
	if !ok {
		return false, fmt.Errorf("launcher spec not found")
	}

	// Generate launcher pod name
	launcherName := fmt.Sprintf("%s-launcher", mpijob.Name)

	// Check if launcher already exists
	existingPod := &corev1.Pod{}
	err := r.Get(ctx, types.NamespacedName{Name: launcherName, Namespace: mpijob.Namespace}, existingPod)
	if err == nil {
		// Launcher already exists
		log.V(1).Info("Launcher pod already exists", "name", launcherName)
		return false, nil
	}
	if !errors.IsNotFound(err) {
		return false, fmt.Errorf("failed to get launcher pod: %w", err)
	}

	// Create launcher pod
	launcherPod := r.createLauncherPod(mpijob, &launcherSpec, launcherName)
	if err := controllerutil.SetControllerReference(mpijob, launcherPod, r.Scheme); err != nil {
		return false, fmt.Errorf("failed to set owner reference: %w", err)
	}

	if err := r.Create(ctx, launcherPod); err != nil {
		return false, fmt.Errorf("failed to create launcher pod: %w", err)
	}

	log.Info("Created launcher pod", "name", launcherName)
	r.StatusTracker.UpdateCondition(mpijob, trainingv2.JobRunning, "LauncherCreated", "Launcher pod created")
	return true, nil
}

// reconcileWorkers creates or updates worker pods
func (r *MPIJobReconciler) reconcileWorkers(ctx context.Context, mpijob *trainingv2.MPIJob, log logr.Logger) (bool, error) {
	workerSpec, ok := mpijob.Spec.MPIReplicaSpecs["Worker"]
	if !ok {
		return false, fmt.Errorf("worker spec not found")
	}

	replicas := int32(1)
	if workerSpec.Replicas != nil {
		replicas = *workerSpec.Replicas
	}

	anyCreated := false
	for i := int32(0); i < replicas; i++ {
		workerName := fmt.Sprintf("%s-worker-%d", mpijob.Name, i)

		// Check if worker already exists
		existingPod := &corev1.Pod{}
		err := r.Get(ctx, types.NamespacedName{Name: workerName, Namespace: mpijob.Namespace}, existingPod)
		if err == nil {
			// Worker already exists
			continue
		}
		if !errors.IsNotFound(err) {
			return anyCreated, fmt.Errorf("failed to get worker pod %s: %w", workerName, err)
		}

		// Create worker pod
		workerPod := r.createWorkerPod(mpijob, &workerSpec, workerName, i)
		if err := controllerutil.SetControllerReference(mpijob, workerPod, r.Scheme); err != nil {
			return anyCreated, fmt.Errorf("failed to set owner reference: %w", err)
		}

		if err := r.Create(ctx, workerPod); err != nil {
			return anyCreated, fmt.Errorf("failed to create worker pod %s: %w", workerName, err)
		}

		log.Info("Created worker pod", "name", workerName, "index", i)
		anyCreated = true
	}

	if anyCreated {
		r.StatusTracker.UpdateCondition(mpijob, trainingv2.JobRunning, "WorkersCreated", fmt.Sprintf("Created %d worker pods", replicas))
	}

	return anyCreated, nil
}

// reconcileService creates or updates the headless service for worker discovery
func (r *MPIJobReconciler) reconcileService(ctx context.Context, mpijob *trainingv2.MPIJob, log logr.Logger) (bool, error) {
	serviceName := fmt.Sprintf("%s-worker", mpijob.Name)

	// Check if service already exists
	existingService := &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{Name: serviceName, Namespace: mpijob.Namespace}, existingService)
	if err == nil {
		// Service already exists
		log.V(1).Info("Service already exists", "name", serviceName)
		return false, nil
	}
	if !errors.IsNotFound(err) {
		return false, fmt.Errorf("failed to get service: %w", err)
	}

	// Create service
	service := r.createWorkerService(mpijob, serviceName)
	if err := controllerutil.SetControllerReference(mpijob, service, r.Scheme); err != nil {
		return false, fmt.Errorf("failed to set owner reference: %w", err)
	}

	if err := r.Create(ctx, service); err != nil {
		return false, fmt.Errorf("failed to create service: %w", err)
	}

	log.Info("Created worker service", "name", serviceName)
	return true, nil
}

// reconcileConfigMap creates or updates the ConfigMap for MPI hostfile
func (r *MPIJobReconciler) reconcileConfigMap(ctx context.Context, mpijob *trainingv2.MPIJob, log logr.Logger) (bool, error) {
	configMapName := fmt.Sprintf("%s-config", mpijob.Name)

	// Check if configmap already exists
	existingCM := &corev1.ConfigMap{}
	err := r.Get(ctx, types.NamespacedName{Name: configMapName, Namespace: mpijob.Namespace}, existingCM)
	if err == nil {
		// ConfigMap already exists, update hostfile if needed
		return r.updateConfigMap(ctx, mpijob, existingCM, log)
	}
	if !errors.IsNotFound(err) {
		return false, fmt.Errorf("failed to get configmap: %w", err)
	}

	// Create configmap
	configMap := r.createConfigMap(mpijob, configMapName)
	if err := controllerutil.SetControllerReference(mpijob, configMap, r.Scheme); err != nil {
		return false, fmt.Errorf("failed to set owner reference: %w", err)
	}

	if err := r.Create(ctx, configMap); err != nil {
		return false, fmt.Errorf("failed to create configmap: %w", err)
	}

	log.Info("Created configmap", "name", configMapName)
	return true, nil
}

// reconcileDelete handles the deletion of an MPIJob
func (r *MPIJobReconciler) reconcileDelete(ctx context.Context, mpijob *trainingv2.MPIJob, log logr.Logger) (ctrl.Result, error) {
	log.Info("Deleting MPIJob")

	// Cleanup is handled by Kubernetes garbage collection via owner references
	// Just remove the finalizer
	controllerutil.RemoveFinalizer(mpijob, FinalizerName)
	if err := r.Update(ctx, mpijob); err != nil {
		log.Error(err, "Failed to remove finalizer")
		return ctrl.Result{}, err
	}

	log.Info("Successfully deleted MPIJob")
	return ctrl.Result{}, nil
}

// handleJobCompletion handles completion of the job (success or failure)
func (r *MPIJobReconciler) handleJobCompletion(ctx context.Context, mpijob *trainingv2.MPIJob, log logr.Logger) (ctrl.Result, error) {
	// Set completion time if not set
	if mpijob.Status.CompletionTime == nil {
		now := metav1.Now()
		mpijob.Status.CompletionTime = &now
	}

	// Handle TTL cleanup if configured
	if mpijob.Spec.RunPolicy != nil && mpijob.Spec.RunPolicy.TTLSecondsAfterFinished != nil {
		ttl := time.Duration(*mpijob.Spec.RunPolicy.TTLSecondsAfterFinished) * time.Second
		if time.Since(mpijob.Status.CompletionTime.Time) > ttl {
			log.Info("TTL expired, deleting MPIJob")
			if err := r.Delete(ctx, mpijob); err != nil {
				return ctrl.Result{}, fmt.Errorf("failed to delete MPIJob: %w", err)
			}
			return ctrl.Result{}, nil
		}
		// Requeue to check TTL again
		return ctrl.Result{RequeueAfter: ttl}, nil
	}

	// Job is finished, no need to requeue
	return ctrl.Result{}, nil
}

// isJobFinished checks if the job has finished (succeeded or failed)
func (r *MPIJobReconciler) isJobFinished(mpijob *trainingv2.MPIJob) bool {
	for _, condition := range mpijob.Status.Conditions {
		if condition.Type == string(trainingv2.JobSucceeded) && condition.Status == metav1.ConditionTrue {
			return true
		}
		if condition.Type == string(trainingv2.JobFailed) && condition.Status == metav1.ConditionTrue {
			return true
		}
	}
	return false
}

// getJobPhase returns the current phase of the job
func (r *MPIJobReconciler) getJobPhase(mpijob *trainingv2.MPIJob) string {
	if len(mpijob.Status.Conditions) == 0 {
		return "Unknown"
	}
	latestCondition := mpijob.Status.Conditions[len(mpijob.Status.Conditions)-1]
	return latestCondition.Type
}

// createLauncherPod creates a launcher pod spec
func (r *MPIJobReconciler) createLauncherPod(mpijob *trainingv2.MPIJob, spec *trainingv2.ReplicaSpec, name string) *corev1.Pod {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: mpijob.Namespace,
			Labels: map[string]string{
				"mpijob-name":   mpijob.Name,
				"mpijob-role":   "launcher",
				"app.kubernetes.io/name": "mpijob",
				"app.kubernetes.io/component": "launcher",
			},
		},
		Spec: *spec.Template.Spec.DeepCopy(),
	}

	// Add environment variables
	for i := range pod.Spec.Containers {
		pod.Spec.Containers[i].Env = append(pod.Spec.Containers[i].Env,
			corev1.EnvVar{Name: "MPIJOB_NAME", Value: mpijob.Name},
			corev1.EnvVar{Name: "MPIJOB_ROLE", Value: "launcher"},
		)
	}

	return pod
}

// createWorkerPod creates a worker pod spec
func (r *MPIJobReconciler) createWorkerPod(mpijob *trainingv2.MPIJob, spec *trainingv2.ReplicaSpec, name string, index int32) *corev1.Pod {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: mpijob.Namespace,
			Labels: map[string]string{
				"mpijob-name":   mpijob.Name,
				"mpijob-role":   "worker",
				"mpijob-index":  fmt.Sprintf("%d", index),
				"app.kubernetes.io/name": "mpijob",
				"app.kubernetes.io/component": "worker",
			},
		},
		Spec: *spec.Template.Spec.DeepCopy(),
	}

	// Add environment variables
	for i := range pod.Spec.Containers {
		pod.Spec.Containers[i].Env = append(pod.Spec.Containers[i].Env,
			corev1.EnvVar{Name: "MPIJOB_NAME", Value: mpijob.Name},
			corev1.EnvVar{Name: "MPIJOB_ROLE", Value: "worker"},
			corev1.EnvVar{Name: "MPIJOB_WORKER_INDEX", Value: fmt.Sprintf("%d", index)},
		)
	}

	return pod
}

// createWorkerService creates a headless service for worker discovery
func (r *MPIJobReconciler) createWorkerService(mpijob *trainingv2.MPIJob, name string) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: mpijob.Namespace,
			Labels: map[string]string{
				"mpijob-name": mpijob.Name,
				"app.kubernetes.io/name": "mpijob",
				"app.kubernetes.io/component": "service",
			},
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: "None", // Headless service
			Selector: map[string]string{
				"mpijob-name": mpijob.Name,
				"mpijob-role": "worker",
			},
			Ports: []corev1.ServicePort{
				{
					Name:     "mpi",
					Port:     22,
					Protocol: corev1.ProtocolTCP,
				},
			},
		},
	}
}

// createConfigMap creates a ConfigMap with MPI hostfile
func (r *MPIJobReconciler) createConfigMap(mpijob *trainingv2.MPIJob, name string) *corev1.ConfigMap {
	hostfile := r.generateHostfile(mpijob)

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: mpijob.Namespace,
			Labels: map[string]string{
				"mpijob-name": mpijob.Name,
				"app.kubernetes.io/name": "mpijob",
				"app.kubernetes.io/component": "config",
			},
		},
		Data: map[string]string{
			"hostfile": hostfile,
		},
	}
}

// updateConfigMap updates the ConfigMap with current hostfile
func (r *MPIJobReconciler) updateConfigMap(ctx context.Context, mpijob *trainingv2.MPIJob, cm *corev1.ConfigMap, log logr.Logger) (bool, error) {
	newHostfile := r.generateHostfile(mpijob)
	if cm.Data["hostfile"] == newHostfile {
		return false, nil
	}

	cm.Data["hostfile"] = newHostfile
	if err := r.Update(ctx, cm); err != nil {
		return false, fmt.Errorf("failed to update configmap: %w", err)
	}

	log.Info("Updated hostfile in configmap")
	return true, nil
}

// generateHostfile generates MPI hostfile content
func (r *MPIJobReconciler) generateHostfile(mpijob *trainingv2.MPIJob) string {
	workerSpec, ok := mpijob.Spec.MPIReplicaSpecs["Worker"]
	if !ok {
		return ""
	}

	replicas := int32(1)
	if workerSpec.Replicas != nil {
		replicas = *workerSpec.Replicas
	}

	slotsPerWorker := 1
	if mpijob.Spec.SlotsPerWorker != nil {
		slotsPerWorker = int(*mpijob.Spec.SlotsPerWorker)
	}

	hostfile := ""
	for i := int32(0); i < replicas; i++ {
		workerName := fmt.Sprintf("%s-worker-%d", mpijob.Name, i)
		hostfile += fmt.Sprintf("%s slots=%d\n", workerName, slotsPerWorker)
	}

	return hostfile
}

// SetupWithManager sets up the controller with the Manager
func (r *MPIJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// Create predicate to filter events
	pred := predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			// Ignore updates to status subresource
			if e.ObjectOld.GetGeneration() == e.ObjectNew.GetGeneration() {
				return false
			}
			return true
		},
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&trainingv2.MPIJob{}).
		Owns(&corev1.Pod{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.ConfigMap{}).
		WithEventFilter(pred).
		Complete(r)
}
