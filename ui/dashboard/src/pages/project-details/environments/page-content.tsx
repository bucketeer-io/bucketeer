import { useEffect } from 'react';
import { IconAddOutlined } from 'react-icons-material-design';
import { environmentArchive, environmentUnArchive } from '@api/environment';
import { invalidateEnvironments } from '@queries/environments';
import { useQueryClient } from '@tanstack/react-query';
import { usePartialState } from 'hooks';
import { useTranslation } from 'i18n';
import pickBy from 'lodash/pickBy';
import { CollectionStatusType, Environment } from '@types';
import { isEmptyObject, isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import Button from 'components/button';
import Icon from 'components/icon';
import { Tabs, TabsContent, TabsList, TabsTrigger } from 'components/tabs';
import Filter from 'elements/filter';
import CollectionLoader from './collection-loader';
import { EnvironmentFilters } from './types';

const PageContent = ({
  onAdd,
  onEdit
}: {
  onAdd: () => void;
  onEdit: (v: Environment) => void;
}) => {
  const queryClient = useQueryClient();
  const { t } = useTranslation(['common']);
  const { searchOptions, onChangSearchParams } = useSearchParams();
  const searchFilters: Partial<EnvironmentFilters> = searchOptions;

  const defaultFilters = {
    page: 1,
    orderBy: 'CREATED_AT',
    orderDirection: 'DESC',
    status: 'ACTIVE',
    ...searchFilters
  } as EnvironmentFilters;

  const [filters, setFilters] =
    usePartialState<EnvironmentFilters>(defaultFilters);

  const onChangeFilters = (values: Partial<EnvironmentFilters>) => {
    const options = pickBy({ ...filters, ...values }, v => isNotEmpty(v));
    onChangSearchParams(options);
    setFilters({ ...values });
  };

  const onArchivedEnvironment = (environment: Environment) => {
    environmentArchive({
      id: environment.id,
      command: {}
    }).then(() => {
      invalidateEnvironments(queryClient);
    });
  };

  const onUnArchiveEnvironment = (environment: Environment) => {
    environmentUnArchive({
      id: environment.id,
      command: {}
    }).then(() => {
      invalidateEnvironments(queryClient);
    });
  };

  const onActionHandler = (type: string, environment: Environment) => {
    if (type === 'ARCHIVED_ENVIRONMENT') {
      onArchivedEnvironment(environment);
    } else if (type === 'UNARCHIVE_ENVIRONMENT') {
      onUnArchiveEnvironment(environment);
    } else {
      onEdit(environment);
    }
  };

  useEffect(() => {
    if (isEmptyObject(searchOptions)) {
      setFilters({ ...defaultFilters });
    }
  }, [searchOptions]);

  return (
    <>
      <Filter
        action={
          <Button className="flex-1 lg:flex-none" onClick={onAdd}>
            <Icon icon={IconAddOutlined} size="sm" />
            {t(`new-env`)}
          </Button>
        }
        searchValue={filters.searchQuery}
        onSearchChange={searchQuery => onChangeFilters({ searchQuery })}
      />
      <Tabs
        className="flex-1 flex h-full flex-col mt-6"
        value={filters.status}
        onValueChange={value => {
          const status = value as CollectionStatusType;
          onChangeFilters({ status, searchQuery: '' });
        }}
      >
        <TabsList>
          <TabsTrigger value="ACTIVE">{t(`active`)}</TabsTrigger>
          <TabsTrigger value="ARCHIVED">{t(`archived`)}</TabsTrigger>
        </TabsList>

        <TabsContent value={filters.status}>
          <CollectionLoader
            onAdd={onAdd}
            filters={filters}
            setFilters={onChangeFilters}
            onActionHandler={onActionHandler}
          />
        </TabsContent>
      </Tabs>
    </>
  );
};

export default PageContent;
