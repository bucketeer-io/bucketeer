import * as crypto from 'crypto';
import { RolloutStrategy, Strategy } from './proto/feature/strategy_pb';
import { Variation } from './proto/feature/variation_pb';
//
const MAX = 0xffffffffffffffffn;

class StrategyEvaluator {
  evaluate(
    strategy: Strategy,
    userID: string,
    variations: Variation[],
    featureID: string,
    samplingSeed: string,
  ): Variation {
    switch (strategy.getType()) {
      case Strategy.Type.FIXED:
        return this.findVariation(strategy.getFixedStrategy()?.getVariation() || '', variations);
      case Strategy.Type.ROLLOUT:
        const rolloutStrategy = strategy.getRolloutStrategy();
        if (rolloutStrategy !== undefined) {
          const variationID = this.rollout(rolloutStrategy, userID, featureID, samplingSeed);
          return this.findVariation(variationID, variations);
        }
        throw new Error('Missing rollout strategy');
      default:
        throw new Error('Unsupported strategy');
    }
  }

  private rollout(
    strategy: RolloutStrategy,
    userID: string,
    featureID: string,
    samplingSeed: string,
  ): string {
    const bucket = this.bucket(userID, featureID, samplingSeed);

    let sum = 0.0;
    for (const variation of strategy.getVariationsList()) {
      sum += variation.getWeight() / 100000.0;
      if (bucket < sum) {
        return variation.getVariation();
      }
    }
    throw new Error('Variation not found');
  }

  private bucket(userID: string, featureID: string, samplingSeed: string): number {
    try {
      const hash = this.hash(userID, featureID, samplingSeed);
      const intVal = BigInt('0x' + hash.slice(0, 16)); // Convert the first 16 hex characters to BigInt

      // Divide the BigInt value by `max` and convert it to a number. Use Number() since we need a float.
      return Number(intVal) / Number(MAX);
    } catch (error) {
      console.error('Failed to calculate bucket value:', error);
      throw error;
    }
  }

  private hash(userID: string, featureID: string, samplingSeed: string): string {
    const concat = `${featureID}-${userID}${samplingSeed}`;
    return crypto.createHash('md5').update(concat).digest('hex');
  }

  private findVariation(variationID: string, variations: Variation[]): Variation {
    for (const variation of variations) {
      if (variation.getId() === variationID) {
        return variation;
      }
    }
    throw new Error('Variation not found');
  }
}

export { StrategyEvaluator };
