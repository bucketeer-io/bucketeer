import { useScreen } from 'hooks';
import { UserSegmentsFilters } from 'pages/user-segments/types';
import FilterUserSegmentPopup from './popup-filter';
import FilterUserSegmentSlideModal from './slide-filter';

export type FilterProps = {
  onSubmit: (v: Partial<UserSegmentsFilters>) => void;
  isOpen: boolean;
  onClose: () => void;
  onClearFilters: () => void;
  filters?: Partial<UserSegmentsFilters>;
};

const FilterUserSegmentModal = ({
  onSubmit,
  isOpen,
  onClose,
  onClearFilters,
  filters
}: FilterProps) => {
  const { fromMobileScreen } = useScreen();

  return fromMobileScreen ? (
    <FilterUserSegmentPopup
      isOpen={isOpen}
      filters={filters}
      onClose={onClose}
      onSubmit={onSubmit}
      onClearFilters={onClearFilters}
    />
  ) : (
    <FilterUserSegmentSlideModal
      isOpen={isOpen}
      filters={filters}
      onClose={onClose}
      onSubmit={onSubmit}
      onClearFilters={onClearFilters}
    />
  );
};

export default FilterUserSegmentModal;
