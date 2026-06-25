import { ReactElement } from 'react';
import { UserSegment } from '@types';
import { UserSegmentCard } from 'components/mobile-card/user-segment-card';
import PageLayout from 'elements/page-layout';
import { UserSegmentsActionsType } from '../types';

interface CardCollectionProps {
  data: UserSegment[];
  getUploadingStatus: (segment: UserSegment) => boolean | undefined;
  onActions: (value: UserSegment, type: UserSegmentsActionsType) => void;
  emptyCollection?: ReactElement;
  isLoading?: boolean;
}

export const CardCollection = ({
  data,
  emptyCollection,
  isLoading,
  getUploadingStatus,
  onActions
}: CardCollectionProps) => {
  return isLoading ? (
    <PageLayout.LoadingState className="py-10" />
  ) : (
    <div className="flex flex-col gap-3">
      {data.length
        ? data.map(segment => (
            <UserSegmentCard
              key={segment.id}
              getUploadingStatus={getUploadingStatus}
              onActionHandler={onActions}
              data={segment}
            />
          ))
        : emptyCollection}
    </div>
  );
};
