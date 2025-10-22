package training

import (
	"context"
	"fmt"
	"io"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"

	trainingv2 "github.com/kubeflow/training-operator/pkg/apis/kubeflow.org/v2"
)

// MPIJobServiceServer implements the gRPC MPIJobService
type MPIJobServiceServer struct {
	UnimplementedMPIJobServiceServer
	k8sClient   client.Client
	kubeClient  kubernetes.Interface
}

// NewMPIJobServiceServer creates a new gRPC service server
func NewMPIJobServiceServer(k8sClient client.Client, kubeClient kubernetes.Interface) *MPIJobServiceServer {
	return &MPIJobServiceServer{
		k8sClient:  k8sClient,
		kubeClient: kubeClient,
	}
}

// CreateMPIJob creates a new MPIJob
func (s *MPIJobServiceServer) CreateMPIJob(ctx context.Context, req *CreateMPIJobRequest) (*CreateMPIJobResponse, error) {
	if req.Namespace == "" {
		return nil, status.Error(codes.InvalidArgument, "namespace is required")
	}
	if req.Spec == nil {
		return nil, status.Error(codes.InvalidArgument, "spec is required")
	}

	// Convert proto spec to CRD
	mpijob := &trainingv2.MPIJob{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "mpijob-",
			Namespace:    req.Namespace,
		},
		Spec: s.convertProtoSpecToCRD(req.Spec),
	}

	// Create the MPIJob
	if err := s.k8sClient.Create(ctx, mpijob); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create MPIJob: %v", err)
	}

	return &CreateMPIJobResponse{
		Name:      mpijob.Name,
		Namespace: mpijob.Namespace,
		Uid:       string(mpijob.UID),
		CreatedAt: timestamppb.New(mpijob.CreationTimestamp.Time),
	}, nil
}

// GetMPIJob retrieves an MPIJob by name
func (s *MPIJobServiceServer) GetMPIJob(ctx context.Context, req *GetMPIJobRequest) (*GetMPIJobResponse, error) {
	if req.Namespace == "" {
		return nil, status.Error(codes.InvalidArgument, "namespace is required")
	}
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	mpijob := &trainingv2.MPIJob{}
	if err := s.k8sClient.Get(ctx, types.NamespacedName{
		Name:      req.Name,
		Namespace: req.Namespace,
	}, mpijob); err != nil {
		if errors.IsNotFound(err) {
			return nil, status.Errorf(codes.NotFound, "MPIJob not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to get MPIJob: %v", err)
	}

	return &GetMPIJobResponse{
		Mpijob: s.convertCRDToProto(mpijob),
	}, nil
}

// ListMPIJobs lists MPIJobs in a namespace
func (s *MPIJobServiceServer) ListMPIJobs(ctx context.Context, req *ListMPIJobsRequest) (*ListMPIJobsResponse, error) {
	if req.Namespace == "" {
		return nil, status.Error(codes.InvalidArgument, "namespace is required")
	}

	mpijobList := &trainingv2.MPIJobList{}
	listOpts := &client.ListOptions{
		Namespace: req.Namespace,
	}

	if req.Limit > 0 {
		listOpts.Limit = int64(req.Limit)
	}
	if req.ContinueToken != "" {
		listOpts.Continue = req.ContinueToken
	}

	if err := s.k8sClient.List(ctx, mpijobList, listOpts); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list MPIJobs: %v", err)
	}

	items := make([]*MPIJob, len(mpijobList.Items))
	for i, mpijob := range mpijobList.Items {
		items[i] = s.convertCRDToProto(&mpijob)
	}

	return &ListMPIJobsResponse{
		Items:         items,
		ContinueToken: mpijobList.Continue,
	}, nil
}

// DeleteMPIJob deletes an MPIJob
func (s *MPIJobServiceServer) DeleteMPIJob(ctx context.Context, req *DeleteMPIJobRequest) (*DeleteMPIJobResponse, error) {
	if req.Namespace == "" {
		return nil, status.Error(codes.InvalidArgument, "namespace is required")
	}
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	mpijob := &trainingv2.MPIJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
		},
	}

	if err := s.k8sClient.Delete(ctx, mpijob); err != nil {
		if errors.IsNotFound(err) {
			return nil, status.Errorf(codes.NotFound, "MPIJob not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to delete MPIJob: %v", err)
	}

	return &DeleteMPIJobResponse{
		Message: fmt.Sprintf("MPIJob %s deleted successfully", req.Name),
	}, nil
}

