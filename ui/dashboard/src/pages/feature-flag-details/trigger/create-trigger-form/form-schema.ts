import { requiredMessage } from 'constants/message';
import * as yup from 'yup';
import { TriggerActionType, TriggerType } from '@types';

export const formSchema = yup.object().shape({
  type: yup.mixed<TriggerType>().required(requiredMessage),
  action: yup.mixed<TriggerActionType>().required(requiredMessage),
  description: yup.string()
});

export type CreateTriggerSchema = yup.InferType<typeof formSchema>;
