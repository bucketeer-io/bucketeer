import { Bucketeer } from './bucketeer';
import { RolloutStrategy, Strategy } from './proto/feature/strategy_pb';
import { Variation } from './proto/feature/variation_pb';

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
    const bucketeer = new Bucketeer();
    
    // Audience control: Check if user should be included in experiment
    const audience = strategy.getAudience();
    if (audience !== undefined) {
      if (audience.getPercentage() > 0 && audience.getPercentage() < 100) {
        // Use different hash input for audience control to ensure independence from A/B split
        const trafficInput = `traffic-${featureID}-${userID}-${samplingSeed}`;
        const trafficBucket = bucketeer.bucket(trafficInput);
        const trafficThreshold = audience.getPercentage() / 100.0;
        
        // If user is not in experiment traffic, return default variation
        if (trafficBucket >= trafficThreshold) {
          if (audience.getDefaultVariation() === '') {
            throw new Error('Variation not found');
          }
          return audience.getDefaultVariation();
        }
      }
    }
    
    // Original A/B split logic for users in experiment traffic
    const input = `${featureID}-${userID}-${samplingSeed}`;
    const bucket = bucketeer.bucket(input);

    let rangeEnd = 0.0;
    for (const variation of strategy.getVariationsList()) {
      rangeEnd += variation.getWeight() / 100000.0;
      if (bucket < rangeEnd) {
        return variation.getVariation();
      }
    }
    throw new Error('Variation not found');
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
