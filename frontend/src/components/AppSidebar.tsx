import React from 'react';
import { NavLink, useLocation } from 'react-router-dom';
import { Nav, NavList, NavItem, NavExpandable } from '@patternfly/react-core';

export const AppSidebar: React.FC = () => {
  const location = useLocation();
  const [activeGroup, setActiveGroup] = React.useState('training');
  const [activeItem, setActiveItem] = React.useState(location.pathname);

  const onSelect = (result: { groupId?: string; itemId?: string }) => {
    if (result.groupId) {
      setActiveGroup(result.groupId);
    }
    if (result.itemId) {
      setActiveItem(result.itemId);
    }
  };

  React.useEffect(() => {
    setActiveItem(location.pathname);

    // Set active group based on path
    if (location.pathname.includes('/training')) {
      setActiveGroup('training');
    } else if (location.pathname.includes('/models')) {
      setActiveGroup('models');
    } else if (location.pathname.includes('/projects')) {
      setActiveGroup('projects');
    } else if (location.pathname.includes('/settings')) {
      setActiveGroup('settings');
    }
  }, [location]);

  return (
    <Nav onSelect={onSelect} aria-label="Nav">
      <NavList>
        <NavItem itemId="/dashboard" isActive={activeItem === '/dashboard'}>
          <NavLink to="/dashboard">Dashboard</NavLink>
        </NavItem>
        <NavExpandable
          title="Training"
          groupId="training"
          isActive={activeGroup === 'training'}
          isExpanded={activeGroup === 'training'}
        >
          <NavItem itemId="/training/jobs" isActive={activeItem === '/training/jobs'}>
            <NavLink to="/training/jobs">Jobs</NavLink>
          </NavItem>
          <NavItem itemId="/training/jobs/create" isActive={activeItem === '/training/jobs/create'}>
            <NavLink to="/training/jobs/create">Create Job</NavLink>
          </NavItem>
        </NavExpandable>
        <NavExpandable
          title="Models"
          groupId="models"
          isActive={activeGroup === 'models'}
          isExpanded={activeGroup === 'models'}
        >
          <NavItem itemId="/models/registry" isActive={activeItem === '/models/registry'}>
            <NavLink to="/models/registry">Model Registry</NavLink>
          </NavItem>
          <NavItem itemId="/models/serving" isActive={activeItem === '/models/serving'}>
            <NavLink to="/models/serving">Model Serving</NavLink>
          </NavItem>
        </NavExpandable>
        <NavExpandable
          title="Projects"
          groupId="projects"
          isActive={activeGroup === 'projects'}
          isExpanded={activeGroup === 'projects'}
        >
          <NavItem itemId="/projects/list" isActive={activeItem === '/projects/list'}>
            <NavLink to="/projects/list">Project List</NavLink>
          </NavItem>
          <NavItem itemId="/projects/create" isActive={activeItem === '/projects/create'}>
            <NavLink to="/projects/create">Create Project</NavLink>
          </NavItem>
        </NavExpandable>
        <NavItem itemId="/settings" isActive={activeItem === '/settings'}>
          <NavLink to="/settings">Settings</NavLink>
        </NavItem>
      </NavList>
    </Nav>
  );
};