import { ReactElement } from 'react';
import { Account, Team } from '@types';
import { MemberCard } from 'components/mobile-card/member-card';
import PageLayout from 'elements/page-layout';
import { MemberActionsType, MembersFilters } from '../types';

interface CardCollectionProps {
  isLoading: boolean;
  emptyCollection?: ReactElement;
  data: Account[];
  filters: MembersFilters;
  teams: Team[];
  onActions: (item: Account, type: MemberActionsType) => void;
  setFilters: (values: Partial<MembersFilters>) => void;
}

export const CardCollection = ({
  isLoading,
  emptyCollection,
  data,
  filters,
  teams,
  setFilters,
  onActions
}: CardCollectionProps) => {
  return isLoading ? (
    <PageLayout.LoadingState className="py-10" />
  ) : (
    <div className="flex flex-col gap-3">
      {data.length
        ? data.map(account => (
            <MemberCard
              key={account.email}
              onActions={onActions}
              data={account}
              filters={filters}
              teams={teams}
              setFilters={setFilters}
            />
          ))
        : emptyCollection}
    </div>
  );
};
