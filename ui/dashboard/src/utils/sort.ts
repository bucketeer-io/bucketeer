import { SortingType } from 'containers/pages';
import { OrderBy } from '@types';

type SortingProps = {
  accessorKey: string;
  sortingKey?: OrderBy;
  sortingState: SortingType;
  setSortingState: (state: SortingType) => void;
  cb?: () => void;
};
export const sortingFn = ({
  accessorKey,
  sortingKey,
  sortingState,
  setSortingState,
  cb
}: SortingProps) => {
  if (cb) cb();
  const { id, orderBy, orderDirection } = sortingState;
  const isDesc = orderDirection === 'DESC';

  if (sortingKey) {
    return setSortingState({
      id: isDesc ? 'default' : accessorKey,
      orderBy: isDesc ? 'DEFAULT' : sortingKey,
      orderDirection: orderBy === 'DEFAULT' ? 'ASC' : isDesc ? 'ASC' : 'DESC'
    });
  }
  if (id === accessorKey) {
    return setSortingState({
      id: isDesc ? 'default' : accessorKey,
      orderBy: isDesc ? 'DEFAULT' : orderBy,
      orderDirection: isDesc ? 'ASC' : 'DESC'
    });
  }

  const orderKey = handleOrderKey(accessorKey);
  setSortingState({
    id: accessorKey,
    orderBy: orderKey,
    orderDirection: 'ASC'
  });
};

const handleOrderKey = (accessorKey: string): OrderBy => {
  switch (accessorKey) {
    case 'id':
    case 'name':
    case 'email':
      return accessorKey.toUpperCase() as OrderBy;

    case 'createdAt':
    case 'updatedAt': {
      const replaceText = accessorKey.replace('At', '');

      return `${replaceText}_at`.toUpperCase() as OrderBy;
    }
    case 'urlCode':
      return 'URL_CODE';
    case 'userCount':
    case 'environmentCount':
    case 'projectCount': {
      const replaceText = accessorKey.replace('Count', '');
      return `${replaceText}_count`.toUpperCase() as OrderBy;
    }
    case 'featureFlagCount':
      return 'FEATURE_COUNT';
    default:
      return 'ID';
  }
};
