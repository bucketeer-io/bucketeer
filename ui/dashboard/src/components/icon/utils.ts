import type { FunctionComponent } from 'react';
import compact from 'lodash/compact';
import type { Color } from '@types';

export const colorsx = (...inputs: (Color | boolean)[]): Color => {
  return compact(inputs)[0] as Color;
};

export const iconsx = (
  ...inputs: (FunctionComponent | boolean)[]
): FunctionComponent => {
  return compact(inputs)[0] as FunctionComponent;
};
