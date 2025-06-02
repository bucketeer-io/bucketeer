import { requiredMessage } from 'constants/message';
import * as yup from 'yup';

export const addDebuggerFormSchema = yup.object().shape({
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

export type AddDebuggerFormType = yup.InferType<typeof addDebuggerFormSchema>;
