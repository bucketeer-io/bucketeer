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
    <PageLayout.Content className="p-6 min-w-[900px]">
      <SegmentForm
        isUpdate={isUpdate}
        userSegment={userSegment}
        isDisabled={isDisabled}
      />
    </PageLayout.Content>
  );
};

export default PageContent;
