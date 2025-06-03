import { memo } from 'react';
import { createRoute, Outlet, redirect } from '@tanstack/react-router';
import Navigation from 'components/navigation';
import { Route as RootRoute } from '../__root';

export const EnvironmentRoot = memo(() => {
  return (
    <div className="flex size-full">
      <Navigation onClickNavLink={() => {}} />
      <div className="w-full ml-[248px] shadow-lg overflow-y-auto">
        <Outlet />
      </div>
    </div>
  );
});

export const Route = createRoute({
  path: '$envUrlCode',
  // beforeLoad: async ({ params, location }) => {
  //   console.log(location.pathname);
  //   if (location.pathname === `/v3/${params.envUrlCode}`) {
  //     throw redirect({ to: '/$envUrlCode/features', params });
  //   }
  // },
  getParentRoute: () => RootRoute,
  component: EnvironmentRoot
});
