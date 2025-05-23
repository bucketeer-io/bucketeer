import { useMemo } from 'react';
import { AutoOpsRule } from '@types';
import { Rollout } from '@types';
import { OperationModalState } from '../..';
import { OperationCombinedType } from '../../types';
import Operation from '../operation';

const CompletedContent = ({
  operations,
  rollouts,
  onOperationActions
}: {
  operations: AutoOpsRule[];
  rollouts: Rollout[];
  onOperationActions: (data: OperationModalState) => void;
}) => {
  const completedStatuses = useMemo(() => ['STOPPED', 'FINISHED'], []);

  const scheduleCompletedData = useMemo(
    () =>
      operations?.filter(item =>
        completedStatuses.includes(item.autoOpsStatus)
      ),
    [operations, completedStatuses]
  );
  const rolloutCompletedData = useMemo(
    () =>
      rollouts?.filter(
        item =>
          completedStatuses.includes(item.status) &&
          ['MANUAL_SCHEDULE', 'TEMPLATE_SCHEDULE'].includes(item.type)
      ),
    [rollouts, completedStatuses]
  );

  const operationData = useMemo(
    () =>
      [
        ...scheduleCompletedData,
        ...rolloutCompletedData
      ] as OperationCombinedType[],
    [scheduleCompletedData, rolloutCompletedData]
  );

  return (
    <div className="flex flex-col w-full gap-y-6">
      {operationData?.map((item, index) => (
        <Operation
          key={index}
          isFinished={true}
          operation={item}
          onActions={onOperationActions}
        />
      ))}
    </div>
  );
};

export default CompletedContent;
