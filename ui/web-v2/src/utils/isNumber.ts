export const isNumber = (value: unknown): boolean => {
  return typeof value === 'number' && Number.isFinite(value);
};
