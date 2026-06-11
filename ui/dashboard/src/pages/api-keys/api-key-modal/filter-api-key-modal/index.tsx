import { useScreen } from 'hooks';
import { APIKeysFilters } from 'pages/api-keys/types';
import FilterAPIKeyPopup from './popup-filter';
import FilterAPIKeySlide from './slide-filter';

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
  const { isMobile } = useScreen();
  return !isMobile ? (
    <FilterAPIKeyPopup
      isOpen={isOpen}
      onClearFilters={onClearFilters}
      onClose={onClose}
      onSubmit={onSubmit}
      filters={filters}
    />
  ) : (
    <FilterAPIKeySlide
      isOpen={isOpen}
      onClearFilters={onClearFilters}
      onClose={onClose}
      onSubmit={onSubmit}
      filters={filters}
    />
  );
};

export default FilterAPIKeyModal;
