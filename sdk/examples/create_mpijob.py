#!/usr/bin/env python3
"""Example of creating and monitoring an MPIJob using the OpenShift AI SDK."""

import argparse
import time
from typing import Dict, Optional

from odh.training import MPIJobClient, MPIJobStatus


def create_sample_mpijob(
    name: str, namespace: str, workers: int, gpus: int, cpu: str, memory: str, image: str
) -> None:
    """Create a sample MPIJob and monitor its execution.

    Args:
        name: Name of the job
        namespace: Namespace to create the job in
        workers: Number of worker replicas
        gpus: Number of GPUs per worker
        cpu: CPU request per worker
        memory: Memory request per worker
        image: Container image to use
    """
    # Initialize client
    client = MPIJobClient(namespace=namespace)

    # Define worker spec
    worker_spec = {
        "replicas": workers,
        "resources": {
            "requests": {
                "cpu": cpu,
                "memory": memory,
                "nvidia.com/gpu": gpus,
            },
            "limits": {
                "cpu": cpu,
                "memory": memory,
                "nvidia.com/gpu": gpus,
            },
        },
        "image": image,
        "command": [
            "mpirun",
            "--allow-run-as-root",
            "-np",
            str(workers * gpus),
            "-bind-to",
            "none",
            "-map-by",
            "slot",
            "-x",
            "NCCL_DEBUG=INFO",
            "python",
            "/train.py",
            "--epochs",
            "10",
        ],
    }

    print(f"Creating MPIJob {name} in namespace {namespace}")
    print(f"  Workers: {workers}")
    print(f"  GPUs per worker: {gpus}")
    print(f"  CPU per worker: {cpu}")
    print(f"  Memory per worker: {memory}")
    print(f"  Image: {image}")

    # Create MPIJob
    try:
        job = client.create_mpijob(name=name, namespace=namespace, worker_spec=worker_spec)
        print(f"Created MPIJob: {job.name}")
    except Exception as e:
        print(f"Failed to create MPIJob: {e}")
        return

    # Monitor job status
    def status_callback(status: MPIJobStatus) -> None:
        """Callback for job status updates.

        Args:
            status: Current job status
        """
        print(f"Job status changed: {status}")

    try:
        print("Monitoring job status...")
        job.monitor(callback=status_callback, poll_interval=5)

        # Check final status
        job.refresh()
        if job.is_succeeded:
            print("Job completed successfully!")
        else:
            print(f"Job failed: {job.phase}")

        # Get logs from launcher
        print("\nLauncher logs:")
        logs = job.get_logs(worker="launcher")
        for pod_name, log_content in logs.items():
            print(f"Pod {pod_name}:")
            print(log_content[:500] + "..." if len(log_content) > 500 else log_content)

    except KeyboardInterrupt:
        print("Monitoring interrupted. Job will continue running.")
    except Exception as e:
        print(f"Error monitoring job: {e}")


def parse_args() -> Dict[str, Optional[str]]:
    """Parse command line arguments.

    Returns:
        Dict[str, Optional[str]]: Dictionary of arguments
    """
    parser = argparse.ArgumentParser(description="Create a sample MPIJob")
    parser.add_argument("--name", required=True, help="Name of the MPIJob")
    parser.add_argument("--namespace", default="default", help="Namespace to create the job in")
    parser.add_argument("--workers", type=int, default=2, help="Number of worker replicas")
    parser.add_argument("--gpus", type=int, default=1, help="Number of GPUs per worker")
    parser.add_argument("--cpu", default="1", help="CPU request per worker")
    parser.add_argument("--memory", default="4Gi", help="Memory request per worker")
    parser.add_argument(
        "--image",
        default="mpioperator/tensorflow-benchmarks:latest",
        help="Container image to use",
    )
    return vars(parser.parse_args())


def main() -> None:
    """Run the example."""
    args = parse_args()
    create_sample_mpijob(
        name=args["name"],
        namespace=args["namespace"],
        workers=args["workers"],
        gpus=args["gpus"],
        cpu=args["cpu"],
        memory=args["memory"],
        image=args["image"],
    )


if __name__ == "__main__":
    main()