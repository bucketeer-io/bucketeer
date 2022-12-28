import { yupLocale } from '@/lang/yup';
import * as yup from 'yup';

import { PROJECT_ID_MAX_LENGTH } from '../../../constants/project';
import { intl } from '../../../lang';
import { messages } from '../../../lang/messages';

yup.setLocale(yupLocale);

const regex = new RegExp('^[a-zA-Z0-9-]+$');
const idSchema = yup
  .string()
  .required()
  .matches(regex, intl.formatMessage(messages.input.error.invalidId))
  .test(
    'maxLength',
    intl.formatMessage(messages.input.error.maxLength, {
      max: `${PROJECT_ID_MAX_LENGTH}`,
    }),
    function (value) {
      return value.length <= PROJECT_ID_MAX_LENGTH;
    }
  );

export const addFormSchema = yup.object().shape({
  id: idSchema,
});
