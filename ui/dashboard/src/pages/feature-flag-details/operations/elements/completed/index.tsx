import { useMemo } from 'react';
import { AutoOpsCount, AutoOpsRule, Rollout } from '@types';
import { OperationModalState } from '../..';
import { OperationCombinedType } from '../../types';
import { EmptyCollection } from '../collection-layout/empty-collection';
import Operation from '../operation';

const CompletedContent = ({
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
  const completedStatuses = useMemo(() => ['STOPPED', 'FINISHED'], []);

  const scheduleCompletedData = useMemo(
    () =>
      operations?.filter(
        item =>
          completedStatuses.includes(item.autoOpsStatus) &&
          item.opsType === 'SCHEDULE'
      ),
    [operations, completedStatuses]
  );

  const eventRateCompletedData = useMemo(
    () =>
      operations?.filter(
        item =>
          completedStatuses.includes(item.autoOpsStatus) &&
          item.opsType === 'EVENT_RATE'
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
        ...eventRateCompletedData,
        ...scheduleCompletedData,
        ...rolloutCompletedData
      ] as OperationCombinedType[],
    [scheduleCompletedData, eventRateCompletedData, rolloutCompletedData]
  );

  if (!operationData.length) return <EmptyCollection />;

  return (
    <div className="flex flex-col w-full gap-y-6 pb-6">
      {operationData?.map((item, index) => (
        <Operation
          key={index}
          isFinished={true}
          operation={item}
          opsCounts={opsCounts}
          onActions={onOperationActions}
        />
      ))}
    </div>
  );
};

export default CompletedContent;
