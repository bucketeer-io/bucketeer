import { useCallback, useEffect } from 'react';
import { IconAddOutlined } from 'react-icons-material-design';
import { getAccountAccess, getCurrentEnvironment, useAuth } from 'auth';
import { usePartialState, useToggleOpen } from 'hooks';
import { useTranslation } from 'i18n';
import pickBy from 'lodash/pickBy';
import { Project } from '@types';
import { isEmptyObject, isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import Button from 'components/button';
import Icon from 'components/icon';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import Filter from 'elements/filter';
import PageLayout from 'elements/page-layout';
import TableListContainer from 'elements/table-list-container';
import CollectionLoader from './collection-loader';
import FilterProjectModal from './project-modal/filter-project-modal';
import { ProjectFilters } from './types';

const PageContent = ({
  onAdd,
  onEdit
}: {
  onAdd: () => void;
  onEdit: (v: Project) => void;
}) => {
  const { t } = useTranslation(['common']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { envEditable, isOrganizationAdmin } = getAccountAccess(
    consoleAccount!
  );

  const { searchOptions, onChangSearchParams } = useSearchParams();
  const searchFilters: Partial<ProjectFilters> = searchOptions;

  const defaultFilters = {
    page: 1,
    orderBy: 'CREATED_AT',
    orderDirection: 'DESC',
    ...searchFilters
  } as ProjectFilters;

  const [openFilterModal, onOpenFilterModal, onCloseFilterModal] =
    useToggleOpen(false);

  const [filters, setFilters] = usePartialState<ProjectFilters>(defaultFilters);

  const onChangeFilters = useCallback(
    (values: Partial<ProjectFilters>) => {
      const options = pickBy({ ...filters, ...values }, v => isNotEmpty(v));
      onChangSearchParams(options);
      setFilters({ ...values });
    },
    [filters]
  );

  const onActionHandler = useCallback((project: Project) => {
    onEdit(project);
  }, []);

  const onClearFilters = useCallback(
    () => setFilters({ searchQuery: '', disabled: undefined }),
    []
  );

  useEffect(() => {
    if (isEmptyObject(searchOptions)) {
      setFilters({ ...defaultFilters });
    }
  }, [searchOptions]);

  return (
    <PageLayout.Content>
      <Filter
        isShowDocumentation={false}
        onOpenFilter={onOpenFilterModal}
        action={
          <DisabledButtonTooltip
            hidden={envEditable && isOrganizationAdmin}
            type={!isOrganizationAdmin ? 'admin' : 'editor'}
            trigger={
              <Button
                className="flex-1 lg:flex-none"
                onClick={onAdd}
                disabled={!envEditable || !isOrganizationAdmin}
              >
                <Icon icon={IconAddOutlined} size="sm" />
                {t(`new-project`)}
              </Button>
            }
          />
        }
        searchValue={filters.searchQuery}
        filterCount={isNotEmpty(filters.disabled) ? 1 : undefined}
        onSearchChange={searchQuery => onChangeFilters({ searchQuery })}
      />
      {openFilterModal && (
        <FilterProjectModal
          isOpen={openFilterModal}
          filters={filters}
          onClose={onCloseFilterModal}
          onSubmit={value => {
            onChangeFilters(value);
            onCloseFilterModal();
          }}
          onClearFilters={() => {
            onChangeFilters({ disabled: undefined });
            onCloseFilterModal();
          }}
        />
      )}
      <TableListContainer>
        <CollectionLoader
          onAdd={onAdd}
          filters={filters}
          setFilters={onChangeFilters}
          onActionHandler={onActionHandler}
          organizationId={currentEnvironment.organizationId}
          onClearFilters={onClearFilters}
        />
      </TableListContainer>
    </PageLayout.Content>
  );
};

export default PageContent;
