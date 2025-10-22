package controllers

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"

	trainingv2 "github.com/kubeflow/training-operator/pkg/apis/kubeflow.org/v2"
)

// MPIJobStatusTracker tracks and updates MPIJob status
type MPIJobStatusTracker struct {
	client client.Client
}

// NewMPIJobStatusTracker creates a new status tracker
func NewMPIJobStatusTracker(client client.Client) *MPIJobStatusTracker {
	return &MPIJobStatusTracker{
		client: client,
	}
}

// UpdateJobStatus updates the status of an MPIJob based on pod states
func (st *MPIJobStatusTracker) UpdateJobStatus(ctx context.Context, mpijob *trainingv2.MPIJob, client client.Client) error {
	// Get all pods for this MPIJob
	launcherPods, workerPods, err := st.getPods(ctx, mpijob, client)
	if err != nil {
		return fmt.Errorf("failed to get pods: %w", err)
	}

	// Update replica statuses
	st.updateReplicaStatuses(mpijob, launcherPods, workerPods)

	// Update job conditions based on pod states
	st.updateJobConditions(mpijob, launcherPods, workerPods)

	// Set start time if not set and pods are running
	if mpijob.Status.StartTime == nil && st.hasRunningPods(launcherPods, workerPods) {
		now := metav1.Now()
		mpijob.Status.StartTime = &now
	}

	return nil
}

// getPods retrieves all launcher and worker pods for an MPIJob
func (st *MPIJobStatusTracker) getPods(ctx context.Context, mpijob *trainingv2.MPIJob, client client.Client) ([]corev1.Pod, []corev1.Pod, error) {
	// Get launcher pods
	launcherPods := &corev1.PodList{}
	launcherSelector := labels.SelectorFromSet(map[string]string{
		"mpijob-name": mpijob.Name,
		"mpijob-role": "launcher",
	})
	if err := client.List(ctx, launcherPods, &client.ListOptions{
		Namespace:     mpijob.Namespace,
		LabelSelector: launcherSelector,
	}); err != nil {
		return nil, nil, fmt.Errorf("failed to list launcher pods: %w", err)
	}

	// Get worker pods
	workerPods := &corev1.PodList{}
	workerSelector := labels.SelectorFromSet(map[string]string{
		"mpijob-name": mpijob.Name,
		"mpijob-role": "worker",
	})
	if err := client.List(ctx, workerPods, &client.ListOptions{
		Namespace:     mpijob.Namespace,
		LabelSelector: workerSelector,
	}); err != nil {
		return nil, nil, fmt.Errorf("failed to list worker pods: %w", err)
	}

	return launcherPods.Items, workerPods.Items, nil
}

// updateReplicaStatuses updates the replica status counts
func (st *MPIJobStatusTracker) updateReplicaStatuses(mpijob *trainingv2.MPIJob, launcherPods, workerPods []corev1.Pod) {
	if mpijob.Status.ReplicaStatuses == nil {
		mpijob.Status.ReplicaStatuses = make(map[string]*trainingv2.ReplicaStatus)
	}

	// Update launcher status
	launcherStatus := &trainingv2.ReplicaStatus{
		Active:    0,
		Succeeded: 0,
		Failed:    0,
	}
	for _, pod := range launcherPods {
		switch pod.Status.Phase {
		case corev1.PodRunning, corev1.PodPending:
			launcherStatus.Active++
		case corev1.PodSucceeded:
			launcherStatus.Succeeded++
		case corev1.PodFailed:
			launcherStatus.Failed++
		}
	}
	mpijob.Status.ReplicaStatuses["Launcher"] = launcherStatus

	// Update worker status
	workerStatus := &trainingv2.ReplicaStatus{
		Active:    0,
		Succeeded: 0,
		Failed:    0,
	}
	for _, pod := range workerPods {
		switch pod.Status.Phase {
		case corev1.PodRunning, corev1.PodPending:
			workerStatus.Active++
		case corev1.PodSucceeded:
			workerStatus.Succeeded++
		case corev1.PodFailed:
			workerStatus.Failed++
		}
	}
	mpijob.Status.ReplicaStatuses["Worker"] = workerStatus
}

