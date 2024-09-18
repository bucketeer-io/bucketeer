import PageHeader from 'containers/page-header';

const SettingsPage = () => {
  return (
    <div className="flex flex-col size-full overflow-auto">
      <PageHeader
        title="Settings"
        description="You can see all your clients data"
      />
    </div>
  );
};

export default SettingsPage;
