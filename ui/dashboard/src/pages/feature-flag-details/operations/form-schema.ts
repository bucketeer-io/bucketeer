import { AUTOOPS_MAX_MIN_COUNT } from 'constants/autoops';
import { FormSchemaProps } from 'hooks/use-form-schema';
import * as yup from 'yup';
import { OpsEventRateClauseOperator, RecurrenceFrequency } from '@types';
import {
  areIntervalsApart,
  hasDuplicateTimestamps,
  isArraySorted,
  isTimestampArraySorted
} from 'utils/function';
import {
  ActionTypeMap,
  EndConditionType,
  IntervalMap,
  RolloutTypeMap,
  ScheduleType
} from './types';

export interface SchedulesListType {
  scheduleId: string;
  weight: number;
  executeAt: Date;
  triggeredAt?: string;
}

export const schedulesListSchema = ({
  requiredMessage,
  translation
}: FormSchemaProps) => {
  const laterThanCurrentTimeMessage = translation(
    'message:validation.operation.later-than-current-time'
  );

  return yup
    .array()
    .of(
      yup.object().shape({
        scheduleId: yup.string().required(),
        weight: yup
          .number()
          .transform(value => (isNaN(value) ? undefined : value))
          .required(requiredMessage)
          .min(
            1,
            translation('message:operation.must-be-greater-value', {
              name: translation('message:operation.the-weight'),
              value: 1
            })
          )
          .max(
            100,
            translation('message:operation.must-be-less-value', {
              name: translation('message:operation.the-weight'),
              value: 100
            })
          ),
        executeAt: yup
          .date()
          .required(requiredMessage)
          .test('isLaterThanNow', (value, context) => {
            if (
              value &&
              context.from &&
              context.from[3].value.progressiveRolloutType ===
                RolloutTypeMap.MANUAL_SCHEDULE
            ) {
              const isValidDate = value.getTime() > new Date().getTime();
              if (!isValidDate)
                return context.createError({
                  message: laterThanCurrentTimeMessage,
                  path: context.path
                });
            }
            return true;
          }),

        triggeredAt: yup.string()
      })
    )
    .required(requiredMessage)
    .test(
      'isAscending',
      translation('message:validation.operation.weight-increasing-order'),
      (_, context) => {
        if (
          context.from &&
          context.from[2].value.progressiveRolloutType ===
            RolloutTypeMap.MANUAL_SCHEDULE
        ) {
          return isArraySorted(
            context.from[2].value.progressiveRollout.manual.schedulesList.map(
              (d: yup.AnyObject) => Number(d.weight)
            )
          );
        }
        return true;
      }
    )
    .test(
      'isAscending',
      translation('message:validation.operation.date-increasing-order'),
      (_, context) => {
        if (
          context.from &&
          context.from[2].value.progressiveRolloutType ===
            RolloutTypeMap.MANUAL_SCHEDULE
        ) {
          return isTimestampArraySorted(
            context.from[2].value.progressiveRollout.manual.schedulesList.map(
              (d: yup.AnyObject) => d.executeAt.getTime()
            )
          ).isSorted;
        }
        return true;
      }
    )
    .test(
      'timeIntervals',
      translation('message:validation.operation.schedule-interval'),
      (_, context) => {
        if (
          context.from &&
          context.from[2].value.progressiveRolloutType ===
            RolloutTypeMap.MANUAL_SCHEDULE
        ) {
          const isValidIntervals = areIntervalsApart(
            context.from[2].value.progressiveRollout.manual.schedulesList.map(
              (d: yup.AnyObject) => d.executeAt.getTime()
            ),
            5
          );
          if (!isValidIntervals) return false;
        }

        return true;
      }
    );
};

export interface RecurringClauseItem {
  id?: string;
  actionType: ActionTypeMap;
  time: Date;
  wasExecuted?: boolean;
}

export interface RecurringScheduleFormType {
  startDate: Date;
  frequency: RecurrenceFrequency;
  daysOfWeek: number[];
  dayOfMonth: number;
  endCondition: EndConditionType;
  endDate?: Date;
  maxOccurrences?: number;
  recurringClauses: RecurringClauseItem[];
}

export interface ScheduleOperationFormType {
  scheduleType: ScheduleType;
  datetimeClausesList: {
    id?: string;
    actionType: ActionTypeMap;
    wasPassed?: boolean;
    time: Date;
  }[];
  recurring: RecurringScheduleFormType;
}

export interface DateTimeClauseListType {
  datetimeClausesList: {
    id?: string;
    actionType: ActionTypeMap;
    wasPassed?: boolean;
    time: Date;
  }[];
}

