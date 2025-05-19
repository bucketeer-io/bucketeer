import { useMemo } from 'react';
import { AutoOpsRule } from '@types';
import { Rollout } from '@types';
import { OperationModalState } from '../..';
import { OperationCombinedType } from '../../types';
import Operation from '../operation';

const ActiveContent = ({
  operations,
  rollouts,
  onOperationActions
}: {
  operations: AutoOpsRule[];
  rollouts: Rollout[];
  onOperationActions: (data: OperationModalState) => void;
}) => {
  const activeStatuses = useMemo(() => ['WAITING', 'RUNNING'], []);

  const scheduleActiveData = useMemo(
    () =>
      operations?.filter(
        item =>
          activeStatuses.includes(item.autoOpsStatus) &&
          item.opsType === 'SCHEDULE'
      ),
    [operations, activeStatuses]
  );

  const rolloutActiveData = useMemo(
    () =>
      rollouts?.filter(
        item =>
          activeStatuses.includes(item.status) &&
          ['MANUAL_SCHEDULE', 'TEMPLATE_SCHEDULE'].includes(item.type)
      ),
    [rollouts, activeStatuses]
  );

  const operationData = useMemo(
    () =>
      [...scheduleActiveData, ...rolloutActiveData] as OperationCombinedType[],
    [scheduleActiveData, rolloutActiveData]
  );

  return (
    <div className="flex flex-col w-full gap-y-6">
      {operationData?.map((item, index) => (
        <Operation
          key={index}
          isFinished={false}
          operation={item}
          onActions={onOperationActions}
        />
      ))}
    </div>
  );
};

export default ActiveContent;
