import { ReactNode, useMemo } from 'react';
import { Line } from 'react-chartjs-2';
import { InsightsTimeSeries } from '@types';
import Spinner from 'components/spinner';
import { LegendTable } from '../LegendTable';
import { ChartCard, formatYAxis, getColor } from '../chart-utils';

interface TimeSeriesLineChartProps {
  title: string;
  description?: string | ReactNode;
  legendTitle?: string;
  timeseries: InsightsTimeSeries[];
  isLoading: boolean;
  timeUnit?: 'minute' | 'hour' | 'day';
  startAt: string;
  endAt: string;
  yAxisFormatter?: (value: number) => string;
  environmentNameMap?: Record<string, string>;
  labelBuilder?: (series: InsightsTimeSeries) => string;
}

const TimeSeriesLineChart = ({
  title,
  description,
  legendTitle,
  timeseries,
  isLoading,
  timeUnit,
  startAt,
  endAt,
  yAxisFormatter = formatYAxis,
  environmentNameMap,
  labelBuilder
}: TimeSeriesLineChartProps) => {
  const datasets = useMemo(() => {
    const uniqueEnvIds = new Set(timeseries.map(s => s.environmentId));
    const singleEnv = uniqueEnvIds.size <= 1;

    return timeseries.map((series, i) => {
      let label: string;
      if (labelBuilder) {
        label = labelBuilder(series);
      } else {
        const labelValues = series.labels ? Object.values(series.labels) : [];
        if (labelValues.length > 0) {
          label = labelValues.join(' / ');
        } else {
          const labelParts: string[] = [];
          if (!singleEnv) {
            labelParts.push(
              environmentNameMap?.[series.environmentId] ?? series.environmentId
            );
          }
          labelParts.push(series.sourceId);
          if (series.apiId && series.apiId !== 'UNKNOWN_API')
            labelParts.push(series.apiId);
          label = labelParts.filter(Boolean).join(' / ');
        }
      }
      return {
        label,
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
    });
  }, [timeseries, environmentNameMap, labelBuilder]);

  const legendData = useMemo(
    () =>
      datasets.map(ds => ({ label: ds.label, data: ds.data.map(d => d.y) })),
    [datasets]
  );

  const options = useMemo(
    () => ({
      responsive: true,
      maintainAspectRatio: false,
      plugins: { legend: { display: false } },
      scales: {
        x: {
          type: 'time' as const,
          min: Number(startAt) * 1000,
          max: Number(endAt) * 1000,
          time: {
            unit: timeUnit,
            displayFormats: { minute: 'HH:mm', hour: 'HH:mm', day: 'MMM d' }
          },
          grid: { display: false },
          border: { display: false },
          ticks: {
            align: 'center' as const,
            source: 'auto' as const,
            color: '#94A3B8',
            font: { family: 'Sofia Pro', size: 14, weight: 400 }
          }
        },
        y: {
          min: 0,
          grid: { color: '#E2E8F0', drawTicks: false },
          title: {
            display: true,
            text: legendTitle,
            color: '#94A3B8',
            font: { size: 12, weight: 'bold' as const }
          },
          ticks: {
            color: '#94A3B8',
            font: { size: 12 },
            callback: (value: number | string) => yAxisFormatter(Number(value))
          }
        }
      }
    }),
    [startAt, endAt, timeUnit, legendTitle, yAxisFormatter]
  );

  return (
    <ChartCard title={title} description={description}>
      {isLoading ? (
        <div className="h-[250px] flex items-center justify-center">
          <Spinner />
        </div>
      ) : (
        <div className="flex-1">
          <div className="h-[250px]">
            <Line data={{ datasets }} options={options} />
          </div>
          <LegendTable datasets={legendData} formatter={yAxisFormatter} />
        </div>
      )}
    </ChartCard>
  );
};

export default TimeSeriesLineChart;
