import { createRoute } from '@tanstack/react-router';
import FeatureFlagsPage from 'pages/feature-flags';
import { Route as FeaturesRoute } from '.';

export const Route = createRoute({
  path: '/clone/$featureId',
  component: FeatureFlagsPage,
  getParentRoute: () => FeaturesRoute
});
