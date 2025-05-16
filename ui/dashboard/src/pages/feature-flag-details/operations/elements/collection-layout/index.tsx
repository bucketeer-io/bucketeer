import { AutoOpsRule } from '@types';
import { OperationModalState } from '../..';
import { OperationTab } from '../../types';
import ActiveContent from '../active';
import CompletedContent from '../completed';
import Overview from '../overview';

const CollectionLayout = ({
  currentTab,
  operations,
  onOperationActions
}: {
  currentTab: OperationTab;
  operations: AutoOpsRule[];
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
          operations={operations}
          onOperationActions={onOperationActions}
        />
      )}
      {currentTab === OperationTab.COMPLETED && (
        <CompletedContent
          operations={operations}
          onOperationActions={onOperationActions}
        />
      )}
    </div>
  );
};

export default CollectionLayout;
