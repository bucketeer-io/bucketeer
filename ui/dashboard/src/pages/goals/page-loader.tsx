import { useState } from 'react';
import { Trans, useTranslation } from 'react-i18next';
import { goalDeleter } from '@api/goal';
import { goalUpdater, GoalUpdaterPayload } from '@api/goal/goal-updater';
import { invalidateGoals } from '@queries/goals';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_GOALS } from 'constants/routing';
import { useToast, useToggleOpen } from 'hooks';
import useActionWithURL from 'hooks/use-action-with-url';
import { Goal } from '@types';
import DeleteGoalModal from 'pages/goal-details/elements/delete-goal-modal';
import ConfirmModal from 'elements/confirm-modal';
import PageLayout from 'elements/page-layout';
import { EmptyCollection } from './collection-layout/empty-collection';
import { useFetchGoals } from './collection-loader/use-fetch-goals';
import AddGoalModal from './goals-modal/add-goal-modal';
import ConnectionsModal from './goals-modal/connections-modal';
import PageContent from './page-content';
import { GoalActions } from './types';

const PageLoader = () => {
  const { t } = useTranslation(['common', 'table']);
  const { notify } = useToast();
  const queryClient = useQueryClient();

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

  const [isOpenDeleteModal, onOpenDeleteModal, onCloseDeleteModal] =
    useToggleOpen(false);

  const [openConfirmModal, onOpenConfirmModal, onCloseConfirmModal] =
    useToggleOpen(false);

  const onHandleActions = (goal: Goal, type: GoalActions) => {
    setSelectedGoal(goal);
    if (type === 'CONNECTION') {
      onOpenConnectionModal();
    }
    if (type === 'DELETE') {
      onOpenDeleteModal();
    }
    if (['ARCHIVE', 'UNARCHIVE'].includes(type)) {
      onOpenConfirmModal();
    }
  };

  const mutation = useMutation({
    mutationFn: async (goal: Goal) => {
      return goalDeleter({
        id: goal.id,
        environmentId: currentEnvironment.id
      });
    },
    onSuccess: () => {
      onCloseDeleteModal();
      invalidateGoals(queryClient);
      notify({
        toastType: 'toast',
        messageType: 'success',
        message: (
          <span>
            <b>{selectedGoal?.name}</b>
            {` has been deleted successfully!`}
          </span>
        )
      });
    }
  });

  const onDeleteGoal = () => {
    if (selectedGoal?.id) {
      mutation.mutate(selectedGoal);
    }
  };

  const mutationState = useMutation({
    mutationFn: async (payload: GoalUpdaterPayload) => {
      return goalUpdater(payload);
    },
    onSuccess: data => {
      onCloseConfirmModal();
      invalidateGoals(queryClient);
      notify({
        message: (
          <span>
            <b>{data?.goal?.name}</b> {`has been successfully updated!`}
          </span>
        )
      });
      mutationState.reset();
    },
    onError: error =>
      notify({
        messageType: 'error',
        message: error?.message || 'Something went wrong.'
      })
  });

  const onUpdateGoal = async (payload: GoalUpdaterPayload) =>
    mutationState.mutate(payload);

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
      {isOpenDeleteModal && (
        <DeleteGoalModal
          goal={selectedGoal!}
          isOpen={isOpenDeleteModal}
          loading={mutation.isPending}
          onClose={onCloseDeleteModal}
          onSubmit={onDeleteGoal}
        />
      )}
      {openConfirmModal && selectedGoal && (
        <ConfirmModal
          isOpen={openConfirmModal}
          loading={mutationState.isPending}
          title={
            selectedGoal.archived
              ? t(`table:popover.unarchive-goal`)
              : t(`table:popover.archive-goal`)
          }
          description={
            <Trans
              i18nKey={
                selectedGoal.archived
                  ? 'table:goals.confirm-unarchive-desc'
                  : 'table:goals.confirm-archive-desc'
              }
              values={{ name: selectedGoal?.name }}
              components={{ bold: <strong /> }}
            />
          }
          onClose={onCloseConfirmModal}
          onSubmit={() =>
            onUpdateGoal({
              id: selectedGoal.id,
              name: selectedGoal.name,
              environmentId: currentEnvironment.id,
              description: selectedGoal.description,
              archived: selectedGoal.archived ? false : true
            })
          }
        />
      )}
    </>
  );
};

export default PageLoader;
