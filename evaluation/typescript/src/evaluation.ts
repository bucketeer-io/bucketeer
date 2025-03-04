import { EVALUATOR_ERRORS } from './errors';
import { Clause } from './proto/feature/clause_pb';
import { Evaluation, UserEvaluations } from './proto/feature/evaluation_pb';
import { Feature } from './proto/feature/feature_pb';
import { Reason } from './proto/feature/reason_pb';
import { SegmentUser } from './proto/feature/segment_pb';
import { Variation } from './proto/feature/variation_pb';
import { User } from './proto/user/user_pb';
import { RuleEvaluator } from './ruleEvaluator';
import { StrategyEvaluator } from './strategyEvaluator';
import { NewUserEvaluations, UserEvaluationsID } from './userEvaluation';
import { createReason } from './modelFactory';

const SECONDS_TO_RE_EVALUATE_ALL = 30 * 24 * 60 * 60; // 30 days
const SECONDS_FOR_ADJUSTMENT = 10; // 10 seconds

export function EvaluationID(featureID: string, featureVersion: number, userID: string): string {
  return `${featureID}:${featureVersion}:${userID}`;
}

class Evaluator {
  private ruleEvaluator: RuleEvaluator;
  private strategyEvaluator: StrategyEvaluator;

  constructor() {
    this.ruleEvaluator = new RuleEvaluator();
    this.strategyEvaluator = new StrategyEvaluator();
  }

  async evaluateFeatures(
    features: Feature[],
    user: User,
    mapSegmentUsers: Map<string, SegmentUser[]>,
    targetTag: string,
  ): Promise<UserEvaluations> {
    return this.evaluate(features, user, mapSegmentUsers, false, targetTag);
  }

  // GetPrerequisiteDownwards gets the features specified as prerequisite by the targetFeatures.
  getPrerequisiteDownwards(targetFeatures: Feature[], allFeatures: Feature[]): Feature[] {
    const allFeaturesMap = new Map<string, Feature>();
    for (const feature of allFeatures) {
      allFeaturesMap.set(feature.getId(), feature);
    }
    const dependedOnTargets = getFeaturesDependedOnTargets(targetFeatures, allFeaturesMap);
    return Object.values(dependedOnTargets);
  }

  evaluateFeaturesByEvaluatedAt(
    features: Feature[],
    user: User,
    mapSegmentUsers: Map<string, SegmentUser[]>,
    prevUEID: string,
    evaluatedAt: number,
    userAttributesUpdated: boolean,
    targetTag: string,
  ): UserEvaluations {
    if (prevUEID === '') {
      return this.evaluate(features, user, mapSegmentUsers, true, targetTag);
    }

    const now = Math.floor(Date.now() / 1000);
    if (evaluatedAt < now - SECONDS_TO_RE_EVALUATE_ALL) {
      return this.evaluate(features, user, mapSegmentUsers, true, targetTag);
    }

    const adjustedEvalAt = evaluatedAt - SECONDS_FOR_ADJUSTMENT;
    const updatedFeatures = features.filter(
      (feature) =>
        feature.getUpdatedAt() > adjustedEvalAt ||
        (userAttributesUpdated && feature.getRulesList().length > 0),
    );

    // If the UserEvaluationsID has changed, but both User Attributes and Feature Flags have not been updated,
    // it is considered unusual and a force update should be performed.
    if (updatedFeatures.length === 0) {
      return this.evaluate(features, user, mapSegmentUsers, true, targetTag);
    }

    const evalTargets = this.getEvalFeatures(updatedFeatures, features);
    return this.evaluate(evalTargets, user, mapSegmentUsers, false, targetTag);
  }

  private evaluate(
    features: Feature[],
    user: User,
    mapSegmentUsers: Map<string, SegmentUser[]>,
    forceUpdate: boolean,
    targetTag: string,
  ): UserEvaluations {
    const flagVariations: { [key: string]: string } = {};
    // fs need to be sorted in order from upstream to downstream.
    const sortedFeatures = topologicalSort(features);

    const evaluations: Evaluation[] = [];
    const archivedIDs: string[] = [];

    for (const feature of sortedFeatures) {
      if (feature.getArchived()) {
        // To keep response size small, the feature flags archived long time ago are excluded.
        if (!this.isArchivedBeforeLastThirtyDays(feature)) {
          archivedIDs.push(feature.getId());
        }
        continue;
      }

      const segmentUsers = this.listSegmentIDs(feature).flatMap(
        (id) => mapSegmentUsers.get(id) || [],
      );

      const [reason, variation] = this.assignUser(feature, user, segmentUsers, flagVariations);
      // VariationId is used to check if prerequisite flag's result is what user expects it to be.
      flagVariations[feature.getId()] = variation.getId();
      // When the tag is set in the request,
      // it will return only the evaluations of flags that match the tag configured on the dashboard.
      // When empty, it will return all the evaluations of the flags in the environment.
      if (targetTag !== '' && !this.tagExist(feature.getTagsList(), targetTag)) {
        continue;
      }

      const evaluationID = EvaluationID(feature.getId(), feature.getVersion(), user.getId());
      const evaluation = new Evaluation();
      evaluation.setId(evaluationID);
      evaluation.setFeatureId(feature.getId());
      evaluation.setFeatureVersion(feature.getVersion());
      evaluation.setUserId(user.getId());
      evaluation.setVariationId(variation.getId());
      evaluation.setVariationName(variation.getName());
      evaluation.setVariationValue(variation.getValue());
      // Deprecated
      // FIXME: Remove the Variation when is no longer being used.
      // For security reasons, we should remove the variation description.
      // We copy the variation object to avoid race conditions when removing
      // the description directly from the `variation`
      const copyVariation = new Variation();
      copyVariation.setId(variation.getId());
      copyVariation.setName(variation.getName());
      copyVariation.setValue(variation.getValue());
      evaluation.setVariation(copyVariation);
      evaluation.setReason(reason);

      evaluations.push(evaluation);
    }

    const id = UserEvaluationsID(
      user.getId(),
      arrayToRecord(user.getDataMap().toArray()),
      features,
    );
    const userEvaluations = NewUserEvaluations(id, evaluations, archivedIDs, forceUpdate);
    return userEvaluations;
  }

