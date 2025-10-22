"""Client for interacting with MPIJob resources."""

import os
from typing import Any, Dict, List, Optional, Union

import yaml
from kubernetes import client, config
from pydantic import BaseModel, Field

from odh.training.models import (
    MPIImplementationType,
    MPIJob as MPIJobModel,
    MPIJobSpec,
    NetworkPolicy,
    PodTemplateSpec,
    ReplicaSpec,
    RunPolicy,
)
from odh.training.mpijob import MPIJob
from odh.training.resources import ResourceSpec


class MPIJobClient(BaseModel):
    """Client for managing MPIJob resources in OpenShift AI."""

    kubeconfig: Optional[str] = Field(
        None, description="Path to kubeconfig file, uses default if not specified"
    )
    namespace: str = Field("default", description="Default namespace for operations")
    _client: Optional[Any] = Field(None, description="Kubernetes client", exclude=True)

    class Config:
        """Pydantic model configuration."""

        arbitrary_types_allowed = True

    def __init__(self, **data: Any):
        """Initialize the client and load kubernetes configuration.

        Args:
            **data: Initialization parameters
        """
        super().__init__(**data)
        self._initialize_client()

    def _initialize_client(self) -> None:
        """Initialize the Kubernetes client."""
        if self.kubeconfig:
            config.load_kube_config(self.kubeconfig)
        else:
            try:
                # Try to load from default location
                config.load_kube_config()
            except config.config_exception.ConfigException:
                # If that fails, try in-cluster configuration
                try:
                    config.load_incluster_config()
                except config.config_exception.ConfigException:
                    raise ValueError(
                        "Could not load kubeconfig. "
                        "Specify a valid path or run from within a Kubernetes cluster."
                    )

        self._client = client.ApiClient()

    def create_mpijob(
        self,
        name: str,
        namespace: Optional[str] = None,
        worker_spec: Optional[Dict[str, Any]] = None,
        launcher_spec: Optional[Dict[str, Any]] = None,
        from_file: Optional[str] = None,
        from_dict: Optional[Dict[str, Any]] = None,
        mpi_implementation: str = "OpenMPI",
        slots_per_worker: int = 1,
        dry_run: bool = False,
    ) -> MPIJob:
        """Create a new MPIJob.

        Args:
            name: Name of the MPIJob
            namespace: Namespace for the MPIJob, uses default if not specified
            worker_spec: Worker specification
            launcher_spec: Launcher specification
            from_file: Load job specification from a file
            from_dict: Load job specification from a dictionary
            mpi_implementation: MPI implementation to use
            slots_per_worker: Number of slots per worker
            dry_run: Only validate the request, don't actually create the job

        Returns:
            MPIJob: Created MPIJob object
        """
        if sum(x is not None for x in [from_file, from_dict, worker_spec]) != 1:
            raise ValueError(
                "Exactly one of from_file, from_dict, or worker_spec must be provided"
            )

        namespace = namespace or self.namespace
        custom_api = client.CustomObjectsApi(self._client)

        # Build MPIJob from file
        if from_file:
            with open(from_file, "r") as f:
                mpijob_dict = yaml.safe_load(f)
            return self._create_from_dict(custom_api, mpijob_dict, namespace, dry_run)

        # Build MPIJob from dict
        if from_dict:
            return self._create_from_dict(custom_api, from_dict, namespace, dry_run)

        # Build MPIJob from worker_spec
        if worker_spec:
            mpijob_dict = self._build_mpijob_dict(
                name=name,
                namespace=namespace,
                worker_spec=worker_spec,
                launcher_spec=launcher_spec,
                mpi_implementation=mpi_implementation,
                slots_per_worker=slots_per_worker,
            )
            return self._create_from_dict(custom_api, mpijob_dict, namespace, dry_run)

        raise ValueError("No job specification provided")

    def _create_from_dict(
        self, custom_api: Any, mpijob_dict: Dict[str, Any], namespace: str, dry_run: bool
    ) -> MPIJob:
        """Create MPIJob from a dictionary.

        Args:
            custom_api: Kubernetes CustomObjects API
            mpijob_dict: MPIJob dictionary
            namespace: Namespace for the job
            dry_run: Only validate the request, don't actually create the job

        Returns:
            MPIJob: Created MPIJob object
        """
        # Validate with pydantic model
        job_model = MPIJobModel.model_validate(mpijob_dict)

        # Set namespace if not already set
        if "metadata" not in mpijob_dict:
            mpijob_dict["metadata"] = {}
        if "namespace" not in mpijob_dict["metadata"]:
            mpijob_dict["metadata"]["namespace"] = namespace
        if "apiVersion" not in mpijob_dict:
            mpijob_dict["apiVersion"] = "kubeflow.org/v2"
        if "kind" not in mpijob_dict:
            mpijob_dict["kind"] = "MPIJob"

        # Create the job
        result = custom_api.create_namespaced_custom_object(
            group="kubeflow.org",
            version="v2",
            namespace=mpijob_dict["metadata"]["namespace"],
            plural="mpijobs",
            body=mpijob_dict,
            dry_run="All" if dry_run else None,
        )

        # Create MPIJob object
        return MPIJob(
            name=result["metadata"]["name"],
            namespace=result["metadata"]["namespace"],
            client=self._client,
            raw=result,
        )

    def _build_mpijob_dict(
        self,
        name: str,
        namespace: str,
        worker_spec: Dict[str, Any],
        launcher_spec: Optional[Dict[str, Any]] = None,
        mpi_implementation: str = "OpenMPI",
        slots_per_worker: int = 1,
    ) -> Dict[str, Any]:
        """Build MPIJob dictionary from components.

        Args:
            name: Name of the MPIJob
            namespace: Namespace for the MPIJob
            worker_spec: Worker specification
            launcher_spec: Launcher specification
            mpi_implementation: MPI implementation to use
            slots_per_worker: Number of slots per worker

        Returns:
            Dict[str, Any]: MPIJob dictionary
        """
        # Create default launcher spec if not provided
        if launcher_spec is None:
            launcher_spec = {
                "replicas": 1,
                "resources": {"requests": {"cpu": "1", "memory": "1Gi"}},
                "image": worker_spec.get("image"),
            }

        # Build MPIJob spec
        mpi_replica_specs = {
            "Worker": self._build_replica_spec("Worker", worker_spec),
            "Launcher": self._build_replica_spec("Launcher", launcher_spec),
        }

        # Build complete MPIJob
        return {
            "apiVersion": "kubeflow.org/v2",
            "kind": "MPIJob",
            "metadata": {
                "name": name,
                "namespace": namespace,
            },
            "spec": {
                "slotsPerWorker": slots_per_worker,
                "mpiImplementation": mpi_implementation,
                "mpiReplicaSpecs": mpi_replica_specs,
            },
        }

    def _build_replica_spec(self, replica_type: str, spec: Dict[str, Any]) -> Dict[str, Any]:
        """Build a ReplicaSpec for an MPIJob component.

        Args:
            replica_type: Type of replica (Worker or Launcher)
            spec: Specification for the replica

        Returns:
            Dict[str, Any]: ReplicaSpec dictionary
        """
        container = {
            "name": f"mpi-{replica_type.lower()}",
            "image": spec["image"],
            "resources": spec["resources"],
        }

        if "command" in spec and spec["command"]:
            container["command"] = spec["command"]

        if "args" in spec and spec["args"]:
            container["args"] = spec["args"]

        if "env" in spec and spec["env"]:
            container["env"] = spec["env"]

        if "working_dir" in spec and spec["working_dir"]:
            container["workingDir"] = spec["working_dir"]

        if "volume_mounts" in spec and spec["volume_mounts"]:
            container["volumeMounts"] = spec["volume_mounts"]

        pod_spec = {"containers": [container]}

        if "volumes" in spec and spec["volumes"]:
            pod_spec["volumes"] = spec["volumes"]

        return {
            "replicas": spec["replicas"],
            "template": {
                "spec": pod_spec,
            },
        }

    def get_mpijob(self, name: str, namespace: Optional[str] = None) -> MPIJob:
        """Get an MPIJob by name.

        Args:
            name: Name of the MPIJob
            namespace: Namespace of the MPIJob, uses default if not specified

        Returns:
            MPIJob: MPIJob object

        Raises:
            ValueError: If job is not found
        """
        namespace = namespace or self.namespace
        custom_api = client.CustomObjectsApi(self._client)

        try:
            job = custom_api.get_namespaced_custom_object(
                group="kubeflow.org",
                version="v2",
                namespace=namespace,
                plural="mpijobs",
                name=name,
            )

            return MPIJob(name=name, namespace=namespace, client=self._client, raw=job)

        except client.rest.ApiException as e:
            if e.status == 404:
                raise ValueError(f"MPIJob {name} not found in namespace {namespace}")
            raise

    def list_mpijobs(
        self, namespace: Optional[str] = None, label_selector: Optional[str] = None
    ) -> List[MPIJob]:
        """List MPIJobs in a namespace.

        Args:
            namespace: Namespace to list jobs in, uses default if not specified
            label_selector: Label selector to filter jobs

        Returns:
            List[MPIJob]: List of MPIJob objects
        """
        namespace = namespace or self.namespace
        custom_api = client.CustomObjectsApi(self._client)

        jobs = custom_api.list_namespaced_custom_object(
            group="kubeflow.org",
            version="v2",
            namespace=namespace,
            plural="mpijobs",
            label_selector=label_selector,
        )

        return [
            MPIJob(
                name=job["metadata"]["name"],
                namespace=namespace,
                client=self._client,
                raw=job,
            )
            for job in jobs.get("items", [])
        ]

    def delete_mpijob(
        self, name: str, namespace: Optional[str] = None, wait: bool = False
    ) -> bool:
        """Delete an MPIJob.

        Args:
            name: Name of the MPIJob
            namespace: Namespace of the MPIJob, uses default if not specified
            wait: Wait for job to be deleted

        Returns:
            bool: True if deletion succeeded, False otherwise
        """
        namespace = namespace or self.namespace

        try:
            job = self.get_mpijob(name, namespace)
            return job.delete(wait=wait)
        except ValueError:
            # Job already deleted
            return True