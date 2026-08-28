const KEY = 'theme';

export type Theme = 'light' | 'dark';

export const getThemeStorage = (): Theme | null => {
  try {
    const stored = window.localStorage.getItem(KEY);
    if (stored === 'dark' || stored === 'light') return stored;
  } catch (error) {
    console.error(error);
  }
  return null;
};

export const setThemeStorage = (theme: Theme): void => {
  try {
    window.localStorage.setItem(KEY, theme);
  } catch (error) {
    console.error(error);
  }
};