  private tagExist(tags: string[], target: string): boolean {
    return tags.includes(target);
  }

  /*
  IsArchivedBeforeLastThirtyDays returns a bool value
  indicating whether the feature flag was archived within the last thirty days.
  */
  private isArchivedBeforeLastThirtyDays(feature: Feature): boolean {
    if (!feature.getArchived()) {
      return false;
    }
    const now = Math.floor(Date.now() / 1000);
    return feature.getUpdatedAt() < now - SECONDS_TO_RE_EVALUATE_ALL;
  }

  listSegmentIDs(feature: Feature): string[] {
    const mapIDs = new Set<string>();
    for (const rule of feature.getRulesList()) {
      for (const clause of rule.getClausesList()) {
        if (clause.getOperator() === Clause.Operator.SEGMENT) {
          clause.getValuesList().forEach((value) => mapIDs.add(value));
        }
      }
    }
    return Array.from(mapIDs);
  }

  assignUser(
    feature: Feature,
    user: User,
    segmentUsers: SegmentUser[],
    flagVariations: { [key: string]: string },
  ): [Reason, Variation] {
    for (const pf of feature.getPrerequisitesList()) {
      const variationId = flagVariations[pf.getFeatureId()];
      if (!variationId) {
        throw EVALUATOR_ERRORS.PrerequisiteVariationNotFound;
      }
      if (pf.getVariationId() !== variationId) {
        if (feature.getOffVariation() !== '') {
          const variation = this.findVariation(
            feature.getOffVariation(),
            feature.getVariationsList(),
          );
          const reason = createReason('', Reason.Type.PREREQUISITE);
          return [reason, variation];
        }
      }
    }
    // It doesn't assign the user in case of the feature is disabled and OffVariation is not set
    if (!feature.getEnabled() && feature.getOffVariation() !== '') {
      const variation = this.findVariation(feature.getOffVariation(), feature.getVariationsList());
      const reason = createReason('', Reason.Type.OFF_VARIATION);
      return [reason, variation];
    }

    // evaluate from top to bottom, return if one rule matches
    // evaluate targeting rules
    for (const target of feature.getTargetsList()) {
      if (target.getUsersList().includes(user.getId())) {
        const variation = this.findVariation(target.getVariation(), feature.getVariationsList());
        const reason = createReason('', Reason.Type.TARGET);
        return [reason, variation];
      }
    }

    // evaluate ruleset
    const rule = this.ruleEvaluator.evaluate(
      feature.getRulesList(),
      user,
      segmentUsers,
      flagVariations,
    );
    if (rule) {
      const strategy = rule.getStrategy();
      if (strategy == undefined) {
        throw EVALUATOR_ERRORS.StrategyNotFound;
      }
      const variation = this.strategyEvaluator.evaluate(
        strategy,
        user.getId(),
        feature.getVariationsList(),
        feature.getId(),
        feature.getSamplingSeed(),
      );
      const reason = createReason(rule.getId(), Reason.Type.RULE);
      return [reason, variation];
    }

    // use default strategy
    const defaultStrategy = feature.getDefaultStrategy();
    if (defaultStrategy == undefined) {
      throw EVALUATOR_ERRORS.DefaultStrategyNotFound;
    }

    const variation = this.strategyEvaluator.evaluate(
      defaultStrategy,
      user.getId(),
      feature.getVariationsList(),
      feature.getId(),
      feature.getSamplingSeed(),
    );
    const reason = createReason('', Reason.Type.DEFAULT);
    return [reason, variation];
  }

