import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { Experiment, GoalResult } from '@types';
import { cn } from 'utils/style';
import { IconInfo } from '@icons';
import { Polygon } from 'pages/experiment-details/elements/header-details';
import Icon from 'components/icon';

const headerList = [
  {
    name: 'variation',
    tooltip: '',
    minSize: 270
  },
  {
    name: 'evaluation-user',
    tooltip: '',
    minSize: 143
  },
  {
    name: 'goal-total',
    tooltip: '',
    minSize: 119
  },
  {
    name: 'goal-user',
    tooltip: '',
    minSize: 123
  },
  {
    name: 'conversion-rate',
    tooltip: '',
    minSize: 147
  },
  {
    name: 'value-total',
    tooltip: '',
    minSize: 125
  },
  {
    name: 'value-user',
    tooltip: '',
    minSize: 123
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
  const formatText = text.replace(' ', '<br />');
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
          __html: formatText
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
          <HeaderItem
            key={index}
            text={t(`table:results.${item.name}`)}
            isShowIcon={index > 0}
            minSize={item.minSize}
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
              <RowItem
                isFirstItem={true}
                value={item?.variationName || ''}
                minSize={270}
              />
              <RowItem value={evaluationCount?.userCount} minSize={143} />
              <RowItem value={experimentCount?.eventCount} minSize={123} />
              <RowItem value={experimentCount?.userCount} minSize={119} />
              <RowItem
                value={
                  isNaN(conversionRate) ? 'N/A' : conversionRate.toFixed(2)
                }
                minSize={147}
              />
              <RowItem value={experimentCount?.valueSum} minSize={125} />
              <RowItem
                value={isNaN(valuePerUser) ? 'N/A' : valuePerUser.toFixed(2)}
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
