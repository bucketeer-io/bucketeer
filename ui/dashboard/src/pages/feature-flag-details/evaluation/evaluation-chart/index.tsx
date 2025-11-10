import { forwardRef, Ref, useImperativeHandle, useRef } from 'react';
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
  ChartDataset,
  Legend,
  TimeUnit
} from 'chart.js';
import 'chartjs-adapter-luxon';
import { COLORS } from 'constants/styles';
import { formatTooltipLabel } from 'utils/chart';
import { formatLongDateTime } from 'utils/date-time';
import { getVariationColor } from 'utils/style';
import {
  ChartToggleLegendRef,
  DatasetReduceType
} from 'pages/experiment-details/collection-loader/results/goal-results/timeseries-area-line-chart';
import { Option } from 'components/creatable-select';

interface TimeseriesStackedLineChartProps {
  variationValues: Option[];
  timeseries: Array<string>;
  data: Array<Array<number>>;
  unit: string;
  setDataSets: (datasets: DatasetReduceType[]) => void;
}

interface DatasetType extends ChartDataset<'line'> {
  value: string;
}

ChartJS.register(
  LineElement,
  CategoryScale,
  LinearScale,
  TimeScale,
  PointElement,
  TimeSeriesScale,
  Legend,
  Tooltip,
  Filler
);

export const EvaluationChart = forwardRef(
  (
    {
      variationValues,
      timeseries,
      data,
      unit = 'day',
      setDataSets
    }: TimeseriesStackedLineChartProps,
    ref: Ref<ChartToggleLegendRef>
  ) => {
    const chartRef = useRef<ChartJS<'line'> | null>(null);

    const chartData: ChartData<'line', (number | Point | null)[], Date> = {
      labels: timeseries.map(t => new Date(Number(t) * 1000)),
      datasets: variationValues.map((e, i) => {
        const color = getVariationColor(i % COLORS.length);
        return {
          label:
            e.label.length > 40 ? `${e.label.substring(0, 40)}...` : e.label,
          data: data[i],
          backgroundColor: color,
          borderColor: color,
          fill: false,
          tension: 0.2,
          value: e.value
        };
      })
    };

    // Format large numbers: 1000 → "1K", 1000000 → "1M", 1000000000 → "1B"
    const formatNumber = (value: number): string => {
      if (value >= 1_000_000_000) {
        return `${(value / 1_000_000_000).toFixed(value % 1_000_000_000 === 0 ? 0 : 1)}B`;
      }
      if (value >= 1_000_000) {
        return `${(value / 1_000_000).toFixed(value % 1_000_000 === 0 ? 0 : 1)}M`;
      }
      if (value >= 1_000) {
        return `${(value / 1_000).toFixed(value % 1_000 === 0 ? 0 : 1)}K`;
      }
      return value.toLocaleString();
    };

    const options: ChartOptions<'line'> = {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        title: {
          display: false
        },
        legend: { display: false },
        tooltip: {
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
        }
      },
      scales: {
        x: {
          type: 'time',
          time: {
            unit: unit as TimeUnit,
            displayFormats: { hour: 'HH:mm' }
          },
          ticks: {
            align: 'center' as const,
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
          display: true,
          beginAtZero: true,
          ticks: {
            font: {
              family: 'Sofia Pro',
              size: 14,
              weight: 400
            },
            color: '#94A3B8',
            callback: value => {
              return formatNumber(Number(value));
            }
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

    useImperativeHandle(ref, () => {
      return {
        toggleLegend(label: string) {
          toggleDataset(label);
        }
      };
    }, [chartRef]);

    const toggleDataset = (value: string) => {
      const chart = chartRef.current;
      if (chart) {
        const datasets: DatasetType[] = chart.data.datasets as DatasetType[];

        datasets.forEach((dataset, index) => {
          if (dataset?.value === value)
            datasets[index].hidden = !datasets[index].hidden;
        });

        chart.update();
        setDataSets(
          datasets.map(dataset => ({
            label: dataset.value,
            hidden: dataset.hidden || false
          }))
        );
      }
    };

    return (
      <Line
        ref={chartRef}
        style={{
          height: 'auto',
          maxHeight: '600px',
          width: '100%',
          minWidth: 650
        }}
        data={chartData}
        options={options}
      />
    );
  }
);
