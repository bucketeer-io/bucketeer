import { ReactElement } from 'react';
import { Account, AutoOpsRule, Feature, Rollout } from '@types';
import { FeatureCard } from 'components/mobile-card/feature-card';
import PageLayout from 'elements/page-layout';
import { FlagActionType } from '../types';

interface CardCollectionProps {
  data: Feature[];
  onActions: (item: Feature, type: FlagActionType) => void;
  accounts: Account[];
  filterTags?: string[];
  rollouts: Rollout[];
  autoOpsRules: AutoOpsRule[];
  handleTagFilters: (tag: string) => void;
  emptyCollection?: ReactElement;
  isLoading?: boolean;
}

export const CardCollection = ({
  isLoading,
  emptyCollection,
  data,
  accounts,
  filterTags,
  rollouts,
  autoOpsRules,
  handleTagFilters,
  onActions
}: CardCollectionProps) => {
  return isLoading ? (
    <PageLayout.LoadingState className="py-10" />
  ) : (
    <div className="flex flex-col gap-3">
      {data.length
        ? data.map(feature => (
            <FeatureCard
              key={feature.id}
              data={feature}
              accounts={accounts}
              filterTags={filterTags}
              rollouts={rollouts}
              autoOpsRules={autoOpsRules}
              handleTagFilters={handleTagFilters}
              onActions={onActions}
            />
          ))
        : emptyCollection}
    </div>
  );
};
