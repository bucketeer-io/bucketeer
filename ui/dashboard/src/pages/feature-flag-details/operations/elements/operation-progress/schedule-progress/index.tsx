import { useCallback, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { AutoOpsRuleClause, DatetimeClause } from '@types';
import {
  ActionTypeMap,
  OperationCombinedType
} from 'pages/feature-flag-details/operations/types';
import OperationPagination from '../../operation-pagination';
import { ProgressDateTimePoint } from '../progress-date-time-point';

const ScheduleProgress = ({
  operation
}: {
  operation: OperationCombinedType;
}) => {
  const { t } = useTranslation(['form', 'common', 'table']);

  const [page, setPage] = useState(0);

  const { clauses, createdAt, updatedAt } = operation;

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

  const isUpdated = useMemo(
    () => updatedAt > createdAt,
    [updatedAt, createdAt]
  );

  const currentClause = useMemo(() => clauses[page * 10 - 1], [clauses, page]);

  const { displayTime, displayLabel } = useMemo(() => {
    if (page === 0) {
      return {
        displayTime: isUpdated ? updatedAt : createdAt,
        displayLabel: isUpdated ? t('table:updated-at') : t('table:created-at')
      };
    }
    return {
      displayTime: getTimeClause(currentClause),
      displayLabel:
        stateOptions.find(o => o.value === currentClause?.actionType)?.label ||
        ''
    };
  }, [isUpdated, updatedAt, createdAt, page, currentClause]);

  const handlePageChange = (page: number) => {
    setPage(page);
  };

  return (
    <div>
      <div className="p-12 pb-16 bg-gray-100 rounded-lg">
        <div className="flex relative h-1">
          <ProgressDateTimePoint
            displayLabel={displayLabel}
            displayTime={displayTime}
          />
          {paginatedClausesList.map(scheduleClause => {
            const time = getTimeClause(scheduleClause);
            return (
              <ProgressDateTimePoint
                key={scheduleClause.id}
                className="flex flex-1 justify-end items-center relative"
                displayLabel={
                  stateOptions.find(o => o.value === scheduleClause.actionType)
                    ?.label || ''
                }
                displayTime={time}
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
  );
};

export default ScheduleProgress;
