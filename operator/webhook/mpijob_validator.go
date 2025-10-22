package webhook

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	trainingv2 "github.com/kubeflow/training-operator/pkg/apis/kubeflow.org/v2"
)

// MPIJobValidator validates MPIJob resources
type MPIJobValidator struct {
	Client  client.Client
	decoder *admission.Decoder
}

// +kubebuilder:webhook:path=/validate/mpijob,mutating=false,failurePolicy=fail,groups=kubeflow.org,resources=mpijobs,verbs=create;update,versions=v2,name=mpijob-validation.kubeflow.org

// Handle validates the MPIJob resource
func (v *MPIJobValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	mpijob := &trainingv2.MPIJob{}

	err := v.decoder.Decode(req, mpijob)
	if err != nil {
		return admission.Errored(400, err)
	}

	// Run all validation checks
	if err := v.validateMPIJob(ctx, mpijob); err != nil {
		return admission.Denied(err.Error())
	}

	return admission.Allowed("")
}

// InjectDecoder injects the decoder
func (v *MPIJobValidator) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}

// validateMPIJob performs comprehensive validation of the MPIJob resource
func (v *MPIJobValidator) validateMPIJob(ctx context.Context, mpijob *trainingv2.MPIJob) error {
	// 1. Validate replica specifications
	if err := v.validateReplicaSpecs(mpijob); err != nil {
		return fmt.Errorf("replica spec validation failed: %w", err)
	}

	// 2. Validate resource specifications
	if err := v.validateResourceSpecs(mpijob); err != nil {
		return fmt.Errorf("resource spec validation failed: %w", err)
	}

	// 3. Validate slots per worker
	if err := v.validateSlotsPerWorker(mpijob); err != nil {
		return fmt.Errorf("slots per worker validation failed: %w", err)
	}

	// 4. Validate MPI implementation
	if err := v.validateMPIImplementation(mpijob); err != nil {
		return fmt.Errorf("MPI implementation validation failed: %w", err)
	}

	// 5. Validate run policy
	if err := v.validateRunPolicy(mpijob); err != nil {
		return fmt.Errorf("run policy validation failed: %w", err)
	}

	// 6. Validate scheduling policy
	if err := v.validateSchedulingPolicy(ctx, mpijob); err != nil {
		return fmt.Errorf("scheduling policy validation failed: %w", err)
	}

	// 7. Validate container specs
	if err := v.validateContainerSpecs(mpijob); err != nil {
		return fmt.Errorf("container spec validation failed: %w", err)
	}

	// 8. Validate resource quotas
	if err := v.validateResourceQuotas(ctx, mpijob); err != nil {
		return fmt.Errorf("resource quota validation failed: %w", err)
	}

	// 9. Validate security context
	if err := v.validateSecurityContext(mpijob); err != nil {
		return fmt.Errorf("security context validation failed: %w", err)
	}

	// 10. Validate network policy
	if err := v.validateNetworkPolicy(mpijob); err != nil {
		return fmt.Errorf("network policy validation failed: %w", err)
	}

	return nil
}

// validateReplicaSpecs validates the replica specifications
func (v *MPIJobValidator) validateReplicaSpecs(mpijob *trainingv2.MPIJob) error {
	if mpijob.Spec.MPIReplicaSpecs == nil {
		return fmt.Errorf("mpiReplicaSpecs is required")
	}

	// Validate Launcher
	launcher, ok := mpijob.Spec.MPIReplicaSpecs["Launcher"]
	if !ok {
		return fmt.Errorf("Launcher replica spec is required")
	}
	if launcher.Replicas == nil || *launcher.Replicas != 1 {
		return fmt.Errorf("Launcher replicas must be exactly 1, got %v", launcher.Replicas)
	}

	// Validate Worker
	worker, ok := mpijob.Spec.MPIReplicaSpecs["Worker"]
	if !ok {
		return fmt.Errorf("Worker replica spec is required")
	}
	if worker.Replicas == nil || *worker.Replicas < 1 {
		return fmt.Errorf("Worker replicas must be >= 1, got %v", worker.Replicas)
	}

	// Validate template is specified
	if launcher.Template == nil {
		return fmt.Errorf("Launcher template is required")
	}
	if worker.Template == nil {
		return fmt.Errorf("Worker template is required")
	}

	return nil
}

// validateResourceSpecs validates resource requests and limits
func (v *MPIJobValidator) validateResourceSpecs(mpijob *trainingv2.MPIJob) error {
	// Validate launcher resources
	launcherSpec := mpijob.Spec.MPIReplicaSpecs["Launcher"]
	if err := v.validatePodResourceSpecs(&launcherSpec.Template.Spec, "Launcher"); err != nil {
		return err
	}

	// Validate worker resources
	workerSpec := mpijob.Spec.MPIReplicaSpecs["Worker"]
	if err := v.validatePodResourceSpecs(&workerSpec.Template.Spec, "Worker"); err != nil {
		return err
	}

	return nil
}

