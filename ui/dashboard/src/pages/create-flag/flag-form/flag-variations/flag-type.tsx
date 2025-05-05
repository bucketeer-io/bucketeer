import { useCallback } from 'react';
import { useFormContext } from 'react-hook-form';
import { useTranslation } from 'i18n';
import { cloneDeep } from 'lodash';
import { v4 as uuid } from 'uuid';
import { FeatureVariation, FeatureVariationType } from '@types';
import {
  IconFlagJSON,
  IconFlagNumber,
  IconFlagString,
  IconFlagSwitch,
  IconInfo
} from '@icons';
import { AddFlagForm } from 'pages/create-flag/form-schema';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import { Tooltip } from 'components/tooltip';

export const flagTypeOptions = [
  {
    label: 'Boolean',
    value: 'BOOLEAN',
    icon: IconFlagSwitch
  },
  {
    label: 'String',
    value: 'STRING',
    icon: IconFlagString
  },
  {
    label: 'Number',
    value: 'NUMBER',
    icon: IconFlagNumber
  },
  {
    label: 'JSON',
    value: 'JSON',
    icon: IconFlagJSON
  }
];

const defaultVariations: FeatureVariation[] = [
  {
    id: uuid(),
    value: 'true',
    name: '',
    description: ''
  },
  {
    id: uuid(),
    value: 'false',
    name: '',
    description: ''
  }
];

const FlagType = () => {
  const { t } = useTranslation(['form', 'common']);
  const { control, watch, setValue, setFocus, resetField } =
    useFormContext<AddFlagForm>();

  const variationType = watch('variationType');

  const currentFlagOption = flagTypeOptions.find(
    item => item.value === variationType
  );

  const handleOnChangeVariationType = useCallback(
    (
      value: FeatureVariationType,
      onChange: (value: FeatureVariationType) => void
    ) => {
      const cloneVariations = cloneDeep(defaultVariations);
      const newVariations =
        value === 'BOOLEAN'
          ? cloneVariations
          : cloneVariations.map(item => ({
              ...item,
              value: ''
            }));
      resetField('variations');
      setValue('variations', newVariations);
      setValue('defaultOnVariation', newVariations[0].id);
      setValue('defaultOffVariation', newVariations[1].id);

      onChange(value);
      let timerId: NodeJS.Timeout | null = null;
      if (timerId) clearTimeout(timerId);
      timerId = setTimeout(() => setFocus('variations.0.value'), 100);
    },
    [defaultVariations]
  );

  return (
    <Form.Field
      control={control}
      name={`variationType`}
      render={({ field }) => (
        <Form.Item className="py-0">
          <Form.Label required className="relative w-fit !mb-2">
            {t('feature-flags.flag-type')}
            <Tooltip
              align="start"
              alignOffset={-76}
              trigger={
                <div className="flex-center absolute top-0 -right-6">
                  <Icon icon={IconInfo} size={'sm'} color="gray-500" />
                </div>
              }
              content={t('flag-type-tooltip')}
              className="!z-[100] max-w-[400px]"
            />
          </Form.Label>
          <Form.Control>
            <DropdownMenu>
              <DropdownMenuTrigger
                placeholder={t(`feature-flags.flag-type`)}
                trigger={
                  <div className="flex items-center gap-x-2">
                    {currentFlagOption?.icon && (
                      <Icon
                        icon={currentFlagOption?.icon}
                        size={'md'}
                        className="flex-center"
                      />
                    )}
                    <p>{currentFlagOption?.label}</p>
                  </div>
                }
                variant="secondary"
                className="w-full"
              />
              <DropdownMenuContent
                className="w-[502px]"
                align="start"
                {...field}
              >
                {flagTypeOptions.map((item, index) => (
                  <DropdownMenuItem
                    {...field}
                    key={index}
                    icon={item.icon}
                    value={item.value}
                    label={item.label}
                    onSelectOption={value =>
                      handleOnChangeVariationType(
                        value as FeatureVariationType,
                        field.onChange
                      )
                    }
                  />
                ))}
              </DropdownMenuContent>
            </DropdownMenu>
          </Form.Control>
          <Form.Message />
        </Form.Item>
      )}
    />
  );
};

export default FlagType;
