import { useTranslation as i18nextUseTranslation } from 'react-i18next';
import type { Namespace } from './types';

const useTranslation = (ns: Namespace | Namespace[]) => {
  return i18nextUseTranslation<Namespace | Namespace[]>(ns);
};

export default useTranslation;
