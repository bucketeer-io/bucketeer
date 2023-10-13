import { yupLocale } from '@/lang/yup';
import * as yup from 'yup';

import {
  PROJECT_NAME_MAX_LENGTH,
  PROJECT_URL_CODE_MAX_LENGTH,
} from '../../../constants/project';
import { intl } from '../../../lang';
import { messages } from '../../../lang/messages';

yup.setLocale(yupLocale);

const urlCodeRegex = /^[a-zA-Z0-9][a-zA-Z0-9-]*$/;

const nameSchema = yup.string().max(PROJECT_NAME_MAX_LENGTH).required();
const urlCodeSchema = yup
  .string()
  .required()
  .matches(urlCodeRegex, intl.formatMessage(messages.input.error.invalidUrlCode))
  .test(
    'maxLength',
    intl.formatMessage(messages.input.error.maxLength, {
      max: `${PROJECT_URL_CODE_MAX_LENGTH}`,
    }),
    function (value) {
      return value.length <= PROJECT_URL_CODE_MAX_LENGTH;
    }
  );

export const addFormSchema = yup.object().shape({
  name: nameSchema,
  urlCode: urlCodeSchema
});

export const updateFormSchema = yup.object().shape({
  name: nameSchema,
});
