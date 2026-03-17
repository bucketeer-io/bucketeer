import { useScreen } from 'hooks';
import { FlagFilters } from 'pages/feature-flags/types';
import FilterFlagPopup from './popup-filter';
import FilterFlagSlide from './slide-filter';

export type FilterProps = {
  isOpen: boolean;
  filters?: Partial<FlagFilters>;
  onSubmit: (v: Partial<FlagFilters>) => void;
  onClose: () => void;
  onClearFilters: () => void;
};

const FilterFlagModal = ({
  isOpen,
  filters,
  onSubmit,
  onClose,
  onClearFilters
}: FilterProps) => {
  const { fromMobileScreen } = useScreen();

  return fromMobileScreen ? (
    <FilterFlagPopup
      isOpen={isOpen}
      onClearFilters={onClearFilters}
      onClose={onClose}
      onSubmit={onSubmit}
      filters={filters}
    />
  ) : (
    <FilterFlagSlide
      isOpen={isOpen}
      onClearFilters={onClearFilters}
      onClose={onClose}
      onSubmit={onSubmit}
      filters={filters}
    />
  );
};

export default FilterFlagModal;
