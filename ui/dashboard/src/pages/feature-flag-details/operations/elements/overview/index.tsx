import { FunctionComponent } from 'react';
import { useTranslation } from 'i18n';
import { FeatureCountByStatus, IconSize } from '@types';
import { cn } from 'utils/style';
import {
  IconCalendarXL,
  IconFlagOperationXL,
  IconOperationArrowXL
} from '@icons';
import OverviewCard, { OverviewIconColor } from 'elements/overview-card';
import { OpsTypeMap } from '../../types';

interface OverviewOption {
  titleKey: string;
  countKey?: keyof FeatureCountByStatus;
  color: OverviewIconColor;
  icon: FunctionComponent;
  iconSize: IconSize;
  opsType: OpsTypeMap;
  description: string;
}

const overviewOptions: OverviewOption[] = [
  {
    titleKey: 'feature-flags.schedule',
    countKey: 'total',
    color: 'brand',
    icon: IconCalendarXL,
    iconSize: 'xl',
    opsType: OpsTypeMap.SCHEDULE,
    description: 'table:feature-flags.operations-schedule-desc'
  },
  {
    titleKey: 'feature-flags.event-rate',
    countKey: 'active',
    color: 'pink',
    icon: IconFlagOperationXL,
    iconSize: 'xl',
    opsType: OpsTypeMap.EVENT_RATE,
    description: 'table:feature-flags.operations-event-rate-desc'
  },
  {
    titleKey: 'feature-flags.progressive-rollout',
    countKey: 'inactive',
    color: 'blue',
    icon: IconOperationArrowXL,
    iconSize: 'xl',
    opsType: OpsTypeMap.ROLLOUT,
    description: 'table:feature-flags.operations-rollout-desc'
  }
];

const Overview = ({
  disabled,
  onOperationActions
}: {
  disabled: boolean;
  onOperationActions: (operationType: OpsTypeMap) => void;
}) => {
  const { t } = useTranslation(['form', 'table']);

  return (
    <div className="flex flex-wrap items-center w-full gap-6 pb-4 px-3 sm:px-6">
      {overviewOptions.map(
        ({ titleKey, color, icon, iconSize, opsType, description }, index) => (
          <OverviewCard
            key={index}
            title={
              <p className="typo-head-bold-medium text-gray-900">
                {t(titleKey)}
              </p>
            }
            description={<p>{t(description)}</p>}
            color={color}
            icon={icon}
            iconSize={iconSize}
            className={cn(
              'items-start border border-transparent min-h-full self-stretch min-w-[300px]',
              {
                'pointer-events-none': disabled
              }
            )}
            iconClassName={'p-4'}
            onClick={() => {
              if (!disabled) onOperationActions(opsType);
            }}
          />
        )
      )}
    </div>
  );
};

export default Overview;
