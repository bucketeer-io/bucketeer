import { FormSchemaProps } from 'hooks/use-form-schema';
import * as yup from 'yup';

export const addDebuggerFormSchema = ({ requiredMessage }: FormSchemaProps) =>
  yup.object().shape({
    flags: yup
      .array()
      .of(yup.string().required(requiredMessage))
      .min(1, requiredMessage)
      .required(requiredMessage),
    userIds: yup
      .array()
      .of(yup.string().required(requiredMessage))
      .min(1, requiredMessage)
      .required(requiredMessage),
    attributes: yup.array().of(
      yup.object().shape({
        key: yup.string(),
        value: yup.string()
      })
    )
  });

export interface AddDebuggerFormType {
  flags: string[];
  userIds: string[];
  attributes?: { key?: string; value?: string }[];
}
