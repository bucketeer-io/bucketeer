import { Evaluation, UserEvaluations } from './proto/feature/evaluation_pb';
import { Feature } from './proto/feature/feature_pb';
import { UNAVAILABLE_UEID_PREFIX, holdoutUnavailable } from './bucketeer';
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

function sortedHoldoutTargets(group: HoldoutGroup): Array<{
  featureId: string;
  variationId: string;
  addedAt: number;
  id: string;
}> {
  const features = group.getFeaturesList?.() ?? [];
  return features
    .filter((feature) => feature && feature.getFeatureId() !== '' && feature.getRemovedAt() === 0)
    .map((feature) => ({
      featureId: feature.getFeatureId(),
      variationId: feature.getHoldoutVariationId(),
      addedAt: feature.getAddedAt(),
      id: feature.getId(),
    }))
    .sort((a, b) => {
      if (a.featureId !== b.featureId) {
        return a.featureId < b.featureId ? -1 : 1;
      }
      return a.id < b.id ? -1 : a.id > b.id ? 1 : 0;
    });
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

// holdoutGroupsHashInput is the part of a change-detection hash that covers
// the live holdout groups. A group with no ID adds nothing, so an ID computed
// without a live group is the ID computed before holdout existed.
//
// Targets are part of the hash: two edits in the same second share UpdatedAt,
// and the ID still has to change when a target or its holdout variation changes.
// Matches writeHoldoutGroups in evaluation/go.
function holdoutGroupsHashInput(groups: HoldoutGroup[] | undefined): string {
  let hashInput = '';
  sortedHoldoutGroups(groups).forEach((group) => {
    hashInput += `${group.getId()}:${group.getUpdatedAt()}`;
    sortedHoldoutTargets(group).forEach((target) => {
      hashInput += `:${target.featureId}:${target.variationId}:${target.addedAt}:${target.id}`;
    });
  });
  return hashInput;
}

// FEATURES_ID_HOLDOUT_SEPARATOR separates the flags part of a features ID from
// its holdout part. Mirrors featuresIDHoldoutSeparator in Go.
const FEATURES_ID_HOLDOUT_SEPARATOR = ':';

// GenerateFeaturesID generates a hash ID for the server SDK to detect feature flag changes.
// It uses FNV-1a 64-bit hash algorithm which is:
// - Fast enough for SDK polling frequency
// - 64-bit provides negligible collision risk for feature flag counts
// - No cryptographic requirements for change detection
//
// A locally evaluating SDK caches the live holdout groups next to the flags, so
// starting or stopping a group, or changing its targets, has to change the ID
// even though no flag changed. The ID therefore has two parts,
// "<flags>:<holdout>". Without a live group it is the flags part alone, which
// is the ID computed before holdout existed.
//
// The parts are kept apart, rather than hashed together, so the server can tell
// a holdout-only change from a flag change. Matches GenerateFeaturesID in
// evaluation/go.
function GenerateFeaturesID(features: Feature[], holdoutGroups?: HoldoutGroup[]): string {
  // Sort features based on the 'id'
  features.sort((a, b) => (a.getId() < b.getId() ? -1 : 1));

  // Initialize FNV-1a 64-bit hash string
  let hashInput = '';

  // Concatenate each feature's 'id' and 'updatedAt' into the hash input
  features.forEach((feature) => {
    hashInput += `${feature.getId()}:${feature.getUpdatedAt()}`;
  });

  // Generate the FNV-1a 64-bit hash as a decimal string
  const flagsID = fnv.hash(hashInput, 64).dec();
  if (sortedHoldoutGroups(holdoutGroups).length === 0) {
    return flagsID;
  }
  // The holdout part hashes the same input as UserEvaluationsID.
  const holdoutID = fnv.hash(holdoutGroupsHashInput(holdoutGroups), 64).dec();
  return flagsID + FEATURES_ID_HOLDOUT_SEPARATOR + holdoutID;
}

// FeaturesIDFlagsPart returns the part of a features ID that identifies the
// flags alone. Two IDs with the same flags part were issued for the same flags
// and differ only in the holdout groups. Matches FeaturesIDFlagsPart in
// evaluation/go.
function FeaturesIDFlagsPart(id: string): string {
  const index = id.indexOf(FEATURES_ID_HOLDOUT_SEPARATOR);
  return index === -1 ? id : id.slice(0, index);
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
  hashInput += holdoutGroupsHashInput(holdoutGroups);

  // Generate the FNV-1a 64-bit hash
  const hash = fnv.hash(hashInput, 64);

  // An ID issued without a readable holdout snapshot is marked so the next
  // evaluation after the snapshot is back runs in full.
  if (holdoutUnavailable(holdoutGroups)) {
    return UNAVAILABLE_UEID_PREFIX + hash.dec();
  }
  // Return the hashed value as a decimal string
  return hash.dec();
}

export {
  UserEvaluationsID,
  GenerateFeaturesID,
  FeaturesIDFlagsPart,
  sortMapKeys,
  NewUserEvaluations,
};
