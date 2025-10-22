# OpenShift AI Frontend for MPIJob Support

React-based frontend for managing MPIJob training jobs in OpenShift AI.

## Features

- Create, view, and manage distributed training jobs
- Support for MPIJobs using MPI protocol
- Job monitoring with real-time status updates
- Worker topology visualization
- Resource usage monitoring
- Log viewer for launcher and worker pods
- Responsive design using PatternFly components

## Getting Started

### Prerequisites

- Node.js 16 or later
- npm 7 or later

### Installation

```bash
# Install dependencies
npm install
```

### Development

```bash
# Start the development server
npm start
```

This will start the development server at http://localhost:3000.

### Building for Production

```bash
# Build for production
npm run build
```

The build artifacts will be stored in the `build/` directory.

## Project Structure

```
src/
├── components/         # React components
│   ├── training/       # Training-specific components
│   │   ├── MpiJobForm.tsx     # MPIJob creation form
│   │   ├── JobList.tsx        # Job listing page
│   │   ├── JobDetailView.tsx  # Job details page
│   │   └── ...
│   └── ...
├── services/           # API services and utilities
├── validation/         # Form validation logic
├── styles/             # CSS and styling
├── App.tsx             # Main application component
└── index.tsx           # Entry point
```

## UI Components

### MpiJobForm

Form for creating MPIJobs with the following features:
- Basic configuration (name, namespace, image, command)
- Worker configuration (count, GPUs, CPU, memory)
- MPI implementation selection
- Advanced options (slots per worker, launcher resources)
- Resource validation
- Progressive disclosure pattern

### JobList

Job listing page with features:
- Filtering by status and type
- Search by name or namespace
- Pagination
- Status indicators
- Navigation to job details

### JobDetailView

Job details page with features:
- Job status and information
- Worker topology visualization
- Resource utilization charts
- Log viewer for launcher and workers
- Events timeline
- Quick actions

## Technologies

- React 18
- TypeScript
- PatternFly 4 (React components)
- React Router 6
- Formik (Form management)
- Yup (Form validation)
- Axios (API client)

## Development Guidelines

### Code Style

The project uses ESLint and Prettier for code formatting:

```bash
# Run linting
npm run lint

# Format code
npm run format
```

### Testing

```bash
# Run tests
npm test
```

### Accessibility

The UI components are designed to meet WCAG 2.1 Level AA standards, including:
- Keyboard navigation
- Screen reader compatibility
- Color contrast compliance
- Focus management

## Integration with Backend

The frontend interacts with the OpenShift AI API server, which provides:
- Job management endpoints (create, read, update, delete)
- Status updates via API polling
- Log streaming
- Resource metrics

## License

Apache License 2.0