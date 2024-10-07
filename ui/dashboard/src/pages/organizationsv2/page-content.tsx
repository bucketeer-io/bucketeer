import PageLayout from 'elements/page-layout';
// import SubPageFooter from '~/elements/sub-page-footer';
// import SubPageHeader from '~/elements/sub-page-header';
import PageBody from './page-body';

const PageContent = () => {
  return (
    <PageLayout.Content>
      <PageLayout.Header></PageLayout.Header>

      <PageLayout.Body>
        <PageBody />
      </PageLayout.Body>
    </PageLayout.Content>
  );
};

export default PageContent;
