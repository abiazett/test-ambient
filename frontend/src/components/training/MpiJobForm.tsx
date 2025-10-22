import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Form,
  FormGroup,
  TextInput,
  Select,
  SelectOption,
  SelectVariant,
  Switch,
  Button,
  ActionGroup,
  NumberInput,
  TextArea,
  Card,
  CardBody,
  CardTitle,
  Divider,
  Alert,
  FormHelperText,
  HelperText,
  HelperTextItem,
} from '@patternfly/react-core';
import { useFormik } from 'formik';
import * as Yup from 'yup';

// Mock MPIJob creation function
const createMPIJob = async (jobData: any): Promise<boolean> => {
  console.log('Creating MPIJob with data:', jobData);
  // In a real implementation, this would make an API call
  return new Promise((resolve) => {
    setTimeout(() => resolve(true), 1000);
  });
};

const MPIImplementations = ['OpenMPI', 'IntelMPI', 'MPICH'] as const;
type MPIImplementationType = typeof MPIImplementations[number];

export interface MPIJobFormValues {
  name: string;
  namespace: string;
  image: string;
  command: string;
  workerCount: number;
  gpuCount: number;
  cpuRequest: string;
  memoryRequest: string;
  mpiImplementation: MPIImplementationType;
  advancedOptions: boolean;
  slotCount: number;
  launcherCpu: string;
  launcherMemory: string;
}

