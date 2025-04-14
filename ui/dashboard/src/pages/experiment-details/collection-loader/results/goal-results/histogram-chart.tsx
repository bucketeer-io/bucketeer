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
import { formatLongDateTime } from 'utils/date-time';
import { getVariationColor } from 'utils/style';

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
  dataLabels: Array<string>;
  hist: Array<Array<number | string>>;
  bins: Array<number>;
}

export const HistogramChart = memo(
  ({ label, dataLabels, hist, bins }: HistogramChartProps) => {
    const chartData = {
      labels: bins,
      datasets: dataLabels.map((e, i) => {
        return {
          label: e,
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
          display: label ? true : false,
          text: label
        },
        tooltip: {
          enabled: true,
          callbacks: {
            title: (tooltipItems: { label: string }[]) => {
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
          display: false,
          barPercentage: 1.5,
          ticks: {
            max: hist[0].length - 1,
            autoSkip: true,
            beginAtZero: true
          }
        },
        y: {
          display: false,
          ticks: {
            autoSkip: true,
            beginAtZero: true,
            max: hist[0].length
          }
        }
      }
    };

    return <Bar data={chartData} options={options} />;
  }
);
