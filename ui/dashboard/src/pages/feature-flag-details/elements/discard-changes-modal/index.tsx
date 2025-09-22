import { Trans } from 'react-i18next';
import { IconRemoveOutlined } from 'react-icons-material-design';
import { useTranslation } from 'i18n';
import { capitalize, cn } from 'utils/style';
import { IconArrowUpDown, IconPlus, IconWarningOutline } from '@icons';
import {
  DiscardChangesStateData,
  DiscardChangesType,
  RuleOrders,
  VariationPercent
} from 'pages/feature-flag-details/targeting/types';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';
import VariationLabel from 'elements/variation-label';

interface Props {
  isOpen: boolean;
  type: DiscardChangesType | undefined;
  data: DiscardChangesStateData[];
  ruleIndex?: number;
  reorderRule?: boolean;
  actionRule?: 'new-rule' | 'edit-rule' | undefined;
  onClose: () => void;
  onSubmit: (type: DiscardChangesType, index?: number) => void;
}

interface RuleHeader extends Omit<DiscardChangesStateData, 'variationIndex'> {
  formNotify: string;
}

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
            : IconWarningOutline
    }
    size={'sm'}
    color="gray-600"
  />
);

const PrerequisiteDiscardItem = ({
  labelType,
  label,
  variationIndex,
  variation
}: DiscardChangesStateData) => {
  const { t } = useTranslation(['common', 'form']);
  return (
    <div className="flex flex-col w-full gap-1">
      <div className="flex w-full gap-x-2">
        <div className="mt-[3px]">
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
              b: <strong />
            }}
          />
        </div>
      </div>
      <div className="flex items-center gap-x-2 pl-7">
        <div className="flex-center size-fit">
          <FlagVariationPolygon index={variationIndex} />
        </div>
        <p className="typo-para-medium text-gray-700">
          {variation?.name || variation?.value}
        </p>
      </div>
    </div>
  );
};

