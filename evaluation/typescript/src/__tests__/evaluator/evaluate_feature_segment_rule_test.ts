import test from 'ava';
import { Evaluator } from '../../evaluation';
import { createFeature, createUser, createSegmentUser } from '../../modelFactory';
import { Clause } from '../../proto/feature/clause_pb';
import { Rule } from '../../proto/feature/rule_pb';
import { Reason } from '../../proto/feature/reason_pb';
import { Segment, SegmentUser } from '../../proto/feature/segment_pb';
import { Strategy } from '../../proto/feature/strategy_pb';
import { User } from '../../proto/user/user_pb';

// Covers the full evaluation path for rule-based segments: a flag rule with a
// SEGMENT clause referencing a segment whose membership is defined by rules
// on user attributes. Mirrors TestEvaluateFeaturesWithRuleBasedSegment in
// evaluation/go.
function newTestFeature() {
  return createFeature({
    id: 'feature-id',
    name: 'test feature',
    version: 1,
    enabled: true,
    variations: [
      { id: 'variation-A', value: 'A', name: 'Variation A', description: '' },
      { id: 'variation-B', value: 'B', name: 'Variation B', description: '' },
    ],
    rules: [
      {
        id: 'rule-1',
        attribute: '',
        operator: Clause.Operator.SEGMENT,
        values: ['segment-1'],
        fixedVariation: 'variation-B',
      },
    ],
    defaultStrategy: { type: Strategy.Type.FIXED, variation: 'variation-A' },
  });
}

function newRuleBasedSegments(): Map<string, Segment> {
  const clause = new Clause();
  clause.setId('segment-clause-1');
  clause.setAttribute('plan');
  clause.setOperator(Clause.Operator.EQUALS);
  clause.setValuesList(['premium']);
  const rule = new Rule();
  rule.setId('segment-rule-1');
  rule.setClausesList([clause]);
  const segment = new Segment();
  segment.setId('segment-1');
  segment.setRulesList([rule]);
  return new Map([['segment-1', segment]]);
}

interface TestCase {
  desc: string;
  user: User;
  segmentUsers: Map<string, SegmentUser[]>;
  segments: Map<string, Segment> | null;
  expectedVariation: string;
  expectedReasonType: Reason.TypeMap[keyof Reason.TypeMap];
  expectedRuleId: string;
}

const testCases: TestCase[] = [
  {
    desc: 'user matches the segment rule by attribute',
    user: createUser('user-1', { plan: 'premium' }),
    segmentUsers: new Map(),
    segments: newRuleBasedSegments(),
    expectedVariation: 'variation-B',
    expectedReasonType: Reason.Type.RULE,
    expectedRuleId: 'rule-1',
  },
  {
    desc: 'user does not match the segment rule',
    user: createUser('user-1', { plan: 'free' }),
    segmentUsers: new Map(),
    segments: newRuleBasedSegments(),
    expectedVariation: 'variation-A',
    expectedReasonType: Reason.Type.DEFAULT,
    expectedRuleId: '',
  },
  {
    desc: 'mixed: user in the include list matches even when the rule does not',
    user: createUser('listed-user', { plan: 'free' }),
    segmentUsers: new Map([
      ['segment-1', [createSegmentUser('listed-user', 'segment-1', SegmentUser.State.INCLUDED)]],
    ]),
    segments: newRuleBasedSegments(),
    expectedVariation: 'variation-B',
    expectedReasonType: Reason.Type.RULE,
    expectedRuleId: 'rule-1',
  },
  {
    desc: 'backward compat: null segments map evaluates the include list only',
    user: createUser('user-1', { plan: 'premium' }),
    segmentUsers: new Map(),
    segments: null,
    expectedVariation: 'variation-A',
    expectedReasonType: Reason.Type.DEFAULT,
    expectedRuleId: '',
  },
];

testCases.forEach((tc) => {
  test(`rule-based segment: ${tc.desc}`, async (t) => {
    const evaluator = new Evaluator();
    const feature = newTestFeature();
    const result = await evaluator.evaluateFeatures(
      [feature],
      tc.user,
      tc.segmentUsers,
      tc.segments,
      '',
    );
    const evaluation = result
      .getEvaluationsList()
      .find((e) => e.getFeatureId() === feature.getId());
    t.truthy(evaluation);
    t.is(evaluation?.getVariationId(), tc.expectedVariation);
    t.is(evaluation?.getReason()?.getType(), tc.expectedReasonType);
    t.is(evaluation?.getReason()?.getRuleId() ?? '', tc.expectedRuleId);
  });
});