export const dateTimeClauseListSchema = ({
  requiredMessage,
  translation
}: FormSchemaProps) => {
  const laterThanCurrentTimeMessage = translation(
    'message:validation.operation.later-than-current-time'
  );
  const increasingOrderMessage = translation(
    'message:validation.operation.date-increasing-order'
  );
  return yup.object().shape({
    datetimeClausesList: yup
      .array()
      .of(
        yup.object().shape({
          id: yup.string(),
          actionType: yup.mixed<ActionTypeMap>().required(),
          wasPassed: yup.boolean(),
          time: yup
            .date()
            .required(requiredMessage)
            .test('isLaterThanNow', (value, context) => {
              if (
                value &&
                value?.getTime() < new Date().getTime() &&
                context?.from &&
                !context?.from[0]?.value?.wasPassed
              ) {
                return context.createError({
                  message: laterThanCurrentTimeMessage,
                  path: context.path
                });
              }
              return true;
            })
        })
      )
      .required(requiredMessage)
      .test('isAscending', (value, context) => {
        if (value?.length) {
          const hasDuplicate = hasDuplicateTimestamps(
            value?.map(item => item.time?.getTime() ?? 0) || []
          );
          if (hasDuplicate) {
            return context.createError({
              message: translation(
                'message:validation.operation.schedules-same-time'
              ),
              path: context.path
            });
          }
          const sortedState = isTimestampArraySorted(
            value?.map(item => item.time?.getTime() ?? 0) || []
          );

          if (!sortedState.isSorted) {
            return context.createError({
              message: increasingOrderMessage,
              path: context.path
            });
          }
        }
        return true;
      })
  });
};

export const recurringScheduleSchema = ({
  requiredMessage,
  translation
}: FormSchemaProps) => {
  return yup.object().shape({
    scheduleType: yup.mixed<ScheduleType>().required(),
    datetimeClausesList: yup.array().when('scheduleType', {
      is: ScheduleType.ONE_TIME,
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
      then: (schema: yup.AnySchema) =>
        dateTimeClauseListSchema({ requiredMessage, translation }).fields
          .datetimeClausesList as typeof schema,
      otherwise: schema => schema.optional()
    }),
    recurring: yup.object().when('scheduleType', {
      is: ScheduleType.RECURRING,
      then: () =>
        yup.object().shape({
          startDate: yup.date().required(requiredMessage),
          frequency: yup
            .mixed<RecurrenceFrequency>()
            .required(requiredMessage)
            .test(
              'validFrequency',
              translation('message:validation.operation.invalid-frequency'),
              value => ['DAILY', 'WEEKLY', 'MONTHLY'].includes(value as string)
            ),
          daysOfWeek: yup
            .array()
            .of(yup.number().min(0).max(6).required())
            .when('frequency', {
              is: 'WEEKLY',
              then: schema =>
                schema.min(
                  1,
                  translation(
                    'message:validation.operation.days-of-week-required'
                  )
                ),
              otherwise: schema => schema.optional()
            }),
          dayOfMonth: yup.number().when('frequency', {
            is: 'MONTHLY',
            then: schema => schema.min(1).max(31).required(requiredMessage),
            otherwise: schema => schema.optional()
          }),
          endCondition: yup.mixed<EndConditionType>().required(),
          endDate: yup.date().when('endCondition', {
            is: EndConditionType.ON_DATE,
            then: schema =>
              schema
                .required(requiredMessage)
                .test(
                  'isAfterStartDate',
                  translation(
                    'message:validation.operation.end-date-on-after-start-date'
                  ),
                  (value, context) => {
                    const startDate = (context.parent as { startDate?: Date })
                      .startDate;
                    if (value && startDate) {
                      const endDay = new Date(
                        value.getFullYear(),
                        value.getMonth(),
                        value.getDate()
                      );
                      const startDay = new Date(
                        startDate.getFullYear(),
                        startDate.getMonth(),
                        startDate.getDate()
                      );
                      return endDay.getTime() >= startDay.getTime();
                    }
                    return true;
                  }
                ),
            otherwise: schema => schema.optional().nullable()
          }),
          maxOccurrences: yup
            .number()
            .transform(value => (isNaN(value) ? undefined : value))
            .when('endCondition', {
              is: EndConditionType.AFTER,
              then: schema =>
                schema
                  .min(
                    1,
                    translation('message:operation.must-be-greater-value', {
                      name: translation('form:feature-flags.occurrences'),
                      value: 1
                    })
                  )
                  .required(requiredMessage),
              otherwise: schema => schema.optional().nullable()
            }),
          recurringClauses: yup
            .array()
            .of(
              yup.object().shape({
                id: yup.string(),
                actionType: yup.mixed<ActionTypeMap>().required(),
                wasExecuted: yup.boolean(),
                time: yup
                  .date()
                  .required(requiredMessage)
                  .test('isLaterThanNow', (value, context) => {
                    const parent = context.parent as RecurringClauseItem;
                    if (parent.wasExecuted) return true;
                    if (!value) return true;
                    const startDate = (
                      context.from?.[1]?.value as RecurringScheduleFormType
                    )?.startDate;
                    if (!startDate) return true;
                    const now = new Date();
                    const startDay = new Date(
                      startDate.getFullYear(),
                      startDate.getMonth(),
                      startDate.getDate()
                    );
                    const today = new Date(
                      now.getFullYear(),
                      now.getMonth(),
                      now.getDate()
                    );
                    if (startDay.getTime() > today.getTime()) return true;
                    if (
                      startDay.getTime() === today.getTime() &&
                      value.getHours() * 60 + value.getMinutes() <=
                        now.getHours() * 60 + now.getMinutes()
                    ) {
                      return context.createError({
                        message: translation(
                          'message:validation.operation.later-than-current-time'
                        ),
                        path: context.path
                      });
                    }
                    return true;
                  })
              })
            )
            .min(1, requiredMessage)
            .required(requiredMessage)
            .test('noDuplicateOrUnsortedTimes', (value, context) => {
              if (value?.length) {
                const times = value.map(item => {
                  if (!item.time) return 0;
                  return new Date(
                    1970,
                    0,
                    1,
                    item.time.getHours(),
                    item.time.getMinutes(),
                    0,
                    0
                  ).getTime();
                });
                if (hasDuplicateTimestamps(times)) {
                  return context.createError({
                    message: translation(
                      'message:validation.operation.schedules-same-time'
                    ),
                    path: context.path
                  });
                }
                const sortedState = isTimestampArraySorted(times);
                if (!sortedState.isSorted) {
                  return context.createError({
                    message: translation(
                      'message:validation.operation.date-increasing-order'
                    ),
                    path: context.path
                  });
                }
              }
              return true;
            })
        }),
      otherwise: () => yup.object().optional()
    })
  });
};

