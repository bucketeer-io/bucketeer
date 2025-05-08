import { Feature } from '@types';
import ArchiveFlag from './archive-flag';
import CloneFlag from './clone-flag';
import GeneralInfoForm from './general-info-form';

const SettingsPage = ({ feature }: { feature: Feature }) => {
  return (
    <div className="flex flex-col w-full p-5 pt-0 gap-y-6">
      <GeneralInfoForm feature={feature} />
      <CloneFlag feature={feature} />
      <ArchiveFlag feature={feature} />
    </div>
  );
};

export default SettingsPage;
