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
    const labels = [
      'January',
      'February',
      'March',
      'April',
      'May',
      'June',
      'July',
    ];

    const chartData = {
      labels,
      datasets: [
        {
          label: 'Dataset 1',
          data: [20000, 4, 3, 8, 6, 10, 3],
          borderColor: 'rgb(255, 99, 132)',
          backgroundColor: 'rgba(255, 99, 132, 0.5)',
          fill: false,
        },
        {
          label: 'Dataset 2',
          data: [1, 4, 13000, 2, 7, 1000, 5],
          borderColor: 'rgb(53, 162, 235)',
          backgroundColor: 'rgba(53, 162, 235, 0.5)',
          fill: false,
        },
      ],
      // labels: timeseries.map((t) => new Date(t * 1000)),
      // datasets: dataLabels.map((e, i) => {
      //   return {
      //     label: e.length > 40 ? `${e.substring(0, 40)}...` : e,
      //     data: data[i],
      //     backgroundColor: COLORS[i % COLORS.length],
      //     borderColor: COLORS[i % COLORS.length],
      //     fill: false,
      //   };
      // }),
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
              parseInt(tooltipItem.value).toLocaleString()
            );
          },
        },
      },
      scales: {
        // xAxes: [
        //   {
        //     type: 'time',
        //     time: {
        //       unit: 'hour',
        //       displayFormats: { hour: 'HH:mm' },
        //     },
        //   },
        // ],
        yAxes: [
          {
            display: true,
            stacked: false,
            ticks: {
              // beginAtZero: true,
              userCallback: function (value) {
                return value.toLocaleString();
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
