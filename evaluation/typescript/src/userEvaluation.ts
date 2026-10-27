import { Evaluation, UserEvaluations } from './proto/feature/evaluation_pb';
import { Feature } from './proto/feature/feature_pb';
import type { HoldoutGroup } from './bucketeer';
import fnv from 'fnv-plus';

//
function NewUserEvaluations(
  id: string,
  evaluations: Evaluation[],
  archivedFeaturesIds: string[],
  forceUpdate: boolean,
): UserEvaluations {
  const now = Math.floor(Date.now() / 1000); // Equivalent to Unix time

  const userEvaluations = new UserEvaluations();
  userEvaluations.setId(id);
  userEvaluations.setEvaluationsList(evaluations);
  userEvaluations.setCreatedAt(now);
  userEvaluations.setArchivedFeatureIdsList(archivedFeaturesIds);
  userEvaluations.setForceUpdate(forceUpdate);

  return userEvaluations;
}

function sortedHoldoutTargets(
  group: HoldoutGroup,
): Array<{ featureId: string; variationId: string }> {
  const features = group.getFeaturesList?.() ?? [];
  return features
    .filter((feature) => feature && feature.getFeatureId() !== '' && feature.getRemovedAt() === 0)
    .map((feature) => ({
      featureId: feature.getFeatureId(),
      variationId: feature.getHoldoutVariationId(),
    }))
    .sort((a, b) => (a.featureId < b.featureId ? -1 : a.featureId > b.featureId ? 1 : 0));
}

function sortedHoldoutGroups(groups: HoldoutGroup[] | undefined): HoldoutGroup[] {
  if (!groups || groups.length === 0) {
    return [];
  }
  return groups
    .filter((group) => group && group.getId() !== '')
    .sort((a, b) => (a.getId() < b.getId() ? -1 : a.getId() > b.getId() ? 1 : 0));
}

function sortMapKeys(data: Record<string, string>): string[] {
  const keys = Object.keys(data);
  keys.sort(); // Sort keys alphabetically
  return keys;
}

// GenerateFeaturesID generates a hash ID for the server SDK to detect feature flag changes.
// It uses FNV-1a 64-bit hash algorithm which is:
// - Fast enough for SDK polling frequency
// - 64-bit provides negligible collision risk for feature flag counts
// - No cryptographic requirements for change detection
function GenerateFeaturesID(features: Feature[]): string {
  // Sort features based on the 'id'
  features.sort((a, b) => (a.getId() < b.getId() ? -1 : 1));

  // Initialize FNV-1a 64-bit hash string
  let hashInput = '';

  // Concatenate each feature's 'id' and 'updatedAt' into the hash input
  features.forEach((feature) => {
    hashInput += `${feature.getId()}:${feature.getUpdatedAt()}`;
  });

  // Generate the FNV-1a 64-bit hash and return it as a decimal string
  const hash = fnv.hash(hashInput, 64);
  return hash.dec();
}

// UserEvaluationsID generates a hash ID for the client SDK to detect evaluation changes.
// It uses FNV-1a 64-bit hash algorithm which is:
// - Fast enough for SDK polling frequency
// - 64-bit provides negligible collision risk for feature flag counts
// - No cryptographic requirements for change detection
function UserEvaluationsID(
  userID: string,
  userMetadata: Record<string, string>, // equivalent to map[string]string in Go
  features: Feature[],
  holdoutGroups?: HoldoutGroup[],
): string {
  // Sort features by ID
  features.sort((a, b) => (a.getId() < b.getId() ? -1 : 1));

  // Initialize FNV-1a 64-bit hash input
  let hashInput = userID;

  // Sort and append userMetadata to the hash input
  const keys = sortMapKeys(userMetadata);
  keys.forEach((key) => {
    hashInput += `${key}:${userMetadata[key]}`;
  });

  // Append feature details to the hash input
  features.forEach((feature) => {
    hashInput += `${feature.getId()}:${feature.getUpdatedAt()}`;
  });

  // Hash live groups after features so a missing group leaves the ID unchanged.
  // Targets are part of the hash: two edits in the same second share UpdatedAt,
  // and the ID still has to change when a target or its holdout variation changes.
  sortedHoldoutGroups(holdoutGroups).forEach((group) => {
    hashInput += `${group.getId()}:${group.getUpdatedAt()}`;
    sortedHoldoutTargets(group).forEach((target) => {
      hashInput += `:${target.featureId}:${target.variationId}`;
    });
  });

  // Generate the FNV-1a 64-bit hash
  const hash = fnv.hash(hashInput, 64);

  // Return the hashed value as a decimal string
  return hash.dec();
}

export { UserEvaluationsID, GenerateFeaturesID, sortMapKeys, NewUserEvaluations };
