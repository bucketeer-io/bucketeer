import resolveConfig from 'tailwindcss/resolveConfig';
import customTailwindConfig from '../../tailwind.config';

export const tailwindConfig = resolveConfig(customTailwindConfig);

const releaseMode = import.meta.env.VITE_RELEASE_CHANNEL;

declare global {
  interface Window {
    env: {
      DEMO_SIGN_IN_ENABLED?: string;
      DEMO_SIGN_IN_EMAIL?: string;
      VITE_DEMO_SIGN_IN_PASSWORD?: string;
      GOOGLE_TAG_MANAGER_ID?: string;
    };
  }
}

export const urls = {
  GRPC: releaseMode !== 'prod' ? import.meta.env.VITE_WEB_API_ENDPOINT : '',
  AUTH_REDIRECT:
    releaseMode !== 'prod'
      ? `${import.meta.env.VITE_AUTH_REDIRECT_ENDPOINT}/auth/callback`
      : `${window.location.origin}/auth/callback`
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
    ? import.meta.env.VITE_DEMO_SIGN_IN_ENABLED
    : window.env?.VITE_DEMO_SIGN_IN_PASSWORD;
