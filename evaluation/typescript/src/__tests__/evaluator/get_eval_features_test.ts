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
    expectedIDs: ['featureD', 'featureB', 'featureA', 'featureE', 'featureF'],
  },
];

patterns.forEach(({ desc, targets, all, expectedIDs }) => {
  test(desc, (t) => {
    // Test code
    try {
      const evalator = new Evaluator();
      const actual = evalator.getEvalFeatures(targets, all);
      t.deepEqual(
        actual.map((e) => {
          return e.getId();
        }),
        expectedIDs,
      );
    } catch (error) {
      t.fail(`Error: ${error}`);
    }
  });
});
