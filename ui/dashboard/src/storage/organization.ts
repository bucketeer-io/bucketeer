import { type Nullable } from 'option-t/nullable';

const KEY = 'organizationId';

export const getOrgIdStorage = (): Nullable<string> => {
  try {
    return window.localStorage.getItem(KEY);
  } catch (error) {
    console.error(error);
  }
  return null;
};

export const setOrgIdStorage = (organization: string): void => {
  try {
    window.localStorage.setItem(KEY, organization);
  } catch (error) {
    console.error(error);
  }
};

export const clearOrgIdStorage = (): void => {
  try {
    window.localStorage.removeItem(KEY);
  } catch (error) {
    console.error(error);
  }
};
