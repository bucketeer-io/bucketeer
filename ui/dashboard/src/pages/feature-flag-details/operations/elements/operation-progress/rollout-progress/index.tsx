import { useCallback, useMemo, useState } from 'react';
import { useTranslation } from 'i18n';
import { RolloutClause } from '@types';
import {
  OperationCombinedType,
  RolloutTypeMap
} from 'pages/feature-flag-details/operations/types';
import { getDateTimeDisplay } from 'pages/feature-flag-details/operations/utils';
import OperationPagination from '../../operation-pagination';
import { OperationDescription } from '../operation-description';
import { ProgressDateTimePoint } from '../progress-date-time-point';

const RolloutProgress = ({
  operation
}: {
  operation: OperationCombinedType;
}) => {
  const { t } = useTranslation(['form']);
  const [page, setPage] = useState(0);

  const { clause, type, createdAt } = operation;
  const rolloutClause = clause as RolloutClause;
  const isTemplate = useMemo(
    () => type === RolloutTypeMap.TEMPLATE_SCHEDULE,
    [type]
  );

  const schedulesList = useMemo(() => rolloutClause.schedules, [rolloutClause]);
  const count = useMemo(
    () => Math.ceil(schedulesList.length / 10),
    [schedulesList]
  );
  const paginatedScheduleList = useMemo(
    () => schedulesList.slice(page * 10, (page + 1) * 10),
    [page]
  );
  const firstSchedule = useMemo(() => {
    return {
      weight: page === 0 ? 0 : schedulesList[page * 10 - 1].weight / 1000,
      executeAt: page === 0 ? createdAt : schedulesList[page * 10 - 1].executeAt
    };
  }, [schedulesList, page, createdAt]);
  const lastItemWithTriggeredAt = useMemo(
    () =>
      [...schedulesList]
        .reverse()
        .find(s => s.triggeredAt && s.triggeredAt !== '0'),
    [schedulesList]
  );

  const handlePageChange = useCallback((page: number) => {
    setPage(page);
  }, []);

  return (
    <div className="flex flex-col w-full gap-y-4">
      <div className="flex items-center gap-x-2">
        {isTemplate && rolloutClause?.increments && (
          <OperationDescription
            titleKey={'form:increments-value'}
            value={`${rolloutClause.increments}%`}
          />
        )}
        <OperationDescription
          titleKey={'form:start-date-value'}
          value={
            getDateTimeDisplay(rolloutClause.schedules[0]?.executeAt)?.date ||
            ''
          }
          isLastItem={!isTemplate}
        />
        {isTemplate && rolloutClause?.interval && (
          <OperationDescription
            titleKey={'form:frequency-value'}
            value={
              rolloutClause.interval
                ? t(`${rolloutClause.interval?.toLowerCase()}`)
                : ''
            }
            isLastItem
            className="[&>p>span]:capitalize"
          />
        )}
      </div>
      <div className="p-12 pb-16 bg-gray-100 rounded-lg">
        <div className="flex relative h-1">
          <ProgressDateTimePoint
            displayLabel={`${firstSchedule.weight || 0}%`}
            displayTime={firstSchedule?.executeAt}
          />
          {paginatedScheduleList.map((item, index) => {
            const isCurrentActive =
              item.triggeredAt !== '0' &&
              (paginatedScheduleList[index + 1]?.triggeredAt === '0' ||
                !paginatedScheduleList[index + 1]);
            return (
              <ProgressDateTimePoint
                key={item.scheduleId}
                className="flex flex-1 justify-end items-center relative"
                displayLabel={`${item.weight / 1000}%`}
                displayTime={item.executeAt}
                conditionDate={
                  lastItemWithTriggeredAt
                    ? new Date(+lastItemWithTriggeredAt.executeAt * 1000)
                    : undefined
                }
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
  );
};

export default RolloutProgress;
