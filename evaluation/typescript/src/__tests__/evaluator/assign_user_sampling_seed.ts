import test from 'ava';
import { createRolloutStrategy, createStrategy, createUser } from '../../modelFactory';
import { newTestFeature } from './evaluate_feature_test';
import { Evaluator } from '../../evaluation';
import { Reason } from '../../proto/feature/reason_pb';
import { Strategy } from '../../proto/feature/strategy_pb';

test('assign user with sampling seed', (t) => {
  const user = createUser('uid', {});
  const f = newTestFeature('fid');
  const rolloutStrategy = createRolloutStrategy({
    variations: [
      { variation: 'variation-A', weight: 30000 },
      { variation: 'variation-B', weight: 40000 },
      { variation: 'variation-C', weight: 30000 },
    ],
  });
  const strategy = createStrategy({
    rolloutStrategy: rolloutStrategy,
    type: Strategy.Type.ROLLOUT,
  });
  f.setDefaultStrategy(strategy);

  const evalator = new Evaluator();
  let [reason, variation] = evalator.assignUser(f, user, [], {});
  t.is(reason.getType(), Reason.Type.DEFAULT);
  t.is(variation.getId(), 'variation-B'); //rolloutStrategy.getVariationsList()[1].getVariation()

  // Channge sampling seed to change assigned variation.
  f.setSamplingSeed('test');
  [reason, variation] = evalator.assignUser(f, user, [], {});
  t.is(reason.getType(), Reason.Type.DEFAULT);
  t.is(variation.getId(), 'variation-A'); //rolloutStrategy.getVariationsList()[0].getVariation()
});
