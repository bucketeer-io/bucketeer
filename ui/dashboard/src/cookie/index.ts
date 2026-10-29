import Cookies from 'js-cookie';

const STATE_AVAILABLE_DAY = 1 / 24;

enum Keys {
  STATE = 'state'
}

export const getCookieState = (): string | undefined => {
  return Cookies.get(Keys.STATE);
};

export const setCookieState = (state: string): void => {
  Cookies.set(Keys.STATE, state, {
    expires: STATE_AVAILABLE_DAY,
    // The OAuth provider returns to the callback with a cross-site top-level
    // navigation, which `lax` still allows.
    sameSite: 'lax',
    // `secure` cookies are dropped on plain-http origins, so set it only when
    // the dashboard itself is served over https.
    secure: window.location.protocol === 'https:'
  });
};

export const removeCookieState = (): void => {
  Cookies.remove(Keys.STATE);
};
