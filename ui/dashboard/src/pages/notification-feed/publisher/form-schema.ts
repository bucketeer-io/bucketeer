import { FormSchemaProps } from 'hooks/use-form-schema';
import * as yup from 'yup';
import { NotificationLocalizationInput } from '../types';

export interface PublishFormValues {
  localizations: NotificationLocalizationInput[];
}

// A notification is one or more localizations, each with a title and content.
export const formSchema = ({ requiredMessage }: FormSchemaProps) =>
  yup.object().shape({
    localizations: yup
      .array()
      .of(
        yup.object().shape({
          language: yup.string().required(),
          title: yup.string().required(requiredMessage),
          content: yup.string().required(requiredMessage),
          tags: yup
            .array()
            .of(
              yup.object().shape({
                name: yup.string().required(),
                color: yup.string().required()
              })
            )
            .required()
        })
      )
      .min(1)
      .required()
  }) as yup.ObjectSchema<PublishFormValues>;
