import test from 'ava';
import { Evaluator } from '../evaluation';
import { HoldoutAssignment } from '../bucketeer';
import type { HoldoutFeature, HoldoutGroup } from '../bucketeer';
import { Feature } from '../proto/feature/feature_pb';
import { Reason } from '../proto/feature/reason_pb';
import { Strategy, FixedStrategy } from '../proto/feature/strategy_pb';
import { Variation } from '../proto/feature/variation_pb';
import { createUser } from '../modelFactory';

// Matches SECONDS_FOR_ADJUSTMENT in evaluation.ts.
const GRACE_SECONDS = 10;

function flag(id: string, updatedAt: number): Feature {
  const feature = new Feature();
  feature.setId(id);
  feature.setEnabled(true);
  feature.setOffVariation('off');
  feature.setUpdatedAt(updatedAt);
  feature.setVariationsList(
    ['off', 'hold', 'on'].map((variationId) => {
      const variation = new Variation();
      variation.setId(variationId);
      variation.setValue(variationId);
      return variation;
    }),
  );
  const strategy = new Strategy();
  strategy.setType(Strategy.Type.FIXED);
  const fixed = new FixedStrategy();
  fixed.setVariation('on');
  strategy.setFixedStrategy(fixed);
  feature.setDefaultStrategy(strategy);
  return feature;
}

function group(id: string, updatedAt: number, targets: boolean): HoldoutGroup {
  const features: HoldoutFeature[] = targets
    ? [
        {
          getId: () => '',
          getFeatureId: () => 'stale',
          getHoldoutVariationId: () => 'hold',
          getAddedAt: () => 0,
          getRemovedAt: () => 0,
        },
      ]
    : [];
  return {
    getId: () => id,
    getEnvironmentId: () => 'env',
    getAssignmentType: () => HoldoutAssignment.HASH,
    getSamplingSeed: () => 'seed',
    getBucketStart: () => 0,
    getPercentage: () => 100,
    getSegmentId: () => '',
    getUpdatedAt: () => updatedAt,
    getFeaturesList: () => features,
  };
}

test('holdout cursor uses the feature-diff grace window and a stop revision', (t) => {
  const now = Math.floor(Date.now() / 1000);
  const evaluatedAt = now - 100;
  const features = [flag('stale', now - 1000), flag('fresh', now - 10)];
  const user = createUser('user-1', null);
  const evaluator = new Evaluator();

  const cases: Array<{
    name: string;
    holdout: HoldoutGroup;
    force: boolean;
    ids: string[];
    staleReason?: Reason.TypeMap[keyof Reason.TypeMap];
    staleVariation?: string;
  }> = [
    {
      name: 'holdout newer than the cursor',
      holdout: group('hg', now - 50, true),
      force: true,
      ids: ['stale', 'fresh'],
      staleReason: Reason.Type.HOLDOUT,
      staleVariation: 'hold',
    },
    {
      name: 'same second as evaluatedAt',
      holdout: group('hg', evaluatedAt, true),
      force: true,
      ids: ['stale', 'fresh'],
      staleReason: Reason.Type.HOLDOUT,
      staleVariation: 'hold',
    },
    {
      name: 'update inside the grace window',
      holdout: group('hg', evaluatedAt - (GRACE_SECONDS - 1), true),
      force: true,
      ids: ['stale', 'fresh'],
      staleReason: Reason.Type.HOLDOUT,
      staleVariation: 'hold',
    },
    {
      name: 'update older than the grace window',
      holdout: group('hg', evaluatedAt - GRACE_SECONDS, true),
      force: false,
      ids: ['fresh'],
    },
    {
      name: 'serving stopped',
      holdout: group('', now - 50, false),
      force: true,
      ids: ['stale', 'fresh'],
      staleReason: Reason.Type.DEFAULT,
      staleVariation: 'on',
    },
    {
      name: 'environment never stopped a group',
      holdout: group('', 0, false),
      force: false,
      ids: ['fresh'],
    },
  ];

  cases.forEach((tc) => {
    const actual = evaluator.evaluateFeaturesByEvaluatedAt(
      [...features],
      user,
      new Map(),
      null,
      'previous',
      evaluatedAt,
      false,
      '',
      [tc.holdout],
    );
    t.is(actual.getForceUpdate(), tc.force, tc.name);
    const evaluations = actual.getEvaluationsList();
    t.deepEqual(
      evaluations.map((evaluation) => evaluation.getFeatureId()),
      tc.ids,
      tc.name,
    );
    const stale = evaluations.find((evaluation) => evaluation.getFeatureId() === 'stale');
    if (tc.staleReason === undefined) {
      t.is(stale, undefined, tc.name);
      return;
    }
    t.is(stale?.getReason()?.getType(), tc.staleReason, tc.name);
    t.is(stale?.getVariationId(), tc.staleVariation, tc.name);
  });
});
