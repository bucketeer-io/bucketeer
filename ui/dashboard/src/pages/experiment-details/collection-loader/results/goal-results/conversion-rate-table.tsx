import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { Experiment, GoalResult } from '@types';
import { isNumber } from 'utils/chart';
import { GoalResultState } from '..';
import { ResultHeaderCell, ResultCell } from './goal-results-table-element';
import { DatasetReduceType } from './timeseries-area-line-chart';

const ConversionRateTable = ({
  goalResultState,
  experiment,
  goalResult,
  conversionRateDataSets,
  onToggleShowData
}: {
  goalResultState: GoalResultState;
  experiment: Experiment;
  goalResult: GoalResult;
  conversionRateDataSets: DatasetReduceType[];
  onToggleShowData: (label: string) => void;
}) => {
  const { t } = useTranslation(['common', 'table']);

  const headerList = useMemo(
    () => [
      {
        name: 'variation',
        tooltipKey: '',
        minSize: 270
      },
      goalResultState?.chartType === 'conversion-rate'
        ? {
            name: 'conversion-rate',
            tooltipKey: 'conversion-rate-tooltip',
            minSize: 210
          }
        : {
            name: 'value-user',
            tooltipKey: 'value-user-tooltip',
            minSize: 210
          },
      {
        name: 'improvement',
        tooltipKey: 'improvement-tooltip',
        minSize: 210
      },
      {
        name: 'probability-to-beat-baseline',
        tooltipKey: 'probability-to-beat-baseline-tooltip',
        minSize: 210
      },
      {
        name: 'probability-to-best',
        tooltipKey: 'probability-to-best-tooltip',
        minSize: 210
      }
    ],
    [goalResultState]
  );

  const conversionRateData = useMemo(
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

  const baseVariationResult = useMemo(
    () =>
      goalResult?.variationResults?.find(
        el => el.variationId === experiment.baseVariationId
      ),
    [goalResult, experiment]
  );

  const baseConversionRate = useMemo(() => {
    const experimentUserCount = Number(
      baseVariationResult?.experimentCount?.userCount
    );
    const evaluationUserCount = Number(
      baseVariationResult?.evaluationCount?.userCount
    );
    return evaluationUserCount > 0
      ? (experimentUserCount / evaluationUserCount) * 100
      : 0;
  }, [baseVariationResult]);

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
          const isConversionRateChart =
            goalResultState.chartType === 'conversion-rate';

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

          const isSameVariationId =
            item.variationId === experiment.baseVariationId;

          const improvementValue = isSameVariationId
            ? 'Baseline'
            : (isNumber(conversionRate - baseConversionRate)
                ? conversionRate - baseConversionRate
                : 0
              ).toFixed(1) + ' %';

          const probBeatBaseline = isConversionRateChart
            ? cvrProbBeatBaseline
            : goalValueSumPerUserProbBeatBaseline;

          const probBeatBaselineValue = isSameVariationId
            ? 'Baseline'
            : isNumber(probBeatBaseline?.mean)
              ? (probBeatBaseline.mean * 100).toFixed(1) + ' %'
              : '-';

          const probBest = isConversionRateChart
            ? cvrProbBest
            : goalValueSumPerUserProbBest;

          const probBestValue = isNumber(probBest?.mean)
            ? (probBest.mean * 100).toFixed(1) + ' %'
            : '-';

          const isHidden = conversionRateDataSets.find(
            dataset => dataset.label === item?.variationName
          )?.hidden;

          return (
            <div key={i} className="flex items-center w-full">
              <ResultCell
                variationId={item.variationId}
                isFirstItem
                value={item?.variationName || ''}
                minSize={270}
                currentIndex={i}
                isChecked={!isHidden}
                onToggleShowData={onToggleShowData}
              />
              <ResultCell
                value={
                  isConversionRateChart
                    ? conversionRate.toFixed(1) + ' %'
                    : valuePerUser.toFixed(2)
                }
                minSize={210}
              />
              <ResultCell value={improvementValue} minSize={210} />
              <ResultCell value={probBeatBaselineValue} minSize={210} />
              <ResultCell value={probBestValue} minSize={210} />
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default ConversionRateTable;
