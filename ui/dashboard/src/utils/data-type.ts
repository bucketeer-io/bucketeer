import omit from 'lodash/omit';
import pick from 'lodash/pick';

export const isEmpty = (
  data: number | string | undefined | null | boolean | string[]
) =>
  data === undefined ||
  data === '' ||
  data === null ||
  (Array.isArray(data) && !data.length);

export const isNotEmpty = (
  data: number | string | undefined | null | boolean | string[]
) =>
  (data !== undefined && data !== '' && data !== null) ||
  (Array.isArray(data) && data.length > 0);

export const isEmptyObject = (data: object) => Object.keys(data).length === 0;

export const isNotEmptyObject = (data: object) => Object.keys(data).length > 0;

export const toString = (value: number | null | undefined) => {
  return isNotEmpty(value) ? String(value) : undefined;
};

export const toNumber = (value: string | undefined) => {
  const number = Number(value);
  return Number.isNaN(number) ? 0 : number;
};

export const pickValues = <T extends object>(
  response: T,
  values: object
): Partial<T> => pick(response, Object.keys(values));

export const omitValues = <T extends object>(
  response: T,
  values: object
): Partial<T> => omit(response, Object.keys(values));

type Keys<T> = Array<keyof T>;

type Values<T> = {
  [K in keyof T]: T[K];
}[keyof T][];

type ValuesNonNullable<T> = {
  [K in keyof T]: NonNullable<T[K]>;
}[keyof T][];

export type Entries<T> = {
  [K in keyof T]: [K, T[K]];
}[keyof T][];

export type ListEntries<T> = {
  [K in keyof T]: [number, NonNullable<T[K]>];
}[keyof T][];

export const getObjectKeys = <T extends object>(obj: T) =>
  Object.keys(obj) as Keys<T>;

export const getObjectValues = <T extends object>(obj: T) =>
  Object.values(obj) as Values<T>;

export const getObjectEntries = <T extends object>(obj: T) =>
  Object.entries(obj) as Entries<T>;

export const getListKeys = <T extends object>(obj: T) =>
  getObjectKeys(obj)
    .filter(key => obj[key] !== null)
    .map(key => Number(key));

export const getListValues = <T extends object>(obj: T) =>
  Object.values(obj).filter(item => item !== null) as ValuesNonNullable<T>;

export const getListEntries = <T extends object>(obj: T) =>
  Object.entries(obj)
    .filter(([, item]) => item !== null)
    .map(([id, item]) => [Number(id), item]) as ListEntries<T>;
