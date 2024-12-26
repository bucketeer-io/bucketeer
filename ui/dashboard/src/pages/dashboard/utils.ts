import { TempChartData } from './page-content';

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

export const findMinMax = (
  array: TempChartData[],
  field: keyof TempChartData
): { min: number; max: number } => {
  if (array.length === 0) {
    throw new Error('Array is empty');
  }

  return array.reduce(
    (acc, item) => {
      const value = item[field] as number;
      acc.min = Math.min(acc.min, value);
      acc.max = Math.max(acc.max, value);
      return acc;
    },
    { min: Infinity, max: -Infinity }
  );
};

export const formatNumber = (
  num: number,
  options?: Intl.NumberFormatOptions
): string => {
  return new Intl.NumberFormat('en-US', options).format(num);
};
