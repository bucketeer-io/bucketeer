import resolveConfig from 'tailwindcss/resolveConfig';
import customTailwindConfig from '../../tailwind.config';

export const tailwindConfig = resolveConfig(customTailwindConfig);

const releaseMode = import.meta.env.VITE_RELEASE_CHANNEL;

export const urls = {
  GRPC: releaseMode !== 'prod' ? import.meta.env.VITE_WEB_API_ENDPOINT : '',
  AUTH_REDIRECT:
    releaseMode !== 'prod'
      ? `${import.meta.env.VITE_AUTH_REDIRECT_ENDPOINT}/auth/callback`
      : `${window.location.origin}/auth/callback`
};

export const DEMO_SIGN_IN_ENABLED = import.meta.env.VITE_DEMO_SIGN_IN_ENABLED;
export const DEMO_SIGN_IN_EMAIL = import.meta.env.VITE_DEMO_SIGN_IN_EMAIL;
export const DEMO_SIGN_IN_PASSWORD = import.meta.env.VITE_DEMO_SIGN_IN_PASSWORD;
