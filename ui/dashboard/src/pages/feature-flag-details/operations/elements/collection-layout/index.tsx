import { AutoOpsCount, AutoOpsRule } from '@types';
import { OperationModalState } from '../..';
import { OperationTab } from '../../types';
import ActiveContent from '../active';
import CompletedContent from '../completed';
import Overview from '../overview';

const CollectionLayout = ({
  currentTab,
  operations,
  opsCounts,
  onOperationActions
}: {
  currentTab: OperationTab;
  operations: AutoOpsRule[];
  opsCounts: AutoOpsCount[];
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
          opsCounts={opsCounts}
          operations={operations}
          onOperationActions={onOperationActions}
        />
      )}
      {currentTab === OperationTab.FINISHED && (
        <CompletedContent
          opsCounts={opsCounts}
          operations={operations}
          onOperationActions={onOperationActions}
        />
      )}
    </div>
  );
};

export default CollectionLayout;
