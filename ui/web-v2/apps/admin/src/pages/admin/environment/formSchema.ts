import { yupLocale } from '@/lang/yup';
import * as yup from 'yup';

import {
  ENVIRONMENT_NAME_MAX_LENGTH,
  ENVIRONMENT_URL_CODE_MAX_LENGTH,
} from '../../../constants/environment';
import { intl } from '../../../lang';
import { messages } from '../../../lang/messages';

yup.setLocale(yupLocale);

const nameRegex = /^[a-zA-Z0-9][a-zA-Z0-9\s-]*$/;
const urlCodeRegex = /^[a-zA-Z0-9][a-zA-Z0-9-]*$/;

const nameSchema = yup.string().max(ENVIRONMENT_NAME_MAX_LENGTH).required();

const urlCodeSchema = yup
  .string()
  .required()
  .matches(urlCodeRegex, intl.formatMessage(messages.input.error.invalidUrlCode))
  .test(
    'maxLength',
    intl.formatMessage(messages.input.error.maxLength, {
      max: `${ENVIRONMENT_URL_CODE_MAX_LENGTH}`,
    }),
    function (value) {
      return value.length <= ENVIRONMENT_URL_CODE_MAX_LENGTH;
    }
  );

export const addFormSchema = yup.object().shape({
  name: nameSchema,
  urlCode: urlCodeSchema,
  projectId: yup.string().required(),
});

export const updateFormSchema = yup.object().shape({
  name: nameSchema,
});
