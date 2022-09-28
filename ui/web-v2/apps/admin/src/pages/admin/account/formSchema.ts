import * as yup from 'yup';

import { ACCOUNT_EMAIL_MAX_LENGTH } from '../../../constants/account';
import { intl } from '../../../lang';
import { messages } from '../../../lang/messages';
import { localJp } from '../../../lang/yup/jp';

yup.setLocale(localJp);

const emailSchema = yup
  .string()
  .required()
  .email(intl.formatMessage(messages.input.error.invalidEmailAddress))
  .max(ACCOUNT_EMAIL_MAX_LENGTH);

export const addFormSchema = yup.object().shape({
  email: emailSchema,
});
