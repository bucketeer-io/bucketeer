import { yupLocale } from '../../lang/yup';
import * as yup from 'yup';

import { APIKEY_NAME_MAX_LENGTH } from '../../constants/apiKey';
import { APIKey } from '../../proto/account/api_key_pb';

yup.setLocale(yupLocale);

const nameSchema = yup.string().required().max(APIKEY_NAME_MAX_LENGTH);

export const apiKeyRole = yup
  .mixed<APIKey.RoleMap[keyof APIKey.RoleMap]>()
  .required();
export type ApiKeyRole = yup.InferType<typeof apiKeyRole>;

export const addApiKeyFormSchema = yup.object().shape({
  name: nameSchema,
  role: apiKeyRole
});
export type AddApiKeyForm = yup.InferType<typeof addApiKeyFormSchema>;

export interface AddApiKeyFormSchema
  extends yup.Asserts<typeof addApiKeyFormSchema> {}

export const updateApiKeyFormSchema = yup.object().shape({
  name: nameSchema
});
export type UpdateApiKeyForm = yup.InferType<typeof updateApiKeyFormSchema>;

export interface UpdateApiKeyFormSchema
  extends yup.Asserts<typeof updateApiKeyFormSchema> {}