// validatePodResourceSpecs validates resources for a pod spec
func (v *MPIJobValidator) validatePodResourceSpecs(podSpec *corev1.PodSpec, replicaType string) error {
	if len(podSpec.Containers) == 0 {
		return fmt.Errorf("%s must have at least one container", replicaType)
	}

	for i, container := range podSpec.Containers {
		// Validate resource requests are specified
		if container.Resources.Requests == nil {
			return fmt.Errorf("%s container[%d] must specify resource requests", replicaType, i)
		}

		// Validate CPU and memory are specified
		if _, ok := container.Resources.Requests[corev1.ResourceCPU]; !ok {
			return fmt.Errorf("%s container[%d] must specify CPU request", replicaType, i)
		}
		if _, ok := container.Resources.Requests[corev1.ResourceMemory]; !ok {
			return fmt.Errorf("%s container[%d] must specify memory request", replicaType, i)
		}

		// Validate limits >= requests if limits are specified
		if container.Resources.Limits != nil {
			if err := v.validateResourceLimits(container.Resources.Requests, container.Resources.Limits, replicaType, i); err != nil {
				return err
			}
		}

		// Validate GPU requests are non-negative
		if gpuReq, ok := container.Resources.Requests["nvidia.com/gpu"]; ok {
			if gpuReq.Sign() < 0 {
				return fmt.Errorf("%s container[%d] GPU request must be >= 0", replicaType, i)
			}
		}
	}

	return nil
}

// validateResourceLimits ensures limits are greater than or equal to requests
func (v *MPIJobValidator) validateResourceLimits(requests, limits corev1.ResourceList, replicaType string, containerIndex int) error {
	for resourceName, requestQty := range requests {
		if limitQty, ok := limits[resourceName]; ok {
			if limitQty.Cmp(requestQty) < 0 {
				return fmt.Errorf("%s container[%d] resource limit for %s (%s) must be >= request (%s)",
					replicaType, containerIndex, resourceName, limitQty.String(), requestQty.String())
			}
		}
	}
	return nil
}

// validateSlotsPerWorker validates the slotsPerWorker field
func (v *MPIJobValidator) validateSlotsPerWorker(mpijob *trainingv2.MPIJob) error {
	if mpijob.Spec.SlotsPerWorker != nil {
		if *mpijob.Spec.SlotsPerWorker < 1 {
			return fmt.Errorf("slotsPerWorker must be >= 1, got %d", *mpijob.Spec.SlotsPerWorker)
		}
	}
	return nil
}

// validateMPIImplementation validates the MPI implementation field
func (v *MPIJobValidator) validateMPIImplementation(mpijob *trainingv2.MPIJob) error {
	if mpijob.Spec.MPIImplementation != nil {
		impl := *mpijob.Spec.MPIImplementation
		validImpls := []string{"OpenMPI", "IntelMPI", "MPICH"}

		valid := false
		for _, validImpl := range validImpls {
			if impl == validImpl {
				valid = true
				break
			}
		}

		if !valid {
			return fmt.Errorf("mpiImplementation must be one of %v, got %s", validImpls, impl)
		}
	}
	return nil
}

// validateRunPolicy validates the run policy configuration
func (v *MPIJobValidator) validateRunPolicy(mpijob *trainingv2.MPIJob) error {
	if mpijob.Spec.RunPolicy == nil {
		return nil
	}

	runPolicy := mpijob.Spec.RunPolicy

	// Validate cleanPodPolicy
	if runPolicy.CleanPodPolicy != nil {
		validPolicies := []string{"All", "Running", "None"}
		policy := string(*runPolicy.CleanPodPolicy)

		valid := false
		for _, validPolicy := range validPolicies {
			if policy == validPolicy {
				valid = true
				break
			}
		}

		if !valid {
			return fmt.Errorf("cleanPodPolicy must be one of %v, got %s", validPolicies, policy)
		}
	}

	// Validate ttlSecondsAfterFinished
	if runPolicy.TTLSecondsAfterFinished != nil && *runPolicy.TTLSecondsAfterFinished < 0 {
		return fmt.Errorf("ttlSecondsAfterFinished must be >= 0, got %d", *runPolicy.TTLSecondsAfterFinished)
	}

	// Validate activeDeadlineSeconds
	if runPolicy.ActiveDeadlineSeconds != nil && *runPolicy.ActiveDeadlineSeconds < 0 {
		return fmt.Errorf("activeDeadlineSeconds must be >= 0, got %d", *runPolicy.ActiveDeadlineSeconds)
	}

	// Validate backoffLimit
	if runPolicy.BackoffLimit != nil && *runPolicy.BackoffLimit < 0 {
		return fmt.Errorf("backoffLimit must be >= 0, got %d", *runPolicy.BackoffLimit)
	}

	return nil
}

