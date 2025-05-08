import DefaultVariations from './default-variations';
import FlagType from './flag-type';
import Variations from './variations';
import VariationsSwitch from './variations-switch';

const FlagVariations = () => {
  return (
    <div className="flex flex-col w-full p-5 gap-y-6 bg-white rounded-lg shadow-card">
      <VariationsSwitch />
      <FlagType />
      <Variations />
      <DefaultVariations />
    </div>
  );
};

export default FlagVariations;
