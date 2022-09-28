import * as yup from 'yup';

import {
  PUSH_NAME_MAX_LENGTH,
  PUSH_FCM_API_KEY_MAX_LENGTH,
  PUSH_TAG_LIST_MIN_LENGTH,
} from '../../constants/push';
import { localJp } from '../../lang/yup/jp';

yup.setLocale(localJp);

const nameSchema = yup.string().required().max(PUSH_NAME_MAX_LENGTH);
const fcmApiKeySchema = yup
  .string()
  .required()
  .max(PUSH_FCM_API_KEY_MAX_LENGTH);
const tagsSchema = yup.array().required().min(PUSH_TAG_LIST_MIN_LENGTH);

export const addFormSchema = yup.object().shape({
  name: nameSchema,
  fcmApiKey: fcmApiKeySchema,
  tags: tagsSchema,
});

export const updateFormSchema = yup.object().shape({
  name: nameSchema,
  fcmApiKey: fcmApiKeySchema,
  tags: tagsSchema,
});
