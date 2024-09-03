import Cookies from 'js-cookie';

const STATE_AVAILABLE_DAY = 1 / 24;

enum Keys {
  STATE = 'state'
}

export const getCookieState = (): string | undefined => {
  return Cookies.get(Keys.STATE);
};

export const setCookieState = (state: string): void => {
  Cookies.set(Keys.STATE, state, { expires: STATE_AVAILABLE_DAY });
};
