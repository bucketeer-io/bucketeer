import { AUTOOPS_MAX_MIN_COUNT } from 'constants/autoops';
import { i18n } from 'i18n';
import * as yup from 'yup';
import { OpsEventRateClauseOperator } from '@types';
import {
  areIntervalsApart,
  hasDuplicateTimestamps,
  isArraySorted,
  isTimestampArraySorted
} from 'utils/function';
import { ActionTypeMap, IntervalMap, RolloutTypeMap } from './types';

const translation = i18n.t;

const requiredMessage = translation('message:required-field');
const laterThanCurrentTimeMessage = translation(
  'message:validation.operation.later-than-current-time'
);
const increasingOrderMessage = translation(
  'message:validation.operation.date-increasing-order'
);

export const schedulesListSchema = yup
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
        )
        .test('isAscending', (_, context) => {
          if (
            context.from &&
            context.from[3].value.progressiveRolloutType ===
              RolloutTypeMap.MANUAL_SCHEDULE
          ) {
            const isValidWeight = isArraySorted(
              context.from[3].value.progressiveRollout.manual.schedulesList.map(
                (d: yup.AnyObject) => Number(d.weight)
              )
            );
            if (!isValidWeight)
              return context.createError({
                message: translation('message:operation.must-be-increasing', {
                  name: translation('message:operation.the-weight')
                }),
                path: context.path
              });
          }
          return true;
        }),
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
        })
        .test('isAscending', (_, context) => {
          if (
            context.from &&
            context.from[3].value.progressiveRolloutType ===
              RolloutTypeMap.MANUAL_SCHEDULE
          ) {
            const isValidDates = isTimestampArraySorted(
              context.from[3].value.progressiveRollout.manual.schedulesList.map(
                (d: yup.AnyObject) => d.executeAt.getTime()
              )
            ).isSorted;
            if (!isValidDates)
              return context.createError({
                message: increasingOrderMessage,
                path: context.path
              });
          }
          return true;
        })
        .test('timeIntervals', (_, context) => {
          if (
            context.from &&
            context.from[3].value.progressiveRolloutType ===
              RolloutTypeMap.MANUAL_SCHEDULE
          ) {
            const isValidIntervals = areIntervalsApart(
              context.from[3].value.progressiveRollout.manual.schedulesList.map(
                (d: yup.AnyObject) => d.executeAt.getTime()
              ),
              5
            );
            if (!isValidIntervals)
              return context.createError({
                message: translation(
                  'message:validation.operation.schedule-interval'
                ),
                path: context.path
              });
          }
          return true;
        }),
      triggeredAt: yup.string()
    })
  )
  .required(requiredMessage);

export const dateTimeClauseListSchema = yup.object().shape({
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

export type DateTimeClauseListType = yup.InferType<
  typeof dateTimeClauseListSchema
>;

export const eventRateSchema = yup.object().shape({
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

export type EventRateSchemaType = yup.InferType<typeof eventRateSchema>;

export const rolloutSchema = yup.object().shape({
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
      schedulesList: schedulesListSchema,
      interval: yup
        .mixed<IntervalMap[keyof IntervalMap]>()
        .required(requiredMessage)
    }),
    manual: yup.object().shape({
      variationId: yup.string().required(requiredMessage),
      schedulesList: schedulesListSchema
    })
  })
});
export type RolloutSchemaType = yup.InferType<typeof rolloutSchema>;
