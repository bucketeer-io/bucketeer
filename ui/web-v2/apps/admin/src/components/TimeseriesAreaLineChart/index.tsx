import { ChartOptions } from 'chart.js';
import { FC } from 'react';
import { Line } from 'react-chartjs-2';

import { COLORS } from '../../constants/colorPattern';

interface TimeseriesAreaLineChartProps {
  label: string;
  dataLabels: Array<string>;
  timeseries: Array<number>;
  upperBoundaries: Array<Array<number>>;
  lowerBoundaries: Array<Array<number>>;
  representatives: Array<Array<number>>;
  height: number;
}

// TimeseriesAreaLineChart takes 3 sets of data for each metric.
// 1. upper boundary data (e.g. 97.5 percentile)
// 2. lower boundary data (e.g. 2.5 percentile)
// 2. representative data (e.g. median)
export const TimeseriesAreaLineChart: FC<TimeseriesAreaLineChartProps> = ({
  label,
  dataLabels,
  timeseries,
  upperBoundaries,
  lowerBoundaries,
  representatives,
  height,
}) => {
  const datasets = Array<any>();
  dataLabels.forEach((l, i) => {
    const color = COLORS[i % COLORS.length];
    datasets.push({
      label: null,
      data: upperBoundaries[i],
      borderWidth: 0,
      backgroundColor: hexToRgba(color, 0.2),
      pointRadius: 0,
      fill: '+1',
    });
    datasets.push({
      label: null,
      data: lowerBoundaries[i],
      borderWidth: 0,
      backgroundColor: hexToRgba(color, 0.2),
      pointRadius: 0,
      fill: '-1',
    });
    datasets.push({
      label: l,
      data: representatives[i],
      borderColor: COLORS[i % COLORS.length],
      pointRadius: 0,
      fill: false,
    });
  });
  const chartData = {
    labels: timeseries.map((t) => new Date(t * 1000)),
    datasets: datasets,
  };
  const options: ChartOptions = {
    title: {
      display: label == '' ? false : true,
      text: label,
      fontStyle: 'normal',
    },
    legend: {
      display: true,
      labels: {
        filter: (legendItem: any, _: any) => {
          return !!legendItem.text;
        },
      },
      onClick: undefined,
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

  return (
    <div>
      <Line
        data={chartData}
        options={options}
        height={height}
        datasetKeyProvider={datasetKeyProvider}
      />
    </div>
  );
};

const hexToRgba = (hex: string, alpha: number) => {
  const r = parseInt(hex.slice(1, 3), 16),
    g = parseInt(hex.slice(3, 5), 16),
    b = parseInt(hex.slice(5, 7), 16);
  return 'rgba(' + r + ', ' + g + ', ' + b + ', ' + alpha + ')';
};

const datasetKeyProvider = () => {
  return Math.random();
};
