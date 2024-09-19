import { CellValueType, TableSignature } from '@types';

export type SortedType = 'asc' | 'desc' | '';

type SortedProps<T> = {
  data: T[];
  sortedType: SortedType;
  fieldName: string;
};
type CompareProps = {
  a: CellValueType;
  b: CellValueType;
  sortedType: SortedType;
};

const compareFunc = ({ a, b, sortedType }: CompareProps) => {
  if (a instanceof Date || b instanceof Date) return 0;
  if (a < b) return sortedType === 'asc' ? -1 : 1;
  if (a > b) return sortedType === 'asc' ? 1 : -1;
  return 0;
};

const sortedDataFunc = <T extends TableSignature>({
  data,
  sortedType,
  fieldName
}: SortedProps<T>) => {
  const sortedData = data.sort((a, b) => {
    const aValue = a[fieldName];
    const bValue = b[fieldName];
    if (!aValue || !bValue) return 0;
    if (aValue instanceof Date && bValue instanceof Date) {
      const dateA = new Date(aValue).getTime();
      const dateB = new Date(bValue).getTime();

      return compareFunc({ a: dateA, b: dateB, sortedType });
    }
    return compareFunc({ a: aValue, b: bValue, sortedType });
  });
  return sortedData;
};

export { sortedDataFunc };
