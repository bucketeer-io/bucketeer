import PageLayout from 'elements/page-layout';
import FlagForm from './flag-form';

const PageContent = () => {
  return (
    <PageLayout.Content className="p-6 min-w-[900px]">
      <FlagForm />
    </PageLayout.Content>
  );
};

export default PageContent;
