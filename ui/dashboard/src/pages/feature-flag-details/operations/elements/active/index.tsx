import { useMemo } from 'react';
import { AutoOpsRule } from '@types';
import { OperationModalState } from '../..';
import { OperationCombinedType } from '../../types';
import Operation from '../operation';

const ActiveContent = ({
  operations,
  onOperationActions
}: {
  operations: AutoOpsRule[];
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

  const sortedData: OperationCombinedType[] = useMemo(
    () =>
      [...scheduleActiveData].sort(
        (a, b) => Number(a.createdAt) - Number(b.createdAt)
      ) as OperationCombinedType[],
    [scheduleActiveData]
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
