import React from 'react';
import {
  Brand,
  Button,
  ButtonVariant,
  Dropdown,
  DropdownToggle,
  KebabToggle,
  PageHeader,
  PageHeaderTools,
  PageHeaderToolsGroup,
  PageHeaderToolsItem,
} from '@patternfly/react-core';
import { BellIcon, CogIcon, HelpIcon } from '@patternfly/react-icons';

interface AppHeaderProps {
  onNavToggle: () => void;
}

export const AppHeader: React.FC<AppHeaderProps> = ({ onNavToggle }) => {
  const [isKebabDropdownOpen, setIsKebabDropdownOpen] = React.useState(false);
  const [isDropdownOpen, setIsDropdownOpen] = React.useState(false);

  const onKebabDropdownToggle = () => {
    setIsKebabDropdownOpen(!isKebabDropdownOpen);
  };

  const onDropdownToggle = () => {
    setIsDropdownOpen(!isDropdownOpen);
  };

  const kebabDropdownItems = [
    <Button variant="link">Settings</Button>,
    <Button variant="link">Help</Button>,
  ];

  const userDropdownItems = [
    <Button variant="link">Profile</Button>,
    <Button variant="link">Preferences</Button>,
    <Button variant="link">Logout</Button>,
  ];

  const headerTools = (
    <PageHeaderTools>
      <PageHeaderToolsGroup
        visibility={{
          default: 'hidden',
          lg: 'visible',
        }}
      >
        <PageHeaderToolsItem>
          <Button
            aria-label="Notifications"
            variant={ButtonVariant.plain}
            icon={<BellIcon />}
          />
        </PageHeaderToolsItem>
        <PageHeaderToolsItem>
          <Button
            aria-label="Settings"
            variant={ButtonVariant.plain}
            icon={<CogIcon />}
          />
        </PageHeaderToolsItem>
        <PageHeaderToolsItem>
          <Button
            aria-label="Help"
            variant={ButtonVariant.plain}
            icon={<HelpIcon />}
          />
        </PageHeaderToolsItem>
      </PageHeaderToolsGroup>
      <PageHeaderToolsGroup>
        <PageHeaderToolsItem
          visibility={{
            lg: 'hidden',
          }}
        >
          <Dropdown
            isPlain
            position="right"
            toggle={<KebabToggle onToggle={onKebabDropdownToggle} />}
            isOpen={isKebabDropdownOpen}
            dropdownItems={kebabDropdownItems}
          />
        </PageHeaderToolsItem>
        <PageHeaderToolsItem
          visibility={{
            default: 'hidden',
            md: 'visible',
          }}
        >
          <Dropdown
            isPlain
            position="right"
            onSelect={() => setIsDropdownOpen(false)}
            isOpen={isDropdownOpen}
            toggle={<DropdownToggle onToggle={onDropdownToggle}>User Name</DropdownToggle>}
            dropdownItems={userDropdownItems}
          />
        </PageHeaderToolsItem>
      </PageHeaderToolsGroup>
    </PageHeaderTools>
  );

  return (
    <PageHeader
      logo={<Brand src="#" alt="OpenShift AI Logo" />}
      headerTools={headerTools}
      showNavToggle
      onNavToggle={onNavToggle}
    />
  );
};