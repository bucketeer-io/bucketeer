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
  unit: string;
}

export const TimeseriesStackedLineChart: FC<TimeseriesStackedLineChartProps> =
  ({ label, dataLabels, timeseries, data, unit }) => {
    const chartData = {
      labels: timeseries.map((t) => new Date(t * 1000)),
      datasets: dataLabels.map((e, i) => {
        return {
          label: e.length > 40 ? `${e.substring(0, 40)}...` : e,
          data: data[i],
          backgroundColor: COLORS[i % COLORS.length],
          borderColor: COLORS[i % COLORS.length],
          fill: false,
        };
      }),
    };
    const options: ChartOptions = {
      title: {
        display: true,
        text: label,
        fontStyle: 'normal',
      },
      tooltips: {
        callbacks: {
          label: function (tooltipItem, data) {
            return (
              data.datasets[tooltipItem.datasetIndex].label +
              ': ' +
              Number(tooltipItem.value).toLocaleString()
            );
          },
        },
      },
      scales: {
        xAxes: [
          {
            type: 'time',
            time: {
              unit,
              displayFormats: { hour: 'HH:mm' },
            },
          },
        ],
        yAxes: [
          {
            display: true,
            stacked: false,
            ticks: {
              // beginAtZero: true,
              userCallback: function (value) {
                return Number(value).toLocaleString();
              },
            },
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
