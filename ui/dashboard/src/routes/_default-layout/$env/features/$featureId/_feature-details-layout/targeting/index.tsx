import { createRoute } from '@tanstack/react-router';
import { Route as FeatureDetailsLayoutRoute } from '../../_feature-details-layout';

export const Route = createRoute({
  path: 'targeting',
  getParentRoute: () => FeatureDetailsLayoutRoute,
  component: RouteComponent
});

function RouteComponent() {
  return (
    <div>
      Hello
      "/_default-layout/$env/features/$featureId/_feature-details-layout/targeting/"!
    </div>
  );
}
