import { useMemo } from 'react';
import { Trans } from 'react-i18next';
import { IconRemoveOutlined } from 'react-icons-material-design';
import { useTranslation } from 'i18n';
import { isNil } from 'lodash';
import { capitalize, cn } from 'utils/style';
import { IconArrowUpDown, IconPlus, IconSwitchUpdate } from '@icons';
import {
  DiscardChangesStateData,
  DiscardChangesType,
  VariationPercent
} from 'pages/feature-flag-details/targeting/types';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import AutoWrapText from 'components/auto-wrap-text';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';

interface Props {
  isOpen: boolean;
  type: DiscardChangesType | undefined;
  data: DiscardChangesStateData[];
  ruleIndex?: number;
  actionSegmentRule?: 'new-rule' | 'edit-rule' | undefined;
  ruleDiscardChange?: DiscardChangesType;
  onClose: () => void;
  onSubmit: (type: DiscardChangesType, index?: number) => void;
}

interface RuleHeader extends Omit<DiscardChangesStateData, 'variationIndex'> {
  formNotify: string;
  variations?: VariationPercent[];
}

const notifyMap: Record<string, (action?: string) => string> = {
  clause: action =>
    `form:${action?.toLocaleLowerCase()}-custom-rule-discard-desc`,
  strategy: () => 'form:custom-rule-strategy-discard-desc',
  audience: action =>
    `form:${action?.toLocaleLowerCase()}-custom-rule-audience-discard-desc`,
  'new-rule': () => 'form:custom-rule-strategy-add-new-discard-desc',
  'default-strategy': () => 'form:custom-rule-strategy-discard-desc',
  'default-audience': () => 'form:custom-rule-default-audience-discard-desc'
};

const ActionIcon = ({
  labelType
}: {
  labelType: DiscardChangesStateData['labelType'];
}) => (
  <Icon
    icon={
      labelType === 'ADD'
        ? IconPlus
        : labelType === 'REMOVE'
          ? IconRemoveOutlined
          : labelType === 'REORDER'
            ? IconArrowUpDown
            : IconSwitchUpdate
    }
    size={'sm'}
    color={
      labelType === 'UPDATE'
        ? 'gray-600'
        : labelType === 'REMOVE'
          ? 'accent-red-500'
          : 'accent-green-500'
    }
  />
);

export const PrerequisiteDiscardItem = ({
  labelType,
  label,
  variationIndex,
  variation
}: DiscardChangesStateData) => {
  return (
    <div className="flex flex-col w-full pl-4 gap-1">
      <div className="flex w-full gap-x-2">
        <div className="flex items-start mt-1">
          <ActionIcon labelType={labelType} />
        </div>
        <div className="typo-para-medium text-gray-700">
          <Trans
            i18nKey={`form:${labelType.toLowerCase()}_prerequisite-discard-desc`}
            values={{
              flagName: label
            }}
            components={{
              b: <strong />,
              variantElement: (
                <div className="inline-flex items-center gap-x-1">
                  <div className="flex-center size-fit">
                    <FlagVariationPolygon index={variationIndex} />
                  </div>
                  <p>{variation?.name || variation?.value}</p>
                </div>
              )
            }}
          />
        </div>
      </div>
    </div>
  );
};

export const IndividualDiscardItem = ({
  labelType,
  groupLabel,
  variationIndex,
  variation
}: DiscardChangesStateData) => {
  const { t } = useTranslation(['common', 'form']);
  const formatLabel = groupLabel?.map((item, index) => (
    <div className="inline-flex flex-wrap items-center" key={index}>
      <AutoWrapText
        text={item}
        width={450}
        sparate={index !== groupLabel.length - 1}
      />
    </div>
  ));
  return (
    <div className="flex flex-col w-full pl-4 gap-1">
      <div className="flex w-full gap-x-2">
        <div className="flex items-start mt-1">
          <ActionIcon labelType={labelType} />
        </div>
        <div className="inline-flex flex-wrap gap-1 typo-para-medium text-gray-700 ">
          <Trans
            i18nKey={`form:${labelType.toLowerCase()}_individual-discard-desc`}
            values={{
              conjection: labelType === 'REMOVE' ? t('form:from') : t('form:to')
            }}
            components={{
              labelElement: (
                <div className="inline-flex flex-wrap items-center gap-1">
                  {formatLabel}
                </div>
              ),
              variantElement: (
                <div className="inline-flex items-center gap-x-1">
                  <div className="flex-center size-fit">
                    <FlagVariationPolygon index={variationIndex} />
                  </div>
                  <p>{variation?.name || variation?.value}</p>
                </div>
              )
            }}
          />
        </div>
      </div>
    </div>
  );
};

