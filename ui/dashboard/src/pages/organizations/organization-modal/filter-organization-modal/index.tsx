import { useScreen } from 'hooks';
import { OrganizationFilters } from 'pages/organizations/types';
import FilterOrganizationPopup from './popup-filter';
import FilterOrganizationSlideModal from './side-filter';

export type FilterProps = {
  onSubmit: (v: Partial<OrganizationFilters>) => void;
  isOpen: boolean;
  onClose: () => void;
  onClearFilters: () => void;
  filters?: Partial<OrganizationFilters>;
};

const FilterOrganizationModal = ({
  onSubmit,
  isOpen,
  onClose,
  onClearFilters,
  filters
}: FilterProps) => {
  const { fromMobileScreen } = useScreen();

  return fromMobileScreen ? (
    <FilterOrganizationPopup
      onSubmit={onSubmit}
      isOpen={isOpen}
      onClose={onClose}
      onClearFilters={onClearFilters}
      filters={filters}
    />
  ) : (
    <FilterOrganizationSlideModal
      onSubmit={onSubmit}
      onClose={onClose}
      filters={filters}
      onClearFilters={onClearFilters}
      isOpen={isOpen}
    />
  );
};

export default FilterOrganizationModal;
