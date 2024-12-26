import { useMemo } from 'react';
import { Line } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  LineElement,
  CategoryScale,
  LinearScale,
  TimeScale,
  PointElement,
  Tooltip,
  Filler,
  ChartData,
  ChartOptions
} from 'chart.js';

ChartJS.register(
  LineElement,
  CategoryScale,
  LinearScale,
  TimeScale,
  PointElement,
  Tooltip,
  Filler
);

export type TicksType = {
  label: string;
  value: string;
};

export type ChartDataType = ChartData<'line', number[], string>;

type Props = {
  chartData: ChartDataType;
  formatLabel?: (index: number) => string;
};

const LineChart = ({ chartData, formatLabel }: Props) => {
  const options: ChartOptions<'line'> = useMemo(
    () => ({
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          display: false
        },
        tooltip: {
          enabled: true
        }
      },
      scales: {
        x: {
          title: {
            display: false
          },
          grid: {
            display: true,
            color: 'rgba(192, 192, 192, 0.2)', // Custom grid line color for x-axis
            lineWidth: 1 // Custom line width
          },
          ticks: {
            align: 'center' as const,
            callback: (index: string | number) => {
              return formatLabel && formatLabel(+index);
            },
            autoSkip: false,
            minRotation: 0,
            maxRotation: 0
          }
        },
        y: {
          title: {
            display: false
          },
          grid: {
            display: true,
            color: 'rgba(192, 192, 192, 0.2)',
            lineWidth: 1
          },
          beginAtZero: true
        }
      }
    }),
    []
  );

  return (
    <div className="flex flex-1 w-full xl:w-[60%] min-w-[650px] h-fit pr-5 border-r border-gray-600/15">
      <Line
        style={{
          width: '100%'
        }}
        data={chartData}
        options={options}
      />
    </div>
  );
};

export default LineChart;
