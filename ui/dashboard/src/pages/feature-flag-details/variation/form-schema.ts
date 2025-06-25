import { FormSchemaProps } from 'hooks/use-form-schema';
import * as yup from 'yup';
import { FeatureVariationType } from '@types';
import {
  createVariationsSchema,
  VariationSchema
} from 'pages/create-flag/form-schema';

export const variationsFormSchema = ({
  requiredMessage,
  translation
}: FormSchemaProps) =>
  yup.object().shape({
    variationType: yup.mixed<FeatureVariationType>().required(requiredMessage),
    onVariation: yup.string(),
    offVariation: yup.string(),
    variations: createVariationsSchema({ requiredMessage, translation })
  });
export interface VariationForm {
  variationType: FeatureVariationType;
  onVariation: string;
  offVariation: string;
  variations: VariationSchema[];
}
