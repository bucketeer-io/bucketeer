import { NotificationLanguage } from '@types';

interface LanguageItem {
  readonly label: string;
  readonly value: NotificationLanguage;
}

export const languageList: LanguageItem[] = [
  { label: '日本語', value: 'JAPANESE' },
  { label: 'English', value: 'ENGLISH' }
];
