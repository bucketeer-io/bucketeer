import { FormSchemaProps } from 'hooks/use-form-schema';
import * as yup from 'yup';
import {
  SEGMENT_MAX_FILE_SIZE,
  SEGMENT_SUPPORTED_FORMATS
} from 'pages/user-segments/constants';

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
      )
  });
