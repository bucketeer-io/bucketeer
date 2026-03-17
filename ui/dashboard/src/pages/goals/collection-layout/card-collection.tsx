import { ReactElement } from 'react';
import { Goal } from '@types';
import { GoalCard } from 'components/mobile-card/goal-card';
import PageLayout from 'elements/page-layout';
import { GoalActions } from '../types';

interface CardCollectionProps {
  data: Goal[];
  onActions: (item: Goal, type: GoalActions) => void;
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
            <GoalCard key={goal.id} onActions={onActions} data={goal} />
          ))
        : emptyCollection}
    </div>
  );
};
