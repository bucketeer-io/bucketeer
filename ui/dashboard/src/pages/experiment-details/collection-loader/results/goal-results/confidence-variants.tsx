import { useCallback, useMemo } from 'react';
import { Trans } from 'react-i18next';
import { useTranslation } from 'i18n';
import { BestVariation, FeatureVariation } from '@types';
import {
  IconExperiment,
  IconInfo,
  IconOperationArrow,
  IconOutperformed
} from '@icons';
import Button from 'components/button';
import Icon from 'components/icon';
import { Tooltip } from 'components/tooltip';

const ConfidenceVariants = ({
  bestVariations,
  variations,
  safeToStop,
  onOpenRolloutVariant
}: {
  bestVariations: BestVariation[];
  variations: FeatureVariation[];
  safeToStop: boolean;
  onOpenRolloutVariant: () => void;
}) => {
  const { t } = useTranslation(['table']);
  const bestVariation = useMemo(
    () => bestVariations.find(item => item.isBest),
    [bestVariations]
  );
  const getPercentage = useCallback((probability?: number) => {
    const percent = probability ? probability * 100 : 0;
    return `${percent.toFixed(1)}%`;
  }, []);

  const variants = useMemo(() => {
    const results = bestVariation ? [bestVariation] : [];
    bestVariations.forEach(item => {
      if (!results.find(v => v?.id === item.id)) {
        results.push(item);
      }
    });
    return results;
  }, [bestVariations, bestVariation]);

  const getVariationName = useCallback(
    (id?: string) => {
      if (!id) return null;
      const variation = variations.find(item => item.id === id);
      return { name: variation?.name, value: variation?.value };
    },
    [variations]
  );

  return (
    <div className="flex items-center justify-between w-full gap-x-2 px-4 py-2 rounded-lg bg-gray-100">
      {/* flex-1 min-w-0 allows the left section to shrink when space is tight */}
      <div className="flex items-center gap-x-3 flex-1 min-w-0">
        {/* ── Status badge — primary signal ── */}
        {safeToStop ? (
          <span className="flex items-center gap-x-2 typo-head-bold-small text-primary-500 whitespace-nowrap">
            <Icon icon={IconExperiment} color="primary-500" size={'sm'} />
            {t('results.status-safe')}
          </span>
        ) : (
          <Tooltip
            content={
              <p className="typo-para-small text-white max-w-xs">
                <Trans
                  i18nKey={'table:results.monitoring-tooltip'}
                  values={{
                    percent: getPercentage(bestVariation?.probability)
                  }}
                />
              </p>
            }
            trigger={
              <button
                type="button"
                className="flex items-center gap-x-2 typo-head-bold-small text-gray-600 whitespace-nowrap cursor-default bg-transparent border-0 p-0"
              >
                <Icon icon={IconExperiment} color="gray-600" size={'sm'} />
                {t('results.status-monitoring')}
                <Icon icon={IconInfo} size="xxs" color="gray-600" />
              </button>
            }
          />
        )}

        {/* ── Per-variation chips — probability once, clearly labeled ── */}
        {variants?.slice(0, 2)?.map((item, index) => {
          const variation = getVariationName(item?.id);
          return (
            <div
              key={item?.id}
              className="flex items-center gap-x-2 pl-3 border-l border-gray-400 typo-para-small text-gray-600 whitespace-nowrap"
            >
              <Trans
                i18nKey={'table:results.variant-chance-beats-baseline'}
                values={{
                  name: variation?.name || variation?.value,
                  percent: getPercentage(item.probability)
                }}
              />
              {index === 0 && (
                <div className="flex-center p-1 rounded bg-primary-100/30">
                  <Icon
                    icon={IconOutperformed}
                    size="xxs"
                    color="primary-500"
                  />
                </div>
              )}
            </div>
          );
        })}

        {variants.length > 2 && (
          <Tooltip
            content={
              <div>
                {variants.slice(2, variants.length).map((item, index) => {
                  const variation = getVariationName(item?.id);
                  return (
                    <div key={index} className="typo-para-medium text-white">
                      <Trans
                        i18nKey={'table:results.variant-chance-beats-baseline'}
                        values={{
                          name: variation?.name || variation?.value,
                          percent: getPercentage(item.probability)
                        }}
                      />
                    </div>
                  );
                })}
              </div>
            }
            trigger={
              <div className="flex items-center gap-x-2 typo-para-small text-gray-600 pl-3 border-l border-gray-400 whitespace-nowrap">
                <Trans
                  i18nKey={'table:results.more-variants'}
                  values={{
                    quantity: `+${variants.slice(2, variants.length).length}`
                  }}
                />
                <Icon icon={IconInfo} size={'xxs'} />
              </div>
            }
          />
        )}
      </div>

      {/* ── Rollout CTA — shrink-0 keeps it fixed regardless of left width ── */}
      <div className="shrink-0">
        {safeToStop ? (
          <Button
            variant={'text'}
            onClick={onOpenRolloutVariant}
            className="typo-para-small"
          >
            <Icon icon={IconOperationArrow} color="primary-500" size={'sm'} />
            {t('results.rollout-variant')}
          </Button>
        ) : (
          <Tooltip
            content={
              <p className="typo-para-small text-white max-w-xs">
                {t('results.rollout-not-safe-tooltip', {
                  percent: getPercentage(bestVariation?.probability)
                })}
              </p>
            }
            trigger={
              <Button
                variant={'text'}
                aria-disabled="true"
                className="typo-para-small text-gray-500 cursor-not-allowed hover:text-gray-500"
              >
                <Icon icon={IconOperationArrow} color="gray-600" size={'sm'} />
                {t('results.rollout-variant')}
              </Button>
            }
          />
        )}
      </div>
    </div>
  );
};

export default ConfidenceVariants;
