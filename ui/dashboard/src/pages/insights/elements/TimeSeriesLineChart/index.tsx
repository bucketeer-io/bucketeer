import { useMemo } from 'react';
import { Line } from 'react-chartjs-2';
import { useTranslation } from 'i18n';
import { InsightsTimeSeries } from '@types';
import Spinner from 'components/spinner';
import EmptyState from 'elements/empty-state';
import { ChartCard, LegendTable, getColor } from '../chart-utils';

interface TimeSeriesLineChartProps {
  title: string;
  legendTitle?: string;
  timeseries: InsightsTimeSeries[];
  isLoading: boolean;
  timeUnit?: 'hour' | 'day';
  yAxisFormatter?: (value: number) => string;
  environmentNameMap?: Record<string, string>;
}

const TimeSeriesLineChart = ({
  title,
  legendTitle,
  timeseries,
  isLoading,
  timeUnit,
  yAxisFormatter,
  environmentNameMap
}: TimeSeriesLineChartProps) => {
  const { t } = useTranslation(['common']);
  const datasets = useMemo(
    () =>
      timeseries.map((series, i) => {
        const labelParts = [
          environmentNameMap?.[series.environmentId] ?? series.environmentId,
          series.sourceId
        ];
        if (series.apiId && series.apiId !== 'UNKNOWN_API')
          labelParts.push(series.apiId);
        if (series.labels) {
          Object.entries(series.labels).forEach(([k, v]) =>
            labelParts.push(`${k}:${v}`)
          );
        }
        return {
          label: labelParts.join(' / '),
          data: series.data.map(d => ({
            x: Number(d.timestamp) * 1000,
            y: d.value
          })),
          borderColor: getColor(i),
          backgroundColor: getColor(i),
          fill: false,
          tension: 0.2,
          pointRadius: 0
        };
      }),
    [timeseries]
  );

  const legendData = useMemo(
    () =>
      datasets.map(ds => ({ label: ds.label, data: ds.data.map(d => d.y) })),
    [datasets]
  );

  return (
    <ChartCard title={title}>
      {isLoading ? (
        <div className="h-[300px] flex items-center justify-center">
          <Spinner />
        </div>
      ) : !datasets.length ? (
        <div className="h-[300px] flex items-center justify-center">
          <EmptyState.Root variant="no-data" size="sm">
            <EmptyState.Illustration />
            <EmptyState.Body>
              <EmptyState.Title>{t('no-data')}</EmptyState.Title>
            </EmptyState.Body>
          </EmptyState.Root>
        </div>
      ) : (
        <div className="flex-1">
          <div className="h-[300px]">
            <Line
              data={{ datasets }}
              options={{
                responsive: true,
                maintainAspectRatio: false,
                plugins: { legend: { display: false } },
                scales: {
                  x: {
                    type: 'time',
                    time: {
                      unit: timeUnit,
                      displayFormats: { hour: 'HH:mm', day: 'MMM d' }
                    },
                    grid: { display: false },
                    border: { display: false },
                    ticks: {
                      color: '#94A3B8',
                      font: { size: 12 },
                      maxTicksLimit: timeUnit === 'day' ? 10 : undefined
                    }
                  },
                  y: {
                    grid: { color: '#E2E8F0', drawTicks: false },
                    title: {
                      display: true,
                      text: legendTitle,
                      color: '#94A3B8',
                      font: { size: 12, weight: 'bold' }
                    },
                    ticks: {
                      color: '#94A3B8',
                      font: { size: 12 },
                      ...(yAxisFormatter && {
                        callback: (value: number | string) =>
                          yAxisFormatter(Number(value))
                      })
                    }
                  }
                }
              }}
            />
          </div>
          <LegendTable datasets={legendData} formatter={yAxisFormatter} />
        </div>
      )}
    </ChartCard>
  );
};

export default TimeSeriesLineChart;
