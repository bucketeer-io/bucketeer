import { AnyObject } from 'yup';
import { VariationResult } from '@types';
import {
  ChartDataType,
  GoalResultTab
} from 'pages/experiment-details/collection-loader/results';
import { formatLongDateTime } from './date-time';

export const getTimeSeries = (
  variationResults: VariationResult[],
  type: ChartDataType,
  tab: GoalResultTab
) => {
  if (tab === 'CONVERSION') {
    return variationResults[0]?.goalEventCountTimeseries?.timestamps || [];
  }
  switch (type) {
    case 'goal-total':
      return variationResults[0]?.goalEventCountTimeseries?.timestamps || [];
    case 'goal-user':
      return variationResults[0]?.goalUserCountTimeseries?.timestamps || [];
    case 'conversion-rate':
      return variationResults[0]?.cvrTimeseries?.timestamps || [];
    case 'value-total':
      return variationResults[0]?.goalValueSumTimeseries?.timestamps || [];
    case 'value-user':
      return (
        variationResults[0]?.goalValueSumPerUserTimeseries?.timestamps || []
      );
    case 'evaluation-user':
    default:
      return (
        variationResults[0]?.evaluationUserCountTimeseries?.timestamps || []
      );
  }
};

export const getData = (
  variationResults: VariationResult[],
  type: ChartDataType
) => {
  switch (type) {
    case 'goal-total':
      return (
        variationResults?.map(item => item?.goalEventCountTimeseries?.values) ||
        []
      );
    case 'goal-user':
      return (
        variationResults?.map(item => item?.goalUserCountTimeseries?.values) ||
        []
      );
    case 'conversion-rate':
      return variationResults?.map(item => item?.cvrTimeseries?.values) || [];
    case 'value-total':
      return (
        variationResults?.map(item => item?.goalValueSumTimeseries?.values) ||
        []
      );
    case 'value-user':
      return (
        variationResults?.map(
          item => item?.goalValueSumPerUserTimeseries?.values
        ) || []
      );
    case 'evaluation-user':
    default:
      return (
        variationResults?.map(
          item => item?.evaluationUserCountTimeseries?.values
        ) || []
      );
  }
};

export const isNumber = (n: number) =>
  typeof n === 'number' && !isNaN(n) && n !== Infinity;

export const formatTooltipLabel = ({
  formattedValue,
  dataset
}: {
  formattedValue: string;
  dataset: AnyObject;
}) => {
  const value = Number(formattedValue);
  return (
    (dataset?.label ? `${dataset.label}: ` : '') +
    (isNumber(value) ? value.toFixed(1) : formattedValue)
  );
};

export const formatXAxisLabel = (index: number, labels: Date[]) => {
  const date = labels[index];
  if (date) {
    return formatLongDateTime({
      value: String(date.getTime() / 1000),
      overrideOptions: {
        day: '2-digit',
        month: '2-digit',
        year: undefined
      },
      locale: 'en-US'
    }).replace(/(\d{2}) (\w{3})/, '$2 $1');
  }
  return '';
};
