import { createRoute } from '@tanstack/react-router';
import { Route as EnvRoute } from './__env';

export const Route = createRoute({
  path: 'settings',
  getParentRoute: () => EnvRoute,
  component: () => <div>Settings</div>
});
