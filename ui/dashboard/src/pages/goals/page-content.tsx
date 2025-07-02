import { useEffect } from 'react';
import { IconAddOutlined } from 'react-icons-material-design';
import { DOCUMENTATION_LINKS } from 'constants/documentation-links';
import { usePartialState } from 'hooks';
import { useTranslation } from 'i18n';
import pickBy from 'lodash/pickBy';
import { CollectionStatusType, Goal } from '@types';
import { isEmptyObject, isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import Button from 'components/button';
import Icon from 'components/icon';
import { Tabs, TabsList, TabsTrigger, TabsContent } from 'components/tabs';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import Filter from 'elements/filter';
import PageLayout from 'elements/page-layout';
import TableListContainer from 'elements/table-list-container';
import CollectionLoader from './collection-loader';
import { GoalActions, GoalFilters } from './types';

// import Overview from './overview';

const PageContent = ({
  editable,
  onAdd,
  onHandleActions
}: {
  editable: boolean;
  onAdd: () => void;
  onHandleActions: (item: Goal, type: GoalActions) => void;
}) => {
  const { t } = useTranslation(['common']);

  const { searchOptions, onChangSearchParams } = useSearchParams();
  const searchFilters: Partial<GoalFilters> = searchOptions;

  const defaultFilters = {
    page: 1,
    orderBy: 'CREATED_AT',
    orderDirection: 'DESC',
    status: 'ACTIVE',
    ...searchFilters
  } as GoalFilters;

  const [filters, setFilters] = usePartialState<GoalFilters>(defaultFilters);

  const onChangeFilters = (values: Partial<GoalFilters>) => {
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
      {/* <Overview /> */}
      <Filter
        link={DOCUMENTATION_LINKS.GOALS}
        action={
          <DisabledButtonTooltip
            hidden={editable}
            trigger={
              <Button
                className="flex-1 lg:flex-none"
                disabled={!editable}
                onClick={onAdd}
              >
                <Icon icon={IconAddOutlined} size="sm" />
                {t(`new-goal`)}
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
          onChangeFilters({
            searchQuery: '',
            status
          });
        }}
      >
        <TabsList className="px-6">
          <TabsTrigger value="ACTIVE">{t(`active`)}</TabsTrigger>
          <TabsTrigger value="ARCHIVED">{t(`archived`)}</TabsTrigger>
        </TabsList>

        <TabsContent value={filters.status} className="mt-0">
          <TableListContainer>
            <CollectionLoader
              onAdd={onAdd}
              filters={filters}
              setFilters={onChangeFilters}
              onActions={onHandleActions}
            />
          </TableListContainer>
        </TabsContent>
      </Tabs>
    </PageLayout.Content>
  );
};

export default PageContent;
