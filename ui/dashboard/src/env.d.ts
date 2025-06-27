interface ImportMetaEnv {
  VITE_ENV: 'dev' | 'prod';
  VITE_WEB_API_ENDPOINT: string;
  VITE_AUTH_REDIRECT_ENDPOINT: string;
  VITE_API_ENDPOINT: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
