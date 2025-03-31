import { AUTOOPS_MAX_MIN_COUNT } from 'constants/autoops';
import * as yup from 'yup';
import {
  areIntervalsApart,
  isArraySorted,
  isTimestampArraySorted
} from 'utils/function';
import {
  ActionTypeMap,
  IntervalMap,
  OpsTypeMap,
  RolloutTypeMap
} from './types';

const schedulesListSchema = yup.array().of(
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
            );
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

export const operationFormSchema = yup.object().shape({
  opsType: yup.mixed<OpsTypeMap>().required(),
  datetimeClausesList: yup.array().of(
    yup.object().shape({
      id: yup.string(),
      actionType: yup.mixed<ActionTypeMap>().required(),
      time: yup.date()
    })
  ),
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
        time: yup.date().test(
          'isLaterThanNow',

          (value, context) => {
            if (
              value &&
              context.from &&
              context.from[3].value.progressiveRolloutType ===
                'TEMPLATE_SCHEDULE'
            ) {
              return value.getTime() > new Date().getTime();
            }
            return true;
          }
        )
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
