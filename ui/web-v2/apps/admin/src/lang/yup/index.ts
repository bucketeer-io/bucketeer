import { isLanguageJapanese } from '../getSelectedLanguage';

import { localEn } from './en';
import { localJp } from './jp';

export const yupLocale = isLanguageJapanese ? localJp : localEn;
