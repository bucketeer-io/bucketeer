import * as yup from 'yup';
import { FeatureVariationType } from '@types';
import { variationsSchema } from 'pages/feature-flags/flags-modal/add-flag-modal/formSchema';

export const variationsFormSchema = yup.object().shape({
  variationType: yup.mixed<FeatureVariationType>().required(),
  onVariation: yup.string(),
  offVariation: yup.string(),
  variations: variationsSchema,
  resetSampling: yup.bool(),
  requireComment: yup.bool(),
  comment: yup.string().when('requireComment', {
    is: (requireComment: boolean) => requireComment,
    then: schema => schema.required()
  })
});
export type VariationForm = yup.InferType<typeof variationsFormSchema>;
