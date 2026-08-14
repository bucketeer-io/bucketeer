import { UserSegment } from '@types';
import PageLayout from 'elements/page-layout';
import SegmentForm from './segment-form';

const PageContent = ({
  isUpdate,
  userSegment,
  isDisabled
}: {
  isUpdate: boolean;
  userSegment?: UserSegment;
  isDisabled: boolean;
}) => {
  return (
    <PageLayout.Content className="p-4 sm:p-6">
      <SegmentForm
        isUpdate={isUpdate}
        userSegment={userSegment}
        isDisabled={isDisabled}
      />
    </PageLayout.Content>
  );
};

export default PageContent;
