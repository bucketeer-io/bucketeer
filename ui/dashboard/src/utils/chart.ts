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
  const {
    goalEventCountTimeseries,
    goalUserCountTimeseries,
    cvrTimeseries,
    goalValueSumTimeseries,
    goalValueSumPerUserTimeseries,
    evaluationUserCountTimeseries
  } = variationResults[0] || {};
  if (tab === 'CONVERSION') {
    return goalEventCountTimeseries?.timestamps || [];
  }
  switch (type) {
    case 'goal-total':
      return goalEventCountTimeseries?.timestamps || [];
    case 'goal-user':
      return goalUserCountTimeseries?.timestamps || [];
    case 'conversion-rate':
      return cvrTimeseries?.timestamps || [];
    case 'value-total':
      return goalValueSumTimeseries?.timestamps || [];
    case 'value-user':
      return goalValueSumPerUserTimeseries?.timestamps || [];
    case 'evaluation-user':
    default:
      return evaluationUserCountTimeseries?.timestamps || [];
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
        day: 'numeric',
        month: 'short',
        year: undefined
      },
      locale: 'en-US'
    }).replace(/(\d{2}) (\w{3})/, '$2 $1');
  }
  return '';
};
