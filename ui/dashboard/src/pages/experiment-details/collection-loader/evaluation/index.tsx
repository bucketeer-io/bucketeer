import { useCallback, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import LineChart, { ChartDataType } from './chart';
import EvaluationTable from './evaluation-table';

type Option = {
  label: string;
  value: string;
};

const evaluationOptions: Option[] = [
  {
    label: 'Option 1',
    value: 'option-1'
  },
  {
    label: 'Option 2',
    value: 'option-2'
  },
  {
    label: 'Option 3',
    value: 'option-3'
  }
];

export const clientRequestDataPoints = [
  { x: '2024-01-01', y: 10, type: 'GetEvaluations' },
  { x: '2024-01-29', y: 15, type: 'GetEvaluations' },
  { x: '2024-02-03', y: 30, type: 'GetEvaluations' },
  { x: '2024-02-27', y: 25, type: 'GetEvaluations' },
  { x: '2024-03-03', y: 2, type: 'GetEvaluations' },
  { x: '2024-03-08', y: 19, type: 'GetEvaluations' },
  { x: '2024-03-19', y: 10, type: 'GetEvaluations' },
  { x: '2024-03-28', y: 44, type: 'GetEvaluations' },
  { x: '2024-04-01', y: 30, type: 'GetEvaluations' },
  { x: '2024-04-07', y: 10, type: 'GetEvaluations' },
  { x: '2024-04-28', y: 26, type: 'GetEvaluations' },
  { x: '2024-05-01', y: 11, type: 'GetEvaluations' },
  { x: '2024-05-21', y: 10, type: 'GetEvaluations' },
  { x: '2024-06-06', y: 10, type: 'GetEvaluations' },
  { x: '2024-06-25', y: 15, type: 'GetEvaluations' },
  { x: '2024-07-04', y: 30, type: 'GetEvaluations' },
  { x: '2024-07-17', y: 40, type: 'GetEvaluations' },
  { x: '2024-07-29', y: 25, type: 'GetEvaluations' },
  { x: '2024-08-02', y: 2, type: 'GetEvaluations' },
  { x: '2024-08-28', y: 10, type: 'GetEvaluations' },
  { x: '2024-09-03', y: 44, type: 'GetEvaluations' },
  { x: '2024-09-27', y: 30, type: 'GetEvaluations' },
  { x: '2024-10-12', y: 10, type: 'GetEvaluations' },
  { x: '2024-10-29', y: 26, type: 'GetEvaluations' },
  { x: '2024-11-21', y: 11, type: 'GetEvaluations' },
  { x: '2024-12-11', y: 10, type: 'GetEvaluations' },
  { x: '2024-12-28', y: 10, type: 'GetEvaluations' }
];

export const clientRequestDataPoints2 = [
  { x: '2024-01-01', y: 20, type: 'RegisterEvents' },
  { x: '2024-01-29', y: 25, type: 'RegisterEvents' },
  { x: '2024-02-03', y: 40, type: 'RegisterEvents' },
  { x: '2024-02-12', y: 50, type: 'RegisterEvents' },
  { x: '2024-03-03', y: 12, type: 'RegisterEvents' },
  { x: '2024-03-08', y: 29, type: 'RegisterEvents' },
  { x: '2024-03-19', y: 20, type: 'RegisterEvents' },
  { x: '2024-03-28', y: 54, type: 'RegisterEvents' },
  { x: '2024-04-01', y: 40, type: 'RegisterEvents' },
  { x: '2024-04-07', y: 20, type: 'RegisterEvents' },
  { x: '2024-04-28', y: 36, type: 'RegisterEvents' },
  { x: '2024-05-08', y: 20, type: 'RegisterEvents' },
  { x: '2024-05-21', y: 26, type: 'RegisterEvents' },
  { x: '2024-06-06', y: 20, type: 'RegisterEvents' },
  { x: '2024-06-25', y: 25, type: 'RegisterEvents' },
  { x: '2024-07-04', y: 40, type: 'RegisterEvents' },
  { x: '2024-07-17', y: 50, type: 'RegisterEvents' },
  { x: '2024-07-29', y: 35, type: 'RegisterEvents' },
  { x: '2024-08-16', y: 29, type: 'RegisterEvents' },
  { x: '2024-08-28', y: 20, type: 'RegisterEvents' },
  { x: '2024-09-03', y: 54, type: 'RegisterEvents' },
  { x: '2024-09-27', y: 40, type: 'RegisterEvents' },
  { x: '2024-10-12', y: 20, type: 'RegisterEvents' },
  { x: '2024-10-29', y: 36, type: 'RegisterEvents' },
  { x: '2024-11-21', y: 21, type: 'RegisterEvents' },
  { x: '2024-12-11', y: 20, type: 'RegisterEvents' },
  { x: '2024-12-28', y: 15, type: 'RegisterEvents' }
];

export const getMonthName = (month: number) => {
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
};

const Evaluation = () => {
  const { t } = useTranslation(['common', 'form']);

  const [selectedEvaluation, setSelectedEvaluation] = useState<Option | null>(
    null
  );
  const saveXAxisRef = useRef('');

  const labels = clientRequestDataPoints.map(item => {
    const date = new Date(item.x);
    const month = getMonthName(date?.getMonth());
    return month;
  });

  const formatLabel = useCallback(
    (index: number) => {
      const currentValue = labels[index];
      if (currentValue !== saveXAxisRef.current) {
        saveXAxisRef.current = currentValue;
        return saveXAxisRef.current;
      }
      return '';
    },
    [labels]
  );
  const chartData: ChartDataType = {
    labels,
    datasets: [
      {
        data: clientRequestDataPoints.map(p => p.y),
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
        data: clientRequestDataPoints2.map(p => p.y),
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

  return (
    <div className="flex flex-col h-full gap-y-6">
      <DropdownMenu>
        <DropdownMenuTrigger
          label={selectedEvaluation?.label || ''}
          placeholder={t(`form:select-evaluation`)}
          variant="secondary"
          className="w-[300px] lg:w-[528px]"
        />
        <DropdownMenuContent className="w-[235px]" align="start">
          {evaluationOptions.map((item, index) => (
            <DropdownMenuItem
              key={index}
              value={item.value}
              label={item.label}
              onSelectOption={() => setSelectedEvaluation(item)}
            />
          ))}
        </DropdownMenuContent>
      </DropdownMenu>
      <EvaluationTable />
      <LineChart chartData={chartData} formatLabel={formatLabel} />
    </div>
  );
};

export default Evaluation;
