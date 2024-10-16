import PageHeader from 'elements/page-header';

const FeatureFlagsPage = () => {
  return (
    <div className="flex flex-col size-full overflow-auto">
      <PageHeader
        title="Feature Flags"
        description="Select a flag to manage the environment-specific targeting and rollout rules"
      />
    </div>
  );
};

export default FeatureFlagsPage;
