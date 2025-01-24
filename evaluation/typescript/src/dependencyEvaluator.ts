import { EVALUATOR_ERRORS } from './errors';

class DependencyEvaluator {
  evaluate(
    featureID: string,
    variationIDs: string[],
    flagVariations: { [key: string]: string },
  ): boolean {
    const targetVarID = flagVariations[featureID];

    if (!targetVarID) {
      throw EVALUATOR_ERRORS.FeatureNotFound; // Throwing an error if feature is not found
    }

    for (const varID of variationIDs) {
      if (varID === targetVarID) {
        return true; // Match found
      }
    }
    return false; // No match found
  }
}

export { DependencyEvaluator };
