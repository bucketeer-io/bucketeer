import {
  forwardRef,
  memo,
  Ref,
  useEffect,
  useImperativeHandle,
  useRef
} from 'react';
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
  Legend
} from 'chart.js';
import 'chartjs-adapter-luxon';
import { formatTooltipLabel } from 'utils/chart';
import { formatLongDateTime } from 'utils/date-time';
import { getVariationColor } from 'utils/style';
import { DataLabel } from './timeseries-line-chart';

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

export interface DatasetReduceType {
  label?: string;
  hidden: boolean;
}

export interface ChartToggleLegendRef {
  toggleLegend: (variationId: string) => void;
}

interface DatasetType extends ChartDataset<'line'> {
  value: string;
}

interface TimeseriesAreaLineChartProps {
  chartType: string;
  label?: string;
  dataLabels: Array<DataLabel>;
  timeseries: Array<number | string>;
  upperBoundaries: Array<Array<number>>;
  lowerBoundaries: Array<Array<number>>;
  representatives: Array<Array<number>>;
  setDataSets: (datasets: DatasetReduceType[]) => void;
}

export const TimeseriesAreaLineChart = memo(
  forwardRef(
    (
      {
        chartType,
        label,
        dataLabels,
        timeseries,
        upperBoundaries,
        lowerBoundaries,
        representatives,
        setDataSets
      }: TimeseriesAreaLineChartProps,
      ref: Ref<ChartToggleLegendRef>
    ) => {
      const labels = timeseries.map(t => new Date(Number(t) * 1000));
      const datasets: DatasetType[] = [];

      dataLabels.forEach((l, i) => {
        const color = getVariationColor(i);
        const hexColor = hexToRgba(color, 0.2);

        datasets.push({
          label: undefined,
          data: upperBoundaries[i],
          backgroundColor: hexColor,
          borderWidth: 0,
          pointRadius: 0,
          pointHoverRadius: 0,
          pointHoverBorderWidth: 0,
          pointHitRadius: 0,
          fill: '+1',
          value: l.value,
          tension: 0.2
        });
        datasets.push({
          label: undefined,
          data: lowerBoundaries[i],
          backgroundColor: hexColor,
          borderWidth: 0,
          pointRadius: 0,
          pointHoverRadius: 0,
          pointHoverBorderWidth: 0,
          pointHitRadius: 0,
          fill: '-1',
          value: l.value,
          tension: 0.2
        });
        datasets.push({
          label: l.label,
          data: representatives[i],
          borderColor: color,
          fill: false,
          value: l.value,
          tension: 0.2
        });
      });

      const chartRef = useRef<ChartJS<'line'> | null>(null);

      useImperativeHandle(ref, () => {
        return {
          toggleLegend(variationId: string) {
            toggleDataset(variationId);
          }
        };
      }, [chartRef]);

      const toggleDataset = (variationId: string) => {
        const chart = chartRef.current;
        if (chart) {
          const datasets: DatasetType[] = chart.data.datasets as DatasetType[];
          datasets.forEach((dataset, index) => {
            if (dataset?.value === variationId) {
              datasets[index].hidden = !datasets[index].hidden;
            }
          });

          chart.update();
          setDataSets(
            datasets.map(dataset => ({
              label: dataset.label,
              hidden: dataset.hidden || false
            }))
          );
        }
      };

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
              },
              label: formatTooltipLabel
            }
          }
        },
        scales: {
          x: {
            type: 'time',
            time: {
              unit: 'day',
              displayFormats: {
                day: 'MMM d'
              }
            },
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

      useEffect(() => {
        if (chartRef.current) {
          const datasets = chartRef.current?.data?.datasets;
          setDataSets(
            datasets.map(dataset => ({
              label: dataset.label,
              hidden: dataset.hidden || false
            }))
          );
        }
      }, [chartRef, chartType]);

      return (
        <div className="flex flex-1 w-full min-w-[650px] h-fit pr-5">
          <Line
            ref={chartRef}
            height={300}
            data={chartData}
            options={options}
            datasetIdKey={datasetKeyProvider().toString()}
          />
        </div>
      );
    }
  )
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
