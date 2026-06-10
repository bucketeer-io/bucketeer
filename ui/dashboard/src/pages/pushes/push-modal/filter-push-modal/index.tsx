import { useScreen } from 'hooks';
import { PushFilters } from 'pages/pushes/types';
import FilterPushPopup from './popup-filter';
import FilterPushSlide from './side-filter';

export type FilterProps = {
  onSubmit: (v: Partial<PushFilters>) => void;
  isOpen: boolean;
  onClose: () => void;
  onClearFilters: () => void;
  filters?: Partial<PushFilters>;
};

const FilterPushModal = ({
  onSubmit,
  isOpen,
  onClose,
  onClearFilters,
  filters
}: FilterProps) => {
  const { isMobile } = useScreen();
  return !isMobile ? (
    <FilterPushPopup
      isOpen={isOpen}
      onClearFilters={onClearFilters}
      onClose={onClose}
      onSubmit={onSubmit}
      filters={filters}
    />
  ) : (
    <FilterPushSlide
      isOpen={isOpen}
      onClearFilters={onClearFilters}
      onClose={onClose}
      onSubmit={onSubmit}
      filters={filters}
    />
  );
};

export default FilterPushModal;
