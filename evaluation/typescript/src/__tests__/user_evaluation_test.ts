import test from 'ava';
import * as userEvaluation from '../userEvaluation';
import { Evaluation } from '../proto/feature/evaluation_pb';
import {
  FeaturesIDFlagsPart,
  GenerateFeaturesID,
  UserEvaluationsID,
  sortMapKeys,
} from '../userEvaluation';
import { Feature } from '../proto/feature/feature_pb';
import type { HoldoutFeature, HoldoutGroup } from '../bucketeer';
import { createHoldoutFeature, createHoldoutGroup } from '../modelFactory';

// Helper function to create an Evaluation object
function NewEvaluation(id: string): Evaluation {
  const evaluation = new Evaluation();
  evaluation.setId(id);
  return evaluation;
}

// Helper function to create a Feature object
function NewFeature(id: string, updatedAt: number): Feature {
  const feature = new Feature();
  feature.setId(id);
  feature.setUpdatedAt(updatedAt);
  return feature;
}

// Define the interface for the test case
type NewUserEvaluationsTestCase = {
  id: string;
  evaluations: Evaluation[];
  archivedFeaturesIds: string[];
  forceUpdate: boolean;
  expected: {
    id: string;
    evaluations: Evaluation[];
    archivedFeaturesIds: string[];
    forceUpdate: boolean;
  };
};

// Test cases pattern similar to the Go version
const NewUserEvaluationsTestCases: NewUserEvaluationsTestCase[] = [
  {
    id: '1234',
    evaluations: [NewEvaluation('test-id1')],
    archivedFeaturesIds: ['test-id2'],
    forceUpdate: false,
    expected: {
      id: '1234',
      evaluations: [NewEvaluation('test-id1')],
      archivedFeaturesIds: ['test-id2'],
      forceUpdate: false,
    },
  },
  {
    id: '5678',
    evaluations: [NewEvaluation('test-id3')],
    archivedFeaturesIds: [],
    forceUpdate: true,
    expected: {
      id: '5678',
      evaluations: [NewEvaluation('test-id3')],
      archivedFeaturesIds: [],
      forceUpdate: true,
    },
  },
];

// Iterate over each test case
NewUserEvaluationsTestCases.forEach(
  ({ id, evaluations, archivedFeaturesIds, forceUpdate, expected }) => {
    test(`NewUserEvaluations - ${id}`, (t) => {
      const actual = userEvaluation.NewUserEvaluations(
        id,
        evaluations,
        archivedFeaturesIds,
        forceUpdate,
      );

      t.is(actual.getId(), expected.id);
      t.deepEqual(actual.getEvaluationsList(), expected.evaluations);
      t.deepEqual(actual.getArchivedFeatureIdsList(), expected.archivedFeaturesIds);
      t.is(actual.getForceUpdate(), expected.forceUpdate);
      t.truthy(actual.getCreatedAt()); // Check if CreatedAt is set, similar to NotZero in Go
      //TODO: Check me - did Golang test is missing this
      //t.deepEqual(actual.toObject(), expected);
    });
  },
);

// Test cases similar to the Go version
type SortMapKeysTestCase = {
  input: { [key: string]: string } | null;
  expected: string[];
  desc: string;
};

const SortMapTestCases: SortMapKeysTestCase[] = [
  {
    input: null,
    expected: [],
    desc: 'nil',
  },
  {
    input: {},
    expected: [],
    desc: 'empty',
  },
  {
    input: { b: 'value-b', c: 'value-c', a: 'value-a', d: 'value-d' },
    expected: ['a', 'b', 'c', 'd'],
    desc: 'success',
  },
];

// Run each test case
SortMapTestCases.forEach(({ input, expected, desc }) => {
  test(`sortMapKeys - ${desc}`, (t) => {
    const actual = sortMapKeys(input ?? {});
    t.deepEqual(actual, expected, desc);
  });
});

// Define test case structure for UserEvaluationsID
// Note: These test values must match the Go implementation exactly
type UserEvaluationsIDTestCase = {
  desc: string;
  userID: string;
  userMetadata: Record<string, string>;
  features: Feature[] | null;
  expected: string;
};

