import { getCurrentEnvironment, useAuth } from 'auth';
import { SCHEDULED_FLAG_CHANGES_ENABLED } from 'configs';
import { Feature } from '@types';
import ScheduledChangesBanner from '../elements/scheduled-changes-banner';
import ArchiveFlag from './archive-flag';
import CloneFlag from './clone-flag';
import GeneralInfoForm from './general-info-form';

const SettingsPage = ({
  feature,
  editable
}: {
  feature: Feature;
  editable: boolean;
}) => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  return (
    <div className="flex flex-col w-full p-5 pt-0 gap-y-6">
      {SCHEDULED_FLAG_CHANGES_ENABLED && (
        <ScheduledChangesBanner
          featureId={feature.id}
          environmentId={currentEnvironment.id}
        />
      )}
      <GeneralInfoForm feature={feature} disabled={!editable} />
      <CloneFlag feature={feature} disabled={!editable} />
      <ArchiveFlag feature={feature} disabled={!editable} />
    </div>
  );
};

export default SettingsPage;
