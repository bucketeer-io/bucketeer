import { yupLocale } from '@/lang/yup';
import * as yup from 'yup';

import {
  ACCOUNT_EMAIL_MAX_LENGTH,
  ACCOUNT_NAME_MAX_LENGTH,
} from '../../constants/account';
import { intl } from '../../lang';
import { messages } from '../../lang/messages';

yup.setLocale(yupLocale);

const nameSchema = yup.string().required().max(ACCOUNT_NAME_MAX_LENGTH);

const emailSchema = yup
  .string()
  .required()
  .email(intl.formatMessage(messages.input.error.invalidEmailAddress))
  .max(ACCOUNT_EMAIL_MAX_LENGTH);

const roleSchema = yup.string().required();

export const addFormSchema = yup.object().shape({
  name: nameSchema,
  email: emailSchema,
  role: roleSchema,
});

export const updateFormSchema = yup.object().shape({
  name: nameSchema,
  role: roleSchema,
});
