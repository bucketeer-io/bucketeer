import { OrganizationFilters } from 'pages/organizations/types';
import FilterModal from 'elements/filter-modal';
import { organizationFilterConfig } from './config';

export type FilterProps = {
  onSubmit: (v: Partial<OrganizationFilters>) => void;
  isOpen: boolean;
  onClose: () => void;
  onClearFilters: () => void;
  filters?: Partial<OrganizationFilters>;
};

const FilterOrganizationModal = (props: FilterProps) => (
  <FilterModal config={organizationFilterConfig} {...props} />
);

export default FilterOrganizationModal;