// validateSchedulingPolicy validates the scheduling policy configuration
func (v *MPIJobValidator) validateSchedulingPolicy(ctx context.Context, mpijob *trainingv2.MPIJob) error {
	if mpijob.Spec.RunPolicy == nil || mpijob.Spec.RunPolicy.SchedulingPolicy == nil {
		return nil
	}

	schedulingPolicy := mpijob.Spec.RunPolicy.SchedulingPolicy

	// Validate priority class exists if specified
	if schedulingPolicy.PriorityClass != nil && *schedulingPolicy.PriorityClass != "" {
		// TODO: Check if priority class exists in cluster
		// This requires a cluster API call which may not be appropriate in a webhook
		// Consider moving this to a separate admission controller
	}

	// Validate queue for gang scheduling
	if schedulingPolicy.Queue != nil && *schedulingPolicy.Queue != "" {
		// Queue name should be valid (no special characters except hyphens)
		queueName := *schedulingPolicy.Queue
		if !isValidQueueName(queueName) {
			return fmt.Errorf("invalid queue name: %s (must contain only alphanumeric characters and hyphens)", queueName)
		}
	}

	return nil
}

// validateContainerSpecs validates container specifications
func (v *MPIJobValidator) validateContainerSpecs(mpijob *trainingv2.MPIJob) error {
	// Validate launcher container specs
	launcherSpec := mpijob.Spec.MPIReplicaSpecs["Launcher"]
	if err := v.validatePodContainerSpecs(&launcherSpec.Template.Spec, "Launcher"); err != nil {
		return err
	}

	// Validate worker container specs
	workerSpec := mpijob.Spec.MPIReplicaSpecs["Worker"]
	if err := v.validatePodContainerSpecs(&workerSpec.Template.Spec, "Worker"); err != nil {
		return err
	}

	return nil
}

// validatePodContainerSpecs validates containers in a pod spec
func (v *MPIJobValidator) validatePodContainerSpecs(podSpec *corev1.PodSpec, replicaType string) error {
	if len(podSpec.Containers) == 0 {
		return fmt.Errorf("%s must have at least one container", replicaType)
	}

	for i, container := range podSpec.Containers {
		// Validate image is specified
		if container.Image == "" {
			return fmt.Errorf("%s container[%d] must specify an image", replicaType, i)
		}

		// Validate command or args is specified
		if len(container.Command) == 0 && len(container.Args) == 0 {
			return fmt.Errorf("%s container[%d] must specify command or args", replicaType, i)
		}

		// Validate volume mounts reference existing volumes
		if err := v.validateVolumeMounts(podSpec.Volumes, container.VolumeMounts, replicaType, i); err != nil {
			return err
		}
	}

	return nil
}

// validateVolumeMounts ensures all volume mounts reference existing volumes
func (v *MPIJobValidator) validateVolumeMounts(volumes []corev1.Volume, mounts []corev1.VolumeMount, replicaType string, containerIndex int) error {
	volumeMap := make(map[string]bool)
	for _, vol := range volumes {
		volumeMap[vol.Name] = true
	}

	for _, mount := range mounts {
		if !volumeMap[mount.Name] {
			return fmt.Errorf("%s container[%d] references non-existent volume: %s", replicaType, containerIndex, mount.Name)
		}
	}

	return nil
}

// validateResourceQuotas validates resource requests against namespace quotas
func (v *MPIJobValidator) validateResourceQuotas(ctx context.Context, mpijob *trainingv2.MPIJob) error {
	// Calculate total resource requirements
	totalResources := v.calculateTotalResources(mpijob)

	// Get namespace resource quota
	quotaList := &corev1.ResourceQuotaList{}
	err := v.Client.List(ctx, quotaList, client.InNamespace(mpijob.Namespace))
	if err != nil {
		// If we can't get quotas, allow the request to proceed
		// The resource quota admission controller will handle enforcement
		return nil
	}

	// Check if resources exceed quota
	for _, quota := range quotaList.Items {
		if err := v.checkQuotaCompliance(totalResources, &quota); err != nil {
			return err
		}
	}

	return nil
}

