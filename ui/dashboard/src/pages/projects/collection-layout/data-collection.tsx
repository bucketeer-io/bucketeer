import {
  IconEditOutlined,
  IconMoreHorizOutlined
} from 'react-icons-material-design';
import { Link } from 'react-router-dom';
import type { ColumnDef } from '@tanstack/react-table';
import { getCurrentEnvironment, hasEditable, useAuth } from 'auth';
import { PAGE_PATH_PROJECTS } from 'constants/routing';
import { useTranslation } from 'i18n';
import { Project } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import { Popover } from 'components/popover';
import DateTooltip from 'elements/date-tooltip';
import NameWithTooltip from 'elements/name-with-tooltip';

export const useColumns = ({
  onActionHandler
}: {
  onActionHandler: (value: Project) => void;
}): ColumnDef<Project>[] => {
  const { t } = useTranslation(['common', 'table']);
  const formatDateTime = useFormatDateTime();

  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const editable = hasEditable(consoleAccount!);

  return [
    {
      accessorKey: 'name',
      header: `${t('name')}`,
      size: 250,
      cell: ({ row }) => {
        const project = row.original;
        const { id, name } = project;
        return (
          <NameWithTooltip
            id={id}
            content={<NameWithTooltip.Content content={name} id={id} />}
            trigger={
              <Link
                to={`/${currentEnvironment.urlCode}${PAGE_PATH_PROJECTS}/${id}`}
              >
                <NameWithTooltip.Trigger
                  id={id}
                  name={name}
                  maxLines={1}
                  className="min-w-[300px]"
                />
              </Link>
            }
            maxLines={1}
          />
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
        const { id, creatorEmail } = project;
        return (
          <NameWithTooltip
            id={`maintainer-${id}`}
            content={
              <NameWithTooltip.Content
                content={creatorEmail}
                id={`maintainer-${project.id}`}
              />
            }
            trigger={
              <NameWithTooltip.Trigger
                id={`maintainer-${project.id}`}
                name={creatorEmail}
                maxLines={1}
                className="min-w-[300px]"
                haveAction={false}
              />
            }
            maxLines={1}
          />
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
          editable && (
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
          )
        );
      }
    }
  ];
};
