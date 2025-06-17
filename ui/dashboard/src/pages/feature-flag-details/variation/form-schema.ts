import * as yup from 'yup';
import { FeatureVariationType } from '@types';
import { variationsSchema } from 'pages/create-flag/form-schema';

export const variationsFormSchema = yup.object().shape({
  variationType: yup.mixed<FeatureVariationType>().required(),
  onVariation: yup.string(),
  offVariation: yup.string(),
  variations: variationsSchema
});
export type VariationForm = yup.InferType<typeof variationsFormSchema>;
