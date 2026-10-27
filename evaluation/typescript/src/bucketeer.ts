import murmurHash3 from 'murmurhash3js';
import { SegmentUser } from './proto/feature/segment_pb';
import { User } from './proto/user/user_pb';

// Matches proto HoldoutGroup.AssignmentType. The evaluator takes this shape so it
// does not import holdout.proto, which depends on google.api.field_behavior.
export const HoldoutAssignment = {
  HASH: 0,
  SEGMENT: 1,
} as const;

export interface HoldoutFeature {
  getFeatureId(): string;
  getHoldoutVariationId(): string;
  getRemovedAt(): number;
}

export interface HoldoutGroup {
  getId(): string;
  getEnvironmentId(): string;
  getAssignmentType(): number;
  getSamplingSeed(): string;
  getBucketStart(): number;
  getPercentage(): number;
  getSegmentId(): string;
  getUpdatedAt(): number;
  getFeaturesList(): HoldoutFeature[];
}

// Thousandths of a percent on [0, 100). Matches DECIMAL(6,3).
const HOLDOUT_PERCENT_SCALE = 100000;

class Bucketeer {
  // Calculate the input hash of the target property and map it to a float64 between [0,1]
  bucket(input: string): number {
    const { high, low } = this.murmur128(input);
    return this.toFloat64(high, low);
  }

  // Computes MurmurHash3 (128-bit) hash and returns the high and low parts.
  // The high and low parts of the hash are returned in big-endian format,
  // which is crucial for consistency with the Go implementation that also uses big-endian.
  private murmur128(input: string): { high: bigint; low: bigint } {
    const hash = murmurHash3.x64.hash128(input);
    // Convert the hash from hexadecimal to BigInt.
    // The first 16 characters represent the high part (most significant bits),
    // and the last 16 characters represent the low part (least significant bits).
    const high = BigInt('0x' + hash.slice(0, 16)); // High part
    const low = BigInt('0x' + hash.slice(16, 32)); // Low part

    return { high, low };
  }

  // Convert a 128-bit hash (two uint64 values) to a float64 in the range [0,1].
  // Because we bucket millions of users, we use the full 128-bit hash to avoid collisions.
  private toFloat64(high: bigint, low: bigint): number {
    // Combine the high and low parts into a single floating-point number.
    // This maintains the full 128-bit range, ensuring a normalized value between [0,1].
    const full = Number(high) * Math.pow(2, 64) + Number(low);

    // Calculate the maximum value for a 128-bit number.
    const maxValue = Math.pow(2, 128) - 1;

    // Normalize the combined value to the range [0,1].
    return full / maxValue;
  }
}

const holdoutBucketeer = new Bucketeer();

// isHeldOut matches evaluation/go IsHeldOut. HASH uses one murmur3 hash.
// SEGMENT checks included segment users and ignores segment rules.
export function isHeldOut(
  group: HoldoutGroup | null | undefined,
  user: User | null | undefined,
  segmentUsers: SegmentUser[],
): boolean {
  if (!group || !user) {
    return false;
  }
  if (group.getAssignmentType() === HoldoutAssignment.SEGMENT) {
    return containsIncludedSegmentUser(group.getSegmentId(), user.getId(), segmentUsers);
  }
  if (group.getAssignmentType() !== HoldoutAssignment.HASH) {
    return false;
  }
  if (group.getSamplingSeed() === '' || group.getPercentage() <= 0) {
    return false;
  }
  const bucket = holdoutBucketeer.bucket(
    `holdout-${group.getEnvironmentId()}-${user.getId()}-${group.getSamplingSeed()}`,
  );
  return holdoutArcContains(bucket, group.getBucketStart(), group.getPercentage());
}

export function holdoutArcContains(
  bucket: number,
  bucketStart: number,
  percentage: number,
): boolean {
  const start = Math.round(bucketStart * 1000);
  const end = Math.round((bucketStart + percentage) * 1000);
  let pos = Math.round(bucket * HOLDOUT_PERCENT_SCALE);
  if (pos >= HOLDOUT_PERCENT_SCALE) {
    pos = 0;
  }
  if (end <= HOLDOUT_PERCENT_SCALE) {
    return pos >= start && pos < end;
  }
  return pos >= start || pos < end - HOLDOUT_PERCENT_SCALE;
}

