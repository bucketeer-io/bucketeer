import { useScreen } from 'hooks';
import { APIKeysFilters } from 'pages/api-keys/types';
import FilterProjectSlide from 'pages/projects/project-modal/filter-project-modal/slide-filter';
import FilterAPIKeyPopup from './popup-filter';

export type FilterProps = {
  onSubmit: (v: Partial<APIKeysFilters>) => void;
  isOpen: boolean;
  onClose: () => void;
  onClearFilters: () => void;
  filters?: Partial<APIKeysFilters>;
};

export interface Option {
  value: string;
  label: string;
}

const FilterAPIKeyModal = ({
  onSubmit,
  isOpen,
  onClose,
  onClearFilters,
  filters
}: FilterProps) => {
  const { fromMobileScreen } = useScreen();
  return fromMobileScreen ? (
    <FilterAPIKeyPopup
      isOpen={isOpen}
      onClearFilters={onClearFilters}
      onClose={onClose}
      onSubmit={onSubmit}
      filters={filters}
    />
  ) : (
    <FilterProjectSlide
      isOpen={isOpen}
      onClearFilters={onClearFilters}
      onClose={onClose}
      onSubmit={onSubmit}
      filters={filters}
    />
  );
};

export default FilterAPIKeyModal;
