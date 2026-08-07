import { useCallback, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router';
import { useQueryUserSegment } from '@queries/user-segment-details';
import { getCurrentEnvironment, hasEditable, useAuth } from 'auth';
import { PAGE_PATH_USER_SEGMENTS } from 'constants/routing';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import FormLoading from 'elements/form-loading';
import PageDetailsHeader from 'elements/page-details-header';
import PageLayout from 'elements/page-layout';
import PageContent from './page-content';

const PageLoader = () => {
  const { t } = useTranslation(['common']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const editable = hasEditable(consoleAccount!);
  const navigate = useNavigate();
  const { segmentId } = useParams();
  const isUpdate = !!segmentId;
  const { errorNotify } = useToast();

  const {
    data: segmentCollection,
    isLoading,
    isError,
    error
  } = useQueryUserSegment({
    params: {
      environmentId: currentEnvironment.id,
      id: segmentId as string
    },
    enabled: isUpdate
  });
  const userSegment = segmentCollection?.segment;

  const onBack = useCallback(
    () => navigate(`/${currentEnvironment.urlCode}${PAGE_PATH_USER_SEGMENTS}`),
    [currentEnvironment]
  );

  useEffect(() => {
    if (isError && error) {
      errorNotify(error);
      onBack();
    }
  }, [isError, error]);

  return (
    <div className="w-full min-h-screen !max-w-[1192px]">
      <PageDetailsHeader
        onBack={onBack}
        title={t(isUpdate ? 'update-user-segment' : 'new-user-segment')}
      />
      {isUpdate && (isLoading || !userSegment) ? (
        <PageLayout.Content className="p-6">
          <FormLoading />
        </PageLayout.Content>
      ) : (
        <PageContent
          isUpdate={isUpdate}
          userSegment={userSegment}
          isDisabled={!editable}
        />
      )}
    </div>
  );
};

export default PageLoader;
