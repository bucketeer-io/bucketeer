import { useMemo } from 'react';
import { Control, FieldErrors } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { Link } from 'react-router';
import { PAGE_PATH_USER_SEGMENTS } from 'constants/routing';
import { useUserSegmentsLoader } from 'hooks/use-user-segments-loading-more';
import { useTranslation } from 'i18n';
import get from 'lodash/get';
import omit from 'lodash/omit';
import { Feature } from '@types';
import { truncateBySide } from 'utils/converts';
import { isNotEmptyObject } from 'utils/data-type';
import { cn } from 'utils/style';
import { IconInfo, IconTrash } from '@icons';
import { UserMessage } from 'pages/feature-flag-details/targeting/individual-rule';
import { RuleClauseType } from 'pages/feature-flag-details/targeting/types';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import Button from 'components/button';
import { CreatableSelect } from 'components/creatable-select';
import { ReactDatePicker } from 'components/date-time-picker';
import Dropdown, { DropdownOption } from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import { Popover } from 'components/popover';
import Spinner from 'components/spinner';
import { Tooltip } from 'components/tooltip';
import DropdownMenuWithSearch from 'elements/dropdown-with-search';
import FeatureFlagStatus from 'elements/feature-flag-status';
import AttributeKeySelect from './attribute-key-select';
import { getSegmentSummary, SituationOption } from './index';

interface ClauseValue {
  id?: string;
  type?: string;
  attribute?: string;
  operator?: string;
  values?: string[];
}

interface Props {
  clause: ClauseValue & { clauseId?: string };
  clauseIndex: number;
  clausesName: string;
  clausesLength: number;
  situationOptions: SituationOption[];
  environmentId: string;
  environmentUrlCode: string;
  isLanguageJapanese: boolean;
  control: Control;
  errors: FieldErrors;
  feature?: Feature;
  features?: Feature[];
  flagOptions: DropdownOption[];
  createdOptionList: string[];
  sdkAttributeKeys: string[];
  conditionerDateOptions: DropdownOption[];
  conditionerCompareOptions: DropdownOption[];
  onCreateOption: (value: string) => void;
  onChangeConditioner: (
    value: RuleClauseType,
    index: number,
    onChange: (value: RuleClauseType) => void
  ) => void;
  onRemoveClause: (index: number) => void;
  getFieldName: (name: string, index: number) => string;
}

