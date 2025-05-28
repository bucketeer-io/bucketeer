import test from 'ava';
import { StrategyEvaluator } from '../strategyEvaluator';
import { Strategy, FixedStrategy, RolloutStrategy, Audience } from '../proto/feature/strategy_pb';
import { Variation } from '../proto/feature/variation_pb';

test('StrategyEvaluator evaluate fixed strategy', (t) => {
  const evaluator = new StrategyEvaluator();
  const variations = [
    createVariation('variation-a', 'a'),
    createVariation('variation-b', 'b'),
  ];
  
  const strategy = new Strategy();
  strategy.setType(Strategy.Type.FIXED);
  const fixedStrategy = new FixedStrategy();
  fixedStrategy.setVariation('variation-a');
  strategy.setFixedStrategy(fixedStrategy);

  const result = evaluator.evaluate(strategy, 'user-1', variations, 'feature-1', 'seed');
  t.is(result.getId(), 'variation-a');
});

test('StrategyEvaluator evaluate rollout strategy without audience', (t) => {
  const evaluator = new StrategyEvaluator();
  const variations = [
    createVariation('variation-a', 'a'),
    createVariation('variation-b', 'b'),
  ];
  
  const strategy = new Strategy();
  strategy.setType(Strategy.Type.ROLLOUT);
  const rolloutStrategy = new RolloutStrategy();
  
  const variationA = new RolloutStrategy.Variation();
  variationA.setVariation('variation-a');
  variationA.setWeight(50000);
  
  const variationB = new RolloutStrategy.Variation();
  variationB.setVariation('variation-b');
  variationB.setWeight(50000);
  
  rolloutStrategy.setVariationsList([variationA, variationB]);
  // No audience configuration (undefined means no audience control)
  strategy.setRolloutStrategy(rolloutStrategy);

  const result = evaluator.evaluate(strategy, 'user-1', variations, 'feature-1', 'seed');
  t.true(result.getId() === 'variation-a' || result.getId() === 'variation-b');
});

test('StrategyEvaluator evaluate rollout strategy with audience', (t) => {
  const evaluator = new StrategyEvaluator();
  const variations = [
    createVariation('variation-a', 'a'),
    createVariation('variation-b', 'b'),
    createVariation('variation-default', 'default'),
  ];
  
  const strategy = new Strategy();
  strategy.setType(Strategy.Type.ROLLOUT);
  const rolloutStrategy = new RolloutStrategy();
  
  const variationA = new RolloutStrategy.Variation();
  variationA.setVariation('variation-a');
  variationA.setWeight(50000);
  
  const variationB = new RolloutStrategy.Variation();
  variationB.setVariation('variation-b');
  variationB.setWeight(50000);
  
  rolloutStrategy.setVariationsList([variationA, variationB]);
  
  // 10% audience configuration
  const audience = new Audience();
  audience.setPercentage(10);
  audience.setDefaultVariation('variation-default');
  rolloutStrategy.setAudience(audience);
  
  strategy.setRolloutStrategy(rolloutStrategy);

  // Test multiple users to verify audience control
  let inExperimentCount = 0;
  let outOfExperimentCount = 0;
  const totalUsers = 1000;

  for (let i = 0; i < totalUsers; i++) {
    const userID = `user-${i}`;
    const result = evaluator.evaluate(strategy, userID, variations, 'feature-1', 'seed');
    
    if (result.getId() === 'variation-default') {
      outOfExperimentCount++;
    } else if (result.getId() === 'variation-a' || result.getId() === 'variation-b') {
      inExperimentCount++;
    } else {
      t.fail(`Unexpected variation ${result.getId()} for user ${userID}`);
    }
  }

  // Verify approximately 10% are in experiment (allow some variance)
  const expectedInExperiment = Math.floor(totalUsers * 10 / 100);
  const tolerance = Math.floor(totalUsers * 5 / 100); // 5% tolerance

  t.true(
    inExperimentCount >= expectedInExperiment - tolerance && 
    inExperimentCount <= expectedInExperiment + tolerance,
    `Expected approximately ${expectedInExperiment} users in experiment, got ${inExperimentCount} (out of ${totalUsers} total)`
  );

  t.true(
    outOfExperimentCount >= totalUsers - expectedInExperiment - tolerance && 
    outOfExperimentCount <= totalUsers - expectedInExperiment + tolerance,
    `Expected approximately ${totalUsers - expectedInExperiment} users out of experiment, got ${outOfExperimentCount} (out of ${totalUsers} total)`
  );
});

