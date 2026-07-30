const KEY = 'walkthrough_pending';

export const getWalkthroughPendingStorage = (): boolean => {
  try {
    return window.localStorage.getItem(KEY) === 'true';
  } catch (error) {
    console.error(error);
  }
  return false;
};

export const setWalkthroughPendingStorage = (): void => {
  try {
    window.localStorage.setItem(KEY, 'true');
  } catch (error) {
    console.error(error);
  }
};

export const clearWalkthroughPendingStorage = (): void => {
  try {
    window.localStorage.removeItem(KEY);
  } catch (error) {
    console.error(error);
  }
};
