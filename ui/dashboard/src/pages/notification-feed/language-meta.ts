import { FunctionComponent } from 'react';
import { Language } from 'i18n';
import { IconEnglishFlag, IconJapanFlag } from '@icons';

export interface LanguageMeta {
  label: string;
  englishName: string;
  icon: FunctionComponent;
}

export const LANGUAGE_META: Record<string, LanguageMeta> = {
  [Language.ENGLISH]: {
    label: 'English',
    englishName: 'English',
    icon: IconEnglishFlag
  },
  [Language.JAPANESE]: {
    label: '日本語',
    englishName: 'Japanese',
    icon: IconJapanFlag
  }
};

export const FORM_LANGUAGES: Language[] = [Language.ENGLISH, Language.JAPANESE];
