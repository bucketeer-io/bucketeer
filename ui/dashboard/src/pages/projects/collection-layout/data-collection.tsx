import { useCallback, useMemo } from 'react';
import { IconEditOutlined } from 'react-icons-material-design';
import { Link } from 'react-router-dom';
import type { ColumnDef } from '@tanstack/react-table';
import primaryAvatar from 'assets/avatars/primary.svg';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_PROJECTS } from 'constants/routing';
import { useTranslation } from 'i18n';
import { Project } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import { joinName } from 'utils/name';
import { useFetchMembers } from 'pages/members/collection-loader/use-fetch-members';
import { AvatarImage } from 'components/avatar';
import DateTooltip from 'elements/date-tooltip';
import DisabledPopoverTooltip from 'elements/disabled-popover-tooltip';
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
  const { data: accountCollection } = useFetchMembers({
    organizationId: currentEnvironment.organizationId,
    pageSize: 0
  });

  const accounts = useMemo(
    () => accountCollection?.accounts || [],
    [accountCollection]
  );

  const handleGetCurrentAccount = useCallback(
    (email: string) => {
      const account = accounts.find(
        item =>
          item.email === email &&
          item.organizationId === currentEnvironment.organizationId
      );
      return account;
    },
    [accounts]
  );

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
        const { firstName, lastName, name, avatarImageUrl } =
          handleGetCurrentAccount(creatorEmail) || {};
        const accountName = joinName(firstName, lastName) || name;

        return (
          <div className="flex gap-2">
            <AvatarImage
              image={avatarImageUrl || primaryAvatar}
              alt="member-avatar"
            />
            <div className="flex flex-col gap-0.5">
              <NameWithTooltip
                id={creatorEmail}
                content={
                  <NameWithTooltip.Content
                    content={accountName}
                    id={creatorEmail}
                  />
                }
                trigger={
                  <NameWithTooltip.Trigger
                    id={creatorEmail}
                    name={accountName}
                    maxLines={1}
                    className="min-w-[300px]"
                    haveAction={false}
                  />
                }
                maxLines={1}
              />
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
            </div>
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
          <DisabledPopoverTooltip
            isNeedAdminAccess
            options={[
              {
                label: `${t('table:popover.edit-project')}`,
                icon: IconEditOutlined,
                value: 'EDIT_PROJECT'
              }
            ]}
            onClick={() => onActionHandler(project)}
          />
        );
      }
    }
  ];
};
