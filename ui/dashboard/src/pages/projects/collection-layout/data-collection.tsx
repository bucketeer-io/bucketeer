import {
  IconEditOutlined,
  IconMoreHorizOutlined
} from 'react-icons-material-design';
import type { ColumnDef } from '@tanstack/react-table';
import { useTranslation } from 'i18n';
import { Project } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import { Popover } from 'components/popover';

export const useColumns = ({
  onActionHandler
}: {
  onActionHandler: (value: Project) => void;
}): ColumnDef<Project>[] => {
  const { t } = useTranslation(['common', 'table']);
  const formatDateTime = useFormatDateTime();

  return [
    {
      accessorKey: 'name',
      header: `${t('name')}`,
      size: 250,
      cell: ({ row }) => {
        const project = row.original;
        return (
          <div className="underline text-primary-500 typo-para-medium">
            {project.name}
          </div>
        );
      }
    },
    {
      accessorKey: 'creatorEmail',
      header: `${t('maintainer')}`,
      size: 350,
      cell: ({ row }) => {
        const project = row.original;
        return (
          <div className="text-gray-700 typo-para-medium">
            {project.creatorEmail}
          </div>
        );
      }
    },
    {
      accessorKey: 'environmentCount',
      header: `${t('environments')}`,
      size: 120,
      cell: ({ row }) => {
        const project = row.original;
        return (
          <div className="text-gray-700 typo-para-medium">
            {project.environmentCount}
          </div>
        );
      }
    },
    {
      accessorKey: 'featureFlagCount',
      header: t('table:flags'),
      size: 100,
      cell: ({ row }) => {
        const project = row.original;
        return (
          <div className="text-gray-700 typo-para-medium">
            {project.featureFlagCount}
          </div>
        );
      }
    },
    {
      accessorKey: 'createdAt',
      header: `${t('table:created-at')}`,
      size: 160,
      cell: ({ row }) => {
        const project = row.original;
        return (
          <div className="text-gray-700 typo-para-medium">
            {formatDateTime(project.createdAt)}
          </div>
        );
      }
    },
    {
      accessorKey: 'action',
      size: 60,
      header: '',
      meta: {
        align: 'center',
        style: { textAlign: 'center', fitContent: true }
      },
      enableSorting: false,
      cell: ({ row }) => {
        const project = row.original;

        return (
          <Popover
            options={[
              {
                label: `${t('table:popover.edit-project')}`,
                icon: IconEditOutlined,
                value: 'EDIT_PROJECT'
              }
            ]}
            icon={IconMoreHorizOutlined}
            onClick={() => onActionHandler(project)}
            align="end"
          />
        );
      }
    }
  ];
};
