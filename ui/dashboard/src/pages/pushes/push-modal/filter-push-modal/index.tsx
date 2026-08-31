import { PushFilters } from 'pages/pushes/types';
import FilterModal from 'elements/filter-modal';
import { pushFilterConfig } from './config';

export type FilterProps = {
  isOpen: boolean;
  filters?: Partial<PushFilters>;
  onSubmit: (v: Partial<PushFilters>) => void;
  onClose: () => void;
  onClearFilters: () => void;
};

const FilterPushModal = (props: FilterProps) => (
  <FilterModal config={pushFilterConfig} {...props} />
);

export default FilterPushModal;
