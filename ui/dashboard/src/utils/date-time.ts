import { Locale } from 'date-fns';
import { enUS as enCH, ja as jaCH } from 'date-fns/locale';
import { Language, getLanguage } from 'i18n';
import { format } from 'timeago.js';

export type FormatDateTimeOptions = { pattern: DTPattern };

export enum DTPattern {
  /**
   * Format: MM/dd/yyyy HH:mm
   */
  FullDateTime
}

const locales: Record<Language, Locale> = {
  en: enCH,
  ja: jaCH
};

const patterns: Record<DTPattern, Record<Language, string>> = {
  [DTPattern.FullDateTime]: { en: 'dd.MM.yyyy HH:mm', ja: 'dd.MM.yyyy HH:mm' }
};

export const getDateTimePattern = (pattern: DTPattern) => {
  const language = getLanguage();
  return patterns[pattern][language] || '';
};

export type FormatDateTime = (
  value: string,
  overrideOptions?: FormatDateTimeOptions
) => string;

export const useFormatDateTime = () => {
  const formatDateTime = (value: string) => {
    try {
      const date = new Date(Number(value) * 1000);

      return format(date);
    } catch (e) {
      console.error(e);
      return value;
    }
  };

  return formatDateTime;
};

export const getDateTimeLocale = (language: Language) => {
  return locales[language];
};

export const formatLongDateTime = (
  date: Date,
  overrideOptions?: Intl.DateTimeFormatOptions
) => {
  const options: Intl.DateTimeFormatOptions = {
    weekday: 'short',
    month: 'long',
    day: 'numeric',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    hour12: true,
    ...overrideOptions
  };

  return new Intl.DateTimeFormat('en-US', options).format(date);
};
