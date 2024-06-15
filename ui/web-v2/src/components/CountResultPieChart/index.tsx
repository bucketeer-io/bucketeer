import { jsx } from '@emotion/react';
import { FC } from 'react';
import { Pie } from 'react-chartjs-2';

import { COLORS } from '../../constants/colorPattern';

interface CountResultPieChartProps {
  label: string;
  variationValues: string[];
  data: number[];
}

export const CountResultPieChart: FC<CountResultPieChartProps> = ({
  label,
  variationValues,
  data,
}) => {
  const chartData = {
    labels: variationValues,
    datasets: [
      {
        data: data,
        backgroundColor: COLORS.slice(
          0,
          variationValues.length % COLORS.length
        ),
      },
    ],
  };
  const options = {
    title: {
      display: true,
      text: label,
      fontStyle: 'normal',
    },
  };

  return (
    <div>
      <Pie data={chartData} options={options} />
    </div>
  );
};
