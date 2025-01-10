import { useNavigate, useParams } from 'react-router-dom';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_GOALS } from 'constants/routing';
import { useToast } from 'hooks';
import { truncateTextCenter } from 'utils/converts';
import { useFormatDateTime } from 'utils/date-time';
import { copyToClipBoard } from 'utils/function';
import { cn } from 'utils/style';
import { IconCopy } from '@icons';
import { mocks } from 'pages/goals/page-loader';
import Icon from 'components/icon';
import PageDetailsHeader from 'elements/page-details-header';
import PageLayout from 'elements/page-layout';
import Status from 'elements/status';

const PageLoader = () => {
  const navigate = useNavigate();
  const formatDateTime = useFormatDateTime();
  const { notify } = useToast();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const { goalId } = useParams();
  const isLoading = false;
  const isError = false;

  const goal = mocks.find(item => item.id === goalId);
  const isErrorState = isError || !goal;

  const handleCopyId = (id: string) => {
    copyToClipBoard(id);
    notify({
      toastType: 'toast',
      messageType: 'success',
      message: (
        <span>
          <b>ID</b> {` has been successfully copied!`}
        </span>
      )
    });
  };

  return (
    <>
      {isLoading ? (
        <PageLayout.LoadingState />
      ) : isErrorState ? (
        <PageLayout.ErrorState onRetry={() => {}} />
      ) : (
        <>
          <PageDetailsHeader
            title={goal.name}
            description={`Created ${formatDateTime(goal.createdAt)}`}
            onBack={() =>
              navigate(`/${currentEnvironment.urlCode}${PAGE_PATH_GOALS}`)
            }
          >
            <div className="flex flex-col w-full gap-y-4 mt-3">
              <div className="flex items-center w-full gap-x-2">
                <h1 className="typo-head-bold-huge leading-6 text-gray-900">
                  {goal.name}
                </h1>
                <Status
                  status={goal?.isInUseStatus ? 'In Use' : 'Not In Use'}
                  className={cn({
                    'bg-accent-green-50 text-accent-green-500':
                      goal.isInUseStatus
                  })}
                />
              </div>
              <div className="flex items-center h-5 gap-x-2 typo-para-tiny text-gray-500 select-none">
                {truncateTextCenter(goal.id)}
                <div onClick={() => handleCopyId(goal.id)}>
                  <Icon
                    icon={IconCopy}
                    size={'sm'}
                    className="opacity-100 cursor-pointer"
                  />
                </div>
              </div>
            </div>
          </PageDetailsHeader>
          {/* <PageContent organization={organization} /> */}
        </>
      )}
    </>
  );
};

export default PageLoader;
