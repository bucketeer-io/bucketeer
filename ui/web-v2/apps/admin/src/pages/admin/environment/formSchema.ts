import * as yup from 'yup';

import { ENVIRONMENT_ID_MAX_LENGTH } from '../../../constants/environment';
import { intl } from '../../../lang';
import { messages } from '../../../lang/messages';
import { localJp } from '../../../lang/yup/jp';

yup.setLocale(localJp);

const regex = new RegExp('^[a-zA-Z0-9-]+$');
const idSchema = yup
  .string()
  .required()
  .matches(regex, intl.formatMessage(messages.input.error.invalidId))
  .test(
    'maxLength',
    intl.formatMessage(messages.input.error.maxLength, {
      max: `${ENVIRONMENT_ID_MAX_LENGTH}`,
    }),
    function (value) {
      return value.length <= ENVIRONMENT_ID_MAX_LENGTH;
    }
  );

export const addFormSchema = yup.object().shape({
  id: idSchema,
  projectId: yup.string().required(),
});
