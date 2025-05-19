import { useMemo } from 'react';
import { useTranslation } from 'i18n';
import { AutoOpsCount, OpsEventRateClause } from '@types';
import { cn } from 'utils/style';
import { IconQuestion } from '@icons';
import { OperationCombinedType } from 'pages/feature-flag-details/operations/types';
import Icon from 'components/icon';
import { Tooltip } from 'components/tooltip';
import { OperationDescription } from '../operation-description';

const PercentItem = ({
  isActive,
  percent,
  className
}: {
  isActive: boolean;
  percent: number;
  className?: string;
}) => {
  return (
    <div className={cn('flex items-center h-[4px]', className)}>
      <div
        className={cn(
          'size-2 rounded-full relative',
          isActive ? 'bg-accent-pink-500' : 'border border-gray-400 bg-gray-50'
        )}
      >
        <span className="absolute -top-8 left-1/2 -translate-x-1/2 typo-head-light-small text-gray-700">
          {percent}%
        </span>
      </div>
    </div>
  );
};

const EventRateProgress = ({
  operation,
  opsCounts
}: {
  operation: OperationCombinedType;
  opsCounts: AutoOpsCount[];
}) => {
  const { t } = useTranslation(['form', 'table']);
  const clause: OpsEventRateClause = useMemo(
    () => (operation.clauses[0]?.clause || {}) as OpsEventRateClause,
    [operation]
  );
  const { goalId, minCount, threadsholdRate } = clause;

  const opsCount = opsCounts.find(
    opsCount => opsCount.autoOpsRuleId === operation.id
  );

  const currentEventRate = useMemo(() => {
    if (opsCount && (+opsCount?.opsEventCount || 0) >= (+minCount || 0)) {
      const { opsEventCount, evaluationCount } = opsCount || {};
      const eventCount =
        +opsEventCount > +evaluationCount ? +evaluationCount : +opsEventCount;
      return Math.round((+eventCount / +evaluationCount) * 100 * 100) / 100;
    }
    return 0;
  }, [opsCount, minCount]);

  const numberOfSteps = useMemo(
    () =>
      Math.round(threadsholdRate * 100) > 10
        ? 10
        : Math.round(threadsholdRate * 100),
    [threadsholdRate]
  );

  const step = useMemo(
    () => (threadsholdRate * 100) / numberOfSteps,
    [threadsholdRate, numberOfSteps]
  );

  const stepArray = useMemo(
    () =>
      Array.from({ length: numberOfSteps }, (_, index) =>
        Math.round(step + index * step)
      ),
    [numberOfSteps, step]
  );

  const barWidth = useMemo(
    () => Math.min((currentEventRate / (threadsholdRate * 100)) * 100, 100),
    [currentEventRate, threadsholdRate]
  );

  return (
    <div className="flex flex-col w-full gap-y-5">
      <div className="flex items-center w-full gap-x-2">
        <OperationDescription
          titleKey={'form:feature-flags.progress-goal-value'}
          value={goalId}
        />
        <OperationDescription
          titleKey={'form:feature-flags.progress-min-count'}
          value={minCount}
        />
        <OperationDescription
          titleKey={'form:feature-flags.progress-current-goal'}
          value={`${currentEventRate}/100 (${currentEventRate}%)`}
          isLastItem
        />
        <Tooltip
          content={t('table:current-event-rate-tooltip')}
          trigger={
            <div className="flex-center size-4">
              <Icon icon={IconQuestion} size={'xxs'} />
            </div>
          }
        />
      </div>

      <div className="bg-gray-100 rounded px-12 pt-14 pb-8 relative">
        <div className="flex h-[4px] bg-gray-200 relative">
          <div
            className="bg-accent-pink-500 absolute h-1 "
            style={{
              width: `${barWidth}%`
            }}
          />
          <PercentItem isActive={currentEventRate > 0} percent={0} />

          {stepArray.map(percentage => (
            <PercentItem
              key={percentage}
              isActive={
                percentage <= currentEventRate && currentEventRate !== 0
              }
              percent={percentage}
              className={'flex justify-end flex-1 items-center h-[4px]'}
            />
          ))}
        </div>
      </div>
    </div>
  );
};

export default EventRateProgress;
