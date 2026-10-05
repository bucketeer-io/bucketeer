import { RouteObject, useRoutes } from 'react-router';
import { ID_CLONE, ID_NEW } from 'constants/routing';
import CreateFlagPage from 'pages/create-flag';
import CreateUpdateSegmentPage from 'pages/create-update-segment';
import ExperimentDetailsPage from 'pages/experiment-details';
import ExperimentsPage from 'pages/experiments';
import FeatureFlagDetailsPage from 'pages/feature-flag-details';
import FeatureFlagsPage from 'pages/feature-flags';
import GoalDetailsPage from 'pages/goal-details';
import GoalsPage from 'pages/goals';
import MembersPage from 'pages/members';
import OrganizationDetailPage from 'pages/organization-details';
import OrganizationsPage from 'pages/organizations';
import ProjectDetailsPage from 'pages/project-details';
import ProjectsPage from 'pages/projects';
import UserSegmentsPage from 'pages/user-segments';

export const organizationsRoutes: RouteObject[] = [
  { index: true, element: <OrganizationsPage /> },
  { path: 'new', element: <OrganizationsPage /> },
  { path: ':organizationId', element: <OrganizationsPage /> },
  { path: ':organizationId/*', element: <OrganizationDetailPage /> }
];

export const projectsRoutes: RouteObject[] = [
  { index: true, element: <ProjectsPage /> },
  { path: 'new', element: <ProjectsPage /> },
  { path: ':projectId', element: <ProjectsPage /> },
  { path: ':projectId/*', element: <ProjectDetailsPage /> }
];

export const experimentsRoutes: RouteObject[] = [
  { index: true, element: <ExperimentsPage /> },
  { path: 'new', element: <ExperimentsPage /> },
  { path: ':experimentId', element: <ExperimentsPage /> },
  { path: ':experimentId/:tab', element: <ExperimentDetailsPage /> }
];

export const goalsRoutes: RouteObject[] = [
  { index: true, element: <GoalsPage /> },
  { path: 'new', element: <GoalsPage /> },
  { path: ':goalId/*', element: <GoalDetailsPage /> }
];

export const featureFlagsRoutes: RouteObject[] = [
  { index: true, element: <FeatureFlagsPage /> },
  { path: ID_NEW, element: <CreateFlagPage /> },
  { path: `${ID_CLONE}/:flagId`, element: <FeatureFlagsPage /> },
  { path: ':flagId/*', element: <FeatureFlagDetailsPage /> }
];

export const userSegmentsRoutes: RouteObject[] = [
  { index: true, element: <UserSegmentsPage /> },
  { path: ID_NEW, element: <CreateUpdateSegmentPage /> },
  { path: ':segmentId', element: <CreateUpdateSegmentPage /> }
];

export const memberRoutes: RouteObject[] = [
  { index: true, element: <MembersPage /> },
  { path: ID_NEW, element: <MembersPage /> },
  { path: ':memberId', element: <MembersPage /> }
];

export const OrganizationsRoot = () => useRoutes(organizationsRoutes);
export const ProjectsRoot = () => useRoutes(projectsRoutes);
export const ExperimentsRoot = () => useRoutes(experimentsRoutes);
export const GoalsRoot = () => useRoutes(goalsRoutes);
export const FeatureFlagsRoot = () => useRoutes(featureFlagsRoutes);
export const UserSegmentsRoot = () => useRoutes(userSegmentsRoutes);
export const MemberRoot = () => useRoutes(memberRoutes);
