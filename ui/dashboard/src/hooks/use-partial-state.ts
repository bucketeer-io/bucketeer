import { useCallback, useState } from 'react';

export const usePartialState = <T extends object>(initialFilters: T) => {
  const [state, setState] = useState<T>(initialFilters);
  const onChange = useCallback((values: Partial<T>) => {
    setState(old => ({ ...old, ...values }));
  }, []);

  return [state, onChange] as [T, (values: Partial<T>) => void];
};
