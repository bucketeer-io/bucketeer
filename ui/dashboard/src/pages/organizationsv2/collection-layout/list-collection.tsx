import { ListTableCollection } from './list-table-collection';
import type { CollectionProps } from './types';

export const ListCollection = ({
  organizations,
  onSortingChange,
  isLoading
}: CollectionProps) => {
  return (
    <ListTableCollection
      organizations={organizations}
      isLoading={isLoading}
      onSortingChange={onSortingChange}
    />
  );
};
