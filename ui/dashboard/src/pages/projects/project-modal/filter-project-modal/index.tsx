import { useScreen } from 'hooks';
import { ProjectFilters } from 'pages/projects/types';
import FilterProjectPopup from './popup-filter';
import FilterProjectSlide from './slide-filter';

export type FilterProps = {
  onSubmit: (v: Partial<ProjectFilters>) => void;
  isOpen: boolean;
  onClose: () => void;
  onClearFilters: () => void;
  filters?: Partial<ProjectFilters>;
};

const FilterProjectModal = ({
  onSubmit,
  isOpen,
  onClose,
  onClearFilters,
  filters
}: FilterProps) => {
  const { fromMobileScreen } = useScreen();
  return fromMobileScreen ? (
    <FilterProjectPopup
      isOpen={isOpen}
      onSubmit={onSubmit}
      onClearFilters={onClearFilters}
      onClose={onClose}
      filters={filters}
    />
  ) : (
    <FilterProjectSlide
      isOpen={isOpen}
      onSubmit={onSubmit}
      onClearFilters={onClearFilters}
      onClose={onClose}
      filters={filters}
    />
  );
};

export default FilterProjectModal;
