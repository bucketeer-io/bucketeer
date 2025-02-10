import test from 'ava';
import { createUser } from '../../modelFactory';
import { newTestFeature } from './evaluate_feature_test';
import { Evaluator } from '../../evaluation';
import { Reason } from '../../proto/feature/reason_pb';

test('no default strategy', (t) => {
  const user = createUser('user-id1', { name3: 'user3' });
  const f = newTestFeature('test-feature');
  f.clearDefaultStrategy();
  const evalator = new Evaluator();
  try {
    evalator.assignUser(f, user, [], {});
    t.fail('should throw an error "evaluator: default strategy not found"');
  } catch (error) {
    t.deepEqual(error, new Error('evaluator: default strategy not found'));
  }
});

test('with default strategy', (t) => {
  const user = createUser('user-id1', { name3: 'user3' });
  const f = newTestFeature('test-feature');
  const evalator = new Evaluator();
  const [reason, variation] = evalator.assignUser(f, user, [], {});
  t.is(reason.getType(), Reason.Type.DEFAULT);
  t.is(variation.getId(), 'variation-B');
});
