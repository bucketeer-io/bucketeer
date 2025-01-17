import { useTranslation } from 'i18n';
import { Goal } from '@types';
import PageLayout from 'elements/page-layout';
import DeleteWarning from './elements/delete-warning';
import GoalActions from './elements/goal-actions';
import GoalConnections from './elements/goal-connections';
import GoalUpdateForm from './elements/goal-update-form';

const PageContent = ({ goal }: { goal: Goal }) => {
  const { t } = useTranslation(['common', 'form']);

  return (
    <PageLayout.Content className="gap-y-6 overflow-auto">
      <GoalUpdateForm goal={goal} />
      {goal.experiments?.length > 0 && <GoalConnections goal={goal} />}
      <GoalActions
        title={t('archive-goal')}
        description={t('form:goal-details.archive-desc')}
        btnText={t('archive-goal')}
        onClick={() => {}}
      />
      <GoalActions
        title={t('delete-goal')}
        description={t('form:goal-details.delete-desc')}
        btnText={t('delete-goal')}
        onClick={() => {}}
      >
        {(goal.experiments?.length > 0 || goal.isInUseStatus) && (
          <DeleteWarning />
        )}
      </GoalActions>
    </PageLayout.Content>
  );
};

export default PageContent;
