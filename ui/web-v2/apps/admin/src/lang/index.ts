import en from '@/assets/lang/en.json';
import ja from '@/assets/lang/ja.json';
import { createIntl, createIntlCache } from 'react-intl';

import { getSelectedLanguage, isLanguageJapanese } from './getSelectedLanguage';

let messages = en;

if (isLanguageJapanese) {
  messages = ja;
}

const locale = getSelectedLanguage();
const defaultLocale = getSelectedLanguage();

const cache = createIntlCache();
export const intl = createIntl(
  {
    locale,
    defaultLocale: defaultLocale,
    messages,
  },
  cache
);
