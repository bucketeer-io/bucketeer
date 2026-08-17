const KEY = 'navigation-collapsed';

export const getNavigationCollapsedStorage = (): boolean => {
  try {
    return window.localStorage.getItem(KEY) === 'true';
  } catch (error) {
    console.error(error);
  }
  return false;
};

export const setNavigationCollapsedStorage = (collapsed: boolean): void => {
  try {
    window.localStorage.setItem(KEY, String(collapsed));
  } catch (error) {
    console.error(error);
  }
};
