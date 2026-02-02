import PageLayout from 'elements/page-layout';
import FlagForm from './flag-form';

const PageContent = () => {
  return (
    <PageLayout.Content className="p-3 sm:p-6">
      <FlagForm />
    </PageLayout.Content>
  );
};

export default PageContent;
