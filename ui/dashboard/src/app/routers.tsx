import { Route, Routes } from 'react-router-dom';
import { ID_CLONE, ID_NEW } from 'constants/routing';
import FeatureFlagDetailsPage from 'pages/feature-flag-details';
import FeatureFlagsPage from 'pages/feature-flags';
import GoalDetailsPage from 'pages/goal-details';
import GoalsPage from 'pages/goals';
import NotFoundPage from 'pages/not-found';
import OrganizationDetailPage from 'pages/organization-details';
import OrganizationsPage from 'pages/organizations';
import ProjectDetailsPage from 'pages/project-details';
import ProjectsPage from 'pages/projects';

export const OrganizationsRoot = () => {
  return (
    <Routes>
      <Route index element={<OrganizationsPage />} />
      <Route path=":organizationId/*" element={<OrganizationDetailPage />} />
    </Routes>
  );
};

export const ProjectsRoot = () => {
  return (
    <Routes>
      <Route index element={<ProjectsPage />} />
      <Route path=":projectId/*" element={<ProjectDetailsPage />} />
    </Routes>
  );
};

export const GoalsRoot = () => {
  return (
    <Routes>
      <Route index element={<GoalsPage />} />
      <Route path={ID_NEW} element={<GoalsPage />} />
      <Route path=":goalId/*" element={<GoalDetailsPage />} />
    </Routes>
  );
};

export const FeatureFlagsRoot = () => {
  return (
    <Routes>
      <Route index element={<FeatureFlagsPage />} />
      <Route path={ID_NEW} element={<FeatureFlagsPage />} />
      <Route path={`${ID_CLONE}/:flagId`} element={<FeatureFlagsPage />} />
      <Route path=":flagId/:tab" element={<FeatureFlagDetailsPage />} />
      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  );
};
