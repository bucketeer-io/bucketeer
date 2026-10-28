import test from 'ava';
import * as fs from 'fs';
import * as path from 'path';
import { Evaluator } from '../evaluation';
import { HoldoutAssignment } from '../bucketeer';
import type { HoldoutFeature, HoldoutGroup } from '../bucketeer';
import { Feature } from '../proto/feature/feature_pb';
import { Prerequisite } from '../proto/feature/prerequisite_pb';
import { Reason } from '../proto/feature/reason_pb';
import { SegmentUser } from '../proto/feature/segment_pb';
import { Strategy, FixedStrategy } from '../proto/feature/strategy_pb';
import { Variation } from '../proto/feature/variation_pb';
import { createUser } from '../modelFactory';

// Shared with evaluation/go. See evaluation/testdata/holdout_conformance.json.
interface ConformanceTarget {
  featureId: string;
  variationId: string;
}

interface ConformanceGroup {
  id: string;
  assignmentType: string;
  bucketStart?: number;
  percentage?: number;
  segmentId?: string;
  targets: ConformanceTarget[];
}

interface ConformanceFeature {
  id: string;
  enabled: boolean;
  offVariation: string;
  variations: string[];
  defaultVariation: string;
  prerequisites?: Array<{ featureId: string; variationId: string }>;
}

interface ConformanceExpected {
  featureId: string;
  variationId: string;
  reason: string;
  holdoutGroupId: string;
}

interface ConformanceCase {
  desc: string;
  userId: string;
  group: ConformanceGroup;
  features: ConformanceFeature[];
  expected: ConformanceExpected[];
}

interface ConformanceFixture {
  environmentId: string;
  seed: string;
  segmentUsers: Array<{ segmentId: string; userId: string; state: string }>;
  cases: ConformanceCase[];
}

function loadFixture(): ConformanceFixture {
  const candidates = [
    path.join(__dirname, '../../../testdata/holdout_conformance.json'),
    path.join(__dirname, '../../testdata/holdout_conformance.json'),
  ];
  for (const candidate of candidates) {
    if (fs.existsSync(candidate)) {
      return JSON.parse(fs.readFileSync(candidate, 'utf-8'));
    }
  }
  throw new Error('holdout_conformance.json not found');
}

function toState(name: string): SegmentUser.StateMap[keyof SegmentUser.StateMap] {
  const state = SegmentUser.State[name as keyof SegmentUser.StateMap];
  if (state === undefined) {
    throw new Error(`unknown segment user state: ${name}`);
  }
  return state;
}

function toAssignment(name: string): number {
  if (name === 'HASH') {
    return HoldoutAssignment.HASH;
  }
  if (name === 'SEGMENT') {
    return HoldoutAssignment.SEGMENT;
  }
  throw new Error(`unknown assignment type: ${name}`);
}

function toReason(name: string): Reason.TypeMap[keyof Reason.TypeMap] {
  const reason = Reason.Type[name as keyof Reason.TypeMap];
  if (reason === undefined) {
    throw new Error(`unknown reason: ${name}`);
  }
  return reason;
}

const fixture = loadFixture();

const segmentUsers = new Map<string, SegmentUser[]>();
fixture.segmentUsers.forEach((su) => {
  const segmentUser = new SegmentUser();
  segmentUser.setSegmentId(su.segmentId);
  segmentUser.setUserId(su.userId);
  segmentUser.setState(toState(su.state));
  const users = segmentUsers.get(su.segmentId) ?? [];
  users.push(segmentUser);
  segmentUsers.set(su.segmentId, users);
});

function buildGroup(group: ConformanceGroup): HoldoutGroup {
  const features: HoldoutFeature[] = group.targets.map((target) => ({
    getId: () => '',
    getFeatureId: () => target.featureId,
    getHoldoutVariationId: () => target.variationId,
    getAddedAt: () => 0,
    getRemovedAt: () => 0,
  }));
  return {
    getId: () => group.id,
    getEnvironmentId: () => fixture.environmentId,
    getAssignmentType: () => toAssignment(group.assignmentType),
    getSamplingSeed: () => fixture.seed,
    getBucketStart: () => group.bucketStart ?? 0,
    getPercentage: () => group.percentage ?? 0,
    getSegmentId: () => group.segmentId ?? '',
    getUpdatedAt: () => 0,
    getFeaturesList: () => features,
  };
}

function buildFeatures(features: ConformanceFeature[]): Feature[] {
  return features.map((spec) => {
    const feature = new Feature();
    feature.setId(spec.id);
    feature.setEnabled(spec.enabled);
    feature.setOffVariation(spec.offVariation);
    feature.setVariationsList(
      spec.variations.map((id) => {
        const variation = new Variation();
        variation.setId(id);
        variation.setValue(id);
        return variation;
      }),
    );
    feature.setPrerequisitesList(
      (spec.prerequisites ?? []).map((prerequisite) => {
        const proto = new Prerequisite();
        proto.setFeatureId(prerequisite.featureId);
        proto.setVariationId(prerequisite.variationId);
        return proto;
      }),
    );
    const strategy = new Strategy();
    strategy.setType(Strategy.Type.FIXED);
    const fixed = new FixedStrategy();
    fixed.setVariation(spec.defaultVariation);
    strategy.setFixedStrategy(fixed);
    feature.setDefaultStrategy(strategy);
    return feature;
  });
}

fixture.cases.forEach((tc) => {
  test(`holdout conformance: ${tc.desc}`, async (t) => {
    const evaluator = new Evaluator();
    const evaluations = await evaluator.evaluateFeatures(
      buildFeatures(tc.features),
      createUser(tc.userId, null),
      segmentUsers,
      null,
      '',
      [buildGroup(tc.group)],
    );
    const actual = new Map(
      evaluations.getEvaluationsList().map((evaluation) => [evaluation.getFeatureId(), evaluation]),
    );
    t.is(actual.size, tc.expected.length);
    tc.expected.forEach((expected) => {
      const evaluation = actual.get(expected.featureId);
      t.truthy(evaluation, expected.featureId);
      t.is(evaluation?.getVariationId(), expected.variationId);
      t.is(evaluation?.getReason()?.getType(), toReason(expected.reason));
      t.is(evaluation?.getReason()?.getHoldoutGroupId(), expected.holdoutGroupId);
    });
  });
});
