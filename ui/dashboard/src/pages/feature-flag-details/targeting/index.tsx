import AddRuleDropdown from './add-rule-dropdown';
import DefaultRule from './default-rule';
import TargetSegmentRule from './target-segment-rule';
import TargetingState from './targeting-state';

const Targeting = () => {
  return (
    <div className="flex flex-col size-full gap-y-6 overflow-visible">
      <TargetingState />
      <AddRuleDropdown />
      <TargetSegmentRule />
      <DefaultRule />
    </div>
  );
};

export default Targeting;
