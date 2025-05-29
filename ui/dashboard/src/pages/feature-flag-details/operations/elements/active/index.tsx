import { useMemo } from 'react';
import { AutoOpsCount, AutoOpsRule } from '@types';
import { OperationModalState } from '../..';
import { OperationCombinedType } from '../../types';
import Operation from '../operation';

const ActiveContent = ({
  operations,
  opsCounts,
  onOperationActions
}: {
  operations: AutoOpsRule[];
  opsCounts: AutoOpsCount[];
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

  const operationData = useMemo(
    () =>
      [
        ...eventRateActiveData,
        ...scheduleActiveData
      ] as OperationCombinedType[],
    [eventRateActiveData, scheduleActiveData]
  );

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
