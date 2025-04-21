import { useState } from 'react';
import ActionBar from './action-bar';
import ResultList from './result-list';

export type GroupByType = 'FLAG' | 'USER';

interface Props {
  onResetFields: () => void;
  onEditFields: () => void;
}

const DebuggerResults = ({ onResetFields, onEditFields }: Props) => {
  const [groupBy, setGroupBy] = useState<GroupByType>('FLAG');

  return (
    <div className="flex flex-col w-full gap-y-6">
      <ActionBar
        groupBy={groupBy}
        setGroupBy={setGroupBy}
        onResetFields={onResetFields}
        onEditFields={onEditFields}
      />
      <ResultList groupBy={groupBy} evaluationFeatures={[]} />
    </div>
  );
};

export default DebuggerResults;
