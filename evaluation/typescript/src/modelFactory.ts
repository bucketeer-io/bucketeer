import { Clause } from './proto/feature/clause_pb';
import { Evaluation } from './proto/feature/evaluation_pb';
import { Feature } from './proto/feature/feature_pb';
import { Prerequisite } from './proto/feature/prerequisite_pb';
import { Reason } from './proto/feature/reason_pb';
import { Rule } from './proto/feature/rule_pb';
import { SegmentUser } from './proto/feature/segment_pb';
import { FixedStrategy, RolloutStrategy, Strategy } from './proto/feature/strategy_pb';
import { Target } from './proto/feature/target_pb';
import { Variation } from './proto/feature/variation_pb';
import { User } from './proto/user/user_pb';

// Helper function to create a User instance
export function createUser(id: string, data: Record<string, string> | null): User {
  const user = new User();
  user.setId(id);
  let map = user.getDataMap();
  if (data != null) {
    Object.entries(data).forEach(([key, value]) => map.set(key, value));
  }
  return user;
}

export function createFeature(
  options: {
    id?: string;
    name?: string;
    version?: number;
    enabled?: boolean;
    createdAt?: number;
    updatedAt?: number;
    variationType?: Feature.VariationTypeMap[keyof Feature.VariationTypeMap];
    variations?: Array<{ id: string; value: string; name: string; description: string }>;
    targets?: Array<{ variation: string; users: string[] }>;
    rules?: Array<{
      id: string;
      attribute: string;
      operator: Clause.OperatorMap[keyof Clause.OperatorMap];
      values: string[];
      fixedVariation: string;
    }>;
    defaultStrategy?:
      | { type: Strategy.TypeMap[keyof Strategy.TypeMap]; variation: string }
      | undefined;
    prerequisitesList?: Array<Prerequisite>;
    tagList?: Array<string>;
    offVariation?: string;
  } = {},
): Feature {
  const defaultOptions = {
    id: '',
    name: '',
    version: 0,
    enabled: false,
    createdAt: Date.now(),
    updatedAt: Date.now(),
    variationType: Feature.VariationType.STRING,
    variations: [],
    targets: [],
    rules: [],
    defaultStrategy: undefined,
    prerequisitesList: [],
    tagList: [],
    offVariation: undefined,
  };

  const finalOptions = { ...defaultOptions, ...options };

  const feature = new Feature();
  feature.setId(finalOptions.id);
  feature.setName(finalOptions.name);
  feature.setVersion(finalOptions.version);
  feature.setEnabled(finalOptions.enabled);
  feature.setCreatedAt(finalOptions.createdAt);
  feature.setUpdatedAt(finalOptions.updatedAt);
  feature.setVariationType(finalOptions.variationType);

  // Set variations
  const variationList = finalOptions.variations.map((v) =>
    createVariation(v.id, v.value, v.name, v.description),
  );
  feature.setVariationsList(variationList);

  // Set targets
  const targetList = finalOptions.targets.map((t) => createTarget(t.variation, t.users));
  feature.setTargetsList(targetList);

  // Set rules
  const ruleList = finalOptions.rules.map((r) =>
    createRule(r.id, r.attribute, r.operator, r.values, r.fixedVariation),
  );
  feature.setRulesList(ruleList);

  // Set default strategy
  const defaultStrategy = finalOptions.defaultStrategy;
  if (defaultStrategy !== undefined) {
    const defaultFixedStrategy = createFixedStrategy(defaultStrategy.variation);
    const strategy = createStrategy({
      type: defaultStrategy.type,
      fixedStrategy: defaultFixedStrategy,
    });
    feature.setDefaultStrategy(strategy);
  }

  feature.setPrerequisitesList(finalOptions.prerequisitesList);
  feature.setTagsList(finalOptions.tagList);

  if (finalOptions.offVariation !== undefined) {
    feature.setOffVariation(finalOptions.offVariation);
  }

  return feature;
}

export function createTarget(variation: string, users: string[]): Target {
  const target = new Target();
  target.setVariation(variation);
  target.setUsersList(users);
  return target;
}

export function createFixedStrategy(variation: string): FixedStrategy {
  const fixedStrategy = new FixedStrategy();
  fixedStrategy.setVariation(variation);
  return fixedStrategy;
}

export function createRolloutStrategy(options: {
  variations: Array<{ variation: string; weight: number }>;
}): RolloutStrategy {
  const rolloutStrategy = new RolloutStrategy();
  const variations = options.variations.map((v) => {
    const variation = new RolloutStrategy.Variation();
    variation.setVariation(v.variation);
    variation.setWeight(v.weight);
    return variation;
  });
  rolloutStrategy.setVariationsList(variations);
  return rolloutStrategy;
}

export function createStrategy(options: {
  type: Strategy.TypeMap[keyof Strategy.TypeMap];
  fixedStrategy?: FixedStrategy;
  rolloutStrategy?: RolloutStrategy;
}): Strategy {
  const strategy = new Strategy();
  strategy.setType(options.type);
  strategy.setFixedStrategy(options.fixedStrategy);
  strategy.setRolloutStrategy(options.rolloutStrategy);
  return strategy;
}

//TODO: FIX ME - Rule should have a constructor that accepts all the parameters
//Current is missing many parameters
export function createRule(
  id: string,
  attribute: string,
  operator: Clause.OperatorMap[keyof Clause.OperatorMap],
  values: string[],
  fixedVariation: string,
): Rule {
  const rule = new Rule();
  rule.setId(id);

  const fixedStrategy = createFixedStrategy(fixedVariation);
  const strategy = createStrategy({ type: Strategy.Type.FIXED, fixedStrategy });
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
  description: string,
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

//TODO: should we set the ruleId to empty string as default?
//TODO: create optional constructor for Reason
export function createReason(ruleId: string, type: Reason.TypeMap[keyof Reason.TypeMap]): Reason {
  const reason = new Reason();
  reason.setType(type);
  reason.setRuleId(ruleId);
  return reason;
}
