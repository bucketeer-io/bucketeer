import { useMemo } from 'react';
import { Trans } from 'react-i18next';
import { useTranslation } from 'i18n';
import { IconCollapse, IconExpand } from '@icons';
import Button from 'components/button';
import Dropdown from 'components/dropdown';
import Icon from 'components/icon';
import { GroupByType } from '../page-content';

interface Props {
  isExpandAll: boolean;
  groupBy: GroupByType;
  onChangeGroupBy: (val: GroupByType) => void;
  onResetFields: () => void;
  onEditFields: () => void;
  onToggleExpandAll: () => void;
}

const ActionBar = ({
  isExpandAll,
  groupBy,
  onChangeGroupBy,
  onResetFields,
  onEditFields,
  onToggleExpandAll
}: Props) => {
  const { t } = useTranslation(['common']);

  const groupByOptions = useMemo(
    () => [
      {
        label: t('flag'),
        value: 'FLAG'
      },
      {
        label: t('user'),
        value: 'USER'
      }
    ],
    []
  );

  return (
    <div className="flex items-center w-full justify-between gap-x-4">
      <p className="typo-head-bold-small text-gray-800 whitespace-nowrap">
        {t('debugger-results')}
      </p>
      <div className="flex items-center gap-x-4">
        <Dropdown
          labelCustom={
            <Trans
              i18nKey="common:group-by-type"
              values={{
                type: t(groupBy === 'FLAG' ? 'flag' : 'user')
              }}
            />
          }
          value={groupBy}
          options={groupByOptions}
          onChange={value => onChangeGroupBy(value as GroupByType)}
          alignContent="end"
          contentClassName="w-[173px]"
        />

        <Button
          variant={'secondary'}
          className="max-w-[154px]"
          onClick={onToggleExpandAll}
        >
          <Icon
            icon={isExpandAll ? IconCollapse : IconExpand}
            size="sm"
            color="primary-500"
          />
          {t(isExpandAll ? 'collapse-all' : 'expand-all')}
        </Button>
        <Button variant="secondary" onClick={onEditFields}>
          {t('edit-fields')}
        </Button>
        <Button onClick={onResetFields}>{t('clear-all-fields')}</Button>
      </div>
    </div>
  );
};

export default ActionBar;
