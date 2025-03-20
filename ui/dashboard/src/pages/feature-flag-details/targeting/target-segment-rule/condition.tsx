import { forwardRef, Ref, useCallback, useMemo } from 'react';
import { useFormContext } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { useQueryUserSegments } from '@queries/user-segments';
import { getCurrentEnvironment, useAuth } from 'auth';
import { LIST_PAGE_SIZE } from 'constants/app';
import { useTranslation } from 'i18n';
import { omit } from 'lodash';
import { Feature, FeatureRuleClauseOperator } from '@types';
import { cn } from 'utils/style';
import { IconTrash } from '@icons';
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
import { RuleClauseSchema, RuleClauseType } from '../form-schema';

interface Props {
  features: Feature[];
  type: 'if' | 'and';
  clause: RuleClauseSchema;
  isDisabledDelete: boolean;
  segmentIndex: number;
  clauseIndex: number;
  onDeleteCondition: () => void;
  onChangeFormField: (
    field: keyof RuleClauseSchema,
    value: string | string[]
  ) => void;
}

const ConditionForm = forwardRef(
  (
    {
      features,
      type,
      clause,
      isDisabledDelete,
      segmentIndex,
      clauseIndex,
      onDeleteCondition,
      onChangeFormField
    }: Props,
    ref: Ref<HTMLDivElement>
  ) => {
    const { t } = useTranslation(['form', 'common', 'table']);
    const { consoleAccount } = useAuth();
    const currentEnvironment = getCurrentEnvironment(consoleAccount!);

    const isCompare = useMemo(
      () => clause.type === RuleClauseType.COMPARE,
      [clause]
    );
    const isUserSegment = useMemo(
      () => clause.type === RuleClauseType.SEGMENT,
      [clause]
    );
    const isDate = useMemo(() => clause.type === RuleClauseType.DATE, [clause]);
    const isFlag = useMemo(
      () => clause.type === RuleClauseType.FEATURE_FLAG,
      [clause]
    );

    const methods = useFormContext();
    const { control, watch, setValue } = methods;

    const situationOptions = useMemo(
      () => [
        {
          label: t('feature-flags.compare'),
          value: RuleClauseType.COMPARE
        },
        {
          label: t('feature-flags.user-segment'),
          value: RuleClauseType.SEGMENT
        },
        {
          label: t('feature-flags.date'),
          value: RuleClauseType.DATE
        },
        {
          label: t('feature-flags.feature-flag'),
          value: RuleClauseType.FEATURE_FLAG
        }
      ],
      []
    );

    const conditionerCompareOptions = useMemo(
      () => [
        {
          label: '=',
          value: FeatureRuleClauseOperator.EQUALS
        },
        {
          label: '>=',
          value: FeatureRuleClauseOperator.GREATER_OR_EQUAL
        },
        {
          label: '>',
          value: FeatureRuleClauseOperator.GREATER
        },
        {
          label: '<=',
          value: FeatureRuleClauseOperator.LESS_OR_EQUAL
        },
        {
          label: '<',
          value: FeatureRuleClauseOperator.LESS
        },
        {
          label: t('contains'),
          value: FeatureRuleClauseOperator.IN
        },
        {
          label: t('partially-matches'),
          value: FeatureRuleClauseOperator.PARTIALLY_MATCH
        },
        {
          label: t('starts-with'),
          value: FeatureRuleClauseOperator.STARTS_WITH
        },
        {
          label: t('ends-with'),
          value: FeatureRuleClauseOperator.ENDS_WITH
        }
      ],
      []
    );

    const conditionerDateOptions = useMemo(
      () => [
        {
          label: t('before'),
          value: FeatureRuleClauseOperator.BEFORE
        },
        {
          label: t('after'),
          value: FeatureRuleClauseOperator.AFTER
        }
      ],
      []
    );

    const commonName = useMemo(
      () => `rules.${segmentIndex}.clauses.${clauseIndex}.`,
      [segmentIndex, clauseIndex]
    );
    const featureId = isFlag ? watch(`${commonName}attribute`) : '';

    const flagOptions = useMemo(
      () =>
        features?.map(item => ({
          label: item.name,
          value: item.id
        })),
      []
    );

    const variationOptions = useMemo(
      () =>
        features
          ?.find(item => item.id === featureId)
          ?.variations?.map(v => ({
            label: v.name || v.value,
            value: v.id
          })),
      [featureId, features]
    );

    const { data: segmentCollection } = useQueryUserSegments({
      params: {
        cursor: String(0),
        pageSize: LIST_PAGE_SIZE,
        environmentId: currentEnvironment.id
      },
      enabled: isUserSegment
    });

    const userSegments = segmentCollection?.segments || [];

    const segmentOptions = userSegments?.map(item => ({
      label: item.name,
      value: item.id
    }));

    const handleChangeConditioner = useCallback(
      (value: RuleClauseType) => {
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
        onChangeFormField('operator', _value);
        setValue(`${commonName}operator`, _value);
      },
      [commonName]
    );

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
              name={`${commonName}type`}
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
                            item => item.value === clause.type
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
                              onChangeFormField('type', value as string);
                              handleChangeConditioner(value as RuleClauseType);
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
            {!isUserSegment && (
              <Form.Field
                control={control}
                name={`${commonName}attribute`}
                render={({ field }) => (
                  <Form.Item className="flex flex-col flex-1 self-stretch py-0 min-w-[170px] order-2">
                    <Form.Label required>
                      {isFlag ? (
                        t('feature-flags.feature-flag')
                      ) : (
                        <Trans
                          i18nKey={'form:feature-flags.value-type'}
                          values={{
                            type: isCompare ? 'First' : ''
                          }}
                        />
                      )}
                    </Form.Label>
                    <Form.Control>
                      {isFlag ? (
                        <DropdownMenu>
                          <DropdownMenuTrigger
                            label={
                              flagOptions?.find(
                                item =>
                                  item.value ===
                                  (field.value || clause?.attribute)
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
                                  onChangeFormField(
                                    'attribute',
                                    value as string
                                  );
                                }}
                              />
                            ))}
                          </DropdownMenuContent>
                        </DropdownMenu>
                      ) : (
                        <Input
                          {...field}
                          onChange={value => {
                            field.onChange(value);
                            onChangeFormField('attribute', value);
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
              name={`${commonName}operator`}
              render={({ field }) => (
                <Form.Item className="flex flex-col flex-1 self-stretch py-0 min-w-[170px] order-2">
                  <Form.Label required>
                    {t('feature-flags.conditioner')}
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
                              onSelectOption={value => {
                                field.onChange(value);
                                onChangeFormField('operator', value as string);
                              }}
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
              name={`${commonName}values`}
              render={({ field }) => {
                const { value, ...rest } = field;
                const fieldValue = isDate ? Number(value[0]) * 1000 : value;
                return (
                  <Form.Item className="flex flex-col flex-1 self-stretch py-0 min-w-[170px] order-2">
                    <Form.Label required>
                      {isFlag ? (
                        t('table:feature-flags.variation')
                      ) : isDate ? (
                        t('feature-flags.date')
                      ) : (
                        <Trans
                          i18nKey={'form:feature-flags.value-type'}
                          values={{
                            type: isCompare ? 'Second' : ''
                          }}
                        />
                      )}
                    </Form.Label>
                    <Form.Control>
                      {isDate ? (
                        <ReactDatePicker
                          {...omit(rest, 'ref')}
                          selected={fieldValue ? new Date(fieldValue) : null}
                          onChange={date => {
                            if (date) {
                              const value =
                                (date.getTime() / 1000)?.toString() || '';
                              field.onChange([value]);
                              onChangeFormField('values', [value]);
                            }
                          }}
                        />
                      ) : isFlag || isUserSegment ? (
                        <DropdownMenu>
                          <DropdownMenuTrigger
                            disabled={
                              isFlag
                                ? !variationOptions?.length
                                : !segmentOptions.length
                            }
                            label={
                              (isFlag
                                ? variationOptions
                                : segmentOptions
                              )?.find(item => item.value === value[0])?.label ||
                              ''
                            }
                            placeholder={t('common:select-value')}
                            className="w-full [&>div>p]:truncate [&>div]:max-w-[calc(100%-36px)]"
                          />
                          <DropdownMenuContent align="start" {...field}>
                            {(isFlag ? variationOptions : segmentOptions)?.map(
                              (item, index) => (
                                <DropdownMenuItem
                                  key={index}
                                  label={item.label}
                                  value={item.value}
                                  onSelectOption={value => {
                                    field.onChange([value]);
                                    onChangeFormField('values', [
                                      value as string
                                    ]);
                                  }}
                                />
                              )
                            )}
                          </DropdownMenuContent>
                        </DropdownMenu>
                      ) : (
                        <CreatableSelect
                          value={value?.map((item: string) => ({
                            label: item,
                            value: item
                          }))}
                          onChange={options => {
                            const _value = options.map(item => item.value);
                            field.onChange(_value);
                            onChangeFormField('values', _value);
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
