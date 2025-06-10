import { createRoute } from '@tanstack/react-router';
import FeatureFlagsPage from 'pages/feature-flags';
import { Route as EnvRoute } from '../__env';

export const Route = createRoute({
  id: 'features',
  component: FeatureFlagsPage,
  getParentRoute: () => EnvRoute
});
