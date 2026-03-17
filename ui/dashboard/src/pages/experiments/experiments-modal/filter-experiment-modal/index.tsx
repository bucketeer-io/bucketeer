import { useScreen } from 'hooks';
import { ExperimentFilters } from 'pages/experiments/types';
import FilterExperimentPopup from './popup-filter';
import FilterExperimentSlideModal from './slide-filter';

export type FilterProps = {
  onSubmit: (v: Partial<ExperimentFilters>) => void;
  isOpen: boolean;
  onClose: () => void;
  onClearFilters: () => void;
  filters?: Partial<ExperimentFilters>;
};

const FilterExperimentModal = ({
  onSubmit,
  isOpen,
  onClose,
  onClearFilters,
  filters
}: FilterProps) => {
  const { fromMobileScreen } = useScreen();

  return fromMobileScreen ? (
    <FilterExperimentPopup
      onSubmit={onSubmit}
      isOpen={isOpen}
      onClose={onClose}
      onClearFilters={onClearFilters}
      filters={filters}
    />
  ) : (
    <FilterExperimentSlideModal
      isOpen={isOpen}
      onSubmit={onSubmit}
      filters={filters}
      onClose={onClose}
      onClearFilters={onClearFilters}
    />
  );
};

export default FilterExperimentModal;
