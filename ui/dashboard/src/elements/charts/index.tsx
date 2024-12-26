import { useCallback, useRef } from 'react';
import { cn } from 'utils/style';
import { DataTable } from 'elements/data-table';
import ChartHeader, { ChartHeaderProps } from './chart-header';
import { TempTableDataType, useColumns } from './data-collection';
import LineChart, { ChartDataType } from './line-chart';

type ChartWrapperProps = ChartHeaderProps & {};

const dataPoints = [
  { x: '2024-01-02', y: 10 },
  { x: '2024-01-10', y: 20 },
  { x: '2024-01-29', y: 15 },
  { x: '2024-02-03', y: 30 },
  { x: '2024-02-12', y: 40 },
  { x: '2024-02-27', y: 25 },
  { x: '2024-03-03', y: 2 },
  { x: '2024-03-08', y: 19 },
  { x: '2024-03-19', y: 10 },
  { x: '2024-03-28', y: 44 },
  { x: '2024-04-01', y: 30 },
  { x: '2024-04-07', y: 10 },
  { x: '2024-04-28', y: 26 },
  { x: '2024-05-01', y: 11 },
  { x: '2024-05-08', y: 10 },
  { x: '2024-05-21', y: 10 },
  { x: '2024-06-06', y: 10 },
  { x: '2024-06-19', y: 20 },
  { x: '2024-06-25', y: 15 },
  { x: '2024-07-04', y: 30 },
  { x: '2024-07-17', y: 40 },
  { x: '2024-07-29', y: 25 },
  { x: '2024-08-02', y: 2 },
  { x: '2024-08-16', y: 19 },
  { x: '2024-08-28', y: 10 },
  { x: '2024-09-03', y: 44 },
  { x: '2024-09-27', y: 30 },
  { x: '2024-10-12', y: 10 },
  { x: '2024-10-29', y: 26 },
  { x: '2024-11-21', y: 11 },
  { x: '2024-12-11', y: 10 },
  { x: '2024-12-28', y: 10 }
];

const dataPoints2 = [
  { x: '2024-01-02', y: 5 },
  { x: '2024-01-10', y: 15 },
  { x: '2024-01-29', y: 10 },
  { x: '2024-02-03', y: 20 },
  { x: '2024-02-12', y: 30 },
  { x: '2024-02-27', y: 15 },
  { x: '2024-03-03', y: 1 },
  { x: '2024-03-08', y: 10 },
  { x: '2024-03-19', y: 5 },
  { x: '2024-03-28', y: 35 },
  { x: '2024-04-01', y: 25 },
  { x: '2024-04-07', y: 5 },
  { x: '2024-04-28', y: 20 },
  { x: '2024-05-01', y: 5 },
  { x: '2024-05-08', y: 4 },
  { x: '2024-05-21', y: 6 },
  { x: '2024-06-06', y: 7 },
  { x: '2024-06-19', y: 15 },
  { x: '2024-06-25', y: 10 },
  { x: '2024-07-04', y: 25 },
  { x: '2024-07-17', y: 30 },
  { x: '2024-07-29', y: 15 },
  { x: '2024-08-02', y: 1 },
  { x: '2024-08-16', y: 10 },
  { x: '2024-08-28', y: 4 },
  { x: '2024-09-03', y: 35 },
  { x: '2024-09-27', y: 20 },
  { x: '2024-10-12', y: 4 },
  { x: '2024-10-29', y: 17 },
  { x: '2024-11-21', y: 7 },
  { x: '2024-12-11', y: 8 },
  { x: '2024-12-28', y: 9 }
];

const chartData: ChartDataType = {
  // labels: ['Jan', 'Feb', "March", 'April', 'May', 'June', 'July', 'Aug', 'Sep', 'Oct', "Nov", "Dec"]
  datasets: [
    {
      data: dataPoints,
      borderColor: 'rgba(228, 57, 172, 1)',
      fill: true,
      backgroundColor: 'rgba(228, 57, 172, 0.1)',
      borderWidth: 1,
      pointBackgroundColor: 'transparent',
      pointBorderColor: 'transparent',
      pointHoverRadius: 6,
      pointHoverBackgroundColor: 'rgba(228, 57, 172, 1)'
      // tension: 0.1, --> smooth the line
    },
    {
      data: dataPoints2,
      borderColor: 'rgba(87, 55, 146, 1)',
      fill: true,
      backgroundColor: 'rgba(87, 55, 146, 0.1)',
      borderWidth: 1,
      pointBackgroundColor: 'transparent',
      pointBorderColor: 'transparent',
      pointHoverRadius: 6,
      pointHoverBackgroundColor: 'rgba(87, 55, 146, 1)'
      // tension: 0.1, --> smooth the line
    }
  ]
};

const tableData: TempTableDataType[] = [
  {
    name: 'GetEvaluations',
    min: '2,495 K',
    max: '6,362 K',
    current: '3,123 K'
  },
  {
    name: 'RegisterEvents',
    min: '1,232 K',
    max: '5,374 K',
    current: '3,314 K'
  }
];

const ChartWrapper = ({ ...props }: ChartWrapperProps) => {
  const saveXAxisRef = useRef('');

  const renderName = useCallback((tempData: TempTableDataType) => {
    return (
      <div className="flex items-center w-full gap-x-2">
        <div
          className={cn('size-3 rounded-sm bg-accent-pink-500', {
            'bg-primary-500': tempData.name === 'RegisterEvents'
          })}
        ></div>
        <p>{tempData.name}</p>
      </div>
    );
  }, []);

  const columns = useColumns({ renderName });

  const getMonthName = useCallback((month: number) => {
    switch (month) {
      case 1:
        return 'Jan';
      case 2:
        return 'Feb';
      case 3:
        return 'March';
      case 4:
        return 'April';
      case 5:
        return 'May';
      case 6:
        return 'June';
      case 7:
        return 'July';
      case 8:
        return 'Aug';
      case 9:
        return 'Sep';
      case 10:
        return 'Oct';
      case 11:
        return 'Nov';
      case 12:
        return 'Dec';
      default:
        return '';
    }
  }, []);

  const formatLabel = useCallback(
    (index: number) => {
      const month = new Date(dataPoints[index]?.x)?.getMonth();
      const currentValue = getMonthName(month + 1);
      if (currentValue !== saveXAxisRef.current) {
        saveXAxisRef.current = currentValue;
        return saveXAxisRef.current;
      }
      return '';
    },
    [dataPoints]
  );

  return (
    <div className="flex flex-col w-fit border border-gray-200 rounded-2xl">
      <ChartHeader {...props} />
      <div className="flex w-full divide-x divide-gray-600/15">
        <LineChart chartData={chartData} formatLabel={formatLabel} />
        <div className="flex flex-1 min-w-[500px] p-5">
          <DataTable
            data={tableData}
            columns={columns}
            onSortingChange={() => {}}
          />
        </div>
      </div>
    </div>
  );
};

export default ChartWrapper;
