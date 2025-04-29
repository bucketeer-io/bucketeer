import { useMemo } from 'react';
import { Trans } from 'react-i18next';
import { useTranslation } from 'i18n';
import { FeatureVariation } from '@types';
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
  variations,
  onOpenRolloutVariant
}: {
  variations: FeatureVariation[];
  onOpenRolloutVariant: () => void;
}) => {
  const { t } = useTranslation(['table']);
  const variants = useMemo(
    () => variations.slice(0, variations.length > 2 ? 1 : variations.length),
    [variations]
  );
  return (
    <div className="flex items-center justify-between w-full gap-x-2 px-4 py-2 rounded-lg bg-gray-100">
      <div className="flex items-center gap-x-3">
        <div className="flex items-center gap-x-2 typo-head-bold-small text-gray-700">
          <Icon icon={IconExperiment} color="primary-500" size={'sm'} />
          <Trans
            i18nKey={'table:results.confidence-percent'}
            values={{
              percent: '100%'
            }}
          />
        </div>
        {variants.map((item, index) => (
          <div
            key={item.id}
            className="flex items-center gap-x-2 pl-3 border-l border-gray-400 typo-para-small text-gray-600"
          >
            <Trans
              i18nKey={'table:results.variant-outperformed-percent'}
              values={{
                name: item.name || item.value,
                percent: '100%'
              }}
            />
            {index === variants.length - 1 && (
              <div className="flex-center p-1 rounded bg-primary-100/30">
                <Icon icon={IconOutperformed} size="xxs" color="primary-500" />
              </div>
            )}
          </div>
        ))}

        {variations.length > 2 && (
          <Tooltip
            content={
              <div>
                {variations.slice(1, variations.length).map((item, index) => (
                  <div key={index} className="typo-para-medium text-white">
                    <Trans
                      i18nKey={'table:results.variant-outperformed-percent'}
                      values={{
                        name: item.name || item.value,
                        percent: '100%'
                      }}
                    />
                  </div>
                ))}
              </div>
            }
            trigger={
              <div className="flex items-center gap-x-2 typo-para-small text-gray-600 pl-3 border-l border-gray-400">
                <Trans
                  i18nKey={'table:results.more-variants'}
                  values={{
                    quantity: `+${variations.length - 1}`
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
