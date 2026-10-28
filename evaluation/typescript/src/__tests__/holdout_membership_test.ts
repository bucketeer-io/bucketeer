import test from 'ava';
import { HoldoutAssignment, HoldoutMembership } from '../bucketeer';
import type { HoldoutFeature, HoldoutGroup } from '../bucketeer';
import { SegmentUser } from '../proto/feature/segment_pb';
import { User } from '../proto/user/user_pb';

const user = new User();
user.setId('user-1');

function target(featureId: string, variationId: string, removedAt = 0): HoldoutFeature {
  return {
    getId: () => '',
    getFeatureId: () => featureId,
    getHoldoutVariationId: () => variationId,
    getAddedAt: () => 0,
    getRemovedAt: () => removedAt,
  };
}

function segmentGroup(id: string, segmentId: string, features: HoldoutFeature[]): HoldoutGroup {
  return {
    getId: () => id,
    getEnvironmentId: () => 'env',
    getAssignmentType: () => HoldoutAssignment.SEGMENT,
    getSamplingSeed: () => '',
    getBucketStart: () => 0,
    getPercentage: () => 0,
    getSegmentId: () => segmentId,
    getUpdatedAt: () => 1,
    getFeaturesList: () => features,
  };
}

function includedUser(segmentId: string): SegmentUser {
  const segmentUser = new SegmentUser();
  segmentUser.setSegmentId(segmentId);
  segmentUser.setUserId(user.getId());
  segmentUser.setState(SegmentUser.State.INCLUDED);
  return segmentUser;
}

function segmentUsers(): Map<string, SegmentUser[]> {
  return new Map<string, SegmentUser[]>([
    ['seg', [includedUser('seg')]],
    // Another flag's segment also lists the user; it must not be consulted.
    ['other', [includedUser('seg')]],
  ]);
}

test('HoldoutMembership.create returns undefined without groups', (t) => {
  t.is(HoldoutMembership.create(undefined, user, segmentUsers()), undefined);
  t.is(HoldoutMembership.create([], user, segmentUsers()), undefined);
});

test('member of the group own segment is served the holdout variation', (t) => {
  const group = segmentGroup('hg', 'seg', [target('flag', 'hold')]);
  const held = HoldoutMembership.create([group], user, segmentUsers())?.variationFor('flag');
  t.truthy(held);
  t.is(held?.group, group);
  t.is(held?.variationId, 'hold');
});

test('only the group own segment is scanned', (t) => {
  const group = segmentGroup('hg', 'seg', [target('flag', 'hold')]);
  // The user row is filed under key "other", not under the group's segment.
  const users = new Map<string, SegmentUser[]>([['other', [includedUser('seg')]]]);
  t.is(HoldoutMembership.create([group], user, users)?.variationFor('flag'), undefined);
});

test('untargeted flag does not resolve membership', (t) => {
  const group = segmentGroup('hg', 'seg', [target('flag', 'hold')]);
  // Remove the user's row after an untargeted lookup: if that lookup had
  // resolved membership, the later targeted lookup would still be a hit.
  const users = segmentUsers();
  const membership = HoldoutMembership.create([group], user, users);
  t.is(membership?.variationFor('other-flag'), undefined);
  users.delete('seg');
  t.is(membership?.variationFor('flag'), undefined);
});

test('membership is resolved once and reused across flags', (t) => {
  const group = segmentGroup('hg', 'seg', [target('flag', 'hold'), target('flag2', 'h2')]);
  // Remove the user's row after the first resolution: a second resolution
  // would observe the removal, so a hit on flag2 proves the memo.
  const users = segmentUsers();
  const membership = HoldoutMembership.create([group], user, users);
  t.is(membership?.variationFor('flag')?.variationId, 'hold');
  users.delete('seg');
  t.is(membership?.variationFor('flag2')?.variationId, 'h2');
});

test('removed target and clearance snapshot are not targets', (t) => {
  const removed = segmentGroup('hg', 'seg', [target('flag', 'hold', 1)]);
  t.is(HoldoutMembership.create([removed], user, segmentUsers())?.variationFor('flag'), undefined);
  const clearance = segmentGroup('', '', []);
  t.is(
    HoldoutMembership.create([clearance], user, segmentUsers())?.variationFor('flag'),
    undefined,
  );
});

test('a miss on the targeting group does not fall through to another group', (t) => {
  const miss = segmentGroup('miss', 'none', [target('flag', 'hold')]);
  const hit = segmentGroup('hit', 'seg', [target('flag', 'hold')]);
  t.is(
    HoldoutMembership.create([miss, hit], user, segmentUsers())?.variationFor('flag'),
    undefined,
  );
});
