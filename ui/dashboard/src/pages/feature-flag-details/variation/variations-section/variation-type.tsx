import { useMemo } from 'react';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import { IconFlagSwitch, IconInfo } from '@icons';
import { VariationTypeTooltip } from 'pages/feature-flags/collection-layout/elements';
import Icon from 'components/icon';
import { VariationProps } from '..';

const VariationType = ({ feature }: VariationProps) => {
  const { t } = useTranslation(['table']);
  const isJSON = useMemo(() => feature.variationType === 'JSON', [feature]);
  const isBoolean = useMemo(
    () => feature.variationType === 'BOOLEAN',
    [feature]
  );

  return (
    <div className="flex items-center gap-x-2">
      <p className="typo-para-medium text-gray-700">
        {t('table:feature-flags.variation-type')}
      </p>
      <VariationTypeTooltip
        align="center"
        trigger={
          <div className="flex-center h-full">
            <Icon
              icon={IconInfo}
              size={'xxs'}
              color="gray-500"
              className="flex-center size-4"
            />
          </div>
        }
        variationType={feature.variationType}
      />
      {isBoolean && <Icon icon={IconFlagSwitch} />}
      <p
        className={cn('typo-para-small text-gray-600 capitalize', {
          uppercase: isJSON
        })}
      >
        {feature?.variationType?.toLowerCase()}
      </p>
    </div>
  );
};

export default VariationType;
