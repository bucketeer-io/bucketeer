import { useCallback, useMemo, useState } from 'react';
import { Trans } from 'react-i18next';
import { useTranslation } from 'i18n';
import { RolloutClause, RuleStrategyVariation } from '@types';
import { OperationCombinedType } from 'pages/feature-flag-details/operations/types';
import {
  getDateTimeDisplay,
  getVariationIndex
} from 'pages/feature-flag-details/operations/utils';
import PercentageBar from 'pages/feature-flag-details/targeting/segment-rule/percentage-bar';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import { Tooltip } from 'components/tooltip';
import VariationLabel from 'elements/variation-label';
import OperationPagination from '../../operation-pagination';
import { OperationDescription } from '../operation-description';
import { ProgressDateTimePoint } from '../progress-date-time-point';

const RolloutProgress = ({
  operation,
  currentAllocationPercentage
}: {
  currentAllocationPercentage: RuleStrategyVariation[];
  operation: OperationCombinedType;
}) => {
  const { t } = useTranslation(['form']);
  const [page, setPage] = useState(0);

  const { clause, createdAt } = operation;
  const rolloutClause = clause as RolloutClause;
  const isActive =
    operation.status === 'RUNNING' || operation.status === 'WAITING';

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

  const alignLeft = useMemo(() => {
    const max = currentAllocationPercentage
      .flatMap(i => i.weight)
      .reduce((a, b) => Math.max(a, b), 0);
    return max.toString().length >= 3
      ? '40px'
      : max.toString().length === 2
        ? '30px'
        : '20px';
  }, [currentAllocationPercentage]);

  // hash code for target variation and control variation, derived from rollout data
  const targetVariationHash = {
    variation: currentAllocationPercentage[0]?.variation ?? 'variation 1',
    index: Math.max(
      getVariationIndex(
        currentAllocationPercentage,
        currentAllocationPercentage[0]?.variation
      ),
      0
    )
  };
  const controlVariationHash = {
    variation: currentAllocationPercentage[1]?.variation ?? 'variation 2',
    index: Math.max(
      getVariationIndex(
        currentAllocationPercentage,
        currentAllocationPercentage[1]?.variation
      ),
      1
    )
  };

  return (
    <div className="flex flex-col w-full gap-y-4">
      <div className="flex items-center gap-x-2">
        <OperationDescription
          titleKey={'form:start-date-value'}
          value={
            getDateTimeDisplay(rolloutClause.schedules[0]?.executeAt)?.date ||
            ''
          }
          isLastItem={!isActive}
        />
        {isActive && (
          <>
            <OperationDescription
              titleKey={'form:target-variation-value'}
              value={
                <VariationLabel
                  label={targetVariationHash.variation}
                  index={targetVariationHash.index}
                  className="max-w-[200px] mt-[2px]"
                />
              }
              isLastItem={false}
            />
            <OperationDescription
              titleKey={'form:control-variation-value'}
              value={
                <VariationLabel
                  label={controlVariationHash.variation}
                  index={controlVariationHash.index}
                  className="max-w-[200px] mt-[2px]"
                />
              }
              isLastItem={true}
            />
          </>
        )}
      </div>
      {!!currentAllocationPercentage.length && isActive && (
        <div className="hover:cursor-pointer typo-para-medium pb-3">
          <p className="text-gray-600 pb-2">
            {t('form:operation.current-variation-title')}
          </p>
          <Tooltip
            content={
              <div>
                {currentAllocationPercentage.map(
                  (item: RuleStrategyVariation, index: number) => (
                    <div
                      key={index}
                      className="flex items-center justify-between mt-[2px] gap-x-2"
                    >
                      <Trans
                        i18nKey={'form:operation.percent-variation'}
                        values={{
                          percent: parseFloat(item.weight.toFixed(3)),
                          variation: item.variation
                        }}
                        components={{
                          p: (
                            <p
                              style={{
                                width: alignLeft,
                                textAlign: 'right'
                              }}
                            />
                          ),
                          comp: (
                            <div className="flex items-center gap-x-2">
                              <FlagVariationPolygon index={index} />
                              <div>{item.variation}</div>
                            </div>
                          )
                        }}
                      />
                    </div>
                  )
                )}
              </div>
            }
            trigger={
              <div className="flex items-center w-full p-0.5 border border-gray-400 rounded-full">
                {currentAllocationPercentage.map(
                  (item: RuleStrategyVariation, index: number) => (
                    <PercentageBar
                      key={index}
                      weight={item.weight}
                      currentIndex={index}
                      isRoundedFull={currentAllocationPercentage.length === 1}
                    />
                  )
                )}
              </div>
            }
          />
        </div>
      )}
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