// calculateTotalResources calculates total resource requirements for the MPIJob
func (v *MPIJobValidator) calculateTotalResources(mpijob *trainingv2.MPIJob) corev1.ResourceList {
	totalResources := corev1.ResourceList{}

	// Add launcher resources
	launcherSpec := mpijob.Spec.MPIReplicaSpecs["Launcher"]
	addPodResources(totalResources, &launcherSpec.Template.Spec, 1)

	// Add worker resources
	workerSpec := mpijob.Spec.MPIReplicaSpecs["Worker"]
	workerReplicas := int32(1)
	if workerSpec.Replicas != nil {
		workerReplicas = *workerSpec.Replicas
	}
	addPodResources(totalResources, &workerSpec.Template.Spec, workerReplicas)

	return totalResources
}

// addPodResources adds pod resources to the total
func addPodResources(total corev1.ResourceList, podSpec *corev1.PodSpec, replicas int32) {
	for _, container := range podSpec.Containers {
		for resourceName, quantity := range container.Resources.Requests {
			if existing, ok := total[resourceName]; ok {
				quantity.Add(existing)
			}
			scaledQty := quantity.DeepCopy()
			scaledQty.Set(quantity.Value() * int64(replicas))
			total[resourceName] = scaledQty
		}
	}
}

// checkQuotaCompliance checks if resources comply with quota
func (v *MPIJobValidator) checkQuotaCompliance(required corev1.ResourceList, quota *corev1.ResourceQuota) error {
	for resourceName, requiredQty := range required {
		if hardLimit, ok := quota.Status.Hard[resourceName]; ok {
			if usedQty, ok := quota.Status.Used[resourceName]; ok {
				available := hardLimit.DeepCopy()
				available.Sub(usedQty)

				if requiredQty.Cmp(available) > 0 {
					return fmt.Errorf("insufficient quota for %s: requested %s, available %s (quota: %s, used: %s)",
						resourceName, requiredQty.String(), available.String(), hardLimit.String(), usedQty.String())
				}
			}
		}
	}
	return nil
}

// validateSecurityContext validates security context settings
func (v *MPIJobValidator) validateSecurityContext(mpijob *trainingv2.MPIJob) error {
	// Validate launcher security context
	launcherSpec := mpijob.Spec.MPIReplicaSpecs["Launcher"]
	if err := v.validatePodSecurityContext(&launcherSpec.Template.Spec, "Launcher"); err != nil {
		return err
	}

	// Validate worker security context
	workerSpec := mpijob.Spec.MPIReplicaSpecs["Worker"]
	if err := v.validatePodSecurityContext(&workerSpec.Template.Spec, "Worker"); err != nil {
		return err
	}

	return nil
}

// validatePodSecurityContext validates security context for a pod
func (v *MPIJobValidator) validatePodSecurityContext(podSpec *corev1.PodSpec, replicaType string) error {
	// Check pod-level security context
	if podSpec.SecurityContext != nil {
		// Recommend runAsNonRoot
		if podSpec.SecurityContext.RunAsNonRoot != nil && !*podSpec.SecurityContext.RunAsNonRoot {
			// This is a warning, not an error - allow it but log
			// TODO: Add warning mechanism
		}
	}

	// Check container-level security context
	for i, container := range podSpec.Containers {
		if container.SecurityContext != nil {
			// Check allowPrivilegeEscalation
			if container.SecurityContext.AllowPrivilegeEscalation != nil && *container.SecurityContext.AllowPrivilegeEscalation {
				// This is a security concern but may be needed for some MPI implementations
				// TODO: Make this configurable based on cluster policy
			}

			// Check privileged
			if container.SecurityContext.Privileged != nil && *container.SecurityContext.Privileged {
				return fmt.Errorf("%s container[%d] cannot run as privileged", replicaType, i)
			}
		}
	}

	return nil
}

// validateNetworkPolicy validates network policy settings
func (v *MPIJobValidator) validateNetworkPolicy(mpijob *trainingv2.MPIJob) error {
	if mpijob.Spec.NetworkPolicy == nil {
		return nil
	}

	if mpijob.Spec.NetworkPolicy.Template != nil {
		template := *mpijob.Spec.NetworkPolicy.Template
		validTemplates := []string{"Default", "Restricted"}

		valid := false
		for _, validTemplate := range validTemplates {
			if template == validTemplate {
				valid = true
				break
			}
		}

		if !valid {
			return fmt.Errorf("networkPolicy.template must be one of %v, got %s", validTemplates, template)
		}
	}

	return nil
}

// isValidQueueName checks if a queue name is valid
func isValidQueueName(name string) bool {
	if len(name) == 0 || len(name) > 63 {
		return false
	}

	for _, c := range name {
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-') {
			return false
		}
	}

	// Cannot start or end with hyphen
	if strings.HasPrefix(name, "-") || strings.HasSuffix(name, "-") {
		return false
	}

	return true
}
