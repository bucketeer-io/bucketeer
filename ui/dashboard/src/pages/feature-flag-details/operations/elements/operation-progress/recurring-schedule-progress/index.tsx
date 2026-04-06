import { useCallback, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import dayjs from 'dayjs';
import { AutoOpsRuleClause, DatetimeClause } from '@types';
import { ActionTypeMap, OperationCombinedType } from '../../../types';
import OperationPagination from '../../operation-pagination';
import { OperationDescription } from '../operation-description';
import { ProgressDateTimePoint } from '../progress-date-time-point';

const RecurringScheduleProgress = ({
  operation
}: {
  operation: OperationCombinedType;
}) => {
  const { t } = useTranslation(['form', 'common', 'table']);

  const [page, setPage] = useState(0);

  const { clauses, createdAt } = operation;

  const firstClause = clauses[0]?.clause as DatetimeClause;
  const recurrence = firstClause?.recurrence;

  const frequencyLabel = useMemo(() => {
    if (!recurrence?.frequency) return '';
    const frequencyKeyMap: Record<string, string> = {
      DAILY: 'form:daily',
      WEEKLY: 'form:weekly',
      MONTHLY: 'form:monthly'
    };
    return t(frequencyKeyMap[recurrence.frequency] ?? 'form:unknown');
  }, [recurrence?.frequency, t]);

  const startDate = useMemo(() => {
    if (!recurrence?.startDate || Number(recurrence.startDate) === 0) return '';
    return dayjs(Number(recurrence.startDate) * 1000).format('YYYY/MM/DD');
  }, [recurrence]);

  const completedCycles = useMemo(() => {
    if (!clauses.length) return 0;
    return Math.min(
      ...clauses.map(c => {
        const dc = c.clause as DatetimeClause;
        return dc.executionCount ?? 0;
      })
    );
  }, [clauses]);

  const maxOccurrences = useMemo(() => {
    if (!recurrence?.maxOccurrences || recurrence.maxOccurrences <= 0) return 0;
    return recurrence.maxOccurrences;
  }, [recurrence]);

  const count = useMemo(() => Math.ceil(clauses.length / 10), [clauses]);

  const paginatedClausesList = useMemo(
    () => clauses.slice(page * 10, (page + 1) * 10),
    [clauses, page]
  );

  const stateOptions = [
    {
      label: t('form:experiments.on'),
      value: ActionTypeMap.ENABLE
    },
    {
      label: t('form:experiments.off'),
      value: ActionTypeMap.DISABLE
    }
  ];

  const getDisplayTime = useCallback(
    (clause: AutoOpsRuleClause) => {
      const dc = clause.clause as DatetimeClause;
      if (clause.isRecurring) {
        const thisCount = dc.executionCount ?? 0;
        // When this clause has already advanced past the current cycle
        // (executed more times than another clause in the same rule),
        // show lastExecutedAt to keep the timeline within the same cycle.
        if (
          thisCount > completedCycles &&
          dc.lastExecutedAt &&
          Number(dc.lastExecutedAt) > 0
        ) {
          return dc.lastExecutedAt;
        }
        if (dc.nextExecutionAt && Number(dc.nextExecutionAt) > 0) {
          return dc.nextExecutionAt;
        }
        if (dc.lastExecutedAt && Number(dc.lastExecutedAt) > 0) {
          return dc.lastExecutedAt;
        }
        if (clause.executedAt && Number(clause.executedAt) > 0) {
          return clause.executedAt;
        }
        return createdAt;
      }
      return dc.time;
    },
    [createdAt, completedCycles]
  );

  const currentClause = useMemo(() => clauses[page * 10 - 1], [clauses, page]);

  const { displayTime, displayLabel } = useMemo(() => {
    if (page === 0) {
      return {
        displayTime: createdAt,
        displayLabel: t('table:created-at')
      };
    }
    return {
      displayTime: getDisplayTime(currentClause),
      displayLabel:
        stateOptions.find(o => o.value === currentClause?.actionType)?.label ||
        ''
    };
  }, [createdAt, page, currentClause, t, getDisplayTime, stateOptions]);

  const handlePageChange = (page: number) => {
    setPage(page);
  };

  const executedDisplay = useMemo(() => {
    if (maxOccurrences > 0) return `${completedCycles} / ${maxOccurrences}`;
    return String(completedCycles);
  }, [completedCycles, maxOccurrences]);

  return (
    <div className="flex flex-col w-full gap-y-4">
      <div className="flex items-center flex-wrap gap-x-2 gap-y-1">
        <OperationDescription
          titleKey={'form:frequency-value'}
          value={frequencyLabel}
        />
        {startDate && (
          <OperationDescription
            titleKey={'form:start-date-value'}
            value={startDate}
          />
        )}
        <OperationDescription
          titleKey={'form:feature-flags.executed-value'}
          value={executedDisplay}
          isLastItem
        />
      </div>

      <div>
        <div className="p-12 pb-16 bg-gray-100 rounded-lg">
          <div className="flex relative h-1">
            <ProgressDateTimePoint
              displayLabel={displayLabel}
              displayTime={displayTime}
            />
            {paginatedClausesList.map(scheduleClause => {
              const dc = scheduleClause.clause as DatetimeClause;
              const thisCount = dc.executionCount ?? 0;
              const nextExec = Number(dc.nextExecutionAt || 0);
              const smallestNextExec = Math.min(
                ...paginatedClausesList
                  .map(c =>
                    Number((c.clause as DatetimeClause).nextExecutionAt || 0)
                  )
                  .filter(v => v > 0)
              );
              const isCurrentActive =
                nextExec > 0 &&
                nextExec === smallestNextExec &&
                thisCount === completedCycles;
              const time = getDisplayTime(scheduleClause);
              return (
                <ProgressDateTimePoint
                  key={scheduleClause.id}
                  className="flex flex-1 justify-end items-center relative"
                  displayLabel={
                    stateOptions.find(
                      o => o.value === scheduleClause.actionType
                    )?.label || ''
                  }
                  displayTime={time}
                  isCurrentActive={isCurrentActive}
                />
              );
            })}
          </div>
        </div>
        <OperationPagination
          page={page}
          count={count}
          onPageChange={handlePageChange}
        />
      </div>
    </div>
  );
};

export default RecurringScheduleProgress;
