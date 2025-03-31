import { OperationCombinedType } from '../../types';
import OperationStatus from '../operation-status';

interface Props {
  isCompleted: boolean;
  operation: OperationCombinedType;
  onActions: () => void;
}

const Operation = ({ isCompleted, operation, onActions }: Props) => {
  return (
    <div className="p-5 shadow-card rounded-lg bg-white">
      <OperationStatus
        operation={operation}
        isCompleted={isCompleted}
        onActions={onActions}
      />
    </div>
  );
};

export default Operation;
