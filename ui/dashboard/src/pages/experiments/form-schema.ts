import {
  EXPERIMENT_DESCRIPTION_MAX_LENGTH,
  EXPERIMENT_GOAL_MAX_LENGTH,
  EXPERIMENT_GOAL_MIN_LENGTH,
  EXPERIMENT_NAME_MAX_LENGTH,
  EXPERIMENT_START_AT_OLDEST_DAYS
} from 'constants/experiment';
import { requiredMessage } from 'constants/message';
import { i18n } from 'i18n';
import * as yup from 'yup';

const translation = i18n.t;

export const experimentFormSchema = yup.object().shape({
  id: yup.string().max(EXPERIMENT_NAME_MAX_LENGTH),
  name: yup.string().required(requiredMessage),
  baseVariationId: yup.string().required(requiredMessage),
  startType: yup
    .string()
    .oneOf(['manual', 'schedule'])
    .required(requiredMessage),
  startAt: yup
    .string()
    .required(requiredMessage)
    .test(
      'laterThanStartAt',
      translation('message:validation.later-or-equal-days', {
        count: EXPERIMENT_START_AT_OLDEST_DAYS
      }),
      function (value) {
        const startDate = new Date(+value * 1000);
        const d = new Date();
        d.setDate(d.getDate() - EXPERIMENT_START_AT_OLDEST_DAYS);
        return startDate >= d;
      }
    ),
  stopAt: yup
    .string()
    .required(requiredMessage)
    .test('laterThanStartAt', (value, context) => {
      const endDate = new Date(+value * 1000);
      const startAtValue = context?.from && context?.from[0]?.value?.startAt;
      const startDate = new Date(+startAtValue * 1000);
      const startTime = startDate.getTime();
      const endTime = endDate.getTime();
      if (startTime && endTime && endTime <= startTime) {
        return context.createError({
          message: translation('message:validation.stop-later-than-start'),
          path: context.path
        });
      }
      return true;
    })
    .test('lessThanOrEquals30Days', (value, context) => {
      const maxPeriodSeconds = 60 * 60 * 24 * 30;
      const startAtValue = context?.from && context?.from[0]?.value?.startAt;
      const startDate = new Date(+startAtValue * 1000);
      const endDate = new Date(+value * 1000);
      const startTime = startDate.getTime();
      const endTime = endDate.getTime();
      if (endTime / 1000 - startTime / 1000 >= maxPeriodSeconds) {
        return context.createError({
          message: translation('message:validation.less-or-equal-days', {
            count: EXPERIMENT_START_AT_OLDEST_DAYS
          }),
          path: context.path
        });
      }

      return true;
    }),
  description: yup.string().max(EXPERIMENT_DESCRIPTION_MAX_LENGTH),
  audience: yup.mixed(),
  featureId: yup.string().required(requiredMessage),
  goalIds: yup
    .array()
    .min(EXPERIMENT_GOAL_MIN_LENGTH)
    .max(EXPERIMENT_GOAL_MAX_LENGTH)
    .required()
});
