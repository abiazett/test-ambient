"""MPIJob implementation for OpenShift AI SDK."""

import time
from typing import Any, Callable, Dict, List, Optional, Union

from kubernetes import client
from pydantic import BaseModel, Field

from odh.training.models import JobConditionType, MPIJob as MPIJobModel
from odh.training.models import MPIJobStatus


class MPIJob(BaseModel):
    """High-level MPIJob interface for OpenShift AI SDK."""

    name: str = Field(..., description="Name of the MPIJob")
    namespace: str = Field(..., description="Namespace of the MPIJob")
    client: Any = Field(..., description="Kubernetes client")
    raw: Optional[MPIJobModel] = Field(None, description="Raw MPIJob resource")

    class Config:
        """Pydantic model configuration."""

        arbitrary_types_allowed = True

    def refresh(self) -> "MPIJob":
        """Refresh job status from the API.

        Returns:
            MPIJob: Self with refreshed data
        """
        api_client = client.ApiClient()
        custom_api = client.CustomObjectsApi(api_client)

        self.raw = custom_api.get_namespaced_custom_object(
            group="kubeflow.org",
            version="v2",
            namespace=self.namespace,
            plural="mpijobs",
            name=self.name,
        )

        return self

    @property
    def status(self) -> Optional[MPIJobStatus]:
        """Get current job status.

        Returns:
            Optional[MPIJobStatus]: Job status or None if not available
        """
        if not self.raw or "status" not in self.raw:
            return None

        return MPIJobStatus.model_validate(self.raw["status"])

    @property
    def phase(self) -> str:
        """Get current job phase.

        Returns:
            str: Job phase (Created, Running, Succeeded, Failed, Unknown)
        """
        status = self.status
        if not status or not status.conditions:
            return "Unknown"

        # Check for terminal conditions first
        for condition in status.conditions:
            if (
                condition.type == JobConditionType.SUCCEEDED
                and condition.status.lower() == "true"
            ):
                return "Succeeded"
            if condition.type == JobConditionType.FAILED and condition.status.lower() == "true":
                return "Failed"

        # Then check for running condition
        for condition in status.conditions:
            if condition.type == JobConditionType.RUNNING and condition.status.lower() == "true":
                return "Running"

        # Then check for created condition
        for condition in status.conditions:
            if condition.type == JobConditionType.CREATED and condition.status.lower() == "true":
                return "Created"

        return "Unknown"

    @property
    def workers_running(self) -> int:
        """Get number of currently running workers.

        Returns:
            int: Number of running workers
        """
        status = self.status
        if (
            not status
            or not status.replica_statuses
            or "Worker" not in status.replica_statuses
        ):
            return 0

        return status.replica_statuses["Worker"].active or 0

    @property
    def workers_succeeded(self) -> int:
        """Get number of succeeded workers.

        Returns:
            int: Number of succeeded workers
        """
        status = self.status
        if (
            not status
            or not status.replica_statuses
            or "Worker" not in status.replica_statuses
        ):
            return 0

        return status.replica_statuses["Worker"].succeeded or 0

    @property
    def workers_failed(self) -> int:
        """Get number of failed workers.

        Returns:
            int: Number of failed workers
        """
        status = self.status
        if (
            not status
            or not status.replica_statuses
            or "Worker" not in status.replica_statuses
        ):
            return 0

        return status.replica_statuses["Worker"].failed or 0

    @property
    def is_running(self) -> bool:
        """Check if job is currently running.

        Returns:
            bool: True if job is running
        """
        return self.phase == "Running"

    @property
    def is_completed(self) -> bool:
        """Check if job has completed (succeeded or failed).

        Returns:
            bool: True if job has completed
        """
        return self.phase in ["Succeeded", "Failed"]

    @property
    def is_succeeded(self) -> bool:
        """Check if job has succeeded.

        Returns:
            bool: True if job has succeeded
        """
        return self.phase == "Succeeded"

    @property
    def is_failed(self) -> bool:
        """Check if job has failed.

        Returns:
            bool: True if job has failed
        """
        return self.phase == "Failed"

    def wait_for_completion(self, timeout: int = 3600, poll_interval: int = 10) -> bool:
        """Wait for job to complete.

        Args:
            timeout: Maximum time to wait in seconds
            poll_interval: Time between status checks in seconds

        Returns:
            bool: True if job succeeded, False if failed or timed out
        """
        start_time = time.time()
        while time.time() - start_time < timeout:
            self.refresh()
            if self.is_completed:
                return self.is_succeeded

            time.sleep(poll_interval)

        # Timeout
        return False

    def monitor(
        self, callback: Optional[Callable[[MPIJobStatus], None]] = None, poll_interval: int = 10
    ) -> None:
        """Monitor job status with callback for updates.

        Args:
            callback: Function to call when status changes
            poll_interval: Time between status checks in seconds
        """
        last_phase = None
        while True:
            self.refresh()
            current_phase = self.phase

            if current_phase != last_phase:
                status = self.status
                if callback and status:
                    callback(status)
                last_phase = current_phase

            if self.is_completed:
                break

            time.sleep(poll_interval)

    def delete(self, wait: bool = False, timeout: int = 60) -> bool:
        """Delete the job.

        Args:
            wait: Wait for job to be deleted
            timeout: Maximum time to wait in seconds

        Returns:
            bool: True if deletion succeeded, False otherwise
        """
        api_client = client.ApiClient()
        custom_api = client.CustomObjectsApi(api_client)

        try:
            custom_api.delete_namespaced_custom_object(
                group="kubeflow.org",
                version="v2",
                namespace=self.namespace,
                plural="mpijobs",
                name=self.name,
            )

            if not wait:
                return True

            # Wait for job to be deleted
            start_time = time.time()
            while time.time() - start_time < timeout:
                try:
                    self.refresh()
                    time.sleep(1)
                except client.rest.ApiException as e:
                    if e.status == 404:
                        return True
                    raise

            # Timeout waiting for deletion
            return False

        except client.rest.ApiException as e:
            if e.status == 404:
                # Already deleted
                return True
            raise

    def get_logs(
        self, worker: Optional[Union[int, str]] = None, container: Optional[str] = None
    ) -> Dict[str, str]:
        """Get logs from job pods.

        Args:
            worker: Worker index or 'launcher', None for launcher
            container: Container name, None for default

        Returns:
            Dict[str, str]: Dictionary of pod name to logs
        """
        # Default to launcher if worker not specified
        worker_name = "launcher" if worker is None else worker

        api_client = client.ApiClient()
        core_api = client.CoreV1Api(api_client)

        # Construct label selector for pods
        if worker_name == "launcher":
            label_selector = (
                f"training.kubeflow.org/job-name={self.name},training.kubeflow.org/replica-type=Launcher"
            )
        elif isinstance(worker_name, int):
            label_selector = (
                f"training.kubeflow.org/job-name={self.name},training.kubeflow.org/replica-type=Worker,"
                f"training.kubeflow.org/replica-index={worker_name}"
            )
        else:
            label_selector = (
                f"training.kubeflow.org/job-name={self.name},training.kubeflow.org/replica-type=Worker"
            )

        # Get pods matching the label selector
        pods = core_api.list_namespaced_pod(
            namespace=self.namespace, label_selector=label_selector
        )

        logs = {}
        for pod in pods.items:
            pod_name = pod.metadata.name
            try:
                logs[pod_name] = core_api.read_namespaced_pod_log(
                    name=pod_name, namespace=self.namespace, container=container
                )
            except client.rest.ApiException as e:
                logs[pod_name] = f"Error getting logs: {e.reason}"

        return logs