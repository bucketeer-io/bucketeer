const KEY = 'is_login_first_time';

export const getIsLoginFirstTimeStorage = (): boolean => {
  try {
    const value = window.localStorage.getItem(KEY);
    if (value) {
      return value === 'true';
    }
  } catch (error) {
    console.error(error);
  }
  return false;
};

export const setIsLoginFirstTimeStorage = (value: boolean): void => {
  try {
    window.localStorage.setItem(KEY, JSON.stringify(value));
  } catch (error) {
    console.error(error);
  }
};

export const clearIsLoginFirstTimeStorage = (): void => {
  try {
    window.localStorage.removeItem(KEY);
  } catch (error) {
    console.error(error);
  }
};
