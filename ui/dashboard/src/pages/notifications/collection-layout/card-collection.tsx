import { ReactElement } from 'react';
import { Notification } from '@types';
import { NotificationCard } from 'components/mobile-card/notification-card';
import PageLayout from 'elements/page-layout';
import { NotificationActionsType } from '../types';

interface CardCollectionProps {
  isLoading: boolean;
  emptyCollection?: ReactElement;
  data: Notification[];
  onActions: (item: Notification, type: NotificationActionsType) => void;
}

export const CardCollection = ({
  isLoading,
  emptyCollection,
  data,
  onActions
}: CardCollectionProps) => {
  return isLoading ? (
    <PageLayout.LoadingState className="py-10" />
  ) : (
    <div className="flex flex-col gap-3">
      {data.length
        ? data.map(notification => (
            <NotificationCard
              key={notification.id}
              onActions={onActions}
              data={notification}
            />
          ))
        : emptyCollection}
    </div>
  );
};
