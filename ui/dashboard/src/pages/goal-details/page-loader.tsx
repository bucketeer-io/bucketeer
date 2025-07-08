import { useNavigate, useParams } from 'react-router-dom';
import { useQueryGoalDetails } from '@queries/goal-details';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_GOALS } from 'constants/routing';
import { useTranslation } from 'i18n';
import { useFormatDateTime } from 'utils/date-time';
import { checkEnvironmentEmptyId } from 'utils/function';
import PageDetailsHeader from 'elements/page-details-header';
import PageLayout from 'elements/page-layout';
import HeaderDetails from './elements/header-details';
import PageContent from './page-content';

const PageLoader = () => {
  const navigate = useNavigate();
  const formatDateTime = useFormatDateTime();
  const { t } = useTranslation(['table']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const { goalId } = useParams();

  const { data, isLoading, refetch, isError } = useQueryGoalDetails({
    params: {
      id: goalId!,
      environmentId: checkEnvironmentEmptyId(currentEnvironment.id)
    }
  });
  const goal = data?.goal;
  const isErrorState = isError || !goal;

  return (
    <>
      {isLoading ? (
        <PageLayout.LoadingState />
      ) : isErrorState ? (
        <PageLayout.ErrorState onRetry={refetch} />
      ) : (
        <>
          <PageDetailsHeader
            description={t('created-at-time', {
              time: formatDateTime(goal.createdAt)
            })}
            onBack={() =>
              navigate(`/${currentEnvironment.urlCode}${PAGE_PATH_GOALS}`)
            }
          >
            <HeaderDetails goal={goal} />
          </PageDetailsHeader>
          <PageContent goal={goal} />
        </>
      )}
    </>
  );
};

export default PageLoader;
