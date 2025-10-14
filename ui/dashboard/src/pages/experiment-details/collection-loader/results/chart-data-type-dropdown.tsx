import { useTranslation } from 'i18n';
import Dropdown, { DropdownValue } from 'components/dropdown';
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
    <Dropdown
      options={options}
      value={chartType}
      placeholder={t('common:select-value')}
      onChange={val => onSelectOption(val as DropdownValue)}
      className="max-w-[528px]"
    />
  );
};

export default ChartDataTypeDropdown;
