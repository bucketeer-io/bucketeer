import PageHeader from 'containers/page-header';
import { OrganizationsContent } from 'containers/pages';

const OrganizationsPage = () => {
  return (
    <div className="flex flex-col size-full overflow-auto">
      <PageHeader
        title="Organizations"
        description="You can see all your clients data"
      />
      <OrganizationsContent />
    </div>
  );
};

export default OrganizationsPage;