const ClauseRow = ({
  clause,
  clauseIndex,
  clausesName,
  clausesLength,
  situationOptions,
  environmentId,
  environmentUrlCode,
  isLanguageJapanese,
  control,
  errors,
  features,
  flagOptions,
  createdOptionList,
  sdkAttributeKeys,
  conditionerDateOptions,
  conditionerCompareOptions,
  onCreateOption,
  onChangeConditioner,
  onRemoveClause,
  getFieldName
}: Props) => {
  const { t } = useTranslation(['form', 'common', 'table']);

  const type = clauseIndex === 0 ? 'if' : 'and';
  const isCompare = clause.type === RuleClauseType.COMPARE;
  const isUserSegment = clause.type === RuleClauseType.SEGMENT;
  const isDate = clause.type === RuleClauseType.DATE;
  const isFlag = clause.type === RuleClauseType.FEATURE_FLAG;
  const featureId = isFlag ? clause?.attribute : '';
  const variationOptions = useMemo(
    () =>
      features
        ?.find(item => item.id === featureId)
        ?.variations?.map((v, index) => ({
          label: (
            <div className="flex items-center gap-x-2 pl-0.5">
              <FlagVariationPolygon index={index} />
              <p className="-mt-0.5 truncate">{v.name}</p>
            </div>
          ),
          value: v.id
        })),
    [features, featureId]
  );

  const selectedSegmentIds = clause.values || [];

  const {
    userSegments: loadedUserSegments,
    selectedSegments: selectedUserSegments,
    isLoadingMore: isLoadingMoreSegments,
    hasMore: hasMoreSegments,
    isInitialLoading: isLoadingSegments,
    isResolvingSelection: isResolvingSelectedSegments,
    hasNoSegments,
    loadMore: loadMoreSegments
  } = useUserSegmentsLoader({
    environmentId,
    selectedSegmentIds,
    enabled: isUserSegment
  });

  const segmentOptions = isUserSegment
    ? loadedUserSegments?.map(item => ({
        label: `${item.name} (${getSegmentSummary(item, t)})`,
        value: item.id
      }))
    : [];

  const isEmptySegment = isUserSegment && hasNoSegments;
  const isHaveError = isNotEmptyObject(
    (get(errors, `${clausesName}.${clauseIndex}`) as unknown as object) || {}
  );

  return (
    <div className="flex items-center w-full gap-x-4">
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
      <div className="flex items-center w-full flex-1 pl-4 border-l border-primary-500 gap-x-4">
        <div
          className={cn(
            'grid grid-cols-4 items-end w-full gap-x-4 max-w-full',
            {
              'grid-cols-3': isUserSegment && !isEmptySegment
            }
          )}
        >
          <div
            className={cn('flex flex-1 col-span-1 self-stretch', {
              'flex-initial': isEmptySegment
            })}
          >
            <Form.Field
              control={control}
              name={getFieldName('type', clauseIndex)}
              render={({ field }) => {
                return (
                  <Form.Item
                    className={cn(
                      'flex flex-col w-full self-stretch py-0 min-w-[170px] order-1',
                      {
                        'max-w-[250px]': isEmptySegment
                      }
                    )}
                  >
                    <Form.Label required>
                      {t('feature-flags.context-kind')}
                    </Form.Label>
                    <Form.Control>
                      <Dropdown
                        options={situationOptions}
                        value={field.value}
                        className="w-full"
                        onChange={value => {
                          onChangeConditioner(
                            value as RuleClauseType,
                            clauseIndex,
                            field.onChange
                          );
                        }}
                      />
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                );
              }}
            />
          </div>

          {!isUserSegment && (
            <div className="flex flex-1 col-span-1 self-stretch">
              <Form.Field
                control={control}
                name={getFieldName('attribute', clauseIndex)}
                render={({ field }) => {
                  return (
                    <Form.Item className="flex flex-col w-full self-stretch py-0 min-w-[170px] order-2">
                      <Form.Label required className="relative w-fit">
                        {isFlag
                          ? t(`feature-flags.feature-flag`)
                          : t(`feature-flags.attribute-key`)}

                        {!isFlag && (
                          <Tooltip
                            content={t('targeting.tooltip.attribute')}
                            trigger={
                              <div className="flex-center size-fit absolute top-0.5 -right-5">
                                <Icon icon={IconInfo} size="xxs" />
                              </div>
                            }
                            className="max-w-[300px]"
                          />
                        )}
                      </Form.Label>
                      <Form.Control>
                        {isFlag ? (
                          <DropdownMenuWithSearch
                            align="start"
                            label={truncateBySide(
                              features?.find(item =>
                                [field.value, clause?.attribute].includes(
                                  item.id
                                )
                              )?.name || '',
                              50
                            )}
                            placeholder={t('experiments.select-flag')}
                            isExpand
                            options={flagOptions}
                            selectedOptions={field.value}
                            additionalElement={item => (
                              <FeatureFlagStatus
                                status={t(
                                  item.enabled
                                    ? 'experiments.on'
                                    : 'experiments.off'
                                )}
                                enabled={item.enabled as boolean}
                              />
                            )}
                            onSelectOption={value => {
                              field.onChange(value);
                            }}
                            contentClassName="!w-[500px] !max-w-[500px]"
                          />
                        ) : (
                          <AttributeKeySelect
                            createdOptions={createdOptionList?.map(
                              (item: string) => ({
                                label: item,
                                value: item
                              })
                            )}
                            sdkOptions={sdkAttributeKeys
                              ?.sort()
                              .map((item: string) => ({
                                label: item,
                                value: item
                              }))}
                            onChange={value => field.onChange(value)}
                            onCreateOption={(value: string) => {
                              onCreateOption(value);
                              field.onChange(value);
                            }}
                            value={{
                              label: field.value,
                              value: field.value
                            }}
                          />
                        )}
                      </Form.Control>
                      <Form.Message />
                    </Form.Item>
                  );
                }}
              />
            </div>
          )}
          <div
            className={cn('flex flex-1 col-span-1 self-stretch', {
              'col-span-3': isEmptySegment
            })}
          >
            <Form.Field
              control={control}
              name={getFieldName('operator', clauseIndex)}
              render={({ field }) => (
                <Form.Item className="flex flex-col flex-1 self-stretch py-0 min-w-[170px] order-2">
                  {!isEmptySegment && (
                    <Form.Label required>
                      {t('feature-flags.operator')}
                    </Form.Label>
                  )}
                  <Form.Control>
                    {isDate || isCompare ? (
                      <Dropdown
                        options={
                          isDate
                            ? conditionerDateOptions
                            : conditionerCompareOptions
                        }
                        value={field.value ?? clause.operator}
                        onChange={value => field.onChange(value)}
                        placeholder={t('common:select-condition')}
                        className="w-full"
                        alignContent="start"
                      />
                    ) : isEmptySegment ? (
                      <div className="flex items-end mb-4 h-full typo-para-small text-gray-700">
                        <Trans
                          i18nKey={'message:empty-segment'}
                          components={{
                            comp: (
                              <Link
                                target="_blank"
                                to={`/${environmentUrlCode}${PAGE_PATH_USER_SEGMENTS}`}
                                className={cn('text-primary-500 underline', {
                                  'mx-1': !isLanguageJapanese
                                })}
                              />
                            )
                          }}
                        />
                      </div>
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
          </div>
          {!isEmptySegment && (
            <div className="flex flex-1 col-span-1 self-stretch">
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
                  const selectedSegments = isUserSegment
                    ? selectedUserSegments || []
                    : [];
                  return (
                    <Form.Item className="flex flex-col flex-1 self-stretch py-0 min-w-[170px] order-2">
                      <Form.Label required className="relative w-fit">
                        {isFlag
                          ? t('table:feature-flags.variation')
                          : isDate
                            ? t('feature-flags.value')
                            : t('feature-flags.values')}
                        {!isFlag && !isDate && (
                          <Tooltip
                            content={t(
                              isUserSegment
                                ? 'targeting.tooltip.segment-value'
                                : 'targeting.tooltip.value'
                            )}
                            trigger={
                              <div className="flex-center size-fit absolute top-0.5 -right-5">
                                <Icon icon={IconInfo} size="xxs" />
                              </div>
                            }
                            className="max-w-[310px]"
                          />
                        )}
                      </Form.Label>
                      <Form.Control>
                        {isDate ? (
                          <ReactDatePicker
                            {...omit(rest, 'ref')}
                            timeFormat="HH:mm"
                            selected={fieldValue ? new Date(fieldValue) : null}
                            onChange={date => {
                              if (date) {
                                const value =
                                  (date.getTime() / 1000)?.toString() || '';
                                field.onChange([value]);
                              }
                            }}
                          />
                        ) : isUserSegment ? (
                          <DropdownMenuWithSearch
                            isMultiselect
                            isExpand
                            label={
                              selectedSegments.length
                                ? truncateBySide(
                                    selectedSegments
                                      .map(item => item.name)
                                      .join(', '),
                                    50
                                  )
                                : ''
                            }
                            placeholder={t('common:select-value')}
                            isLoading={isLoadingSegments}
                            isLoadingMore={isLoadingMoreSegments}
                            isHasMore={hasMoreSegments}
                            hideSearchInput
                            onHasMoreOptions={loadMoreSegments}
                            options={segmentOptions || []}
                            selectedOptions={(value as string[]) || []}
                            onSelectOption={val => {
                              const current = (value as string[]) || [];
                              const next = current.includes(val as string)
                                ? current.filter(v => v !== val)
                                : [...current, val as string];
                              field.onChange(next);
                            }}
                            onClear={() => field.onChange([])}
                            showClear
                            triggerClassName="w-full [&>div>p]:truncate [&>div]:max-w-[calc(100%-36px)]"
                          />
                        ) : isFlag ? (
                          <Dropdown
                            options={variationOptions}
                            value={value?.[0] ?? ''}
                            onChange={val => field.onChange([val])}
                            placeholder={t('common:select-value')}
                            disabled={!variationOptions?.length}
                            className="w-full [&>div>p]:truncate [&>div]:max-w-[calc(100%-36px)]"
                          />
                        ) : (
                          <CreatableSelect
                            value={value?.map((item: string) => ({
                              label: item,
                              value: item
                            }))}
                            onChange={options => {
                              const values = options.map(item => item.value);
                              field.onChange(values);
                            }}
                            formatCreateLabel={value => (
                              <p>
                                {`${t('create-option', {
                                  option: value
                                })}`}
                              </p>
                            )}
                            noOptionsMessage={() => (
                              <UserMessage
                                message={t('no-opts-type-to-create')}
                              />
                            )}
                          />
                        )}
                      </Form.Control>
                      {isUserSegment && selectedSegments.length > 0 && (
                        <div className="mt-0.5 flex items-center gap-x-1">
                          <Popover
                            align="start"
                            trigger={
                              <div>
                                <span className="typo-para-small font-medium text-primary-500">
                                  {t('common:show-count-user-segments', {
                                    count: selectedSegments.length
                                  })}
                                </span>
                              </div>
                            }
                            triggerCls="w-fit justify-start"
                            className="flex flex-col w-[300px] gap-y-0.5 p-0 overflow-hidden"
                          >
                            <div className="flex flex-col gap-y-0.5 max-h-[220px] overflow-y-auto small-scroll p-2">
                              {selectedSegments.map(item => (
                                <div
                                  key={item.id}
                                  className="flex items-center w-full gap-x-1 rounded hover:bg-primary-50"
                                >
                                  <Link
                                    target="_blank"
                                    to={`/${environmentUrlCode}${PAGE_PATH_USER_SEGMENTS}/${item.id}`}
                                    className="typo-para-small text-primary-500 hover:underline truncate flex-1 min-w-0 px-1 py-1"
                                  >
                                    {`${item.name} (${getSegmentSummary(item, t)})`}
                                  </Link>
                                </div>
                              ))}
                            </div>
                          </Popover>
                          {isResolvingSelectedSegments && <Spinner size="sm" />}
                        </div>
                      )}
                      <Form.Message />
                    </Form.Item>
                  );
                }}
              />
            </div>
          )}
        </div>
        <div
          className={cn('flex items-center mt-[22px] self-stretch order-5', {
            'items-end mb-4 mt-0': isEmptySegment,
            'mt-0': isHaveError
          })}
        >
          <Button
            type="button"
            disabled={clausesLength <= 1}
            variant={'grey'}
            className="flex-center text-gray-500 hover:text-gray-600 size-fit p-0"
            onClick={() => onRemoveClause(clauseIndex)}
          >
            <Icon icon={IconTrash} size={'sm'} />
          </Button>
        </div>
      </div>
    </div>
  );
};

export default ClauseRow;
