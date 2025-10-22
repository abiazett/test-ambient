package training

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	trainingv2 "github.com/kubeflow/training-operator/pkg/apis/kubeflow.org/v2"
)

// MPIJobAPI provides REST API endpoints for MPIJob operations
type MPIJobAPI struct {
	client  client.Client
	service *MPIJobService
}

// NewMPIJobAPI creates a new MPIJob API handler
func NewMPIJobAPI(client client.Client, service *MPIJobService) *MPIJobAPI {
	return &MPIJobAPI{
		client:  client,
		service: service,
	}
}

// RegisterRoutes registers all MPIJob API routes
func (api *MPIJobAPI) RegisterRoutes(router *mux.Router) {
	// MPIJob CRUD operations
	router.HandleFunc("/api/v1/namespaces/{namespace}/mpijobs", api.CreateMPIJob).Methods("POST")
	router.HandleFunc("/api/v1/namespaces/{namespace}/mpijobs", api.ListMPIJobs).Methods("GET")
	router.HandleFunc("/api/v1/namespaces/{namespace}/mpijobs/{name}", api.GetMPIJob).Methods("GET")
	router.HandleFunc("/api/v1/namespaces/{namespace}/mpijobs/{name}", api.DeleteMPIJob).Methods("DELETE")
	router.HandleFunc("/api/v1/namespaces/{namespace}/mpijobs/{name}/status", api.GetMPIJobStatus).Methods("GET")

	// MPIJob logs
	router.HandleFunc("/api/v1/namespaces/{namespace}/mpijobs/{name}/logs", api.GetMPIJobLogs).Methods("GET")
	router.HandleFunc("/api/v1/namespaces/{namespace}/mpijobs/{name}/logs/launcher", api.GetLauncherLogs).Methods("GET")
	router.HandleFunc("/api/v1/namespaces/{namespace}/mpijobs/{name}/logs/worker/{index}", api.GetWorkerLogs).Methods("GET")

	// MPIJob operations
	router.HandleFunc("/api/v1/namespaces/{namespace}/mpijobs/{name}/restart", api.RestartMPIJob).Methods("POST")
}

// CreateMPIJob creates a new MPIJob
// POST /api/v1/namespaces/{namespace}/mpijobs
func (api *MPIJobAPI) CreateMPIJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]

	// Parse request body
	var mpijob trainingv2.MPIJob
	if err := json.NewDecoder(r.Body).Decode(&mpijob); err != nil {
		api.sendError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}

	// Set namespace from URL
	mpijob.Namespace = namespace

	// Validate MPIJob
	if err := api.service.ValidateMPIJob(&mpijob); err != nil {
		api.sendError(w, http.StatusBadRequest, fmt.Errorf("validation failed: %w", err))
		return
	}

	// Create MPIJob
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	if err := api.client.Create(ctx, &mpijob); err != nil {
		api.sendError(w, http.StatusInternalServerError, fmt.Errorf("failed to create MPIJob: %w", err))
		return
	}

	api.sendJSON(w, http.StatusCreated, mpijob)
}

// ListMPIJobs lists all MPIJobs in a namespace
// GET /api/v1/namespaces/{namespace}/mpijobs
func (api *MPIJobAPI) ListMPIJobs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]

	// Parse query parameters
	labelSelector := r.URL.Query().Get("labelSelector")
	fieldSelector := r.URL.Query().Get("fieldSelector")
	limit := api.parseIntQuery(r, "limit", 100)
	continueToken := r.URL.Query().Get("continue")

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// List MPIJobs
	mpijobList := &trainingv2.MPIJobList{}
	listOpts := &client.ListOptions{
		Namespace: namespace,
		Limit:     int64(limit),
		Continue:  continueToken,
	}

	// Apply label selector if provided
	if labelSelector != "" {
		selector, err := metav1.ParseToLabelSelector(labelSelector)
		if err != nil {
			api.sendError(w, http.StatusBadRequest, fmt.Errorf("invalid label selector: %w", err))
			return
		}
		listOpts.LabelSelector, _ = metav1.LabelSelectorAsSelector(selector)
	}

	if err := api.client.List(ctx, mpijobList, listOpts); err != nil {
		api.sendError(w, http.StatusInternalServerError, fmt.Errorf("failed to list MPIJobs: %w", err))
		return
	}

	api.sendJSON(w, http.StatusOK, mpijobList)
}

// GetMPIJob retrieves a specific MPIJob
// GET /api/v1/namespaces/{namespace}/mpijobs/{name}
func (api *MPIJobAPI) GetMPIJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	name := vars["name"]

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	mpijob := &trainingv2.MPIJob{}
	if err := api.client.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, mpijob); err != nil {
		if client.IgnoreNotFound(err) == nil {
			api.sendError(w, http.StatusNotFound, fmt.Errorf("MPIJob not found"))
		} else {
			api.sendError(w, http.StatusInternalServerError, fmt.Errorf("failed to get MPIJob: %w", err))
		}
		return
	}

	api.sendJSON(w, http.StatusOK, mpijob)
}

