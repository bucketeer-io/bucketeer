import { Trans } from 'react-i18next';
import { Link } from 'react-router-dom';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_FEATURES } from 'constants/routing';
import { Feature } from '@types';
import { cn } from 'utils/style';
import { IconToastWarning } from '@icons';
import Icon from 'components/icon';

const SegmentWarning = ({
  features,
  className
}: {
  features: Feature[];
  className?: string;
}) => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  return (
    <div
      className={cn(
        'flex flex-col w-full px-4 py-3 bg-accent-yellow-50 border-l-4 border-accent-yellow-500 rounded mt-5',
        className
      )}
    >
      <div className="flex gap-x-2 w-full pr-3">
        <Icon
          icon={IconToastWarning}
          size={'xxs'}
          color="accent-yellow-500"
          className="mt-1"
        />
        <Trans
          i18nKey="form:update-user-segment-warning"
          values={{ count: 1 }}
          components={{
            p: <p className="typo-para-medium text-accent-yellow-500" />
          }}
        />
      </div>
      <div className="flex flex-col w-full gap-y-1">
        {features?.map((item, index) => (
          <div
            key={item.id}
            className="flex gap-x-2 w-full pl-6 typo-para-medium text-primary-500"
          >
            <p>{index + 1}.</p>
            <Link
              to={`/${currentEnvironment.urlCode}/${PAGE_PATH_FEATURES}/${item.id}/targeting`}
              className="hover:underline"
            >
              {item.name}
            </Link>
          </div>
        ))}
      </div>
    </div>
  );
};

export default SegmentWarning;
