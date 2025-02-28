import { useTranslation } from 'react-i18next';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import { DefineAudienceProps } from '.';

const DefineAudienceRule = ({ field }: DefineAudienceProps) => {
  const { t } = useTranslation(['form', 'common']);

  const ruleOptions = [
    {
      label: 'Rule 1',
      value: 'rule-1'
    },
    {
      label: 'Rule 2',
      value: 'rule-2'
    }
  ];

  return (
    <div className="flex flex-col w-full gap-y-3 typo-para-small leading-[14px] text-gray-600">
      <div className="flex items-center w-full gap-x-2">
        <p>The</p>
        <DropdownMenu>
          <DropdownMenuTrigger
            placeholder={t(`experiments.select-rule`)}
            label={
              ruleOptions.find(item => item.value === field.value?.rule)
                ?.label || ''
            }
            variant="secondary"
            className="w-full"
          />
          <DropdownMenuContent className="w-[502px]" align="start" {...field}>
            {ruleOptions.map((item, index) => (
              <DropdownMenuItem
                {...field}
                key={index}
                value={item.value}
                label={item.label}
                onSelectOption={value => {
                  field.onChange({
                    ...field.value,
                    rule: value
                  });
                }}
              />
            ))}
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
      <p>{t('experiments.define-audience.rule-defines')}</p>
      <p className="typo-para-medium leading-5 text-gray-500">
        {t('experiments.define-audience.any-traffic')}
      </p>
    </div>
  );
};

export default DefineAudienceRule;
