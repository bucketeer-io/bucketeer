import { useCallback, useState } from 'react';
import AddRuleDropdown from './add-rule-dropdown';
import {
  initialIndividualRule,
  initialPrerequisitesRule,
  initialSegmentCondition
} from './constants';
import DefaultRule from './default-rule';
import TargetSegmentRule from './target-segment-rule';
import TargetingState from './targeting-state';
import {
  RuleCategory,
  TargetIndividualItem,
  TargetPrerequisiteItem,
  TargetSegmentItem
} from './types';

const Targeting = () => {
  const [targetSegmentRules, setTargetSegmentRules] = useState<
    TargetSegmentItem[]
  >([]);
  const [targetIndividualRules, setTargetIndividualRules] = useState<
    TargetIndividualItem[]
  >([]);
  const [prerequisitesRules, setPrerequisitesRules] = useState<
    TargetPrerequisiteItem[]
  >([]);

  const onAddRule = useCallback(
    (type: RuleCategory) => {
      if (type === 'target-segments') {
        return setTargetSegmentRules([
          ...targetSegmentRules,
          {
            index: targetSegmentRules.length + 1,
            rules: [
              {
                variation: true,
                conditions: [initialSegmentCondition]
              }
            ]
          }
        ]);
      }
      if (type === 'target-individuals') {
        return setTargetIndividualRules([
          ...targetIndividualRules,
          {
            index: targetIndividualRules.length + 1,
            rules: [initialIndividualRule]
          }
        ]);
      }
      setPrerequisitesRules([
        ...prerequisitesRules,
        {
          index: prerequisitesRules.length + 1,
          rules: [initialPrerequisitesRule]
        }
      ]);
    },
    [targetSegmentRules, targetIndividualRules, prerequisitesRules]
  );

  return (
    <div className="flex flex-col size-full gap-y-6 overflow-visible">
      <TargetingState />
      <AddRuleDropdown onAddRule={onAddRule} />
      <TargetSegmentRule
        targetSegmentRules={targetSegmentRules}
        setTargetSegmentRules={setTargetSegmentRules}
      />
      <DefaultRule />
    </div>
  );
};

export default Targeting;
