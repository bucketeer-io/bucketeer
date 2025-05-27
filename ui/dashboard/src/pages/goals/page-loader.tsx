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
import AddGoalModal from './goals-modal/add-goal-modal';
import ConnectionsModal from './goals-modal/connections-modal';
import PageContent from './page-content';
import { GoalActions } from './types';

const PageLoader = () => {
  const { t } = useTranslation(['common', 'table', 'message']);
  const { notify, errorNotify } = useToast();
  const queryClient = useQueryClient();

  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const { isAdd, onOpenAddModal, onCloseActionModal } = useActionWithURL({
    closeModalPath: `/${currentEnvironment.urlCode}${PAGE_PATH_GOALS}`
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
    switch (type) {
      case 'CONNECTION':
        return onOpenConnectionModal();
      case 'DELETE':
        return onOpenDeleteModal();
      case 'ARCHIVE':
      case 'UNARCHIVE':
        return onOpenConfirmModal();
      default:
        return;
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
        message: t('message:collection-action-success', {
          collection: t('source-type.goal'),
          action: t('deleted')
        })
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
    onSuccess: () => {
      onCloseConfirmModal();
      invalidateGoals(queryClient);
      notify({
        message: t('message:collection-action-success', {
          collection: t('source-type.goal'),
          action: t('updated')
        })
      });
      mutationState.reset();
    },
    onError: error => errorNotify(error)
  });

  const onUpdateGoal = async (payload: GoalUpdaterPayload) =>
    mutationState.mutate(payload);

  return (
    <>
      <PageContent onAdd={onOpenAddModal} onHandleActions={onHandleActions} />
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