const RuleHeader = ({
  isAddNew,
  labelType,
  label,
  changeType,
  formNotify,
  valueLabel,
  variations
}: RuleHeader) => {
  const { t } = useTranslation(['common', 'form']);
  if (isAddNew) return null;
  const isNewRule = changeType === 'new-rule';
  const isDeleteRule = changeType === 'deleted-rule';

  return (
    <div className="flex w-full gap-x-2">
      <div className="mt-[3px]">
        <ActionIcon labelType={labelType} />
      </div>
      <div className="typo-para-medium text-gray-700">
        {isNewRule || isDeleteRule ? (
          <div className="inline items-center">
            <div className="inline">
              <Trans
                i18nKey={isNewRule ? 'form:add-new-rule' : 'form:delete-rule'}
                values={{
                  clauseLabel: label
                }}
                components={{
                  b: <strong />,
                  variantElement: <StrategyList variations={variations} />
                }}
              />
            </div>
          </div>
        ) : (
          <div className="inline items-center">
            <div className="inline-block mr-1">
              <Trans
                i18nKey={formNotify}
                values={{
                  action: t(capitalize(labelType?.toLowerCase())),
                  value: changeType === 'value' ? valueLabel : ''
                }}
                components={{ b: <strong /> }}
              />
            </div>
            <strong>{label}</strong>
            <StrategyList variations={variations} />
          </div>
        )}
      </div>
    </div>
  );
};

const AudienceChange = ({
  audienceExcluded,
  audienceIncluded
}: {
  audienceExcluded: VariationPercent;
  audienceIncluded?: VariationPercent[];
}) => {
  const { t } = useTranslation(['common', 'form']);
  return (
    <div className="pl-6 text-gray-700">
      <div className="flex w-full gap-x-2 items-center pl-3">
        <div className="inline">
          <div className="inline">
            <Trans
              i18nKey={t('form:custom-rule-audience-not-include-desc')}
              values={{
                percent: audienceExcluded.weight,
                variation: audienceExcluded.variation
              }}
              components={{
                b: <strong />,
                variantElement: (
                  <div className="inline-flex items-center gap-x-1 ml-1">
                    <div className="inline-flex flex-center size-fit">
                      {!isNil(audienceExcluded.variationIndex) && (
                        <FlagVariationPolygon
                          index={audienceExcluded.variationIndex || 0}
                        />
                      )}
                    </div>
                    <p className="inline max-w-[300px] truncate">
                      {audienceExcluded.variation}
                    </p>
                  </div>
                )
              }}
            />
          </div>
        </div>
      </div>
      <div className="flex w-full gap-x-2 items-start">
        <div className="inline pl-3">
          <div className="inline">
            <Trans
              i18nKey={t('form:custom-rule-audience-include-desc')}
              values={{
                percent: Number(100 - (audienceExcluded.weight || 0))
              }}
              components={{
                b: <strong />,
                variantElement: (
                  <div className="pl-3">
                    {audienceIncluded?.map((audience, index) => (
                      <div
                        key={index}
                        className="flex items-center gap-x-1 ml-1"
                      >
                        <div className="inline-flex flex-center size-fit">
                          {!isNil(audience.variationIndex) && (
                            <FlagVariationPolygon
                              index={audience.variationIndex || 0}
                            />
                          )}
                        </div>
                        <div className="inline-flex">
                          <p className="max-w-[450px] truncate">
                            {audience.variation}
                          </p>
                        </div>
                      </div>
                    ))}
                  </div>
                )
              }}
            />
          </div>
        </div>
      </div>
    </div>
  );
};

const StrategyList = ({ variations }: { variations?: VariationPercent[] }) =>
  !!variations?.length && (
    <div className="pl-3">
      {variations.map((v, index) => (
        <div key={index} className="flex items-center gap-x-1 ml-1">
          <div className="flex-center size-fit">
            <FlagVariationPolygon index={v.variationIndex || 0} />
          </div>
          <p className="max-w-[320px] h-full truncate">{v.variation}</p>
          {!isNil(v.weight) && (
            <p className="max-w-[180px] truncate text-gray-700">
              <span>({v.weight?.toString()}%)</span>
            </p>
          )}
        </div>
      ))}
    </div>
  );

