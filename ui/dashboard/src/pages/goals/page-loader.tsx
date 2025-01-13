import { useState } from 'react';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToggleOpen } from 'hooks';
import { Goal } from '@types';
import PageLayout from 'elements/page-layout';
import { EmptyCollection } from './collection-layout/empty-collection';
import { useFetchGoals } from './collection-loader/use-fetch-goals';
import AddGoalModal from './goals-modal/add-goal-modal';
import ConnectionsModal from './goals-modal/connections-modal';
import PageContent from './page-content';
import { GoalActions } from './types';

const PageLoader = () => {
  const { consoleAccount } = useAuth();
  const currenEnvironment = getCurrentEnvironment(consoleAccount!);

  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchGoals({
    pageSize: 1,
    environmentId: currenEnvironment.id
  });

  const [selectedGoal, setSelectedGoal] = useState<Goal>();

  const [isOpenAddModal, onOpenAddModal, onCloseAddModal] =
    useToggleOpen(false);

  const [isOpenConnectionModal, onOpenConnectionModal, onCloseConnectionModal] =
    useToggleOpen(false);

  const onHandleActions = (goal: Goal, type: GoalActions) => {
    setSelectedGoal(goal);
    if (type === 'CONNECTION') {
      return onOpenConnectionModal();
    }
  };

  const isEmpty = collection?.goals.length === 0;

  return (
    <>
      {isLoading ? (
        <PageLayout.LoadingState />
      ) : isError ? (
        <PageLayout.ErrorState onRetry={refetch} />
      ) : isEmpty ? (
        <PageLayout.EmptyState>
          <EmptyCollection onAdd={onOpenAddModal} />
        </PageLayout.EmptyState>
      ) : (
        <PageContent onAdd={onOpenAddModal} onHandleActions={onHandleActions} />
      )}

      {isOpenAddModal && (
        <AddGoalModal isOpen={isOpenAddModal} onClose={onCloseAddModal} />
      )}
      {isOpenConnectionModal && selectedGoal && (
        <ConnectionsModal
          isOpen={isOpenConnectionModal}
          goal={selectedGoal}
          onClose={onCloseConnectionModal}
        />
      )}
    </>
  );
};

export default PageLoader;
