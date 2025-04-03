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
  ChartOptions
} from 'chart.js';
import { formatTooltipLabel, formatXAxisLabel } from 'utils/chart';
import { formatLongDateTime } from 'utils/date-time';
import { getVariationColor } from 'utils/style';
<<<<<<< HEAD
import {
  ChartToggleLegendRef,
  DatasetReduceType
} from './timeseries-area-line-chart';
=======
>>>>>>> 76e7257b (fix: decrease the default delay duration of the tooltip)

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
      const labels = timeseries.map(t => new Date(Number(t) * 1000));

<<<<<<< HEAD
      const chartRef = useRef<ChartJS<
        'line',
        (string | number)[],
        Date
      > | null>(null);

      useImperativeHandle(ref, () => {
=======
    const chartData: ChartData<'line', (string | number)[], Date> = {
      labels,
      datasets: dataLabels.map((e, i) => {
        const color = getVariationColor(i);
>>>>>>> 76e7257b (fix: decrease the default delay duration of the tooltip)
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

      const chartData: ChartData<'line', (string | number)[], Date> = {
        labels,
        datasets: dataLabels.map((e, i) => {
          const color = getVariationColor(i);
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
