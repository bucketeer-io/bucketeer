import { FunctionComponent } from 'react';
import { useTranslation } from 'i18n';
import { FeatureCountByStatus, IconSize } from '@types';
import { IconCalendar, IconFlagOperation, IconOperationArrow } from '@icons';
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
    icon: IconCalendar,
    iconSize: 'xl',
    opsType: OpsTypeMap.SCHEDULE,
    description: 'table:feature-flags.operations-schedule-desc'
  },
  {
    titleKey: 'feature-flags.event-rate',
    countKey: 'active',
    color: 'pink',
    icon: IconFlagOperation,
    iconSize: 'xl',
    opsType: OpsTypeMap.EVENT_RATE,
    description: 'table:feature-flags.operations-event-rate-desc'
  },
  {
    titleKey: 'feature-flags.progressive-rollout',
    countKey: 'inactive',
    color: 'blue',
    icon: IconOperationArrow,
    iconSize: 'xl',
    opsType: OpsTypeMap.ROLLOUT,
    description: 'table:feature-flags.operations-rollout-desc'
  }
];

const Overview = ({
  onOperationActions
}: {
  onOperationActions: (operationType: OpsTypeMap) => void;
}) => {
  const { t } = useTranslation(['form', 'table']);

  return (
    <div className="flex flex-col w-full gap-6">
      <p className="typo-head-bold-big text-gray-800">
        {t('table:feature-flags:operations-desc')}
      </p>
      <div className="flex flex-wrap items-center w-full gap-6 pb-8">
        {overviewOptions.map(
          (
            { titleKey, color, icon, iconSize, opsType, description },
            index
          ) => (
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
              className="border border-transparent"
              iconClassName={'p-4'}
              onClick={() => onOperationActions(opsType)}
            />
          )
        )}
      </div>
    </div>
  );
};

export default Overview;
