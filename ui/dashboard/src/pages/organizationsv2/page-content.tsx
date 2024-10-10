import { IconAddOutlined } from 'react-icons-material-design';
import Filter from 'containers/filter';
import { usePartialState } from 'hooks';
import { useTranslation } from 'i18n';
import Button from 'components/button';
import Icon from 'components/icon';
import PageLayout from 'elements/page-layout';
import CollectionLoader from './collection-loader';

const PageContent = ({ onAdd }: { onAdd: () => void }) => {
  const { t } = useTranslation(['common']);

  const [filters, setFilters] = usePartialState({
    searchQuery: ''
  });

  return (
    <PageLayout.Content>
      <Filter
        action={
          <Button className="flex-1 lg:flex-none" onClick={onAdd}>
            <Icon icon={IconAddOutlined} size="sm" />
            {t(`new-org`)}
          </Button>
        }
        searchValue={filters.searchQuery}
        onSearchChange={v => setFilters({ ...filters, searchQuery: v })}
      />
      <CollectionLoader />
    </PageLayout.Content>
  );
};

export default PageContent;
