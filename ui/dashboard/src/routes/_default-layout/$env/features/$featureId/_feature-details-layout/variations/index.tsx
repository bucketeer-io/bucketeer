import { createRoute } from '@tanstack/react-router';
import Variations from 'pages/feature-flag-details/variation';
import { ErrorState } from 'elements/empty-state/error';
import { Route as FeatureDetailsLayoutRoute } from '../../_feature-details-layout';

export const Route = createRoute({
  path: 'variations',
  getParentRoute: () => FeatureDetailsLayoutRoute,
  component: Variations,
  errorComponent: error => <ErrorState error={error} />
});
