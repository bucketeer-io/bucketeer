import { useCallback, useMemo } from 'react';
import { useFieldArray, useFormContext, FieldPath } from 'react-hook-form';
import { useTranslation } from 'i18n';
import { omit } from 'lodash';
import { v4 as uuid } from 'uuid';
import { Feature, FeatureRuleClauseOperator, UserSegment } from '@types';
import { truncateBySide } from 'utils/converts';
import { cn } from 'utils/style';
import { IconInfo, IconPlus, IconTrash } from '@icons';
import Button from 'components/button';
import { CreatableSelect } from 'components/creatable-select';
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
import { Tooltip } from 'components/tooltip';
import {
  conditionerCompareOptions,
  conditionerDateOptions,
  situationOptions
} from '../constants';
import { TargetingSchema } from '../form-schema';
import { RuleClauseType } from '../types';

interface Props {
  features: Feature[];
  segmentIndex: number;
  userSegments?: UserSegment[];
}

const RuleForm = ({ features, segmentIndex, userSegments }: Props) => {
  const { t } = useTranslation(['form', 'common', 'table']);

  const methods = useFormContext<TargetingSchema>();
  const { control, watch, setValue } = methods;

  const clausesWatch = watch(`segmentRules.${segmentIndex}.clauses`);

  const {
    fields: clauses,
    append,
    remove
  } = useFieldArray({
    control,
    name: `segmentRules.${segmentIndex}.clauses`,
    keyName: 'clauseId'
  });

  const formatClauses = clausesWatch.map(item => ({
    ...item,
    clauseId: clauses.find(clause => clause.id === item.id)?.clauseId
  }));

  const flagOptions = useMemo(
    () =>
      features?.map(item => ({
        label: item.name,
        value: item.id
      })),
    [features]
  );

  const segmentOptions = userSegments?.map(item => ({
    label: item.name,
    value: item.id
  }));

  const getFieldName = (name: string, index: number) =>
    `segmentRules.${segmentIndex}.clauses.${index}.${name}` as FieldPath<TargetingSchema>;

  const handleChangeConditioner = useCallback(
    (value: RuleClauseType, index: number) => {
      let _value = '';
      switch (value) {
        case RuleClauseType.COMPARE:
          _value = FeatureRuleClauseOperator.EQUALS;
          break;
        case RuleClauseType.SEGMENT:
          _value = FeatureRuleClauseOperator.SEGMENT;
          break;
        case RuleClauseType.FEATURE_FLAG:
          _value = FeatureRuleClauseOperator.FEATURE_FLAG;
          break;
        case RuleClauseType.DATE:
          _value = FeatureRuleClauseOperator.BEFORE;
          break;
        default:
          break;
      }
      setValue(getFieldName('operator', index), _value);
    },
    [clauses]
  );

  return (
    <>
      <div className="flex flex-col w-full gap-y-4">
        {formatClauses.map((clause, clauseIndex) => {
          const type = clauseIndex === 0 ? 'if' : 'and';
          const isCompare = clause.type === RuleClauseType.COMPARE;
          const isUserSegment = clause.type === RuleClauseType.SEGMENT;
          const isDate = clause.type === RuleClauseType.DATE;
          const isFlag = clause.type === RuleClauseType.FEATURE_FLAG;
          const featureId = isFlag ? clause?.attribute : '';
          const variationOptions = features
            ?.find(item => item.id === featureId)
            ?.variations?.map(v => ({
              label: v.name || v.value,
              value: v.id
            }));

          return (
            <div
              key={clause.clauseId}
              className="flex items-center w-full gap-x-4"
            >
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
                <div className="flex items-end w-full gap-x-4 max-w-full">
                  <Form.Field
                    control={control}
                    name={getFieldName('type', clauseIndex)}
                    render={({ field }) => {
                      return (
                        <Form.Item className="flex flex-col flex-1 self-stretch py-0 min-w-[170px] order-1">
                          <Form.Label required>
                            {t('feature-flags.context-kind')}
                          </Form.Label>
                          <Form.Control>
                            <DropdownMenu>
                              <DropdownMenuTrigger
                                label={
                                  situationOptions.find(
                                    item => item.value === field.value
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
                                      handleChangeConditioner(
                                        value as RuleClauseType,
                                        clauseIndex
                                      );
                                    }}
                                  />
                                ))}
                              </DropdownMenuContent>
                            </DropdownMenu>
                          </Form.Control>
                          <Form.Message />
                        </Form.Item>
                      );
                    }}
                  />
                  {!isUserSegment && (
                    <Form.Field
                      control={control}
                      name={getFieldName('attribute', clauseIndex)}
                      render={({ field }) => (
                        <Form.Item className="flex flex-col flex-1 self-stretch py-0 min-w-[170px] order-2">
                          <Form.Label required className="relative w-fit">
                            {isFlag
                              ? t(`feature-flags.feature-flag`)
                              : isCompare
                                ? t(`feature-flags.attribute`)
                                : t(`feature-flags.values`)}

                            {!isFlag && (
                              <Tooltip
                                trigger={
                                  <div className="flex-center size-fit absolute top-0.5 -right-5">
                                    <Icon icon={IconInfo} size="xxs" />
                                  </div>
                                }
                              />
                            )}
                          </Form.Label>
                          <Form.Control>
                            {isFlag ? (
                              <DropdownMenu>
                                <DropdownMenuTrigger
                                  label={truncateBySide(
                                    flagOptions?.find(
                                      item =>
                                        item.value ===
                                        (field.value || clause?.attribute)
                                    )?.label || '',
                                    50
                                  )}
                                  placeholder={t('common:select-value')}
                                  className="w-full"
                                />
                                <DropdownMenuContent
                                  align="start"
                                  {...field}
                                  className="max-w-[500px]"
                                >
                                  {flagOptions?.map((item, index) => (
                                    <DropdownMenuItem
                                      key={index}
                                      label={item.label}
                                      value={item.value}
                                      onSelectOption={value =>
                                        field.onChange(value)
                                      }
                                    />
                                  ))}
                                </DropdownMenuContent>
                              </DropdownMenu>
                            ) : (
                              <Input {...field} />
                            )}
                          </Form.Control>
                          <Form.Message />
                        </Form.Item>
                      )}
                    />
                  )}
                  <Form.Field
                    control={control}
                    name={getFieldName('operator', clauseIndex)}
                    render={({ field }) => (
                      <Form.Item className="flex flex-col flex-1 self-stretch py-0 min-w-[170px] order-2">
                        <Form.Label required>
                          {t('feature-flags.operator')}
                        </Form.Label>
                        <Form.Control>
                          {isDate || isCompare ? (
                            <DropdownMenu>
                              <DropdownMenuTrigger
                                label={
                                  (isDate
                                    ? conditionerDateOptions
                                    : conditionerCompareOptions
                                  ).find(item =>
                                    [field.value, clause.operator].includes(
                                      item.value
                                    )
                                  )?.label
                                }
                                className="w-full"
                              />
                              <DropdownMenuContent align="start" {...field}>
                                {(isDate
                                  ? conditionerDateOptions
                                  : conditionerCompareOptions
                                ).map((item, index) => (
                                  <DropdownMenuItem
                                    key={index}
                                    label={item.label}
                                    value={item.value}
                                    onSelectOption={value =>
                                      field.onChange(value)
                                    }
                                  />
                                ))}
                              </DropdownMenuContent>
                            </DropdownMenu>
                          ) : (
                            <Input
                              {...field}
                              disabled={isUserSegment || isFlag}
                              value={isUserSegment ? t('is-included-in') : '='}
                            />
                          )}
                        </Form.Control>
                        <Form.Message />
                      </Form.Item>
                    )}
                  />
                  <Form.Field
                    control={control}
                    name={getFieldName('values', clauseIndex)}
                    render={({ field }) => {
                      const { value, ...rest } = field;
                      const fieldValue = isDate
                        ? value[0]
                          ? Number(value[0]) * 1000
                          : null
                        : value;
                      return (
                        <Form.Item className="flex flex-col flex-1 self-stretch py-0 min-w-[170px] order-2">
                          <Form.Label required className="relative w-fit">
                            {isFlag
                              ? t('table:feature-flags.variation')
                              : isDate
                                ? t('feature-flags.date')
                                : t('feature-flags.values')}
                            {!isFlag && !isDate && (
                              <Tooltip
                                trigger={
                                  <div className="flex-center size-fit absolute top-0.5 -right-5">
                                    <Icon icon={IconInfo} size="xxs" />
                                  </div>
                                }
                              />
                            )}
                          </Form.Label>
                          <Form.Control>
                            {isDate ? (
                              <ReactDatePicker
                                {...omit(rest, 'ref')}
                                showTimeSelect={false}
                                dateFormat={'yyyy/MM/dd'}
                                selected={
                                  fieldValue ? new Date(fieldValue) : null
                                }
                                onChange={date => {
                                  if (date) {
                                    const value =
                                      (date.getTime() / 1000)?.toString() || '';
                                    field.onChange([value]);
                                  }
                                }}
                              />
                            ) : isFlag || isUserSegment ? (
                              <DropdownMenu>
                                <DropdownMenuTrigger
                                  disabled={
                                    isFlag
                                      ? !variationOptions?.length
                                      : !segmentOptions?.length
                                  }
                                  label={
                                    (isFlag
                                      ? variationOptions
                                      : segmentOptions
                                    )?.find(item => item.value === value[0])
                                      ?.label || ''
                                  }
                                  placeholder={t('common:select-value')}
                                  className="w-full [&>div>p]:truncate [&>div]:max-w-[calc(100%-36px)]"
                                />
                                <DropdownMenuContent align="start" {...field}>
                                  {(isFlag
                                    ? variationOptions
                                    : segmentOptions
                                  )?.map((item, index) => (
                                    <DropdownMenuItem
                                      key={index}
                                      label={item.label}
                                      value={item.value}
                                      onSelectOption={value =>
                                        field.onChange([value])
                                      }
                                    />
                                  ))}
                                </DropdownMenuContent>
                              </DropdownMenu>
                            ) : (
                              <CreatableSelect
                                value={value?.map((item: string) => ({
                                  label: item,
                                  value: item
                                }))}
                                onChange={options => {
                                  const values = options.map(
                                    item => item.value
                                  );
                                  field.onChange(values);
                                }}
                              />
                            )}
                          </Form.Control>
                          <Form.Message />
                        </Form.Item>
                      );
                    }}
                  />

                  <div className="flex items-center self-stretch order-5">
                    <Button
                      type="button"
                      disabled={formatClauses.length <= 1}
                      variant={'grey'}
                      className="flex-center text-gray-500 hover:text-gray-600"
                      onClick={() => remove(clauseIndex)}
                    >
                      <Icon icon={IconTrash} size={'sm'} />
                    </Button>
                  </div>
                </div>
              </div>
            </div>
          );
        })}
      </div>
      <Button
        type="button"
        variant={'text'}
        className="w-fit gap-x-2 h-6 !p-0"
        onClick={() =>
          append({
            id: uuid(),
            type: RuleClauseType.COMPARE,
            attribute: '',
            operator: FeatureRuleClauseOperator.EQUALS,
            values: []
          })
        }
      >
        <Icon
          icon={IconPlus}
          color="primary-500"
          className="flex-center"
          size={'sm'}
        />{' '}
        {t('form:feature-flags.add-condition')}
      </Button>
    </>
  );
};

export default RuleForm;
