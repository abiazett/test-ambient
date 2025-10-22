import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Button,
  Toolbar,
  ToolbarContent,
  ToolbarItem,
  ToolbarGroup,
  InputGroup,
  TextInput,
  Select,
  SelectOption,
  SelectVariant,
  PageSection,
  Title,
  EmptyState,
  EmptyStateIcon,
  EmptyStateBody,
  Bullseye,
  Pagination,
} from '@patternfly/react-core';
import { Table, TableHeader, TableBody } from '@patternfly/react-table';
import { SearchIcon } from '@patternfly/react-icons';
import { CubesIcon } from '@patternfly/react-icons';

// Mock job data
const mockJobs = [
  {
    id: 1,
    name: 'mnist-training',
    namespace: 'default',
    status: 'Running',
    type: 'MPIJob',
    workers: '4/4',
    gpu: '8',
    age: '10m',
    image: 'example/horovod-mnist:latest',
  },
  {
    id: 2,
    name: 'bert-finetuning',
    namespace: 'ai-models',
    status: 'Succeeded',
    type: 'MPIJob',
    workers: '8/8',
    gpu: '16',
    age: '2h',
    image: 'example/bert-finetuning:latest',
  },
  {
    id: 3,
    name: 'resnet-training',
    namespace: 'vision-models',
    status: 'Failed',
    type: 'MPIJob',
    workers: '2/4',
    gpu: '4',
    age: '30m',
    image: 'example/vision-training:v1',
  },
];

