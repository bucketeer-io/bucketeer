import {
  IconEditOutlined,
  IconMoreHorizOutlined
} from 'react-icons-material-design';
import type { ColumnDef } from '@tanstack/react-table';
import primaryAvatar from 'assets/avatars/primary.svg';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useTranslation } from 'i18n';
import compact from 'lodash/compact';
import { Account } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import { joinName } from 'utils/name';
import { IconTrash } from '@icons';
import { AvatarImage } from 'components/avatar';
import { Popover } from 'components/popover';
import Switch from 'components/switch';
import ExpandableTag from 'elements/expandable-tag';
import TruncationWithTooltip from 'elements/truncation-with-tooltip';
import { useFetchTags } from '../collection-loader';
import { MemberActionsType } from '../types';

export const useColumns = ({
  onActions
}: {
  onActions: (item: Account, type: MemberActionsType) => void;
}): ColumnDef<Account>[] => {
  const { t } = useTranslation(['common', 'table']);
  const formatDateTime = useFormatDateTime();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const isOrganizationAdmin =
    consoleAccount?.organizationRole === 'Organization_ADMIN';

  const { data: tagCollection } = useFetchTags({
    organizationId: currentEnvironment.organizationId
  });

  const tagList = tagCollection?.tags || [];

  return compact([
    {
      accessorKey: 'email',
      header: `${t('name')}`,
      size: 350,
      cell: ({ row }) => {
        const account = row.original;
        const isPendingInvite = Number(account.lastSeen) < 1;

        return (
          <div className="flex gap-2">
            <AvatarImage
              image={account?.avatarImageUrl || primaryAvatar}
              alt="member-avatar"
            />
            <div className="flex flex-col gap-0.5">
              {!isPendingInvite && (
                <button
                  onClick={() =>
                    onActions(account, isOrganizationAdmin ? 'EDIT' : 'DETAILS')
                  }
                  className="underline text-primary-500 typo-para-medium text-left"
                >
                  {joinName(account.firstName, account.lastName) ||
                    account.name}
                </button>
              )}
              <TruncationWithTooltip
                content={account.email}
                maxSize={320}
                elementId={`email-${account.email}`}
                className="w-fit"
                additionalClassName={['line-clamp-1', 'break-all']}
              >
                <div
                  id={`email-${account.email}`}
                  className="typo-para-medium text-gray-700"
                >
                  {account.email}
                </div>
              </TruncationWithTooltip>
              {isPendingInvite && (
                <div className="py-[3px] px-2 w-fit rounded bg-accent-orange-50 typo-para-small text-accent-orange-500">
                  {`Pending invite`}
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
      size: 180,
      cell: ({ row }) => {
        const account = row.original;
        return (
          <div className="text-gray-700 typo-para-medium">
            {String(account.organizationRole).split('_')[1]}
          </div>
        );
      }
    },
    {
      accessorKey: 'tags',
      header: `${t('tags')}`,
      size: 300,
      cell: ({ row }) => {
        const account = row.original;
        const formattedTags = account.tags?.map(
          item => tagList.find(tag => tag.id === item)?.name || item
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
    isOrganizationAdmin && {
      accessorKey: 'state',
      header: `${t('state')}`,
      size: 120,
      cell: ({ row }) => {
        const account = row.original;
        const isPendingInvite = Number(account.lastSeen) < 1;

        return (
          <Switch
            checked={isPendingInvite ? false : !account.disabled}
            disabled={isPendingInvite}
            onCheckedChange={value =>
              onActions(account, value ? 'ENABLE' : 'DISABLE')
            }
          />
        );
      }
    },
    isOrganizationAdmin && {
      accessorKey: 'action',
      size: 60,
      header: '',
      meta: {
        align: 'center',
        style: { textAlign: 'center', fitContent: true }
      },
      enableSorting: false,
      cell: ({ row }) => {
        const account = row.original;

        return (
          <Popover
            options={compact([
              Number(account.lastSeen) > 0 && {
                label: `${t('table:popover.edit-member')}`,
                icon: IconEditOutlined,
                value: 'EDIT'
              },
              {
                label: `${t('table:popover.delete-member')}`,
                icon: IconTrash,
                value: 'DELETE'
              }
            ])}
            icon={IconMoreHorizOutlined}
            onClick={value => onActions(account, value as MemberActionsType)}
            align="end"
          />
        );
      }
    }
  ]);
};
