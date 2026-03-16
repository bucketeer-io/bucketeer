import { useCallback, useMemo } from 'react';
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
import { useTranslation } from 'i18n';
import { DateTime } from 'luxon';
import {
  Environment,
  InsightApiId,
  InsightsMonthlySummaryResponse,
  InsightSourceId,
  Project
} from '@types';
import { IconInfo, IconUpload } from '@icons';
import Dropdown, { DropdownOption } from 'components/dropdown';
import Icon from 'components/icon';
import DropdownMenuWithSearch from 'elements/dropdown-with-search';
import MonthlyBarChart from './elements/MonthlyBarChart';
import TimeSeriesLineChart from './elements/TimeSeriesLineChart';
import { InsightsFilters, TimeRangePreset } from './page-loader';

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
  DateTime.fromFormat(m, 'yyyyMM').toFormat('MMM yy')
);

interface PageContentProps {
  projects: Project[];
  environments: Environment[];
  monthlySummary?: InsightsMonthlySummaryResponse;
  monthlySummaryLoading: boolean;
  timeRangeParams: InsightsTimeSeriesFetcherParams;
  filters: InsightsFilters;
  onFiltersChange: (filters: InsightsFilters) => void;
}

const SDK_OPTION_VALUES = [
  { label: 'Android', value: 'ANDROID' },
  { label: 'iOS', value: 'IOS' },
  { label: 'Web', value: 'WEB' },
  { label: 'Go Server', value: 'GO_SERVER' },
  { label: 'Node Server', value: 'NODE_SERVER' },
  { label: 'JavaScript', value: 'JAVASCRIPT' },
  { label: 'Flutter', value: 'FLUTTER' },
  { label: 'React', value: 'REACT' },
  { label: 'React Native', value: 'REACT_NATIVE' },
  { label: 'OpenFeature Kotlin', value: 'OPEN_FEATURE_KOTLIN' },
  { label: 'OpenFeature Swift', value: 'OPEN_FEATURE_SWIFT' },
  { label: 'OpenFeature JavaScript', value: 'OPEN_FEATURE_JAVASCRIPT' },
  { label: 'OpenFeature Go', value: 'OPEN_FEATURE_GO' },
  { label: 'OpenFeature Node', value: 'OPEN_FEATURE_NODE' },
  { label: 'OpenFeature React', value: 'OPEN_FEATURE_REACT' },
  { label: 'OpenFeature React Native', value: 'OPEN_FEATURE_REACT_NATIVE' }
];

const API_OPTION_VALUES = [
  { label: 'GetEvaluation', value: 'GET_EVALUATION' },
  { label: 'GetEvaluations', value: 'GET_EVALUATIONS' },
  { label: 'RegisterEvents', value: 'REGISTER_EVENTS' },
  { label: 'GetFeatureFlags', value: 'GET_FEATURE_FLAGS' },
  { label: 'GetSegmentUsers', value: 'GET_SEGMENT_USERS' },
  { label: 'SdkGetVariation', value: 'SDK_GET_VARIATION' }
];