const IndividualDiscardItem = ({
  labelType,
  label,
  variationIndex,
  variation
}: DiscardChangesStateData) => {
  const { t } = useTranslation(['common', 'form']);
  return (
    <div className="flex flex-col w-full gap-1">
      <div className="flex w-full gap-x-2">
        <div className="mt-[3px]">
          <ActionIcon labelType={labelType} />
        </div>
        <div className="typo-para-medium text-gray-700">
          <Trans
            i18nKey={'form:individual-discard-desc'}
            values={{
              action: t(capitalize(labelType.toLowerCase())),
              flagName: label
            }}
            components={{
              b: <strong />,
              variantElement: (
                <div className="flex items-center gap-x-2">
                  <div className="flex-center size-fit">
                    <FlagVariationPolygon index={variationIndex} />
                  </div>
                  <p className="typo-para-medium text-gray-700">
                    {variation?.name || variation?.value}
                  </p>
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
  valueLabel
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
          <p>
            {t('common:add-rule')}
            <strong className="px-1">{label}</strong>
            {t('common:server').toLowerCase()}:
          </p>
        ) : (
          <p>
            <Trans
              i18nKey={formNotify}
              values={{
                action: t(capitalize(labelType.toLowerCase())),
                value: changeType === 'value' ? valueLabel : ''
              }}
              components={{ b: <strong /> }}
            />
            <strong>{label}</strong>
          </p>
        )}
      </div>
    </div>
  );
};

const ReorderList = ({ ruleOrders }: { ruleOrders: RuleOrders }) => {
  const { t } = useTranslation(['common', 'form']);
  const RuleLabel = (ruleLabel: string[]) => (
    <p>
      {ruleLabel.reduce<React.ReactNode[]>((acc, label, i) => {
        if (i > 0) {
          acc.push(<b key={`and-${i}`}> {t('common:and').toLowerCase()} </b>);
        }
        acc.push(<span key={`label-${i}`}>{label}</span>);
        return acc;
      }, [])}
      <strong className="pl-1">{t('common:server').toLowerCase()}</strong>
      <span>:</span>
    </p>
  );

  return (
    <div className="pl-7">
      <ol className="leading-7 space-y-2">
        {ruleOrders.labels.map((ruleLabel: string[], index: number) => (
          <li key={`label-${index}`}>
            <div className="flex gap-1">
              <span>{index + 1}.</span>
              {RuleLabel(ruleLabel)}
            </div>
            {ruleOrders.variations[index].map((v, vIndex: number) => (
              <div
                className="flex items-center gap-1 ml-4"
                key={`variation-${index}-${vIndex}`}
              >
                <VariationLabel
                  label={v.variation}
                  index={v.variationIndex || 0}
                />
                {v.weight !== null && (
                  <p className="text-gray-700"> - ({v.weight}%)</p>
                )}
              </div>
            ))}
          </li>
        ))}
      </ol>
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
    <div className="pl-6">
      <div className="flex w-full gap-x-2">
        <div className="mt-[3px]">
          <ActionIcon labelType={'REMOVE'} />
        </div>
        <p>
          <Trans
            i18nKey={t('form:custom-rule-audience-not-include-desc')}
            values={{
              percent: audienceExcluded.weight,
              variation: audienceExcluded.variation
            }}
            components={{ b: <strong /> }}
          />
        </p>
      </div>
      <div className="flex w-full gap-x-2">
        <div className="mt-[3px]">
          <ActionIcon labelType={'ADD'} />
        </div>
        <Trans
          i18nKey={t('form:custom-rule-audience-include-desc')}
          values={{
            percent: Number(100 - (audienceExcluded.weight || 0))
          }}
          components={{ b: <strong /> }}
        />
      </div>
    </div>
  );
};

const StrategyList = ({
  variationPercent
}: {
  variationPercent: VariationPercent[];
}) => (
  <div className="pl-7 font-thin leading-7">
    {variationPercent.map((vp, index: number) => (
      <div className="flex items-center gap-1" key={index}>
        <VariationLabel label={vp.variation} index={vp.variationIndex || 0} />
        {vp.weight !== null && (
          <p className="text-gray-700">
            {' '}
            - <span>({vp.weight?.toString()}%)</span>
          </p>
        )}
      </div>
    ))}
  </div>
);

const CustomRuleDiscardItem = ({
  labelType,
  label,
  changeType,
  valueLabel,
  ruleOrders,
  variationPercent,
  audienceExcluded,
  isAddNew
}: DiscardChangesStateData) => {
  const { t } = useTranslation(['common', 'form']);
  const isReorder = changeType === 'reorder';

  const notifyMap: Record<string, string> = {
    clause: 'form:custom-rule-discard-desc',
    strategy: 'form:custom-rule-strategy-discard-desc',
    reorder: 'form:custom-rule-reorder-discard-desc',
    audience: 'form:custom-rule-audience-discard-desc',
    'new-rule': 'form:custom-rule-strategy-add-new-discard-desc',
    'default-strategy': 'form:custom-rule-default-strategy-discard-desc',
    'default-audience': 'form:custom-rule-default-audience-discard-desc'
  };

  const formNotify = t(
    notifyMap[changeType || 'clause'] ||
      'form:custom-rule-clause-value-discard-desc'
  );

  const showReorder =
    changeType === 'reorder' && !isAddNew && !!ruleOrders?.labels?.length;
  const showAudience =
    ['default-audience', 'audience'].includes(changeType || '') &&
    audienceExcluded;
  const showVariationPercent =
    ['strategy', 'default-strategy', 'new-rule'].includes(changeType || '') &&
    !!variationPercent?.length;
  return (
    <div className={cn('flex flex-col w-full gap-1 pl-4', isReorder && 'pl-0')}>
      <RuleHeader
        isAddNew={isAddNew}
        labelType={labelType}
        label={label}
        changeType={changeType}
        formNotify={formNotify}
        valueLabel={valueLabel}
      />

      {showReorder && <ReorderList ruleOrders={ruleOrders} />}
      {showAudience && <AudienceChange audienceExcluded={audienceExcluded} />}
      {showVariationPercent && (
        <StrategyList variationPercent={variationPercent} />
      )}
    </div>
  );
};

const DiscardChangeModal = ({
  isOpen,
  type,
  data,
  ruleIndex,
  actionRule,
  reorderRule = false,
  onClose,
  onSubmit
}: Props) => {
  const { t } = useTranslation(['common', 'form']);
  const isEdit = data.find(
    item => item.changeType !== 'new-rule' && item.changeType !== 'reorder'
  );
  const isAddNew = data.find(item => item.changeType === 'new-rule');
  return (
    <DialogModal
      className="w-[500px]"
      title={t('form:discard-unsaved-changes')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="flex flex-col w-full gap-y-4 p-5 max-h-[500px] overflow-y-auto small-scroll">
        {reorderRule &&
          data.map(
            (item, index) =>
              item.changeType === 'reorder' && (
                <CustomRuleDiscardItem
                  key={index}
                  {...item}
                  ruleIndex={Number(ruleIndex)}
                  isAddNew={!!isAddNew}
                />
              )
          )}
        <>
          {ruleIndex! >= 0 && actionRule === 'edit-rule' && isEdit && (
            <div className="flex gap-1 items-center">
              <Trans i18nKey={'common:edit-rule'} />
              <Trans
                i18nKey={'table:feature-flags.rule-index'}
                values={{
                  index: ruleIndex! + 1
                }}
              />
            </div>
          )}
          {data.map((item, index) => {
            const { PREREQUISITE, INDIVIDUAL, CUSTOM, DEFAULT } =
              DiscardChangesType;
            if (type === PREREQUISITE)
              return <PrerequisiteDiscardItem key={index} {...item} />;
            if (type === INDIVIDUAL)
              return <IndividualDiscardItem key={index} {...item} />;
            if (
              (type === CUSTOM || type === DEFAULT) &&
              item.changeType !== 'reorder'
            )
              return (
                <CustomRuleDiscardItem
                  key={index}
                  {...item}
                  ruleIndex={Number(ruleIndex)}
                />
              );
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
