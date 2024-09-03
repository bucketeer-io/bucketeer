import { initReactI18next } from 'react-i18next';
import i18n from 'i18next';
import LanguageDetector from 'i18next-browser-languagedetector';
import Backend from 'i18next-http-backend';

const savedLanguage = localStorage.getItem('i18nextLng') || 'en';

i18n
  .use(Backend)
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    debug: true,
    supportedLngs: ['en', 'ja'],
    lng: savedLanguage,
    fallbackLng: 'en',
    ns: 'common',
    defaultNS: 'common',
    interpolation: {
      escapeValue: false
    }
  });

export default i18n;
