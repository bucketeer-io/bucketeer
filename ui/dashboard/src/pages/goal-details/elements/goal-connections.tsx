import { Link } from 'react-router-dom';
import { getCurrentEnvironment, useAuth } from 'auth';
import {
  PAGE_PATH_EXPERIMENTS,
  PAGE_PATH_FEATURE_AUTOOPS,
  PAGE_PATH_FEATURES
} from 'constants/routing';
import { useTranslation } from 'i18n';
import { ExperimentStatus, Goal, OperationStatus } from '@types';
import { cn } from 'utils/style';
import Divider from 'components/divider';
import Status from 'elements/status';

const headerCls =
    'typo-para-small font-medium text-gray-500 uppercase whitespace-nowrap',
  experimentCls = 'col-span-4 min-w-[200px] truncate',
  experimentStatusCls = 'col-span-3 min-w-[150px]';

const GoalConnections = ({ goal }: { goal: Goal }) => {
  const { t } = useTranslation(['common', 'table']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const isExperimentType = goal.connectionType === 'EXPERIMENT';

  return (
    <div className="flex flex-col w-full p-5 gap-y-5 shadow-card rounded-lg bg-white">
      <p className="text-gray-800 typo-head-bold-small">
        {t('goal-connections')}
      </p>
      <Divider />
      <div className="flex flex-col w-full gap-y-3">
        <div className="grid grid-cols-12 w-full gap-x-4">
          <div className={cn(headerCls, experimentCls)}>
            {isExperimentType
              ? t('table:goals.experiment')
              : t('table:goals.auto-operations')}
          </div>
          <div className={cn(headerCls, experimentStatusCls)}>
            {t('status')}
          </div>
        </div>
        <div className="flex flex-col w-full gap-y-2">
          {isExperimentType
            ? goal.experiments.map((item, index) => (
                <ConnectionItem
                  key={index}
                  url={`/${currentEnvironment.urlCode}${PAGE_PATH_EXPERIMENTS}/${item.id}`}
                  name={item.featureName}
                  status={item.status}
                />
              ))
            : goal.autoOpsRules.map((item, index) => (
                <ConnectionItem
                  key={index}
                  url={`/${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${item.featureId}${PAGE_PATH_FEATURE_AUTOOPS}`}
                  name={item.featureName}
                  status={item.autoOpsStatus}
                />
              ))}
        </div>
      </div>
    </div>
  );
};

const ConnectionItem = ({
  url,
  name,
  status
}: {
  url: string;
  name: string;
  status: ExperimentStatus | OperationStatus;
}) => {
  const { t } = useTranslation(['table']);
  return (
    <div className="grid grid-cols-12 w-full gap-x-4">
      <div className={experimentCls}>
        <Link
          to={url}
          className={cn('w-fit typo-para-medium text-primary-500 underline')}
        >
          {name}
        </Link>
      </div>
      <div className={experimentStatusCls}>
        <Status
          status={status}
          text={t(`experiment.${status?.toLowerCase()?.replace('_', '-')}`)}
        />
      </div>
    </div>
  );
};

export default GoalConnections;
