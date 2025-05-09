import { useTranslation } from 'i18n';
import { IconInfo } from '@icons';
import { ReactDateRangePicker } from 'components/date-range-picker';
import Icon from 'components/icon';
import { Tooltip } from 'components/tooltip';

const FilterBar = () => {
  const { t } = useTranslation(['common']);

  return (
    <div className="flex items-center w-full justify-between">
      <div className="flex items-center gap-x-2">
        <p className="typo-head-bold-small text-gray-800">{t('evaluation')}</p>
        <Tooltip
          // Need to update
          content="Evaluation Content"
          trigger={
            <div className="flex-center -mb-1">
              <Icon icon={IconInfo} size="xxs" color="gray-500" />
            </div>
          }
        />
      </div>
      <ReactDateRangePicker onChange={value => console.log(value)} />
    </div>
  );
};

export default FilterBar;
