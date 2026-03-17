import { useScreen } from 'hooks';
import { MembersFilters } from 'pages/members/types';
import FilterMemberPopup from './popup-filter';
import FilterMemberSlide from './slide-filter';

export type FilterProps = {
  onSubmit: (v: Partial<MembersFilters>) => void;
  isOpen: boolean;
  onClose: () => void;
  onClearFilters: () => void;
  filters?: Partial<MembersFilters>;
};

const FilterMemberModal = ({
  isOpen,
  filters,
  onSubmit,
  onClose,
  onClearFilters
}: FilterProps) => {
  const { fromMobileScreen } = useScreen();
  return fromMobileScreen ? (
    <FilterMemberPopup
      isOpen={isOpen}
      onClearFilters={onClearFilters}
      onSubmit={onSubmit}
      onClose={onClose}
      filters={filters}
    />
  ) : (
    <FilterMemberSlide
      isOpen={isOpen}
      onClearFilters={onClearFilters}
      onSubmit={onSubmit}
      onClose={onClose}
      filters={filters}
    />
  );
};

export default FilterMemberModal;
