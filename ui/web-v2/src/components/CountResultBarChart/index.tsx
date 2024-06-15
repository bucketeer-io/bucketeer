import { jsx } from '@emotion/react';
import { FC } from 'react';
import { Bar } from 'react-chartjs-2';

import { COLORS } from '../../constants/colorPattern';

interface CountResultBarChartProps {
  label: string;
  variationValues: string[];
  data: number[];
}

export const CountResultBarChart: FC<CountResultBarChartProps> = ({
  label,
  variationValues,
  data,
}) => {
  const chartData = {
    labels: variationValues,
    datasets: [
      {
        label: '',
        backgroundColor: COLORS.slice(
          0,
          variationValues.length % COLORS.length
        ),
        borderWidth: 1,
        data: data,
      },
    ],
  };
  const options = {
    legend: {
      display: false,
    },
    title: {
      display: true,
      text: label,
      fontStyle: 'normal',
    },
    scales: {
      yAxes: [
        {
          ticks: {
            beginAtZero: true,
          },
        },
      ],
    },
  };

  return (
    <div>
      <Bar data={chartData} options={options} />
    </div>
  );
};
