import { Evaluation, UserEvaluations } from './proto/feature/evaluation_pb';
import { Feature } from './proto/feature/feature_pb';
import * as crypto from 'crypto';
//
function NewUserEvaluations(
  id: string,
  evaluations: Evaluation[],
  archivedFeaturesIds: string[],
  forceUpdate: boolean
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

function sortMapKeys(data: Record<string, string>): string[] {
  const keys = Object.keys(data);
  keys.sort(); // Sort keys alphabetically
  return keys;
}

function GenerateFeaturesID(features: Feature[]): string {
  // Sort features by ID
  features.sort((a, b) => a.getId().localeCompare(b.getId()));

  // Create a new hash (fnv64a replacement)
  const hash = crypto.createHash('sha256');

  // Append feature details to the hash
  features.forEach((feature) => {
    hash.update(`${feature.getId()}:${feature.getVersion()}`);
  });

  // Return the hashed value as a string
  return BigInt('0x' + hash.digest('hex')).toString(10);
}

function UserEvaluationsID(
  userID: string,
  userMetadata: Record<string, string>, // equivalent to map[string]string in Go
  features: Feature[]
): string {
  // Sort features by ID
  features.sort((a, b) => a.getId().localeCompare(b.getId()));

  // Use Node.js crypto module to generate a 64-bit hash (as Go's fnv64a hash)
  const hash = crypto.createHash('sha256'); // Using sha256 as a more common example
  hash.update(userID);

  // Sort userMetadata keys
  const keys = sortMapKeys(userMetadata);
  keys.forEach((key) => {
    hash.update(`${key}:${userMetadata[key]}`);
  });

  // Append feature details to the hash
  features.forEach((feature) => {
    hash.update(`${feature.getId()}:${feature.getVersion()}`);
  });

  // Return the hashed value
  return BigInt('0x' + hash.digest('hex')).toString(10);
}

export {
  UserEvaluationsID,
  GenerateFeaturesID,
  sortMapKeys,
  NewUserEvaluations
}