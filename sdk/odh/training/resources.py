"""Resource specification utilities for MPIJob support."""

from typing import Dict, Optional, Union

from pydantic import BaseModel, Field


class ResourceSpec(BaseModel):
    """Resource specification for MPIJob components."""

    cpu: Optional[str] = Field(None, description="CPU resources (e.g., '1', '500m')")
    memory: Optional[str] = Field(None, description="Memory resources (e.g., '1Gi', '512Mi')")
    gpu: Optional[int] = Field(None, description="Number of GPUs")
    gpu_vendor: str = Field("nvidia.com", description="GPU vendor")
    gpu_product: str = Field("gpu", description="GPU product name")

    def to_resource_requirements(self) -> Dict[str, Dict[str, str]]:
        """Convert to Kubernetes resource requirements format.

        Returns:
            Dict[str, Dict[str, str]]: Resource requirements with requests and limits
        """
        requests = {}
        limits = {}

        if self.cpu:
            requests["cpu"] = self.cpu
            limits["cpu"] = self.cpu

        if self.memory:
            requests["memory"] = self.memory
            limits["memory"] = self.memory

        if self.gpu:
            gpu_resource = f"{self.gpu_vendor}/{self.gpu_product}"
            requests[gpu_resource] = str(self.gpu)
            limits[gpu_resource] = str(self.gpu)

        return {"requests": requests, "limits": limits}

    @classmethod
    def from_resource_requirements(
        cls, resources: Dict[str, Dict[str, Union[str, int]]]
    ) -> "ResourceSpec":
        """Create ResourceSpec from Kubernetes resource requirements.

        Args:
            resources: Resource requirements with requests and limits

        Returns:
            ResourceSpec: ResourceSpec instance
        """
        cpu = None
        memory = None
        gpu = None
        gpu_vendor = "nvidia.com"
        gpu_product = "gpu"

        # Extract from requests section
        requests = resources.get("requests", {})
        if "cpu" in requests:
            cpu = str(requests["cpu"])
        if "memory" in requests:
            memory = str(requests["memory"])

        # Look for GPU in both requests and limits
        for section in [requests, resources.get("limits", {})]:
            for key, value in section.items():
                if "gpu" in key.lower():
                    # Extract vendor and product from resource name
                    parts = key.split("/")
                    if len(parts) == 2:
                        gpu_vendor = parts[0]
                        gpu_product = parts[1]
                    gpu = int(value)
                    break

        return cls(
            cpu=cpu, memory=memory, gpu=gpu, gpu_vendor=gpu_vendor, gpu_product=gpu_product
        )

    def calculate_total(self, replicas: int) -> "ResourceSpec":
        """Calculate total resources for multiple replicas.

        Args:
            replicas: Number of replicas

        Returns:
            ResourceSpec: New ResourceSpec with total resources
        """
        result = ResourceSpec()

        # Handle CPU (simple multiplication not accurate for CPU values like "500m")
        if self.cpu:
            try:
                # If it's a simple number, multiply
                cpu_value = float(self.cpu)
                result.cpu = str(cpu_value * replicas)
            except ValueError:
                # Handle CPU millicores (e.g., "500m")
                if self.cpu.endswith("m"):
                    cpu_millicores = int(self.cpu[:-1])
                    total_millicores = cpu_millicores * replicas
                    if total_millicores >= 1000:
                        result.cpu = f"{total_millicores // 1000}.{(total_millicores % 1000) // 100}"
                    else:
                        result.cpu = f"{total_millicores}m"
                else:
                    # For other formats, just indicate it's a multiple
                    result.cpu = f"{self.cpu}*{replicas}"

        # Handle memory
        if self.memory:
            # Check for common memory units
            memory_units = ["Ki", "Mi", "Gi", "Ti"]
            unit = ""
            value = 0

            for unit_suffix in memory_units:
                if self.memory.endswith(unit_suffix):
                    unit = unit_suffix
                    value = int(self.memory[:-len(unit_suffix)])
                    break

            if unit:
                result.memory = f"{value * replicas}{unit}"
            else:
                try:
                    # If it's a simple number, multiply
                    memory_value = float(self.memory)
                    result.memory = str(memory_value * replicas)
                except ValueError:
                    # For other formats, just indicate it's a multiple
                    result.memory = f"{self.memory}*{replicas}"

        # Handle GPU
        if self.gpu:
            result.gpu = self.gpu * replicas
            result.gpu_vendor = self.gpu_vendor
            result.gpu_product = self.gpu_product

        return result

    def __str__(self) -> str:
        """Get string representation of resources.

        Returns:
            str: String representation
        """
        parts = []
        if self.cpu:
            parts.append(f"{self.cpu} CPU")
        if self.memory:
            parts.append(f"{self.memory} Memory")
        if self.gpu:
            parts.append(f"{self.gpu} GPU ({self.gpu_vendor}/{self.gpu_product})")

        return ", ".join(parts) if parts else "No resources specified"