import { createRoute } from '@tanstack/react-router';
import Operations from 'pages/feature-flag-details/operations';
import { ErrorState } from 'elements/empty-state/error';
import { Route as FeatureDetailsLayoutRoute } from '../../_feature-details-layout';

export const Route = createRoute({
  path: 'operations',
  getParentRoute: () => FeatureDetailsLayoutRoute,
  component: Operations,
  errorComponent: error => <ErrorState error={error} />
});