test('StrategyEvaluator evaluate rollout strategy with audience but no default variation', (t) => {
  const evaluator = new StrategyEvaluator();
  const variations = [
    createVariation('variation-a', 'a'),
    createVariation('variation-b', 'b'),
  ];
  
  const strategy = new Strategy();
  strategy.setType(Strategy.Type.ROLLOUT);
  const rolloutStrategy = new RolloutStrategy();
  
  const variationA = new RolloutStrategy.Variation();
  variationA.setVariation('variation-a');
  variationA.setWeight(50000);
  
  const variationB = new RolloutStrategy.Variation();
  variationB.setVariation('variation-b');
  variationB.setWeight(50000);
  
  rolloutStrategy.setVariationsList([variationA, variationB]);
  
  const audience = new Audience();
  audience.setPercentage(10);
  audience.setDefaultVariation(''); // No default variation
  rolloutStrategy.setAudience(audience);
  
  strategy.setRolloutStrategy(rolloutStrategy);

  // Find a user that would be outside the experiment traffic
  let foundError = false;
  for (let i = 0; i < 100; i++) {
    const userID = `user-${i}`;
    try {
      evaluator.evaluate(strategy, userID, variations, 'feature-1', 'seed');
    } catch (error) {
      if (error instanceof Error && error.message === 'Variation not found') {
        foundError = true;
        break;
      } else {
        t.fail(`Unexpected error for user ${userID}: ${error}`);
      }
    }
  }
  
  t.true(foundError, 'Expected to find at least one user that gets "Variation not found" error');
});

test('StrategyEvaluator evaluate rollout strategy with 100% audience', (t) => {
  const evaluator = new StrategyEvaluator();
  const variations = [
    createVariation('variation-a', 'a'),
    createVariation('variation-b', 'b'),
    createVariation('variation-default', 'default'),
  ];
  
  const strategy = new Strategy();
  strategy.setType(Strategy.Type.ROLLOUT);
  const rolloutStrategy = new RolloutStrategy();
  
  const variationA = new RolloutStrategy.Variation();
  variationA.setVariation('variation-a');
  variationA.setWeight(50000);
  
  const variationB = new RolloutStrategy.Variation();
  variationB.setVariation('variation-b');
  variationB.setWeight(50000);
  
  rolloutStrategy.setVariationsList([variationA, variationB]);
  
  // 100% audience configuration
  const audience = new Audience();
  audience.setPercentage(100);
  audience.setDefaultVariation('variation-default');
  rolloutStrategy.setAudience(audience);
  
  strategy.setRolloutStrategy(rolloutStrategy);

  // With 100% audience, all users should be in experiment
  for (let i = 0; i < 10; i++) {
    const userID = `user-${i}`;
    const result = evaluator.evaluate(strategy, userID, variations, 'feature-1', 'seed');
    
    // Should never get default variation with 100% audience
    t.not(result.getId(), 'variation-default', `Unexpected default variation for user ${userID} with 100% audience`);
    t.true(
      result.getId() === 'variation-a' || result.getId() === 'variation-b',
      `Expected variation-a or variation-b for user ${userID}, got ${result.getId()}`
    );
  }
});

function createVariation(id: string, value: string): Variation {
  const variation = new Variation();
  variation.setId(id);
  variation.setValue(value);
  return variation;
} 