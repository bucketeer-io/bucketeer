import { forwardRef, Ref } from 'react';
import { useFormContext } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { useQueryUserSegments } from '@queries/user-segments';
import { getCurrentEnvironment, useAuth } from 'auth';
import { LIST_PAGE_SIZE } from 'constants/app';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import { IconTrash } from '@icons';
import {
  booleanVariations,
  flagOptions,
  jsonVariations,
  numberVariations,
  stringVariations
} from 'pages/feature-flag-details/mocks';
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
import Input from 'components/input';
import { SegmentConditionType, SituationType } from '../types';

interface Props {
  type: 'if' | 'and';
  condition: SegmentConditionType;
  isDisabledDelete: boolean;
  segmentIndex: number;
  ruleIndex: number;
  conditionIndex: number;
  onDeleteCondition: () => void;
  onChangeFormField: (field: string, value: string | number | boolean) => void;
}

const ConditionForm = forwardRef(
  (
    {
      type,
      condition,
      isDisabledDelete,
      segmentIndex,
      ruleIndex,
      conditionIndex,
      onDeleteCondition,
      onChangeFormField
    }: Props,
    ref: Ref<HTMLDivElement>
  ) => {
    const { t } = useTranslation(['form', 'common', 'table']);
    const { consoleAccount } = useAuth();
    const currentEnvironment = getCurrentEnvironment(consoleAccount!);

    const isCompare = condition.situation === 'compare';
    const isUserSegment = condition.situation === 'user-segment';
    const isDate = condition.situation === 'date';
    const isFlag = condition.situation === 'feature-flag';

    const methods = useFormContext();
    const { control, watch, setValue } = methods;

    const situationOptions = [
      {
        label: t('feature-flags.compare'),
        value: 'compare'
      },
      {
        label: t('feature-flags.user-segment'),
        value: 'user-segment'
      },
      {
        label: t('feature-flags.date'),
        value: 'date'
      },
      {
        label: t('feature-flags.feature-flag'),
        value: 'feature-flag'
      }
    ];

    const conditionerDateOptions = [
      {
        label: 'Before',
        value: 'before'
      },
      {
        label: 'After',
        value: 'after'
      }
    ];

    const { data: segmentCollection, isLoading: segmentLoading } =
      useQueryUserSegments({
        params: {
          cursor: String(0),
          pageSize: LIST_PAGE_SIZE,
          environmentId: currentEnvironment.id
        },
        enabled: isUserSegment
      });

    const userSegments = segmentCollection?.segments || [];

    const flagId = watch('flagId');

    const isStringVariation = flagId?.includes('string');
    const isNumberVariation = flagId?.includes('number');
    const isBooleanVariation = flagId?.includes('boolean');

    const variationOptions = isStringVariation
      ? stringVariations
      : isNumberVariation
        ? numberVariations
        : isBooleanVariation
          ? booleanVariations
          : jsonVariations;

    const commonName = `targetSegmentRules.${segmentIndex}.rules.${ruleIndex}.conditions.${conditionIndex}.`;

    return (
      <div ref={ref} className="flex items-center w-full gap-x-4">
        <div
          className={cn(
            'flex-center w-[42px] h-[26px] rounded-[3px] typo-para-small leading-[14px]',
            {
              'bg-accent-pink-50 text-accent-pink-500': type === 'if',
              'bg-gray-200 text-gray-600': type === 'and'
            }
          )}
        >
          {type === 'if' ? t('common:if') : t('common:and')}
        </div>
        <div className="flex items-center w-full flex-1 pl-4 border-l border-primary-500">
          <div className="flex items-end w-full gap-x-4">
            <Form.Field
              control={control}
              name={`${commonName}situation`}
              render={({ field }) => (
                <Form.Item className="flex flex-col flex-1 self-stretch py-0 min-w-[170px] order-1">
                  <Form.Label required>
                    {t('feature-flags.situation')}
                  </Form.Label>
                  <Form.Control>
                    <DropdownMenu>
                      <DropdownMenuTrigger
                        label={
                          situationOptions.find(
                            item => item.value === condition.situation
                          )?.label
                        }
                        className="w-full"
                      />
                      <DropdownMenuContent align="start" {...field}>
                        {situationOptions.map((item, index) => (
                          <DropdownMenuItem
                            key={index}
                            label={item.label}
                            value={item.value}
                            onSelectOption={value => {
                              field.onChange(value);
                              onChangeFormField('situation', value);
                              const isEqualConditioner = [
                                'compare',
                                'feature-flag'
                              ].includes(value as SituationType);
                              const _isDate = value === 'date';
                              onChangeFormField(
                                'conditioner',
                                isEqualConditioner
                                  ? '='
                                  : _isDate
                                    ? 'before'
                                    : 'Is included in'
                              );
                              setValue(
                                `${commonName}conditioner`,
                                isEqualConditioner
                                  ? '='
                                  : _isDate
                                    ? 'before'
                                    : 'Is included in'
                              );
                            }}
                          />
                        ))}
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            {isCompare && (
              <Form.Field
                control={control}
                name={`${commonName}firstValue`}
                render={({ field }) => (
                  <Form.Item className="flex flex-col flex-1 self-stretch py-0 min-w-[170px] order-2">
                    <Form.Label required>
                      <Trans
                        i18nKey={'form:feature-flags.value-type'}
                        values={{
                          type: 'First'
                        }}
                      />
                    </Form.Label>
                    <Form.Control>
                      <Input
                        {...field}
                        value={field.value || condition.firstValue}
                        onChange={value => {
                          field.onChange(value);
                          onChangeFormField('firstValue', value);
                        }}
                      />
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />
            )}
            {isFlag && (
              <Form.Field
                control={control}
                name={`${commonName}flagId`}
                render={({ field }) => (
                  <Form.Item className="flex flex-col flex-1 self-stretch py-0 min-w-[170px] order-2">
                    <Form.Label required>
                      {t('feature-flags.feature-flag')}
                    </Form.Label>
                    <Form.Control>
                      <DropdownMenu>
                        <DropdownMenuTrigger
                          label={
                            flagOptions?.find(
                              item =>
                                item.value ===
                                (field.value || condition?.flagId)
                            )?.label || ''
                          }
                          placeholder={t('common:select-value')}
                          className="w-full [&>div>p]:truncate [&>div]:max-w-[calc(100%-36px)]"
                        />
                        <DropdownMenuContent align="start" {...field}>
                          {flagOptions?.map((item, index) => (
                            <DropdownMenuItem
                              key={index}
                              label={item.label}
                              value={item.value}
                              onSelectOption={value => {
                                field.onChange(value);
                                onChangeFormField('flagId', value);
                              }}
                            />
                          ))}
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />
            )}
            {(isUserSegment || isDate) && (
              <Form.Field
                control={control}
                name={`${commonName}value`}
                render={({ field }) => (
                  <Form.Item
                    className={cn(
                      'flex flex-col flex-1 self-stretch py-0 min-w-[170px] order-2',
                      { 'order-3': isUserSegment }
                    )}
                  >
                    <Form.Label required>
                      <Trans
                        i18nKey={'form:feature-flags.value'}
                        values={{
                          type: 'First'
                        }}
                      />
                    </Form.Label>
                    <Form.Control>
                      {isUserSegment ? (
                        <DropdownMenu>
                          <DropdownMenuTrigger
                            disabled={segmentLoading}
                            label={
                              userSegments?.find(
                                item =>
                                  item.id === (field.value || condition?.value)
                              )?.name || ''
                            }
                            placeholder={t('common:select-value')}
                            className="w-full"
                          />
                          <DropdownMenuContent align="start" {...field}>
                            {userSegments?.map((item, index) => (
                              <DropdownMenuItem
                                key={index}
                                label={item.name}
                                value={item.id}
                                onSelectOption={value => {
                                  field.onChange(value);
                                  onChangeFormField('value', value);
                                }}
                              />
                            ))}
                          </DropdownMenuContent>
                        </DropdownMenu>
                      ) : (
                        <Input
                          {...field}
                          value={field.value || condition.value}
                          onChange={value => {
                            field.onChange(value);
                            onChangeFormField('value', value);
                          }}
                        />
                      )}
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />
            )}
            <Form.Field
              control={control}
              name={`${commonName}conditioner`}
              render={({ field }) => (
                <Form.Item
                  className={cn(
                    'flex flex-col flex-1 py-0 self-stretch min-w-[170px] order-3',
                    { 'order-2': isUserSegment }
                  )}
                >
                  <Form.Label required>
                    {t('feature-flags.conditioner')}
                  </Form.Label>
                  <Form.Control>
                    {isDate ? (
                      <DropdownMenu>
                        <DropdownMenuTrigger
                          label={
                            conditionerDateOptions.find(item =>
                              [field.value, condition.conditioner].includes(
                                item.value
                              )
                            )?.label
                          }
                          className="w-full"
                        />
                        <DropdownMenuContent align="start" {...field}>
                          {conditionerDateOptions.map((item, index) => (
                            <DropdownMenuItem
                              key={index}
                              label={item.label}
                              value={item.value}
                              onSelectOption={value => {
                                field.onChange(value);
                                onChangeFormField('conditioner', value);
                              }}
                            />
                          ))}
                        </DropdownMenuContent>
                      </DropdownMenu>
                    ) : (
                      <Input
                        disabled={isUserSegment && isFlag}
                        {...field}
                        value={field.value || condition.conditioner}
                        onChange={value => {
                          field.onChange(value);
                          onChangeFormField('conditioner', value);
                        }}
                      />
                    )}
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            {isCompare && (
              <Form.Field
                control={control}
                name={`${commonName}secondValue`}
                render={({ field }) => (
                  <Form.Item className="flex flex-col flex-1 self-stretch py-0 min-w-[170px] order-4">
                    <Form.Label required>
                      <Trans
                        i18nKey={'form:feature-flags.value-type'}
                        values={{
                          type: 'Second'
                        }}
                      />
                    </Form.Label>
                    <Form.Control>
                      <Input
                        {...field}
                        value={field.value || condition.secondValue}
                        onChange={value => {
                          field.onChange(value);
                          onChangeFormField('secondValue', value);
                        }}
                      />
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />
            )}
            {isDate && (
              <Form.Field
                control={control}
                name={`${commonName}date`}
                render={({ field }) => {
                  const { value, ...rest } = field;
                  const fieldValue = Number(value ?? condition?.date) * 1000;

                  return (
                    <Form.Item className="flex flex-col flex-1 self-stretch py-0 min-w-[170px] order-4">
                      <Form.Label required>
                        {t('feature-flags.date')}
                      </Form.Label>
                      <Form.Control>
                        <ReactDatePicker
                          {...rest}
                          selected={fieldValue ? new Date(fieldValue) : null}
                          onChange={date => {
                            if (date) {
                              const value =
                                (date.getTime() / 1000)?.toString() || '';
                              field.onChange(value);
                              onChangeFormField('date', value);
                            }
                          }}
                        />
                      </Form.Control>
                      <Form.Message />
                    </Form.Item>
                  );
                }}
              />
            )}
            {isFlag && (
              <Form.Field
                control={control}
                name={`${commonName}variation`}
                render={({ field }) => (
                  <Form.Item className="flex flex-col flex-1 self-stretch py-0 min-w-[170px] order-4">
                    <Form.Label required>
                      {t('table:feature-flags.variation')}
                    </Form.Label>
                    <Form.Control>
                      <DropdownMenu>
                        <DropdownMenuTrigger
                          label={
                            variationOptions?.find(
                              item =>
                                item.value ===
                                (field.value || condition?.variation)
                            )?.label || ''
                          }
                          placeholder={t('common:select-value')}
                          className="w-full [&>div>p]:truncate [&>div]:max-w-[calc(100%-36px)]"
                        />
                        <DropdownMenuContent align="start" {...field}>
                          {variationOptions?.map((item, index) => (
                            <DropdownMenuItem
                              key={index}
                              label={item.label}
                              value={item.value}
                              onSelectOption={value => {
                                field.onChange(value);
                                onChangeFormField('variation', value);
                              }}
                            />
                          ))}
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />
            )}
            <div className="flex items-center self-stretch order-5">
              <Button
                type="button"
                disabled={isDisabledDelete}
                variant={'grey'}
                className="flex-center text-gray-500 hover:text-gray-600"
                onClick={onDeleteCondition}
              >
                <Icon icon={IconTrash} size={'sm'} />
              </Button>
            </div>
          </div>
        </div>
      </div>
    );
  }
);

export default ConditionForm;
