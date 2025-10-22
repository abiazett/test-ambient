import React from 'react';
import { Route, Routes, Navigate } from 'react-router-dom';
import { Page, PageSidebar, PageSection, PageSectionVariants } from '@patternfly/react-core';
import { AppHeader } from './components/AppHeader';
import { AppSidebar } from './components/AppSidebar';
import { Dashboard } from './components/Dashboard';
import { JobList } from './components/training/JobList';
import { CreateJob } from './components/training/CreateJob';
import { JobDetail } from './components/training/JobDetailView';

const App: React.FC = () => {
  const [isNavOpen, setIsNavOpen] = React.useState(true);

  const onNavToggle = () => {
    setIsNavOpen(!isNavOpen);
  };

  const sidebar = <PageSidebar nav={<AppSidebar />} isNavOpen={isNavOpen} />;

  return (
    <Page header={<AppHeader onNavToggle={onNavToggle} />} sidebar={sidebar}>
      <PageSection variant={PageSectionVariants.light}>
        <Routes>
          <Route path="/" element={<Navigate to="/dashboard" replace />} />
          <Route path="/dashboard" element={<Dashboard />} />
          <Route path="/training/jobs" element={<JobList />} />
          <Route path="/training/jobs/create" element={<CreateJob />} />
          <Route path="/training/jobs/:name" element={<JobDetail />} />
        </Routes>
      </PageSection>
    </Page>
  );
};

export default App;