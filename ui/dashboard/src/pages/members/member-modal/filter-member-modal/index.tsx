import { MembersFilters } from 'pages/members/types';
import FilterModal from 'elements/filter-modal';
import { memberFilterConfig } from './config';

export type FilterProps = {
  onSubmit: (v: Partial<MembersFilters>) => void;
  isOpen: boolean;
  onClose: () => void;
  onClearFilters: () => void;
  filters?: Partial<MembersFilters>;
};

const FilterMemberModal = (props: FilterProps) => (
  <FilterModal config={memberFilterConfig} {...props} />
);

export default FilterMemberModal;
