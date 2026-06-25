import { useScreen } from 'hooks';
import { MembersFilters } from 'pages/members/types';
import FilterNotificationPopup from './popup-filter';
import FilterNotificationSlide from './slide-filter';

export type FilterProps = {
  onSubmit: (v: Partial<MembersFilters>) => void;
  isOpen: boolean;
  onClose: () => void;
  onClearFilters: () => void;
  filters?: Partial<MembersFilters>;
};

const FilterNotificationModal = ({
  isOpen,
  filters,
  onSubmit,
  onClose,
  onClearFilters
}: FilterProps) => {
  const { isMobile } = useScreen();
  return !isMobile ? (
    <FilterNotificationPopup
      isOpen={isOpen}
      onClearFilters={onClearFilters}
      onSubmit={onSubmit}
      onClose={onClose}
      filters={filters}
    />
  ) : (
    <FilterNotificationSlide
      isOpen={isOpen}
      onClearFilters={onClearFilters}
      onSubmit={onSubmit}
      onClose={onClose}
      filters={filters}
    />
  );
};

export default FilterNotificationModal;
