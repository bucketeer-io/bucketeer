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

  const flagOptions = [
    {
      label: 'Flag 1',
      value: 'flag-1'
    },
    {
      label: 'Flag 2',
      value: 'flag-2'
    }
  ];
  
  return (
    <div className="flex flex-col w-full gap-y-3 typo-para-small leading-[14px] text-gray-600">
      <div className="flex items-center w-full gap-x-2">
        <p>The</p>
        <DropdownMenu>
          <DropdownMenuTrigger
            placeholder={t(`experiments.select-flag`)}
            label={''}
            variant="secondary"
            className="w-full"
          />
          <DropdownMenuContent className="w-[502px]" align="start" {...field}>
            {flagOptions.map((item, index) => (
              <DropdownMenuItem
                {...field}
                key={index}
                value={item.value}
                label={item.label}
                onSelectOption={value => {
                  field.onChange(value);
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