// GetMPIJobStatus retrieves the status of an MPIJob
func (s *MPIJobServiceServer) GetMPIJobStatus(ctx context.Context, req *GetMPIJobStatusRequest) (*GetMPIJobStatusResponse, error) {
	if req.Namespace == "" {
		return nil, status.Error(codes.InvalidArgument, "namespace is required")
	}
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	mpijob := &trainingv2.MPIJob{}
	if err := s.k8sClient.Get(ctx, types.NamespacedName{
		Name:      req.Name,
		Namespace: req.Namespace,
	}, mpijob); err != nil {
		if errors.IsNotFound(err) {
			return nil, status.Errorf(codes.NotFound, "MPIJob not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to get MPIJob: %v", err)
	}

	return &GetMPIJobStatusResponse{
		Status: s.convertStatusToProto(&mpijob.Status),
	}, nil
}

// WatchMPIJobStatus watches for status changes of an MPIJob
func (s *MPIJobServiceServer) WatchMPIJobStatus(req *WatchMPIJobStatusRequest, stream MPIJobService_WatchMPIJobStatusServer) error {
	if req.Namespace == "" {
		return status.Error(codes.InvalidArgument, "namespace is required")
	}
	if req.Name == "" {
		return status.Error(codes.InvalidArgument, "name is required")
	}

	// Create a watch on the MPIJob
	watcher, err := s.k8sClient.Watch(stream.Context(), &trainingv2.MPIJobList{}, &client.ListOptions{
		Namespace: req.Namespace,
	})
	if err != nil {
		return status.Errorf(codes.Internal, "failed to watch MPIJob: %v", err)
	}
	defer watcher.Stop()

	// Stream events
	for {
		select {
		case event, ok := <-watcher.ResultChan():
			if !ok {
				return nil
			}

			mpijob, ok := event.Object.(*trainingv2.MPIJob)
			if !ok {
				continue
			}

			// Filter by name
			if mpijob.Name != req.Name {
				continue
			}

			eventType := string(event.Type)
			statusEvent := &MPIJobStatusEvent{
				EventType: eventType,
				Status:    s.convertStatusToProto(&mpijob.Status),
				Timestamp: timestamppb.Now(),
			}

			if err := stream.Send(statusEvent); err != nil {
				return status.Errorf(codes.Internal, "failed to send event: %v", err)
			}

			// Stop watching if job is deleted
			if event.Type == watch.Deleted {
				return nil
			}

		case <-stream.Context().Done():
			return stream.Context().Err()
		}
	}
}

// GetMPIJobLogs retrieves logs from an MPIJob
func (s *MPIJobServiceServer) GetMPIJobLogs(req *GetMPIJobLogsRequest, stream MPIJobService_GetMPIJobLogsServer) error {
	if req.Namespace == "" {
		return status.Error(codes.InvalidArgument, "namespace is required")
	}
	if req.Name == "" {
		return status.Error(codes.InvalidArgument, "name is required")
	}

	// Determine pod name based on type
	var podName string
	if req.PodType == "launcher" {
		podName = fmt.Sprintf("%s-launcher", req.Name)
	} else if req.PodType == "worker" {
		podName = fmt.Sprintf("%s-worker-%d", req.Name, req.PodIndex)
	} else {
		return status.Error(codes.InvalidArgument, "pod_type must be 'launcher' or 'worker'")
	}

	// Build log options
	logOpts := &corev1.PodLogOptions{
		Follow: req.Follow,
	}
	if req.TailLines > 0 {
		tailLines := int64(req.TailLines)
		logOpts.TailLines = &tailLines
	}
	if req.SinceTime != nil {
		logOpts.SinceTime = &metav1.Time{Time: req.SinceTime.AsTime()}
	}

	// Get pod logs
	logStream, err := s.kubeClient.CoreV1().Pods(req.Namespace).GetLogs(podName, logOpts).Stream(stream.Context())
	if err != nil {
		return status.Errorf(codes.Internal, "failed to get logs: %v", err)
	}
	defer logStream.Close()

	// Stream logs
	buf := make([]byte, 4096)
	for {
		n, err := logStream.Read(buf)
		if n > 0 {
			logMsg := &LogMessage{
				PodName:       podName,
				ContainerName: "training",
				Timestamp:     timestamppb.Now(),
				Message:       string(buf[:n]),
			}
			if err := stream.Send(logMsg); err != nil {
				return status.Errorf(codes.Internal, "failed to send log message: %v", err)
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return status.Errorf(codes.Internal, "error reading logs: %v", err)
		}
	}

	return nil
}

// Helper methods for converting between proto and CRD types

func (s *MPIJobServiceServer) convertProtoSpecToCRD(protoSpec *MPIJobSpec) trainingv2.MPIJobSpec {
	spec := trainingv2.MPIJobSpec{
		MPIReplicaSpecs: make(map[string]*trainingv2.ReplicaSpec),
	}

	if protoSpec.SlotsPerWorker > 0 {
		slots := protoSpec.SlotsPerWorker
		spec.SlotsPerWorker = &slots
	}

	if protoSpec.MpiImplementation != "" {
		spec.MPIImplementation = &protoSpec.MpiImplementation
	}

	// Convert replica specs
	for replicaType, protoReplicaSpec := range protoSpec.MpiReplicaSpecs {
		spec.MPIReplicaSpecs[replicaType] = s.convertProtoReplicaSpecToCRD(protoReplicaSpec)
	}

	// Convert run policy
	if protoSpec.RunPolicy != nil {
		spec.RunPolicy = s.convertProtoRunPolicyToCRD(protoSpec.RunPolicy)
	}

	// Convert network policy
	if protoSpec.NetworkPolicy != nil {
		spec.NetworkPolicy = &trainingv2.NetworkPolicy{
			Template: &protoSpec.NetworkPolicy.Template,
		}
	}

	return spec
}

func (s *MPIJobServiceServer) convertProtoReplicaSpecToCRD(protoSpec *ReplicaSpec) *trainingv2.ReplicaSpec {
	spec := &trainingv2.ReplicaSpec{
		Replicas: &protoSpec.Replicas,
	}

	if protoSpec.RestartPolicy != "" {
		spec.RestartPolicy = (*trainingv2.RestartPolicy)(&protoSpec.RestartPolicy)
	}

	// Convert pod template (simplified for now)
	spec.Template = &corev1.PodTemplateSpec{
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{},
		},
	}

	if protoSpec.Template != nil && protoSpec.Template.Spec != nil {
		for _, protoContainer := range protoSpec.Template.Spec.Containers {
			container := corev1.Container{
				Name:    protoContainer.Name,
				Image:   protoContainer.Image,
				Command: protoContainer.Command,
				Args:    protoContainer.Args,
			}

			// Convert environment variables
			for _, protoEnv := range protoContainer.Env {
				container.Env = append(container.Env, corev1.EnvVar{
					Name:  protoEnv.Name,
					Value: protoEnv.Value,
				})
			}

			// Convert resources (simplified)
			if protoContainer.Resources != nil {
				container.Resources = corev1.ResourceRequirements{
					Requests: make(corev1.ResourceList),
					Limits:   make(corev1.ResourceList),
				}
			}

			spec.Template.Spec.Containers = append(spec.Template.Spec.Containers, container)
		}
	}

	return spec
}

func (s *MPIJobServiceServer) convertProtoRunPolicyToCRD(protoPolicy *RunPolicy) *trainingv2.RunPolicy {
	policy := &trainingv2.RunPolicy{}

	if protoPolicy.CleanPodPolicy != "" {
		cleanPolicy := trainingv2.CleanPodPolicy(protoPolicy.CleanPodPolicy)
		policy.CleanPodPolicy = &cleanPolicy
	}

	if protoPolicy.TtlSecondsAfterFinished > 0 {
		policy.TTLSecondsAfterFinished = &protoPolicy.TtlSecondsAfterFinished
	}

	if protoPolicy.ActiveDeadlineSeconds > 0 {
		policy.ActiveDeadlineSeconds = &protoPolicy.ActiveDeadlineSeconds
	}

	if protoPolicy.BackoffLimit > 0 {
		policy.BackoffLimit = &protoPolicy.BackoffLimit
	}

	if protoPolicy.SchedulingPolicy != nil {
		policy.SchedulingPolicy = &trainingv2.SchedulingPolicy{
			PriorityClass: &protoPolicy.SchedulingPolicy.PriorityClass,
			Queue:         &protoPolicy.SchedulingPolicy.Queue,
		}
	}

	return policy
}

func (s *MPIJobServiceServer) convertCRDToProto(mpijob *trainingv2.MPIJob) *MPIJob {
	proto := &MPIJob{
		Name:      mpijob.Name,
		Namespace: mpijob.Namespace,
		Uid:       string(mpijob.UID),
		Labels:    mpijob.Labels,
		CreatedAt: timestamppb.New(mpijob.CreationTimestamp.Time),
		Spec:      s.convertSpecToProto(&mpijob.Spec),
		Status:    s.convertStatusToProto(&mpijob.Status),
	}

	return proto
}

func (s *MPIJobServiceServer) convertSpecToProto(spec *trainingv2.MPIJobSpec) *MPIJobSpec {
	protoSpec := &MPIJobSpec{
		MpiReplicaSpecs: make(map[string]*ReplicaSpec),
	}

	if spec.SlotsPerWorker != nil {
		protoSpec.SlotsPerWorker = *spec.SlotsPerWorker
	}

	if spec.MPIImplementation != nil {
		protoSpec.MpiImplementation = *spec.MPIImplementation
	}

	// Convert replica specs (simplified)
	for replicaType, replicaSpec := range spec.MPIReplicaSpecs {
		protoReplicaSpec := &ReplicaSpec{}
		if replicaSpec.Replicas != nil {
			protoReplicaSpec.Replicas = *replicaSpec.Replicas
		}
		protoSpec.MpiReplicaSpecs[replicaType] = protoReplicaSpec
	}

	return protoSpec
}

func (s *MPIJobServiceServer) convertStatusToProto(status *trainingv2.JobStatus) *MPIJobStatus {
	protoStatus := &MPIJobStatus{
		Conditions:       make([]*Condition, len(status.Conditions)),
		ReplicaStatuses:  make(map[string]*ReplicaStatus),
	}

	// Convert conditions
	for i, cond := range status.Conditions {
		protoStatus.Conditions[i] = &Condition{
			Type:               cond.Type,
			Status:             string(cond.Status),
			Reason:             cond.Reason,
			Message:            cond.Message,
			LastTransitionTime: timestamppb.New(cond.LastTransitionTime.Time),
			LastUpdateTime:     timestamppb.New(cond.LastUpdateTime.Time),
		}
	}

	// Convert replica statuses
	for replicaType, replicaStatus := range status.ReplicaStatuses {
		protoStatus.ReplicaStatuses[replicaType] = &ReplicaStatus{
			Active:    replicaStatus.Active,
			Succeeded: replicaStatus.Succeeded,
			Failed:    replicaStatus.Failed,
		}
	}

	// Convert timestamps
	if status.StartTime != nil {
		protoStatus.StartTime = timestamppb.New(status.StartTime.Time)
	}
	if status.CompletionTime != nil {
		protoStatus.CompletionTime = timestamppb.New(status.CompletionTime.Time)
	}

	return protoStatus
}

// UnimplementedMPIJobServiceServer is a placeholder for future implementation
type UnimplementedMPIJobServiceServer struct{}

func (UnimplementedMPIJobServiceServer) CreateMPIJob(context.Context, *CreateMPIJobRequest) (*CreateMPIJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMPIJob not implemented")
}

func (UnimplementedMPIJobServiceServer) GetMPIJob(context.Context, *GetMPIJobRequest) (*GetMPIJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMPIJob not implemented")
}

func (UnimplementedMPIJobServiceServer) ListMPIJobs(context.Context, *ListMPIJobsRequest) (*ListMPIJobsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListMPIJobs not implemented")
}

func (UnimplementedMPIJobServiceServer) DeleteMPIJob(context.Context, *DeleteMPIJobRequest) (*DeleteMPIJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMPIJob not implemented")
}

func (UnimplementedMPIJobServiceServer) GetMPIJobStatus(context.Context, *GetMPIJobStatusRequest) (*GetMPIJobStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMPIJobStatus not implemented")
}

func (UnimplementedMPIJobServiceServer) WatchMPIJobStatus(*WatchMPIJobStatusRequest, MPIJobService_WatchMPIJobStatusServer) error {
	return status.Errorf(codes.Unimplemented, "method WatchMPIJobStatus not implemented")
}

func (UnimplementedMPIJobServiceServer) GetMPIJobLogs(*GetMPIJobLogsRequest, MPIJobService_GetMPIJobLogsServer) error {
	return status.Errorf(codes.Unimplemented, "method GetMPIJobLogs not implemented")
}
