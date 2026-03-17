import { ReactElement } from 'react';
import { APIKey } from '@types';
import { ApiCard } from 'components/mobile-card/api-card';
import PageLayout from 'elements/page-layout';
import { APIKeyActionsType } from '../types';

interface CardCollectionProps {
  data: APIKey[];
  onActions: (item: APIKey, type: APIKeyActionsType) => void;
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
            <ApiCard
              key={experiment.id}
              onActions={onActions}
              data={experiment}
            />
          ))
        : emptyCollection}
    </div>
  );
};
