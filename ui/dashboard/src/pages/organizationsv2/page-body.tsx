// import { usePartialState } from 'hooks/use-partial-state';
import CollectionLoader from './collection-loader';

// import type { CompaniesFilters } from './types';

const PageBody = () =>
  // { onAdd }: { onAdd: () => void }
  {
    // const [filters, setFilters] = usePartialState<CompaniesFilters>({
    //   searchQuery: ''
    // });

    return (
      <CollectionLoader
      // filters={filters} setFilters={setFilters} onAdd={onAdd}
      />
    );
  };

export default PageBody;
