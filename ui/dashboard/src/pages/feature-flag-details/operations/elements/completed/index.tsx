import { useMemo } from 'react';
import { AutoOpsCount, AutoOpsRule } from '@types';
import { OperationModalState } from '../..';
import { OperationCombinedType } from '../../types';
import Operation from '../operation';

const CompletedContent = ({
  operations,
  opsCounts,
  onOperationActions
}: {
  operations: AutoOpsRule[];
  opsCounts: AutoOpsCount[];
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

  const operationData = useMemo(
    () =>
      [
        ...eventRateCompletedData,
        ...scheduleCompletedData
      ] as OperationCombinedType[],
    [scheduleCompletedData, eventRateCompletedData]
  );

  return (
    <div className="flex flex-col w-full gap-y-6">
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
