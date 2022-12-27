import { yupLocale } from '@/lang/yup';
import * as yup from 'yup';

import {
  GOAL_ID_MAX_LENGTH,
  GOAL_NAME_MAX_LENGTH,
  GOAL_DESCRIPTION_MAX_LENGTH,
} from '../../constants/goal';

yup.setLocale(yupLocale);

const idSchema = yup.string().required().max(GOAL_ID_MAX_LENGTH);
const nameSchema = yup.string().required().max(GOAL_NAME_MAX_LENGTH);
const descriptionSchema = yup.string().max(GOAL_DESCRIPTION_MAX_LENGTH);

export const addFormSchema = yup.object().shape({
  id: idSchema,
  name: nameSchema,
  description: descriptionSchema,
});

export const updateFormSchema = yup.object().shape({
  name: nameSchema,
  description: descriptionSchema,
});