const PageContent = ({
  projects,
  environments,
  monthlySummary,
  monthlySummaryLoading,
  timeRangeParams,
  filters,
  onFiltersChange
}: PageContentProps) => {
  const { t } = useTranslation(['common']);

  const sourceIdOptions: DropdownOption[] = useMemo(
    () => [{ label: t('insights.all-sdks'), value: '' }, ...SDK_OPTION_VALUES],
    [t]
  );

  const apiIdOptions: DropdownOption[] = useMemo(
    () => [{ label: t('insights.all-apis'), value: '' }, ...API_OPTION_VALUES],
    [t]
  );

  const timeRangeOptions: DropdownOption[] = useMemo(
    () => [
      { label: t('insights.last-1-hour'), value: '1h' },
      { label: t('insights.last-6-hours'), value: '6h' },
      { label: t('insights.last-24-hours'), value: '24h' },
      { label: t('insights.last-7-days'), value: '7d' },
      { label: t('insights.last-30-days'), value: '30d' },
      { label: t('insights.this-month'), value: 'this_month' }
    ],
    [t]
  );

  const projectOptions: DropdownOption[] = useMemo(
    () => [
      { label: t('insights.all-projects'), value: '' },
      ...projects.map(p => ({ label: p.name, value: p.id }))
    ],
    [projects, t]
  );

  const environmentOptions: DropdownOption[] = useMemo(() => {
    const filtered = filters.projectId
      ? environments.filter(e => e.projectId === filters.projectId)
      : environments;
    return [
      { label: t('insights.all-environments'), value: '' },
      ...filtered.map(e => ({ label: e.name, value: e.id }))
    ];
  }, [environments, filters.projectId, t]);

  const environmentNameMap = useMemo(
    () =>
      Object.fromEntries(environments.map(e => [e.id, e.name])) as Record<
        string,
        string
      >,
    [environments]
  );

  const { data: latencyData, isLoading: latencyLoading } =
    useQueryInsightsLatency({ params: timeRangeParams });

  const { data: requestsData, isLoading: requestsLoading } =
    useQueryInsightsRequests({ params: timeRangeParams });

  const { data: evaluationsData, isLoading: evaluationsLoading } =
    useQueryInsightsEvaluations({ params: timeRangeParams });

  const { data: errorRatesData, isLoading: errorRatesLoading } =
    useQueryInsightsErrorRates({ params: timeRangeParams });

  const handleProjectChange = useCallback(
    (value: string) => {
      onFiltersChange({ ...filters, projectId: value, environmentId: '' });
    },
    [filters, onFiltersChange]
  );

  const handleEnvironmentChange = useCallback(
    (value: string) => {
      onFiltersChange({ ...filters, environmentId: value });
    },
    [filters, onFiltersChange]
  );

  const handleSourceChange = useCallback(
    (value: string) => {
      onFiltersChange({ ...filters, sourceId: value as InsightSourceId | '' });
    },
    [filters, onFiltersChange]
  );

  const handleApiChange = useCallback(
    (value: string) => {
      onFiltersChange({ ...filters, apiId: value as InsightApiId | '' });
    },
    [filters, onFiltersChange]
  );

  const handleTimeRangeChange = useCallback(
    (value: string | string[]) => {
      const preset = (
        Array.isArray(value) ? value[0] : value
      ) as TimeRangePreset;
      onFiltersChange({ ...filters, timeRange: preset || '24h' });
    },
    [filters, onFiltersChange]
  );

  const handleExportCSV = useCallback(() => {
    if (!monthlySummary?.series?.length) return;

    const rows: string[] = [
      ['Project', 'Environment', 'SDK', 'Month', 'MAU', 'Requests'].join(',')
    ];

    for (const series of monthlySummary.series) {
      const dataMap = new Map(series.data.map(d => [d.yearmonth, d]));
      for (const m of MONTHS) {
        const dp = dataMap.get(m);
        rows.push(
          [
            series.projectName,
            series.environmentName,
            series.sourceId,
            DateTime.fromFormat(m, 'yyyyMM').toFormat('yyyy-MM'),
            dp?.mau ?? '0',
            dp?.requests ?? '0'
          ].join(',')
        );
      }
    }

    const blob = new Blob([rows.join('\n')], { type: 'text/csv' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `insights-monthly-${DateTime.now().toFormat('yyyyMMdd')}.csv`;
    a.click();
    URL.revokeObjectURL(url);
  }, [monthlySummary]);

  const timeUnit: 'day' | 'hour' = ['7d', '30d', 'this_month'].includes(
    filters.timeRange
  )
    ? 'day'
    : 'hour';

  const selectedProjectLabel = filters.projectId
    ? projectOptions.find(o => o.value === filters.projectId)?.label
    : undefined;
  const selectedEnvLabel = filters.environmentId
    ? environmentOptions.find(o => o.value === filters.environmentId)?.label
    : undefined;
  const selectedSourceLabel = filters.sourceId
    ? sourceIdOptions.find(o => o.value === filters.sourceId)?.label
    : undefined;
  const selectedApiLabel = filters.apiId
    ? apiIdOptions.find(o => o.value === filters.apiId)?.label
    : undefined;

  return (
    <div className="p-6 flex flex-col gap-6 min-w-[1200px]">
      <div className="bg-white p-4 flex flex-wrap items-end gap-4">
        <div className="flex flex-col gap-1 min-w-[180px]">
          <label className="typo-para-small text-gray-500 font-medium">
            {t('project')}
          </label>
          <DropdownMenuWithSearch
            options={projectOptions}
            label={selectedProjectLabel}
            itemSelected={filters.projectId}
            selectedOptions={filters.projectId ? [filters.projectId] : []}
            placeholder={t('insights.all-projects')}
            onSelectOption={v => handleProjectChange(String(v))}
            triggerClassName="min-w-[180px]"
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
              filters.environmentId ? [filters.environmentId] : []
            }
            placeholder={t('insights.all-environments')}
            onSelectOption={v => handleEnvironmentChange(String(v))}
            triggerClassName="min-w-[200px]"
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
            selectedOptions={filters.sourceId ? [filters.sourceId] : []}
            placeholder={t('insights.all-sdks')}
            onSelectOption={v => handleSourceChange(String(v))}
            triggerClassName="min-w-[160px]"
          />
        </div>
      </div>

      <div className="shadow-card-secondary rounded-xl border border-gray-200 p-6">
        <div className="flex justify-between items-center mb-4">
          <div className="flex items-center gap-x-1">
            <p className="typo-head-bold-small">{t('insights.monthly-use')}</p>
            <Icon icon={IconInfo} size="xs" />
          </div>
          <button
            type="button"
            onClick={handleExportCSV}
            className="flex items-center gap-2 typo-para-medium text-gray-700 transition-colors"
          >
            <Icon icon={IconUpload} size="sm" color="gray-600" />
            {t('insights.export-csv')}
          </button>
        </div>
        <div className="grid grid-cols-2 gap-6">
          <MonthlyBarChart
            title={t('insights.estimated-mau')}
            summary={monthlySummary}
            isLoading={monthlySummaryLoading}
            field="mau"
            months={MONTHS}
            labels={MONTH_LABELS}
          />
          <MonthlyBarChart
            title={t('insights.requests')}
            summary={monthlySummary}
            isLoading={monthlySummaryLoading}
            field="requests"
            months={MONTHS}
            labels={MONTH_LABELS}
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
              selectedOptions={filters.apiId ? [filters.apiId] : []}
              placeholder={t('insights.all-apis')}
              onSelectOption={v => handleApiChange(String(v))}
              triggerClassName="min-w-[160px]"
            />
          </div>
          <div className="flex items-center gap-x-3 min-w-[180px]">
            <Dropdown
              options={timeRangeOptions}
              value={filters.timeRange}
              onChange={v => handleTimeRangeChange(v as string)}
              className="min-w-[180px]"
            />
          </div>
        </div>

        <div className="grid grid-cols-2 gap-6">
          <TimeSeriesLineChart
            title={t('insights.average-latency')}
            legendTitle={t('insights.latency')}
            timeseries={latencyData?.timeseries ?? []}
            isLoading={latencyLoading}
            timeUnit={timeUnit}
            yAxisFormatter={v => `${v.toFixed(1)}ms`}
            environmentNameMap={environmentNameMap}
          />
          <TimeSeriesLineChart
            title={t('insights.request-per-minute')}
            legendTitle={t('insights.request-per-min')}
            timeseries={requestsData?.timeseries ?? []}
            isLoading={requestsLoading}
            timeUnit={timeUnit}
            environmentNameMap={environmentNameMap}
          />
          <TimeSeriesLineChart
            title={t('insights.evaluations-second')}
            legendTitle={t('insights.evaluations-per-sec')}
            timeseries={evaluationsData?.timeseries ?? []}
            isLoading={evaluationsLoading}
            timeUnit={timeUnit}
            yAxisFormatter={v => v.toFixed(1)}
            environmentNameMap={environmentNameMap}
          />
          <TimeSeriesLineChart
            title={t('insights.error-rate-title')}
            legendTitle={t('insights.error-rate-legend')}
            timeseries={errorRatesData?.timeseries ?? []}
            isLoading={errorRatesLoading}
            timeUnit={timeUnit}
            yAxisFormatter={v => `${v.toFixed(2)}%`}
            environmentNameMap={environmentNameMap}
          />
        </div>
      </div>
    </div>
  );
};

export default PageContent;
