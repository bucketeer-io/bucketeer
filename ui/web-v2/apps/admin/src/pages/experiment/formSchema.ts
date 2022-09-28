import * as yup from 'yup';

import {
  EXPERIMENT_DESCRIPTION_MAX_LENGTH,
  EXPERIMENT_NAME_MAX_LENGTH,
  EXPERIMENT_GOAL_MIN_LENGTH,
  EXPERIMENT_GOAL_MAX_LENGTH,
  EXPERIMENT_START_AT_OLDEST_DAYS,
} from '../../constants/experiment';
import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { localJp } from '../../lang/yup/jp';

yup.setLocale(localJp);

const nameSchema = yup.string().required().max(EXPERIMENT_NAME_MAX_LENGTH);
const descriptionSchema = yup.string().max(EXPERIMENT_DESCRIPTION_MAX_LENGTH);

export const addFormSchema = yup.object().shape({
  name: nameSchema,
  description: descriptionSchema,
  featureId: yup.string().required(),
  baselineVariation: yup.string().required(),
  goalIds: yup
    .array()
    .min(EXPERIMENT_GOAL_MIN_LENGTH)
    .max(EXPERIMENT_GOAL_MAX_LENGTH)
    .of(yup.string().required()),
  startAt: yup
    .date()
    .required()
    .test(
      'laterThanStartAt',
      intl.formatMessage(messages.input.error.notLaterThanOrEqualDays, {
        days: `${EXPERIMENT_START_AT_OLDEST_DAYS}`,
      }),
      function (value) {
        const d = new Date();
        d.setDate(d.getDate() - EXPERIMENT_START_AT_OLDEST_DAYS);
        return value >= d;
      }
    ),
  stopAt: yup
    .date()
    .required()
    .test(
      'laterThanStartAt',
      intl.formatMessage(messages.input.error.notLaterThanStartAt),
      function (value) {
        const { from } = this as any;
        return from[0].value.startAt.getTime() < value.getTime();
      }
    )
    .test(
      'lessThanOrEquals30Days',
      intl.formatMessage(messages.input.error.notLessThanOrEquals30Days),
      function (value) {
        const { from } = this as any;
        const maxPeriodSeconds = 60 * 60 * 24 * 30;
        return (
          value.getTime() / 1000 - from[0].value.startAt.getTime() / 1000 <=
          maxPeriodSeconds
        );
      }
    ),
});

export const updateFormSchema = yup.object().shape({
  name: nameSchema,
  description: descriptionSchema,
});
