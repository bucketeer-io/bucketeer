import { useEffect } from 'react';
import { IconAddOutlined } from 'react-icons-material-design';
import { getCurrentEnvironment, hasEditable, useAuth } from 'auth';
import { usePartialState, useToggleOpen } from 'hooks';
import { useTranslation } from 'i18n';
import pickBy from 'lodash/pickBy';
import { Project } from '@types';
import { isEmptyObject, isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import Button from 'components/button';
import Icon from 'components/icon';
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
  const editable = hasEditable(consoleAccount!);

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

  const onChangeFilters = (values: Partial<ProjectFilters>) => {
    const options = pickBy({ ...filters, ...values }, v => isNotEmpty(v));
    onChangSearchParams(options);
    setFilters({ ...values });
  };

  const onActionHandler = (project: Project) => {
    onEdit(project);
  };

  useEffect(() => {
    if (isEmptyObject(searchOptions)) {
      setFilters({ ...defaultFilters });
    }
  }, [searchOptions]);

  return (
    <PageLayout.Content>
      <Filter
        onOpenFilter={onOpenFilterModal}
        action={
          <Button
            className="flex-1 lg:flex-none"
            onClick={onAdd}
            disabled={!editable}
          >
            <Icon icon={IconAddOutlined} size="sm" />
            {t(`new-project`)}
          </Button>
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
          onClearFilters={() =>
            setFilters({ searchQuery: '', disabled: undefined })
          }
        />
      </TableListContainer>
    </PageLayout.Content>
  );
};

export default PageContent;
