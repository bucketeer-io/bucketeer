import { useMemo } from 'react';
import { AutoOpsRule, Rollout } from '@types';
import { OperationModalState } from '../..';
import { OperationCombinedType } from '../../types';
import Operation from '../operation';

const ActiveContent = ({
  rollouts,
  operations,
  onOperationActions
}: {
  rollouts: Rollout[];
  operations: AutoOpsRule[];
  onOperationActions: (data: OperationModalState) => void;
}) => {
  const activeStatuses = useMemo(() => ['WAITING', 'RUNNING'], []);

  const operationActiveData = useMemo(
    () =>
      operations?.filter(item => activeStatuses.includes(item.autoOpsStatus)),
    [operations, activeStatuses]
  );

  const rolloutActiveData = useMemo(
    () => rollouts?.filter(item => activeStatuses.includes(item.status)),
    [rollouts, activeStatuses]
  );
  const sortedData: OperationCombinedType[] = useMemo(
    () =>
      [...operationActiveData, ...rolloutActiveData].sort(
        (a, b) => Number(a.createdAt) - Number(b.createdAt)
      ) as OperationCombinedType[],
    [operationActiveData, rolloutActiveData]
  );

  return (
    <div className="flex flex-col w-full gap-y-6">
      {sortedData?.map((item, index) => (
        <Operation
          key={index}
          isCompleted={false}
          operation={item}
          onActions={onOperationActions}
        />
      ))}
    </div>
  );
};

export default ActiveContent;
