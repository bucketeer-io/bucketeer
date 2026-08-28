import { FormSchemaProps } from 'hooks/use-form-schema';
import * as yup from 'yup';
import { FeatureVariationType, VariationValueSchema } from '@types';
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
    variationValueSchema: yup
      .mixed<VariationValueSchema>()
      .nullable()
      .default(null),
    variations: createVariationsSchema({ requiredMessage, translation })
  });
export interface VariationForm {
  variationType: FeatureVariationType;
  onVariation: string;
  offVariation: string;
  variationValueSchema?: VariationValueSchema | null;
  variations: VariationSchema[];
}
