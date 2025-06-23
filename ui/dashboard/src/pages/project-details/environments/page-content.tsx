import { useEffect } from 'react';
import { IconAddOutlined } from 'react-icons-material-design';
import { useAuthAccess } from 'auth';
import { usePartialState } from 'hooks';
import { useTranslation } from 'i18n';
import pickBy from 'lodash/pickBy';
import { CollectionStatusType, Environment } from '@types';
import { isEmptyObject, isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import Button from 'components/button';
import Icon from 'components/icon';
import { Tabs, TabsContent, TabsList, TabsTrigger } from 'components/tabs';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import Filter from 'elements/filter';
import CollectionLoader from './collection-loader';
import { EnvironmentActionsType, EnvironmentFilters } from './types';

const PageContent = ({
  onAdd,
  onActionHandler
}: {
  onAdd: () => void;
  onActionHandler: (item: Environment, type: EnvironmentActionsType) => void;
}) => {
  const { t } = useTranslation(['common']);
  const { envEditable, isOrganizationAdmin } = useAuthAccess();
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

  useEffect(() => {
    if (isEmptyObject(searchOptions)) {
      setFilters({ ...defaultFilters });
    }
  }, [searchOptions]);

  return (
    <>
      <Filter
        action={
          <DisabledButtonTooltip
            type={!isOrganizationAdmin ? 'admin' : 'editor'}
            hidden={envEditable && isOrganizationAdmin}
            trigger={
              <Button
                className="flex-1 lg:flex-none"
                onClick={onAdd}
                disabled={!envEditable || !isOrganizationAdmin}
              >
                <Icon icon={IconAddOutlined} size="sm" />
                {t(`new-env`)}
              </Button>
            }
          />
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
        <TabsList className="px-6">
          <TabsTrigger value="ACTIVE">{t(`active`)}</TabsTrigger>
          <TabsTrigger value="ARCHIVED">{t(`archived`)}</TabsTrigger>
        </TabsList>

        <TabsContent
          value={filters.status}
          className="px-6 pb-6 overflow-y-hidden overflow-x-auto"
        >
          <CollectionLoader
            onAdd={onAdd}
            filters={filters}
            setFilters={onChangeFilters}
            onActions={onActionHandler}
          />
        </TabsContent>
      </Tabs>
    </>
  );
};

export default PageContent;
