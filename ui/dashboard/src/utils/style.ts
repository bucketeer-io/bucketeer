import type { FunctionComponent } from 'react';
import { clsx, type ClassValue } from 'clsx';
import { COLORS } from 'constants/styles';
import compact from 'lodash/compact';
import { twMerge } from 'tailwind-merge';
import type { Color } from '@types';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export const colorsx = (...inputs: (Color | boolean)[]): Color => {
  return compact(inputs)[0] as Color;
};

export const iconsx = (
  ...inputs: (FunctionComponent | boolean)[]
): FunctionComponent => {
  return compact(inputs)[0] as FunctionComponent;
};

export const getVariationColor = (index: number) =>
  COLORS[index % COLORS.length];

export const capitalize = (str: string) => {
  if (!str) return '';
  return str.charAt(0).toUpperCase() + str.slice(1);
};
