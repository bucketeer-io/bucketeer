import test from 'ava';
import { Evaluator } from '../evaluation';
import { HoldoutAssignment } from '../bucketeer';
import type { HoldoutGroup } from '../bucketeer';
import { Feature } from '../proto/feature/feature_pb';
import { Prerequisite } from '../proto/feature/prerequisite_pb';
import { Reason } from '../proto/feature/reason_pb';
import { Strategy, FixedStrategy } from '../proto/feature/strategy_pb';
import { Variation } from '../proto/feature/variation_pb';
import { createUser } from '../modelFactory';

function disabledFlag(): Feature {
  const feature = new Feature();
  feature.setId('flag');
  feature.setEnabled(false);
  feature.setOffVariation('');
  const hold = new Variation();
  hold.setId('hold');
  hold.setValue('hold');
  const on = new Variation();
  on.setId('on');
  on.setValue('on');
  feature.setVariationsList([hold, on]);
  const strategy = new Strategy();
  strategy.setType(Strategy.Type.FIXED);
  const fixed = new FixedStrategy();
  fixed.setVariation('on');
  strategy.setFixedStrategy(fixed);
  feature.setDefaultStrategy(strategy);
  return feature;
}

function sibling(): Feature {
  const feature = new Feature();
  feature.setId('sibling');
  feature.setEnabled(true);
  feature.setOffVariation('off');
  const off = new Variation();
  off.setId('off');
  off.setValue('off');
  const on = new Variation();
  on.setId('on');
  on.setValue('on');
  feature.setVariationsList([off, on]);
  const strategy = new Strategy();
  strategy.setType(Strategy.Type.FIXED);
  const fixed = new FixedStrategy();
  fixed.setVariation('on');
  strategy.setFixedStrategy(fixed);
  feature.setDefaultStrategy(strategy);
  return feature;
}

function group(percentage: number): HoldoutGroup {
  return {
    getId: () => 'hg',
    getEnvironmentId: () => 'env',
    getAssignmentType: () => HoldoutAssignment.HASH,
    getSamplingSeed: () => 'seed',
    getBucketStart: () => 0,
    getPercentage: () => percentage,
    getSegmentId: () => '',
    getUpdatedAt: () => 1,
    getFeaturesList: () => [
      {
        getFeatureId: () => 'flag',
        getHoldoutVariationId: () => 'hold',
        getRemovedAt: () => 0,
      },
    ],
  };
}

test('a disabled holdout target without an off variation is assigned only to members', async (t) => {
  const evaluator = new Evaluator();
  const user = createUser('user-1', null);
  const features = [disabledFlag(), sibling()];

  const member = await evaluator.evaluateFeatures(features, user, new Map(), null, '', [
    group(100),
  ]);
  const memberById = new Map(
    member.getEvaluationsList().map((evaluation) => [evaluation.getFeatureId(), evaluation]),
  );
  t.is(memberById.size, 2);
  t.is(memberById.get('flag')?.getReason()?.getType(), Reason.Type.HOLDOUT);
  t.is(memberById.get('flag')?.getVariationId(), 'hold');

  const nonMember = await evaluator.evaluateFeatures(
    [disabledFlag(), sibling()],
    user,
    new Map(),
    null,
    '',
    [group(0)],
  );
  const ids = nonMember.getEvaluationsList().map((evaluation) => evaluation.getFeatureId());
  t.deepEqual(ids, ['sibling']);
  t.deepEqual(nonMember.getArchivedFeatureIdsList(), ['flag']);

  const now = Math.floor(Date.now() / 1000);
  const stale = disabledFlag();
  stale.setUpdatedAt(now);
  const partial = evaluator.evaluateFeaturesByEvaluatedAt(
    [stale, sibling()],
    user,
    new Map(),
    null,
    'prev',
    now - 5,
    false,
    '',
    [group(0)],
  );
  t.false(partial.getForceUpdate());
  t.deepEqual(partial.getArchivedFeatureIdsList(), ['flag']);
});

test('an omitted holdout target does not abort its dependents', async (t) => {
  const evaluator = new Evaluator();
  const user = createUser('user-1', null);
  const target = disabledFlag();
  const dependent = sibling();
  dependent.setId('dependent');
  const prerequisite = new Prerequisite();
  prerequisite.setFeatureId('flag');
  prerequisite.setVariationId('on');
  dependent.setPrerequisitesList([prerequisite]);
  const unrelated = sibling();

  const evaluations = await evaluator.evaluateFeatures(
    [dependent, target, unrelated],
    user,
    new Map(),
    null,
    '',
    [group(0)],
  );
  const byId = new Map(
    evaluations.getEvaluationsList().map((evaluation) => [evaluation.getFeatureId(), evaluation]),
  );
  t.false(byId.has('flag'));
  t.deepEqual(evaluations.getArchivedFeatureIdsList(), ['flag']);
  t.is(byId.get('dependent')?.getReason()?.getType(), Reason.Type.PREREQUISITE);
  t.is(byId.get('dependent')?.getVariationId(), 'off');
  t.is(byId.get('sibling')?.getReason()?.getType(), Reason.Type.DEFAULT);
});
