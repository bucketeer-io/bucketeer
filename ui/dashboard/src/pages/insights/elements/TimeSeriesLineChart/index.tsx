import { useMemo } from 'react';
import { Line } from 'react-chartjs-2';
import { InsightsTimeSeries } from '@types';
import Spinner from 'components/spinner';
import { LegendTable } from '../LegendTable';
import { ChartCard, getColor } from '../chart-utils';

interface TimeSeriesLineChartProps {
  title: string;
  legendTitle?: string;
  timeseries: InsightsTimeSeries[];
  isLoading: boolean;
  timeUnit?: 'minute' | 'hour' | 'day';
  startAt: string;
  endAt: string;
  yAxisFormatter?: (value: number) => string;
  environmentNameMap?: Record<string, string>;
}

const TimeSeriesLineChart = ({
  title,
  legendTitle,
  timeseries,
  isLoading,
  timeUnit,
  startAt,
  endAt,
  yAxisFormatter,
  environmentNameMap
}: TimeSeriesLineChartProps) => {
  const datasets = useMemo(
    () =>
      timeseries.map((series, i) => {
        const labelValues = series.labels ? Object.values(series.labels) : [];
        let label: string;
        if (labelValues.length > 0) {
          label = labelValues.join(' / ');
        } else {
          const labelParts = [
            environmentNameMap?.[series.environmentId] ?? series.environmentId,
            series.sourceId
          ];
          if (series.apiId && series.apiId !== 'UNKNOWN_API')
            labelParts.push(series.apiId);
          label = labelParts.join(' / ');
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
      }),
    [timeseries, environmentNameMap]
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
                    min: Number(startAt) * 1000,
                    max: Number(endAt) * 1000,
                    time: {
                      unit: timeUnit,
                      displayFormats: {
                        minute: 'HH:mm',
                        hour: 'HH:mm',
                        day: 'MMM d'
                      }
                    },
                    grid: { display: false },
                    border: { display: false },
                    ticks: {
                      align: 'center',
                      source: 'auto',
                      color: '#94A3B8',
                      font: { family: 'Sofia Pro', size: 14, weight: 400 }
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
