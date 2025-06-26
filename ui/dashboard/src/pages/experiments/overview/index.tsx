import { FunctionComponent } from 'react';
import { useTranslation } from 'i18n';
import { ExperimentCollection, ExperimentStatus } from '@types';
import { cn } from 'utils/style';
import {
  IconExperiment, // IconNotStartedExperiment,
  IconStoppedExperiment,
  IconWaitingExperiment
} from '@icons';
import OverviewCard, { OverviewIconColor } from 'elements/overview-card';
import { SummaryType } from '../types';

interface OverviewOption {
  titleKey: string;
  countKey?: keyof ExperimentCollection['summary'];
  color: OverviewIconColor;
  icon: FunctionComponent;
  summaryFilterValue: SummaryType;
  filterValues: ExperimentStatus[];
}

const overviewOptions: OverviewOption[] = [
  {
    titleKey: 'experiment.scheduled',
    countKey: 'totalWaitingCount',
    color: 'orange',
    icon: IconWaitingExperiment,
    summaryFilterValue: 'scheduled',
    filterValues: ['WAITING']
  },
  {
    titleKey: 'experiment.running',
    countKey: 'totalRunningCount',
    color: 'brand',
    icon: IconExperiment,
    summaryFilterValue: 'running',
    filterValues: ['RUNNING']
  },
  {
    titleKey: 'experiment.stopped',
    countKey: 'totalStoppedCount',
    color: 'red',
    icon: IconStoppedExperiment,
    summaryFilterValue: 'stopped',
    filterValues: ['STOPPED', 'FORCE_STOPPED']
  }
  // {
  //   titleKey: 'experiment.not-started',
  //   countKey: undefined,
  //   color: 'gray',
  //   icon: IconNotStartedExperiment,
  //   summaryFilterValue: 'not-started',
  //   filterValues: ['NOT_STARTED']
  // }
];

const Overview = ({
  summary,
  filterBySummary,
  onChangeFilters
}: {
  summary?: ExperimentCollection['summary'];
  filterBySummary?: SummaryType;
  onChangeFilters: (
    statuses: ExperimentStatus[],
    summaryFilterValue: SummaryType
  ) => void;
}) => {
  const { t } = useTranslation(['table']);

  return (
    <div className="flex flex-wrap items-center w-full gap-6 pb-8 px-6 lg:pr-2">
      {overviewOptions.map((item, index) => (
        <OverviewCard
          key={index}
          title={t(item.titleKey)}
          count={summary && item.countKey ? Number(summary[item.countKey]) : 0}
          color={item.color}
          icon={item.icon}
          className={cn('border border-transparent', {
            'border-gray-300':
              filterBySummary && item.summaryFilterValue === filterBySummary
          })}
          onClick={() =>
            onChangeFilters(item.filterValues, item.summaryFilterValue)
          }
        />
      ))}
    </div>
  );
};

export default Overview;
