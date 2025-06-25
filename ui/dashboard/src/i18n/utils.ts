import i18n from './i18n';
import { Language } from './types';

export const getLanguage = () => i18n.language as Language;
export const setLanguage = async (language: Language) => {
  await i18n.changeLanguage(language);
  await i18n.loadLanguages(language);
};
