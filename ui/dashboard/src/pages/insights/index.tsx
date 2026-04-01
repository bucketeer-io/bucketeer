import PageHeader from 'elements/page-header';
import PageLayout from 'elements/page-layout';
import PageLoader from './page-loader';

const InsightsPage = () => {
  return (
    <PageLayout.Root title="Insights">
      <PageHeader title="Insights" description="Metrics overview" />
      <PageLoader />
    </PageLayout.Root>
  );
};

export default InsightsPage;
