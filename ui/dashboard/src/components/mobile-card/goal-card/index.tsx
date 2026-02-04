import {
  IconArchiveOutlined,
  IconMoreVertOutlined
} from 'react-icons-material-design';
import { Link } from 'react-router-dom';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_GOALS } from 'constants/routing';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import { compact } from 'lodash';
import { ConnectionType, Goal } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import { copyToClipBoard } from 'utils/function';
import { useSearchParams } from 'utils/search-params';
import { cn } from 'utils/style';
import {
  IconArrowDown,
  IconCopy,
  IconGoal,
  IconProton,
  IconTrash,
  IconWatch
} from '@icons';
import { GoalActions } from 'pages/goals/types';
import Divider from 'components/divider';
import Icon from 'components/icon';
import DateTooltip from 'elements/date-tooltip';
import DisabledPopoverTooltip from 'elements/disabled-popover-tooltip';
import NameWithTooltip from 'elements/name-with-tooltip';
import { Card } from '../elements';

export const Tag = ({ tag, type }: { tag: string; type: ConnectionType }) => {
  return (
    <div
      className={cn(
        'flex-center w-fit px-2 py-1.5 typo-para-small leading-[14px] text-center rounded-[3px] capitalize whitespace-nowrap',
        {
          'px-4 text-gray-600 bg-gray-100': type === 'UNKNOWN',
          'text-primary-500 bg-primary-50': type === 'EXPERIMENT',
          'text-accent-pink-500 bg-accent-pink-50': type === 'OPERATION'
        }
      )}
    >
      {tag}
    </div>
  );
};

const TagGoal = ({
  goal,
  onActions
}: {
  goal: Goal;
  onActions: (item: Goal, type: GoalActions) => void;
}) => {
  const { t } = useTranslation(['common', 'table']);
  const connectionCount =
    goal.connectionType === 'EXPERIMENT'
      ? goal.experiments?.length
      : goal?.autoOpsRules?.length;

  if (!goal.isInUseStatus && goal.connectionType === 'UNKNOWN')
    return <Tag tag={'not in use'} type="UNKNOWN" />;
  return (
    <button
      disabled={!connectionCount}
      onClick={() => connectionCount && onActions(goal, 'CONNECTION')}
      className="flex items-center gap-1"
    >
      <Tag
        tag={
          goal.connectionType === 'EXPERIMENT'
            ? t('form:experiment', { count: connectionCount })
            : t('form:operation', { count: connectionCount })
        }
        type={goal.connectionType}
      />
      {connectionCount > 0 && <Icon icon={IconArrowDown} />}
    </button>
  );
};

interface GoalCardProps {
  data: Goal;
  onActions: (item: Goal, type: GoalActions) => void;
}

export const GoalCard: React.FC<GoalCardProps> = ({ data, onActions }) => {
  const { t } = useTranslation(['common', 'table']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const { searchOptions } = useSearchParams();
  const { isInUseStatus } = data;
  const formatDateTime = useFormatDateTime();
  const { notify } = useToast();
  const handleCopyId = (id: string) => {
    copyToClipBoard(id);
    notify({
      message: t('message:copied')
    });
  };
  return (
    <Card>
      <Card.Header
        icon={<Icon icon={IconGoal} />}
        triger={
          <div className="flex flex-col gap-0.5 max-w-fit">
            <NameWithTooltip
              id={data.id}
              content={
                <NameWithTooltip.Content content={data.name} id={data.id} />
              }
              trigger={
                <Link
                  to={`/${currentEnvironment.urlCode}${PAGE_PATH_GOALS}/${data.id}`}
                >
                  <NameWithTooltip.Trigger id={data.id} name={data.name} />
                </Link>
              }
              maxLines={1}
            />
            <div className="flex items-center h-5 gap-x-2 typo-para-tiny text-gray-500 group select-none">
              {data.id}
              <div onClick={() => handleCopyId(data.id)}>
                <Icon
                  icon={IconCopy}
                  size={'sm'}
                  className="opacity-0 group-hover:opacity-100 cursor-pointer"
                />
              </div>
            </div>
          </div>
        }
      >
        <Card.Action>
          <DisabledPopoverTooltip
            icon={IconMoreVertOutlined}
            options={compact([
              searchOptions.status === 'ARCHIVED'
                ? {
                    label: `${t('table:popover.unarchive-goal')}`,
                    icon: IconArchiveOutlined,
                    value: 'UNARCHIVE'
                  }
                : {
                    label: `${t('table:popover.archive-goal')}`,
                    icon: IconArchiveOutlined,
                    value: 'ARCHIVE',
                    disabled: isInUseStatus,
                    tooltip: isInUseStatus
                      ? t('form:goal-details.archive-warning-desc')
                      : ''
                  },
              {
                label: `${t('table:popover.delete-goal')}`,
                icon: IconTrash,
                value: 'DELETE',
                disabled: isInUseStatus,
                tooltip: isInUseStatus
                  ? t('form:goal-details.delete-warning-desc')
                  : ''
              }
            ])}
            onClick={value => onActions(data, value as GoalActions)}
          />
        </Card.Action>
      </Card.Header>

      <Card.Meta>
        <div className="flex h-full w-full items-stretch justify-between gap-3 pb-3">
          <div className="flex-1 typo-para-tiny font-bold bg-gray-100 p-3 rounded-xl">
            <p className="flex items-center gap-1 uppercase text-gray-500">
              <span>{t('table:goals.connections')}</span>
            </p>
            <div className="mt-3 flex items-center gap-1">
              <Icon icon={IconProton} size="sm" className="text-primary-500" />
              <TagGoal goal={data} onActions={onActions} />
            </div>
          </div>
        </div>
        <Divider />
      </Card.Meta>
      <Card.Footer
        left={
          <DateTooltip
            trigger={
              <div className="flex items-center gap-1 text-gray-500 typo-para-small whitespace-nowrap">
                <Icon icon={IconWatch} size={'xxs'} />
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
