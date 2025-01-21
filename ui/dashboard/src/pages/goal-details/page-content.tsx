import { Trans } from 'react-i18next';
import { goalUpdater, GoalUpdaterPayload } from '@api/goal/goal-updater';
import { invalidateGoalDetails } from '@queries/goal-details';
import { invalidateGoals } from '@queries/goals';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast, useToggleOpen } from 'hooks';
import { useTranslation } from 'i18n';
import { Goal } from '@types';
import ConfirmModal from 'elements/confirm-modal';
import PageLayout from 'elements/page-layout';
import DeleteWarning from './elements/delete-warning';
import GoalActions from './elements/goal-actions';
import GoalConnections from './elements/goal-connections';
import GoalUpdateForm from './elements/goal-update-form';
import DeleteGoalModal from './goal-details-modal/delete-goal-modal';

const PageContent = ({ goal }: { goal: Goal }) => {
  const { t } = useTranslation(['common', 'form', 'table']);

  const { notify } = useToast();
  const queryClient = useQueryClient();

  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const [isOpenDeleteModal, onOpenDeleteModal, onCloseDeleteModal] =
    useToggleOpen(false);
  const [openConfirmModal, onOpenConfirmModal, onCloseConfirmModal] =
    useToggleOpen(false);

  const mutationState = useMutation({
    mutationFn: async (payload: GoalUpdaterPayload) => {
      return goalUpdater(payload);
    },
    onSuccess: data => {
      onCloseConfirmModal();
      invalidateGoalDetails(queryClient, {
        id: goal.id,
        environmentId: currentEnvironment.id
      });
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

  return (
    <PageLayout.Content className="gap-y-6 overflow-auto">
      <GoalUpdateForm goal={goal} onSubmit={onUpdateGoal} />
      {goal.experiments?.length > 0 && <GoalConnections goal={goal} />}
      <GoalActions
        title={
          goal.archived
            ? t(`table:popover.unarchive-goal`)
            : t(`table:popover.archive-goal`)
        }
        description={goal?.archived ? '' : t('form:goal-details.archive-desc')}
        btnText={
          goal.archived
            ? t(`table:popover.unarchive-goal`)
            : t(`table:popover.archive-goal`)
        }
        onClick={onOpenConfirmModal}
      />
      <GoalActions
        title={t('delete-goal')}
        description={t('form:goal-details.delete-desc')}
        btnText={t('delete-goal')}
        disabled={goal.experiments?.length > 0 || goal.isInUseStatus}
        onClick={onOpenDeleteModal}
      >
        {(goal.experiments?.length > 0 || goal.isInUseStatus) && (
          <DeleteWarning />
        )}
      </GoalActions>
      {isOpenDeleteModal && (
        <DeleteGoalModal
          goal={goal}
          isOpen={isOpenDeleteModal}
          loading={false}
          onClose={onCloseDeleteModal}
          onSubmit={() => {}}
        />
      )}
      {openConfirmModal && (
        <ConfirmModal
          isOpen={openConfirmModal}
          loading={mutationState.isPending}
          title={
            goal.archived
              ? t(`table:popover.unarchive-goal`)
              : t(`table:popover.archive-goal`)
          }
          description={
            <Trans
              i18nKey={
                goal.archived
                  ? 'table:goals.confirm-unarchive-desc'
                  : 'table:goals.confirm-archive-desc'
              }
              values={{ name: goal?.name }}
              components={{ bold: <strong /> }}
            />
          }
          onClose={onCloseConfirmModal}
          onSubmit={() =>
            onUpdateGoal({
              id: goal.id,
              name: goal.name,
              environmentId: currentEnvironment.id,
              description: goal.description,
              archived: !goal.archived
            })
          }
        />
      )}
    </PageLayout.Content>
  );
};

export default PageContent;
