import { FormSchemaProps } from 'hooks/use-form-schema';
import * as yup from 'yup';
import { TriggerActionType, TriggerType } from '@types';

export const formSchema = ({ requiredMessage }: FormSchemaProps) =>
  yup.object().shape({
    type: yup.mixed<TriggerType>().required(requiredMessage),
    action: yup.mixed<TriggerActionType>().required(requiredMessage),
    description: yup.string()
  });

export interface CreateTriggerSchema {
  type: TriggerType;
  action: TriggerActionType;
  description?: string;
}
