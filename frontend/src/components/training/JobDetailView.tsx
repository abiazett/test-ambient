import React, { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  PageSection,
  Title,
  Card,
  CardBody,
  Grid,
  GridItem,
  Button,
  Tabs,
  Tab,
  TabTitleText,
  Label,
  List,
  ListItem,
  Split,
  SplitItem,
  Spinner,
  Divider,
  Flex,
  FlexItem,
  Alert,
  Breadcrumb,
  BreadcrumbItem,
} from '@patternfly/react-core';
import { ChartDonutUtilization, ChartThemeColor } from '@patternfly/react-charts';

// Mock job data
const mockJobDetail = {
  name: 'mnist-training',
  namespace: 'default',
  status: 'Running',
  type: 'MPIJob',
  startTime: '2025-10-21T13:45:00Z',
  image: 'example/horovod-mnist:latest',
  command: 'python /workspace/train.py --epochs 10 --batch-size 64',
  workerCount: 4,
  workersRunning: 4,
  gpuCount: 2,
  totalGPUs: 8,
  cpuRequest: '4',
  memoryRequest: '16Gi',
  mpiImplementation: 'OpenMPI',
  events: [
    {
      type: 'Normal',
      reason: 'Created',
      message: 'MPIJob created, waiting for launcher pod',
      age: '10m',
    },
    {
      type: 'Normal',
      reason: 'LauncherCreated',
      message: 'Created launcher pod',
      age: '9m',
    },
    {
      type: 'Normal',
      reason: 'WorkersCreated',
      message: 'Created 4 worker pods',
      age: '9m',
    },
    {
      type: 'Normal',
      reason: 'Running',
      message: 'All pods running, starting MPI training',
      age: '8m',
    },
  ],
  workers: [
    {
      name: 'mnist-training-worker-0',
      status: 'Running',
      ready: true,
      cpuUsage: 85,
      memoryUsage: 70,
      gpuUsage: 95,
    },
    {
      name: 'mnist-training-worker-1',
      status: 'Running',
      ready: true,
      cpuUsage: 80,
      memoryUsage: 65,
      gpuUsage: 90,
    },
    {
      name: 'mnist-training-worker-2',
      status: 'Running',
      ready: true,
      cpuUsage: 82,
      memoryUsage: 68,
      gpuUsage: 92,
    },
    {
      name: 'mnist-training-worker-3',
      status: 'Running',
      ready: true,
      cpuUsage: 78,
      memoryUsage: 72,
      gpuUsage: 88,
    },
  ],
  launcher: {
    name: 'mnist-training-launcher',
    status: 'Running',
    ready: true,
  },
  logs: {
    launcher:
      '2025-10-21 13:45:05 INFO Initializing launcher\n2025-10-21 13:45:10 INFO All workers ready\n2025-10-21 13:45:15 INFO Starting MPI training\n2025-10-21 13:45:20 INFO Epoch 1/10, Loss: 0.4532, Accuracy: 0.8730\n2025-10-21 13:45:30 INFO Epoch 2/10, Loss: 0.2341, Accuracy: 0.9120\n2025-10-21 13:45:40 INFO Epoch 3/10, Loss: 0.1453, Accuracy: 0.9450',
    worker0:
      '2025-10-21 13:45:07 INFO Worker 0 initializing\n2025-10-21 13:45:12 INFO Worker 0 ready\n2025-10-21 13:45:17 INFO Worker 0 starting training\n2025-10-21 13:45:22 INFO Worker 0 processing batch 1/100\n2025-10-21 13:45:32 INFO Worker 0 processing batch 50/100\n2025-10-21 13:45:42 INFO Worker 0 processing batch 100/100',
  },
  metrics: {
    accuracy: [
      { x: 1, y: 0.8730 },
      { x: 2, y: 0.9120 },
      { x: 3, y: 0.9450 },
    ],
    loss: [
      { x: 1, y: 0.4532 },
      { x: 2, y: 0.2341 },
      { x: 3, y: 0.1453 },
    ],
  },
};

