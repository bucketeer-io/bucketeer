import { useCallback, useEffect, useMemo, useState } from 'react';
import { featureUpdater } from '@api/features';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import { Environment, Feature, StrategyType } from '@types';
import { IconPercentage } from '@icons';
import { DefaultRuleSchema } from 'pages/feature-flag-details/targeting/form-schema';
import { handleGetDefaultRuleStrategy } from 'pages/feature-flag-details/targeting/utils';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Dropdown from 'components/dropdown';
import Icon from 'components/icon';
import OperationInfoCard from './operation-info-card';

export type OperationActiveModalProps = {
  editable: boolean;
  isDeleting?: boolean;
  environment: Environment;
  feature: Feature;
  loading?: boolean;
  onClose: () => void;
  onActionOperation: () => void;
  refetchFeature: () => void;
};

const CURRENT_PERCENTAGE = 'CURRENT_PERCENTAGE' as const;

const OperationActiveModal = ({
  editable,
  isDeleting = false,
  environment,
  feature,
  loading,
  onClose,
  onActionOperation,
  refetchFeature
}: OperationActiveModalProps) => {
  const { t } = useTranslation(['common', 'table', 'form']);
  const { errorNotify } = useToast();

  const [defaultRule, setDefaultRule] = useState<DefaultRuleSchema>({
    ...feature.defaultStrategy,
    currentOption: CURRENT_PERCENTAGE
  });

  const variationOptions = useMemo(() => {
    const variations = feature.variations.map((item, index) => ({
      label: (
        <div className="flex items-center gap-x-2 pl-0.5">
          <FlagVariationPolygon index={index} />
          <p className="-mt-0.5 truncate">{item.name}</p>
        </div>
      ),
      value: item.id,
      type: StrategyType.FIXED,
      variationName: item.name
    }));
    return [
      ...variations,
      {
        label: (
          <div className="flex items-center gap-x-2 pl-0.5">
            <Icon icon={IconPercentage} />
            <p className="-mt-0.5 truncate">{t('form:current-percentage')}</p>
          </div>
        ),
        value: CURRENT_PERCENTAGE,
        type: CURRENT_PERCENTAGE,
        variationName: t('form:current-percentage')
      }
    ];
  }, [feature.variations]);

  const handleSelectStrategy = (value: string) => {
    setDefaultRule(prev => {
      if (value === CURRENT_PERCENTAGE) {
        return {
          ...prev,
          currentOption: value
        };
      }
      return {
        ...prev,
        type: StrategyType.FIXED,
        currentOption: value,
        fixedStrategy: {
          variation: value
        }
      };
    });
  };

  const onConfirm = useCallback(async () => {
    if (!editable) return;
    try {
      onActionOperation();
      const isCurrentPercentage =
        defaultRule.currentOption === CURRENT_PERCENTAGE;

      const sameStrategyFixed =
        feature.defaultStrategy.type === StrategyType.FIXED
          ? feature.defaultStrategy.fixedStrategy.variation ===
            defaultRule.fixedStrategy.variation
          : false;

      if (isCurrentPercentage || sameStrategyFixed) return;
      await featureUpdater({
        id: feature.id,
        environmentId: environment.id,
        comment: t('form:feature-flags.update-default-strategy'),
        defaultStrategy: handleGetDefaultRuleStrategy(defaultRule)
      });
      refetchFeature();
    } catch (error) {
      errorNotify(error);
    }
  }, [
    editable,
    feature,
    environment,
    defaultRule,
    onActionOperation,
    refetchFeature
  ]);

  useEffect(() => {
    setDefaultRule({
      ...feature.defaultStrategy,
      currentOption: CURRENT_PERCENTAGE
    });
  }, [feature.defaultStrategy]);

  return (
    <>
      <div className="flex flex-col w-full items-start px-5 py-4 gap-1">
        <p className="capitalize">{t('table:feature-flags.serve')}</p>
        <Dropdown
          className="w-full max-w-[250px]"
          contentClassName="max-w-[550px]"
          value={defaultRule.currentOption}
          options={variationOptions}
          onChange={value => handleSelectStrategy(value as string)}
        />

        <OperationInfoCard
          className="my-4"
          title={t(
            isDeleting
              ? 'form:operation.confirm-delete-active-rollout-title'
              : 'form:operation.confirm-stop-title'
          )}
          description={t(
            isDeleting
              ? 'form:operation.confirm-delete-active-rollout-desc'
              : 'form:operation.confirm-stop-rollout-desc'
          )}
        />
      </div>
      <ButtonBar
        secondaryButton={
          <Button loading={loading} onClick={onConfirm} disabled={!editable}>
            {t(`submit`)}
          </Button>
        }
        primaryButton={
          <Button onClick={onClose} variant="secondary">
            {t(`cancel`)}
          </Button>
        }
      />
    </>
  );
};

export default OperationActiveModal;
