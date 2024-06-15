import { Nullable } from 'option-t/lib/Nullable';

const KEY = 'environment';

let cache: Nullable<string> = null;

export const getCurrentEnvironmentId = (): Nullable<string> => {
  if (cache) {
    return cache;
  }
  let id: string;
  try {
    id = window.localStorage.getItem(KEY);
  } catch (error) {
    return null;
  }
  if (id == 'undefined') {
    return null;
  }
  return id;
};

export const setCurrentEnvironmentId = (id: string): void => {
  try {
    window.localStorage.setItem(KEY, id);
  } catch (error) {
    // ignore
  }
  cache = id;
};

export const clearCurrentEnvironmentId = (): void => {
  cache = null;
  try {
    window.localStorage.removeItem(KEY);
  } catch (error) {
    // ignore
  }
};
