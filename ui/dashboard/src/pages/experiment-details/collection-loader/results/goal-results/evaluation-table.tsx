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
          variationName: variation?.name || variation?.value || ''
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
            Number(evaluationCount?.userCount) > 0
              ? (Number(experimentCount?.userCount) /
                  Number(evaluationCount?.userCount)) *
                100
              : 0;

          const valuePerUser =
            Number(experimentCount.userCount) > 0
              ? Number(experimentCount.valueSum) /
                Number(experimentCount.userCount)
              : 0;

          return (
            <div key={i} className="flex items-center w-full">
              <ResultCell
                currentIndex={i}
                variationId={item.variationId}
                isFirstItem={true}
                value={item?.variationName || ''}
                minSize={270}
              />
              <ResultCell
                value={Number(evaluationCount?.userCount)}
                minSize={143}
              />
              <ResultCell
                value={Number(experimentCount?.eventCount)}
                minSize={123}
              />
              <ResultCell
                value={Number(experimentCount?.userCount)}
                minSize={119}
              />
              <ResultCell
                value={conversionRate.toFixed(1) + ' %'}
                minSize={147}
              />
              <ResultCell
                value={Number(experimentCount?.valueSum)}
                minSize={125}
              />
              <ResultCell value={valuePerUser.toFixed(2)} minSize={123} />
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default EvaluationTable;
