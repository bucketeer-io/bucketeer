import { useMemo } from 'react';
import { useTranslation } from 'i18n';
import { Experiment, GoalResult } from '@types';
import { ResultCell, ResultHeaderCell } from './goal-results-table-element';

const headerList = [
  {
    name: 'variation',
    tooltipKey: '',
    minSize: 270
  },
  {
    name: 'evaluation-user',
    tooltipKey: 'evaluation-user-tooltip',
    minSize: 143
  },
  {
    name: 'goal-total',
    tooltipKey: 'goal-total-tooltip',
    minSize: 119
  },
  {
    name: 'goal-user',
    tooltipKey: 'goal-user-tooltip',
    minSize: 123
  },
  {
    name: 'conversion-rate',
    tooltipKey: 'conversion-rate-tooltip',
    minSize: 147
  },
  {
    name: 'value-total',
    tooltipKey: 'value-total-tooltip',
    minSize: 125
  },
  {
    name: 'value-user',
    tooltipKey: 'value-user-tooltip',
    minSize: 123
  }
];

const EvaluationTable = ({
  experiment,
  goalResult
}: {
  experiment: Experiment;
  goalResult: GoalResult;
}) => {
  const { t } = useTranslation(['common', 'table']);

  const evaluationData = useMemo(
    () =>
      goalResult?.variationResults?.map(item => {
        const variation = experiment?.variations?.find(
          variation => item.variationId === variation.id
        );
        return {
          ...item,
          variationName: variation?.value || variation?.name || ''
        };
      }),
    [goalResult, experiment]
  );

  return (
    <div className="min-w-fit">
      <div className="flex w-full">
        {headerList.map((item, index) => (
          <ResultHeaderCell
            key={index}
            text={t(`table:results.${item.name}`)}
            tooltip={
              item.tooltipKey ? t(`table:results.${item.tooltipKey}`) : ''
            }
            isShowIcon={index > 0}
            minSize={item.minSize}
            isFormatText={true}
          />
        ))}
      </div>
      <div className="divide-y divide-gray-300">
        {evaluationData?.map((item, i) => {
          const { experimentCount, evaluationCount } = item;
          const conversionRate =
            (Number(experimentCount?.userCount) /
              Number(evaluationCount?.userCount)) *
            100;
          const valuePerUser =
            Number(experimentCount.valueSum) /
            Number(experimentCount.userCount);

          return (
            <div key={i} className="flex items-center w-full">
              <ResultCell
                isFirstItem={true}
                value={item?.variationName || ''}
                minSize={270}
              />
              <ResultCell value={evaluationCount?.userCount} minSize={143} />
              <ResultCell value={experimentCount?.eventCount} minSize={123} />
              <ResultCell value={experimentCount?.userCount} minSize={119} />
              <ResultCell
                value={
                  (isNaN(conversionRate) ? 0 : conversionRate.toFixed(1)) + ' %'
                }
                minSize={147}
              />
              <ResultCell value={experimentCount?.valueSum} minSize={125} />
              <ResultCell
                value={isNaN(valuePerUser) ? '0.00' : valuePerUser.toFixed(2)}
                minSize={123}
              />
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default EvaluationTable;
