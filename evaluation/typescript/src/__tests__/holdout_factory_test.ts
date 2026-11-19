import test from 'ava';
import {
  HoldoutAssignment,
  createHoldoutFeature,
  createHoldoutGroup,
  createSegmentUser,
  createUser,
  holdoutSegmentIds,
  isHeldOut,
  SegmentUser,
} from '../index';

test('createHoldoutFeature defaults to a current target with no ids', (t) => {
  const feature = createHoldoutFeature();
  t.is(feature.getId(), '');
  t.is(feature.getFeatureId(), '');
  t.is(feature.getHoldoutVariationId(), '');
  t.is(feature.getAddedAt(), 0);
  t.is(feature.getRemovedAt(), 0);
});

test('createHoldoutFeature keeps every field', (t) => {
  const feature = createHoldoutFeature({
    id: 'row-1',
    featureId: 'flag-1',
    holdoutVariationId: 'variation-off',
    addedAt: 10,
    removedAt: 20,
  });
  t.is(feature.getId(), 'row-1');
  t.is(feature.getFeatureId(), 'flag-1');
  t.is(feature.getHoldoutVariationId(), 'variation-off');
  t.is(feature.getAddedAt(), 10);
  t.is(feature.getRemovedAt(), 20);
});

test('createHoldoutGroup defaults to a HASH group that holds nobody out', (t) => {
  const group = createHoldoutGroup();
  t.is(group.getId(), '');
  t.is(group.getEnvironmentId(), '');
  t.is(group.getAssignmentType(), HoldoutAssignment.HASH);
  t.is(group.getSamplingSeed(), '');
  t.is(group.getBucketStart(), 0);
  t.is(group.getPercentage(), 0);
  t.is(group.getSegmentId(), '');
  t.is(group.getUpdatedAt(), 0);
  t.deepEqual(group.getFeaturesList(), []);
  t.false(isHeldOut(group, createUser('user-1', null), []));
  t.deepEqual(holdoutSegmentIds([group]), []);
});

test('createHoldoutGroup builds a HASH group the evaluator can bucket', (t) => {
  const features = [createHoldoutFeature({ featureId: 'flag-1', holdoutVariationId: 'off' })];
  const everyone = createHoldoutGroup({
    id: 'hg-1',
    environmentId: 'env-1',
    assignmentType: HoldoutAssignment.HASH,
    samplingSeed: 'seed',
    bucketStart: 25,
    percentage: 100,
    updatedAt: 50,
    features,
  });
  t.is(everyone.getId(), 'hg-1');
  t.is(everyone.getEnvironmentId(), 'env-1');
  t.is(everyone.getSamplingSeed(), 'seed');
  t.is(everyone.getBucketStart(), 25);
  t.is(everyone.getPercentage(), 100);
  t.is(everyone.getUpdatedAt(), 50);
  t.is(everyone.getFeaturesList(), features);
  // A full arc holds out every user, wherever it starts.
  t.true(isHeldOut(everyone, createUser('user-1', null), []));
  t.true(isHeldOut(everyone, createUser('user-2', null), []));
  t.deepEqual(holdoutSegmentIds([everyone]), []);
});

test('createHoldoutGroup builds a SEGMENT group whose segment the SDK must fetch', (t) => {
  const group = createHoldoutGroup({
    id: 'hg-1',
    assignmentType: HoldoutAssignment.SEGMENT,
    segmentId: 'segment-1',
  });
  t.deepEqual(holdoutSegmentIds([group]), ['segment-1']);
  t.deepEqual(holdoutSegmentIds([group, group]), ['segment-1']);

  const included = [createSegmentUser('user-1', 'segment-1', SegmentUser.State.INCLUDED)];
  t.true(isHeldOut(group, createUser('user-1', null), included));
  t.false(isHeldOut(group, createUser('user-2', null), included));
  t.false(isHeldOut(group, createUser('user-1', null), []));
});
