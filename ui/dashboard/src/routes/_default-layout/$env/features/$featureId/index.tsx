import { createRoute, redirect } from '@tanstack/react-router';
import { Route as FeaturesRoute } from '../index';

export const Route = createRoute({
  path: '$featureId',
  beforeLoad: ({ params }) => {
    throw redirect({
      to: '/$env/features/$featureId/targeting',
      params
    });
  },
  getParentRoute: () => FeaturesRoute
});
