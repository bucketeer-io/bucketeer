import { useState } from 'react';
import { useToggleOpen } from 'hooks/use-toggle-open';
import PageLayout from 'elements/page-layout';
import { EmptyCollection } from './collection-layout/empty-collection';
import PageContent from './page-content';
import AddUserSegmentModal from './user-segment-modal/add-segment-modal';
import EditUserSegmentModal from './user-segment-modal/edit-segment-modal';

export type UserSegments = {
  id: string;
  name: string;
  description: string;
  rules: string;
  createdAt: string;
  updatedAt: string;
  version: string;
  deleted: boolean;
  includedUserCount: number;
  excludedUserCount: number;
  status: string;
  isInUseStatus: boolean;
  features: unknown;
  connections: number;
};

export const mocks = [
  {
    id: 'segment_1',
    name: 'User Segment 1',
    description: 'User Segment 1 description',
    rules: '',
    createdAt: '1/7/2025, 11:19:54 AM',
    updatedAt: '1/7/2025, 11:19:54 AM',
    version: '',
    deleted: false,
    includedUserCount: 10,
    excludedUserCount: 10,
    status: 'not-in-use',
    isInUseStatus: false,
    features: null,
    connections: 1
  },
  {
    id: 'segment_2',
    name: 'User Segment 2',
    description: 'User Segment 2 description',
    rules: '',
    createdAt: '1/6/2025, 11:19:54 AM',
    updatedAt: '1/6/2025, 11:19:54 AM',
    version: '',
    deleted: false,
    includedUserCount: 20,
    excludedUserCount: 10,
    status: 'new',
    isInUseStatus: false,
    features: null,
    connections: 0
  },
  {
    id: 'segment_3',
    name: 'User Segment 3',
    description: 'User Segment 3 description',
    rules: '',
    createdAt: '1/5/2025, 11:19:54 AM',
    updatedAt: '1/5/2025, 11:19:54 AM',
    version: '',
    deleted: false,
    includedUserCount: 7,
    excludedUserCount: 10,
    status: 'in-use',
    isInUseStatus: true,
    features: null,
    connections: 5
  }
];

export const collection = {
  userSegments: mocks,
  totalCount: 3
};

const PageLoader = () => {
  const isLoading = false;
  const isError = false;

  const [selectedSegment, setSelectedSegment] = useState<UserSegments>();

  const [isOpenAddModal, onOpenAddModal, onCloseAddModal] =
    useToggleOpen(false);
  const [isOpenEditModal, onOpenEditModal, onCloseEditModal] =
    useToggleOpen(false);

  const isEmpty = collection?.userSegments.length === 0;

  return (
    <>
      {isLoading ? (
        <PageLayout.LoadingState />
      ) : isError ? (
        <PageLayout.ErrorState onRetry={() => {}} />
      ) : isEmpty ? (
        <PageLayout.EmptyState>
          <EmptyCollection onAdd={onOpenAddModal} />
        </PageLayout.EmptyState>
      ) : (
        <PageContent
          onAdd={onOpenAddModal}
          onEdit={value => {
            setSelectedSegment(value);
            onOpenEditModal();
          }}
        />
      )}
      {isOpenAddModal && (
        <AddUserSegmentModal
          isOpen={isOpenAddModal}
          onClose={onCloseAddModal}
        />
      )}
      {isOpenEditModal && (
        <EditUserSegmentModal
          isOpen={isOpenEditModal}
          onClose={onCloseEditModal}
          userSegment={selectedSegment!}
        />
      )}
    </>
  );
};

export default PageLoader;
