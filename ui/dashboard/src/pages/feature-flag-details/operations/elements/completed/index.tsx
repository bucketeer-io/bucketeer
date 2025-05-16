import { useMemo } from 'react';
import { AutoOpsRule, Rollout } from '@types';
import { OperationModalState } from '../..';
import { OperationCombinedType } from '../../types';
import Operation from '../operation';

const CompletedContent = ({
  rollouts,
  operations,
  onOperationActions
}: {
  rollouts: Rollout[];
  operations: AutoOpsRule[];
  onOperationActions: (data: OperationModalState) => void;
}) => {
  const completedStatuses = useMemo(() => ['STOPPED', 'FINISHED'], []);

  const operationCompletedData = useMemo(
    () =>
      operations?.filter(item =>
        completedStatuses.includes(item.autoOpsStatus)
      ),
    [operations, completedStatuses]
  );

  const rolloutCompletedData = useMemo(
    () => rollouts?.filter(item => completedStatuses.includes(item.status)),
    [rollouts, completedStatuses]
  );

  const sortedData: OperationCombinedType[] = useMemo(
    () =>
      [...operationCompletedData, ...rolloutCompletedData].sort(
        (a, b) => Number(a.createdAt) - Number(b.createdAt)
      ) as OperationCombinedType[],
    [operationCompletedData, rolloutCompletedData]
  );

  return (
    <div className="flex flex-col w-full gap-y-6">
      {sortedData?.map((item, index) => (
        <Operation
          key={index}
          isCompleted={true}
          operation={item}
          onActions={onOperationActions}
        />
      ))}
    </div>
  );
};

export default CompletedContent;
