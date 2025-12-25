import { AUTOOPS_MAX_MIN_COUNT } from 'constants/autoops';
import { FormSchemaProps } from 'hooks/use-form-schema';
import * as yup from 'yup';
import { OpsEventRateClauseOperator } from '@types';
import {
  areIntervalsApart,
  hasDuplicateTimestamps,
  isArraySorted,
  isTimestampArraySorted
} from 'utils/function';
import { ActionTypeMap, IntervalMap, RolloutTypeMap } from './types';

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
      targetVariationId: string;
      controlVariationId: string;
      increments: number;
      startDate: Date;
      schedulesList: SchedulesListType[];
      interval: IntervalMap[keyof IntervalMap];
    };
    manual: {
      schedulesList: SchedulesListType[];
      targetVariationId: string;
      controlVariationId: string;
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
        targetVariationId: yup.string().required(),
        controlVariationId: yup.string().required(),
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
        targetVariationId: yup.string().required(requiredMessage),
        controlVariationId: yup.string().required(requiredMessage),
        schedulesList: schedulesListSchema({ requiredMessage, translation })
      })
    })
  });
};
