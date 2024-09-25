import { useMemo } from 'react';
import { TargetTab } from 'containers/pages/organization-details';
import { Project } from '@types';
import { ProjectEnvironments } from '../environments';
import { Settings } from '../settings';

type Props = TargetTab & {
  projectId?: string;
  projectData?: Project;
};

export const ProjectDetailsContent = ({
  targetTab,
  projectId,
  projectData
}: Props) => {
  const render = useMemo(() => {
    switch (targetTab) {
      case 'environments':
        return <ProjectEnvironments projectId={projectId} />;
      case 'settings':
        return <Settings projectData={projectData} />;
    }
  }, [targetTab]);

  return <div className="p-6">{render}</div>;
};
