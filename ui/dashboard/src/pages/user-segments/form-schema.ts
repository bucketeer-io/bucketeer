import { FormSchemaProps } from 'hooks/use-form-schema';
import * as yup from 'yup';
import { RuleClauseType } from 'pages/feature-flag-details/targeting/types';
import {
  SEGMENT_MAX_FILE_SIZE,
  SEGMENT_SUPPORTED_FORMATS
} from 'pages/user-segments/constants';

const segmentRuleClauseSchema = (requiredMessage: string) =>
  yup.object().shape({
    id: yup.string().required(),
    // Segment rules only support attribute comparisons and date conditions:
    // the SEGMENT and FEATURE_FLAG operators are rejected by the server.
    type: yup
      .string()
      .oneOf([RuleClauseType.COMPARE, RuleClauseType.DATE])
      .required(requiredMessage),
    attribute: yup.string().required(requiredMessage),
    operator: yup.string().required(requiredMessage),
    values: yup
      .array()
      .of(yup.string())
      .min(1, requiredMessage)
      .required(requiredMessage)
  });

export const formSchema = ({ requiredMessage, translation }: FormSchemaProps) =>
  yup.object().shape({
    name: yup.string().required(requiredMessage),
    description: yup.string(),
    id: yup.string(),
    userIds: yup.string(),
    file: yup
      .mixed()
      .nullable()
      .test(
        'fileSize',
        translation('message:max-size-file', {
          size: 2
        }),
        value => {
          return !value || (value as File)?.size <= SEGMENT_MAX_FILE_SIZE;
        }
      )
      .test(
        'fileType',
        translation('message:format-file-not-supported'),
        value => {
          return (
            !value || SEGMENT_SUPPORTED_FORMATS.includes((value as File).type)
          );
        }
      ),
    rules: yup.array().of(
      yup.object().shape({
        id: yup.string().required(),
        clauses: yup
          .array()
          .of(segmentRuleClauseSchema(requiredMessage))
          .min(1, requiredMessage)
          .required(requiredMessage)
      })
    )
  });