// updateJobConditions updates job conditions based on pod states
func (st *MPIJobStatusTracker) updateJobConditions(mpijob *trainingv2.MPIJob, launcherPods, workerPods []corev1.Pod) {
	// Check if job has already finished
	if st.hasCondition(mpijob, trainingv2.JobSucceeded) || st.hasCondition(mpijob, trainingv2.JobFailed) {
		return
	}

	// Check for failures
	if st.hasFailedPods(launcherPods, workerPods) {
		reason, message := st.getFailureReason(launcherPods, workerPods)
		st.UpdateCondition(mpijob, trainingv2.JobFailed, reason, message)
		return
	}

	// Check for success (launcher succeeded and all workers succeeded)
	if st.allPodsSucceeded(launcherPods, workerPods) {
		st.UpdateCondition(mpijob, trainingv2.JobSucceeded, "JobSucceeded", "All pods completed successfully")
		return
	}

	// Check if all pods are running
	if st.allPodsRunning(launcherPods, workerPods) {
		if !st.hasCondition(mpijob, trainingv2.JobRunning) {
			st.UpdateCondition(mpijob, trainingv2.JobRunning, "JobRunning", "All pods are running")
		}
		return
	}

	// Check if pods are pending
	if st.hasPendingPods(launcherPods, workerPods) {
		reason, message := st.getPendingReason(launcherPods, workerPods)
		if !st.hasCondition(mpijob, trainingv2.JobCreated) {
			st.UpdateCondition(mpijob, trainingv2.JobCreated, reason, message)
		}
	}
}

// UpdateCondition updates or adds a condition to the MPIJob status
func (st *MPIJobStatusTracker) UpdateCondition(mpijob *trainingv2.MPIJob, conditionType trainingv2.JobConditionType, reason, message string) {
	now := metav1.Now()

	// Check if condition already exists
	for i, condition := range mpijob.Status.Conditions {
		if condition.Type == string(conditionType) {
			// Update existing condition
			mpijob.Status.Conditions[i].Status = metav1.ConditionTrue
			mpijob.Status.Conditions[i].LastTransitionTime = now
			mpijob.Status.Conditions[i].LastUpdateTime = now
			mpijob.Status.Conditions[i].Reason = reason
			mpijob.Status.Conditions[i].Message = message
			return
		}
	}

	// Add new condition
	newCondition := metav1.Condition{
		Type:               string(conditionType),
		Status:             metav1.ConditionTrue,
		LastTransitionTime: now,
		LastUpdateTime:     now,
		Reason:             reason,
		Message:            message,
	}
	mpijob.Status.Conditions = append(mpijob.Status.Conditions, newCondition)
}

// hasCondition checks if a condition exists with status True
func (st *MPIJobStatusTracker) hasCondition(mpijob *trainingv2.MPIJob, conditionType trainingv2.JobConditionType) bool {
	for _, condition := range mpijob.Status.Conditions {
		if condition.Type == string(conditionType) && condition.Status == metav1.ConditionTrue {
			return true
		}
	}
	return false
}

// hasFailedPods checks if any pods have failed
func (st *MPIJobStatusTracker) hasFailedPods(launcherPods, workerPods []corev1.Pod) bool {
	// Check launcher pods
	for _, pod := range launcherPods {
		if pod.Status.Phase == corev1.PodFailed {
			return true
		}
		// Check container statuses for failures
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.State.Terminated != nil && containerStatus.State.Terminated.ExitCode != 0 {
				return true
			}
		}
	}

	// Check worker pods
	for _, pod := range workerPods {
		if pod.Status.Phase == corev1.PodFailed {
			return true
		}
		// Check container statuses for failures
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.State.Terminated != nil && containerStatus.State.Terminated.ExitCode != 0 {
				return true
			}
		}
	}

	return false
}

