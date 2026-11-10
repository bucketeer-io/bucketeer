import { FlagFilters } from 'pages/feature-flags/types';
import FilterModal from 'elements/filter-modal';
import { flagFilterConfig } from './config';

export type FilterProps = {
  isOpen: boolean;
  filters?: Partial<FlagFilters>;
  onSubmit: (v: Partial<FlagFilters>) => void;
  onClose: () => void;
  onClearFilters: () => void;
};

const FilterFlagModal = (props: FilterProps) => (
  <FilterModal config={flagFilterConfig} {...props} />
);

export default FilterFlagModal;
