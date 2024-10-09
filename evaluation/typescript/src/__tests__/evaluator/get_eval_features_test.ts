import test from 'ava';
import { createPrerequisite, creatFeature } from '../../modelFactory';
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
      creatFeature({
        id: 'featureA',
      }),
    ],
    all: [
      creatFeature({
        id: 'featureA',
      }),
      creatFeature({
        id: 'featureB',
      }),
    ],
    expectedIDs: ['featureA'],
  },
  {
    desc: 'success: one feature depends on target',
    targets: [
      creatFeature({
        id: 'featureA',
      }),
    ],
    all: [
      creatFeature({
        id: 'featureA',
      }),
      creatFeature({
        id: 'featureB',
        prerequisitesList: [createPrerequisite('featureA', '')],
      }),
      creatFeature({
        id: 'featureC',
      }),
    ],
    expectedIDs: ['featureA', 'featureB'],
  },
  {
    desc: 'success: multiple features depends on target',
    targets: [
      creatFeature({
        id: 'featureA',
      }),
    ],
    all: [
      creatFeature({
        id: 'featureA',
      }),
      creatFeature({
        id: 'featureB',
        prerequisitesList: [createPrerequisite('featureA', '')],
      }),
      creatFeature({
        id: 'featureC',
        prerequisitesList: [createPrerequisite('featureB', '')],
      }),
      creatFeature({
        id: 'featureD',
        prerequisitesList: [createPrerequisite('featureA', '')],
      }),
      creatFeature({
        id: 'featureE',
      }),
    ],
    expectedIDs: ['featureA', 'featureB', 'featureC', 'featureD'],
  },
  {
    desc: 'success: target depends on one feature',
    targets: [
      creatFeature({
        id: 'featureA',
        prerequisitesList: [createPrerequisite('featureB', '')],
      }),
    ],
    all: [
      creatFeature({
        id: 'featureA',
        prerequisitesList: [createPrerequisite('featureB', '')],
      }),
      creatFeature({
        id: 'featureB',
      }),
      creatFeature({
        id: 'featureC',
      }),
    ],
    expectedIDs: ['featureA', 'featureB'],
  },
  {
    desc: 'success: target depends on multiple features',
    targets: [
      creatFeature({
        id: 'featureA',
        prerequisitesList: [createPrerequisite('featureB', ''), createPrerequisite('featureC', '')],
      }),
    ],
    all: [
      creatFeature({
        id: 'featureA',
        prerequisitesList: [createPrerequisite('featureB', '')],
      }),
      creatFeature({
        id: 'featureB',
        prerequisitesList: [createPrerequisite('featureD', '')],
      }),
      creatFeature({
        id: 'featureC',
      }),
      creatFeature({
        id: 'featureD',
      }),
      creatFeature({
        id: 'featureE',
      }),
    ],
    // order is different with golang test but the result is same
    expectedIDs: ['featureA', 'featureB', 'featureD', 'featureC'],
  },
  {
    desc: 'success: complex pattern 1',
    targets: [
      creatFeature({
        id: 'featureD',
        prerequisitesList: [createPrerequisite('featureB', '')],
      }),
    ],
    all: [
      creatFeature({
        id: 'featureA',
      }),
      creatFeature({
        id: 'featureB',
        prerequisitesList: [createPrerequisite('featureA', '')],
      }),
      creatFeature({
        id: 'featureC',
        prerequisitesList: [createPrerequisite('featureB', '')],
      }),
      creatFeature({
        id: 'featureD',
        prerequisitesList: [createPrerequisite('featureB', '')],
      }),
      creatFeature({
        id: 'featureE',
        prerequisitesList: [createPrerequisite('featureC', ''), createPrerequisite('featureD', '')],
      }),
      creatFeature({
        id: 'featureF',
        prerequisitesList: [createPrerequisite('featureE', '')],
      }),
      creatFeature({
        id: 'featureG',
        prerequisitesList: [createPrerequisite('featureA', '')],
      }),
      creatFeature({
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
