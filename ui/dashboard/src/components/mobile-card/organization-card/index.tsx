import {
  IconArchiveOutlined,
  IconEditOutlined,
  IconMoreVertOutlined
} from 'react-icons-material-design';
import { Link } from 'react-router-dom';
import { PAGE_PATH_ORGANIZATIONS, PAGE_PATH_PROJECTS } from 'constants/routing';
import { useTranslation } from 'i18n';
import { ClockIcon } from 'lucide-react';
import { Organization } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import { useSearchParams } from 'utils/search-params';
import { IconFolder, IconMember, IconUser } from '@icons';
import { OrganizationActionsType } from 'pages/organizations/types';
import Icon from 'components/icon';
import { Popover } from 'components/popover';
import DateTooltip from 'elements/date-tooltip';
import NameWithTooltip from 'elements/name-with-tooltip';
import { Card } from '../elements';

interface UserSegmentCardProps {
  data: Organization;
  onActions: (item: Organization, type: OrganizationActionsType) => void;
}

export const OrganizationCard: React.FC<UserSegmentCardProps> = ({
  data,
  onActions
}) => {
  const { t } = useTranslation(['common', 'table']);
  const { searchOptions } = useSearchParams();
  const formatDateTime = useFormatDateTime();
  return (
    <Card>
      <Card.Header
        icon={<Icon icon={IconUser} />}
        triger={
          <div className="flex flex-col gap-0.5 max-w-fit">
            <NameWithTooltip
              id={data.id}
              content={
                <NameWithTooltip.Content content={data.name} id={data.id} />
              }
              trigger={
                <Link
                  to={`${PAGE_PATH_ORGANIZATIONS}/${data.id}${PAGE_PATH_PROJECTS}`}
                >
                  <NameWithTooltip.Trigger id={data.id} name={data.name} />
                </Link>
              }
            />
          </div>
        }
      >
        <Card.Action>
          <Popover
            icon={IconMoreVertOutlined}
            options={[
              {
                label: `${t('table:popover.edit-org')}`,
                icon: IconEditOutlined,
                value: 'EDIT'
              },
              searchOptions.status === 'ARCHIVED'
                ? {
                    label: `${t('table:popover.unarchive-org')}`,
                    icon: IconArchiveOutlined,
                    value: 'UNARCHIVE'
                  }
                : {
                    label: `${t('table:popover.archive-org')}`,
                    icon: IconArchiveOutlined,
                    value: 'ARCHIVE'
                  }
            ]}
            onClick={value => onActions(data, value as OrganizationActionsType)}
            align="end"
          />
        </Card.Action>
      </Card.Header>

      <Card.Meta>
        <div className="grid grid-cols-2 h-full w-full flex-wrap items-stretch justify-between gap-3">
          <div className="flex-1 font-bold typo-para-tiny bg-gray-100 p-3 rounded-xl">
            <p className="flex items-center gap-1 uppercase text-gray-500">
              <span>{t('projects')}</span>
            </p>
            <div className="mt-3 flex items-center gap-1 text-primary-500">
              <Icon icon={IconFolder} size="sm" />
              <p className="text-nowrap typo-para-small">{data.projectCount}</p>
            </div>
          </div>
          <div className="flex-1 font-bold typo-para-tiny bg-gray-100 p-3 rounded-xl">
            <p className="flex items-center gap-1 uppercase text-gray-500">
              <span>{t('environments')}</span>
            </p>
            <div className="mt-3 flex items-center gap-1 text-primary-500">
              <Icon icon={IconFolder} size="sm" />
              <p className="text-nowrap typo-para-small">
                {data.environmentCount}
              </p>
            </div>
          </div>
          <div className="flex-1 font-bold typo-para-tiny bg-gray-100 p-3 rounded-xl">
            <p className="flex items-center gap-1 uppercase text-gray-500">
              <span>{t('users')}</span>
            </p>
            <div className="mt-3 flex items-center gap-1 text-primary-500 font-bold">
              <Icon icon={IconMember} size="sm" />
              <p className="text-nowrap font-bold typo-para-small">
                {data.userCount}
              </p>
            </div>
          </div>
          <div className="flex-1 font-bold typo-para-tiny bg-gray-100 p-3 rounded-xl">
            <p className="flex items-center gap-1 uppercase text-gray-500">
              <span>{t('table:created-at')}</span>
            </p>
            <div className="mt-3 flex items-center gap-1 text-gray-500">
              <Icon icon={ClockIcon} size="xs" />
              <DateTooltip
                trigger={
                  <div className="typo-para-small">
                    {formatDateTime(data.createdAt)}
                  </div>
                }
                date={data.createdAt}
              />
            </div>
          </div>
        </div>
      </Card.Meta>
    </Card>
  );
};
