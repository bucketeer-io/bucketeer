import test from 'ava';
import { createPrerequisite, createReason, createUser } from '../../modelFactory';
import { Prerequisite } from '../../proto/feature/prerequisite_pb';
import { Reason } from '../../proto/feature/reason_pb';
import { Variation } from '../../proto/feature/variation_pb';
import { newTestFeature, TestVariations } from './evaluate_feature_test';
import { Evaluator } from '../../evaluation';

interface TestCases {
  enabled: boolean;
  offVariation: string;
  userID: string;
  flagVariations: Record<string, string>;
  prerequisite: Prerequisite[];
  expectedReason: Reason;
  expectedVariation: Variation | null;
  expectedError: Error | null;
}

const TestCases: TestCases[] = [
  {
    enabled: false,
    offVariation: 'variation-C',
    userID: 'user5',
    flagVariations: {},
    prerequisite: [],
    expectedReason: createReason('', Reason.Type.OFF_VARIATION),
    expectedVariation: TestVariations.variationC,
    expectedError: null,
  },
  {
    enabled: false,
    offVariation: '',
    userID: 'user5',
    flagVariations: {},
    prerequisite: [],
    expectedReason: createReason('', Reason.Type.DEFAULT),
    expectedVariation: TestVariations.variationB,
    expectedError: null,
  },
  {
    enabled: false,
    offVariation: 'variation-E',
    userID: 'user5',
    flagVariations: {},
    prerequisite: [],
    expectedReason: createReason('', Reason.Type.OFF_VARIATION),
    expectedVariation: null,
    expectedError: new Error('evaluator: variation not found'),
  },
  {
    enabled: true,
    offVariation: '',
    userID: 'user4',
    flagVariations: {},
    prerequisite: [],
    expectedReason: createReason('', Reason.Type.DEFAULT),
    expectedVariation: TestVariations.variationB,
    expectedError: null,
  },
  {
    enabled: true,
    offVariation: 'variation-C',
    userID: 'user4',
    flagVariations: {},
    prerequisite: [],
    expectedReason: createReason('', Reason.Type.DEFAULT),
    expectedVariation: TestVariations.variationB,
    expectedError: null,
  },
  {
    enabled: true,
    offVariation: 'variation-C',
    userID: 'user4',
    flagVariations: {
      'test-feature2': 'variation A', // not matched with expected prerequisites variations
    },
    prerequisite: [createPrerequisite('test-feature2', 'variation D')],
    expectedReason: createReason('', Reason.Type.PREREQUISITE),
    expectedVariation: TestVariations.variationC,
    expectedError: null,
  },
  {
    enabled: true,
    offVariation: 'variation-C',
    userID: 'user4',
    flagVariations: {
      'test-feature2': 'variation D', // matched with expected prerequisites variations
    },
    prerequisite: [createPrerequisite('test-feature2', 'variation D')],
    expectedReason: createReason('', Reason.Type.DEFAULT),
    expectedVariation: TestVariations.variationB,
    expectedError: null,
  },
  {
    enabled: true,
    offVariation: 'variation-C',
    userID: 'user4',
    flagVariations: {}, // not found prerequisite variation
    prerequisite: [createPrerequisite('test-feature2', 'variation D')],
    expectedReason: createReason('', Reason.Type.PREREQUISITE),
    expectedVariation: null,
    expectedError: new Error('evaluator: prerequisite variation not found'),
  },
];

TestCases.forEach((tc, index) => {
  test(`Test Case ${index}`, async (t) => {
    const user = createUser(tc.userID, {});
    const f = newTestFeature('test-feature');
    f.setEnabled(tc.enabled);
    f.setOffVariation(tc.offVariation);
    f.setPrerequisitesList(tc.prerequisite);

    const evaluator = new Evaluator();
    try {
      const [reason, variation] = evaluator.assignUser(f, user, [], tc.flagVariations);
      //TODO: check with deep equal ?
      t.deepEqual(reason, tc.expectedReason);
      t.deepEqual(variation, tc.expectedVariation);
    } catch (err) {
      t.deepEqual(err, tc.expectedError);
    }
  });
});
