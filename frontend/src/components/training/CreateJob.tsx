import React, { useState } from 'react';
import {
  Title,
  Tabs,
  Tab,
  TabTitleText,
  PageSection,
  PageSectionVariants,
  Card,
  CardBody,
  Text,
  Button,
  Grid,
  GridItem,
} from '@patternfly/react-core';
import { MpiJobForm } from './MpiJobForm';

export const CreateJob: React.FC = () => {
  const [activeTabKey, setActiveTabKey] = useState<string | number>(1);

  const handleTabClick = (
    event: React.MouseEvent<any> | React.KeyboardEvent | MouseEvent,
    tabIndex: string | number
  ) => {
    setActiveTabKey(tabIndex);
  };

  return (
    <>
      <PageSection variant={PageSectionVariants.light}>
        <Title headingLevel="h1" size="2xl">
          Create Training Job
        </Title>
        <Text component="p">
          Select the type of training job to create, then configure the job parameters.
        </Text>
      </PageSection>

      <PageSection>
        <Tabs activeKey={activeTabKey} onSelect={handleTabClick} isBox>
          <Tab eventKey={0} title={<TabTitleText>TFJob</TabTitleText>}>
            <Card>
              <CardBody>
                <Title headingLevel="h2" size="xl">
                  TensorFlow Training Job (TFJob)
                </Title>
                <Text component="p">
                  Create a TensorFlow training job using distributed TensorFlow.
                </Text>
                <Button variant="primary" isDisabled>
                  Not yet implemented
                </Button>
              </CardBody>
            </Card>
          </Tab>
          <Tab eventKey={1} title={<TabTitleText>MPIJob</TabTitleText>}>
            <Grid hasGutter>
              <GridItem span={8}>
                <MpiJobForm />
              </GridItem>
              <GridItem span={4}>
                <Card>
                  <CardBody>
                    <Title headingLevel="h3" size="xl">
                      MPI Job Information
                    </Title>
                    <Text component="p">
                      MPIJob is used for distributed training using the Message Passing Interface (MPI) protocol.
                      This is particularly useful for:
                    </Text>
                    <ul>
                      <li>Distributed deep learning with Horovod</li>
                      <li>Training with multiple GPUs across multiple nodes</li>
                      <li>Scaling your training to large clusters</li>
                    </ul>
                    <Title headingLevel="h4" size="lg">
                      When to use MPIJob?
                    </Title>
                    <Text component="p">Use MPIJob when:</Text>
                    <ul>
                      <li>You need to use Horovod for distributed training</li>
                      <li>You need tight coordination between workers</li>
                      <li>You require all-reduce operations for gradient synchronization</li>
                      <li>You want to use NCCL for GPU-to-GPU communication</li>
                    </ul>
                  </CardBody>
                </Card>
              </GridItem>
            </Grid>
          </Tab>
          <Tab eventKey={2} title={<TabTitleText>PyTorchJob</TabTitleText>}>
            <Card>
              <CardBody>
                <Title headingLevel="h2" size="xl">
                  PyTorch Training Job (PyTorchJob)
                </Title>
                <Text component="p">
                  Create a PyTorch training job using distributed PyTorch.
                </Text>
                <Button variant="primary" isDisabled>
                  Not yet implemented
                </Button>
              </CardBody>
            </Card>
          </Tab>
        </Tabs>
      </PageSection>
    </>
  );
};