export const JobDetail: React.FC = () => {
  const { name } = useParams<{ name: string }>();
  const navigate = useNavigate();
  const [activeTabKey, setActiveTabKey] = useState<string | number>(0);
  const [loading] = useState(false);
  const [job] = useState(mockJobDetail); // In real implementation, load job by name

  const handleTabClick = (
    event: React.MouseEvent<any> | React.KeyboardEvent | MouseEvent,
    tabIndex: string | number
  ) => {
    setActiveTabKey(tabIndex);
  };

  if (loading) {
    return (
      <PageSection>
        <Spinner />
      </PageSection>
    );
  }

  if (!job) {
    return (
      <PageSection>
        <Alert variant="danger" title={`Job ${name} not found`}>
          <p>The requested job could not be found. It may have been deleted or you may not have permission to view it.</p>
          <Button variant="primary" onClick={() => navigate('/training/jobs')}>
            Back to Job List
          </Button>
        </Alert>
      </PageSection>
    );
  }

  const getStatusColor = (status: string): string => {
    switch (status) {
      case 'Running':
        return 'blue';
      case 'Succeeded':
        return 'green';
      case 'Failed':
        return 'red';
      case 'Pending':
        return 'orange';
      default:
        return 'grey';
    }
  };

  return (
    <>
      <PageSection>
        <Breadcrumb>
          <BreadcrumbItem to="/dashboard">Dashboard</BreadcrumbItem>
          <BreadcrumbItem to="/training/jobs">Training Jobs</BreadcrumbItem>
          <BreadcrumbItem isActive>{job.name}</BreadcrumbItem>
        </Breadcrumb>

        <Flex justifyContent={{ default: 'justifyContentSpaceBetween' }} alignItems={{ default: 'alignItemsCenter' }}>
          <FlexItem>
            <Title headingLevel="h1" size="2xl">
              {job.name}
            </Title>
          </FlexItem>
          <FlexItem>
            <Label color={getStatusColor(job.status)}>{job.status}</Label>
          </FlexItem>
          <FlexItem>
            <Button variant="secondary" onClick={() => navigate('/training/jobs/create')}>
              Clone and Modify
            </Button>
          </FlexItem>
          <FlexItem>
            <Button variant="danger">Delete Job</Button>
          </FlexItem>
        </Flex>
      </PageSection>

      <PageSection>
        <Grid hasGutter>
          <GridItem span={8}>
            <Card>
              <CardBody>
                <Tabs activeKey={activeTabKey} onSelect={handleTabClick}>
                  <Tab eventKey={0} title={<TabTitleText>Details</TabTitleText>}>
                    <Grid hasGutter className="pf-u-p-md">
                      <GridItem span={6}>
                        <Title headingLevel="h3" size="md">
                          Job Information
                        </Title>
                        <List>
                          <ListItem>
                            <b>Name:</b> {job.name}
                          </ListItem>
                          <ListItem>
                            <b>Namespace:</b> {job.namespace}
                          </ListItem>
                          <ListItem>
                            <b>Status:</b>{' '}
                            <Label color={getStatusColor(job.status)}>{job.status}</Label>
                          </ListItem>
                          <ListItem>
                            <b>Type:</b> {job.type}
                          </ListItem>
                          <ListItem>
                            <b>Start Time:</b> {job.startTime}
                          </ListItem>
                          <ListItem>
                            <b>MPI Implementation:</b> {job.mpiImplementation}
                          </ListItem>
                        </List>
                      </GridItem>

                      <GridItem span={6}>
                        <Title headingLevel="h3" size="md">
                          Resources
                        </Title>
                        <List>
                          <ListItem>
                            <b>Workers:</b> {job.workersRunning}/{job.workerCount} running
                          </ListItem>
                          <ListItem>
                            <b>GPUs per Worker:</b> {job.gpuCount}
                          </ListItem>
                          <ListItem>
                            <b>Total GPUs:</b> {job.totalGPUs}
                          </ListItem>
                          <ListItem>
                            <b>CPU per Worker:</b> {job.cpuRequest}
                          </ListItem>
                          <ListItem>
                            <b>Memory per Worker:</b> {job.memoryRequest}
                          </ListItem>
                        </List>
                      </GridItem>

                      <GridItem span={12}>
                        <Divider className="pf-u-my-md" />
                        <Title headingLevel="h3" size="md">
                          Container
                        </Title>
                        <List>
                          <ListItem>
                            <b>Image:</b> {job.image}
                          </ListItem>
                          <ListItem>
                            <b>Command:</b> <code>{job.command}</code>
                          </ListItem>
                        </List>
                      </GridItem>

                      <GridItem span={12}>
                        <Divider className="pf-u-my-md" />
                        <Title headingLevel="h3" size="md">
                          Events
                        </Title>
                        <table className="pf-c-table pf-m-compact">
                          <thead>
                            <tr>
                              <th>Type</th>
                              <th>Reason</th>
                              <th>Age</th>
                              <th>Message</th>
                            </tr>
                          </thead>
                          <tbody>
                            {job.events.map((event, index) => (
                              <tr key={index}>
                                <td>{event.type}</td>
                                <td>{event.reason}</td>
                                <td>{event.age}</td>
                                <td>{event.message}</td>
                              </tr>
                            ))}
                          </tbody>
                        </table>
                      </GridItem>
                    </Grid>
                  </Tab>

                  <Tab eventKey={1} title={<TabTitleText>Worker Topology</TabTitleText>}>
                    <Grid hasGutter className="pf-u-p-md">
                      <GridItem span={12}>
                        <Title headingLevel="h3" size="md">
                          Worker Pods
                        </Title>
                        <table className="pf-c-table pf-m-compact">
                          <thead>
                            <tr>
                              <th>Name</th>
                              <th>Status</th>
                              <th>Ready</th>
                              <th>CPU Usage</th>
                              <th>Memory Usage</th>
                              <th>GPU Usage</th>
                            </tr>
                          </thead>
                          <tbody>
                            {job.workers.map((worker, index) => (
                              <tr key={index}>
                                <td>{worker.name}</td>
                                <td>
                                  <Label color={getStatusColor(worker.status)}>
                                    {worker.status}
                                  </Label>
                                </td>
                                <td>{worker.ready ? 'Yes' : 'No'}</td>
                                <td>
                                  <Split hasGutter>
                                    <SplitItem>
                                      <ChartDonutUtilization
                                        ariaDesc="CPU usage"
                                        data={{ x: 'CPU', y: worker.cpuUsage }}
                                        labels={({ datum }) => (datum.x ? `${datum.x}: ${datum.y}%` : null)}
                                        size={50}
                                        themeColor={ChartThemeColor.blue}
                                      />
                                    </SplitItem>
                                    <SplitItem>{worker.cpuUsage}%</SplitItem>
                                  </Split>
                                </td>
                                <td>
                                  <Split hasGutter>
                                    <SplitItem>
                                      <ChartDonutUtilization
                                        ariaDesc="Memory usage"
                                        data={{ x: 'Memory', y: worker.memoryUsage }}
                                        labels={({ datum }) => (datum.x ? `${datum.x}: ${datum.y}%` : null)}
                                        size={50}
                                        themeColor={ChartThemeColor.green}
                                      />
                                    </SplitItem>
                                    <SplitItem>{worker.memoryUsage}%</SplitItem>
                                  </Split>
                                </td>
                                <td>
                                  <Split hasGutter>
                                    <SplitItem>
                                      <ChartDonutUtilization
                                        ariaDesc="GPU usage"
                                        data={{ x: 'GPU', y: worker.gpuUsage }}
                                        labels={({ datum }) => (datum.x ? `${datum.x}: ${datum.y}%` : null)}
                                        size={50}
                                        themeColor={ChartThemeColor.gold}
                                      />
                                    </SplitItem>
                                    <SplitItem>{worker.gpuUsage}%</SplitItem>
                                  </Split>
                                </td>
                              </tr>
                            ))}
                          </tbody>
                        </table>
                      </GridItem>

                      <GridItem span={12}>
                        <Divider className="pf-u-my-md" />
                        <Title headingLevel="h3" size="md">
                          Launcher Pod
                        </Title>
                        <table className="pf-c-table pf-m-compact">
                          <thead>
                            <tr>
                              <th>Name</th>
                              <th>Status</th>
                              <th>Ready</th>
                            </tr>
                          </thead>
                          <tbody>
                            <tr>
                              <td>{job.launcher.name}</td>
                              <td>
                                <Label color={getStatusColor(job.launcher.status)}>
                                  {job.launcher.status}
                                </Label>
                              </td>
                              <td>{job.launcher.ready ? 'Yes' : 'No'}</td>
                            </tr>
                          </tbody>
                        </table>
                      </GridItem>
                    </Grid>
                  </Tab>

                  <Tab eventKey={2} title={<TabTitleText>Logs</TabTitleText>}>
                    <Grid hasGutter className="pf-u-p-md">
                      <GridItem span={12}>
                        <Title headingLevel="h3" size="md">
                          Launcher Logs
                        </Title>
                        <pre
                          style={{
                            backgroundColor: '#f0f0f0',
                            padding: '10px',
                            border: '1px solid #ccc',
                            height: '200px',
                            overflowY: 'auto',
                          }}
                        >
                          {job.logs.launcher}
                        </pre>
                      </GridItem>

                      <GridItem span={12}>
                        <Title headingLevel="h3" size="md">
                          Worker 0 Logs
                        </Title>
                        <pre
                          style={{
                            backgroundColor: '#f0f0f0',
                            padding: '10px',
                            border: '1px solid #ccc',
                            height: '200px',
                            overflowY: 'auto',
                          }}
                        >
                          {job.logs.worker0}
                        </pre>
                      </GridItem>
                    </Grid>
                  </Tab>
                </Tabs>
              </CardBody>
            </Card>
          </GridItem>

          <GridItem span={4}>
            <Card>
              <CardBody>
                <Title headingLevel="h3" size="lg">
                  Job Status
                </Title>
                <ChartDonutUtilization
                  ariaDesc="Workers running"
                  data={{ x: 'Workers', y: (job.workersRunning / job.workerCount) * 100 }}
                  labels={({ datum }) => (datum.x ? `${datum.x}: ${job.workersRunning}/${job.workerCount}` : null)}
                  themeColor={ChartThemeColor.blue}
                />

                <Divider className="pf-u-my-md" />

                <Title headingLevel="h3" size="lg">
                  Quick Actions
                </Title>
                <Button variant="primary" isBlock className="pf-u-mb-md">
                  View YAML Definition
                </Button>
                <Button variant="secondary" isBlock className="pf-u-mb-md">
                  Scale Workers
                </Button>
                <Button variant="link" isBlock>
                  View Raw Metrics
                </Button>
              </CardBody>
            </Card>
          </GridItem>
        </Grid>
      </PageSection>
    </>
  );
};