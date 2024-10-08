import { ListTableCollection } from './list-table-collection';
import type { CollectionProps } from './types';

export const ListCollection = ({
  organizations,
  onSortingChange
}: CollectionProps) => {
  return (
    <ListTableCollection
      organizations={organizations}
      onSortingChange={onSortingChange}
    />
  );
};
