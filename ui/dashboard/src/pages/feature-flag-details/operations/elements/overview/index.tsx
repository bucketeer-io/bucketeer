import { FunctionComponent } from 'react';
import { useTranslation } from 'i18n';
import { ExperimentStatus, FeatureCountByStatus, IconSize } from '@types';
import { cn } from 'utils/style';
import { IconCalendar, IconFlagOperation, IconOperationArrow } from '@icons';
import OverviewCard, { OverviewIconColor } from 'elements/overview-card';
import { OpsTypeMap } from '../../types';

interface OverviewOption {
  titleKey: string;
  countKey?: keyof FeatureCountByStatus;
  color: OverviewIconColor;
  icon: FunctionComponent;
  iconSize: IconSize;
  summaryFilterValue: OpsTypeMap;
  filterValues: ExperimentStatus[];
  description: string;
}

const overviewOptions: OverviewOption[] = [
  {
    titleKey: 'feature-flags.schedule',
    countKey: 'total',
    color: 'brand',
    icon: IconCalendar,
    iconSize: 'xl',
    summaryFilterValue: OpsTypeMap.SCHEDULE,
    filterValues: [],
    description: 'table:feature-flags.operations-schedule-desc'
  },
  {
    titleKey: 'feature-flags.event-rate',
    countKey: 'active',
    color: 'pink',
    icon: IconFlagOperation,
    iconSize: 'xl',
    summaryFilterValue: OpsTypeMap.EVENT_RATE,
    filterValues: [],
    description: 'table:feature-flags.operations-event-rate-desc'
  },
  {
    titleKey: 'feature-flags.progressive-rollout',
    countKey: 'inactive',
    color: 'blue',
    icon: IconOperationArrow,
    iconSize: 'xl',
    summaryFilterValue: OpsTypeMap.ROLLOUT,
    filterValues: [],
    description: 'table:feature-flags.operations-rollout-desc'
  }
];

const Overview = ({
  summary,
  filterBySummary,
  onChangeFilters
}: {
  summary?: FeatureCountByStatus;
  filterBySummary?: OpsTypeMap;
  onChangeFilters: (
    statuses: ExperimentStatus[],
    summaryFilterValue: OpsTypeMap
  ) => void;
}) => {
  const { t } = useTranslation(['form', 'table']);
  console.log({ summary });
  return (
    <div className="flex flex-col w-full gap-6">
      <p className="typo-head-bold-big text-gray-800">
        {t('table:feature-flags:operations-desc')}
      </p>
      <div className="flex flex-wrap items-center w-full gap-6 pb-8">
        {overviewOptions.map(
          (
            {
              titleKey,
              color,
              icon,
              iconSize,
              summaryFilterValue,
              filterValues,
              description
            },
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
              className={cn('border border-transparent', {
                'border-gray-300':
                  filterBySummary && summaryFilterValue === filterBySummary
              })}
              iconClassName={'p-4'}
              onClick={() => onChangeFilters(filterValues, summaryFilterValue)}
            />
          )
        )}
      </div>
    </div>
  );
};

export default Overview;
