export const isEmpty = (data: number | string | undefined | null) =>
  data === undefined || data === '' || data === null;

export const isNotEmpty = (data: number | string | undefined | null) =>
  data !== undefined && data !== '' && data !== null;

export const isEmptyObject = (data: object) => Object.keys(data).length === 0;

export const isNotEmptyObject = (data: object) => Object.keys(data).length > 0;

export const toString = (value: number | null | undefined) => {
  return isNotEmpty(value) ? String(value) : undefined;
};

export const toNumber = (value: string | undefined) => {
  const number = Number(value);
  return Number.isNaN(number) ? 0 : number;
};
