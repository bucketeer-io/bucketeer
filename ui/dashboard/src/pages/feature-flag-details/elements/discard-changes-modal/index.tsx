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

const notifyMap: Record<string, string> = {
  clause: 'form:custom-rule-discard-desc',
  strategy: 'form:custom-rule-strategy-discard-desc',
  audience: 'form:custom-rule-audience-discard-desc',
  'new-rule': 'form:custom-rule-strategy-add-new-discard-desc',
  'default-strategy': 'form:custom-rule-default-strategy-discard-desc',
  'default-audience': 'form:custom-rule-default-audience-discard-desc'
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
  const { t } = useTranslation(['common', 'form']);
  return (
    <div className="flex flex-col w-full pl-4 gap-1">
      <div className="flex w-full gap-x-2">
        <div className="flex items-start mt-1">
          <ActionIcon labelType={labelType} />
        </div>
        <div className="typo-para-medium text-gray-700">
          <Trans
            i18nKey={'form:prerequisite-discard-desc'}
            values={{
              action: t(capitalize(labelType.toLowerCase())),
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
  label,
  variationIndex,
  variation
}: DiscardChangesStateData) => {
  const { t } = useTranslation(['common', 'form']);
  return (
    <div className="flex flex-col w-full pl-4 gap-1">
      <div className="flex w-full gap-x-2">
        <div className="flex items-start mt-1">
          <ActionIcon labelType={labelType} />
        </div>
        <div className="typo-para-medium text-gray-700 ">
          <Trans
            i18nKey={'form:individual-discard-desc'}
            values={{
              action: t(capitalize(labelType.toLowerCase())),
              flagName: label
            }}
            components={{
              b: <strong className="break-all" />,
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

  return (
    <div className="flex w-full gap-x-2">
      <div className="mt-[3px]">
        <ActionIcon labelType={labelType} />
      </div>
      <div className="typo-para-medium text-gray-700">
        {changeType === 'new-rule' ? (
          <div className="inline items-center">
            {t('common:add-rule')}
            <strong className="px-1">{label}</strong>
            {t('common:serve').toLowerCase()}
            <StrategyList variations={variations} />
          </div>
        ) : (
          <p>
            <Trans
              i18nKey={formNotify}
              values={{
                action: t(capitalize(labelType?.toLowerCase())),
                value: changeType === 'value' ? valueLabel : ''
              }}
              components={{ b: <strong /> }}
            />
            <strong>{label}</strong>
            <StrategyList variations={variations} />
          </p>
        )}
      </div>
    </div>
  );
};

const AudienceChange = ({
  audienceExcluded
}: {
  audienceExcluded: VariationPercent;
}) => {
  const { t } = useTranslation(['common', 'form']);
  return (
    <div className="pl-6 text-gray-700">
      <div className="flex w-full gap-x-2 items-center">
        <div className="mt-[6px]">
          <ActionIcon labelType={'REMOVE'} />
        </div>
        <p className="inline">
          <Trans
            i18nKey={'form:custom-rule-audience-not-include-desc'}
            values={{
              percent: audienceExcluded.weight,
              variation: audienceExcluded.variation
            }}
            components={{
              b: <strong />,
              variantElement: (
                <div className="inline-flex items-center gap-x-1">
                  <div className="flex-center size-fit">
                    {!isNil(audienceExcluded.variationIndex) && (
                      <FlagVariationPolygon
                        index={audienceExcluded.variationIndex || 0}
                      />
                    )}
                  </div>
                  <p>{audienceExcluded.variation}</p>
                </div>
              )
            }}
          />
        </p>
      </div>
      <div className="flex w-full gap-x-2 items-center">
        <div className="mt-[6px]">
          <ActionIcon labelType={'ADD'} />
        </div>
        <p className="inline">
          <Trans
            i18nKey={t('form:custom-rule-audience-include-desc')}
            values={{
              percent: Number(100 - (audienceExcluded.weight || 0))
            }}
            components={{ b: <strong /> }}
          />
        </p>
      </div>
    </div>
  );
};

const StrategyList = ({ variations }: { variations?: VariationPercent[] }) =>
  !!variations?.length && (
    <div className="ml-1 inline-flex items-center gap-x-1 leading-[1px]">
      {variations.map((v, index) => (
        <div key={index} className="inline-flex items-center gap-x-1">
          <div className="flex-center size-fit">
            <FlagVariationPolygon index={v.variationIndex || 0} />
          </div>
          <p>{v.variation}</p>
          {!isNil(v.weight) && (
            <p className="text-gray-700">
              <span>
                ({v.weight?.toString()}%)
                {index !== variations.length - 1 && ','}
              </span>
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
  isAddNew
}: DiscardChangesStateData) => {
  const { t } = useTranslation(['common', 'form']);

  const formNotify = t(
    notifyMap[changeType || 'clause'] ||
      'form:custom-rule-clause-value-discard-desc'
  );

  const showAudience =
    ['default-audience', 'audience'].includes(changeType || '') &&
    audienceExcluded;
  const showVariationPercent =
    ['strategy', 'default-strategy', 'new-rule'].includes(changeType || '') &&
    !!variationPercent?.length;
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

      {showAudience && <AudienceChange audienceExcluded={audienceExcluded} />}
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
      className="w-full max-w-[600px]"
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
