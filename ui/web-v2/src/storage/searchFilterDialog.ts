const KEY = 'isSearchFilterDialogShown';

export const getSearchFilterDialogShown = () => {
  try {
    const value = window.localStorage.getItem(KEY);
    return value ? JSON.parse(value) : false;
  } catch (error) {
    // ignore
  }
  return false;
};

export const setSearchFilterDialogShown = (
  isSearchFilterDialogShown: boolean
): void => {
  try {
    window.localStorage.setItem(KEY, JSON.stringify(isSearchFilterDialogShown));
  } catch (error) {
    // ignore
  }
};
