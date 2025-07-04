import {
  forwardRef,
  memo,
  Ref,
  useEffect,
  useImperativeHandle,
  useMemo,
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
  ChartDataset
} from 'chart.js';
import 'chartjs-adapter-luxon';
import { FeatureVariationType } from '@types';
import { formatTooltipLabel } from 'utils/chart';
import { formatLongDateTime } from 'utils/date-time';
import { getVariationColor } from 'utils/style';
import {
  ChartToggleLegendRef,
  DatasetReduceType
} from './timeseries-area-line-chart';

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

export interface DataLabel {
  label: string;
  value: string;
  variationType?: FeatureVariationType;
}

interface DatasetType extends ChartDataset<'line'> {
  value: string;
}

interface TimeseriesLineChartProps {
  label?: string;
  dataLabels: Array<DataLabel>;
  timeseries: Array<string | number>;
  data: Array<Array<number | string>>;
  setDataSets: (datasets: DatasetReduceType[]) => void;
}

const TimeseriesLineChart = memo(
  forwardRef(
    (
      {
        label,
        dataLabels,
        timeseries,
        data,
        setDataSets
      }: TimeseriesLineChartProps,
      ref: Ref<ChartToggleLegendRef>
    ) => {
      const labels = timeseries?.map(t => new Date(Number(t) * 1000));

      const chartRef = useRef<ChartJS<
        'line',
        (string | number)[],
        Date
      > | null>(null);

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
          const datasets: DatasetType[] = chart?.data
            ?.datasets as DatasetType[];
          const toggleIndex = datasets?.findIndex(
            dataset => dataset?.value === variationId
          );
          if (toggleIndex === -1) return;
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

      const chartData: ChartData<'line', (string | number)[], Date> = {
        labels,
        datasets: dataLabels?.map((e, i) => {
          const color = getVariationColor(i);
          return {
            label: e.label,
            data: [...(data[i] || [])],
            borderColor: color,
            backgroundColor: color,
            fill: false,
            value: e.value,
            tension: 0.2
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
                  const dateString = tooltipItems[0]?.label;
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
        }),
        []
      );

      useEffect(() => {
        if (chartRef.current) {
          const datasets = chartRef.current?.data?.datasets;
          setDataSets(
            datasets?.map(dataset => ({
              label: dataset.label,
              hidden: dataset.hidden || false
            }))
          );
        }
      }, []);

      return (
        <div className="flex flex-1 w-full min-w-[650px] h-fit pr-5">
          <Line
            ref={chartRef}
            height={300}
            data={chartData}
            options={options}
          />
        </div>
      );
    }
  )
);

export default TimeseriesLineChart;
