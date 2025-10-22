# MPIJob Test Manifests

This directory contains test manifests for validating MPIJob functionality in OpenShift AI.

## Test Manifests

### 1. simple-mpijob.yaml
**Purpose**: Basic MPIJob functionality test
**Configuration**:
- 1 Launcher pod
- 2 Worker pods
- 1 slot per worker
- OpenMPI implementation
- Minimal resource requirements

**Use Case**: Smoke test to verify basic MPIJob creation, scheduling, and execution

**Expected Behavior**:
- All pods should be created successfully
- Launcher should initialize MPI environment
- Workers should accept connections
- Job should complete successfully

### 2. gpu-mpijob.yaml
**Purpose**: GPU-enabled distributed training test
**Configuration**:
- 1 Launcher pod
- 4 Worker pods
- 2 GPUs per worker (8 total)
- 2 slots per worker
- OpenMPI with NCCL
- GPU resource requirements

**Use Case**: Validate GPU allocation, NCCL communication, and multi-GPU training

**Expected Behavior**:
- GPU resources allocated correctly
- NCCL communication established between GPUs
- Distributed training completes successfully
- GPU utilization metrics available

### 3. intel-mpi-job.yaml
**Purpose**: Intel MPI implementation compatibility test
**Configuration**:
- 1 Launcher pod
- 4 Worker pods
- 4 slots per worker
- Intel MPI implementation
- IMB-MPI1 benchmark workload

**Use Case**: Verify Intel MPI implementation support and compatibility

**Expected Behavior**:
- Intel MPI environment configured correctly
- MPI communication using Intel MPI fabrics
- Benchmark completes successfully
- Performance metrics captured

### 4. volcano-scheduling.yaml
**Purpose**: Gang scheduling with Volcano integration test
**Configuration**:
- 1 Launcher pod
- 4 Worker pods
- Volcano scheduler
- Priority class and queue
- Gang scheduling annotations

**Use Case**: Validate gang scheduling prevents partial job placement

**Expected Behavior**:
- All pods scheduled together or none at all
- No deadlocks from partial scheduling
- Queue priorities respected
- Resource allocation efficient

### 5. network-policy-restricted.yaml
**Purpose**: Network isolation and security policy test
**Configuration**:
- 1 Launcher pod
- 2 Worker pods
- Restricted NetworkPolicy
- Security context enforced
- Non-root user execution

**Use Case**: Validate network security and pod security policies

**Expected Behavior**:
- Network traffic restricted per policy
- MPI communication still functional
- Security contexts enforced
- No privilege escalation

## Running Tests

### Prerequisites
```bash
# Ensure kubectl/oc is configured
kubectl config current-context

# Verify namespace exists
kubectl get namespace default

# Check for required resources (GPUs for gpu-mpijob.yaml)
kubectl describe nodes | grep nvidia.com/gpu
```

### Run Individual Tests
```bash
# Simple test
kubectl apply -f simple-mpijob.yaml
kubectl get mpijob simple-mpijob-test -w

# Check status
kubectl describe mpijob simple-mpijob-test

# View logs
kubectl logs -f simple-mpijob-test-launcher
kubectl logs -f simple-mpijob-test-worker-0

# Cleanup
kubectl delete mpijob simple-mpijob-test
```

### Run All Tests
```bash
# Apply all manifests
kubectl apply -f .

# Wait for completion
kubectl get mpijob -w

# Verify all jobs succeeded
kubectl get mpijob -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.status.conditions[-1:].type}{"\n"}{end}'

# Cleanup all
kubectl delete mpijob --all
```

## Validation Criteria

### Success Criteria
- [ ] All pods created with correct specifications
- [ ] Launcher pod runs to completion
- [ ] All worker pods respond to MPI commands
- [ ] Job status updates correctly (Created -> Running -> Succeeded)
- [ ] No orphaned pods after job completion
- [ ] Resources cleaned up per cleanPodPolicy
- [ ] TTL cleanup works if configured

### Performance Criteria
- [ ] Job submission latency < 5 seconds
- [ ] Pod initialization < 30 seconds
- [ ] MPI communication latency acceptable
- [ ] GPU utilization (if applicable) > 80%
- [ ] No resource contention or deadlocks

### Failure Scenarios to Test
- [ ] Invalid configuration rejected by validation webhook
- [ ] Resource quota exceeded handled gracefully
- [ ] Worker pod failure causes job failure
- [ ] Network policy blocks disallowed traffic
- [ ] Security context violations prevented

## Test Automation

These manifests can be integrated into CI/CD pipelines:

```yaml
# Example GitHub Actions workflow
name: MPIJob Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Set up k8s cluster
        run: kind create cluster
      - name: Install KubeFlow Training Operator
        run: kubectl apply -k manifests/kubeflow-training-operator
      - name: Run MPIJob tests
        run: |
          for manifest in test/manifests/mpijob/*.yaml; do
            kubectl apply -f $manifest
            kubectl wait --for=condition=Succeeded mpijob --all --timeout=600s
            kubectl delete -f $manifest
          done
```

## Debugging

### Common Issues

**Pods stuck in Pending**:
```bash
kubectl describe mpijob <name>
kubectl get events --sort-by='.lastTimestamp'
kubectl describe pod <pod-name>
```

**MPI communication failures**:
```bash
kubectl logs <launcher-pod>
kubectl exec -it <worker-pod> -- mpirun --version
kubectl exec -it <worker-pod> -- ssh worker-0 hostname
```

**Network policy issues**:
```bash
kubectl describe networkpolicy mpijob-restricted-network-policy
kubectl exec -it <launcher-pod> -- nc -zv worker-0 22
```

## Contributing

When adding new test manifests:
1. Follow the naming convention: `<feature>-mpijob.yaml`
2. Add comprehensive documentation in this README
3. Include expected behavior and validation criteria
4. Test in a clean environment before committing
5. Update CI/CD pipelines if needed
