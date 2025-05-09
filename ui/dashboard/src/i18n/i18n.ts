import { initReactI18next } from 'react-i18next';
import enResources from '@locales/en';
import jaResources from '@locales/ja';
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
    ns: ['common'],
    defaultNS: 'common',
    backend: {
      loadPath: '/v3/src/@locales/{{lng}}/{{ns}}.json'
    },
    interpolation: {
      escapeValue: false
    },
    resources: {
      en: {
        ...enResources
      },
      ja: {
        ...jaResources
      }
    }
  });

export default i18n;
