import { useMemo } from 'react';
import { AutoOpsCount, AutoOpsRule, Rollout } from '@types';
import { OperationModalState } from '../..';
import { OperationCombinedType } from '../../types';
import { EmptyCollection } from '../collection-layout/empty-collection';
import Operation from '../operation';

const ActiveContent = ({
  operations,
  opsCounts,
  rollouts,
  onOperationActions
}: {
  operations: AutoOpsRule[];
  opsCounts: AutoOpsCount[];
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

  const eventRateActiveData = useMemo(
    () =>
      operations?.filter(
        item =>
          activeStatuses.includes(item.autoOpsStatus) &&
          item.opsType === 'EVENT_RATE'
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
      [
        ...eventRateActiveData,
        ...scheduleActiveData,
        ...rolloutActiveData
      ] as OperationCombinedType[],
    [eventRateActiveData, scheduleActiveData, rolloutActiveData]
  );

  if (!operationData.length) return <EmptyCollection />;

  return (
    <div className="flex flex-col w-full gap-y-6 pb-6">
      {operationData?.map((item, index) => (
        <Operation
          key={index}
          isFinished={false}
          operation={item}
          opsCounts={opsCounts}
          onActions={onOperationActions}
        />
      ))}
    </div>
  );
};

export default ActiveContent;
