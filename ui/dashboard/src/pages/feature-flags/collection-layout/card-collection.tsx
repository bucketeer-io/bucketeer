import { ReactElement } from 'react';
import { Account, AutoOpsRule, Feature, Rollout } from '@types';
import { FeatureCard } from 'components/mobile-card/feature-card';
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
  data,
  accounts,
  filterTags,
  rollouts,
  autoOpsRules,
  handleTagFilters,
  onActions
}: CardCollectionProps) => {
  return (
    <div className="flex flex-col gap-3">
      {data.map(feature => (
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
      ))}
    </div>
  );
};
