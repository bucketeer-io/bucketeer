import { useNavigate } from 'react-router-dom';
import { DataTable } from 'elements/data-table';
import type { CollectionProps } from '../types';
import { useColumns } from './data-table';

export const ListTableCollection = ({ organizations }: CollectionProps) => {
  const navigate = useNavigate();
  const columns = useColumns();

  return (
    <DataTable
      data={organizations}
      columns={columns}
      onRowClick={organization =>
        navigate(`/organization-details/${organization.id}`)
      }
    />
  );
};
