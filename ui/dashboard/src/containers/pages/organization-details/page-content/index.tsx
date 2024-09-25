import { useMemo } from 'react';
import { ProjectContent } from '../projects';
import { Setting } from '../settings';
import { User } from '../users';

export type ContentDetailsProps = {
  organizationId?: string;
};

export type TargetTab = {
  targetTab: string;
};

type PageDetailsProps = ContentDetailsProps & TargetTab;

export const OrganizationDetailsContent = ({
  targetTab,
  organizationId
}: PageDetailsProps) => {
  const render = useMemo(() => {
    switch (targetTab) {
      case 'projects':
        return <ProjectContent organizationId={organizationId} />;
      case 'users':
        return <User organizationId={organizationId} />;
      case 'settings':
        return <Setting organizationId={organizationId} />;
    }
  }, [targetTab]);

  return <div className="p-6">{render}</div>;
};
