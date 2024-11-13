import {
  IconAddOutlined,
  IconEditOutlined,
  IconMoreHorizOutlined
} from 'react-icons-material-design';
import type { ColumnDef } from '@tanstack/react-table';
import primaryAvatar from 'assets/avatars/primary.svg';
import { useTranslation } from 'i18n';
import { Account } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import { joinName } from 'utils/name';
import { IconTrash } from '@icons';
import { AvatarImage } from 'components/avatar';
import { Popover } from 'components/popover';
import { MemberActionsType } from '../types';

export const useColumns = ({
  onActions
}: {
  onActions: (item: Account, type: MemberActionsType) => void;
}): ColumnDef<Account>[] => {
  const { t } = useTranslation(['common', 'table']);
  const formatDateTime = useFormatDateTime();

  return [
    {
      accessorKey: 'email',
      header: `${t('name')}`,
      size: 350,
      cell: ({ row }) => {
        const account = row.original;
        return (
          <div className="flex gap-2">
            <AvatarImage image={account?.avatarImageUrl || primaryAvatar} />
            <div className="flex flex-col gap-0.5">
              <div className="underline text-primary-500 typo-para-medium">
                {joinName(account.firstName, account.lastName) || account.name}
              </div>
              <div className="typo-para-medium text-gray-700">
                {account.email}
              </div>
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
            {String(account.organizationRole).split('_')[1]}
          </div>
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
      accessorKey: 'createdAt',
      header: `${t('table:created-at')}`,
      size: 180,
      cell: ({ row }) => {
        const account = row.original;
        return (
          <div className="text-gray-700 typo-para-medium">
            {formatDateTime(account.createdAt)}
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
        const organization = row.original;

        return (
          <Popover
            options={[
              {
                label: `${t('table:popover.edit-member')}`,
                icon: IconEditOutlined,
                value: 'EDIT'
              },
              {
                label: `${t('table:popover.add-to-env')}`,
                icon: IconAddOutlined,
                value: 'ADD_ENV'
              },
              {
                label: `${t('table:popover.delete-member')}`,
                icon: IconTrash,
                value: 'DELETE'
              }
            ]}
            icon={IconMoreHorizOutlined}
            onClick={value =>
              onActions(organization, value as MemberActionsType)
            }
            align="end"
          />
        );
      }
    }
  ];
};
