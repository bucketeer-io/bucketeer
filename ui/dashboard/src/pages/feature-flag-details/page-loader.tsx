import { useLoaderData, useNavigate } from '@tanstack/react-router';
import { PAGE_PATH_FEATURES } from 'constants/routing';
import { useTranslation } from 'i18n';
import {
  Route as featureDetailsLayoutRoute,
  FeatureDetailsLoaderData
} from 'routes/_default-layout/$env/features/$featureId/_feature-details-layout';
import { FlagStatus } from 'pages/feature-flags/collection-layout/elements';
import { getFlagStatus } from 'pages/feature-flags/collection-layout/elements/utils';
import { Tooltip } from 'components/tooltip';
import PageDetailsHeader from 'elements/page-details-header';
import PageContent from './page-content';

const PageLoader = () => {
  const { t } = useTranslation(['table']);
  const navigate = useNavigate();

  const loaderData: FeatureDetailsLoaderData = useLoaderData({
    from: featureDetailsLayoutRoute.id
  });

  const feature = loaderData?.feature;
  const flagStatus = getFlagStatus(feature);

  return (
    <>
      <PageDetailsHeader
        onBack={() =>
          navigate({
            to: `${PAGE_PATH_FEATURES}`
          })
        }
        title={feature.name}
        additionElement={
          <>
            <Tooltip
              asChild={false}
              align="start"
              trigger={<FlagStatus status={flagStatus} />}
              content={t(
                `feature-flags.${flagStatus === 'active' ? 'active-description' : flagStatus === 'in-active' ? 'inactive-description' : 'new-description'}`
              )}
            />
          </>
        }
      />
      <PageContent />
    </>
  );
};

export default PageLoader;
