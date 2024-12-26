import { ReactNode } from 'react';
import type { ColumnDef } from '@tanstack/react-table';
import { useTranslation } from 'i18n';

export type TempTableDataType = {
  name: string;
  min: string;
  max: string;
  current: string;
};

export const useColumns = ({
  renderName
}: {
  renderName: (name: TempTableDataType) => ReactNode;
}): ColumnDef<TempTableDataType>[] => {
  const { t } = useTranslation(['table', 'common']);

  return [
    {
      accessorKey: 'name',
      header: `${t('common:name')}`,
      size: 250,
      cell: ({ row }) => {
        const temp = row.original;
        return renderName(temp);
      }
    },
    {
      accessorKey: 'min',
      header: `${t('min')}`,
      size: 200,
      cell: ({ row }) => {
        return (
          <div className="text-gray-700 typo-para-medium">
            {row.original?.min}
          </div>
        );
      }
    },
    {
      accessorKey: 'max',
      header: `${t('max')}`,
      size: 200,
      cell: ({ row }) => {
        return (
          <div className="text-gray-700 typo-para-medium">
            {row.original?.max}
          </div>
        );
      }
    },
    {
      accessorKey: 'current',
      header: `${t('current')}`,
      size: 200,
      cell: ({ row }) => {
        return (
          <div className="text-gray-700 typo-para-medium">
            {row.original?.current}
          </div>
        );
      }
    }
  ];
};
