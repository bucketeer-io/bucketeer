import { useCallback, useEffect, useMemo } from 'react';
import { useFieldArray, useFormContext } from 'react-hook-form';
import { useTranslation } from 'i18n';
import { IconInfo, IconPlus, IconTrash, IconWatch } from '@icons';
import { OperationForm } from 'pages/feature-flag-details/operations/form-schema';
import { ActionTypeMap } from 'pages/feature-flag-details/operations/types';
import { createDatetimeClausesList } from 'pages/feature-flag-details/operations/utils';
import Button from 'components/button';
import { ReactDatePicker } from 'components/date-time-picker';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';

const ScheduleList = ({ isDateSorted }: { isDateSorted: boolean }) => {
  const { t } = useTranslation(['form', 'common']);
  const { control, clearErrors, setError } = useFormContext<OperationForm>();
  const {
    fields: scheduleData,
    append,
    remove
  } = useFieldArray({
    name: 'datetimeClausesList',
    control,
    keyName: 'scheduleId'
  });

  const stateOptions = useMemo(
    () => [
      {
        label: t('form:experiments.on'),
        value: ActionTypeMap.ENABLE
      },
      {
        label: t('form:experiments.off'),
        value: ActionTypeMap.DISABLE
      }
    ],
    [scheduleData]
  );

  const handleAddDate = useCallback(() => {
    const lastTime = scheduleData.at(-1)?.time?.getTime();
    const dateTimeClause = createDatetimeClausesList(lastTime);
    append({
      ...dateTimeClause
    });
  }, [scheduleData]);

  useEffect(() => {
    if (!isDateSorted && scheduleData.length > 1) {
      setError('datetimeClausesList', {
        message: 'test message'
      });
    } else {
      clearErrors();
    }
  }, [isDateSorted, scheduleData]);

  return (
    <>
      <p className="typo-head-bold-small text-gray-800">
        {t('feature-flags.schedule')}
      </p>
      <Form.Field
        control={control}
        name="datetimeClausesList"
        render={() => (
          <Form.Item className="flex flex-col gap-y-4 py-0">
            <Form.Control>
              <>
                {scheduleData.map((item, index) => (
                  <div className="flex w-full gap-x-4" key={item.scheduleId}>
                    <Form.Field
                      name={`datetimeClausesList.${index}.actionType`}
                      control={control}
                      render={({ field }) => (
                        <Form.Item className="py-0">
                          <Form.Label required className="relative w-fit">
                            {t('common:state')}
                            <Icon
                              icon={IconInfo}
                              size="xs"
                              color="gray-500"
                              className="absolute -right-6"
                            />
                          </Form.Label>
                          <Form.Control>
                            <DropdownMenu>
                              <DropdownMenuTrigger
                                label={
                                  stateOptions.find(
                                    item => item.value === field.value
                                  )?.label
                                }
                                className="w-[124px] uppercase"
                              />
                              <DropdownMenuContent
                                align="start"
                                className="min-w-[124px]"
                              >
                                {stateOptions.map(({ label, value }, index) => (
                                  <DropdownMenuItem
                                    key={index}
                                    label={label}
                                    value={value}
                                    onSelectOption={value =>
                                      field.onChange(value)
                                    }
                                    className="uppercase"
                                  />
                                ))}
                              </DropdownMenuContent>
                            </DropdownMenu>
                          </Form.Control>
                          <Form.Message />
                        </Form.Item>
                      )}
                    />
                    <Form.Field
                      name={`datetimeClausesList.${index}.time`}
                      control={control}
                      render={({ field }) => (
                        <>
                          <Form.Item className="py-0">
                            <Form.Control>
                              <div className="flex gap-x-4">
                                <div>
                                  <Form.Label required>
                                    {t('feature-flags.start-date')}
                                  </Form.Label>
                                  <ReactDatePicker
                                    dateFormat={'yyyy/MM/dd'}
                                    selected={field.value ?? null}
                                    showTimeSelect={false}
                                    className="w-[186px]"
                                    onChange={date => {
                                      if (date) field.onChange(date);
                                    }}
                                  />
                                </div>
                                <div>
                                  <Form.Label required>
                                    {t('feature-flags.time')}
                                  </Form.Label>
                                  <Form.Control>
                                    <ReactDatePicker
                                      dateFormat={'HH:mm'}
                                      selected={field.value ?? null}
                                      showTimeSelectOnly={true}
                                      className="w-[124px]"
                                      onChange={date => {
                                        if (date) field.onChange(date);
                                      }}
                                      icon={
                                        <Icon
                                          icon={IconWatch}
                                          className="flex-center"
                                        />
                                      }
                                    />
                                  </Form.Control>
                                </div>
                              </div>
                            </Form.Control>
                            <Form.Message />
                          </Form.Item>
                        </>
                      )}
                    />
                    <Button
                      variant={'grey'}
                      size={'icon-sm'}
                      className="self-end mb-2"
                      disabled={scheduleData.length <= 1}
                      onClick={() => remove(index)}
                    >
                      <Icon icon={IconTrash} size={'sm'} />
                    </Button>
                  </div>
                ))}
              </>
            </Form.Control>
            <Form.Message />
          </Form.Item>
        )}
      />
      <Button
        variant={'text'}
        className="flex items-center h-6 self-start p-0"
        onClick={handleAddDate}
      >
        <Icon
          icon={IconPlus}
          color="primary-500"
          size="md"
          className="flex-center"
        />
        {t('feature-flags.date')}
      </Button>
    </>
  );
};

export default ScheduleList;
