import { useCallback } from 'react';
import {
  IconEditOutlined,
  IconMoreVertOutlined
} from 'react-icons-material-design';
import { Link } from 'react-router-dom';
import primaryAvatar from 'assets/avatars/primary.svg';
import { PAGE_PATH_ENVIRONMENTS, PAGE_PATH_PROJECTS } from 'constants/routing';
import { useTranslation } from 'i18n';
import { ClockIcon } from 'lucide-react';
import { Account, Environment, Project } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import { joinName } from 'utils/name';
import { IconGoal } from '@icons';
import { AvatarImage } from 'components/avatar';
import Divider from 'components/divider';
import Icon from 'components/icon';
import DateTooltip from 'elements/date-tooltip';
import DisabledPopoverTooltip from 'elements/disabled-popover-tooltip';
import NameWithTooltip from 'elements/name-with-tooltip';
import { Card } from '../elements';

const ProjectAvatar = ({
  project,
  handleGetCurrentAccount
}: {
  project: Project;
  handleGetCurrentAccount: (value: string) => Account | undefined;
}) => {
  const { id, creatorEmail } = project;
  const { firstName, lastName, name, avatarImageUrl } =
    handleGetCurrentAccount(creatorEmail) || {};
  const accountName = joinName(firstName, lastName) || name;

  return (
    <div className="flex gap-2 p-3 bg-gray-100 rounded-xl mb-3">
      <AvatarImage
        image={avatarImageUrl || primaryAvatar}
        alt="member-avatar"
        className="min-w-[50px] h-[50px] rounded-full"
      />
      <div className="flex flex-col gap-0.5">
        <NameWithTooltip
          id={creatorEmail}
          content={
            <NameWithTooltip.Content
              className="text-[16px]"
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
};

interface ProjectCardProps {
  data: Project;
  accounts: Account[];
  currentEnvironment: Environment;
  organizationId?: string;
  onActions: (item: Project) => void;
}

export const ProjectCard: React.FC<ProjectCardProps> = ({
  data,
  accounts,
  organizationId,
  currentEnvironment,
  onActions
}) => {
  const { t } = useTranslation(['common', 'table']);
  const formatDateTime = useFormatDateTime();

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

  return (
    <Card>
      <Card.Header
        triger={
          <NameWithTooltip
            id={data.id}
            content={
              <NameWithTooltip.Content content={data.name} id={data.id} />
            }
            trigger={
              <Link
                to={`/${currentEnvironment.urlCode}${PAGE_PATH_PROJECTS}/${data.id}${PAGE_PATH_ENVIRONMENTS}?organizationId=${organizationId}`}
              >
                <NameWithTooltip.Trigger
                  id={data.id}
                  name={data.name}
                  maxLines={1}
                  className="min-w-[230px]"
                />
              </Link>
            }
            maxLines={1}
          />
        }
      >
        <Card.Action>
          <DisabledPopoverTooltip
            icon={IconMoreVertOutlined}
            isNeedAdminAccess
            options={[
              {
                label: `${t('table:popover.edit-project')}`,
                icon: IconEditOutlined,
                value: 'EDIT_PROJECT'
              }
            ]}
            onClick={() => onActions(data)}
          />
        </Card.Action>
      </Card.Header>

      <Card.Meta>
        <ProjectAvatar
          project={data}
          handleGetCurrentAccount={handleGetCurrentAccount}
        />
        <div className="flex h-full w-full items-stretch justify-between gap-3 pb-3">
          <div className="flex-1 typo-para-tiny font-bold bg-gray-100 p-3 rounded-xl">
            <p className="flex items-center gap-1 uppercase text-gray-500">
              <span>{t('environment')}</span>
            </p>
            <div className="mt-3 flex items-center gap-2">
              <Icon icon={IconGoal} size="sm" /> {data.environmentCount}
            </div>
          </div>
          <div className="flex-1 typo-para-tiny font-bold bg-gray-100 p-3 rounded-xl">
            <p className="flex items-center gap-1 uppercase text-gray-500">
              <span>{t('common:flag')}</span>
            </p>
            <div className="mt-3 flex items-center gap-2">
              <Icon icon={IconGoal} size="sm" /> {data.featureFlagCount}
            </div>
          </div>
        </div>
        <Divider />
      </Card.Meta>
      <Card.Footer
        left={
          <DateTooltip
            trigger={
              <div className="text-gray-500 typo-para-medium flex items-center gap-1">
                <Icon icon={ClockIcon} size="xs" />
                {Number(data.updatedAt) === 0
                  ? t('never')
                  : formatDateTime(data.updatedAt)}
              </div>
            }
            date={data.updatedAt}
          />
        }
      />
    </Card>
  );
};
