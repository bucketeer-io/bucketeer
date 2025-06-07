import { createRoute } from '@tanstack/react-router';
import { Route as DefaultRoute } from '../../_default-layout';

export const Route = createRoute({
  path: '$env',
  getParentRoute: () => DefaultRoute
});
