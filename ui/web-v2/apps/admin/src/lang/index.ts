import { createIntl, createIntlCache } from 'react-intl';

import ja from '../assets/lang/ja.json';

const locale = 'ja';
const defaultLocale = 'ja';
const cache = createIntlCache();
export const intl = createIntl(
  {
    locale,
    defaultLocale: defaultLocale,
    messages: ja,
  },
  cache
);
