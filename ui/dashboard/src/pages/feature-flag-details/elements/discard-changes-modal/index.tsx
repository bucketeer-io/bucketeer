import { Trans } from 'react-i18next';
import {
  IconRemoveOutlined,
  IconUpdateOutlined
} from 'react-icons-material-design';
import { useTranslation } from 'i18n';
import { capitalize } from 'utils/style';
import { IconPlus } from '@icons';
import {
  DiscardChangesStateData,
  DiscardChangesType
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
  onClose: () => void;
  onSubmit: (type: DiscardChangesType, index?: number) => void;
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
          : IconUpdateOutlined
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

const CustomRuleDiscardItem = ({
  labelType,
  label,
  changeField,
  valueLabel,
  variationPercent
}: DiscardChangesStateData) => {
  const { t } = useTranslation(['common', 'form']);
  const formNotify =
    changeField === 'clause'
      ? t('form:custom-rule-discard-desc')
      : changeField === 'strategy'
        ? t('form:custom-rule-strategy-discard-desc')
        : t('form:custom-rule-clause-value-discard-desc');
  return (
    <div className="flex flex-col w-full gap-1 pl-4">
      <div className="flex w-full gap-x-2">
        <div className="mt-[3px]">
          <ActionIcon labelType={labelType} />
        </div>
        <div className="typo-para-medium text-gray-700">
          <Trans
            i18nKey={formNotify}
            values={{
              action: t(capitalize(labelType.toLowerCase())),
              value: changeField === 'value' ? valueLabel : '',
              clauseLabel: label
            }}
            components={{
              b: <strong />
            }}
          />
        </div>
      </div>
      {changeField === 'strategy' && variationPercent?.length && (
        <div className="pl-7">
          {variationPercent.map(vp => (
            <div className="flex items-center gap-1" key={vp.variation}>
              <VariationLabel label={vp.variation} index={0} />
              {vp.weight !== undefined && (
                <p className="text-gray-700"> - ({vp.weight}%)</p>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

const DiscardChangeModal = ({
  isOpen,
  type,
  data,
  ruleIndex,
  onClose,
  onSubmit
}: Props) => {
  const { t } = useTranslation(['common', 'form']);
  return (
    <DialogModal
      className="w-[500px]"
      title={t('form:discard-unsaved-changes')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="flex flex-col w-full gap-y-4 p-5 max-h-[500px] overflow-y-auto small-scroll">
        {!!ruleIndex && (
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
          if (type === DiscardChangesType.PREREQUISITE)
            return <PrerequisiteDiscardItem key={index} {...item} />;
          if (type === DiscardChangesType.INDIVIDUAL)
            return <IndividualDiscardItem key={index} {...item} />;
          if (type === DiscardChangesType.CUSTOM)
            return <CustomRuleDiscardItem key={index} {...item} />;
        })}
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
