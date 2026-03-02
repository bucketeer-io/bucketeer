import { i18n } from 'i18n';
import * as yup from 'yup';

const translation = i18n.t;
const requiredMessage = translation('message:required-field');

export const SCHEDULE_TYPE_UPDATE_NOW = 'UPDATE_NOW';
export const SCHEDULE_TYPE_SCHEDULE = 'SCHEDULE';

export const formSchema = yup.object().shape({
  requireComment: yup.boolean(),
  resetSampling: yup.boolean(),
  comment: yup.string().when('requireComment', {
    is: (requireComment: boolean) => requireComment,
    then: schema => schema.required(requiredMessage)
  }),
  scheduleType: yup.string(),
  scheduleAt: yup.string().test('validate', function (value, context) {
    const scheduleType = context.from && context.from[0].value.scheduleType;
    if (scheduleType === SCHEDULE_TYPE_SCHEDULE) {
      if (!value)
        return context.createError({
          message: requiredMessage,
          path: context.path
        });
      if (+value * 1000 < new Date().getTime())
        return context.createError({
          message: translation(
            'message:validation.operation.later-than-current-time'
          ),
          path: context.path
        });
    }
    return true;
  })
});
