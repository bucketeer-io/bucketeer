// @ts-ignore
import Cookies from 'js-cookie';

const STATE_AVAILABLE_DAY = 1 / 24;

enum Keys {
  STATE = 'state',
}

export const getState = (): string => {
  return Cookies.get(Keys.STATE, '');
};

export const setState = (state: string): void => {
  Cookies.set(Keys.STATE, state, { expires: STATE_AVAILABLE_DAY });
};
