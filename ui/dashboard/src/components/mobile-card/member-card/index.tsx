import { useCallback } from 'react';
import { Trans } from 'react-i18next';
import {
  IconEditOutlined,
  IconMoreVertOutlined
} from 'react-icons-material-design';
import { useAuth, useAuthAccess } from 'auth';
import { useTranslation } from 'i18n';
import { compact } from 'lodash';
import { ClockIcon } from 'lucide-react';
import { Account, Team } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import { IconMember, IconTrash } from '@icons';
import { MemberActionsType, MembersFilters } from 'pages/members/types';
import Divider from 'components/divider';
import Icon from 'components/icon';
import Switch from 'components/switch';
import DateTooltip from 'elements/date-tooltip';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import DisabledPopoverTooltip from 'elements/disabled-popover-tooltip';
import ExpandableTag from 'elements/expandable-tag';
import NameWithTooltip from 'elements/name-with-tooltip';
import { Card } from '../elements';

interface MemberCardProps {
  data: Account;
  filters: MembersFilters;
  teams: Team[];
  onActions: (item: Account, type: MemberActionsType) => void;
  setFilters: (values: Partial<MembersFilters>) => void;
}

export const MemberCard: React.FC<MemberCardProps> = ({
  data,
  filters,
  teams,
  setFilters,
  onActions
}) => {
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
  const hasEnableEdit = onGetActionsType(data) === 'EDIT';
  const isPendingInvite = Number(data.lastSeen) < 1;

  const { email } = data || {};

  // const accountName = joinName(firstName, lastName) || name;
  const formattedTeams = data.teams?.map(
    item => teams.find(team => team.id === item)?.name || item
  );
  return (
    <Card>
      <Card.Header
        icon={<Icon icon={IconMember} />}
        triger={
          <div>
            <NameWithTooltip
              id={email}
              content={<NameWithTooltip.Content content={email} id={email} />}
              trigger={
                <NameWithTooltip.Trigger
                  id={email}
                  name={email}
                  maxLines={1}
                  className="min-w-[100px]"
                  haveAction={false}
                />
              }
              maxLines={1}
            />
            <div className="flex gap-2 items-center">
              {isPendingInvite && (
                <div className="py-[2px] px-1 w-fit rounded bg-accent-orange-50 typo-para-small text-accent-orange-500">
                  {t(`table:pending-invite`)}
                </div>
              )}
              <div className="text-gray-700 typo-para-small capitalize">
                {t(String(data.organizationRole).split('_')[1]?.toLowerCase())}
              </div>
            </div>
          </div>
        }
      >
        <Card.Action>
          <DisabledPopoverTooltip
            icon={IconMoreVertOutlined}
            isNeedAdminAccess
            options={compact([
              Number(data.lastSeen) > 0 &&
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
            onClick={value => onActions(data, value as MemberActionsType)}
          />
        </Card.Action>
      </Card.Header>

      <Card.Meta>
        <div className="flex flex-col gap-y-2 items-center justify-between">
          <ExpandableTag
            tags={formattedTeams}
            rowId={data.email}
            filterTags={filters.teams}
            className="!max-w-[200px] sm:!max-w-[250px] truncate cursor-pointer"
            onTagClick={team => handleFilterTeams(team)}
          />
        </div>
      </Card.Meta>
      <Divider className="pb-3" />
      <Card.Footer
        left={
          <div className="flex flex-col gap-1">
            <p className="typo-para-tiny font-bold uppercase text-gray-500">
              {t('state')}
            </p>
            <DisabledButtonTooltip
              type={!isOrganizationAdmin ? 'admin' : 'editor'}
              hidden={envEditable && isOrganizationAdmin}
              trigger={
                <div className="w-fit">
                  <Switch
                    checked={isPendingInvite ? false : !data.disabled}
                    disabled={
                      isPendingInvite || !envEditable || !isOrganizationAdmin
                    }
                    onCheckedChange={value =>
                      onActions(data, value ? 'ENABLE' : 'DISABLE')
                    }
                  />
                </div>
              }
            />
          </div>
        }
        right={
          <div className="flex flex-col gap-1">
            <p className="typo-para-tiny font-bold uppercase text-gray-500">
              {t('last-seen')}
            </p>
            <div className="text-gray-500 flex-center gap-2">
              <Icon icon={ClockIcon} size={'xxs'} />
              <DateTooltip
                trigger={
                  <div className="typo-para-small whitespace-nowrap">
                    {Number(data.lastSeen) === 0 ? (
                      t('never')
                    ) : (
                      <Trans
                        i18nKey={'common:time-updated'}
                        values={{
                          time: formatDateTime(data.lastSeen)
                        }}
                      />
                    )}
                  </div>
                }
                date={Number(data.lastSeen) === 0 ? null : data.lastSeen}
              />
            </div>
          </div>
        }
      />
    </Card>
  );
};
