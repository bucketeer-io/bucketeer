import { memo } from 'react';
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
  ChartOptions,
  Point,
  ChartDataset
} from 'chart.js';
import { CHART_COLORS } from 'constants/styles';
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

interface TimeseriesAreaLineChartProps {
  label?: string;
  dataLabels: Array<string>;
  timeseries: Array<number | string>;
  upperBoundaries: Array<Array<number>>;
  lowerBoundaries: Array<Array<number>>;
  representatives: Array<Array<number>>;
}

export const TimeseriesAreaLineChart = memo(
  ({
    label,
    dataLabels,
    timeseries,
    upperBoundaries,
    lowerBoundaries,
    representatives
  }: TimeseriesAreaLineChartProps) => {
    const labels = timeseries.map(t => new Date(Number(t) * 1000));

    const formatLabel = (index: number) => {
      const date = labels[index];
      if (date) {
        return formatLongDateTime({
          value: String(date.getTime() / 1000),
          overrideOptions: {
            day: '2-digit',
            month: '2-digit',
            year: undefined
          },
          locale: 'en-GB'
        });
      }
      return '';
    };

    const datasets: ChartDataset<'line'>[] = [];

    dataLabels.forEach((l, i) => {
      const color = CHART_COLORS[i % CHART_COLORS.length];
      const hexColor = hexToRgba(color, 0.2);
      datasets.push({
        label: undefined,
        data: upperBoundaries[i],
        borderColor: color,
        backgroundColor: hexColor,
        fill: '+1'
      });
      datasets.push({
        label: undefined,
        data: lowerBoundaries[i],
        backgroundColor: hexColor,
        fill: '-1'
      });
      datasets.push({
        label: l,
        data: representatives[i],
        borderColor: color,
        backgroundColor: hexColor,
        fill: false
      });
    });
    const chartData: ChartData<'line', (number | Point | null)[], Date> = {
      labels,
      datasets: datasets
    };
    const options: ChartOptions<'line'> = {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          display: false
        },
        title: {
          display: label ? true : false,
          text: label
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
            }
          }
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
              return formatLabel ? formatLabel(Number(index)) : '';
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
    };

    return (
      <div className="flex flex-1 w-full min-w-[650px] h-fit pr-5">
        <Line
          style={{
            width: '100%'
          }}
          data={chartData}
          options={options}
          datasetIdKey={datasetKeyProvider().toString()}
        />
      </div>
    );
  }
);

const hexToRgba = (hex: string, alpha: number) => {
  const r = parseInt(hex.slice(1, 3), 16),
    g = parseInt(hex.slice(3, 5), 16),
    b = parseInt(hex.slice(5, 7), 16);
  return 'rgba(' + r + ', ' + g + ', ' + b + ', ' + alpha + ')';
};

const datasetKeyProvider = () => {
  return Math.random();
};
