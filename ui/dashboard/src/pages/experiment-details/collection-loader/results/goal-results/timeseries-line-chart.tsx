import { memo, useMemo } from 'react';
import { Line } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LineElement,
  LinearScale,
  TimeScale,
  TimeSeriesScale,
  PointElement,
  Tooltip,
  Filler,
  ChartData,
  ChartOptions
} from 'chart.js';
import { CHART_COLORS } from 'constants/styles';
import { formatTooltipLabel, formatXAxisLabel } from 'utils/chart';
import { formatLongDateTime } from 'utils/date-time';

ChartJS.register(
  LineElement,
  CategoryScale,
  LinearScale,
  TimeScale,
  PointElement,
  TimeSeriesScale,
  Tooltip,
  Filler
);

export type TicksType = {
  label: string;
  value: string;
};

interface TimeseriesLineChartProps {
  label?: string;
  dataLabels: Array<string>;
  timeseries: Array<string | number>;
  data: Array<Array<number | string>>;
}

const TimeseriesLineChart = memo(
  ({ label, dataLabels, timeseries, data }: TimeseriesLineChartProps) => {
    const labels = timeseries.map(t => new Date(Number(t) * 1000));

    const chartData: ChartData<'line', (string | number)[], Date> = {
      labels,
      datasets: dataLabels.map((e, i) => {
        const color = CHART_COLORS[i % CHART_COLORS.length];
        return {
          label: e,
          data: [...data[i]],
          borderColor: color,
          backgroundColor: color,
          fill: false
        };
      })
    };
    const options: ChartOptions<'line'> = useMemo(
      () => ({
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: {
            display: false
          },
          tooltip: {
            enabled: true,
            callbacks: {
              title: tooltipItems => {
                const dateString = tooltipItems[0].label;

                const date = new Date(dateString);
                if (date instanceof Date) {
                  return formatLongDateTime({
                    value: String(date.getTime() / 1000)
                  });
                }
                return tooltipItems[0].label;
              },
              label: formatTooltipLabel
            }
          },
          title: {
            display: label ? true : false,
            text: label
          }
        },
        scales: {
          x: {
            title: {
              display: false
            },
            grid: {
              display: true,
              color: '#E2E8F0',
              lineWidth: 2,
              tickWidth: 0
            },
            border: {
              dash: [5, 5],
              color: '#E2E8F0'
            },
            ticks: {
              align: 'center' as const,
              callback: function (index: string | number) {
                return formatXAxisLabel(Number(index), labels);
              },
              autoSkip: false,
              minRotation: 0,
              maxRotation: 0,
              font: {
                family: 'Sofia Pro',
                size: 14,
                weight: 400
              },
              color: '#94A3B8'
            }
          },
          y: {
            title: {
              display: false
            },
            border: {
              dash: [5, 5]
            },
            ticks: {
              font: {
                family: 'Sofia Pro',
                size: 14,
                weight: 400
              },
              color: '#94A3B8'
            },
            grid: {
              display: true,
              color: '#E2E8F0',
              lineWidth: 2,
              tickWidth: 0
            }
          }
        }
      }),
      []
    );

    return (
      <div className="flex flex-1 w-full min-w-[650px] h-fit pr-5">
        <Line
          style={{
            width: '100%'
          }}
          data={chartData}
          options={options}
        />
      </div>
    );
  }
);

export default TimeseriesLineChart;
