import { APIKeysFilters } from 'pages/api-keys/types';
import FilterModal from 'elements/filter-modal';
import { apiKeyFilterConfig } from './config';

export type FilterProps = {
  isOpen: boolean;
  filters?: Partial<APIKeysFilters>;
  onSubmit: (v: Partial<APIKeysFilters>) => void;
  onClose: () => void;
  onClearFilters: () => void;
};

const FilterAPIKeyModal = (props: FilterProps) => (
  <FilterModal config={apiKeyFilterConfig} {...props} />
);

export default FilterAPIKeyModal;
