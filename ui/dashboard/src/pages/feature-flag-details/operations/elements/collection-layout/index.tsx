import { AutoOpsRule, Feature, Rollout } from '@types';
import TableListContent from 'elements/table-list-content';
import { OperationModalState } from '../..';
import { OperationTab } from '../../types';
import ActiveContent from '../active';
import CompletedContent from '../completed';
import Overview from '../overview';

const CollectionLayout = ({
  currentTab,
  feature,
  rollouts,
  operations,
  onOperationActions
}: {
  currentTab: OperationTab;
  feature: Feature;
  rollouts: Rollout[];
  operations: AutoOpsRule[];
  onOperationActions: (data: OperationModalState) => void;
}) => {
  return (
    <TableListContent>
      <Overview onChangeFilters={() => {}} />
      {currentTab === OperationTab.ACTIVE && (
        <ActiveContent
          rollouts={rollouts}
          operations={operations}
          onOperationActions={onOperationActions}
        />
      )}
      {currentTab === OperationTab.COMPLETED && (
        <CompletedContent rollouts={rollouts} operations={operations} />
      )}
    </TableListContent>
  );
};

export default CollectionLayout;
