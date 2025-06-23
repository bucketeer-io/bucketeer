import { useMemo } from 'react';
import { i18n } from 'i18n';
import { TFunction } from 'i18next';
import * as yup from 'yup';

export interface FormSchemaProps {
  requiredMessage: string;
  translation: TFunction<['translation', ...string[]], undefined>;
}

const useFormSchema = <T extends yup.AnyObject>(
  callback: ({
    requiredMessage,
    translation
  }: FormSchemaProps) => yup.ObjectSchema<T>
) => {
  const translation = i18n.t;
  const requiredMessage = translation('message:required-field');
  const formSchema = useMemo(
    () => callback({ requiredMessage, translation }),
    [callback, requiredMessage, translation]
  );
  return formSchema;
};

export default useFormSchema;
