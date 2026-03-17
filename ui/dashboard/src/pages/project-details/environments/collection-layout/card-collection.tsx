import { ReactElement } from 'react';
import { Environment } from '@types';
import { EnvironmentCard } from 'components/mobile-card/env-card';
import PageLayout from 'elements/page-layout';
import { EnvironmentActionsType } from '../types';

interface CardCollectionProps {
  data: Environment[];
  onActions: (value: Environment, type: EnvironmentActionsType) => void;
  emptyCollection?: ReactElement;
  isLoading?: boolean;
}

export const CardCollection = ({
  data,
  emptyCollection,
  isLoading,
  onActions
}: CardCollectionProps) => {
  return isLoading ? (
    <PageLayout.LoadingState className="py-10" />
  ) : (
    <div className="flex flex-col gap-3">
      {data.length
        ? data.map(project => (
            <EnvironmentCard
              key={project.id}
              onActions={onActions}
              data={project}
            />
          ))
        : emptyCollection}
    </div>
  );
};
