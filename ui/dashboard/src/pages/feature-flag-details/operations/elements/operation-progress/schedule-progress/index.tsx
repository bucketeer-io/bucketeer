import { useCallback, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import dayjs from 'dayjs';
import { AutoOpsRuleClause, DatetimeClause } from '@types';
import {
  ActionTypeMap,
  OperationCombinedType
} from 'pages/feature-flag-details/operations/types';
import OperationPagination from '../../operation-pagination';
import { OperationDescription } from '../operation-description';
import { ProgressDateTimePoint } from '../progress-date-time-point';

const ScheduleProgress = ({
  operation
}: {
  operation: OperationCombinedType;
}) => {
  const { t } = useTranslation(['form', 'common', 'table']);

  const [page, setPage] = useState(0);

  const { clauses, createdAt } = operation;

  const count = useMemo(() => Math.ceil(clauses.length / 10), [clauses]);

  const paginatedClausesList = useMemo(
    () => clauses.slice(page * 10, (page + 1) * 10),
    [clauses, page]
  );

  const stateOptions = useMemo(
    () => [
      {
        label: t('form:experiments.on'),
        value: ActionTypeMap.ENABLE
      },
      {
        label: t('form:experiments.off'),
        value: ActionTypeMap.DISABLE
      }
    ],
    []
  );

  const getTimeClause = useCallback(
    (dateTimeClause: AutoOpsRuleClause) =>
      (dateTimeClause.clause as DatetimeClause).time,
    []
  );

  const currentClause = useMemo(() => clauses[page * 10 - 1], [clauses, page]);

  const startDate = useMemo(() => {
    if (clauses.length === 0) return '';
    const firstTime = (clauses[0].clause as DatetimeClause).time;
    return dayjs(Number(firstTime) * 1000).format('YYYY/MM/DD');
  }, [clauses]);

  const { displayTime, displayLabel } = useMemo(() => {
    if (page === 0) {
      return {
        displayTime: createdAt,
        displayLabel: t('table:created-at')
      };
    }
    return {
      displayTime: getTimeClause(currentClause),
      displayLabel:
        stateOptions.find(o => o.value === currentClause?.actionType)?.label ||
        ''
    };
  }, [createdAt, page, currentClause, t, getTimeClause, stateOptions]);

  const handlePageChange = (page: number) => {
    setPage(page);
  };

  return (
    <div className="flex flex-col w-full gap-y-4">
      <div className="flex items-center flex-wrap gap-x-2 gap-y-1">
        <OperationDescription
          titleKey={'form:frequency-value'}
          value={t('form:feature-flags.one-time')}
        />
        {startDate && (
          <OperationDescription
            titleKey={'form:start-date-value'}
            value={startDate}
            isLastItem
          />
        )}
      </div>

      <div>
        <div className="p-12 pb-16 bg-gray-100 rounded-lg">
          <div className="flex relative h-1">
            <ProgressDateTimePoint
              displayLabel={displayLabel}
              displayTime={displayTime}
            />
            {paginatedClausesList.map((scheduleClause, index) => {
              const isCurrentActive =
                scheduleClause.executedAt !== '0' &&
                (paginatedClausesList[index + 1]?.executedAt === '0' ||
                  !paginatedClausesList[index + 1]);
              const time = getTimeClause(scheduleClause);
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

export default ScheduleProgress;
