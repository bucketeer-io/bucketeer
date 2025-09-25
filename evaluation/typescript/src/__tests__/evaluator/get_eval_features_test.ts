import test from 'ava';
import { createPrerequisite, createFeature } from '../../modelFactory';
import { Feature } from '../../proto/feature/feature_pb';
import { Evaluator } from '../../evaluation';

interface TestCase {
  desc: string;
  targets: Feature[];
  all: Feature[];
  expectedIDs: string[];
}

const patterns: TestCase[] = [
  {
    desc: 'success: No prerequisites',
    targets: [
      createFeature({
        id: 'featureA',
      }),
    ],
    all: [
      createFeature({
        id: 'featureA',
      }),
      createFeature({
        id: 'featureB',
      }),
    ],
    expectedIDs: ['featureA'],
  },
  {
    desc: 'success: one feature depends on target',
    targets: [
      createFeature({
        id: 'featureA',
      }),
    ],
    all: [
      createFeature({
        id: 'featureA',
      }),
      createFeature({
        id: 'featureB',
        prerequisitesList: [createPrerequisite('featureA', '')],
      }),
      createFeature({
        id: 'featureC',
      }),
    ],
    expectedIDs: ['featureA', 'featureB'],
  },
  {
    desc: 'success: multiple features depends on target',
    targets: [
      createFeature({
        id: 'featureA',
      }),
    ],
    all: [
      createFeature({
        id: 'featureA',
      }),
      createFeature({
        id: 'featureB',
        prerequisitesList: [createPrerequisite('featureA', '')],
      }),
      createFeature({
        id: 'featureC',
        prerequisitesList: [createPrerequisite('featureB', '')],
      }),
      createFeature({
        id: 'featureD',
        prerequisitesList: [createPrerequisite('featureA', '')],
      }),
      createFeature({
        id: 'featureE',
      }),
    ],
    expectedIDs: ['featureA', 'featureB', 'featureC', 'featureD'],
  },
  {
    desc: 'success: target depends on one feature',
    targets: [
      createFeature({
        id: 'featureA',
        prerequisitesList: [createPrerequisite('featureB', '')],
      }),
    ],
    all: [
      createFeature({
        id: 'featureA',
        prerequisitesList: [createPrerequisite('featureB', '')],
      }),
      createFeature({
        id: 'featureB',
      }),
      createFeature({
        id: 'featureC',
      }),
    ],
    expectedIDs: ['featureA', 'featureB'],
  },
  {
    desc: 'success: target depends on multiple features',
    targets: [
      createFeature({
        id: 'featureA',
        prerequisitesList: [createPrerequisite('featureB', ''), createPrerequisite('featureC', '')],
      }),
    ],
    all: [
      createFeature({
        id: 'featureA',
        prerequisitesList: [createPrerequisite('featureB', '')],
      }),
      createFeature({
        id: 'featureB',
        prerequisitesList: [createPrerequisite('featureD', '')],
      }),
      createFeature({
        id: 'featureC',
      }),
      createFeature({
        id: 'featureD',
      }),
      createFeature({
        id: 'featureE',
      }),
    ],
    // order is different with golang test but the result is same
    expectedIDs: ['featureA', 'featureB', 'featureD', 'featureC'],
  },
  {
    desc: 'success: complex pattern 1',
    targets: [
      createFeature({
        id: 'featureD',
        prerequisitesList: [createPrerequisite('featureB', '')],
      }),
    ],
    all: [
      createFeature({
        id: 'featureA',
      }),
      createFeature({
        id: 'featureB',
        prerequisitesList: [createPrerequisite('featureA', '')],
      }),
      createFeature({
        id: 'featureC',
        prerequisitesList: [createPrerequisite('featureB', '')],
      }),
      createFeature({
        id: 'featureD',
        prerequisitesList: [createPrerequisite('featureB', '')],
      }),
      createFeature({
        id: 'featureE',
        prerequisitesList: [createPrerequisite('featureC', ''), createPrerequisite('featureD', '')],
      }),
      createFeature({
        id: 'featureF',
        prerequisitesList: [createPrerequisite('featureE', '')],
      }),
      createFeature({
        id: 'featureG',
        prerequisitesList: [createPrerequisite('featureA', '')],
      }),
      createFeature({
        id: 'featureH',
      }),
    ],
    // order is different with golang test but the result is same
    // After transitive closure fix, featureC should be included as it's a dependency of featureE
    expectedIDs: ['featureD', 'featureB', 'featureA', 'featureC', 'featureE', 'featureF'],
  },
];

patterns.forEach(({ desc, targets, all, expectedIDs }) => {
  test(desc, (t) => {
    // Test code
    try {
      const evalator = new Evaluator();
      const actual = evalator.getEvalFeatures(targets, all);
      const actualIDs = actual.map((e) => {
        return e.getId();
      });

      // Use set-based comparison since order doesn't matter for correctness
      t.is(
        actualIDs.length,
        expectedIDs.length,
        `Expected ${expectedIDs.length} features, got ${actualIDs.length}`,
      );

      for (const expectedID of expectedIDs) {
        t.true(actualIDs.includes(expectedID), `Missing expected feature: ${expectedID}`);
      }

      for (const actualID of actualIDs) {
        t.true(expectedIDs.includes(actualID), `Unexpected feature: ${actualID}`);
      }
    } catch (error) {
      t.fail(`Error: ${error}`);
    }
  });
});
