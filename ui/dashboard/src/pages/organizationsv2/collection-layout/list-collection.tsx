import { DataTable } from 'elements/data-table';
import { CollectionProps } from '../types';
import { useColumns } from './data-collection';

export const ListCollection = ({
  organizations,
  onSortingChange,
  isLoading
}: CollectionProps) => {
  const columns = useColumns();

  return (
    <DataTable
      isLoading={isLoading}
      data={organizations}
      columns={columns}
      onSortingChange={onSortingChange}
    />
  );
};
