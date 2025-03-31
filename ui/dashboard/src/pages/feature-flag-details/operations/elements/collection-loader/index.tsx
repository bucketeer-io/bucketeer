import { useQueryAutoOps } from '@queries/auto-ops';
import { useQueryRollouts } from '@queries/rollouts';
import { getCurrentEnvironment, useAuth } from 'auth';
import { Feature } from '@types';
import FormLoading from 'elements/form-loading';
import { OperationTab } from '../../types';
import ActiveContent from '../active';
import Overview from '../overview';

const CollectionLoader = ({
  currentTab,
  feature
}: {
  currentTab: OperationTab;
  feature: Feature;
}) => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { data: rolloutCollection, isLoading: isRolloutLoading } =
    useQueryRollouts({
      params: {
        cursor: String(0),
        featureIds: [feature.id],
        environmentId: currentEnvironment.id
      }
    });

  const { data: operationCollection, isLoading: isOperationLoading } =
    useQueryAutoOps({
      params: {
        cursor: String(0),
        featureIds: [feature.id],
        environmentId: currentEnvironment.id
      }
    });

  const rollouts = rolloutCollection?.progressiveRollouts || [];
  const operations = operationCollection?.autoOpsRules || [];
  return (
    <div>
      <Overview onChangeFilters={() => {}} />
      {isRolloutLoading || isOperationLoading ? (
        <FormLoading />
      ) : (
        <>
          {currentTab && (
            <ActiveContent rollouts={rollouts} operations={operations} />
          )}
        </>
      )}
    </div>
  );
};

export default CollectionLoader;