// allPodsSucceeded checks if all pods have succeeded
func (st *MPIJobStatusTracker) allPodsSucceeded(launcherPods, workerPods []corev1.Pod) bool {
	// Need at least one launcher and one worker
	if len(launcherPods) == 0 || len(workerPods) == 0 {
		return false
	}

	// Check all launcher pods succeeded
	for _, pod := range launcherPods {
		if pod.Status.Phase != corev1.PodSucceeded {
			return false
		}
	}

	// Check all worker pods succeeded
	for _, pod := range workerPods {
		if pod.Status.Phase != corev1.PodSucceeded {
			return false
		}
	}

	return true
}

// allPodsRunning checks if all pods are running
func (st *MPIJobStatusTracker) allPodsRunning(launcherPods, workerPods []corev1.Pod) bool {
	// Need at least one launcher and one worker
	if len(launcherPods) == 0 || len(workerPods) == 0 {
		return false
	}

	// Check all launcher pods running
	for _, pod := range launcherPods {
		if pod.Status.Phase != corev1.PodRunning {
			return false
		}
	}

	// Check all worker pods running
	for _, pod := range workerPods {
		if pod.Status.Phase != corev1.PodRunning {
			return false
		}
	}

	return true
}

// hasPendingPods checks if any pods are pending
func (st *MPIJobStatusTracker) hasPendingPods(launcherPods, workerPods []corev1.Pod) bool {
	for _, pod := range launcherPods {
		if pod.Status.Phase == corev1.PodPending {
			return true
		}
	}
	for _, pod := range workerPods {
		if pod.Status.Phase == corev1.PodPending {
			return true
		}
	}
	return false
}

// hasRunningPods checks if any pods are running
func (st *MPIJobStatusTracker) hasRunningPods(launcherPods, workerPods []corev1.Pod) bool {
	for _, pod := range launcherPods {
		if pod.Status.Phase == corev1.PodRunning {
			return true
		}
	}
	for _, pod := range workerPods {
		if pod.Status.Phase == corev1.PodRunning {
			return true
		}
	}
	return false
}

// getFailureReason determines the reason and message for job failure
func (st *MPIJobStatusTracker) getFailureReason(launcherPods, workerPods []corev1.Pod) (string, string) {
	// Check launcher pods first
	for _, pod := range launcherPods {
		if pod.Status.Phase == corev1.PodFailed {
			reason := "LauncherFailed"
			message := fmt.Sprintf("Launcher pod %s failed", pod.Name)
			if pod.Status.Reason != "" {
				message += fmt.Sprintf(": %s", pod.Status.Reason)
			}
			if pod.Status.Message != "" {
				message += fmt.Sprintf(" - %s", pod.Status.Message)
			}
			return reason, message
		}

		// Check container statuses
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.State.Terminated != nil && containerStatus.State.Terminated.ExitCode != 0 {
				reason := "LauncherContainerFailed"
				message := fmt.Sprintf("Launcher container %s in pod %s failed with exit code %d",
					containerStatus.Name, pod.Name, containerStatus.State.Terminated.ExitCode)
				if containerStatus.State.Terminated.Reason != "" {
					message += fmt.Sprintf(": %s", containerStatus.State.Terminated.Reason)
				}
				if containerStatus.State.Terminated.Message != "" {
					message += fmt.Sprintf(" - %s", containerStatus.State.Terminated.Message)
				}
				return reason, message
			}
		}
	}

	// Check worker pods
	failedWorkers := []string{}
	for _, pod := range workerPods {
		if pod.Status.Phase == corev1.PodFailed {
			failedWorkers = append(failedWorkers, pod.Name)
		}

		// Check container statuses
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.State.Terminated != nil && containerStatus.State.Terminated.ExitCode != 0 {
				reason := "WorkerContainerFailed"
				message := fmt.Sprintf("Worker container %s in pod %s failed with exit code %d",
					containerStatus.Name, pod.Name, containerStatus.State.Terminated.ExitCode)
				if containerStatus.State.Terminated.Reason != "" {
					message += fmt.Sprintf(": %s", containerStatus.State.Terminated.Reason)
				}
				if containerStatus.State.Terminated.Message != "" {
					message += fmt.Sprintf(" - %s", containerStatus.State.Terminated.Message)
				}
				return reason, message
			}
		}
	}

	if len(failedWorkers) > 0 {
		reason := "WorkersFailed"
		message := fmt.Sprintf("%d worker pod(s) failed: %v", len(failedWorkers), failedWorkers)
		return reason, message
	}

	return "JobFailed", "Unknown failure reason"
}

