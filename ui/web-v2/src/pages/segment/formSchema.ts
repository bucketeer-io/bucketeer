import { yupLocale } from '../../lang/yup';
import * as yup from 'yup';

import {
  SEGMENT_DESCRIPTION_MAX_LENGTH,
  SEGMENT_MAX_FILE_SIZE,
  SEGMENT_NAME_MAX_LENGTH,
  SEGMENT_SUPPORTED_FORMATS
} from '../../constants/segment';
import { intl } from '../../lang';
import { messages } from '../../lang/messages';

yup.setLocale(yupLocale);

const nameSchema = yup.string().required().max(SEGMENT_NAME_MAX_LENGTH);
const descriptionSchema = yup.string().max(SEGMENT_DESCRIPTION_MAX_LENGTH);
const fileSchema = yup
  .mixed()
  .nullable()
  .test(
    'fileSize',
    intl.formatMessage(messages.fileUpload.fileMaxSize),
    (value) => {
      return (
        !value ||
        !value[0] ||
        (value[0] && value[0].size <= SEGMENT_MAX_FILE_SIZE)
      );
    }
  )
  .test(
    'fileType',
    intl.formatMessage(messages.fileUpload.unsupportedType),
    (value) => {
      return (
        !value ||
        !value[0] ||
        (value[0] && SEGMENT_SUPPORTED_FORMATS.includes(value[0].type))
      );
    }
  );

export const addFormSchema = yup.object().shape({
  name: nameSchema,
  description: descriptionSchema,
  file: fileSchema
});

export const updateFormSchema = yup.object().shape({
  name: nameSchema,
  description: descriptionSchema,
  file: fileSchema
});
