import test from 'ava';
import { Reason } from '../../proto/feature/reason_pb';
import { createUser } from '../../modelFactory';
import { newTestFeature } from './evaluate_feature_test';
import { Evaluator } from '../../evaluation';

interface TestCase {
  userID: string;
  expectedReason: Reason.TypeMap[keyof Reason.TypeMap];
  expectedVariationID: string;
}

const TestCases: TestCase[] = [
  {
    userID: 'user1',
    expectedReason: Reason.Type.TARGET,
    expectedVariationID: 'variation-A',
  },
  {
    userID: 'user2',
    expectedReason: Reason.Type.TARGET,
    expectedVariationID: 'variation-B',
  },
  {
    userID: 'user3',
    expectedReason: Reason.Type.TARGET,
    expectedVariationID: 'variation-C',
  },
  {
    userID: 'user4',
    expectedReason: Reason.Type.DEFAULT,
    expectedVariationID: 'variation-B',
  },
];

TestCases.forEach((tc, index) => {
  test(`Test Case ${index}`, async (t) => {
    const user = createUser(tc.userID, {});
    const f = newTestFeature('test-feature');

    const evaluator = new Evaluator();
    const [reason, variation] = evaluator.assignUser(f, user, [], {});
    //TODO: check with deep equal ?
    t.deepEqual(reason.getType(), tc.expectedReason);
    t.deepEqual(variation.getId(), tc.expectedVariationID);
  });
});
