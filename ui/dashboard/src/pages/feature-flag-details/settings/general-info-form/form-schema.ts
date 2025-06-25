import { FormSchemaProps } from 'hooks/use-form-schema';
import * as yup from 'yup';

export const generalInfoFormSchema = ({ requiredMessage }: FormSchemaProps) =>
  yup.object().shape({
    maintainer: yup.string().required(requiredMessage),
    name: yup.string().required(requiredMessage),
    flagId: yup.string().required(requiredMessage),
    description: yup.string(),
    tags: yup.array().min(1).required(requiredMessage),
    comment: yup.string()
  });

export interface GeneralInfoFormType {
  maintainer: string;
  name: string;
  flagId: string;
  description?: string;
  tags: string[];
  comment?: string;
}