export interface EventRateSchemaType {
  variationId: string;
  goalId: string;
  minCount: number;
  threadsholdRate: number;
  operator: OpsEventRateClauseOperator;
  actionType: ActionTypeMap;
}

export const eventRateSchema = ({ requiredMessage }: FormSchemaProps) =>
  yup.object().shape({
    variationId: yup.string().required(),
    goalId: yup.string().required(requiredMessage),
    minCount: yup
      .number()
      .transform(value => (isNaN(value) ? undefined : value))
      .required(requiredMessage)
      .min(1)
      .max(AUTOOPS_MAX_MIN_COUNT),
    threadsholdRate: yup
      .number()
      .transform(value => (isNaN(value) ? undefined : value))
      .required(requiredMessage)
      .moreThan(0)
      .max(100),
    operator: yup.mixed<OpsEventRateClauseOperator>().required(),
    actionType: yup.mixed<ActionTypeMap>().required()
  });

export interface RolloutSchemaType {
  progressiveRolloutType: RolloutTypeMap;
  progressiveRollout: {
    template: {
      variationId: string;
      increments: number;
      startDate: Date;
      schedulesList: SchedulesListType[];
      interval: IntervalMap[keyof IntervalMap];
    };
    manual: {
      schedulesList: SchedulesListType[];
      variationId: string;
    };
  };
}

export const rolloutSchema = ({
  requiredMessage,
  translation
}: FormSchemaProps) => {
  const laterThanCurrentTimeMessage = translation(
    'message:validation.operation.later-than-current-time'
  );

  return yup.object().shape({
    progressiveRolloutType: yup.mixed<RolloutTypeMap>().required(),
    progressiveRollout: yup.object().shape({
      template: yup.object().shape({
        variationId: yup.string().required(),
        increments: yup
          .number()
          .transform(value => (isNaN(value) ? undefined : value))
          .required(requiredMessage)
          .min(
            1,
            translation('message:validation.operation.increments-greater-than')
          )
          .max(
            100,
            translation('message:validation.operation.increments-less-than')
          ),
        startDate: yup
          .date()
          .required(requiredMessage)
          .test('isLaterThanNow', (value, context) => {
            if (
              value &&
              context.from &&
              context.from[2].value.progressiveRolloutType ===
                RolloutTypeMap.TEMPLATE_SCHEDULE
            ) {
              if (value.getTime() < new Date().getTime())
                return context.createError({
                  message: laterThanCurrentTimeMessage,
                  path: context.path
                });
            }
            return true;
          }),
        schedulesList: schedulesListSchema({ requiredMessage, translation }),
        interval: yup
          .mixed<IntervalMap[keyof IntervalMap]>()
          .required(requiredMessage)
      }),
      manual: yup.object().shape({
        variationId: yup.string().required(requiredMessage),
        schedulesList: schedulesListSchema({ requiredMessage, translation })
      })
    })
  });
};