// DeleteMPIJob deletes an MPIJob
// DELETE /api/v1/namespaces/{namespace}/mpijobs/{name}
func (api *MPIJobAPI) DeleteMPIJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	name := vars["name"]

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	mpijob := &trainingv2.MPIJob{}
	if err := api.client.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, mpijob); err != nil {
		if client.IgnoreNotFound(err) == nil {
			api.sendError(w, http.StatusNotFound, fmt.Errorf("MPIJob not found"))
		} else {
			api.sendError(w, http.StatusInternalServerError, fmt.Errorf("failed to get MPIJob: %w", err))
		}
		return
	}

	if err := api.client.Delete(ctx, mpijob); err != nil {
		api.sendError(w, http.StatusInternalServerError, fmt.Errorf("failed to delete MPIJob: %w", err))
		return
	}

	api.sendJSON(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("MPIJob %s deleted successfully", name),
	})
}

// GetMPIJobStatus retrieves the status of an MPIJob
// GET /api/v1/namespaces/{namespace}/mpijobs/{name}/status
func (api *MPIJobAPI) GetMPIJobStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	name := vars["name"]

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	status, err := api.service.GetJobStatus(ctx, namespace, name)
	if err != nil {
		if client.IgnoreNotFound(err) == nil {
			api.sendError(w, http.StatusNotFound, fmt.Errorf("MPIJob not found"))
		} else {
			api.sendError(w, http.StatusInternalServerError, fmt.Errorf("failed to get status: %w", err))
		}
		return
	}

	api.sendJSON(w, http.StatusOK, status)
}

// GetMPIJobLogs retrieves logs from all pods of an MPIJob
// GET /api/v1/namespaces/{namespace}/mpijobs/{name}/logs
func (api *MPIJobAPI) GetMPIJobLogs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	name := vars["name"]

	// Parse query parameters
	follow := r.URL.Query().Get("follow") == "true"
	tailLines := api.parseIntQuery(r, "tailLines", 100)
	since := r.URL.Query().Get("since")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Minute)
	defer cancel()

	logs, err := api.service.GetJobLogs(ctx, namespace, name, &LogOptions{
		Follow:    follow,
		TailLines: tailLines,
		Since:     since,
	})
	if err != nil {
		api.sendError(w, http.StatusInternalServerError, fmt.Errorf("failed to get logs: %w", err))
		return
	}

	api.sendJSON(w, http.StatusOK, logs)
}

// GetLauncherLogs retrieves logs from the launcher pod
// GET /api/v1/namespaces/{namespace}/mpijobs/{name}/logs/launcher
func (api *MPIJobAPI) GetLauncherLogs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	name := vars["name"]

	// Parse query parameters
	follow := r.URL.Query().Get("follow") == "true"
	tailLines := api.parseIntQuery(r, "tailLines", 100)
	since := r.URL.Query().Get("since")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Minute)
	defer cancel()

	logs, err := api.service.GetLauncherLogs(ctx, namespace, name, &LogOptions{
		Follow:    follow,
		TailLines: tailLines,
		Since:     since,
	})
	if err != nil {
		api.sendError(w, http.StatusInternalServerError, fmt.Errorf("failed to get launcher logs: %w", err))
		return
	}

	// If follow mode, stream logs
	if follow {
		api.streamLogs(w, logs)
		return
	}

	api.sendText(w, http.StatusOK, logs)
}

// GetWorkerLogs retrieves logs from a specific worker pod
// GET /api/v1/namespaces/{namespace}/mpijobs/{name}/logs/worker/{index}
func (api *MPIJobAPI) GetWorkerLogs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	name := vars["name"]
	indexStr := vars["index"]

	index, err := strconv.Atoi(indexStr)
	if err != nil {
		api.sendError(w, http.StatusBadRequest, fmt.Errorf("invalid worker index: %s", indexStr))
		return
	}

	// Parse query parameters
	follow := r.URL.Query().Get("follow") == "true"
	tailLines := api.parseIntQuery(r, "tailLines", 100)
	since := r.URL.Query().Get("since")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Minute)
	defer cancel()

	logs, err := api.service.GetWorkerLogs(ctx, namespace, name, index, &LogOptions{
		Follow:    follow,
		TailLines: tailLines,
		Since:     since,
	})
	if err != nil {
		api.sendError(w, http.StatusInternalServerError, fmt.Errorf("failed to get worker logs: %w", err))
		return
	}

	// If follow mode, stream logs
	if follow {
		api.streamLogs(w, logs)
		return
	}

	api.sendText(w, http.StatusOK, logs)
}

