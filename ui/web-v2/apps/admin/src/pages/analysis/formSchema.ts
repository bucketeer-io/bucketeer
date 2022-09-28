import * as yup from 'yup';

import { ANALYSIS_USER_METADATA_MAX_LENGTH } from '../../constants/analysis';
import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { localJp } from '../../lang/yup/jp';
import { addDays } from '../../utils/date';

yup.setLocale(localJp);

export const formSchema = yup.object().shape({
  endAt: yup
    .date()
    .required()
    .test(
      'notLaterThanCurrentTime',
      intl.formatMessage(messages.input.error.notLaterThanCurrentTime),
      function (value) {
        return value.getTime() <= new Date().getTime();
      }
    )
    .test(
      'laterThanStartAt',
      intl.formatMessage(messages.input.error.notLaterThanStartAt),
      function (value) {
        const { from } = this as any;
        return from[0].value.startAt.getTime() < value.getTime();
      }
    ),
  goalId: yup.string().required(),
  featureId: yup.string(),
  featureVersion: yup.string(),
  reason: yup.string(),
  startAt: yup
    .date()
    .required()
    .test(
      'notLaterThanCurrentTime',
      intl.formatMessage(messages.input.error.notLessThanOrEquals30Days),
      function (value) {
        return value.getTime() >= addDays(new Date(), -30).getTime();
      }
    ),
  userMetadata: yup
    .array()
    .max(ANALYSIS_USER_METADATA_MAX_LENGTH)
    .of(yup.string().required()),
});
