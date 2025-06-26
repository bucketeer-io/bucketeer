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
      SDK_API_ENDPOINT?: string;
    };
  }
}

export const urls = {
  GRPC: releaseMode !== 'prod' ? import.meta.env.VITE_WEB_API_ENDPOINT : '',
  AUTH_REDIRECT:
    releaseMode !== 'prod'
      ? `${import.meta.env.VITE_AUTH_REDIRECT_ENDPOINT}/auth/callback`
      : `${window.location.origin}/v3/auth/callback`, // TODO: Remove the `/v3` when the new console is released,
  ORIGIN_URL:
    releaseMode !== 'prod'
      ? `${import.meta.env.VITE_AUTH_REDIRECT_ENDPOINT}/`
      : `${window.location.origin}/v3/`, // TODO: Remove the `/v3` when the new console is released
  SDK_API_ENDPOINT:
    releaseMode !== 'prod'
      ? import.meta.env.VITE_SDK_API_ENDPOINT
      : window.env?.SDK_API_ENDPOINT
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
