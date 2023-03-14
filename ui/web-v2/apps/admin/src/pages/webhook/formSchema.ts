import { yupLocale } from '@/lang/yup';
import * as yup from 'yup';

import { WEBHOOK_NAME_MAX_LENGTH } from '../../constants/webhook';

yup.setLocale(yupLocale);

const nameSchema = yup.string().required().max(WEBHOOK_NAME_MAX_LENGTH);

export const addFormSchema = yup.object().shape({
  name: nameSchema,
});

export const updateFormSchema = yup.object().shape({
  name: nameSchema,
});
