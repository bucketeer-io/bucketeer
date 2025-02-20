import test from 'ava';
import { DependencyEvaluator } from '../dependencyEvaluator';
import { EVALUATOR_ERRORS } from '../errors';
// Define the type for the test cases
interface TestCase {
  desc: string;
  targetFeatureID: string;
  variationIDs: string[];
  flagVariations: Record<string, string>;
  expected: boolean;
  expectedErr: Error | null;
}

const patterns: TestCase[] = [
  {
    desc: 'err: depended feature not found',
    targetFeatureID: 'feature-1',
    variationIDs: ['variation-1', 'variation-2'],
    flagVariations: {},
    expected: false,
    expectedErr: EVALUATOR_ERRORS.FeatureNotFound,
  },
  {
    desc: 'not matched',
    targetFeatureID: 'feature-1',
    variationIDs: ['variation-1', 'variation-2'],
    flagVariations: { 'feature-1': 'variation-3' },
    expected: false,
    expectedErr: null,
  },
  {
    desc: 'success',
    targetFeatureID: 'feature-1',
    variationIDs: ['variation-1', 'variation-2'],
    flagVariations: { 'feature-1': 'variation-2' },
    expected: true,
    expectedErr: null,
  },
];

patterns.forEach(
  ({ desc, targetFeatureID, variationIDs, flagVariations, expected, expectedErr }) => {
    test(desc, (t) => {
      const dEval = new DependencyEvaluator();
      try {
        const actual = dEval.evaluate(targetFeatureID, variationIDs, flagVariations);
        t.deepEqual(actual, expected);
      } catch (error) {
        t.deepEqual(error, expectedErr);
      }
    });
  },
);
