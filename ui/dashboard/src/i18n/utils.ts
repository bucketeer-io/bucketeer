import i18n from './i18n';
import { LanguageTypes } from './types';

export const getLanguage = () => i18n.language as LanguageTypes;
