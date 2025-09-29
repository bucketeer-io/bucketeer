import { useTranslation } from 'i18n';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
  DropdownValue
} from 'components/dropdown';
import { GoalResultTab } from '.';

const ChartDataTypeDropdown = ({
  tab,
  chartType,
  onSelectOption
}: {
  tab: GoalResultTab;
  chartType: string;
  onSelectOption: (value: DropdownValue) => void;
}) => {
  const { t } = useTranslation(['table', 'common']);
  const evaluationOptions = [
    {
      label: t('results.evaluation-user'),
      value: 'evaluation-user'
    },
    {
      label: t('results.goal-total'),
      value: 'goal-total'
    },
    {
      label: t('results.goal-user'),
      value: 'goal-user'
    },
    {
      label: t('results.value-total'),
      value: 'value-total'
    },
    {
      label: t('results.value-user'),
      value: 'value-user'
    }
  ];
  const conversionOptions = [
    {
      label: t('results.conversion-rate'),
      value: 'conversion-rate'
    },
    {
      label: t('results.value-user'),
      value: 'value-user'
    }
  ];

  const options = tab === 'EVALUATION' ? evaluationOptions : conversionOptions;

  return (
    <DropdownMenu>
      <DropdownMenuTrigger
        isExpand
        label={options.find(item => item.value === chartType)?.label || ''}
        placeholder={t('common:select-value')}
        className="max-w-[528px]"
      />
      <DropdownMenuContent align="start">
        {options.map(item => (
          <DropdownMenuItem
            key={item.value}
            label={item.label}
            value={item.value}
            isSelectedItem={item.value === chartType}
            onSelectOption={onSelectOption}
          />
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

export default ChartDataTypeDropdown;