export function holdoutSegmentIds(groups: HoldoutGroup[] | undefined): string[] {
  if (!groups || groups.length === 0) {
    return [];
  }
  const ids: string[] = [];
  for (const group of groups) {
    if (
      !group ||
      group.getAssignmentType() !== HoldoutAssignment.SEGMENT ||
      group.getSegmentId() === ''
    ) {
      continue;
    }
    if (!ids.includes(group.getSegmentId())) {
      ids.push(group.getSegmentId());
    }
  }
  return ids;
}

// HoldoutMembership resolves membership lazily and at most once per group for
// one evaluation. Membership is flag-independent (RFC 0048 §4.2.1), so the first
// flag a group targets pays one hash (HASH) or one scan of the group's own
// segment (SEGMENT); later flags reuse the answer, and an evaluation that
// touches no targeted flag never hashes at all. Matches evaluation/go.
export class HoldoutMembership {
  // undefined: unresolved; indexed like groups.
  private readonly state: (boolean | undefined)[];

  constructor(
    private readonly groups: HoldoutGroup[],
    private readonly user: User,
    private readonly mapSegmentUsers: Map<string, SegmentUser[]>,
  ) {
    this.state = new Array<boolean | undefined>(groups.length);
  }

  static create(
    groups: HoldoutGroup[] | undefined,
    user: User,
    mapSegmentUsers: Map<string, SegmentUser[]>,
  ): HoldoutMembership | undefined {
    if (!groups || groups.length === 0) {
      return undefined;
    }
    return new HoldoutMembership(groups, user, mapSegmentUsers);
  }

  // variationFor returns the group and holdout variation to serve for featureId,
  // or undefined when no group the user belongs to targets it.
  variationFor(featureId: string): { group: HoldoutGroup; variationId: string } | undefined {
    for (let i = 0; i < this.groups.length; i++) {
      const group = this.groups[i];
      if (!group) {
        continue;
      }
      const variationId = holdoutTargetVariation(group, featureId);
      if (variationId === undefined) {
        continue;
      }
      if (this.isMember(i)) {
        return { group, variationId };
      }
      // At most one serving group targets a flag. A miss stays on the normal path.
      break;
    }
    return undefined;
  }

  // targets reports whether any group targets featureId, including when the
  // user is not a member.
  targets(featureId: string): boolean {
    return this.groups.some(
      (group) => group !== undefined && holdoutTargetVariation(group, featureId) !== undefined,
    );
  }

  private isMember(i: number): boolean {
    let member = this.state[i];
    if (member === undefined) {
      const group = this.groups[i];
      if (!group) {
        return false;
      }
      const segmentUsers =
        group.getAssignmentType() === HoldoutAssignment.SEGMENT
          ? this.mapSegmentUsers.get(group.getSegmentId()) || []
          : [];
      member = isHeldOut(group, this.user, segmentUsers);
      this.state[i] = member;
    }
    return member;
  }
}

export function holdoutTargetVariation(group: HoldoutGroup, featureId: string): string | undefined {
  for (const feature of group.getFeaturesList()) {
    if (currentHoldoutTarget(feature, featureId)) {
      return feature.getHoldoutVariationId();
    }
  }
  return undefined;
}

function currentHoldoutTarget(feature: HoldoutFeature, featureId: string): boolean {
  return feature.getRemovedAt() === 0 && feature.getFeatureId() === featureId;
}

export function holdoutRequiresFullEval(
  groups: HoldoutGroup[] | undefined,
  evaluatedAt: number,
): boolean {
  if (!groups) {
    return false;
  }
  // UpdatedAt 0 means the environment has never stopped a group.
  // An empty id with a timestamp is the snapshot written when serving stopped.
  return groups.some((group) => group.getUpdatedAt() !== 0 && evaluatedAt < group.getUpdatedAt());
}

function containsIncludedSegmentUser(
  segmentId: string,
  userId: string,
  segmentUsers: SegmentUser[],
): boolean {
  return segmentUsers.some(
    (segmentUser) =>
      segmentUser.getSegmentId() === segmentId &&
      segmentUser.getUserId() === userId &&
      segmentUser.getState() === SegmentUser.State.INCLUDED,
  );
}

export { Bucketeer };
