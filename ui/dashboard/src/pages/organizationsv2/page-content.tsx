import PageLayout from 'elements/page-layout';
// import SubPageFooter from '~/elements/sub-page-footer';
// import SubPageHeader from '~/elements/sub-page-header';
import PageBody from './page-body';

const PageContent = () =>
  // { onAdd }: { onAdd: () => void }
  {
    return (
      <PageLayout.Content>
        <PageLayout.Header>
          {/* <SubPageHeader.Root variant="title">
          <SubPageHeader.Content>
            <SubPageHeader.TitleText>{`Organization management`}</SubPageHeader.TitleText>
            <SubPageHeader.ActionButton variant="primary" onPress={onAdd}>
              <IconAddOutlined />
              {`Add organization`}
            </SubPageHeader.ActionButton>
          </SubPageHeader.Content>
        </SubPageHeader.Root> */}
        </PageLayout.Header>

        <PageLayout.Body>
          <PageBody
          // onAdd={onAdd}
          />
        </PageLayout.Body>
      </PageLayout.Content>
    );
  };

export default PageContent;
