import { ChartOptions } from 'chart.js';
import { FC, memo } from 'react';
import { Line } from 'react-chartjs-2';

import { COLORS } from '../../constants/colorPattern';

interface TimeseriesLineChartProps {
  label: string;
  dataLabels: Array<string>;
  timeseries: Array<number>;
  data: Array<Array<number>>;
  height: number;
}

export const TimeseriesLineChart: FC<TimeseriesLineChartProps> = memo(
  ({ label, dataLabels, timeseries, data, height }) => {
    const chartData = {
      labels: timeseries.map((t) => new Date(t * 1000)),
      datasets: dataLabels.map((e, i) => {
        const color = COLORS[i % COLORS.length];
        return {
          label: e,
          data: [...data[i]], // Copy arrays to avoid  "Uncaught TypeError: Cannot assign to read only property 'length' of object '[object Array]""
          borderColor: color,
          backgroundColor: color,
          fill: false,
        };
      }),
    };
    const options: ChartOptions = {
      title: {
        display: label == '' ? false : true,
        text: label,
        fontStyle: 'normal',
      },
      scales: {
        xAxes: [
          {
            type: 'time',
            time: {
              unit: 'day',
            },
          },
        ],
        yAxes: [
          {
            display: true,
          },
        ],
      },
      responsive: true,
      maintainAspectRatio: false,
    };

    return <Line data={chartData} options={options} height={height} />;
  }
);
