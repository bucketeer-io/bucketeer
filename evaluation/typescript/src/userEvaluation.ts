import { Evaluation, UserEvaluations } from './proto/feature/evaluation_pb';
import { Feature } from './proto/feature/feature_pb';
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

function sortMapKeys(data: Record<string, string>): string[] {
  const keys = Object.keys(data);
  keys.sort(); // Sort keys alphabetically
  return keys;
}

function GenerateFeaturesID(features: Feature[]): string {
  // Sort features based on the 'id'
  features.sort((a, b) => (a.getId() < b.getId() ? -1 : 1));

  // Initialize FNV-1a 64-bit hash string
  let hashInput = '';

  // Concatenate each feature's 'id' and 'version' into the hash input
  features.forEach((feature) => {
    hashInput += `${feature.getId()}:${feature.getVersion()}`;
  });

  // Generate the FNV-1a 64-bit hash and return it as a decimal string
  const hash = fnv.hash(hashInput, 64);
  return hash.dec();
}

function UserEvaluationsID(
  userID: string,
  userMetadata: Record<string, string>, // equivalent to map[string]string in Go
  features: Feature[],
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

  // Generate the FNV-1a 64-bit hash
  const hash = fnv.hash(hashInput, 64);

  // Return the hashed value as a decimal string
  return hash.dec();
}

export { UserEvaluationsID, GenerateFeaturesID, sortMapKeys, NewUserEvaluations };
