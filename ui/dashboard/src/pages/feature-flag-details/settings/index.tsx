import { Feature } from '@types';
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
  return (
    <div className="flex flex-col w-full p-5 pt-0 gap-y-6">
      <GeneralInfoForm feature={feature} disabled={!editable} />
      <CloneFlag feature={feature} disabled={!editable} />
      <ArchiveFlag feature={feature} disabled={!editable} />
    </div>
  );
};

export default SettingsPage;
