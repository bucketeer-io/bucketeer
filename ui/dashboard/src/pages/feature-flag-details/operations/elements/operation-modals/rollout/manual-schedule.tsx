import { useCallback, useMemo } from 'react';
import { useFieldArray, useFormContext } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { getLanguage, useTranslation } from 'i18n';
import { areIntervalsApart } from 'utils/function';
import { cn } from 'utils/style';
import { IconInfo, IconPlus, IconTrash } from '@icons';
import { RolloutSchemaType } from 'pages/feature-flag-details/operations/form-schema';
import {
  handleCreateIncrement,
  numberToJapaneseOrdinal,
  numberToOrdinalWord
} from 'pages/feature-flag-details/operations/utils';
import { ScheduleItem } from 'pages/feature-flag-details/types';
import Button from 'components/button';
import { ReactDatePicker } from 'components/date-time-picker';
import { DropdownOption } from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import InputGroup from 'components/input-group';
import { Tooltip } from 'components/tooltip';
import DropdownMenuWithSearch from 'elements/dropdown-with-search';

const ManualSchedule = ({
  variationOptions,
  isDisableCreateRollout
}: {
  variationOptions: DropdownOption[];
  isDisableCreateRollout: boolean;
}) => {
  const { t } = useTranslation(['form']);
  const isLanguageJapanese = getLanguage() === 'ja';

  const { control, watch, clearErrors } = useFormContext<RolloutSchemaType>();
  const watchScheduleList = [
    ...watch('progressiveRollout.manual.schedulesList')
  ];
  const {
    fields: schedulesList,
    append,
    remove
  } = useFieldArray({
    name: 'progressiveRollout.manual.schedulesList',
    control,
    keyName: 'scheduleItemId'
  });

  const isLastScheduleWeight100 = useMemo(
    () => Number(watchScheduleList.at(-1)?.weight) === 100,
    [watchScheduleList]
  );

  const handleAddIncrement = useCallback(() => {
    const newIncrement = handleCreateIncrement({
      lastSchedule: watchScheduleList.at(-1) as ScheduleItem,
      incrementType: 'minute',
      addValue: 5
    });
    append(newIncrement);
  }, [watchScheduleList]);

  const handleRemoveSchedule = (index: number) => {
    remove(index);
    const filterScheduleList = watchScheduleList.filter(
      (_, scheduleIndex) => scheduleIndex !== index
    );
    const isValidIntervals = areIntervalsApart(
      filterScheduleList.map(item => item.executeAt.getTime()),
      5
    );
    if (isValidIntervals)
      clearErrors('progressiveRollout.manual.schedulesList');
  };

  return (
    <div className="flex flex-col w-full gap-y-4">
      <Form.Field
        control={control}
        name={`progressiveRollout.manual.variationId`}
        render={({ field }) => (
          <Form.Item className="py-0">
            <Form.Label required className="relative w-fit">
              {t('flag-variation')}
              <Tooltip
                align="start"
                alignOffset={-73}
                content={t('rollout-tooltips.manual.variation')}
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

      {schedulesList.map((item, index) => (
        <div key={item.scheduleItemId} className="flex w-full gap-x-4">
          <Form.Field
            name={`progressiveRollout.manual.schedulesList.${index}.weight`}
            render={({ field }) => (
              <Form.Item className="flex flex-col py-0 flex-1 size-full">
                <Form.Label
                  required
                  className={cn('relative w-fit', {
                    capitalize: !isLanguageJapanese
                  })}
                >
                  <Trans
                    i18nKey={'form:ordinal-increment'}
                    values={{
                      ordinal: isLanguageJapanese
                        ? numberToJapaneseOrdinal(index + 1)
                        : numberToOrdinalWord(index + 1)
                    }}
                  />
                  <Tooltip
                    align="start"
                    alignOffset={-73}
                    content={t('rollout-tooltips.manual.weight')}
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
                    className="w-full max-w-full"
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
                <Form.Message className="max-w-full break-all" />
              </Form.Item>
            )}
          />
          <Form.Field
            control={control}
            name={`progressiveRollout.manual.schedulesList.${index}.executeAt`}
            render={({ field }) => (
              <Form.Item className="flex flex-col flex-1 py-0 size-full">
                <Form.Label required className="relative w-fit">
                  {t('feature-flags.start-date')}
                  <Tooltip
                    align="start"
                    alignOffset={-120}
                    content={t('rollout-tooltips.manual.execute-at')}
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
                <Form.Message className="max-w-full break-all" />
              </Form.Item>
            )}
          />
          <Button
            type="button"
            variant={'grey'}
            className="flex-center self-start h-full mt-9 min-w-5"
            disabled={schedulesList.length <= 1 || isDisableCreateRollout}
            onClick={() => handleRemoveSchedule(index)}
          >
            <Icon icon={IconTrash} size={'sm'} />
          </Button>
        </div>
      ))}
      <Button
        type="button"
        variant={'text'}
        className="w-fit px-0 h-6"
        disabled={isLastScheduleWeight100 || isDisableCreateRollout}
        onClick={handleAddIncrement}
      >
        <Icon icon={IconPlus} size={'md'} />
        {t('add-increment')}
      </Button>
    </div>
  );
};

export default ManualSchedule;
