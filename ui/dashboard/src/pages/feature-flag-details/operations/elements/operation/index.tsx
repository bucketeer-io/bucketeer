import { OperationModalState } from '../..';
import { OperationCombinedType } from '../../types';
import OperationProgress from '../operation-progress';
import OperationStatus from '../operation-status';

interface Props {
  isCompleted: boolean;
  operation: OperationCombinedType;
  onActions: (data: OperationModalState) => void;
}

const Operation = ({ isCompleted, operation, onActions }: Props) => {
  return (
    <div className="p-5 shadow-card rounded-lg bg-white">
      <OperationStatus
        operation={operation}
        isCompleted={isCompleted}
        onActions={onActions}
      />
      <OperationProgress />
    </div>
  );
};

export default Operation;
