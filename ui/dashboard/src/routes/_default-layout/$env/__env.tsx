import { createRoute } from '@tanstack/react-router';
import NotFoundPage from 'pages/not-found';
import { Route as DefaultRoute } from '../../_default-layout';

export const Route = createRoute({
  path: '$env',
  getParentRoute: () => DefaultRoute,
  component: NotFoundPage
});
