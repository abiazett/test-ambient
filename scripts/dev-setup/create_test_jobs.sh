#!/bin/bash
set -e

# Create test MPIJobs for development and testing
# This script creates a set of sample MPIJob instances with various configurations

# Ensure we're using the correct cluster
kubectl config use-context kind-mpijob-dev

# Create namespace if it doesn't exist
kubectl create namespace mpijob-test || true

# Basic MPIJob with 4 workers
create_basic_mpijob() {
    echo "Creating basic MPIJob with 4 workers..."

    cat > basic-mpijob.yaml << EOF
apiVersion: kubeflow.org/v2
kind: MPIJob
metadata:
  name: basic-mpijob
  namespace: mpijob-test
spec:
  slotsPerWorker: 1
  mpiImplementation: OpenMPI
  runPolicy:
    cleanPodPolicy: Running
  mpiReplicaSpecs:
    Launcher:
      replicas: 1
      restartPolicy: Never
      template:
        spec:
          containers:
          - name: mpi-launcher
            image: mpioperator/tensorflow-benchmarks:latest
            command:
            - mpirun
            - --allow-run-as-root
            - -np
            - "4"
            - -bind-to
            - none
            - -map-by
            - slot
            - -x
            - NCCL_DEBUG=INFO
            - python
            - scripts/tf_cnn_benchmarks/tf_cnn_benchmarks.py
            - --model=resnet50
            - --batch_size=64
            resources:
              limits:
                cpu: 1
                memory: 2Gi
              requests:
                cpu: 500m
                memory: 1Gi
    Worker:
      replicas: 4
      restartPolicy: Never
      template:
        spec:
          containers:
          - name: mpi-worker
            image: mpioperator/tensorflow-benchmarks:latest
            resources:
              limits:
                cpu: 2
                memory: 4Gi
              requests:
                cpu: 1
                memory: 2Gi
EOF

    kubectl apply -f basic-mpijob.yaml
    echo "Basic MPIJob created successfully."
}

# MPIJob with GPUs
create_gpu_mpijob() {
    echo "Creating MPIJob with GPU resources..."

    cat > gpu-mpijob.yaml << EOF
apiVersion: kubeflow.org/v2
kind: MPIJob
metadata:
  name: gpu-mpijob
  namespace: mpijob-test
spec:
  slotsPerWorker: 2
  mpiImplementation: OpenMPI
  runPolicy:
    cleanPodPolicy: Running
  mpiReplicaSpecs:
    Launcher:
      replicas: 1
      restartPolicy: Never
      template:
        spec:
          containers:
          - name: mpi-launcher
            image: mpioperator/tensorflow-benchmarks:latest
            command:
            - mpirun
            - --allow-run-as-root
            - -np
            - "8"
            - -bind-to
            - none
            - -map-by
            - slot
            - -x
            - NCCL_DEBUG=INFO
            - python
            - scripts/tf_cnn_benchmarks/tf_cnn_benchmarks.py
            - --model=resnet50
            - --batch_size=64
            - --variable_update=horovod
            - --use_fp16
            resources:
              limits:
                cpu: 2
                memory: 4Gi
              requests:
                cpu: 1
                memory: 2Gi
    Worker:
      replicas: 4
      restartPolicy: Never
      template:
        spec:
          containers:
          - name: mpi-worker
            image: mpioperator/tensorflow-benchmarks:latest
            resources:
              limits:
                cpu: 4
                memory: 16Gi
                nvidia.com/gpu: 2
              requests:
                cpu: 2
                memory: 8Gi
                nvidia.com/gpu: 2
EOF

    kubectl apply -f gpu-mpijob.yaml
    echo "GPU MPIJob created successfully."
}

# MPIJob with custom MPI implementation
create_intel_mpi_job() {
    echo "Creating MPIJob with Intel MPI implementation..."

    cat > intel-mpijob.yaml << EOF
apiVersion: kubeflow.org/v2
kind: MPIJob
metadata:
  name: intel-mpijob
  namespace: mpijob-test
spec:
  slotsPerWorker: 1
  mpiImplementation: IntelMPI
  runPolicy:
    cleanPodPolicy: Running
  mpiReplicaSpecs:
    Launcher:
      replicas: 1
      restartPolicy: Never
      template:
        spec:
          containers:
          - name: mpi-launcher
            image: intel/oneapi-hpckit:latest
            command:
            - mpiexec.hydra
            - -n
            - "4"
            - -genv
            - I_MPI_DEBUG=5
            - hostname
            resources:
              limits:
                cpu: 1
                memory: 2Gi
              requests:
                cpu: 500m
                memory: 1Gi
    Worker:
      replicas: 4
      restartPolicy: Never
      template:
        spec:
          containers:
          - name: mpi-worker
            image: intel/oneapi-hpckit:latest
            resources:
              limits:
                cpu: 2
                memory: 4Gi
              requests:
                cpu: 1
                memory: 2Gi
EOF

    kubectl apply -f intel-mpijob.yaml
    echo "Intel MPI job created successfully."
}

# MPIJob with network policy
create_mpijob_with_network_policy() {
    echo "Creating MPIJob with network policy..."

    cat > network-mpijob.yaml << EOF
apiVersion: kubeflow.org/v2
kind: MPIJob
metadata:
  name: network-mpijob
  namespace: mpijob-test
spec:
  slotsPerWorker: 1
  mpiImplementation: OpenMPI
  runPolicy:
    cleanPodPolicy: Running
  networkPolicy:
    template: Restricted
  mpiReplicaSpecs:
    Launcher:
      replicas: 1
      restartPolicy: Never
      template:
        spec:
          containers:
          - name: mpi-launcher
            image: mpioperator/tensorflow-benchmarks:latest
            command:
            - mpirun
            - --allow-run-as-root
            - -np
            - "2"
            - -bind-to
            - none
            - -map-by
            - slot
            - -x
            - NCCL_DEBUG=INFO
            - python
            - scripts/tf_cnn_benchmarks/tf_cnn_benchmarks.py
            - --model=resnet50
            - --batch_size=64
            resources:
              limits:
                cpu: 1
                memory: 2Gi
              requests:
                cpu: 500m
                memory: 1Gi
    Worker:
      replicas: 2
      restartPolicy: Never
      template:
        spec:
          containers:
          - name: mpi-worker
            image: mpioperator/tensorflow-benchmarks:latest
            resources:
              limits:
                cpu: 2
                memory: 4Gi
              requests:
                cpu: 1
                memory: 2Gi
EOF

    kubectl apply -f network-mpijob.yaml
    echo "MPIJob with network policy created successfully."
}

# Main function
main() {
    echo "=== Creating Test MPIJobs ==="

    # Create sample jobs
    create_basic_mpijob
    create_gpu_mpijob
    create_intel_mpi_job
    create_mpijob_with_network_policy

    echo "=== Test Jobs Created ==="
    echo "You can monitor the jobs with:"
    echo "kubectl get mpijobs -n mpijob-test"
    echo "kubectl describe mpijob <job-name> -n mpijob-test"
}

# Execute main function
main