import { Route, Routes } from 'react-router-dom';
import OrganizationDetailPage from 'pages/organization-details';
import OrganizationsPage from 'pages/organizationsv2';

export const OrganizationsRoot = () => {
  return (
    <Routes>
      <Route index element={<OrganizationsPage />} />
      <Route path=":organizationId" element={<OrganizationDetailPage />} />
    </Routes>
  );
};
