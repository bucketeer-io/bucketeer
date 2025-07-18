import {
  PAGE_PATH_AUTH_CALLBACK,
  PAGE_PATH_AUTH_DEMO_CALLBACK
} from 'constants/routing';
import resolveConfig from 'tailwindcss/resolveConfig';
import customTailwindConfig from '../../tailwind.config';

export const tailwindConfig = resolveConfig(customTailwindConfig);

const releaseMode = import.meta.env.VITE_RELEASE_CHANNEL;

declare global {
  interface Window {
    env: {
      DEMO_SIGN_IN_ENABLED?: boolean;
      DEMO_SIGN_IN_EMAIL?: string;
      DEMO_SIGN_IN_PASSWORD?: string;
      GOOGLE_TAG_MANAGER_ID?: string;
      API_ENDPOINT?: string;
      OLD_CONSOLE_ENDPOINT?: string;
    };
  }
}

export const urls = {
  WEB_API_ENDPOINT:
    releaseMode !== 'prod' ? import.meta.env.VITE_WEB_API_ENDPOINT : '',
  AUTH_REDIRECT:
    releaseMode !== 'prod'
      ? `${import.meta.env.VITE_AUTH_REDIRECT_ENDPOINT}${PAGE_PATH_AUTH_CALLBACK}`
      : `${window.location.origin}${PAGE_PATH_AUTH_CALLBACK}`,
  AUTH_DEMO_REDIRECT:
    releaseMode !== 'prod'
      ? `${import.meta.env.VITE_AUTH_REDIRECT_ENDPOINT}${PAGE_PATH_AUTH_DEMO_CALLBACK}`
      : `${window.location.origin}${PAGE_PATH_AUTH_DEMO_CALLBACK}`,
  ORIGIN_URL:
    releaseMode !== 'prod'
      ? `${import.meta.env.VITE_AUTH_REDIRECT_ENDPOINT}`
      : `${window.location.origin}`,
  API_ENDPOINT:
    releaseMode !== 'prod'
      ? import.meta.env.VITE_API_ENDPOINT
      : window.env?.API_ENDPOINT,
  OLD_CONSOLE_ENDPOINT:
    releaseMode !== 'prod'
      ? import.meta.env.VITE_OLD_CONSOLE_ENDPOINT
      : `${window.location.origin}/legacy` // TODO: Remove the `/legacy` when the new console is released
};

export const GOOGLE_TAG_MANAGER_ID = window.env?.GOOGLE_TAG_MANAGER_ID || '';

export const DEMO_SIGN_IN_ENABLED =
  releaseMode !== 'prod'
    ? import.meta.env.VITE_DEMO_SIGN_IN_ENABLED
    : window.env?.DEMO_SIGN_IN_ENABLED;

export const DEMO_SIGN_IN_EMAIL =
  releaseMode !== 'prod'
    ? import.meta.env.VITE_DEMO_SIGN_IN_EMAIL
    : window.env?.DEMO_SIGN_IN_EMAIL;

export const DEMO_SIGN_IN_PASSWORD =
  releaseMode !== 'prod'
    ? import.meta.env.VITE_DEMO_SIGN_IN_PASSWORD
    : window.env?.DEMO_SIGN_IN_PASSWORD;
