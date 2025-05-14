import { i18n } from 'i18n';
import * as yup from 'yup';
import { TriggerActionType, TriggerType } from '@types';

const translation = i18n.t;

const requiredMessage = translation('message.required-field');

export const formSchema = yup.object().shape({
  type: yup.mixed<TriggerType>().required(requiredMessage),
  action: yup.mixed<TriggerActionType>().required(requiredMessage),
  description: yup.string()
});

export type CreateTriggerSchema = yup.InferType<typeof formSchema>;
