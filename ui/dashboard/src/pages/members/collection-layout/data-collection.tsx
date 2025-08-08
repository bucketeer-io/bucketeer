import { useCallback } from 'react';
import { IconEditOutlined } from 'react-icons-material-design';
import type { ColumnDef } from '@tanstack/react-table';
import primaryAvatar from 'assets/avatars/primary.svg';
import { useAuth, useAuthAccess } from 'auth';
import { useTranslation } from 'i18n';
import compact from 'lodash/compact';
import { Account, Team } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import { joinName } from 'utils/name';
import { IconTrash } from '@icons';
import { AvatarImage } from 'components/avatar';
import Switch from 'components/switch';
import DateTooltip from 'elements/date-tooltip';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import DisabledPopoverTooltip from 'elements/disabled-popover-tooltip';
import ExpandableTag from 'elements/expandable-tag';
import NameWithTooltip from 'elements/name-with-tooltip';
import { MemberActionsType, MembersFilters } from '../types';

export const useColumns = ({
  filters,
  teams,
  onActions,
  setFilters
}: {
  filters: MembersFilters;
  teams: Team[];
  onActions: (item: Account, type: MemberActionsType) => void;
  setFilters: (values: Partial<MembersFilters>) => void;
}): ColumnDef<Account>[] => {
  const { t } = useTranslation(['common', 'table']);
  const formatDateTime = useFormatDateTime();
  const { consoleAccount } = useAuth();
  const { envEditable, isOrganizationAdmin } = useAuthAccess();
  const isAccountOwnerRole =
    consoleAccount?.organizationRole === 'Organization_OWNER';
  const isSystemAdmin = consoleAccount?.isSystemAdmin;

  const handleFilterTeams = useCallback(
    (team: string) => {
      const isExisted = filters.teams?.includes(team);
      const newTeams = isExisted
        ? filters.teams?.filter(item => item !== team)
        : [...(filters.teams || []), team];
      setFilters({
        teams: newTeams?.length ? newTeams : undefined
      });
    },
    [filters]
  );

  const onGetActionsType = useCallback(
    (account: Account): MemberActionsType => {
      const isUserOwner = account.organizationRole === 'Organization_OWNER';
      const canEditMember =
        isSystemAdmin ||
        (isUserOwner && isAccountOwnerRole) ||
        (!isUserOwner && isOrganizationAdmin);

      return canEditMember ? 'EDIT' : 'DETAILS';
    },
    [isAccountOwnerRole, isOrganizationAdmin, isSystemAdmin]
  );

  return compact([
    {
      accessorKey: 'email',
      header: `${t('name')}`,
      size: 350,
      cell: ({ row }) => {
        const account = row.original;
        const isPendingInvite = Number(account.lastSeen) < 1;

        const { name, email, firstName, lastName, avatarImageUrl } =
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
                  id={email}
                  content={
                    <NameWithTooltip.Content content={accountName} id={email} />
                  }
                  trigger={
                    <NameWithTooltip.Trigger
                      id={email}
                      name={accountName}
                      maxLines={1}
                      className="min-w-[300px]"
                      onClick={() =>
                        onActions(account, onGetActionsType(account))
                      }
                    />
                  }
                  maxLines={1}
                />
              )}
              <NameWithTooltip
                id={email}
                content={<NameWithTooltip.Content content={email} id={email} />}
                trigger={
                  <NameWithTooltip.Trigger
                    id={email}
                    name={email}
                    maxLines={1}
                    className="min-w-[300px]"
                    haveAction={false}
                  />
                }
                maxLines={1}
              />

              {isPendingInvite && (
                <div className="py-[3px] px-2 w-fit rounded bg-accent-orange-50 typo-para-small text-accent-orange-500">
                  {t(`table:pending-invite`)}
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
          <div className="text-gray-700 typo-para-medium capitalize">
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
        const formattedTeams = account.teams?.map(
          item => teams.find(team => team.id === item)?.name || item
        );
        return (
          <ExpandableTag
            tags={formattedTeams}
            rowId={account.email}
            filterTags={filters.teams}
            className="!max-w-[250px] truncate cursor-pointer"
            onTagClick={team => handleFilterTeams(team)}
          />
        );
      }
    },
    {
      accessorKey: 'lastSeen',
      header: `${t('table:last-seen')}`,
      size: 180,
      cell: ({ row }) => {
        const account = row.original;
        const isNever = Number(account.lastSeen) === 0;
        return (
          <DateTooltip
            trigger={
              <div className="text-gray-700 typo-para-medium">
                {isNever ? t('never') : formatDateTime(account.lastSeen)}
              </div>
            }
            date={isNever ? null : account.lastSeen}
          />
        );
      }
    },
    {
      accessorKey: 'state',
      header: `${t('state')}`,
      size: 120,
      cell: ({ row }) => {
        const account = row.original;
        const isPendingInvite = Number(account.lastSeen) < 1;

        return (
          <DisabledButtonTooltip
            type={!isOrganizationAdmin ? 'admin' : 'editor'}
            hidden={envEditable && isOrganizationAdmin}
            trigger={
              <div className="w-fit">
                <Switch
                  checked={isPendingInvite ? false : !account.disabled}
                  disabled={
                    isPendingInvite || !envEditable || !isOrganizationAdmin
                  }
                  onCheckedChange={value =>
                    onActions(account, value ? 'ENABLE' : 'DISABLE')
                  }
                />
              </div>
            }
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
        const account = row.original;
        const hasEnableEdit = onGetActionsType(account) === 'EDIT';

        return (
          <DisabledPopoverTooltip
            isNeedAdminAccess
            options={compact([
              Number(account.lastSeen) > 0 &&
                hasEnableEdit && {
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
            onClick={value => onActions(account, value as MemberActionsType)}
          />
        );
      }
    }
  ]);
};
