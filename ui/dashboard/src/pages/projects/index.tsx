import PageHeader from 'containers/page-header';
import { ProjectsContent } from 'containers/pages';

const ProjectsPage = () => {
  return (
    <div className="flex flex-col size-full overflow-auto">
      <PageHeader
        title="Projects"
        description="Manage all projects for this environment."
      />
      <ProjectsContent />
    </div>
  );
};

export default ProjectsPage;
