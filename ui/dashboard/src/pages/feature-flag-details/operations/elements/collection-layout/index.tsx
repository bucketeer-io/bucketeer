import { AutoOpsRule } from '@types';
import { Rollout } from '@types';
import { OperationModalState } from '../..';
import { OperationTab } from '../../types';
import ActiveContent from '../active';
import CompletedContent from '../completed';
import Overview from '../overview';

const CollectionLayout = ({
  currentTab,
  operations,
  rollouts,
  onOperationActions
}: {
  currentTab: OperationTab;
  operations: AutoOpsRule[];
  rollouts: Rollout[];
  onOperationActions: (data: OperationModalState) => void;
}) => {
  return (
    <div>
      <Overview
        onOperationActions={operationType =>
          onOperationActions({
            operationType,
            actionType: 'NEW'
          })
        }
      />
      {currentTab === OperationTab.ACTIVE && (
        <ActiveContent
          rollouts={rollouts}
          operations={operations}
          onOperationActions={onOperationActions}
        />
      )}
      {currentTab === OperationTab.FINISHED && (
        <CompletedContent
          operations={operations}
          rollouts={rollouts}
          onOperationActions={onOperationActions}
        />
      )}
    </div>
  );
};

export default CollectionLayout;
