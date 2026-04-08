import { useCallback, useMemo, useState } from 'react';
import { Trans } from 'react-i18next';
import { InsightsTimeSeriesFetcherParams } from '@api/insight';
import {
  useQueryInsightsErrorRates,
  useQueryInsightsEvaluations,
  useQueryInsightsLatency,
  useQueryInsightsRequests
} from '@queries/insights';
import {
  BarElement,
  CategoryScale,
  Chart as ChartJS,
  Filler,
  Legend,
  LinearScale,
  LineElement,
  PointElement,
  TimeScale,
  TimeSeriesScale,
  Tooltip
} from 'chart.js';
import 'chartjs-adapter-luxon';
import { ALL } from 'constants/insight';
import useOptions from 'hooks/use-options';
import { useTranslation } from 'i18n';
import { DateTime } from 'luxon';
import {
  Environment,
  InsightApiId,
  InsightsMonthlySummaryResponse,
  InsightSourceId,
  Project
} from '@types';
import { exportMonthlySummaryCSV } from 'utils/csv-export';
import { IconDownload, IconWatch } from '@icons';
import Button from 'components/button';
import Dropdown, { DropdownOption } from 'components/dropdown';
import Icon from 'components/icon';
import DropdownMenuWithSearch from 'elements/dropdown-with-search';
import DateCustom from './elements/DateCustom';
import ChartDescription from './elements/DescriptionChart';
import MonthlyBarChart from './elements/MonthlyBarChart';
import TimeSeriesLineChart from './elements/TimeSeriesLineChart';
import { InsightsFilters, TimeRangePreset, formatYAxis } from './utils';

ChartJS.register(
  CategoryScale,
  LinearScale,
  BarElement,
  LineElement,
  TimeScale,
  TimeSeriesScale,
  PointElement,
  Tooltip,
  Legend,
  Filler
);

// Computed once at module level — monthly labels don't need to be reactive
const now = DateTime.now();
const MONTHS = Array.from({ length: 12 }, (_, i) =>
  now.minus({ months: 11 - i }).toFormat('yyyyMM')
);
const MONTH_LABELS = MONTHS.map(m =>
  DateTime.fromFormat(m, 'yyyyMM').toFormat('MMM')
);

interface PageContentProps {
  projects: Project[];
  environments: Environment[];
  monthlySummary?: InsightsMonthlySummaryResponse;
  monthlySummaryLoading: boolean;
  timeRangeParams: InsightsTimeSeriesFetcherParams;
  filters: InsightsFilters;
  onFiltersChange: (filters: InsightsFilters) => void;
  onProjectChange: (projectId: string) => void;
  queriesEnabled: boolean;
}

