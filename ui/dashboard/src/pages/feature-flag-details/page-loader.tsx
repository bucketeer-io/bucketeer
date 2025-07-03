import { useNavigate, useParams } from 'react-router-dom';
import { useQueryFeature } from '@queries/feature-details';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_FEATURES } from 'constants/routing';
import { useTranslation } from 'i18n';
import { FlagStatus } from 'pages/feature-flags/collection-layout/elements';
import { getFlagStatus } from 'pages/feature-flags/collection-layout/elements/utils';
import { Tooltip } from 'components/tooltip';
import PageDetailsHeader from 'elements/page-details-header';
import PageLayout from 'elements/page-layout';
import PageContent from './page-content';

const PageLoader = () => {
  const { t } = useTranslation(['table']);
  const params = useParams();
  const navigate = useNavigate();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const {
    data: collection,
    isLoading,
    isError,
    refetch
  } = useQueryFeature({
    params: {
      environmentId: currentEnvironment?.id,
      id: params?.flagId || ''
    },
    enabled:
      !!params?.flagId &&
      (!!currentEnvironment?.id || !!currentEnvironment?.urlCode),
    gcTime: 0
  });

  const feature = collection?.feature;
  const isErrorState = isError || !feature;

  if (isLoading) return <PageLayout.LoadingState />;
  if (isErrorState) return <PageLayout.ErrorState onRetry={refetch} />;

  const flagStatus = getFlagStatus(feature);

  return (
    <>
      <PageDetailsHeader
        onBack={() =>
          navigate(`/${currentEnvironment?.urlCode}${PAGE_PATH_FEATURES}`)
        }
        title={feature.name}
        additionElement={
          <>
            <Tooltip
              asChild={false}
              align="start"
              trigger={<FlagStatus status={flagStatus} />}
              content={t(
                `feature-flags.${flagStatus === 'receiving-traffic' ? 'receiving-traffic-description' : flagStatus === 'no-recent-traffic' ? 'no-recent-traffic-description' : 'never-used-description'}`
              )}
            />
          </>
        }
      />
      <PageContent feature={feature} />
    </>
  );
};

export default PageLoader;
