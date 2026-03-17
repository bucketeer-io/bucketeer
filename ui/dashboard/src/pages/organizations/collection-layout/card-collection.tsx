import { ReactElement } from 'react';
import { Organization } from '@types';
import { OrganizationCard } from 'components/mobile-card/organization-card';
import PageLayout from 'elements/page-layout';
import { OrganizationActionsType } from '../types';

interface CardCollectionProps {
  data: Organization[];
  onActions: (item: Organization, type: OrganizationActionsType) => void;
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
        ? data.map(goal => (
            <OrganizationCard key={goal.id} onActions={onActions} data={goal} />
          ))
        : emptyCollection}
    </div>
  );
};
