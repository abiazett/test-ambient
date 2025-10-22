"""Data models for MPIJob support in OpenShift AI."""

from enum import Enum
from typing import Dict, List, Optional, Union

from pydantic import BaseModel, Field


class MPIImplementationType(str, Enum):
    """MPI implementation types supported by the operator."""

    OPENMPI = "OpenMPI"
    INTELMPI = "IntelMPI"
    MPICH = "MPICH"


class CleanPodPolicy(str, Enum):
    """Clean pod policy options for MPIJob."""

    ALL = "All"
    RUNNING = "Running"
    NONE = "None"


class RestartPolicy(str, Enum):
    """Restart policy options for MPIJob pods."""

    ALWAYS = "Always"
    ON_FAILURE = "OnFailure"
    NEVER = "Never"


class NetworkPolicyTemplate(str, Enum):
    """Network policy template options for MPIJob."""

    DEFAULT = "Default"
    RESTRICTED = "Restricted"


class RunPolicy(BaseModel):
    """Run policy configuration for MPIJob."""

    clean_pod_policy: Optional[CleanPodPolicy] = Field(
        None, alias="cleanPodPolicy", description="Policy for cleaning up pods after job completion"
    )
    ttl_seconds_after_finished: Optional[int] = Field(
        None, alias="ttlSecondsAfterFinished", description="TTL seconds after job is finished"
    )
    active_deadline_seconds: Optional[int] = Field(
        None, alias="activeDeadlineSeconds", description="Active deadline seconds for the job"
    )
    backoff_limit: Optional[int] = Field(
        None, alias="backoffLimit", description="Backoff limit for the job"
    )
    scheduling_policy: Optional[Dict[str, str]] = Field(
        None, alias="schedulingPolicy", description="Scheduling policy for the job"
    )

    class Config:
        """Pydantic model configuration."""

        populate_by_name = True
        validate_assignment = True
        arbitrary_types_allowed = True


class PodTemplateSpec(BaseModel):
    """Kubernetes pod template spec."""

    spec: Dict = Field(..., description="Pod spec definition")

    class Config:
        """Pydantic model configuration."""

        populate_by_name = True
        validate_assignment = True
        arbitrary_types_allowed = True


class ReplicaSpec(BaseModel):
    """Replica specification for a component of the MPIJob."""

    replicas: int = Field(..., description="Number of replicas")
    restart_policy: Optional[RestartPolicy] = Field(
        None, alias="restartPolicy", description="Restart policy for the replicas"
    )
    template: PodTemplateSpec = Field(..., description="Pod template specification")

    class Config:
        """Pydantic model configuration."""

        populate_by_name = True
        validate_assignment = True
        arbitrary_types_allowed = True


class NetworkPolicy(BaseModel):
    """Network policy configuration for MPIJob."""

    template: Optional[NetworkPolicyTemplate] = Field(
        None, description="Network policy template to use"
    )

    class Config:
        """Pydantic model configuration."""

        populate_by_name = True
        validate_assignment = True


class MPIJobSpec(BaseModel):
    """Specification for MPIJob."""

    slots_per_worker: int = Field(
        1, alias="slotsPerWorker", description="Number of slots per worker"
    )
    mpi_implementation: MPIImplementationType = Field(
        MPIImplementationType.OPENMPI,
        alias="mpiImplementation",
        description="MPI implementation to use",
    )
    run_policy: Optional[RunPolicy] = Field(
        None, alias="runPolicy", description="Run policy for the job"
    )
    mpi_replica_specs: Dict[str, ReplicaSpec] = Field(
        ..., alias="mpiReplicaSpecs", description="MPI replica specifications"
    )
    network_policy: Optional[NetworkPolicy] = Field(
        None, alias="networkPolicy", description="Network policy configuration"
    )

    class Config:
        """Pydantic model configuration."""

        populate_by_name = True
        validate_assignment = True
        arbitrary_types_allowed = True


class ReplicaStatus(BaseModel):
    """Status of replicas in an MPIJob."""

    active: Optional[int] = Field(None, description="Number of active replicas")
    succeeded: Optional[int] = Field(None, description="Number of succeeded replicas")
    failed: Optional[int] = Field(None, description="Number of failed replicas")

    class Config:
        """Pydantic model configuration."""

        populate_by_name = True
        validate_assignment = True


class JobConditionType(str, Enum):
    """Condition types for MPIJob status."""

    CREATED = "Created"
    RUNNING = "Running"
    RESTARTING = "Restarting"
    SUCCEEDED = "Succeeded"
    FAILED = "Failed"


class JobCondition(BaseModel):
    """Condition in an MPIJob status."""

    type: JobConditionType = Field(..., description="Type of job condition")
    status: str = Field(..., description="Status of the condition, True, False, Unknown")
    reason: Optional[str] = Field(None, description="Reason for the condition's last transition")
    message: Optional[str] = Field(None, description="Human readable message about the transition")
    last_transition_time: Optional[str] = Field(
        None, alias="lastTransitionTime", description="Last time the condition transitioned"
    )
    last_update_time: Optional[str] = Field(
        None, alias="lastUpdateTime", description="Last time the condition was updated"
    )

    class Config:
        """Pydantic model configuration."""

        populate_by_name = True
        validate_assignment = True


class MPIJobStatus(BaseModel):
    """Status of an MPIJob."""

    conditions: Optional[List[JobCondition]] = Field(None, description="Job conditions")
    start_time: Optional[str] = Field(None, alias="startTime", description="Job start time")
    completion_time: Optional[str] = Field(
        None, alias="completionTime", description="Job completion time"
    )
    replica_statuses: Optional[Dict[str, ReplicaStatus]] = Field(
        None, alias="replicaStatuses", description="Status of replicas"
    )

    class Config:
        """Pydantic model configuration."""

        populate_by_name = True
        validate_assignment = True


class MPIJob(BaseModel):
    """MPIJob resource model."""

    api_version: str = Field("kubeflow.org/v2", alias="apiVersion", description="API version")
    kind: str = Field("MPIJob", description="Resource kind")
    metadata: Dict = Field(..., description="Job metadata")
    spec: MPIJobSpec = Field(..., description="Job specification")
    status: Optional[MPIJobStatus] = Field(None, description="Job status")

    class Config:
        """Pydantic model configuration."""

        populate_by_name = True
        validate_assignment = True
        arbitrary_types_allowed = True


class WorkerSpec(BaseModel):
    """High-level specification for MPIJob workers."""

    replicas: int = Field(..., description="Number of worker replicas")
    resources: Dict[str, Dict[str, Union[str, int]]] = Field(
        ..., description="Resource requests and limits"
    )
    image: str = Field(..., description="Container image to use")
    command: Optional[List[str]] = Field(None, description="Command to run in the container")
    args: Optional[List[str]] = Field(None, description="Arguments to the command")
    env: Optional[List[Dict[str, str]]] = Field(None, description="Environment variables")
    working_dir: Optional[str] = Field(None, description="Working directory in the container")
    volumes: Optional[List[Dict]] = Field(None, description="Volumes to mount")
    volume_mounts: Optional[List[Dict]] = Field(
        None, alias="volumeMounts", description="Volume mounts"
    )

    class Config:
        """Pydantic model configuration."""

        populate_by_name = True
        validate_assignment = True