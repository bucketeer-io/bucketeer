import { FormSchemaProps } from 'hooks/use-form-schema';
import * as yup from 'yup';
import { SCHEDULE_TYPE_SCHEDULE } from '../../elements/confirm-required-modal/form-schema';

export const generalInfoFormSchema = ({
  requiredMessage,
  translation
}: FormSchemaProps) =>
  yup.object().shape({
    maintainer: yup.string().required(requiredMessage),
    name: yup.string().required(requiredMessage),
    flagId: yup.string().required(requiredMessage),
    description: yup.string(),
    tags: yup.array().min(1).required(requiredMessage),
    comment: yup.string(),
    scheduleType: yup.string(),
    scheduleAt: yup.string().test('validate', function (value, context) {
      const scheduleType = context.from && context.from[0].value.scheduleType;
      if (scheduleType === SCHEDULE_TYPE_SCHEDULE) {
        if (!value)
          return context.createError({
            message: requiredMessage,
            path: context.path
          });
        if (+value * 1000 < new Date().getTime())
          return context.createError({
            message: translation(
              'message:validation.operation.later-than-current-time'
            ),
            path: context.path
          });
      }
      return true;
    })
  });

export interface GeneralInfoFormType {
  maintainer: string;
  name: string;
  flagId: string;
  description?: string;
  tags: string[];
  comment?: string;
  scheduleType?: string;
  scheduleAt?: string;
}
