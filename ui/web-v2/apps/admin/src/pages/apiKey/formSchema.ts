import * as yup from 'yup';

import { APIKEY_NAME_MAX_LENGTH } from '../../constants/apiKey';
import { localJp } from '../../lang/yup/jp';

yup.setLocale(localJp);

const nameSchema = yup.string().required().max(APIKEY_NAME_MAX_LENGTH);

export const addApiKeyFormSchema = yup.object().shape({
  name: nameSchema,
});

export interface AddApiKeyFormSchema
  extends yup.Asserts<typeof addApiKeyFormSchema> {}

export const updateApiKeyFormSchema = yup.object().shape({
  name: nameSchema,
});
export interface UpdateApiKeyFormSchema
  extends yup.Asserts<typeof updateApiKeyFormSchema> {}
