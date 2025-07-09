import { useCallback } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useQueryGoalDetails } from '@queries/goal-details';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_GOALS } from 'constants/routing';
import { useTranslation } from 'i18n';
import PageDetailsHeader from 'elements/page-details-header';
import CreatedAtTime from 'elements/page-details-header/created-at-time';
import HeaderDetailsID from 'elements/page-details-header/header-details-id';
import PageLayout from 'elements/page-layout';
import Status from 'elements/status';
import PageContent from './page-content';

const PageLoader = () => {
  const navigate = useNavigate();
  const { t } = useTranslation(['table', 'message', 'common']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const { goalId } = useParams();

  const handleBack = useCallback(
    () => navigate(`/${currentEnvironment.urlCode}${PAGE_PATH_GOALS}`),
    [currentEnvironment]
  );

  const { data, isLoading, refetch, isError } = useQueryGoalDetails({
    params: {
      id: goalId!,
      environmentId: currentEnvironment.id
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
            title={goal.name}
            additionElement={
              <>
                <Status
                  text={t(
                    `common:${goal.isInUseStatus ? 'in-use' : 'not-in-use'}`
                  )}
                  isInUseStatus={goal.isInUseStatus}
                />
                <CreatedAtTime createdAt={goal.createdAt} />
              </>
            }
            onBack={handleBack}
          >
            <HeaderDetailsID id={goal.id} />
          </PageDetailsHeader>
          <PageContent goal={goal} />
        </>
      )}
    </>
  );
};

export default PageLoader;
