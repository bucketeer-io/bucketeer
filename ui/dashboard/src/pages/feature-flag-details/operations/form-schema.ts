import { AUTOOPS_MAX_MIN_COUNT } from 'constants/autoops';
import * as yup from 'yup';
import {
  areIntervalsApart,
  hasDuplicateTimestamps,
  isArraySorted,
  isTimestampArraySorted
} from 'utils/function';
import {
  ActionTypeMap,
  IntervalMap,
  OpsTypeMap,
  RolloutTypeMap
} from './types';

const requiredMessage = 'This field is required.';

export const schedulesListSchema = yup.array().of(
  yup.object().shape({
    weight: yup
      .number()
      .transform(value => (isNaN(value) ? undefined : value))
      .required()
      .min(1)
      .max(100)
      .test('isAscending', (_, context) => {
        if (
          context.from &&
          context.from[3].value.progressiveRolloutType ===
            RolloutTypeMap.MANUAL_SCHEDULE
        ) {
          return isArraySorted(
            context.from[3].value.progressiveRollout.manual.schedulesList.map(
              (d: yup.AnyObject) => Number(d.weight)
            )
          );
        }
        return true;
      }),
    executeAt: yup.object().shape({
      time: yup
        .date()
        .test('isLaterThanNow', (value, context) => {
          if (
            value &&
            context.from &&
            context.from[4].value.progressiveRolloutType ===
              RolloutTypeMap.MANUAL_SCHEDULE
          ) {
            return value.getTime() > new Date().getTime();
          }
          return true;
        })
        .test('isAscending', (_, context) => {
          if (
            context.from &&
            context.from[4].value.progressiveRolloutType ===
              RolloutTypeMap.MANUAL_SCHEDULE
          ) {
            return isTimestampArraySorted(
              context.from[4].value.progressiveRollout.manual.schedulesList.map(
                (d: yup.AnyObject) => d.executeAt.time.getTime()
              )
            ).isSorted;
          }
          return true;
        })
        .test('timeIntervals', (_, context) => {
          if (
            context.from &&
            context.from[4].value.progressiveRolloutType ===
              RolloutTypeMap.MANUAL_SCHEDULE
          ) {
            return areIntervalsApart(
              context.from[4].value.progressiveRollout.manual.schedulesList.map(
                (d: yup.AnyObject) => d.executeAt.time.getTime()
              ),
              5
            );
          }
          return true;
        })
    })
  })
);

export const dateTimeClauseListSchema = yup.object().shape({
  datetimeClausesList: yup
    .array()
    .of(
      yup.object().shape({
        id: yup.string(),
        actionType: yup.mixed<ActionTypeMap>().required(),
        time: yup
          .date()
          .required(requiredMessage)
          .test('isLaterThanNow', (value, context) => {
            if (value && value?.getTime() < new Date().getTime()) {
              return context.createError({
                message: 'This must be later than the current time.',
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
            message:
              'You cannot have multiple schedules at the same date and time.',
            path: context.path
          });
        }
        const isSorted = isTimestampArraySorted(
          value?.map(item => item.time?.getTime() ?? 0) || []
        );

        if (!isSorted) {
          return context.createError({
            message: 'The date must be in increasing order.',
            path: context.path
          });
        }
      }
      return true;
    })
});

export type DateTimeClauseListType = yup.InferType<
  typeof dateTimeClauseListSchema
>;

export const eventRateSchema = yup.object().shape({
  variation: yup.string(),
  goal: yup.string().required(requiredMessage),
  minCount: yup
    .number()
    .required(requiredMessage)
    .min(1)
    .max(AUTOOPS_MAX_MIN_COUNT),
  threadsholdRate: yup.number().required(requiredMessage).moreThan(0).max(100),
  operator: yup.string().required(),
  actionType: yup.mixed<ActionTypeMap>().required()
});

export type EventRateSchemaType = yup.InferType<typeof eventRateSchema>;

export const operationFormSchema = yup.object().shape({
  opsType: yup.mixed<OpsTypeMap>().required(),
  datetimeClausesList: yup
    .array()
    .of(
      yup.object().shape({
        id: yup.string(),
        actionType: yup.mixed<ActionTypeMap>().required(),
        time: yup
          .date()
          .test('isLaterThanNow', (value, context) => {
            if (value && value?.getTime() < new Date().getTime()) {
              return context.createError({
                message: 'This must be later than the current time.',
                path: context.path
              });
            }
            return true;
          })
          .required(requiredMessage)
      })
    )
    .test('isAscending', (value, context) => {
      if (value?.length) {
        const hasDuplicate = hasDuplicateTimestamps(
          value?.map(item => item.time?.getTime() ?? 0) || []
        );
        if (hasDuplicate) {
          return context.createError({
            message:
              'You cannot have multiple schedules at the same date and time.',
            path: context.path
          });
        }
        const isSorted = isTimestampArraySorted(
          value?.map(item => item.time?.getTime() ?? 0) || []
        );

        if (!isSorted) {
          return context.createError({
            message: 'The date must be in increasing order.',
            path: context.path
          });
        }
      }
      return true;
    }),
  eventRate: yup.object().shape({
    variation: yup.string(),
    goal: yup
      .string()
      .nullable()
      .test('required', (value, context) => {
        if (
          context.from &&
          context.from[1].value.opsType === OpsTypeMap.EVENT_RATE
        ) {
          return value != null;
        }
        return true;
      }),
    minCount: yup
      .number()
      .transform(value => (isNaN(value) ? undefined : value))
      .required()
      .min(1)
      .max(AUTOOPS_MAX_MIN_COUNT),
    threadsholdRate: yup
      .number()
      .transform(value => (isNaN(value) ? undefined : value))
      .required()
      .moreThan(0)
      .max(100),
    operator: yup.string()
  }),
  progressiveRolloutType: yup.mixed<RolloutTypeMap>().required(),
  progressiveRollout: yup.object().shape({
    template: yup.object().shape({
      variationId: yup.string().required(),
      increments: yup
        .number()
        .transform(value => (isNaN(value) ? undefined : value))
        .required()
        .min(1)
        .max(100),
      datetime: yup.object().shape({
        time: yup.date().test('isLaterThanNow', (value, context) => {
          if (
            value &&
            context.from &&
            context.from[3].value.progressiveRolloutType === 'TEMPLATE_SCHEDULE'
          ) {
            return value.getTime() > new Date().getTime();
          }
          return true;
        })
      }),
      schedulesList: schedulesListSchema,
      interval: yup.mixed<IntervalMap[keyof IntervalMap]>().required()
    }),
    manual: yup.object().shape({
      variationId: yup.string().required(),
      schedulesList: schedulesListSchema
    })
  })
});
export type OperationForm = yup.InferType<typeof operationFormSchema>;
