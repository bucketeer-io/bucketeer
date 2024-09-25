import { yupLocale } from '../../lang/yup';
import * as yup from 'yup';

yup.setLocale(yupLocale);

export const addFormSchema = yup.object().shape({
  flag: yup.array().of(yup.string().required()).min(1, 'Required').required(),
  userId: yup.string().required()
});
