import { useMemo } from 'react';
import { useQueryAutoOps } from '@queries/auto-ops';
import { useQueryRollouts } from '@queries/rollouts';
import { getCurrentEnvironment, useAuth } from 'auth';
import { Feature } from '@types';
import FormLoading from 'elements/form-loading';
import { OperationTab } from '../../types';
import ActiveContent from '../active';
import CompletedContent from '../completed';
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

  const params = useMemo(
    () => ({
      cursor: String(0),
      featureIds: [feature.id],
      environmentId: currentEnvironment.id
    }),
    [feature, currentEnvironment]
  );

  const { data: rolloutCollection, isLoading: isRolloutLoading } =
    useQueryRollouts({
      params
    });

  const { data: operationCollection, isLoading: isOperationLoading } =
    useQueryAutoOps({
      params
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
          {currentTab === OperationTab.ACTIVE && (
            <ActiveContent rollouts={rollouts} operations={operations} />
          )}
          {currentTab === OperationTab.COMPLETED && (
            <CompletedContent rollouts={rollouts} operations={operations} />
          )}
        </>
      )}
    </div>
  );
};

export default CollectionLoader;
