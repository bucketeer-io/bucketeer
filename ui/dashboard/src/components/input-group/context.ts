import { createContext, useContext } from 'react';

export type InputGroupValue = {
  addonSlot: 'left' | 'right';
  addonSize: 'sm' | 'md' | 'lg';
};

const inputGroupDefaultContext: InputGroupValue = {
  addonSlot: 'right',
  addonSize: 'md'
};

export const InputGroupContext = createContext<InputGroupValue>(
  inputGroupDefaultContext
);

export const useInputGroupContext = () => {
  const inputGroupContext = useContext(InputGroupContext);
  return inputGroupContext;
};
