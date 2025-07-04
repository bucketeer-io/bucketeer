const CONSOLE_KEY = 'console-version';
export interface ConsoleVersion {
  version: 'old' | 'new';
  email: string;
}

export const getConsoleVersion = (): ConsoleVersion | undefined => {
  try {
    const consoleVersion = window.localStorage.getItem(CONSOLE_KEY);
    if (consoleVersion) {
      const version = JSON.parse(consoleVersion);
      return version;
    }
  } catch (error) {
    console.error(error);
  }
};

export const setConsoleVersion = (version: ConsoleVersion) => {
  try {
    window.localStorage.setItem(CONSOLE_KEY, JSON.stringify(version));
  } catch (error) {
    console.error(error);
  }
};

export const clearConsoleVersion = () => {
  try {
    window.localStorage.removeItem(CONSOLE_KEY);
  } catch (error) {
    console.log(error);
  }
};
