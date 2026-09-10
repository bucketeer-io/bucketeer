import { ExperimentFilters } from 'pages/experiments/types';
import FilterModal from 'elements/filter-modal';
import { experimentFilterConfig } from './config';

export type FilterProps = {
  onSubmit: (v: Partial<ExperimentFilters>) => void;
  isOpen: boolean;
  onClose: () => void;
  onClearFilters: () => void;
  filters?: Partial<ExperimentFilters>;
};

const FilterExperimentModal = (props: FilterProps) => (
  <FilterModal config={experimentFilterConfig} {...props} />
);

export default FilterExperimentModal;
