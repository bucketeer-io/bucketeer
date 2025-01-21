import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import EvaluationTable from './evaluation-table';

type Option = {
  label: string;
  value: string;
};

const evaluationOptions: Option[] = [
  {
    label: 'Option 1',
    value: 'option-1'
  },
  {
    label: 'Option 2',
    value: 'option-2'
  },
  {
    label: 'Option 3',
    value: 'option-3'
  }
];

const Evaluation = () => {
  const { t } = useTranslation(['common', 'form']);

  const [selectedEvaluation, setSelectedEvaluation] = useState<Option | null>(
    null
  );

  return (
    <div className="flex flex-col h-full gap-y-6">
      <DropdownMenu>
        <DropdownMenuTrigger
          label={selectedEvaluation?.label || ''}
          placeholder={t(`form:select-evaluation`)}
          variant="secondary"
          className="w-1/2 min-w-[300px] max-w-[528px]"
        />
        <DropdownMenuContent className="w-[235px]" align="start">
          {evaluationOptions.map((item, index) => (
            <DropdownMenuItem
              key={index}
              value={item.value}
              label={item.label}
              onSelectOption={() => setSelectedEvaluation(item)}
            />
          ))}
        </DropdownMenuContent>
      </DropdownMenu>
      <EvaluationTable />
    </div>
  );
};

export default Evaluation;
