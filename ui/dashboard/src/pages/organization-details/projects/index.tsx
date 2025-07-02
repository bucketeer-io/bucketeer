import { useCallback, useState } from 'react';
import { useParams } from 'react-router-dom';
import { usePartialState, useToggleOpen } from 'hooks';
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

  const { searchOptions, onChangSearchParams } = useSearchParams();
  const searchFilters: Partial<ProjectFilters> = searchOptions;
  const [selectedProject, setSelectedProject] = useState<Project>();

  const [openFilterModal, onOpenFilterModal, onCloseFilterModal] =
    useToggleOpen(false);

  const [isOpenEditModal, onOpenEditModal, onCloseEditModal] =
    useToggleOpen(false);

  const defaultFilters = {
    page: 1,
    orderBy: 'CREATED_AT',
    orderDirection: 'DESC',
    ...searchFilters
  } as ProjectFilters;

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
    setSelectedProject(project);
    onOpenEditModal();
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
      {isOpenEditModal && (
        <ProjectCreateUpdateModal
          isOpen={isOpenEditModal}
          onClose={onCloseEditModal}
          project={selectedProject}
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
