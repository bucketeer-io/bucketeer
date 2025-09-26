import { useCallback, useMemo } from 'react';
import { useParams } from 'react-router-dom';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_ORGANIZATIONS, PAGE_PATH_PROJECTS } from 'constants/routing';
import { usePartialState, useToggleOpen } from 'hooks';
import useActionWithURL from 'hooks/use-action-with-url';
import { useTranslation } from 'i18n';
import pickBy from 'lodash/pickBy';
import { Project } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import CollectionLoader from 'pages/projects/collection-loader';
import FilterProjectModal from 'pages/projects/project-modal/filter-project-modal';
import ProjectCreateUpdateModal from 'pages/projects/project-modal/project-create-update-modal/index.tsx';
import { ProjectFilters } from 'pages/projects/types';
import Filter from 'elements/filter';
import TableListContainer from 'elements/table-list-container';

const OrganizationProjects = () => {
  const { organizationId } = useParams();
  const { t } = useTranslation(['form']);
  const { searchOptions, onChangSearchParams } = useSearchParams();
  const searchFilters: Partial<ProjectFilters> = searchOptions;

  const [openFilterModal, onOpenFilterModal, onCloseFilterModal] =
    useToggleOpen(false);

  const defaultFilters = {
    page: 1,
    orderBy: 'CREATED_AT',
    orderDirection: 'DESC',
    ...searchFilters
  } as ProjectFilters;

  const [filters, setFilters] = usePartialState<ProjectFilters>(defaultFilters);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const commonPath = useMemo(
    () => `${PAGE_PATH_ORGANIZATIONS}/${organizationId}${PAGE_PATH_PROJECTS}`,
    [currentEnvironment]
  );
  const { isEdit, onOpenEditModal, onCloseActionModal } = useActionWithURL({
    closeModalPath: commonPath,
    idKey: 'projectId'
  });
  const onChangeFilters = useCallback(
    (values: Partial<ProjectFilters>) => {
      values.page = values?.page || 1;
      const options = pickBy({ ...filters, ...values }, v => isNotEmpty(v));
      onChangSearchParams(options);
      setFilters({ ...values });
    },
    [filters]
  );

  const onActionHandler = useCallback((project: Project) => {
    onOpenEditModal(
      `${PAGE_PATH_ORGANIZATIONS}/${organizationId}${PAGE_PATH_PROJECTS}/${project.id}`
    );
  }, []);

  const onClearFilters = useCallback(() => {
    setFilters({
      searchQuery: '',
      disabled: undefined
    });
  }, []);

  return (
    <div className="flex flex-col flex-1 size-full">
      <Filter
        isShowDocumentation={false}
        placeholder={t('name-email-search-placeholder')}
        name="org-projects-list-search"
        onOpenFilter={onOpenFilterModal}
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
      {!!isEdit && (
        <ProjectCreateUpdateModal
          isOpen={!!isEdit}
          onClose={onCloseActionModal}
        />
      )}
      <TableListContainer className="self-stretch">
        <CollectionLoader
          filters={filters}
          organizationId={organizationId}
          setFilters={onChangeFilters}
          onActionHandler={onActionHandler}
          onClearFilters={onClearFilters}
        />
      </TableListContainer>
    </div>
  );
};

export default OrganizationProjects;