export const CustomRuleDiscardItem = ({
  labelType,
  label,
  changeType,
  valueLabel,
  variationPercent,
  audienceExcluded,
  audienceIncluded,
  isAddNew
}: DiscardChangesStateData) => {
  const { t } = useTranslation(['common', 'form']);

  const formNotify = t(
    notifyMap[changeType || 'clause']?.(
      changeType === 'clause' || changeType === 'audience'
        ? labelType
        : undefined
    ) || `form:${labelType.toLowerCase()}-custom-rule-clause-value-discard-desc`
  );

  const showAudience =
    ['default-audience', 'audience'].includes(changeType || '') &&
    audienceExcluded;
  const showVariationPercent =
    ['strategy', 'default-strategy', 'new-rule', 'deleted-rule'].includes(
      changeType || ''
    ) && !!variationPercent?.length;
  return (
    <div className={cn('flex flex-col w-full pl-4')}>
      <RuleHeader
        isAddNew={isAddNew}
        labelType={labelType}
        label={label}
        changeType={changeType}
        formNotify={formNotify}
        valueLabel={valueLabel}
        variations={showVariationPercent ? variationPercent : []}
      />

      {showAudience && (
        <AudienceChange
          audienceExcluded={audienceExcluded}
          audienceIncluded={audienceIncluded}
        />
      )}
    </div>
  );
};

const DiscardChangeModal = ({
  isOpen,
  type,
  data,
  ruleIndex,
  actionSegmentRule,
  ruleDiscardChange,
  onClose,
  onSubmit
}: Props) => {
  const { t } = useTranslation(['common', 'form']);
  const isEdit = useMemo(
    () => data.some(item => item?.changeType !== 'new-rule'),
    [data]
  );
  const ruleLabel = useMemo(() => {
    if (ruleDiscardChange === DiscardChangesType.INDIVIDUAL)
      return 'common:custom-individual-rule';
    if (ruleDiscardChange === DiscardChangesType.PREREQUISITE) {
      if (!isNil(ruleIndex) && actionSegmentRule === 'edit-rule' && isEdit)
        return {
          key: 'common:custom-segment-rule',
          values: { rule: ruleIndex! + 1 }
        };
      return 'common:custom-prerequises-rule';
    }
    if (ruleDiscardChange === DiscardChangesType.DEFAULT)
      return 'common:custom-default-rule';
    return null;
  }, [ruleDiscardChange, ruleIndex, actionSegmentRule, isEdit]);

  return (
    <DialogModal
      className="w-full max-w-[620px]"
      title={t('form:discard-unsaved-changes')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="flex flex-col w-full gap-y-4 p-5 max-h-[500px] overflow-auto small-scroll">
        <>
          {!isNil(ruleIndex) &&
            actionSegmentRule === 'edit-rule' &&
            ruleDiscardChange === DiscardChangesType.CUSTOM &&
            isEdit && (
              <div className="flex gap-1 items-center typo-para-medium leading-4 text-gray-700 font-bold">
                <Trans
                  i18nKey={'common:custom-segment-rule'}
                  values={{ rule: ruleIndex! + 1 }}
                />
              </div>
            )}
          {ruleLabel && (
            <div className="flex gap-1 items-center typo-para-medium leading-4 text-gray-700 font-bold">
              {typeof ruleLabel === 'string' ? (
                <Trans i18nKey={ruleLabel} />
              ) : (
                <Trans i18nKey={ruleLabel.key} values={ruleLabel.values} />
              )}
            </div>
          )}

          {data.map((item, index) => {
            const { PREREQUISITE, INDIVIDUAL, CUSTOM, DEFAULT } =
              DiscardChangesType;
            if (isNil(item)) return null;
            if (type === PREREQUISITE)
              return <PrerequisiteDiscardItem key={index} {...item} />;
            if (type === INDIVIDUAL)
              return <IndividualDiscardItem key={index} {...item} />;
            if (type === CUSTOM || type === DEFAULT) {
              return (
                <CustomRuleDiscardItem
                  key={index}
                  {...item}
                  ruleIndex={Number(ruleIndex)}
                />
              );
            }
            return null;
          })}
        </>
      </div>

      <ButtonBar
        primaryButton={
          <Button
            type="button"
            variant="secondary"
            className="p-2 h-9 font-bold text-sm rounded-md"
            onClick={onClose}
          >
            {t(`common:cancel`)}
          </Button>
        }
        secondaryButton={
          <Button
            type="button"
            variant="negative"
            className="p-2 h-9 font-bold text-sm rounded-md"
            onClick={() => onSubmit(type!, ruleIndex)}
          >
            {t(`common:discard`)}
          </Button>
        }
      />
    </DialogModal>
  );
};

export default DiscardChangeModal;
