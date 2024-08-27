import resolveConfig from 'tailwindcss/resolveConfig';
import customTailwindConfig from '../../tailwind.config';

export const tailwindConfig = resolveConfig(customTailwindConfig);

export const urls = {
  GRPC:
    import.meta.env.VITE_RELEASE_CHANNEL !== 'prod'
      ? import.meta.env.VITE_WEB_API_ENDPOINT
      : '',
  AUTH_REDIRECT:
    import.meta.env.VITE_RELEASE_CHANNEL !== 'prod'
      ? `${import.meta.env.VITE_AUTH_REDIRECT_ENDPOINT}/auth/callback`
      : `${window.location.origin}/auth/callback`
};
