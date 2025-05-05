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
  onOpenRolloutVariant
}: {
  bestVariations: BestVariation[];
  variations: FeatureVariation[];
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
      return {
        name: variation?.name,
        value: variation?.value
      };
    },
    [variations]
  );
  return (
    <div className="flex items-center justify-between w-full gap-x-2 px-4 py-2 rounded-lg bg-gray-100">
      <div className="flex items-center gap-x-3">
        <div className="flex items-center gap-x-2 typo-head-bold-small text-gray-700 whitespace-nowrap">
          <Icon icon={IconExperiment} color="primary-500" size={'sm'} />
          <Trans
            i18nKey={'table:results.confidence-percent'}
            values={{
              percent: getPercentage(bestVariation?.probability)
            }}
          />
        </div>
        {variants?.slice(0, 2)?.map((item, index) => {
          const variation = getVariationName(item?.id);
          return (
            <div
              key={item?.id}
              className="flex items-center gap-x-2 pl-3 border-l border-gray-400 typo-para-small text-gray-600 whitespace-nowrap"
            >
              <Trans
                i18nKey={'table:results.variant-outperformed-percent'}
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
                        i18nKey={'table:results.variant-outperformed-percent'}
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
      <Button variant={'text'} onClick={onOpenRolloutVariant}>
        <Icon icon={IconOperationArrow} color="primary-500" size={'sm'} />
        {t('results.rollout-variant')}
      </Button>
    </div>
  );
};

export default ConfidenceVariants;
