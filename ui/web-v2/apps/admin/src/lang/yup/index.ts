import { getSelectedLanguage, LanguageTypes } from '../getSelectedLanguage';

import { localEn } from './en';
import { localJp } from './jp';

export const yupLocale =
  getSelectedLanguage() === LanguageTypes.JAPANESE ? localJp : localEn;
