import { OrganizationsFetcherParams } from '@api/organization';
import { useQueryOrganizations } from '@queries/organizations';
import { LIST_PAGE_SIZE } from 'constants/app';
import { useToggleOpen } from 'hooks/use-toggle-open';
import PageLayout from 'elements/page-layout';
import { EmptyCollection } from './collection-layout/empty-collection';
// import { useFetchCompanies } from './collection-loader/use-fetch-organizations';
import PageContent from './page-content';

const PageLoader = () => {
  const defaultParams: OrganizationsFetcherParams = {
    pageSize: LIST_PAGE_SIZE,
    cursor: String(0),
    orderBy: 'DEFAULT',
    orderDirection: 'ASC',
    searchKeyword: '',
    disabled: false,
    archived: false
  };

  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useQueryOrganizations({
    params: defaultParams
  });

  const [, onOpenAddModal] = useToggleOpen(false);

  if (isLoading) {
    return <PageLayout.LoadingState />;
  }

  if (isError) {
    return <PageLayout.ErrorState onRetry={refetch} />;
  }

  const isEmpty = collection?.Organizations.length === 0;

  if (!isEmpty) {
    return (
      <PageLayout.EmptyState>
        <EmptyCollection onAdd={onOpenAddModal} />
      </PageLayout.EmptyState>
    );
  }
  return (
    <>
      <PageContent />
    </>
  );
};

export default PageLoader;
