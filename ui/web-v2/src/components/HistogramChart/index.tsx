import { FC } from 'react';
import { Bar } from 'react-chartjs-2';

import { COLORS } from '../../constants/colorPattern';

interface HistogramChartProps {
  label: string;
  dataLabels: Array<string>;
  hist: Array<Array<number>>;
  bins: Array<number>;
}

export const HistogramChart: FC<HistogramChartProps> = ({
  label,
  dataLabels,
  hist,
  bins,
}) => {
  const chartData = {
    labels: bins,
    datasets: dataLabels.map((e, i) => {
      return {
        label: e,
        data: hist[i],
        backgroundColor: COLORS[i % COLORS.length],
      };
    }),
  };
  const options = {
    title: {
      display: true,
      text: label,
      fontStyle: 'normal',
    },
    scales: {
      xAxes: [
        {
          display: false,
          barPercentage: 1.5,
          ticks: {
            max: hist[0].length - 1,
          },
        },
        {
          display: true,
          ticks: {
            autoSkip: true,
            max: hist[0].length,
          },
        },
      ],
      yAxes: [
        {
          display: false,
          ticks: {
            beginAtZero: true,
          },
        },
      ],
    },
  };

  return <Bar data={chartData} options={options} />;
};
