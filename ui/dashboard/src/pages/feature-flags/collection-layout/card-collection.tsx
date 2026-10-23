import { ReactElement } from 'react';
import { MDIcon } from 'react-icons-material-design';
import { AutoOpsRule, Environment, Feature, Rollout } from '@types';
import { FeatureCard } from 'components/mobile-card/feature-card';
import PageLayout from 'elements/page-layout';
import { FlagActionType } from '../types';

interface CardCollectionProps {
  data: Feature[];
  onActions: (item: Feature, type: FlagActionType) => void;
  filterTags?: string[];
  rollouts: Rollout[];
  autoOpsRules: AutoOpsRule[];
  popoverOptions: {
    label: string;
    icon: MDIcon;
    value: string;
  }[];
  currentEnvironment: Environment;
  editable: boolean;
  handleGetMaintainerInfo: (email: string) => string;
  handleTagFilters: (tag: string) => void;
  emptyCollection?: ReactElement;
  isLoading?: boolean;
}

export const CardCollection = ({
  isLoading,
  emptyCollection,
  data,
  filterTags,
  rollouts,
  autoOpsRules,
  currentEnvironment,
  popoverOptions,
  editable,
  handleGetMaintainerInfo,
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
              currentEnvironment={currentEnvironment}
              popoverOptions={popoverOptions}
              editable={editable}
              handleGetMaintainerInfo={handleGetMaintainerInfo}
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
