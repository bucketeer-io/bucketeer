import { FunctionComponent, ReactNode } from 'react';

export type DropdownValue = number | string;

export type DropdownOption = {
  label: ReactNode;
  value: DropdownValue;
  icon?: FunctionComponent;
  iconElement?: ReactNode;
  additionalElement?: ReactNode;
  description?: string;
  tooltip?: ReactNode;
  disabled?: boolean;
  enabled?: boolean;
  labelText?: string;
};

export type DropdownOptionGroup = {
  options: DropdownOption[];
  value?: DropdownValue | DropdownValue[];
  onChange?: (value: DropdownValue | DropdownValue[]) => void;
};
