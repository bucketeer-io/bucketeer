import { Feature } from '@types';
import PageLayout from 'elements/page-layout';
import TriggerSection from './trigger-section';

const TriggerPage = ({
  feature,
  editable
}: {
  feature: Feature;
  editable: boolean;
}) => {
  return (
    <PageLayout.Content className="p-6 pt-0 gap-y-6 min-w-[900px]">
      <TriggerSection feature={feature} editable={editable} />
    </PageLayout.Content>
  );
};

export default TriggerPage;
