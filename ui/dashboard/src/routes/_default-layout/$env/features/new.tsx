import { createRoute } from '@tanstack/react-router';
import CreateFlagPage from 'pages/create-flag';
import { Route as FeaturesRoute } from '.';

export const Route = createRoute({
  path: 'new',
  component: CreateFlagPage,
  getParentRoute: () => FeaturesRoute
});
