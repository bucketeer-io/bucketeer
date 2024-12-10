import test from 'ava';
import { createUser } from '../../modelFactory';
import { newTestFeature } from './evaluate_feature_test';
import { Evaluator } from '../../evaluation';

test('assign user rule set', (t) => {
  const user = createUser('user-id', { name: 'user3' });
  const f = newTestFeature('test-feature');
  const evalator = new Evaluator();
  const [reason, variation] = evalator.assignUser(f, user, [], {});
  t.is(reason.getRuleId(), 'rule-2');
  t.is(variation.getId(), 'variation-B');
});
