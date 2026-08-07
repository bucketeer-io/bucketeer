import { useMemo } from 'react';
import useOptions from 'hooks/use-options';
import compact from 'lodash/compact';
import { Feature, UserSegment } from '@types';
import RuleClausesForm, { SituationOption } from 'elements/rule-clauses-form';

interface Props {
  feature: Feature;
  features: Feature[];
  segmentIndex: number;
  userSegments?: UserSegment[];
  sdkAttributeKeys: string[];
}

const RuleForm = ({
  feature,
  features,
  segmentIndex,
  userSegments,
  sdkAttributeKeys
}: Props) => {
  const { situationOptions } = useOptions();

  const usedAttributeKeys: string[] = useMemo(
    () =>
      compact(
        feature.rules
          ?.flatMap(item => item?.clauses || [])
          .map(clause => clause.attribute)
      ),
    [feature.rules]
  );

  return (
    <RuleClausesForm
      rulesFieldName="segmentRules"
      ruleIndex={segmentIndex}
      situationOptions={situationOptions as SituationOption[]}
      usedAttributeKeys={usedAttributeKeys}
      feature={feature}
      features={features}
      userSegments={userSegments}
      sdkAttributeKeys={sdkAttributeKeys}
    />
  );
};

export default RuleForm;
