import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { Experiment, GoalResult } from '@types';
import { GoalResultState } from '..';
import { ResultHeaderCell, ResultCell } from './goal-results-table-element';

const headerList = [
  {
    name: 'variation',
    tooltipKey: '',
    minSize: 270
  },
  {
    name: 'value-user',
    tooltipKey: 'value-user-tooltip',
    minSize: 212
  },
  {
    name: 'improvement',
    tooltipKey: 'improvement-tooltip',
    minSize: 212
  },
  {
    name: 'probability-to-beat-baseline',
    tooltipKey: 'probability-to-beat-baseline-tooltip',
    minSize: 212
  },
  {
    name: 'probability-to-best',
    tooltipKey: 'probability-to-best-tooltip',
    minSize: 212
  }
];

const ConversionRateTable = ({
  goalResultState,
  experiment,
  goalResult
}: {
  goalResultState: GoalResultState;
  experiment: Experiment;
  goalResult: GoalResult;
}) => {
  const { t } = useTranslation(['common', 'table']);
  const conversionRateData = useMemo(
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
  const baseVariationResult = useMemo(
    () =>
      goalResult?.variationResults?.find(
        el => el.variationId === experiment.baseVariationId
      ),
    [goalResult, experiment]
  );

  const baseConversionRate = useMemo(
    () =>
      (Number(baseVariationResult?.experimentCount?.userCount) /
        Number(baseVariationResult?.evaluationCount.userCount)) *
      100,
    [baseVariationResult]
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
          />
        ))}
      </div>
      <div className="divide-y divide-gray-300">
        {conversionRateData?.map((item, i) => {
          const {
            experimentCount,
            evaluationCount,
            cvrProbBeatBaseline,
            cvrProbBest,
            goalValueSumPerUserProbBest,
            goalValueSumPerUserProbBeatBaseline
          } = item;
          const conversionRate =
            (Number(experimentCount?.userCount) /
              Number(evaluationCount?.userCount)) *
            100;
          const valuePerUser =
            Number(experimentCount.valueSum) /
            Number(experimentCount.userCount);

          const isSameVariationId =
            item.variationId === experiment.baseVariationId;
          const improvementValue = conversionRate - baseConversionRate;
          const probBeatBaseline =
            goalResultState.chartType === 'conversion-rate'
              ? cvrProbBeatBaseline
              : goalValueSumPerUserProbBeatBaseline;
          const probBest =
            goalResultState.chartType === 'conversion-rate'
              ? cvrProbBest
              : goalValueSumPerUserProbBest;

          return (
            <div key={i} className="flex items-center w-full">
              <ResultCell
                isFirstItem={true}
                value={item?.variationName || ''}
                minSize={270}
              />
              <ResultCell
                value={isNaN(valuePerUser) ? '0.00' : valuePerUser.toFixed(2)}
                minSize={212}
              />
              <ResultCell
                value={
                  isSameVariationId
                    ? 'Baseline'
                    : (isNaN(improvementValue)
                        ? 0
                        : improvementValue.toFixed(1)) + ' %'
                }
                minSize={212}
              />
              <ResultCell
                value={
                  isSameVariationId
                    ? 'Baseline'
                    : probBeatBaseline
                      ? (probBeatBaseline.mean * 100).toFixed(1) + ' %'
                      : '-'
                }
                minSize={212}
              />
              <ResultCell
                value={probBest ? (probBest.mean * 100).toFixed(1) + ' %' : '-'}
                minSize={212}
              />
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default ConversionRateTable;
