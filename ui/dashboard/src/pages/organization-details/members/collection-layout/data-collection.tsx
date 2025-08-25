import { IconMoreHorizOutlined } from 'react-icons-material-design';
import type { ColumnDef } from '@tanstack/react-table';
import primaryAvatar from 'assets/avatars/primary.svg';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useTranslation } from 'i18n';
import { Account } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import { joinName } from 'utils/name';
import { useFetchTags } from 'pages/members/collection-loader';
import { AvatarImage } from 'components/avatar';
import Icon from 'components/icon';
import ExpandableTag from 'elements/expandable-tag';
import NameWithTooltip from 'elements/name-with-tooltip';

export const useColumns = ({
  onActions
}: {
  onActions: (value: Account) => void;
}): ColumnDef<Account>[] => {
  const { t } = useTranslation(['common', 'table']);
  const formatDateTime = useFormatDateTime();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const { data: tagCollection } = useFetchTags({
    organizationId: currentEnvironment.organizationId
  });

  const tagList = tagCollection?.tags || [];

  return [
    {
      accessorKey: 'name',
      header: `${t('name')}`,
      size: 350,
      cell: ({ row }) => {
        const account = row.original;
        const isPendingInvite = Number(account.lastSeen) < 1;

        const { avatarImageUrl, firstName, lastName, name, email } =
          account || {};

        const accountName = joinName(firstName, lastName) || name;

        return (
          <div className="flex gap-2">
            <AvatarImage
              image={avatarImageUrl || primaryAvatar}
              alt="member-avatar"
            />
            <div className="flex flex-col gap-0.5">
              {!isPendingInvite && (
                <NameWithTooltip
                  id={`pending_${email}`}
                  content={
                    <NameWithTooltip.Content
                      content={accountName}
                      id={`pending_${email}`}
                    />
                  }
                  trigger={
                    <NameWithTooltip.Trigger
                      id={`pending_${email}`}
                      name={accountName}
                      maxLines={1}
                      className="min-w-[300px]"
                      onClick={() => onActions(account)}
                    />
                  }
                  maxLines={1}
                />
              )}

              <NameWithTooltip
                id={email}
                content={<NameWithTooltip.Content content={name} id={email} />}
                trigger={
                  <NameWithTooltip.Trigger
                    id={email}
                    name={email}
                    maxLines={1}
                    className="text-gray-700 no-underline cursor-default min-w-[300px]"
                  />
                }
                maxLines={1}
              />
              {isPendingInvite && (
                <div className="py-[3px] px-2 w-fit rounded bg-accent-orange-50 typo-para-small text-accent-orange-500">
                  {t('table:pending-invite')}
                </div>
              )}
            </div>
          </div>
        );
      }
    },
    {
      accessorKey: 'organizationRole',
      header: `${t('role')}`,
      size: 300,
      cell: ({ row }) => {
        const account = row.original;
        return (
          <div className="text-gray-700 typo-para-medium">
            {t(String(account.organizationRole).split('_')[1]?.toLowerCase())}
          </div>
        );
      }
    },
    {
      accessorKey: 'teams',
      header: `${t('teams')}`,
      size: 300,
      cell: ({ row }) => {
        const account = row.original;
        const formattedTags = account.teams?.map(
          item => tagList.find(team => team.id === item)?.name || item
        );

        return (
          <ExpandableTag
            tags={formattedTags}
            rowId={account.email}
            className="!max-w-[250px] truncate"
          />
        );
      }
    },
    {
      accessorKey: 'environmentCount',
      header: `${t('environments')}`,
      size: 120,
      cell: ({ row }) => {
        const account = row.original;
        return (
          <div className="text-gray-700 typo-para-medium">
            {account.environmentRoles.length}
          </div>
        );
      }
    },
    {
      accessorKey: 'lastSeen',
      header: `${t('table:last-seen')}`,
      size: 180,
      cell: ({ row }) => {
        const account = row.original;
        return (
          <div className="text-gray-700 typo-para-medium">
            {Number(account.lastSeen) === 0
              ? t('never')
              : formatDateTime(account.lastSeen)}
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
      cell: () => {
        return (
          <button className="flex-center text-gray-600">
            <Icon icon={IconMoreHorizOutlined} size="sm" />
          </button>
        );
      }
    }
  ];
};
