import { ReactElement } from 'react';
import { Push, Tag } from '@types';
import { PushCard } from 'components/mobile-card/push-card';
import PageLayout from 'elements/page-layout';
import { PushActionsType } from '../types';

interface CardCollectionProps {
  data: Push[];
  tags: Tag[];
  onActions: (item: Push, type: PushActionsType) => void;
  emptyCollection?: ReactElement;
  isLoading?: boolean;
}

export const CardCollection = ({
  data,
  tags,
  isLoading,
  emptyCollection,
  onActions
}: CardCollectionProps) => {
  return isLoading ? (
    <PageLayout.LoadingState className="py-10" />
  ) : (
    <div className="flex flex-col gap-3">
      {data.length
        ? data.map(push => (
            <PushCard
              key={push.id}
              onActions={onActions}
              data={push}
              tags={tags}
            />
          ))
        : emptyCollection}
    </div>
  );
};
