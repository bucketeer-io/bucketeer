import { ReactNode, useMemo } from 'react';
import { Bar } from 'react-chartjs-2';
import { useTranslation } from 'i18n';
import { InsightsMonthlySummaryResponse } from '@types';
import Spinner from 'components/spinner';
import { ChartCard, formatLargeNumber, getColor } from '../chart-utils';

interface MonthlyBarChartProps {
  title: string;
  description?: ReactNode;
  summary?: InsightsMonthlySummaryResponse;
  isLoading: boolean;
  field: 'mau' | 'requests';
  months: string[];
  labels: string[];
  environmentNameMap: Record<string, string>;
}

const BAR_OPTIONS = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: { legend: { display: false } },
  scales: {
    x: {
      stacked: true,
      grid: { display: false },
      ticks: { color: '#94A3B8', font: { size: 12 } }
    },
    y: {
      stacked: true,
      grid: { color: '#E2E8F0' },
      ticks: {
        color: '#94A3B8',
        font: { size: 12 },
        callback: (value: number | string) => formatLargeNumber(Number(value))
      }
    }
  }
} as const;

const MonthlyBarChart = ({
  title,
  description,
  summary,
  isLoading,
  field,
  months,
  labels,
  environmentNameMap
}: MonthlyBarChartProps) => {
  const { t } = useTranslation(['common']);

  const { datasets, isStacked } = useMemo(() => {
    if (!summary?.series?.length) return { datasets: [], isStacked: false };

    // Group series by environmentId
    const envIds = Array.from(
      new Set(summary.series.map(s => s.environmentId))
    );
    const stacked = envIds.length > 1;

    if (stacked) {
      // Only include environments that have at least one non-zero data point
      const activeEnvIds = envIds.filter(envId => {
        const envSeries = summary.series.filter(s => s.environmentId === envId);
        return months.some(m =>
          envSeries.some(s => {
            const dp = s.data.find(d => d.yearmonth === m);
            return dp ? Number(dp[field]) > 0 : false;
          })
        );
      });

      const envDatasets = activeEnvIds.map((envId, i) => {
        const envSeries = summary.series.filter(s => s.environmentId === envId);
        const data = months.map(m =>
          envSeries.reduce((sum, series) => {
            const dp = series.data.find(d => d.yearmonth === m);
            return sum + (dp ? Number(dp[field]) : 0);
          }, 0)
        );
        const color = getColor(i);
        const envName =
          environmentNameMap[envId] || envSeries[0]?.environmentName || envId;
        return {
          label: envName,
          data,
          backgroundColor: color,
          borderRadius: 0,
          stack: 'monthly'
        };
      });
      // Round top bar corners
      if (envDatasets.length > 0) {
        envDatasets[envDatasets.length - 1] = {
          ...envDatasets[envDatasets.length - 1],
          borderRadius: 4
        };
      }
      return { datasets: envDatasets, isStacked: true };
    }

    // Single env — original behaviour
    const totals = months.map(m =>
      summary.series.reduce((sum, series) => {
        const dp = series.data.find(d => d.yearmonth === m);
        return sum + (dp ? Number(dp[field]) : 0);
      }, 0)
    );
    const baseColor = getColor(0);
    const backgroundColor = totals.map((_, i) =>
      i === totals.length - 1 ? baseColor : `${baseColor}66`
    );
    return {
      datasets: [
        {
          label: field === 'mau' ? t('insights.mau') : t('insights.requests'),
          data: totals,
          backgroundColor,
          borderRadius: 4,
          stack: 'monthly'
        }
      ],
      isStacked: false
    };
  }, [summary, field, months, t, environmentNameMap]);

  // For the summary numbers, use the total across all envs for the last two months
  const totals = useMemo(() => {
    if (!datasets.length) return [];
    return months.map((_, mi) =>
      datasets.reduce((sum, ds) => sum + (ds.data[mi] ?? 0), 0)
    );
  }, [datasets, months]);

  const currentMonth =
    totals.length > 0 ? totals[totals.length - 1] : undefined;
  const lastMonth = totals.length > 1 ? totals[totals.length - 2] : undefined;

  return (
    <ChartCard
      title={title}
      description={description}
      currentMonth={currentMonth}
      lastMonth={lastMonth}
      legendDatasets={isStacked ? datasets : undefined}
    >
      {isLoading ? (
        <div className="h-[300px] flex items-center justify-center">
          <Spinner />
        </div>
      ) : (
        <div className="h-[300px]">
          <Bar data={{ labels, datasets }} options={BAR_OPTIONS} />
        </div>
      )}
    </ChartCard>
  );
};

export default MonthlyBarChart;