// RestartMPIJob restarts a failed MPIJob
// POST /api/v1/namespaces/{namespace}/mpijobs/{name}/restart
func (api *MPIJobAPI) RestartMPIJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	name := vars["name"]

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// Get the MPIJob
	mpijob := &trainingv2.MPIJob{}
	if err := api.client.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, mpijob); err != nil {
		if client.IgnoreNotFound(err) == nil {
			api.sendError(w, http.StatusNotFound, fmt.Errorf("MPIJob not found"))
		} else {
			api.sendError(w, http.StatusInternalServerError, fmt.Errorf("failed to get MPIJob: %w", err))
		}
		return
	}

	// Check if job can be restarted
	if !api.service.CanRestart(mpijob) {
		api.sendError(w, http.StatusBadRequest, fmt.Errorf("MPIJob cannot be restarted (must be in Failed state)"))
		return
	}

	// Delete and recreate the job
	newName := fmt.Sprintf("%s-retry-%d", name, time.Now().Unix())
	newMPIJob := mpijob.DeepCopy()
	newMPIJob.Name = newName
	newMPIJob.ResourceVersion = ""
	newMPIJob.UID = ""
	newMPIJob.Status = trainingv2.JobStatus{}

	// Delete old job
	if err := api.client.Delete(ctx, mpijob); err != nil {
		api.sendError(w, http.StatusInternalServerError, fmt.Errorf("failed to delete old MPIJob: %w", err))
		return
	}

	// Create new job
	if err := api.client.Create(ctx, newMPIJob); err != nil {
		api.sendError(w, http.StatusInternalServerError, fmt.Errorf("failed to create new MPIJob: %w", err))
		return
	}

	api.sendJSON(w, http.StatusCreated, newMPIJob)
}

// Helper methods

func (api *MPIJobAPI) sendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (api *MPIJobAPI) sendText(w http.ResponseWriter, statusCode int, text string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)
	w.Write([]byte(text))
}

func (api *MPIJobAPI) sendError(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}

func (api *MPIJobAPI) streamLogs(w http.ResponseWriter, logs string) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		api.sendError(w, http.StatusInternalServerError, fmt.Errorf("streaming not supported"))
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)

	// Send logs in chunks
	w.Write([]byte(logs))
	flusher.Flush()
}

func (api *MPIJobAPI) parseIntQuery(r *http.Request, param string, defaultValue int) int {
	valueStr := r.URL.Query().Get(param)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// LogOptions contains options for log retrieval
type LogOptions struct {
	Follow    bool
	TailLines int
	Since     string
}

// MPIJobService provides business logic for MPIJob operations
type MPIJobService struct {
	client client.Client
}

// NewMPIJobService creates a new MPIJob service
func NewMPIJobService(client client.Client) *MPIJobService {
	return &MPIJobService{
		client: client,
	}
}

// ValidateMPIJob validates an MPIJob specification
func (s *MPIJobService) ValidateMPIJob(mpijob *trainingv2.MPIJob) error {
	if mpijob.Name == "" {
		return fmt.Errorf("name is required")
	}
	if mpijob.Namespace == "" {
		return fmt.Errorf("namespace is required")
	}
	if mpijob.Spec.MPIReplicaSpecs == nil {
		return fmt.Errorf("mpiReplicaSpecs is required")
	}
	if _, ok := mpijob.Spec.MPIReplicaSpecs["Launcher"]; !ok {
		return fmt.Errorf("launcher spec is required")
	}
	if _, ok := mpijob.Spec.MPIReplicaSpecs["Worker"]; !ok {
		return fmt.Errorf("worker spec is required")
	}
	return nil
}

// GetJobStatus retrieves the status of an MPIJob
func (s *MPIJobService) GetJobStatus(ctx context.Context, namespace, name string) (*trainingv2.JobStatus, error) {
	mpijob := &trainingv2.MPIJob{}
	if err := s.client.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, mpijob); err != nil {
		return nil, err
	}
	return &mpijob.Status, nil
}

// GetJobLogs retrieves logs from all pods of an MPIJob
func (s *MPIJobService) GetJobLogs(ctx context.Context, namespace, name string, opts *LogOptions) (string, error) {
	// TODO: Implement log aggregation from all pods
	return "Logs from all pods", nil
}

// GetLauncherLogs retrieves logs from the launcher pod
func (s *MPIJobService) GetLauncherLogs(ctx context.Context, namespace, name string, opts *LogOptions) (string, error) {
	// TODO: Implement log retrieval from launcher pod
	return "Launcher logs", nil
}

// GetWorkerLogs retrieves logs from a specific worker pod
func (s *MPIJobService) GetWorkerLogs(ctx context.Context, namespace, name string, index int, opts *LogOptions) (string, error) {
	// TODO: Implement log retrieval from worker pod
	return fmt.Sprintf("Worker %d logs", index), nil
}

// CanRestart checks if an MPIJob can be restarted
func (s *MPIJobService) CanRestart(mpijob *trainingv2.MPIJob) bool {
	for _, condition := range mpijob.Status.Conditions {
		if condition.Type == string(trainingv2.JobFailed) && condition.Status == metav1.ConditionTrue {
			return true
		}
	}
	return false
}
