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
import { formatTooltipLabel, formatXAxisLabel } from 'utils/chart';
import { formatLongDateTime } from 'utils/date-time';
import { getVariationColor } from 'utils/style';

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
  toggleLegend: (label: string) => void;
}

interface TimeseriesAreaLineChartProps {
  chartType: string;
  label?: string;
  dataLabels: Array<string>;
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

      const datasets: ChartDataset<'line'>[] = [];

      dataLabels.forEach((l, i) => {
        const color = getVariationColor(i);
        const hexColor = hexToRgba(color, 0.2);
        datasets.push({
          label: undefined,
          data: upperBoundaries[i],
          borderWidth: 0,
          backgroundColor: hexColor,
          pointRadius: 0,
          fill: '+1'
        });
        datasets.push({
          label: undefined,
          data: lowerBoundaries[i],
          borderWidth: 0,
          backgroundColor: hexColor,
          pointRadius: 0,
          fill: '-1'
        });
        datasets.push({
          label: l,
          data: representatives[i],
          borderColor: color,
          pointRadius: 0,
          fill: false
        });
      });

      const chartRef = useRef<ChartJS<'line'> | null>(null);

      useImperativeHandle(ref, () => {
        return {
          toggleLegend(label: string) {
            toggleDataset(label);
          }
        };
      }, [chartRef]);

      const toggleDataset = (label: string) => {
        const chart = chartRef.current;
        if (chart) {
          const datasets = chart.data.datasets;
          const toggleIndex = datasets.findIndex(
            dataset => dataset?.label === label
          );

          datasets[toggleIndex].hidden = !datasets[toggleIndex].hidden;
          chart.update();
          if (setDataSets)
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
