import { useMemo } from 'react';
import { Trans } from 'react-i18next';
import { useTranslation } from 'i18n';
import Button from 'components/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import { GroupByType } from '.';

interface Props {
  groupBy: GroupByType;
  setGroupBy: (val: GroupByType) => void;
  onResetFields: () => void;
  onEditFields: () => void;
}

const ActionBar = ({
  groupBy,
  setGroupBy,
  onResetFields,
  onEditFields
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
      <p className="typo-head-bold-small text-gray-800">
        {t('debugger-results')}
      </p>
      <div className="flex items-center gap-x-4">
        <DropdownMenu>
          <DropdownMenuTrigger
            label={
              <Trans
                i18nKey="common:group-by-type"
                values={{
                  type: t(groupBy === 'FLAG' ? 'flag' : 'user')
                }}
              />
            }
          />
          <DropdownMenuContent align="end" className="min-w-[173px]">
            {groupByOptions.map((item, index) => (
              <DropdownMenuItem
                key={index}
                label={item.label}
                value={item.value}
                onSelectOption={value => setGroupBy(value as GroupByType)}
              />
            ))}
          </DropdownMenuContent>
        </DropdownMenu>

        <Button variant="secondary" onClick={onEditFields}>
          {t('edit-fields')}
        </Button>
        <Button onClick={onResetFields}>{t('clear-all-fields')}</Button>
      </div>
    </div>
  );
};

export default ActionBar;
