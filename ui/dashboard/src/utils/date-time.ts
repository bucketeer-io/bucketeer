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

export const formatLongDateTime = ({
  value,
  overrideOptions = {
    month: 'long',
    hour: '2-digit',
    minute: '2-digit',
    weekday: 'short',
    hour12: true
  },
  locale = 'en-US'
}: {
  value: string;
  overrideOptions?: Intl.DateTimeFormatOptions;
  locale?: Intl.LocalesArgument;
}) => {
  try {
    const date = new Date(Number(value) * 1000);

    const options: Intl.DateTimeFormatOptions = {
      day: 'numeric',
      year: 'numeric',
      ...overrideOptions
    };
    return new Intl.DateTimeFormat(locale, options).format(date);
  } catch (error) {
    console.error(error);
    return value;
  }
};
