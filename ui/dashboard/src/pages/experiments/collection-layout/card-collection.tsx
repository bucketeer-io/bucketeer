import { ReactElement } from 'react';
import { Experiment } from '@types';
import { ExperimentCard } from 'components/mobile-card/experiment-card';
import PageLayout from 'elements/page-layout';
import { ExperimentActionsType } from '../types';

interface CardCollectionProps {
  data: Experiment[];
  onActions: (item: Experiment, type: ExperimentActionsType) => void;
  emptyCollection?: ReactElement;
  isLoading?: boolean;
}

export const CardCollection = ({
  data,
  isLoading,
  emptyCollection,
  onActions
}: CardCollectionProps) => {
  return isLoading ? (
    <PageLayout.LoadingState className="py-10" />
  ) : (
    <div className="flex flex-col gap-3">
      {data.length
        ? data.map(experiment => (
            <ExperimentCard
              key={experiment.id}
              onActions={onActions}
              data={experiment}
            />
          ))
        : emptyCollection}
    </div>
  );
};
