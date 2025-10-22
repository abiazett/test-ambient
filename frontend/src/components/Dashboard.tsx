import React from 'react';
import {
  Card,
  CardTitle,
  CardBody,
  Gallery,
  GalleryItem,
  Title,
  EmptyState,
  EmptyStateIcon,
  EmptyStateBody,
  Button,
  Bullseye,
} from '@patternfly/react-core';
import { CubesIcon } from '@patternfly/react-icons';
import { useNavigate } from 'react-router-dom';

export const Dashboard: React.FC = () => {
  const navigate = useNavigate();

  const navigateToCreateJob = () => {
    navigate('/training/jobs/create');
  };

  return (
    <>
      <Title headingLevel="h1" size="2xl" className="pf-u-mb-md">
        Dashboard
      </Title>
      <Gallery hasGutter>
        <GalleryItem>
          <Card>
            <CardTitle>Training Jobs</CardTitle>
            <CardBody>
              <EmptyState>
                <EmptyStateIcon icon={CubesIcon} />
                <Title headingLevel="h4" size="lg">
                  No Training Jobs
                </Title>
                <EmptyStateBody>
                  No training jobs have been created yet. Create your first job to get started.
                </EmptyStateBody>
                <Button variant="primary" onClick={navigateToCreateJob}>
                  Create Training Job
                </Button>
              </EmptyState>
            </CardBody>
          </Card>
        </GalleryItem>
        <GalleryItem>
          <Card>
            <CardTitle>Resources</CardTitle>
            <CardBody>
              <Bullseye>
                <div>
                  <div>CPU: 0/16 cores</div>
                  <div>Memory: 0/64 GB</div>
                  <div>GPU: 0/8</div>
                </div>
              </Bullseye>
            </CardBody>
          </Card>
        </GalleryItem>
        <GalleryItem>
          <Card>
            <CardTitle>Quick Links</CardTitle>
            <CardBody>
              <ul>
                <li>
                  <Button variant="link" component="a" isInline onClick={() => navigate('/training/jobs/create')}>
                    Create Training Job
                  </Button>
                </li>
                <li>
                  <Button variant="link" component="a" isInline onClick={() => navigate('/training/jobs')}>
                    View Training Jobs
                  </Button>
                </li>
                <li>
                  <Button variant="link" component="a" isInline onClick={() => navigate('/models/registry')}>
                    Model Registry
                  </Button>
                </li>
              </ul>
            </CardBody>
          </Card>
        </GalleryItem>
      </Gallery>
    </>
  );
};