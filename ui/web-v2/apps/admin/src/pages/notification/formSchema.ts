import { yupLocale } from '@/lang/yup';
import * as yup from 'yup';

import {
  NOTIFICATION_NAME_MAX_LENGTH,
  NOTIFICATION_SOURCE_TYPES_MIN_LENGTH,
} from '../../constants/notification';
import { intl } from '../../lang';
import { messages } from '../../lang/messages';

yup.setLocale(yupLocale);

const nameSchema = yup.string().required().max(NOTIFICATION_NAME_MAX_LENGTH);
const sourceTypesSchema = yup
  .array()
  .required()
  .min(
    NOTIFICATION_SOURCE_TYPES_MIN_LENGTH,
    intl.formatMessage(messages.input.error.minSelectOptionLength)
  );
const webhookUrlSchema = yup.string().required().url();

export const addFormSchema = yup.object().shape({
  name: nameSchema,
  sourceTypes: sourceTypesSchema,
  webhookUrl: webhookUrlSchema,
});

export const updateFormSchema = yup.object().shape({
  name: nameSchema,
  sourceTypes: sourceTypesSchema,
});
