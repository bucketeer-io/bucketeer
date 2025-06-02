import { Trans } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { goalDeleter } from '@api/goal';
import { goalUpdater, GoalUpdaterPayload } from '@api/goal/goal-updater';
import { invalidateGoalDetails } from '@queries/goal-details';
import { invalidateGoals } from '@queries/goals';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_GOALS } from 'constants/routing';
import { useToast, useToggleOpen } from 'hooks';
import { useTranslation } from 'i18n';
import { Goal } from '@types';
import InfoMessage from 'components/info-message';
import ConfirmModal from 'elements/confirm-modal';
import PageLayout from 'elements/page-layout';
import DeleteGoalModal from './elements/delete-goal-modal';
import GoalActions from './elements/goal-actions';
import GoalConnections from './elements/goal-connections';
import GoalUpdateForm from './elements/goal-update-form';

const PageContent = ({ goal }: { goal: Goal }) => {
  const navigate = useNavigate();
  const { t } = useTranslation(['common', 'form', 'table', 'message']);

  const { notify, errorNotify } = useToast();
  const queryClient = useQueryClient();

  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const [isOpenDeleteModal, onOpenDeleteModal, onCloseDeleteModal] =
    useToggleOpen(false);
  const [openConfirmModal, onOpenConfirmModal, onCloseConfirmModal] =
    useToggleOpen(false);

  const mutation = useMutation({
    mutationFn: async () => {
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
      navigate(`/${currentEnvironment.urlCode}/${PAGE_PATH_GOALS}`);
    }
  });

  const mutationState = useMutation({
    mutationFn: async (payload: GoalUpdaterPayload) => {
      return goalUpdater(payload);
    },
    onSuccess: () => {
      onCloseConfirmModal();
      invalidateGoalDetails(queryClient, {
        id: goal.id,
        environmentId: currentEnvironment.id
      });
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
  const connections = goal.experiments?.length || goal.autoOpsRules?.length;

  return (
    <PageLayout.Content className="p-6 gap-y-6 overflow-auto">
      <GoalUpdateForm goal={goal} onSubmit={onUpdateGoal} />
      {connections > 0 && <GoalConnections goal={goal} />}
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
        disabled={goal.isInUseStatus}
      >
        {(goal.experiments?.length > 0 || goal.isInUseStatus) && (
          <InfoMessage title={t('form:goal-details.archive-warning-desc')} />
        )}
      </GoalActions>
      <GoalActions
        title={t('delete-goal')}
        description={t('form:goal-details.delete-desc')}
        btnText={t('delete-goal')}
        disabled={goal.experiments?.length > 0 || goal.isInUseStatus}
        onClick={onOpenDeleteModal}
      >
        {(goal.experiments?.length > 0 || goal.isInUseStatus) && (
          <InfoMessage title={t('form:goal-details.delete-warning-desc')} />
        )}
      </GoalActions>
      {isOpenDeleteModal && (
        <DeleteGoalModal
          goal={goal}
          isOpen={isOpenDeleteModal}
          loading={mutation.isPending}
          onClose={onCloseDeleteModal}
          onSubmit={() => mutation.mutate()}
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