const PageContent = ({
  projects,
  environments,
  monthlySummary,
  monthlySummaryLoading,
  timeRangeParams,
  filters,
  onFiltersChange,
  onProjectChange,
  queriesEnabled
}: PageContentProps) => {
  const { t } = useTranslation(['common']);
  const [isDatePickerOpen, setIsDatePickerOpen] = useState(false);
  const { sourceIdOptions, apiIdOptions } = useOptions();

  const dateRangeLabel = useMemo(() => {
    if (!filters.customStartAt || !filters.customEndAt) return '';
    return `${DateTime.fromSeconds(Number(filters.customStartAt)).toFormat('MMM d, yyyy')} 
    - ${DateTime.fromSeconds(Number(filters.customEndAt)).toFormat('MMM d, yyyy')}`;
  }, [filters.customStartAt, filters.customEndAt]);

  const timeRangeOptions: DropdownOption[] = useMemo(
    () => [
      { label: t('insights.last-1-hour'), value: '1h' },
      { label: t('insights.last-6-hours'), value: '6h' },
      { label: t('insights.last-24-hours'), value: '24h' },
      { label: t('insights.last-7-days'), value: '7d' },
      { label: t('insights.last-30-days'), value: '30d' },
      { label: t('insights.this-month'), value: 'this_month' },
      { label: t('insights.date-range'), value: 'date_range' }
    ],
    [t]
  );

  const projectOptions: DropdownOption[] = useMemo(
    () => [
      { label: t('insights.all-projects'), value: ALL },
      ...projects.map(p => ({ label: p.name, value: p.id }))
    ],
    [projects, t]
  );

  const environmentOptions: DropdownOption[] = useMemo(() => {
    return [
      { label: t('insights.all-environments'), value: ALL },
      ...environments.map(e => ({
        label: e.name,
        value: e.id === '' ? 'production' : e.id
      }))
    ];
  }, [environments, t]);

  const environmentNameMap = useMemo(
    () =>
      Object.fromEntries(environments.map(e => [e.id, e.name])) as Record<
        string,
        string
      >,
    [environments]
  );

  const { data: latencyData, isLoading: latencyLoading } =
    useQueryInsightsLatency({
      params: timeRangeParams,
      enabled: queriesEnabled
    });

  const { data: requestsData, isLoading: requestsLoading } =
    useQueryInsightsRequests({
      params: timeRangeParams,
      enabled: queriesEnabled
    });

  const { data: evaluationsData, isLoading: evaluationsLoading } =
    useQueryInsightsEvaluations({
      params: timeRangeParams,
      enabled: queriesEnabled
    });

  const { data: errorRatesData, isLoading: errorRatesLoading } =
    useQueryInsightsErrorRates({
      params: timeRangeParams,
      enabled: queriesEnabled
    });

  const isDisableExportCSV = useMemo(
    () => !monthlySummary || monthlySummary.series.length <= 0,
    [monthlySummary]
  );

  const handleEnvironmentChange = useCallback(
    (value: string) => {
      onFiltersChange({ ...filters, environmentId: value });
    },
    [filters, onFiltersChange]
  );

  const handleSourceChange = useCallback(
    (value: string) => {
      onFiltersChange({
        ...filters,
        sourceId: value as InsightSourceId | typeof ALL
      });
    },
    [filters, onFiltersChange]
  );

  const handleApiChange = useCallback(
    (value: string) => {
      onFiltersChange({
        ...filters,
        apiId: value as InsightApiId | typeof ALL
      });
    },
    [filters, onFiltersChange]
  );

  const handleTimeRangeChange = useCallback(
    (value: string | string[]) => {
      if (value === 'date_range') {
        return setIsDatePickerOpen(true);
      }
      const preset = value as TimeRangePreset;

      onFiltersChange({
        ...filters,
        timeRange: preset || '24h',
        customStartAt: undefined,
        customEndAt: undefined
      });
    },
    [filters, onFiltersChange]
  );
  const handleDateRangeApply = useCallback(
    (customStartAt: string, customEndAt: string) => {
      setIsDatePickerOpen(false);
      onFiltersChange({
        ...filters,
        timeRange: 'date_range',
        customStartAt,
        customEndAt
      });
    },
    [filters, onFiltersChange]
  );

  const handleExportCSV = useCallback(() => {
    if (!monthlySummary) return;
    exportMonthlySummaryCSV(monthlySummary, MONTHS);
  }, [monthlySummary]);

  const timeUnit: 'minute' | 'day' | 'hour' = useMemo(() => {
    if (filters.customStartAt && filters.customEndAt) {
      const diffDays = Math.floor(
        (Number(filters.customEndAt) - Number(filters.customStartAt)) / 86400
      );
      return diffDays >= 2 ? 'day' : 'hour';
    }
    if (['7d', '30d', 'this_month'].includes(filters.timeRange)) return 'day';
    if (filters.timeRange === '1h') return 'minute';
    return 'hour';
  }, [filters.timeRange, filters.customStartAt, filters.customEndAt]);

  const selectedProjectLabel =
    filters.projectId !== ALL
      ? projectOptions.find(o => o.value === filters.projectId)?.label
      : t('insights.all-projects');

  const selectedEnvLabel = environmentOptions.find(
    o => o.value === filters.environmentId
  )?.label;

  const selectedSourceLabel =
    filters.sourceId !== ALL
      ? sourceIdOptions.find(o => o.value === filters.sourceId)?.label
      : t('insights.all-sdks');

  const selectedApiLabel =
    filters.apiId !== ALL
      ? apiIdOptions.find(o => o.value === filters.apiId)?.label
      : t('insights.all-apis');

  return (
    <div className="p-6 flex flex-col gap-6 min-w-[950px]">
      <div className="bg-white p-4 flex flex-wrap items-end gap-4">
        <div className="flex flex-col gap-1 min-w-[180px]">
          <label className="typo-para-small text-gray-500 font-medium">
            {t('project')}
          </label>
          <DropdownMenuWithSearch
            options={projectOptions}
            label={selectedProjectLabel}
            itemSelected={filters.projectId}
            selectedOptions={
              filters.projectId !== ALL ? [filters.projectId] : []
            }
            placeholder={t('insights.all-projects')}
            onSelectOption={v => onProjectChange(String(v))}
            triggerClassName="min-w-[180px]"
            align="start"
          />
        </div>
        <div className="flex flex-col gap-1 min-w-[200px]">
          <label className="typo-para-small text-gray-500 font-medium">
            {t('environment')}
          </label>
          <DropdownMenuWithSearch
            options={environmentOptions}
            label={selectedEnvLabel}
            itemSelected={filters.environmentId}
            selectedOptions={
              filters.environmentId !== ALL ? [filters.environmentId] : []
            }
            placeholder={t('insights.all-environments')}
            onSelectOption={v => handleEnvironmentChange(String(v))}
            triggerClassName="min-w-[200px]"
            align="start"
          />
        </div>
        <div className="flex flex-col gap-1 min-w-[160px]">
          <label className="typo-para-small text-gray-500 font-medium">
            SDK
          </label>
          <DropdownMenuWithSearch
            options={sourceIdOptions}
            label={selectedSourceLabel}
            itemSelected={filters.sourceId}
            selectedOptions={filters.sourceId !== ALL ? [filters.sourceId] : []}
            placeholder={t('insights.all-sdks')}
            onSelectOption={v => handleSourceChange(String(v))}
            triggerClassName="min-w-[160px]"
            align="start"
          />
        </div>
      </div>

      <div className="shadow-card-secondary rounded-xl border border-gray-200 p-6">
        <div className="flex justify-between items-center mb-4">
          <p className="typo-head-bold-small">{t('insights.monthly-use')}</p>

          <Button
            disabled={isDisableExportCSV}
            onClick={handleExportCSV}
            className="flex items-center gap-2 typo-para-medium"
          >
            <Icon icon={IconDownload} size="sm" className="text-white" />
            {t('insights.export-csv')}
          </Button>
        </div>
        <div className="grid grid-cols-2 gap-6">
          <MonthlyBarChart
            description={
              <ChartDescription
                title={t('insights.description.mau.title')}
                notes={[
                  t('insights.description.mau.notes.first'),
                  t('insights.description.mau.notes.second'),
                  t('insights.description.mau.notes.third'),
                  t('insights.description.mau.notes.fourth')
                ]}
              />
            }
            title={t('insights.estimated-mau')}
            summary={monthlySummary}
            isLoading={monthlySummaryLoading}
            field="mau"
            months={MONTHS}
            labels={MONTH_LABELS}
            environmentNameMap={environmentNameMap}
          />
          <MonthlyBarChart
            description={
              <ChartDescription
                title={t('insights.description.monthlyRequests.title')}
                notes={[
                  t('insights.description.monthlyRequests.notes.first'),
                  t('insights.description.monthlyRequests.notes.second'),
                  t('insights.description.monthlyRequests.notes.third')
                ]}
              />
            }
            title={t('insights.requests')}
            summary={monthlySummary}
            isLoading={monthlySummaryLoading}
            field="requests"
            months={MONTHS}
            labels={MONTH_LABELS}
            environmentNameMap={environmentNameMap}
          />
        </div>
      </div>

      <div className="bg-white p-4 rounded-xl shadow-card-secondary">
        {/* Time Series Filters */}
        <div className="flex flex-wrap items-end py-4 gap-4">
          <div className="flex items-center gap-3 min-w-[160px]">
            <label className="typo-para-small text-gray-500 font-medium">
              API:
            </label>
            <DropdownMenuWithSearch
              options={apiIdOptions}
              label={selectedApiLabel}
              itemSelected={filters.apiId}
              selectedOptions={filters.apiId !== ALL ? [filters.apiId] : []}
              placeholder={t('insights.all-apis')}
              onSelectOption={v => handleApiChange(String(v))}
              triggerClassName="min-w-[160px]"
              align="start"
            />
          </div>
          <div className="flex items-center gap-x-3 min-w-[180px]">
            <Dropdown
              options={timeRangeOptions}
              value={dateRangeLabel ? 'date_range' : filters.timeRange}
              labelCustom={
                <span className="flex items-center gap-1.5">
                  <Icon icon={IconWatch} size="xs" />
                  <span className="truncate">
                    {dateRangeLabel ||
                      timeRangeOptions.find(o => o.value === filters.timeRange)
                        ?.label}
                  </span>
                </span>
              }
              onChange={v => handleTimeRangeChange(v as string)}
              className="min-w-[180px]"
              contentClassName="!max-h-[300px]"
            />
            <div className="hidden">
              <DateCustom
                onApply={handleDateRangeApply}
                onLabelChange={() => {}}
                isOpen={isDatePickerOpen}
                onClose={() => setIsDatePickerOpen(false)}
              />
            </div>
          </div>
        </div>

        <div className="grid grid-cols-2 gap-6">
          <TimeSeriesLineChart
            title={t('insights.average-latency')}
            description={t('insights.description.latency.title')}
            legendTitle={t('insights.latency')}
            timeseries={latencyData?.timeseries ?? []}
            isLoading={latencyLoading}
            timeUnit={timeUnit}
            startAt={timeRangeParams.startAt}
            endAt={timeRangeParams.endAt}
            yAxisFormatter={v => `${formatYAxis(v)}ms`}
            environmentNameMap={environmentNameMap}
          />
          <TimeSeriesLineChart
            title={t('insights.request-per-second')}
            legendTitle={t('insights.request-per-sec')}
            timeseries={requestsData?.timeseries ?? []}
            isLoading={requestsLoading}
            timeUnit={timeUnit}
            startAt={timeRangeParams.startAt}
            endAt={timeRangeParams.endAt}
            environmentNameMap={environmentNameMap}
          />
          <TimeSeriesLineChart
            title={t('insights.evaluations-second')}
            description={
              <ChartDescription
                title={t('insights.description.evaluations.title')}
                notes={[
                  <Trans
                    i18nKey="insights.description.evaluations.types.diff"
                    components={{
                      b: <b />
                    }}
                  />,
                  <Trans
                    i18nKey="insights.description.evaluations.types.none"
                    components={{
                      b: <b />
                    }}
                  />,
                  <Trans
                    i18nKey="insights.description.evaluations.types.all"
                    components={{
                      b: <b />
                    }}
                  />
                ]}
              />
            }
            legendTitle={t('insights.evaluations-per-sec')}
            timeseries={evaluationsData?.timeseries ?? []}
            isLoading={evaluationsLoading}
            timeUnit={timeUnit}
            startAt={timeRangeParams.startAt}
            endAt={timeRangeParams.endAt}
            environmentNameMap={environmentNameMap}
            labelBuilder={series => {
              const labelValues = series.labels
                ? Object.values(series.labels)
                : [];
              const parts = [...labelValues, series.sourceId];
              return parts.filter(Boolean).join(' / ');
            }}
          />
          <TimeSeriesLineChart
            title={t('insights.error-rate-title')}
            legendTitle={t('insights.error-rate-legend')}
            timeseries={errorRatesData?.timeseries ?? []}
            isLoading={errorRatesLoading}
            timeUnit={timeUnit}
            startAt={timeRangeParams.startAt}
            endAt={timeRangeParams.endAt}
            yAxisFormatter={v => `${formatYAxis(v)}%`}
            environmentNameMap={environmentNameMap}
          />
        </div>
      </div>
    </div>
  );
};

export default PageContent;
