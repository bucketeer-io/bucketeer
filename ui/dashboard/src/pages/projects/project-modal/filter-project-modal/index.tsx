import { ProjectFilters } from 'pages/projects/types';
import FilterModal from 'elements/filter-modal';
import { projectFilterConfig } from './config';

export type FilterProps = {
  onSubmit: (v: Partial<ProjectFilters>) => void;
  isOpen: boolean;
  onClose: () => void;
  onClearFilters: () => void;
  filters?: Partial<ProjectFilters>;
};

const FilterProjectModal = (props: FilterProps) => (
  <FilterModal config={projectFilterConfig} {...props} />
);

export default FilterProjectModal;
