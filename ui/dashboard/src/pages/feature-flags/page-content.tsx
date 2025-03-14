import { useEffect } from 'react';
import { IconAddOutlined } from 'react-icons-material-design';
import { usePartialState } from 'hooks';
import { useTranslation } from 'i18n';
import { pickBy } from 'lodash';
import { AnyObject } from 'yup';
import { isEmptyObject, isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import Button from 'components/button';
import Icon from 'components/icon';
import { Tabs, TabsContent, TabsList, TabsTrigger } from 'components/tabs';
import Filter from 'elements/filter';
import PageLayout from 'elements/page-layout';
import CollectionLoader from './collection-loader';
import { FlagActionType, FlagTabType, FlagsTemp } from './types';

const PageContent = ({
  onAdd,
  onHandleActions
}: {
  onAdd: () => void;
  onHandleActions: (item: FlagsTemp, type: FlagActionType) => void;
}) => {
  const { t } = useTranslation(['common']);
  const { searchOptions, onChangSearchParams } = useSearchParams();
  const searchFilters: Partial<AnyObject> = searchOptions;

  const defaultFilters = {
    page: 1,
    orderBy: 'CREATED_AT',
    orderDirection: 'DESC',
    status: 'ACTIVE',
    ...searchFilters
  } as AnyObject;

  const [filters, setFilters] = usePartialState<AnyObject>(defaultFilters);
  const onChangeFilters = (values: Partial<AnyObject>) => {
    const options = pickBy({ ...filters, ...values }, v => isNotEmpty(v));
    onChangSearchParams(options);
    setFilters({ ...values });
  };

  useEffect(() => {
    if (isEmptyObject(searchOptions)) {
      setFilters({ ...defaultFilters });
    }
  }, [searchOptions]);

  return (
    <PageLayout.Content>
      <Filter
        action={
          <Button className="flex-1 lg:flex-none" onClick={onAdd}>
            <Icon icon={IconAddOutlined} size="sm" />
            {t(`new-flag`)}
          </Button>
        }
        searchValue={filters.searchQuery}
        onSearchChange={searchQuery => onChangeFilters({ searchQuery })}
      />
      <Tabs
        className="flex-1 flex h-full flex-col mt-6"
        value={filters.status}
        onValueChange={value => {
          const status = value as FlagTabType;
          onChangeFilters({
            searchQuery: '',
            status
          });
        }}
      >
        <TabsList>
          <TabsTrigger value="ACTIVE">{t(`active`)}</TabsTrigger>
          <TabsTrigger value="FAVORITES">{t(`favorites`)}</TabsTrigger>
          <TabsTrigger value="ARCHIVED">{t(`archived`)}</TabsTrigger>
        </TabsList>

        <TabsContent value={filters.status}>
          <CollectionLoader onHandleActions={onHandleActions} />
        </TabsContent>
      </Tabs>
    </PageLayout.Content>
  );
};

export default PageContent;
