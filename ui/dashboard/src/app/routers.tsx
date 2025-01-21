import { Route, Routes } from 'react-router-dom';
import ExperimentDetailsPage from 'pages/experiment-details';
import ExperimentsPage from 'pages/experiments';
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

export const ExperimentsRoot = () => {
  return (
    <Routes>
      <Route index element={<ExperimentsPage />} />
      <Route path="new" element={<ExperimentsPage />} />
      <Route path=":experimentId" element={<ExperimentsPage />} />
      <Route path=":experimentId/:tab" element={<ExperimentDetailsPage />} />
    </Routes>
  );
};
