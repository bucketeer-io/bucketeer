import i18n from './i18n';
import { Language } from './types';

export const getLanguage = () => i18n.language as Language;
export const setLanguage = (language: Language) => {
  i18n.changeLanguage(language);
};
