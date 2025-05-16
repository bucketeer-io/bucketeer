import { useMemo } from 'react';
import { AutoOpsRule } from '@types';
import { OperationModalState } from '../..';
import { OperationCombinedType } from '../../types';
import Operation from '../operation';

const CompletedContent = ({
  operations,
  onOperationActions
}: {
  operations: AutoOpsRule[];
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

  const operationData = useMemo(
    () => [...scheduleCompletedData] as OperationCombinedType[],
    [scheduleCompletedData]
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
