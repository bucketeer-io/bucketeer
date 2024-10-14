import { useCallback } from 'react';
import { IconAddOutlined } from 'react-icons-material-design';
import { useNavigate, useLocation } from 'react-router-dom';
import Filter from 'containers/filter';
import { usePartialState } from 'hooks';
import { useTranslation } from 'i18n';
import pickBy from 'lodash/pickby';
import { CollectionStatusType, OrderBy, OrderDirection } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifySearchParams, useSearchParams } from 'utils/search-params';
import Button from 'components/button';
import Icon from 'components/icon';
import { Tabs, TabsList, TabsTrigger, TabsContent } from 'components/tabs';
import PageLayout from 'elements/page-layout';
import CollectionLoader from './collection-loader';
import { OrganizationFilters } from './types';

const PageContent = ({ onAdd }: { onAdd: () => void }) => {
  const { t } = useTranslation(['common']);
  const navigate = useNavigate();
  const { pathname } = useLocation();

  // NOTE: Need improve search options
  const searchOptions = useSearchParams();

  const [filters, setFilters] = usePartialState<OrganizationFilters>({
    orderBy: (searchOptions.orderBy as OrderBy) || 'DEFAULT',
    orderDirection: (searchOptions.orderDirection as OrderDirection) || 'ASC',
    searchQuery: (searchOptions.searchQuery as string) || '',
    status: (searchOptions.status as CollectionStatusType) || 'ACTIVE'
  });

  const onUpdateURL = useCallback(
    (options: Record<string, string | number | boolean>) => {
      navigate(`${pathname}?${stringifySearchParams(options)}`, {
        replace: true
      });
    },
    [navigate]
  );

  const onChangeFilters = (values: Partial<OrganizationFilters>) => {
    const options = pickBy({ ...filters, ...values }, v => isNotEmpty(v));
    onUpdateURL(options);
    setFilters({ ...values });
  };

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
        onSearchChange={searchQuery => onChangeFilters({ searchQuery })}
      />
      <Tabs
        defaultValue={filters.status}
        onValueChange={value =>
          onChangeFilters({ status: value as CollectionStatusType })
        }
      >
        <TabsList>
          <TabsTrigger value="ACTIVE">{t(`active`)}</TabsTrigger>
          <TabsTrigger value="ARCHIVED">{t(`archived`)}</TabsTrigger>
        </TabsList>
        <TabsContent value={filters.status}>
          <CollectionLoader filters={filters} setFilters={onChangeFilters} />
        </TabsContent>
      </Tabs>
    </PageLayout.Content>
  );
};

export default PageContent;
