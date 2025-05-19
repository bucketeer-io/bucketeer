import { useMemo } from 'react';
import { useFormContext } from 'react-hook-form';
import { useTranslation } from 'i18n';
import { IconInfo } from '@icons';
import { RolloutSchemaType } from 'pages/feature-flag-details/operations/form-schema';
import { ReactDatePicker } from 'components/date-time-picker';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
  DropdownOption
} from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import InputGroup from 'components/input-group';
import { Tooltip } from 'components/tooltip';
import DropdownMenuWithSearch from 'elements/dropdown-with-search';

const TemplateSchedule = ({
  variationOptions,
  isDisableCreateRollout
}: {
  variationOptions: DropdownOption[];
  isDisableCreateRollout: boolean;
}) => {
  const { t } = useTranslation(['form', 'table', 'common']);
  const { control } = useFormContext<RolloutSchemaType>();

  const intervalOptions = useMemo(
    () => [
      {
        label: t('hourly'),
        value: 'HOURLY'
      },
      {
        label: t('daily'),
        value: 'DAILY'
      },
      {
        label: t('weekly'),
        value: 'WEEKLY'
      }
    ],
    []
  );

  return (
    <div className="flex flex-col w-full gap-y-4">
      <Form.Field
        control={control}
        name={`progressiveRollout.template.variationId`}
        render={({ field }) => (
          <Form.Item className="py-0">
            <Form.Label required className="relative w-fit">
              {t('table:results.variation')}
              <Tooltip
                align="start"
                alignOffset={-73}
                content={t('rollout-tooltips.template.variation')}
                trigger={
                  <div className="flex-center size-fit absolute top-0 -right-6">
                    <Icon icon={IconInfo} size="xs" color="gray-500" />
                  </div>
                }
                className="max-w-[300px]"
              />
            </Form.Label>
            <Form.Control>
              <DropdownMenuWithSearch
                align="end"
                label={
                  variationOptions.find(item => item.value === field.value)
                    ?.label || ''
                }
                contentClassName="[&>div.wrapper-menu-items>div]:px-4"
                options={variationOptions}
                disabled={isDisableCreateRollout}
                onSelectOption={field.onChange}
              />
            </Form.Control>
            <Form.Message />
          </Form.Item>
        )}
      />

      <Form.Field
        control={control}
        name={`progressiveRollout.template.startDate`}
        render={({ field }) => (
          <Form.Item className="py-0">
            <Form.Label required className="relative w-fit">
              {t('feature-flags.start-date')}
              <Tooltip
                align="start"
                alignOffset={-73}
                content={t('rollout-tooltips.template.start-date')}
                trigger={
                  <div className="flex-center size-fit absolute top-0 -right-6">
                    <Icon icon={IconInfo} size="xs" color="gray-500" />
                  </div>
                }
                className="max-w-[300px]"
              />
            </Form.Label>
            <Form.Control>
              <ReactDatePicker
                selected={field.value ?? null}
                disabled={isDisableCreateRollout}
                onChange={date => {
                  if (date) {
                    field.onChange(date);
                  }
                }}
              />
            </Form.Control>
            <Form.Message />
          </Form.Item>
        )}
      />
      <div className="flex w-full gap-x-4">
        <Form.Field
          control={control}
          name={`progressiveRollout.template.increments`}
          render={({ field }) => (
            <Form.Item className="py-0 flex-1 h-full">
              <Form.Label required className="relative w-fit">
                {t('increments')}
                <Tooltip
                  align="start"
                  alignOffset={-73}
                  content={t('rollout-tooltips.template.increment')}
                  trigger={
                    <div className="flex-center size-fit absolute top-0 -right-6">
                      <Icon icon={IconInfo} size="xs" color="gray-500" />
                    </div>
                  }
                  className="max-w-[300px]"
                />
              </Form.Label>
              <Form.Control>
                <InputGroup
                  className="w-full"
                  addonSlot="right"
                  addonSize="md"
                  addon={'%'}
                >
                  <Input
                    {...field}
                    value={field.value || ''}
                    type="number"
                    className="pr-8"
                    disabled={isDisableCreateRollout}
                    onWheel={e => {
                      e.currentTarget.blur();
                    }}
                    onKeyDown={e => {
                      if (e.key === 'Enter') e.preventDefault();
                    }}
                  />
                </InputGroup>
              </Form.Control>
              <Form.Message />
            </Form.Item>
          )}
        />
        <Form.Field
          control={control}
          name={`progressiveRollout.template.interval`}
          render={({ field }) => (
            <Form.Item className="py-0 flex-1 h-full">
              <Form.Label required className="relative w-fit">
                {t('frequency')}
                <Tooltip
                  align="start"
                  alignOffset={-150}
                  content={t('rollout-tooltips.template.frequency')}
                  trigger={
                    <div className="flex-center size-fit absolute top-0 -right-6">
                      <Icon icon={IconInfo} size="xs" color="gray-500" />
                    </div>
                  }
                  className="max-w-[300px]"
                />
              </Form.Label>
              <Form.Control>
                <DropdownMenu>
                  <DropdownMenuTrigger
                    label={
                      intervalOptions.find(item => item.value === field.value)
                        ?.label || ''
                    }
                    isExpand
                    disabled={isDisableCreateRollout}
                  />
                  <DropdownMenuContent align="end" className="min-w-[243px]">
                    {intervalOptions.map((item, index) => (
                      <DropdownMenuItem
                        key={index}
                        label={item.label}
                        value={item.value}
                        onSelectOption={field.onChange}
                      />
                    ))}
                  </DropdownMenuContent>
                </DropdownMenu>
              </Form.Control>
            </Form.Item>
          )}
        />
      </div>
    </div>
  );
};

export default TemplateSchedule;
