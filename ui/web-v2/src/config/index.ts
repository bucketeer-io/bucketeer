export const urls = {
  GRPC:
    process.env.RELEASE_CHANNEL !== 'prod'
      ? process.env.DEV_WEB_API_ENDPOINT
      : '',
  AUTH_REDIRECT:
    process.env.RELEASE_CHANNEL !== 'prod'
      ? `${process.env.DEV_AUTH_REDIRECT_ENDPOINT}/auth/callback`
      : `${window.location.origin}/legacy/auth/callback`,
  NEW_CONSOLE_ENDPOINT:
    process.env.RELEASE_CHANNEL !== 'prod'
      ? `${process.env.NEW_CONSOLE_ENDPOINT}?fromOldConsole=true`
      : `${window.location.origin}?fromOldConsole=true`
};

export const ENABLE_SETTINGS = true;

declare global {
  interface Window {
    env: {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      [key: string]: any;
    };
  }
}

export const GOOGLE_TAG_MANAGER_ID = window.env?.GOOGLE_TAG_MANAGER_ID || '';

export const DEMO_SIGN_IN_ENABLED =
  process.env.RELEASE_CHANNEL !== 'prod'
    ? process.env.DEMO_SIGN_IN_ENABLED
    : window.env?.DEMO_SIGN_IN_ENABLED;

export const DEMO_SIGN_IN_EMAIL =
  process.env.RELEASE_CHANNEL !== 'prod'
    ? process.env.DEMO_SIGN_IN_EMAIL
    : window.env?.DEMO_SIGN_IN_EMAIL;

export const DEMO_SIGN_IN_PASSWORD =
  process.env.RELEASE_CHANNEL !== 'prod'
    ? process.env.DEMO_SIGN_IN_PASSWORD
    : window.env?.DEMO_SIGN_IN_PASSWORD;
