import { useState } from 'react';
import { useParams } from 'react-router-dom';
import { usePartialState, useToggleOpen } from 'hooks';
import pickBy from 'lodash/pickBy';
import { Project } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import CollectionLoader from 'pages/projects/collection-loader';
import EditProjectModal from 'pages/projects/project-modal/edit-project-modal';
import FilterProjectModal from 'pages/projects/project-modal/filter-project-modal';
import { ProjectFilters } from 'pages/projects/types';
import Filter from 'elements/filter';

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

  const onChangeFilters = (values: Partial<ProjectFilters>) => {
    const options = pickBy({ ...filters, ...values }, v => isNotEmpty(v));
    onChangSearchParams(options);
    setFilters({ ...values });
  };

  return (
    <>
      <Filter
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
        <EditProjectModal
          isOpen={isOpenEditModal}
          onClose={onCloseEditModal}
          project={selectedProject!}
        />
      )}
      <CollectionLoader
        filters={filters}
        organizationIds={[organizationId!]}
        setFilters={onChangeFilters}
        onActionHandler={value => {
          setSelectedProject(value);
          onOpenEditModal();
        }}
      />
    </>
  );
};

export default OrganizationProjects;
