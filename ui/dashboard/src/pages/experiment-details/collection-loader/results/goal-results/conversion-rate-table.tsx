import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { Experiment, GoalResult } from '@types';
import { cn } from 'utils/style';
import { IconInfo } from '@icons';
import { Polygon } from 'pages/experiment-details/elements/header-details';
import Icon from 'components/icon';
import { GoalResultState } from '..';

const headerList = [
  {
    name: 'variation',
    tooltip: '',
    minSize: 270
  },
  {
    name: 'value-user',
    tooltip: '',
    minSize: 212
  },
  {
    name: 'improvement',
    tooltip: '',
    minSize: 212
  },
  {
    name: 'probability-to-beat-baseline',
    tooltip: '',
    minSize: 212
  },

  {
    name: 'probability-to-best',
    tooltip: '',
    minSize: 212
  }
];

const HeaderItem = ({
  text,
  minSize,
  isShowIcon = true,
  className
}: {
  text: string;
  minSize: number;
  isShowIcon?: boolean;
  className?: string;
}) => {
  return (
    <div
      className={cn(
        'flex items-center size-fit w-full p-4 pt-0 gap-x-3 text-[13px] leading-[13px] text-gray-500 uppercase',
        className
      )}
      style={{
        minWidth: minSize
      }}
    >
      <p
        dangerouslySetInnerHTML={{
          __html: text
        }}
      />
      {isShowIcon && <Icon icon={IconInfo} size={'xxs'} color="gray-500" />}
    </div>
  );
};

const RowItem = ({
  value,
  minSize,
  isFirstItem,
  className
}: {
  value: string | number | boolean;
  minSize: number;
  isFirstItem?: boolean;
  className?: string;
}) => {
  const isBooleanValue = ['true', 'false'].includes(value as string);
  return (
    <div
      className={cn(
        'flex items-center size-fit w-full px-4 py-5 gap-x-2 text-gray-500',
        className
      )}
      style={{ minWidth: minSize }}
    >
      {isFirstItem && isBooleanValue && (
        <Polygon
          className={cn('border-none size-3', {
            'bg-accent-blue-500': value === 'true',
            'bg-accent-pink-500': value === 'false'
          })}
        />
      )}
      <p
        className={cn('typo-para-medium leading-4 text-gray-800', {
          capitalize: isBooleanValue
        })}
      >
        {String(value)}
      </p>
    </div>
  );
};

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
          <HeaderItem
            key={index}
            text={t(`table:results.${item.name}`)}
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
              <RowItem
                isFirstItem={true}
                value={item?.variationName || ''}
                minSize={270}
              />
              <RowItem
                value={isNaN(valuePerUser) ? 'N/A' : valuePerUser.toFixed(2)}
                minSize={212}
              />
              <RowItem
                value={
                  isSameVariationId
                    ? 'Baseline'
                    : isNaN(improvementValue)
                      ? 'N/A'
                      : improvementValue.toFixed(2)
                }
                minSize={212}
              />
              <RowItem
                value={
                  isSameVariationId
                    ? 'Baseline'
                    : probBeatBaseline
                      ? (probBeatBaseline.mean * 100).toFixed(1) + ' %'
                      : 'N/A'
                }
                minSize={212}
              />
              <RowItem
                value={
                  probBest ? (probBest.mean * 100).toFixed(1) + ' %' : 'N/A'
                }
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
