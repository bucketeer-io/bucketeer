import {
  IconEditOutlined,
  IconMoreHorizOutlined
} from 'react-icons-material-design';
import { Link } from 'react-router-dom';
import type { ColumnDef } from '@tanstack/react-table';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_PROJECTS } from 'constants/routing';
import { useTranslation } from 'i18n';
import { Project } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import { Popover } from 'components/popover';
import DateTooltip from 'elements/date-tooltip';
import TruncationWithTooltip from 'elements/truncation-with-tooltip';

export const useColumns = ({
  onActionHandler
}: {
  onActionHandler: (value: Project) => void;
}): ColumnDef<Project>[] => {
  const { t } = useTranslation(['common', 'table']);
  const formatDateTime = useFormatDateTime();

  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  return [
    {
      accessorKey: 'name',
      header: `${t('name')}`,
      size: 250,
      cell: ({ row }) => {
        const project = row.original;
        return (
          <Link
            to={`/${currentEnvironment.urlCode}${PAGE_PATH_PROJECTS}/${project.id}`}
            className="underline text-primary-500 typo-para-medium line-clamp-1 break-all"
          >
            {project.name}
          </Link>
        );
      }
    },
    {
      accessorKey: 'creatorEmail',
      header: `${t('maintainer')}`,
      size: 350,
      maxSize: 350,
      cell: ({ row }) => {
        const project = row.original;
        return (
          <TruncationWithTooltip
            elementId={`maintainer-${project.id}`}
            maxSize={350}
            content={project.creatorEmail}
          >
            <div
              id={`maintainer-${project.id}`}
              className="text-gray-700 typo-para-medium w-fit"
            >
              {project.creatorEmail}
            </div>
          </TruncationWithTooltip>
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
          <DateTooltip
            trigger={
              <div className="text-gray-700 typo-para-medium">
                {formatDateTime(project.createdAt)}
              </div>
            }
            date={project.createdAt}
          />
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