// getPendingReason determines the reason for pods being pending
func (st *MPIJobStatusTracker) getPendingReason(launcherPods, workerPods []corev1.Pod) (string, string) {
	allPods := append(launcherPods, workerPods...)

	pendingReasons := make(map[string]int)
	for _, pod := range allPods {
		if pod.Status.Phase == corev1.PodPending {
			// Check pod conditions for specific reasons
			for _, condition := range pod.Status.Conditions {
				if condition.Status == corev1.ConditionFalse {
					pendingReasons[condition.Reason]++
				}
			}

			// Check container statuses
			for _, containerStatus := range pod.Status.ContainerStatuses {
				if containerStatus.State.Waiting != nil {
					pendingReasons[containerStatus.State.Waiting.Reason]++
				}
			}
		}
	}

	// Determine most common reason
	mostCommonReason := ""
	maxCount := 0
	for reason, count := range pendingReasons {
		if count > maxCount {
			mostCommonReason = reason
			maxCount = count
		}
	}

	if mostCommonReason == "" {
		return "PodsPending", "Waiting for pods to be scheduled"
	}

	message := fmt.Sprintf("Pods pending: %s (%d pods)", mostCommonReason, maxCount)
	return "PodsPending", message
}

// GetJobDuration returns the duration of the job
func (st *MPIJobStatusTracker) GetJobDuration(mpijob *trainingv2.MPIJob) time.Duration {
	if mpijob.Status.StartTime == nil {
		return 0
	}

	endTime := time.Now()
	if mpijob.Status.CompletionTime != nil {
		endTime = mpijob.Status.CompletionTime.Time
	}

	return endTime.Sub(mpijob.Status.StartTime.Time)
}

// GetJobPhase returns the current phase of the job
func (st *MPIJobStatusTracker) GetJobPhase(mpijob *trainingv2.MPIJob) string {
	if len(mpijob.Status.Conditions) == 0 {
		return "Created"
	}

	// Find the latest condition with status True
	for i := len(mpijob.Status.Conditions) - 1; i >= 0; i-- {
		if mpijob.Status.Conditions[i].Status == metav1.ConditionTrue {
			return mpijob.Status.Conditions[i].Type
		}
	}

	return "Unknown"
}

// IsJobFinished checks if the job has finished
func (st *MPIJobStatusTracker) IsJobFinished(mpijob *trainingv2.MPIJob) bool {
	return st.hasCondition(mpijob, trainingv2.JobSucceeded) || st.hasCondition(mpijob, trainingv2.JobFailed)
}

// GetWorkerStatsSummary returns a human-readable summary of worker status
func (st *MPIJobStatusTracker) GetWorkerStatsSummary(mpijob *trainingv2.MPIJob) string {
	workerStatus, ok := mpijob.Status.ReplicaStatuses["Worker"]
	if !ok {
		return "No workers"
	}

	total := workerStatus.Active + workerStatus.Succeeded + workerStatus.Failed
	if workerStatus.Failed > 0 {
		return fmt.Sprintf("%d/%d running, %d failed", workerStatus.Active, total, workerStatus.Failed)
	}
	if workerStatus.Succeeded > 0 {
		return fmt.Sprintf("%d/%d succeeded", workerStatus.Succeeded, total)
	}
	return fmt.Sprintf("%d/%d running", workerStatus.Active, total)
}