  getEvalFeatures(targetFeatures: Feature[], allFeatures: Feature[]): Feature[] {
    const allFeaturesMap = new Map<string, Feature>();
    for (const feature of allFeatures) {
      allFeaturesMap.set(feature.getId(), feature);
    }

    const evals1 = getFeaturesDependedOnTargets(targetFeatures, allFeaturesMap);
    const evals2 = getFeaturesDependsOnTargets(targetFeatures, allFeaturesMap);
    const mergedEvals = { ...evals1, ...evals2 };
    return Object.values(mergedEvals);
  }

  private findVariation(v: string, variations: Variation[]): Variation {
    const variation = variations.find((variation) => variation.getId() === v);
    if (!variation) {
      throw EVALUATOR_ERRORS.VariationNotFound;
    }
    return variation;
  }
}

enum Mark {
  Unvisited,
  Temporary,
  Permanently,
}

// FeatureIDsDependsOn returns the ids of the features that this feature depends on.
function getFeatureIDsDependsOn(feature: Feature): Array<string> {
  const ids: Array<string> = [];

  // Iterate over prerequisites and add their FeatureId
  feature.getPrerequisitesList().forEach((p) => {
    ids.push(p.getFeatureId());
  });

  // Iterate over rules and collect ids from clauses where the operator is FEATURE_FLAG
  feature.getRulesList().forEach((rule) => {
    rule.getClausesList().forEach((clause) => {
      if (clause.getOperator() === Clause.Operator.FEATURE_FLAG) {
        ids.push(clause.getAttribute());
      }
    });
  });

  return ids;
}

// Error types
class ErrCycleExists extends Error {
  constructor() {
    super('Cycle exists in the graph');
  }
}

class ErrFeatureNotFound extends Error {
  constructor() {
    super('Feature not found');
  }
}

// Topological sort function
function topologicalSort(features: Array<Feature>): Array<Feature> {
  const marks: { [key: string]: Mark } = {};
  const mapFeatures: { [key: string]: Feature } = {};

  features.forEach((f) => {
    marks[f.getId()] = Mark.Unvisited;
    mapFeatures[f.getId()] = f;
  });

  const sortedFeatures: Array<Feature> = [];

  const sort = (f: Feature) => {
    const fId = f.getId();
    if (marks[fId] === Mark.Permanently) return;

    if (marks[fId] === Mark.Temporary) {
      throw new ErrCycleExists();
    }

    marks[fId] = Mark.Temporary;

    const dependentFeatureIds = getFeatureIDsDependsOn(f);
    for (const fid of dependentFeatureIds) {
      const pf = mapFeatures[fid];
      if (!pf) {
        throw new ErrFeatureNotFound();
      }

      sort(pf);
    }

    marks[fId] = Mark.Permanently;
    sortedFeatures.push(f);
  };

  // Process each feature
  for (const f of features) {
    if (marks[f.getId()] === Mark.Unvisited) {
      sort(f);
    }
  }

  return sortedFeatures;
}

// getFeaturesDependedOnTargets returns the features that are depended on the target features.
// targetFeatures are included in the result.
function getFeaturesDependedOnTargets(
  targets: Array<Feature>,
  all: Map<string, Feature>,
): { [key: string]: Feature } {
  const evals: { [key: string]: Feature } = {};

  const dfs = (f: Feature): void => {
    // If the feature is already in the evals map, skip
    if (evals[f.getId()]) return;

    // Add feature to evals
    evals[f.getId()] = f;

    // Get dependent features recursively
    const featureDependencies = getFeatureIDsDependsOn(f);
    featureDependencies.forEach((fid) => {
      const target = all.get(fid);
      if (target !== undefined) {
        dfs(target);
      }
    });
  };

  // Start DFS for each target feature
  targets.forEach((f) => dfs(f));

  return evals;
}

// getFeaturesDependsOnTargets returns the features that depend on the target features.
// targetFeatures are included in the result.
function getFeaturesDependsOnTargets(
  targets: Array<Feature>,
  all: Map<string, Feature>,
): { [key: string]: Feature } {
  const evals: { [key: string]: Feature } = {};

  // Mark target features in the evals map initially
  targets.forEach((f) => {
    evals[f.getId()] = f;
  });

  // Depth-first search to determine if a feature depends on a target feature
  const dfs = (f: Feature): boolean => {
    if (evals[f.getId()]) {
      return true;
    }

    const featureDependencies = getFeatureIDsDependsOn(f);
    for (const fid of featureDependencies) {
      const dependentFeature = all.get(fid);
      if (dependentFeature && dfs(dependentFeature)) {
        // If the feature depends on one of the target features, add it to evals
        evals[f.getId()] = f;
        return true;
      }
    }
    return false;
  };

  // Apply DFS for all features, except target features
  all.forEach((f) => {
    dfs(f);
  });

  return evals;
}

function arrayToRecord(arr: Array<[string, string]>): Record<string, string> {
  return arr.reduce(
    (acc, [key, value]) => {
      acc[key] = value;
      return acc;
    },
    {} as Record<string, string>,
  );
}

export { Evaluator, getFeatureIDsDependsOn };