export const JobList: React.FC = () => {
  const navigate = useNavigate();
  const [searchValue, setSearchValue] = useState('');
  const [statusIsExpanded, setStatusIsExpanded] = useState(false);
  const [selectedStatus, setSelectedStatus] = useState<string>('All');
  const [typeIsExpanded, setTypeIsExpanded] = useState(false);
  const [selectedType, setSelectedType] = useState<string>('All');
  const [page, setPage] = useState(1);
  const [perPage, setPerPage] = useState(10);

  // For the sake of this demonstration, let's simulate an empty state
  const [hasJobs] = useState(true);

  const onSearchChange = (value: string) => {
    setSearchValue(value);
  };

  const onStatusToggle = (isExpanded: boolean) => {
    setStatusIsExpanded(isExpanded);
  };

  const onStatusSelect = (event: React.MouseEvent | React.ChangeEvent, selection: string) => {
    setSelectedStatus(selection);
    setStatusIsExpanded(false);
  };

  const onTypeToggle = (isExpanded: boolean) => {
    setTypeIsExpanded(isExpanded);
  };

  const onTypeSelect = (event: React.MouseEvent | React.ChangeEvent, selection: string) => {
    setSelectedType(selection);
    setTypeIsExpanded(false);
  };

  const onSetPage = (_event: React.MouseEvent | React.KeyboardEvent, newPage: number) => {
    setPage(newPage);
  };

  const onPerPageSelect = (
    _event: React.MouseEvent | React.KeyboardEvent,
    newPerPage: number,
    newPage: number
  ) => {
    setPerPage(newPerPage);
    setPage(newPage);
  };

  const handleRowClick = (job: any) => {
    navigate(`/training/jobs/${job.name}`);
  };

  const createJob = () => {
    navigate('/training/jobs/create');
  };

  // Filter jobs based on search and filters
  const filteredJobs = mockJobs.filter((job) => {
    const matchesSearch =
      searchValue === '' ||
      job.name.toLowerCase().includes(searchValue.toLowerCase()) ||
      job.namespace.toLowerCase().includes(searchValue.toLowerCase());

    const matchesStatus = selectedStatus === 'All' || job.status === selectedStatus;
    const matchesType = selectedType === 'All' || job.type === selectedType;

    return matchesSearch && matchesStatus && matchesType;
  });

  // Calculate pagination
  const offset = (page - 1) * perPage;
  const paginatedJobs = filteredJobs.slice(offset, offset + perPage);

  const statusOptions = ['All', 'Pending', 'Running', 'Succeeded', 'Failed'];
  const typeOptions = ['All', 'MPIJob', 'TFJob', 'PyTorchJob'];

  const columns = [
    'Name',
    'Namespace',
    'Status',
    'Type',
    'Workers',
    'GPUs',
    'Age',
  ];

  const rows = paginatedJobs.map((job) => ({
    cells: [
      {
        title: (
          <Button variant="link" isInline onClick={() => handleRowClick(job)}>
            {job.name}
          </Button>
        ),
      },
      job.namespace,
      {
        title: (
          <span
            style={{
              color:
                job.status === 'Running'
                  ? 'var(--pf-global--primary-color--100)'
                  : job.status === 'Succeeded'
                  ? 'var(--pf-global--success-color--100)'
                  : job.status === 'Failed'
                  ? 'var(--pf-global--danger-color--100)'
                  : 'inherit',
            }}
          >
            {job.status}
          </span>
        ),
      },
      job.type,
      job.workers,
      job.gpu,
      job.age,
    ],
  }));

  if (!hasJobs) {
    return (
      <PageSection>
        <Title headingLevel="h1" size="2xl">
          Training Jobs
        </Title>
        <EmptyState>
          <EmptyStateIcon icon={CubesIcon} />
          <Title headingLevel="h2" size="lg">
            No training jobs found
          </Title>
          <EmptyStateBody>
            You haven't created any training jobs yet. Create your first job to get started.
          </EmptyStateBody>
          <Button variant="primary" onClick={createJob}>
            Create Training Job
          </Button>
        </EmptyState>
      </PageSection>
    );
  }

  return (
    <PageSection>
      <Title headingLevel="h1" size="2xl">
        Training Jobs
      </Title>
      <Toolbar id="job-list-toolbar">
        <ToolbarContent>
          <ToolbarGroup variant="filter-group">
            <ToolbarItem>
              <InputGroup>
                <TextInput
                  name="search"
                  id="search"
                  type="search"
                  aria-label="search jobs"
                  placeholder="Search"
                  value={searchValue}
                  onChange={onSearchChange}
                />
                <Button variant="control" aria-label="search button">
                  <SearchIcon />
                </Button>
              </InputGroup>
            </ToolbarItem>
            <ToolbarItem>
              <Select
                variant={SelectVariant.single}
                aria-label="Select Status"
                onToggle={onStatusToggle}
                onSelect={onStatusSelect}
                selections={selectedStatus}
                isOpen={statusIsExpanded}
                placeholderText="Status"
              >
                {statusOptions.map((option) => (
                  <SelectOption key={option} value={option} />
                ))}
              </Select>
            </ToolbarItem>
            <ToolbarItem>
              <Select
                variant={SelectVariant.single}
                aria-label="Select Type"
                onToggle={onTypeToggle}
                onSelect={onTypeSelect}
                selections={selectedType}
                isOpen={typeIsExpanded}
                placeholderText="Job Type"
              >
                {typeOptions.map((option) => (
                  <SelectOption key={option} value={option} />
                ))}
              </Select>
            </ToolbarItem>
          </ToolbarGroup>
          <ToolbarGroup variant="button-group">
            <ToolbarItem>
              <Button variant="primary" onClick={createJob}>
                Create Job
              </Button>
            </ToolbarItem>
          </ToolbarGroup>
          <ToolbarItem variant="pagination">
            <Pagination
              itemCount={filteredJobs.length}
              perPage={perPage}
              page={page}
              onSetPage={onSetPage}
              onPerPageSelect={onPerPageSelect}
              widgetId="job-list-pagination-top"
            />
          </ToolbarItem>
        </ToolbarContent>
      </Toolbar>

      <Table aria-label="Job List" cells={columns} rows={rows}>
        <TableHeader />
        <TableBody />
      </Table>

      <Bullseye>
        <Pagination
          itemCount={filteredJobs.length}
          perPage={perPage}
          page={page}
          onSetPage={onSetPage}
          onPerPageSelect={onPerPageSelect}
          widgetId="job-list-pagination-bottom"
        />
      </Bullseye>
    </PageSection>
  );
};