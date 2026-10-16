import { NotificationFilters } from 'pages/notifications/types';
import FilterModal from 'elements/filter-modal';
import { notificationFilterConfig } from './config';

export type FilterProps = {
  isOpen: boolean;
  filters?: Partial<NotificationFilters>;
  onSubmit: (v: Partial<NotificationFilters>) => void;
  onClose: () => void;
  onClearFilters: () => void;
};

const FilterNotificationModal = (props: FilterProps) => (
  <FilterModal config={notificationFilterConfig} {...props} />
);

export default FilterNotificationModal;
