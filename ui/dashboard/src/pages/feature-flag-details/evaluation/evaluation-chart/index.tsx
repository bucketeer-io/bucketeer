import { forwardRef, Ref, useImperativeHandle, useRef } from 'react';
import { Line } from 'react-chartjs-2';
import {
  CategoryScale,
  ChartData,
  ChartDataset,
  Chart as ChartJS,
  ChartOptions,
  Filler,
  Legend,
  LinearScale,
  LineElement,
  Point,
  PointElement,
  TimeScale,
  TimeSeriesScale,
  TimeUnit,
  Tooltip
} from 'chart.js';
import 'chartjs-adapter-luxon';
import { COLORS } from 'constants/styles';
import { formatLongDateTime } from 'utils/date-time';
import { getVariationColor } from 'utils/style';
import {
  ChartToggleLegendRef,
  DatasetReduceType
} from 'pages/experiment-details/collection-loader/results/goal-results/timeseries-area-line-chart';
import { Option } from 'components/creatable-select';
import { RawPoint } from '../types';
import { formatLabel, getLogBase, symlog, symlogInverse } from '../utils';

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
    const isDark = document.documentElement.classList.contains('dark');
    const tickColor = isDark ? '#B5B0C2' : '#94A3B8';
    const gridColor = isDark ? 'rgba(181, 176, 194, 0.25)' : '#E2E8F0';

    const maxValue = Math.max(...data.flat());
    const useSymlog = maxValue > 100;
    const logBase = getLogBase(maxValue);

    const chartData: ChartData<
      'line',
      (number | Point | RawPoint | null)[],
      Date
    > = {
      labels: timeseries.map(t => new Date(Number(t) * 1000)),
      datasets: variationValues.map((e, i) => {
        const color = getVariationColor(i % COLORS.length);
        return {
          label:
            e.label.length > 40 ? `${e.label.substring(0, 40)}...` : e.label,
          data: data[i].map((v, idx) => ({
            x: new Date(Number(timeseries[idx]) * 1000),
            y: useSymlog ? symlog(v, 1, logBase) : v,
            raw: v
          })),
          backgroundColor: color,
          borderColor: color,
          borderWidth: isDark ? 2.5 : 2,
          pointRadius: isDark ? 3 : 2,
          pointHoverRadius: isDark ? 5 : 4,
          fill: false,
          tension: 0.2,
          value: e.value
        };
      })
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
            label: context => {
              const parsedY = context.parsed.y ?? 0;
              const displayValue = useSymlog
                ? symlogInverse(parsedY, 1, logBase)
                : parsedY;
              return `${context.dataset.label ?? ''} : ${displayValue.toLocaleString()}`;
            }
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
            source: 'data',
            font: {
              family: 'Sofia Pro',
              size: 14,
              weight: 400
            },
            color: tickColor
          },
          grid: {
            display: true,
            color: isDark ? 'rgba(181, 176, 194, 0.25)' : '#E2E8F0',
            lineWidth: 1,
            tickWidth: 0
          },
          border: {
            display: true,
            color: isDark ? 'rgba(181, 176, 194, 0.4)' : '#E2E8F0',
            width: 1
          }
        },
        y: {
          max: useSymlog ? symlog(maxValue * 1.1, 1, logBase) : undefined,
          type: 'linear',
          title: {
            display: false
          },
          display: true,
          beginAtZero: true,
          ticks: {
            callback: v =>
              formatLabel(
                useSymlog
                  ? symlogInverse(v as number, 1, logBase)
                  : (v as number)
              ),

            font: {
              family: 'Sofia Pro',
              size: 14,
              weight: 400
            },
            color: tickColor
          },
          afterBuildTicks: axis => {
            if (useSymlog) {
              const tickValues: number[] = [];
              let v = 1;
              while (v <= maxValue * 1.1) {
                tickValues.push(v);
                v *= logBase;
              }
              axis.ticks = tickValues.map(value => ({
                value: symlog(value, 1, logBase),
                label: formatLabel(value)
              }));
            }
          },
          grid: {
            display: true,
            color: gridColor,
            lineWidth: 1,
            tickWidth: 0
          },
          border: {
            display: true,
            color: isDark ? 'rgba(181, 176, 194, 0.4)' : '#E2E8F0',
            width: 1
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