export const MpiJobForm: React.FC = () => {
  const navigate = useNavigate();
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [isNamespaceSelectOpen, setIsNamespaceSelectOpen] = useState(false);
  const [isMPIImplementationSelectOpen, setIsMPIImplementationSelectOpen] = useState(false);

  const initialValues: MPIJobFormValues = {
    name: '',
    namespace: 'default',
    image: '',
    command: '',
    workerCount: 4,
    gpuCount: 2,
    cpuRequest: '4',
    memoryRequest: '16Gi',
    mpiImplementation: 'OpenMPI',
    advancedOptions: false,
    slotCount: 1,
    launcherCpu: '1',
    launcherMemory: '1Gi',
  };

  const validationSchema = Yup.object({
    name: Yup.string()
      .matches(/^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/, 'Must consist of lowercase alphanumeric characters or \'-\'')
      .max(63, 'Cannot be longer than 63 characters')
      .required('Required'),
    namespace: Yup.string().required('Required'),
    image: Yup.string().required('Container image is required'),
    workerCount: Yup.number().min(1, 'At least 1 worker is required').required('Required'),
    gpuCount: Yup.number().min(0, 'Cannot be negative').required('Required'),
    cpuRequest: Yup.string().required('CPU request is required'),
    memoryRequest: Yup.string().required('Memory request is required'),
  });

  const formik = useFormik({
    initialValues,
    validationSchema,
    onSubmit: async (values) => {
      setIsSubmitting(true);
      setError(null);

      try {
        // Calculate total resources
        const totalGPUs = values.workerCount * values.gpuCount;
        const totalCPUs = parseInt(values.cpuRequest) * values.workerCount;
        const totalMemory = parseInt(values.memoryRequest) * values.workerCount;

        console.log(`Total resources: ${totalGPUs} GPUs, ${totalCPUs} CPUs, ${totalMemory} memory`);

        const success = await createMPIJob({
          ...values,
          totalResources: {
            gpu: totalGPUs,
            cpu: totalCPUs,
            memory: totalMemory,
          }
        });

        if (success) {
          navigate('/training/jobs');
        } else {
          setError('Failed to create MPIJob. Please try again.');
        }
      } catch (err) {
        setError(`Error creating MPIJob: ${err}`);
      } finally {
        setIsSubmitting(false);
      }
    },
  });

  const onNamespaceToggle = () => {
    setIsNamespaceSelectOpen(!isNamespaceSelectOpen);
  };

  const onNamespaceSelect = (event: React.MouseEvent | React.ChangeEvent, value: string) => {
    formik.setFieldValue('namespace', value);
    setIsNamespaceSelectOpen(false);
  };

  const onMPIImplementationToggle = () => {
    setIsMPIImplementationSelectOpen(!isMPIImplementationSelectOpen);
  };

  const onMPIImplementationSelect = (event: React.MouseEvent | React.ChangeEvent, value: string) => {
    formik.setFieldValue('mpiImplementation', value);
    setIsMPIImplementationSelectOpen(false);
  };

  const onWorkerCountChange = (value: number) => {
    formik.setFieldValue('workerCount', value);
  };

  const onGPUCountChange = (value: number) => {
    formik.setFieldValue('gpuCount', value);
  };

  const onSlotCountChange = (value: number) => {
    formik.setFieldValue('slotCount', value);
  };

  return (
    <>
      <Card>
        <CardTitle>Create MPIJob</CardTitle>
        <CardBody>
          {error && <Alert variant="danger" title="Error" isInline>{error}</Alert>}
          <Form onSubmit={formik.handleSubmit}>
            {/* Basic Configuration */}
            <FormGroup
              label="Job Name"
              isRequired
              fieldId="name"
              helperText={
                <FormHelperText>
                  <HelperText>
                    <HelperTextItem>
                      Must consist of lowercase alphanumeric characters or '-'
                    </HelperTextItem>
                  </HelperText>
                </FormHelperText>
              }
              validated={formik.touched.name && formik.errors.name ? 'error' : 'default'}
            >
              <TextInput
                isRequired
                type="text"
                id="name"
                name="name"
                value={formik.values.name}
                onChange={formik.handleChange}
                onBlur={formik.handleBlur}
                validated={formik.touched.name && formik.errors.name ? 'error' : 'default'}
              />
              {formik.touched.name && formik.errors.name && (
                <FormHelperText>
                  <HelperText>
                    <HelperTextItem variant="error">{formik.errors.name}</HelperTextItem>
                  </HelperText>
                </FormHelperText>
              )}
            </FormGroup>

            <FormGroup label="Namespace" isRequired fieldId="namespace">
              <Select
                variant={SelectVariant.single}
                onToggle={onNamespaceToggle}
                onSelect={onNamespaceSelect}
                selections={formik.values.namespace}
                isOpen={isNamespaceSelectOpen}
                aria-labelledby="namespace"
              >
                <SelectOption value="default" />
                <SelectOption value="openshift-ai" />
                <SelectOption value="dev" />
                <SelectOption value="production" />
              </Select>
            </FormGroup>

            <FormGroup
              label="Container Image"
              isRequired
              fieldId="image"
              helperText="Container image to use for training (e.g., myregistry.com/training:horovod-latest)"
              validated={formik.touched.image && formik.errors.image ? 'error' : 'default'}
            >
              <TextInput
                isRequired
                type="text"
                id="image"
                name="image"
                value={formik.values.image}
                onChange={formik.handleChange}
                onBlur={formik.handleBlur}
                validated={formik.touched.image && formik.errors.image ? 'error' : 'default'}
              />
              {formik.touched.image && formik.errors.image && (
                <FormHelperText>
                  <HelperText>
                    <HelperTextItem variant="error">{formik.errors.image}</HelperTextItem>
                  </HelperText>
                </FormHelperText>
              )}
            </FormGroup>

            <FormGroup
              label="Command"
              fieldId="command"
              helperText="Command to run in the container (e.g., python /train.py --epochs 10)"
            >
              <TextArea
                id="command"
                name="command"
                value={formik.values.command}
                onChange={formik.handleChange}
                onBlur={formik.handleBlur}
              />
            </FormGroup>

            <Divider />

            {/* Worker Configuration */}
            <FormGroup label="Number of Workers" isRequired fieldId="workerCount">
              <NumberInput
                value={formik.values.workerCount}
                min={1}
                max={100}
                onMinus={() => onWorkerCountChange(Math.max(1, formik.values.workerCount - 1))}
                onPlus={() => onWorkerCountChange(formik.values.workerCount + 1)}
                onChange={(event) => {
                  const value = parseInt(event.currentTarget.value, 10);
                  if (!isNaN(value)) {
                    onWorkerCountChange(value);
                  }
                }}
                inputName="workerCount"
                inputAriaLabel="Number of workers"
              />
            </FormGroup>

            <FormGroup label="GPUs per Worker" isRequired fieldId="gpuCount">
              <NumberInput
                value={formik.values.gpuCount}
                min={0}
                max={8}
                onMinus={() => onGPUCountChange(Math.max(0, formik.values.gpuCount - 1))}
                onPlus={() => onGPUCountChange(formik.values.gpuCount + 1)}
                onChange={(event) => {
                  const value = parseInt(event.currentTarget.value, 10);
                  if (!isNaN(value)) {
                    onGPUCountChange(value);
                  }
                }}
                inputName="gpuCount"
                inputAriaLabel="GPUs per worker"
              />
            </FormGroup>

            <FormGroup
              label="CPU per Worker"
              isRequired
              fieldId="cpuRequest"
              helperText="CPU cores per worker (e.g., 1, 2, 4)"
            >
              <TextInput
                isRequired
                type="text"
                id="cpuRequest"
                name="cpuRequest"
                value={formik.values.cpuRequest}
                onChange={formik.handleChange}
                onBlur={formik.handleBlur}
              />
            </FormGroup>

            <FormGroup
              label="Memory per Worker"
              isRequired
              fieldId="memoryRequest"
              helperText="Memory per worker (e.g., 4Gi, 8Gi, 16Gi)"
            >
              <TextInput
                isRequired
                type="text"
                id="memoryRequest"
                name="memoryRequest"
                value={formik.values.memoryRequest}
                onChange={formik.handleChange}
                onBlur={formik.handleBlur}
              />
            </FormGroup>

            <FormGroup label="MPI Implementation" fieldId="mpiImplementation">
              <Select
                variant={SelectVariant.single}
                onToggle={onMPIImplementationToggle}
                onSelect={onMPIImplementationSelect}
                selections={formik.values.mpiImplementation}
                isOpen={isMPIImplementationSelectOpen}
                aria-labelledby="mpiImplementation"
              >
                {MPIImplementations.map((impl) => (
                  <SelectOption key={impl} value={impl} />
                ))}
              </Select>
            </FormGroup>

            <FormGroup label="Advanced Options" fieldId="advancedOptions">
              <Switch
                id="advancedOptions"
                label="Show Advanced Options"
                isChecked={formik.values.advancedOptions}
                onChange={() => formik.setFieldValue('advancedOptions', !formik.values.advancedOptions)}
              />
            </FormGroup>

            {formik.values.advancedOptions && (
              <>
                <Divider />
                <FormGroup label="Slots per Worker" fieldId="slotCount">
                  <NumberInput
                    value={formik.values.slotCount}
                    min={1}
                    max={8}
                    onMinus={() => onSlotCountChange(Math.max(1, formik.values.slotCount - 1))}
                    onPlus={() => onSlotCountChange(formik.values.slotCount + 1)}
                    onChange={(event) => {
                      const value = parseInt(event.currentTarget.value, 10);
                      if (!isNaN(value)) {
                        onSlotCountChange(value);
                      }
                    }}
                    inputName="slotCount"
                    inputAriaLabel="Slots per worker"
                  />
                </FormGroup>

                <FormGroup
                  label="Launcher CPU"
                  fieldId="launcherCpu"
                  helperText="CPU cores for launcher pod"
                >
                  <TextInput
                    type="text"
                    id="launcherCpu"
                    name="launcherCpu"
                    value={formik.values.launcherCpu}
                    onChange={formik.handleChange}
                    onBlur={formik.handleBlur}
                  />
                </FormGroup>

                <FormGroup
                  label="Launcher Memory"
                  fieldId="launcherMemory"
                  helperText="Memory for launcher pod"
                >
                  <TextInput
                    type="text"
                    id="launcherMemory"
                    name="launcherMemory"
                    value={formik.values.launcherMemory}
                    onChange={formik.handleChange}
                    onBlur={formik.handleBlur}
                  />
                </FormGroup>
              </>
            )}

            <ActionGroup>
              <Button variant="primary" type="submit" isDisabled={isSubmitting}>
                {isSubmitting ? 'Creating...' : 'Create Job'}
              </Button>
              <Button variant="secondary" onClick={() => navigate('/training/jobs')}>
                Cancel
              </Button>
              <Button variant="link">Save as Template</Button>
            </ActionGroup>
          </Form>
        </CardBody>
      </Card>
    </>
  );
};