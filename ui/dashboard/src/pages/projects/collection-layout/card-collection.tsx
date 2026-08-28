import { ReactElement } from 'react';
import { Account, Environment, Project } from '@types';
import { ProjectCard } from 'components/mobile-card/project-card';
import PageLayout from 'elements/page-layout';

interface CardCollectionProps {
  data: Project[];
  accounts: Account[];
  organizationId?: string;
  currentEnvironment: Environment;
  onActions: (item: Project) => void;
  emptyCollection?: ReactElement;
  isLoading?: boolean;
}

export const CardCollection = ({
  data,
  isLoading,
  accounts,
  emptyCollection,
  currentEnvironment,
  organizationId,
  onActions
}: CardCollectionProps) => {
  return isLoading ? (
    <PageLayout.LoadingState className="py-10" />
  ) : (
    <div className="flex flex-col gap-3">
      {data.length
        ? data.map(project => (
            <ProjectCard
              currentEnvironment={currentEnvironment}
              accounts={accounts}
              organizationId={organizationId}
              key={project.id}
              onActions={onActions}
              data={project}
            />
          ))
        : emptyCollection}
    </div>
  );
};
