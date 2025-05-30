import { AutoOpsCount, AutoOpsRule, Rollout } from '@types';
import { OperationModalState } from '../..';
import { OperationTab } from '../../types';
import ActiveContent from '../active';
import CompletedContent from '../completed';

const CollectionLayout = ({
  currentTab,
  operations,
  opsCounts,
  rollouts,
  onOperationActions
}: {
  currentTab: OperationTab;
  operations: AutoOpsRule[];
  opsCounts: AutoOpsCount[];
  rollouts: Rollout[];
  onOperationActions: (data: OperationModalState) => void;
}) => {
  return (
    <div>
      {currentTab === OperationTab.ACTIVE && (
        <ActiveContent
          opsCounts={opsCounts}
          rollouts={rollouts}
          operations={operations}
          onOperationActions={onOperationActions}
        />
      )}
      {currentTab === OperationTab.FINISHED && (
        <CompletedContent
          opsCounts={opsCounts}
          operations={operations}
          rollouts={rollouts}
          onOperationActions={onOperationActions}
        />
      )}
    </div>
  );
};

export default CollectionLayout;
