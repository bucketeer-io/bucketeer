import { Trans } from 'react-i18next';
import { useToggleOpen } from 'hooks';
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

  const [isOpenDeleteModal, onOpenDeleteModal, onCloseDeleteModal] =
    useToggleOpen(false);
  const [openConfirmModal, onOpenConfirmModal, onCloseConfirmModal] =
    useToggleOpen(false);

  return (
    <PageLayout.Content className="gap-y-6 overflow-auto">
      <GoalUpdateForm goal={goal} />
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
          onClose={onCloseConfirmModal}
          onSubmit={() => {}}
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
          loading={false}
        />
      )}
    </PageLayout.Content>
  );
};

export default PageContent;
