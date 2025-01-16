import { useState } from 'react';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_GOALS } from 'constants/routing';
import { useToggleOpen } from 'hooks';
import useActionWithURL from 'hooks/use-action-with-url';
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
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const { isAdd, onOpenAddModal, onCloseActionModal } = useActionWithURL({
    closeModalPath: `/${currentEnvironment.urlCode}${PAGE_PATH_GOALS}`
  });

  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchGoals({
    pageSize: 1,
    environmentId: currentEnvironment.id
  });

  const [selectedGoal, setSelectedGoal] = useState<Goal>();

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

      {isAdd && <AddGoalModal isOpen={isAdd} onClose={onCloseActionModal} />}
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
