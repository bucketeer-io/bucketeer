import { ReactElement } from 'react';
import { UserSegment } from '@types';
import { UserSegmentCard } from 'components/mobile-card/user-segment-card';
import PageLayout from 'elements/page-layout';
import { UserSegmentsActionsType } from '../types';

interface CardCollectionProps {
  data: UserSegment[];
  getUploadingStatus: (segmet: UserSegment) => boolean | undefined;
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
        ? data.map(project => (
            <UserSegmentCard
              key={project.id}
              getUploadingStatus={getUploadingStatus}
              onActionHandler={onActions}
              data={project}
            />
          ))
        : emptyCollection}
    </div>
  );
};
