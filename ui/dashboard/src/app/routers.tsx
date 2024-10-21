import { Navigate, Route, Routes } from 'react-router-dom';
import { PAGE_PATH_ORGANIZATIONS } from 'constants/routing';
import OrganizationDetailPage from 'pages/organization-details';
import OrganizationsPage from 'pages/organizations';

export const OrganizationsRoot = () => {
  return (
    <Routes>
      <Route
        index
        element={<Navigate to={`${PAGE_PATH_ORGANIZATIONS}/active`} replace />}
      />
      <Route path=":organizationStatus/*" element={<OrganizationsPage />} />
      <Route
        path=":organizationStatus/:organizationId/*"
        element={<OrganizationDetailPage />}
      />
    </Routes>
  );
};
