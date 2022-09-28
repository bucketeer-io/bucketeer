import { jsx } from '@emotion/react';
import { ChartOptions } from 'chart.js';
import { FC } from 'react';
import { Line } from 'react-chartjs-2';

import { COLORS } from '../../constants/colorPattern';

interface TimeseriesStackedLineChartProps {
  label: string;
  dataLabels: Array<string>;
  timeseries: Array<number>;
  data: Array<Array<number>>;
}

export const TimeseriesStackedLineChart: FC<TimeseriesStackedLineChartProps> =
  ({ label, dataLabels, timeseries, data }) => {
    const chartData = {
      labels: timeseries.map((t) => new Date(t * 1000)),
      datasets: dataLabels.map((e, i) => {
        return {
          label: e,
          data: data[i],
          backgroundColor: COLORS[i % COLORS.length],
        };
      }),
    };
    const options: ChartOptions = {
      title: {
        display: true,
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
            stacked: true,
          },
        ],
      },
    };

    return (
      <div>
        <Line data={chartData} options={options} />
      </div>
    );
  };
