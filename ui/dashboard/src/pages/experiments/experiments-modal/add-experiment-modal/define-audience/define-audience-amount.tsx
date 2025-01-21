import { useRef, useState } from 'react';
import { Trans, useTranslation } from 'react-i18next';
import { cn } from 'utils/style';
import Divider from 'components/divider';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Input from 'components/input';
import { DefineAudienceProps } from '.';

const experimentOptions = [
  {
    label: '5%',
    value: 5
  },
  {
    label: '10%',
    value: 10
  },
  {
    label: '50%',
    value: 50
  },
  {
    label: '90%',
    value: 90
  },
  {
    label: 'Custom',
    value: 'custom'
  }
];

const servedOptions = [
  {
    label: 'True',
    value: 1
  },
  {
    label: 'False',
    value: 0
  }
];

const ExperimentSelect = ({
  label,
  value,
  isActive,
  onSelect
}: {
  label: string;
  value: string | number;
  isActive: boolean;
  onSelect: (value: string | number) => void;
}) => {
  return (
    <div
      className={cn(
        'flex-center size-fit min-w-[53px] py-[14px] px-3 border border-gray-400 rounded-lg typo-para-medium leading-5 text-gray-700 capitalize cursor-pointer',
        {
          'text-primary-500 border-primary-500': isActive
        }
      )}
      onClick={() => onSelect(value)}
    >
      {label}
    </div>
  );
};

const DefineAudienceAmount = ({ field }: DefineAudienceProps) => {
  const { t } = useTranslation(['form', 'common']);
  const [isCustomExperiment, setIsCustomExperiment] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  const handleSelectExperiment = (value: string | number) => {
    if (typeof value === 'string') {
      let timerId: NodeJS.Timeout | null = null;
      if (timerId) clearTimeout(timerId);

      setIsCustomExperiment(true);
      timerId = setTimeout(() => inputRef.current?.focus(), 100);
      return field.onChange({
        ...field.value,
        inExperiment: 0,
        notInExperiment: 100
      });
    }
    setIsCustomExperiment(false);
    return field.onChange({
      ...field.value,
      inExperiment: value,
      notInExperiment: 100 - +value
    });
  };
  return (
    <div className="w-full p-4 bg-gray-100 rounded-lg">
      <div className="flex flex-col w-full gap-y-4 typo-para-small leading-[14px] text-gray-600">
        <p>{t('experiments.define-audience.audience-amount')}</p>
        <div className="w-full h-3 p-[1px] border border-gray-400 rounded-full bg-gray-100">
          <div
            className={cn('h-full bg-primary-500 rounded-l-full', {
              'rounded-r-full': field.value?.inExperiment === 100
            })}
            style={{
              width: `${field.value?.inExperiment}%`
            }}
          />
        </div>
        <div className="flex items-center w-full gap-x-4">
          <div className="flex items-center gap-x-2">
            <div className="flex-center size-5 m-0.5 rounded bg-primary-500" />
            <p>{`${t('experiments.define-audience.in-this-experiment')} - ${field.value?.inExperiment}%`}</p>
          </div>
          <div className="flex items-center gap-x-2">
            <div className="flex-center size-5 m-0.5 border border-gray-400 rounded bg-gray-100" />
            <p>{`${t('experiments.define-audience.not-in-experiment')} - ${field.value?.notInExperiment}%`}</p>
          </div>
        </div>
      </div>
      <Divider className="my-5" />
      <div className="flex items-center w-full gap-x-2">
        <p className="typo-para-medium leading-5 text-gray-600 whitespace-nowrap">
          {t('experiments.define-audience.in-this-experiment')}:
        </p>
        {experimentOptions.map((item, index) => (
          <ExperimentSelect
            key={index}
            label={item.label}
            value={item.value}
            isActive={
              field.value?.inExperiment === item.value ||
              (item.value === 'custom' && isCustomExperiment)
            }
            onSelect={handleSelectExperiment}
          />
        ))}
      </div>
      {isCustomExperiment && (
        <Input
          ref={inputRef}
          type="number"
          className="mt-5"
          value={field.value?.inExperiment ? field.value?.inExperiment : ''}
          onWheel={e => e.currentTarget.blur()}
          onChange={value =>
            field.onChange({
              ...field.value,
              inExperiment: +value >= 100 ? 100 : +value < 0 ? 0 : +value,
              notInExperiment:
                +value >= 100 ? 0 : +value < 0 ? 100 : 100 - +value
            })
          }
        />
      )}
      <div className="flex items-center w-full gap-x-2 mt-4 typo-para-medium leading-5 text-gray-600 whitespace-nowrap">
        <Trans
          i18nKey={'form:experiments.define-audience.not-in-experiment-served'}
          values={{
            percent: `${field.value?.notInExperiment}%`
          }}
          components={{
            highlight: (
              <div className="flex-center size-fit p-3 rounded-lg typo-para-medium leading-5 text-gray-700 bg-gray-200" />
            )
          }}
        />
        <DropdownMenu>
          <DropdownMenuTrigger
            trigger={
              <div className="flex items-center justify-between w-full gap-x-2">
                <div className="flex items-center gap-x-2">
                  <div
                    className={cn('flex-center size-3 rounded-sm rotate-45', {
                      'bg-accent-blue-500': field.value?.served,
                      'bg-accent-pink-500': !field.value?.served
                    })}
                  />
                  <p className="typo-para-medium leading-5 text-gray-600">
                    {field.value?.served ? 'True' : 'False'}
                  </p>
                </div>
              </div>
            }
            className="!min-w-[120px]"
          />
          <DropdownMenuContent className="w-[100px]" align="end" {...field}>
            {servedOptions.map((item, index) => (
              <DropdownMenuItem
                {...field}
                key={index}
                value={item.value}
                label={item.label}
                onSelectOption={value => {
                  field.onChange({
                    ...field.value,
                    served: !!value
                  });
                }}
              />
            ))}
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
      <div className="flex items-center w-full gap-x-2 mt-4 typo-para-medium leading-5 text-gray-600 whitespace-nowrap">
        <Trans
          i18nKey={'form:experiments.define-audience.in-experiment-target'}
          values={{
            percent: `${field.value?.inExperiment}%`
          }}
          components={{
            highlight: (
              <div className="flex-center size-fit p-3 rounded-lg typo-para-medium leading-5 text-gray-700 bg-gray-200" />
            )
          }}
        />
      </div>
    </div>
  );
};

export default DefineAudienceAmount;
