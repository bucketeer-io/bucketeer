import { memo } from 'react';
import { Bar } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  BarElement,
  CategoryScale,
  LinearScale,
  TimeScale,
  TimeSeriesScale,
  PointElement,
  Tooltip,
  Filler
} from 'chart.js';
import { getVariationColor } from 'utils/style';
import { DataLabel } from './timeseries-line-chart';

ChartJS.register(
  BarElement,
  CategoryScale,
  LinearScale,
  TimeScale,
  PointElement,
  TimeSeriesScale,
  Tooltip,
  Filler
);

interface HistogramChartProps {
  label: string;
  dataLabels: Array<DataLabel>;
  hist: Array<Array<number | string>>;
  bins: Array<number>;
}

export const HistogramChart = memo(
  ({ label, dataLabels, hist, bins }: HistogramChartProps) => {
    const chartData = {
      labels: bins,
      datasets: dataLabels?.map((e, i) => {
        return {
          label: e.label,
          data: hist[i] || [],
          backgroundColor: getVariationColor(i)
        };
      })
    };

    const options = {
      plugins: {
        legend: {
          display: false
        },
        title: {
          display: !!label,
          text: label
        },
        tooltip: {
          enabled: true
        }
      },
      scales: {
        x: {
          display: true,
          barPercentage: 1.5,
          ticks: {
            max: hist[0]?.length - 1,
            autoSkip: true,
            beginAtZero: true
          }
        },
        y: {
          display: true,
          ticks: {
            autoSkip: true,
            beginAtZero: true,
            max: hist[0]?.length
          }
        }
      }
    };

    return <Bar data={chartData} options={options} />;
  }
);
