import test from 'ava';
import { Clause } from '../../proto/feature/clause_pb';
import { Evaluation } from '../../proto/feature/evaluation_pb';
import { Feature } from '../../proto/feature/feature_pb';
import { Prerequisite } from '../../proto/feature/prerequisite_pb';
import { Reason } from '../../proto/feature/reason_pb';
import { Rule } from '../../proto/feature/rule_pb';
import { SegmentUser } from '../../proto/feature/segment_pb';
import { FixedStrategy, Strategy } from '../../proto/feature/strategy_pb';
import { Target } from '../../proto/feature/target_pb';
import { Variation } from '../../proto/feature/variation_pb';

test('exampleTest', (t) => {
  t.pass();
})

export function creatFeature(
  id: string,
  name: string,
  version: number,
  enabled: boolean,
  createdAt: number,
  variationType: Feature.VariationTypeMap[keyof Feature.VariationTypeMap],
  variations: Array<{ id: string; value: string; name: string; description: string }>,
  targets: Array<{ variation: string; users: string[] }>,
  rules: Array<{
    id: string;
    attribute: string;
    operator: Clause.OperatorMap[keyof Clause.OperatorMap];
    values: string[];
    fixedVariation: string;
  }>,
  defaultStrategy: { type: Strategy.TypeMap[keyof Strategy.TypeMap]; variation: string }
): Feature {
  const feature = new Feature();
  feature.setId(id);
  feature.setName(name);
  feature.setVersion(version);
  feature.setEnabled(enabled);
  feature.setCreatedAt(createdAt);
  feature.setVariationType(variationType);

  // Set variations
  const variationList = variations.map(v => createVariation(v.id, v.value, v.name, v.description));
  feature.setVariationsList(variationList);

  // Set targets
  const targetList = targets.map(t => createTarget(t.variation, t.users));
  feature.setTargetsList(targetList);

  // Set rules
  const ruleList = rules.map(r => createRule(r.id, r.attribute, r.operator, r.values));
  feature.setRulesList(ruleList);

  // Set default strategy
  const defaultFixedStrategy = createFixedStrategy(defaultStrategy.variation);
  const strategy = createStrategy(defaultStrategy.type, defaultFixedStrategy);
  feature.setDefaultStrategy(strategy);

  return feature;
}

export function createTarget(variation: string, users: string[]): Target {
  const target = new Target();
  target.setVariation(variation);
  target.setUsersList(users);
  return target;
}

export function createFixedStrategy(
  variation: string,
): FixedStrategy {
  const fixedStrategy = new FixedStrategy();
  fixedStrategy.setVariation(variation);
  return fixedStrategy;
}

export function createStrategy(
  type: Strategy.TypeMap[keyof Strategy.TypeMap],
  fixedStrategy: FixedStrategy,
): Strategy {
  const strategy = new Strategy();
  strategy.setType(type);
  strategy.setFixedStrategy(fixedStrategy);
  return strategy;
}

export function createRule(
  id: string,
  attribute: string,
  operator: Clause.OperatorMap[keyof Clause.OperatorMap],
  values: string[],
): Rule {
  const rule = new Rule();
  rule.setId(id);

  const fixedStrategy = createFixedStrategy('variation-A');
  const strategy = createStrategy(Strategy.Type.FIXED, fixedStrategy);
  rule.setStrategy(strategy);

  rule.setClausesList([createClause(id, attribute, operator, values)]);
  return rule;
}

export function createClause(
  id: string,
  attribute: string,
  operator: Clause.OperatorMap[keyof Clause.OperatorMap],
  values: string[],
): Clause {
  const clause = new Clause();
  clause.setId(id);
  clause.setAttribute(attribute);
  clause.setOperator(operator);
  clause.setValuesList(values);
  return clause;
}

export function createSegmentUser(
  userId: string,
  segmentId: string,
  state: SegmentUser.StateMap[keyof SegmentUser.StateMap],
): SegmentUser {
  const segmentUser = new SegmentUser();
  segmentUser.setUserId(userId);
  segmentUser.setSegmentId(segmentId);
  segmentUser.setState(state);
  return segmentUser;
}

export function createPrerequisite(featureId: string, variationId: string): Prerequisite {
  const prerequisite = new Prerequisite();
  prerequisite.setFeatureId(featureId);
  prerequisite.setVariationId(variationId);
  return prerequisite;
}

// Function to create a Variation instance with the provided parameters
export function createVariation(
  id: string,
  value: string,
  name: string,
  description: string
): Variation {
  const variation = new Variation();
  variation.setId(id);
  variation.setValue(value);
  variation.setName(name);
  variation.setDescription(description);

  return variation;
}

export function createEvaluation(
  id: string,
  featureId: string,
  featureVersion: number,
  userId: string,
  variationId: string,
  variationValue: string,
  variationName: string,
  reason: Reason,
): Evaluation {
  const evaluation = new Evaluation();
  evaluation.setId(id);
  evaluation.setFeatureId(featureId);
  evaluation.setFeatureVersion(featureVersion);
  evaluation.setUserId(userId);
  evaluation.setVariationId(variationId);
  evaluation.setVariationValue(variationValue);
  evaluation.setVariationName(variationName);
  evaluation.setReason(reason);
  return evaluation;
}

export function createReason(
  ruleId: string,
  type: Reason.TypeMap[keyof Reason.TypeMap],
): Reason {
  const reason = new Reason();
  reason.setType(type);
  reason.setRuleId(ruleId);
  return reason;
}