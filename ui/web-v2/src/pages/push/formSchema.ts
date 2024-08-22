import { PUSH_SUPPORTED_FORMATS } from './../../constants/push';
import { yupLocale } from '../../lang/yup';
import * as yup from 'yup';

import {
  PUSH_NAME_MAX_LENGTH,
  PUSH_TAG_LIST_MIN_LENGTH
} from '../../constants/push';
import { messages } from '../../lang/messages';
import { intl } from '../../lang';

yup.setLocale(yupLocale);

const nameSchema = yup.string().required().max(PUSH_NAME_MAX_LENGTH);

const tagsSchema = yup.array().required().min(PUSH_TAG_LIST_MIN_LENGTH);

const fileSchema = yup
  .mixed()
  .required()
  .test(
    'fileType',
    intl.formatMessage(messages.fileUpload.unsupportedType),
    (value) => {
      return (
        !value ||
        !value[0] ||
        (value[0] && PUSH_SUPPORTED_FORMATS.includes(value[0].type))
      );
    }
  )
  .test(
    'validJson',
    intl.formatMessage(messages.fileUpload.invalidJson),
    async (value) => {
      if (!value || !value[0]) {
        return false;
      }

      return await new Promise<boolean>((resolve) => {
        const reader = new FileReader();

        reader.onload = (e) => {
          try {
            JSON.parse(e.target?.result as string);
            resolve(true); // Valid JSON
          } catch (error) {
            resolve(false); // Invalid JSON
          }
        };

        reader.readAsText(value[0]);
      });
    }
  );

export const addFormSchema = yup.object().shape({
  name: nameSchema,
  tags: tagsSchema,
  file: fileSchema
});

export const updateFormSchema = yup.object().shape({
  name: nameSchema,
  tags: tagsSchema
});
