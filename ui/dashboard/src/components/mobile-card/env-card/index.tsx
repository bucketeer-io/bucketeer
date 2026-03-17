import {
  IconArchiveOutlined,
  IconEditOutlined
} from 'react-icons-material-design';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useTranslation } from 'i18n';
import { Environment } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import { useSearchParams } from 'utils/search-params';
import { IconMember } from '@icons';
import { EnvironmentActionsType } from 'pages/project-details/environments/types';
import Icon from 'components/icon';
import DateTooltip from 'elements/date-tooltip';
import DisabledPopoverTooltip from 'elements/disabled-popover-tooltip';
import NameWithTooltip from 'elements/name-with-tooltip';
import { Card } from '../elements';

interface EnvCardProps {
  data: Environment;
  onActions: (item: Environment, type: EnvironmentActionsType) => void;
}

export const EnvironmentCard: React.FC<EnvCardProps> = ({
  data,
  onActions
}) => {
  const { t } = useTranslation(['common', 'table']);
  const formatDateTime = useFormatDateTime();
  const { searchOptions } = useSearchParams();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const isDisabled = currentEnvironment.id === data.id;
  const { id, name, featureFlagCount, createdAt } = data;

  return (
    <Card>
      <Card.Header
        icon={<Icon icon={IconMember} />}
        triger={
          <NameWithTooltip
            id={id}
            content={<NameWithTooltip.Content content={name} id={id} />}
            trigger={
              <NameWithTooltip.Trigger
                id={id}
                name={name}
                maxLines={1}
                onClick={() => onActions(data, 'EDIT')}
              />
            }
            maxLines={1}
          />
        }
      >
        <Card.Action>
          <DisabledPopoverTooltip
            isNeedAdminAccess
            content={
              isDisabled ? t('table:disabled-archive-current-env') : undefined
            }
            options={[
              {
                label: `${t('table:popover.edit-env')}`,
                icon: IconEditOutlined,
                value: 'EDIT'
              },
              searchOptions.status === 'ARCHIVED'
                ? {
                    label: `${t('table:popover.unarchive-env')}`,
                    icon: IconArchiveOutlined,
                    value: 'UNARCHIVE'
                  }
                : {
                    label: `${t('table:popover.archive-env')}`,
                    icon: IconArchiveOutlined,
                    value: 'ARCHIVE',
                    disabled: isDisabled
                  }
            ]}
            onClick={value => onActions(data, value as EnvironmentActionsType)}
          />
        </Card.Action>
      </Card.Header>

      <Card.Meta>
        <div className="flex h-full w-full items-stretch justify-between gap-3 pb-3">
          <div className="flex-1 typo-para-tiny font-bold bg-gray-100 p-3 rounded-xl">
            <p className="flex items-center gap-1 uppercase text-gray-500">
              <span>{t('table:flags')}</span>
            </p>
            <div className="mt-3 flex items-center gap-1">
              <div className="text-gray-700 typo-para-medium">
                {featureFlagCount}
              </div>
            </div>
          </div>
          <div className="flex-1 typo-para-tiny font-bold bg-gray-100 p-3 rounded-xl">
            <p className="flex items-center gap-1 uppercase text-gray-500">
              <span>{t('table:created-at')}</span>
            </p>
            <div className="mt-3 flex items-center gap-1">
              <DateTooltip
                trigger={
                  <div className="text-gray-700 typo-para-medium">
                    {formatDateTime(createdAt)}
                  </div>
                }
                date={createdAt}
              />
            </div>
          </div>
        </div>
      </Card.Meta>
    </Card>
  );
};
