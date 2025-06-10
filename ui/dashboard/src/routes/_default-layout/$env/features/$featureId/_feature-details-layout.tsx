import { featureFetcher } from '@api/features';
import { createRoute } from '@tanstack/react-router';
import { getCurrentEnvironment, meFetcher } from 'auth';
import { getOrgIdStorage } from 'storage/organization';
import { Environment, Feature } from '@types';
import FeatureFlagDetailsPage from 'pages/feature-flag-details';
import { ErrorState } from 'elements/empty-state/error';
import PageLayout from 'elements/page-layout';
import { Route as FeatureDetailsRoute } from './index';

export interface FeatureDetailsLoaderData {
  feature: Feature;
  urlCode: string;
  featureId: string;
}

export const Route = createRoute({
  id: 'feature-details-layout-route',
  loader: async ({ context, params }): Promise<FeatureDetailsLoaderData> => {
    const { featureId, env } = params || {};

    const consoleAccount = context?.auth?.consoleAccount;

    let currentEnvironment: Environment | undefined = undefined;

    currentEnvironment = consoleAccount
      ? getCurrentEnvironment(consoleAccount!)
      : undefined;

    const organizationId = getOrgIdStorage();

    if (!currentEnvironment && organizationId) {
      const response = await meFetcher(organizationId);

      const environmentRoles = response?.account?.environmentRoles;
      currentEnvironment = getCurrentEnvironment(response.account);
      currentEnvironment =
        currentEnvironment?.urlCode === env
          ? currentEnvironment
          : environmentRoles.find(
              item =>
                item.environment.urlCode === env &&
                organizationId === item.environment.organizationId
            )?.environment;
    }

    const featureData = await context.queryClient.ensureQueryData({
      queryKey: ['feature', featureId],
      queryFn: () =>
        featureFetcher({
          environmentId: currentEnvironment?.id || env,
          id: featureId
        })
    });

    return {
      feature: featureData?.feature,
      urlCode: env,
      featureId
    };
  },
  getParentRoute: () => FeatureDetailsRoute,
  component: FeatureFlagDetailsPage,
  pendingComponent: PageLayout.LoadingState,
  errorComponent: error => <ErrorState error={error} />
});
