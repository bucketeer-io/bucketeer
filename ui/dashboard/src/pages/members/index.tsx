import PageHeader from 'elements/page-header';

const ProjectsPage = () => {
  return (
    <div className="flex flex-col size-full overflow-auto">
      <PageHeader
        title="Projects"
        description="Manage all projects for this environment."
      />
    </div>
  );
};

export default ProjectsPage;
