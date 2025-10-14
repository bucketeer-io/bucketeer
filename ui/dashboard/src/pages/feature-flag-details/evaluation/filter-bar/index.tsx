import { useTranslation } from 'i18n';
import { EvaluationTimeRange } from '@types';
import { IconInfo, IconThreeLines } from '@icons';
import Dropdown, { DropdownOption } from 'components/dropdown';
import Icon from 'components/icon';
import { Tooltip } from 'components/tooltip';
import { TimeRangeOption } from '../types';

const FilterBar = ({
  isLoading,
  timeRangeOptions,
  timeRangeLabel,
  currentFilter,
  onChangeTimeRange
}: {
  isLoading: boolean;
  timeRangeOptions: TimeRangeOption[];
  timeRangeLabel: string;
  currentFilter: string;
  onChangeTimeRange: (timeRange: EvaluationTimeRange) => void;
}) => {
  const { t } = useTranslation(['common', 'table']);

  return (
    <div className="flex items-center w-full justify-between">
      <div className="flex items-center gap-x-2">
        <p className="typo-head-bold-small text-gray-800">{t('evaluation')}</p>
        <Tooltip
          align="start"
          alignOffset={-90}
          content={t('table:evaluation.tooltip-content')}
          trigger={
            <div className="flex-center -mb-1">
              <Icon icon={IconInfo} size="xxs" color="gray-500" />
            </div>
          }
          className="max-w-[310px]"
        />
      </div>
      <Dropdown
        trigger={
          <div className="flex items-center gap-x-2">
            <Icon icon={IconThreeLines} size="sm" />
            <p className="text-gray-600">{timeRangeLabel}</p>
          </div>
        }
        value={currentFilter}
        options={timeRangeOptions as DropdownOption[]}
        disabled={isLoading}
        showArrow={false}
        onChange={value => onChangeTimeRange(value as EvaluationTimeRange)}
        alignContent="end"
        className="px-4 py-[13.5px]"
        wrapTriggerStyle="!w-fit"
      />
    </div>
  );
};

export default FilterBar;