const UserEvaluationsIDTestCases: UserEvaluationsIDTestCase[] = [
  {
    desc: 'empty user ID, empty metadata, empty features',
    userID: '',
    userMetadata: {},
    features: null,
    expected: '14695981039346656037',
  },
  {
    desc: 'user ID only',
    userID: 'user-1',
    userMetadata: {},
    features: null,
    expected: '17891572797655370708',
  },
  {
    desc: 'user ID with metadata',
    userID: 'user-1',
    userMetadata: { age: '25', country: 'jp' },
    features: null,
    expected: '15857499200645826216',
  },
  {
    desc: 'user ID with metadata and single feature',
    userID: 'user-1',
    userMetadata: { age: '25', country: 'jp' },
    features: [NewFeature('feature-1', 1000)],
    expected: '10450974209164395423',
  },
  {
    desc: 'user ID with metadata and multiple features',
    userID: 'user-1',
    userMetadata: { age: '25', country: 'jp' },
    features: [NewFeature('feature-1', 1000), NewFeature('feature-2', 2000)],
    expected: '7257619227440290900',
  },
];

// Run each UserEvaluationsID test case
UserEvaluationsIDTestCases.forEach(({ desc, userID, userMetadata, features, expected }) => {
  test(`UserEvaluationsID - ${desc}`, (t) => {
    const actual = UserEvaluationsID(userID, userMetadata, features ?? []);
    t.is(actual, expected, desc);
  });
});

function holdoutGroup(id: string, updatedAt: number): HoldoutGroup {
  return {
    getId: () => id,
    getEnvironmentId: () => '',
    getAssignmentType: () => 0,
    getSamplingSeed: () => '',
    getBucketStart: () => 0,
    getPercentage: () => 0,
    getSegmentId: () => '',
    getUpdatedAt: () => updatedAt,
    getFeaturesList: () => [],
  };
}

test('UserEvaluationsID includes the live holdout group', (t) => {
  const features = [NewFeature('feature-1', 1000)];
  const without = UserEvaluationsID('user-1', {}, features);
  const groupA = holdoutGroup('hg-a', 50);
  const groupB = holdoutGroup('hg-b', 60);
  const withA = UserEvaluationsID('user-1', {}, [NewFeature('feature-1', 1000)], [groupA]);
  t.is(withA, '13300710561474652895');
  t.not(withA, without);
  t.is(
    UserEvaluationsID('user-1', {}, [NewFeature('feature-1', 1000)], [groupA, groupB]),
    UserEvaluationsID('user-1', {}, [NewFeature('feature-1', 1000)], [groupB, groupA]),
  );
});

test('UserEvaluationsID changes when targets change in the same second', (t) => {
  const features = [NewFeature('feature-1', 1000)];
  const group = (
    targets: Array<{ featureId: string; variationId: string; removedAt?: number }>,
  ): HoldoutGroup => ({
    ...holdoutGroup('hg-a', 50),
    getFeaturesList: () =>
      targets.map((target) => ({
        getId: () => '',
        getFeatureId: () => target.featureId,
        getHoldoutVariationId: () => target.variationId,
        getAddedAt: () => 0,
        getRemovedAt: () => target.removedAt ?? 0,
      })),
  });
  const withTargets = UserEvaluationsID('user-1', {}, features, [
    group([
      { featureId: 'flag-b', variationId: 'hold-b' },
      { featureId: 'flag-a', variationId: 'hold-a' },
    ]),
  ]);
  t.not(withTargets, UserEvaluationsID('user-1', {}, features, [holdoutGroup('hg-a', 50)]));
  t.not(
    withTargets,
    UserEvaluationsID('user-1', {}, features, [
      group([
        { featureId: 'flag-a', variationId: 'hold-a' },
        { featureId: 'flag-b', variationId: 'other' },
      ]),
    ]),
  );
  t.is(
    withTargets,
    UserEvaluationsID('user-1', {}, features, [
      group([
        { featureId: 'flag-a', variationId: 'hold-a' },
        { featureId: 'flag-b', variationId: 'hold-b' },
      ]),
    ]),
  );
});

// Define test case structure
type GenerateFeaturesIDTestCase = {
  desc: string;
  input: Feature[] | null;
  expected: string;
};

// Define the test cases
// Note: NewFeature(id, updatedAt) - the second parameter is the updatedAt timestamp
const GenerateFeaturesIDTestCases: GenerateFeaturesIDTestCase[] = [
  {
    desc: 'nil',
    input: null,
    expected: '14695981039346656037',
  },
  {
    desc: 'success: single',
    input: [NewFeature('id-1', 1)], // id-1 with updatedAt=1
    expected: '5476413260388599211',
  },
  {
    desc: 'success: multiple',
    input: [NewFeature('id-1', 1), NewFeature('id-2', 2)], // multiple features with updatedAt values
    expected: '17283374094628184689',
  },
];

