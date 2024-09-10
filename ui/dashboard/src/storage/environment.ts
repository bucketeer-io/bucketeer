import { type Nullable } from 'option-t/nullable';

const KEY = 'environment';

export const getCurrentEnvIdStorage = (): Nullable<string> => {
  try {
    const id = window.localStorage.getItem(KEY);
    if (id) return id;
  } catch (error) {
    console.error(error);
  }
  return null;
};

export const setCurrentEnvIdStorage = (id: string): void => {
  try {
    window.localStorage.setItem(KEY, id);
  } catch (error) {
    console.error(error);
  }
};

export const clearCurrentEnvIdStorage = (): void => {
  try {
    window.localStorage.removeItem(KEY);
  } catch (error) {
    console.error(error);
  }
};
