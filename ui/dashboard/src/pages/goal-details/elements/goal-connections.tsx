import { Link } from 'react-router-dom';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_EXPERIMENTS } from 'constants/routing';
import { useTranslation } from 'i18n';
import { Goal } from '@types';
import { cn } from 'utils/style';
import Divider from 'components/divider';
import Spinner from 'components/spinner';
import Status from 'elements/status';

interface Props {
  goal: Goal;
}

const GoalConnections = ({ goal }: Props) => {
  const { t } = useTranslation(['common', 'table']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const gridRowCls = 'grid grid-cols-12 w-full gap-x-4',
    headerCls =
      'typo-para-small leading-[14px] font-medium text-gray-500 uppercase whitespace-nowrap',
    experimentCls = 'col-span-4 min-w-[200px] truncate',
    experimentStatusCls = 'col-span-3 min-w-[150px]',
    goalStatusCls = 'col-span-5 min-w-[300px]';

  return (
    <div className="flex flex-col w-full min-w-[780px] p-5 gap-y-5 shadow-card rounded-lg bg-white">
      <p className="text-gray-800 typo-head-bold-small">
        {t('goal-connections')}
      </p>
      <Divider />
      <div className="flex flex-col w-full gap-y-3">
        <div className={gridRowCls}>
          <div className={cn(headerCls, experimentCls)}>
            {t('table:goals.experiment')}
          </div>
          <div className={cn(headerCls, experimentStatusCls)}>
            {t('table:goals.experiment-status')}
          </div>
          <div className={cn(headerCls, goalStatusCls)}>
            {t('table:goals.goal-status')}
          </div>
        </div>
        <div className="flex flex-col w-full gap-y-2">
          {goal?.experiments?.map((item, index) => (
            <div key={index} className={gridRowCls}>
              <div className={experimentCls}>
                <Link
                  to={`/${currentEnvironment.urlCode}${PAGE_PATH_EXPERIMENTS}/${item.id}`}
                  className={cn(
                    'w-fit typo-para-medium text-primary-500 underline'
                  )}
                >
                  {item?.name}
                </Link>
              </div>
              <div className={experimentStatusCls}>
                <Status
                  status={item.status}
                  text={item.status?.toLowerCase()?.replace('_', ' ')}
                />
              </div>
              <div className={cn(goalStatusCls, 'flex items-center gap-x-2')}>
                <Spinner className="size-5 min-w-5 min-h-5" />
                <p className="typo-para-small text-gray-800 truncate">
                  This goal has never received an event for this iteration
                </p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default GoalConnections;
