const KEY = 'organizationId';

export const getOrganizationId = () => {
  try {
    return window.localStorage.getItem(KEY);
  } catch (error) {
    // ignore
  }
  return null;
};

export const settOrganizationId = (organizationId): void => {
  try {
    window.localStorage.setItem(KEY, organizationId);
  } catch (error) {
    // ignore
  }
};

export const clearOrganizationId = (): void => {
  try {
    window.localStorage.removeItem(KEY);
  } catch (error) {
    // ignore
  }
};