// Run each test case
GenerateFeaturesIDTestCases.forEach(({ desc, input, expected }) => {
  test(`GenerateFeaturesID - ${desc}`, (t) => {
    const actual = GenerateFeaturesID(input ?? []);
    t.is(actual, expected, desc);
  });
});

// The pinned values match TestGenerateFeaturesIDHoldout in evaluation/go.
test('GenerateFeaturesID includes the live holdout group', (t) => {
  const newFeatures = () => [NewFeature('id-2', 2), NewFeature('id-1', 1)];
  const group = createHoldoutGroup({ id: 'hg-a', updatedAt: 50 });
  const withTarget = createHoldoutGroup({
    id: 'hg-a',
    updatedAt: 50,
    features: [
      createHoldoutFeature({
        id: 'row-2',
        featureId: 'id-2',
        holdoutVariationId: 'v-off',
        addedAt: 20,
      }),
      createHoldoutFeature({
        id: 'row-1',
        featureId: 'id-1',
        holdoutVariationId: 'v-on',
        addedAt: 10,
      }),
    ],
  });

  // With a live group the ID is "<flags>:<holdout>".
  const withoutGroup = GenerateFeaturesID(newFeatures());
  t.is(withoutGroup, '17283374094628184689');
  t.is(GenerateFeaturesID(newFeatures(), [group]), '17283374094628184689:16147347951762488293');
  t.is(GenerateFeaturesID(newFeatures(), [withTarget]), '17283374094628184689:6966737048624406302');
  t.is(GenerateFeaturesID([], [group]), '14695981039346656037:16147347951762488293');

  // The flags part does not depend on the groups.
  t.is(FeaturesIDFlagsPart(withoutGroup), withoutGroup);
  t.is(FeaturesIDFlagsPart(GenerateFeaturesID(newFeatures(), [group])), withoutGroup);
  t.is(FeaturesIDFlagsPart(GenerateFeaturesID(newFeatures(), [withTarget])), withoutGroup);
  t.is(FeaturesIDFlagsPart(''), '');

  // A group with no id is a clearance or an unavailable snapshot, not a live group.
  t.is(GenerateFeaturesID(newFeatures(), []), withoutGroup);
  t.is(GenerateFeaturesID(newFeatures(), [createHoldoutGroup({ updatedAt: 100 })]), withoutGroup);
  t.is(GenerateFeaturesID(newFeatures(), [createHoldoutGroup({ updatedAt: 1 })]), withoutGroup);

  // Restarting the group or starting another one changes the id.
  t.not(
    GenerateFeaturesID(newFeatures(), [createHoldoutGroup({ id: 'hg-a', updatedAt: 51 })]),
    GenerateFeaturesID(newFeatures(), [group]),
  );
  t.not(
    GenerateFeaturesID(newFeatures(), [createHoldoutGroup({ id: 'hg-b', updatedAt: 50 })]),
    GenerateFeaturesID(newFeatures(), [group]),
  );
});

test('GenerateFeaturesID changes when targets change in the same second', (t) => {
  const newFeatures = () => [NewFeature('id-1', 1), NewFeature('id-2', 2)];
  const row1 = createHoldoutFeature({
    id: 'row-1',
    featureId: 'id-1',
    holdoutVariationId: 'v-on',
    addedAt: 10,
  });
  const row2 = createHoldoutFeature({
    id: 'row-2',
    featureId: 'id-2',
    holdoutVariationId: 'v-off',
    addedAt: 20,
  });
  const group = (features: HoldoutFeature[]) =>
    createHoldoutGroup({ id: 'hg-a', updatedAt: 50, features });

  const base = GenerateFeaturesID(newFeatures(), [group([row2, row1])]);
  t.is(base, '17283374094628184689:6966737048624406302');
  t.is(GenerateFeaturesID(newFeatures(), [group([row1, row2])]), base);
  t.not(GenerateFeaturesID(newFeatures(), [group([])]), base);
  t.not(
    GenerateFeaturesID(newFeatures(), [
      group([
        row1,
        createHoldoutFeature({
          id: 'row-2',
          featureId: 'id-2',
          holdoutVariationId: 'v-on',
          addedAt: 20,
        }),
      ]),
    ]),
    base,
  );

  // A removed target no longer counts.
  const removed = createHoldoutFeature({
    id: 'row-2',
    featureId: 'id-2',
    holdoutVariationId: 'v-off',
    addedAt: 20,
    removedAt: 30,
  });
  t.not(GenerateFeaturesID(newFeatures(), [group([row1, removed])]), base);
  t.is(
    GenerateFeaturesID(newFeatures(), [group([row1, removed])]),
    GenerateFeaturesID(newFeatures(), [group([row1])]),
  );
});
