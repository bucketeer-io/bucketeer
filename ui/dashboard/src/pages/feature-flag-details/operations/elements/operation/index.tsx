import { OperationModalState } from '../..';
import { OperationCombinedType } from '../../types';
import OperationProgress from '../operation-progress';
import OperationStatus from '../operation-status';

interface Props {
  isFinished: boolean;
  operation: OperationCombinedType;
  onActions: (data: OperationModalState) => void;
}

const Operation = ({ isFinished, operation, onActions }: Props) => {
  return (
    <div className="flex flex-col p-5 shadow-card rounded-lg bg-white gap-y-4">
      <OperationStatus
        operation={operation}
        isFinished={isFinished}
        onActions={onActions}
      />
      <OperationProgress operation={operation} />
    </div>
  );
};

export default Operation;
