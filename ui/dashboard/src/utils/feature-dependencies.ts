import { Feature, FeatureRuleClauseOperator } from '@types';

/**
 * Returns the IDs of features that this feature depends on.
 * Checks both prerequisites and rules with FEATURE_FLAG clauses.
 */
const getFeatureIdsDependsOn = (feature: Feature): string[] => {
  const ids: string[] = [];

  // Check prerequisites
  for (const prerequisite of feature.prerequisites) {
    ids.push(prerequisite.featureId);
  }

  // Check rules for FEATURE_FLAG clauses
  for (const rule of feature.rules) {
    for (const clause of rule.clauses) {
      if (clause.operator === FeatureRuleClauseOperator.FEATURE_FLAG) {
        // For FEATURE_FLAG clauses, the attribute contains the feature ID
        ids.push(clause.attribute);
      }
    }
  }

  return ids;
};

/**
 * Checks if any feature (not in targets) depends on any of the target features.
 * Used to prevent archiving/deleting flags that other flags depend on.
 *
 * @param targets - The features being checked (e.g., flag to archive)
 * @param allFeatures - All features in the environment
 * @returns true if any other feature depends on any target
 */
export const hasDependentFlags = (
  targets: Feature[],
  allFeatures: Feature[]
): boolean => {
  // Create a set of target IDs for quick lookup
  const targetIds = new Set(targets.map(t => t.id));

  // Check if any feature (not in targets) depends on any target
  for (const feature of allFeatures) {
    // Skip if this feature is one of the targets
    if (targetIds.has(feature.id)) {
      continue;
    }

    // Check if this feature depends on any target
    const dependsOnIds = getFeatureIdsDependsOn(feature);
    for (const depId of dependsOnIds) {
      if (targetIds.has(depId)) {
        return true;
      }
    }
  }

  return false;
};

export const getDependentFlags = (
  targets: Feature[],
  allFeatures: Feature[]
): Feature[] => {
  const targetIds = new Set(targets.map(t => t.id));
  const dependentFlags: Feature[] = [];

  for (const feature of allFeatures) {
    if (targetIds.has(feature.id)) {
      continue;
    }

    const dependsOnIds = getFeatureIdsDependsOn(feature);
    const dependsOnTarget = dependsOnIds.some(depId => targetIds.has(depId));

    if (dependsOnTarget) {
      dependentFlags.push(feature);
    }
  }

  return dependentFlags;
};
