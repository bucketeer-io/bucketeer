import * as yup from 'yup';

const requiredMessage = 'This field is required.';

export const generalInfoFormSchema = yup.object().shape({
  maintainer: yup.string().required(requiredMessage),
  name: yup.string().required(requiredMessage),
  flagId: yup.string().required(requiredMessage),
  description: yup.string(),
  tags: yup.array().min(1).required(requiredMessage),
  comment: yup.string()
});

export type GeneralInfoFormType = yup.InferType<typeof generalInfoFormSchema>;
