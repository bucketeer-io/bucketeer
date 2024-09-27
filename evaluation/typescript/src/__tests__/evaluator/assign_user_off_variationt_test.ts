import test from 'ava';
import { createReason, createUser } from '../../modelFactory';
import { Prerequisite } from '../../proto/feature/prerequisite_pb';
import { Reason } from '../../proto/feature/reason_pb';
import { Variation } from '../../proto/feature/variation_pb';
import { newTestFeature, TestVariations } from './evaluate_feature_test';
import { Evaluator } from '../../evaluation';

interface TestCases {
  enabled: boolean
  offVariation: string
  userID: string
  flagVariations: Record<string, string>
  prerequisite: Prerequisite[]
  expectedReason: Reason
  expectedVariation: Variation
  expectedError: Error | null
}

const TestCases: TestCases[] = [
  {
    enabled: false,
    offVariation: 'variation-C',
    userID:'user5',
    flagVariations: {},
    prerequisite: [],
    expectedReason: createReason('', Reason.Type.OFF_VARIATION),
    expectedVariation: TestVariations.variationC,
    expectedError: null,
  }
]

TestCases.forEach((tc, index) => {
  test(`Test Case ${index}`, async (t) => {
    const user = createUser(tc.userID, {})
    const f = newTestFeature('test-feature')
    f.setEnabled(tc.enabled)
    f.setOffVariation(tc.offVariation)
    f.setPrerequisitesList(tc.prerequisite)

    const evaluator = new Evaluator()
    try {
      const [reason, variation] = evaluator.assignUser(f, user, [], tc.flagVariations)
      t.deepEqual(reason, tc.expectedReason)
      t.deepEqual(variation, tc.expectedVariation)
    } catch (err) {
      t.deepEqual(err, tc.expectedError)
    }
  });
});