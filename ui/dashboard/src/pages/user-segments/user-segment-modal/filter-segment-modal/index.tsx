import { UserSegmentsFilters } from 'pages/user-segments/types';
import FilterModal from 'elements/filter-modal';
import { userSegmentFilterConfig } from './config';

export type FilterProps = {
  onSubmit: (v: Partial<UserSegmentsFilters>) => void;
  isOpen: boolean;
  onClose: () => void;
  onClearFilters: () => void;
  filters?: Partial<UserSegmentsFilters>;
};

const FilterUserSegmentModal = (props: FilterProps) => (
  <FilterModal config={userSegmentFilterConfig} {...props} />
);

export default FilterUserSegmentModal;
