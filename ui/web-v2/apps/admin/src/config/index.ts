export const urls = {
  GRPC:
    process.env.RELEASE_CHANNEL !== 'prod'
      ? process.env.NX_DEV_WEB_API_ENDPOINT
      : '',
  AUTH_REDIRECT:
    process.env.RELEASE_CHANNEL !== 'prod'
      ? `${process.env.NX_DEV_AUTH_REDIRECT_ENDPOINT}/auth/callback`
      : `${window.location.origin}/auth/callback`,
};

export const ENABLE_SETTINGS = true;
