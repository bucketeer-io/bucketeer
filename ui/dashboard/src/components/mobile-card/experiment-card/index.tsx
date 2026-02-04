import {
  IconArchiveOutlined,
  IconEditOutlined,
  IconMoreVertOutlined
} from 'react-icons-material-design';
import { Link } from 'react-router-dom';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_EXPERIMENTS } from 'constants/routing';
import { useTranslation } from 'i18n';
import { Experiment } from '@types';
import { formatLongDateTime } from 'utils/date-time';
import { useSearchParams } from 'utils/search-params';
import {
  IconGoal,
  IconProton,
  IconStartExperiment,
  IconStopExperiment
} from '@icons';
import { ExperimentStatuses } from 'pages/experiments/collection-layout/data-collection';
import { ExperimentActionsType } from 'pages/experiments/types';
import Divider from 'components/divider';
import Icon from 'components/icon';
import DateTooltip from 'elements/date-tooltip';
import DisabledPopoverTooltip from 'elements/disabled-popover-tooltip';
import NameWithTooltip from 'elements/name-with-tooltip';
import { Card } from '../elements';

export interface ExperimentCardViewModel {
  id: string;
  title: string;
  code: string;
  goalsCount: number;
  startDate: string;
  endDate: string;
  status: Experiment['status'];
}

interface ExperimentCardProps {
  data: Experiment;
  onActions: (item: Experiment, type: ExperimentActionsType) => void;
}

export const ExperimentCard: React.FC<ExperimentCardProps> = ({
  data,
  onActions
}) => {
  const { t } = useTranslation(['common', 'table']);

  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const { searchOptions } = useSearchParams();
  const { status } = data;
  const formatDate = (value: string) => {
    return formatLongDateTime({
      value: value,
      overrideOptions: {
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        hour12: false
      },
      locale: 'ja-JP'
    });
  };
  return (
    <Card>
      <Card.Header
        icon={<Icon icon={IconProton} />}
        triger={
          <div className="flex flex-col gap-0.5 max-w-fit">
            <NameWithTooltip
              id={data.id}
              content={
                <NameWithTooltip.Content content={data.name} id={data.id} />
              }
              trigger={
                <Link
                  to={`/${currentEnvironment.urlCode}${PAGE_PATH_EXPERIMENTS}/${data.id}/results`}
                >
                  <NameWithTooltip.Trigger id={data.id} name={data.name} />
                </Link>
              }
              maxLines={1}
            />
          </div>
        }
      >
        <Card.Action>
          <DisabledPopoverTooltip
            onClick={value => onActions(data, value as ExperimentActionsType)}
            icon={IconMoreVertOutlined}
            options={[
              {
                label: `${t('table:popover.edit-experiment')}`,
                icon: IconEditOutlined,
                value: 'EDIT'
              },
              ...(['WAITING', 'RUNNING'].includes(status)
                ? [
                    status === 'WAITING'
                      ? {
                          label: `${t('table:popover.start-experiment')}`,
                          icon: IconStartExperiment,
                          value: 'START'
                        }
                      : {
                          label: `${t('table:popover.stop-experiment')}`,
                          icon: IconStopExperiment,
                          value: 'STOP'
                        }
                  ]
                : []),
              searchOptions.status === 'ARCHIVED'
                ? {
                    label: `${t('table:popover.unarchive-experiment')}`,
                    icon: IconArchiveOutlined,
                    value: 'UNARCHIVE'
                  }
                : {
                    label: `${t('table:popover.archive-experiment')}`,
                    icon: IconArchiveOutlined,
                    value: 'ARCHIVE',
                    disabled: ['RUNNING', 'WAITING'].includes(data.status)
                  }
            ]}
          />
        </Card.Action>
      </Card.Header>

      <Card.Meta>
        <div className="flex flex-wrap h-full w-full items-stretch justify-between gap-3 typo-para-medium pb-5">
          <div className="flex-1 p-3 rounded-xl bg-gray-100">
            <p className="flex items-center gap-2 uppercase typo-para-tiny text-gray-500">
              <span>{t('common:goal')}</span>
            </p>
            <div
              className="mt-2 flex items-center gap-2"
              onClick={() => onActions(data, 'GOALS-CONNECTION')}
            >
              <Icon icon={IconGoal} size="sm" /> {data.goals.length}
              {t('common:goal')}
            </div>
          </div>

          <div className="flex-1 p-3 rounded-xl bg-gray-100 text-nowrap">
            <div className="flex-1">
              <p className="flex items-center gap-2 uppercase typo-para-tiny text-gray-500">
                <span>{t('common:status')}</span>
              </p>
              <div className="mt-2">
                <ExperimentStatuses status={data.status} />
              </div>
            </div>
          </div>
        </div>
        <Divider />
        <div className="flex flex-wrap h-full w-full pt-3 items-stretch justify-between gap-3 typo-para-medium">
          <div className="flex-1 p-3 rounded-xl text-gray-500">
            <p className="flex items-center gap-2 uppercase typo-para-tiny">
              <span>{t('time-start')}</span>
            </p>
            <div className="mt-2 typo-para-small text-nowrap">
              <DateTooltip
                trigger={
                  <div className="text-gray-700 typo-para-medium min-w-[150px]">
                    {formatDate(data.startAt)}
                  </div>
                }
                date={data.startAt}
              />
            </div>
          </div>

          <div className="flex-1 p-3 rounded-xl text-gray-500">
            <div className="flex-1">
              <p className="flex items-center gap-2 uppercase typo-para-tiny">
                <span>{t('common:time-stop')}</span>
              </p>
              <div className="mt-2 typo-para-small text-nowrap">
                <DateTooltip
                  trigger={
                    <div className="text-gray-700 typo-para-medium min-w-[150px]">
                      {formatDate(data.stopAt)}
                    </div>
                  }
                  date={data.stopAt}
                />
              </div>
            </div>
          </div>
        </div>
      </Card.Meta>
    </Card>
  );
};
