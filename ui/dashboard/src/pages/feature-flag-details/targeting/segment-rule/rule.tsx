import { useMemo } from 'react';
import useOptions from 'hooks/use-options';
import compact from 'lodash/compact';
import { Feature } from '@types';
import RuleClausesForm, { SituationOption } from 'elements/rule-clauses-form';

interface Props {
  feature: Feature;
  features: Feature[];
  segmentIndex: number;
  sdkAttributeKeys: string[];
}

const RuleForm = ({
  feature,
  features,
  segmentIndex,
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
      sdkAttributeKeys={sdkAttributeKeys}
    />
  );
};

export default RuleForm;
