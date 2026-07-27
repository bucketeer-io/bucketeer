import test from 'ava';
import { Segment, SegmentUser } from '../proto/feature/segment_pb';
import { Rule } from '../proto/feature/rule_pb';
import { Clause } from '../proto/feature/clause_pb';
import { SegmentEvaluator } from '../segmentEvaluator';
import { User } from '../proto/user/user_pb';
import { createUser } from '../modelFactory';

// Define the type for the test cases
interface SegmentEvaluatorTestCase {
  desc: string;
  segmentIDs: string[];
  user: User;
  segments: Map<string, Segment> | null;
  segmentUsers: SegmentUser[];
  expected: boolean;
}

// A user belongs to a segment when it is in the included-user list OR matches
// any of the segment rules. Multiple segment IDs are evaluated with OR.
export const SegmentEvaluatorTestCases: SegmentEvaluatorTestCase[] = [
  {
    desc: 'user is included in all segments',
    segmentIDs: ['segment-1', 'segment-2'],
    user: createUser('user-1', null),
    segments: null,
    segmentUsers: [
      createSegmentUser('segment-1', 'user-1', SegmentUser.State.INCLUDED),
      createSegmentUser('segment-2', 'user-1', SegmentUser.State.INCLUDED),
    ],
    expected: true,
  },
  {
    desc: 'user is included in one of the segments (OR)',
    segmentIDs: ['segment-1', 'segment-2'],
    user: createUser('user-1', null),
    segments: null,
    segmentUsers: [
      createSegmentUser('segment-1', 'user-1', SegmentUser.State.INCLUDED),
      createSegmentUser('segment-2', 'user-1', SegmentUser.State.EXCLUDED),
    ],
    expected: true,
  },
  {
    desc: 'user does not exist in any segments',
    segmentIDs: ['segment-1', 'segment-2'],
    user: createUser('user-1', null),
    segments: null,
    segmentUsers: [
      createSegmentUser('segment-1', 'user-2', SegmentUser.State.INCLUDED),
      createSegmentUser('segment-2', 'user-2', SegmentUser.State.INCLUDED),
    ],
    expected: false,
  },
  {
    desc: 'empty segment IDs',
    segmentIDs: [],
    user: createUser('user-1', null),
    segments: null,
    segmentUsers: [createSegmentUser('segment-1', 'user-1', SegmentUser.State.INCLUDED)],
    expected: false, // No segments to evaluate means no match
  },
  {
    desc: 'single segment ID, user included',
    segmentIDs: ['segment-1'],
    user: createUser('user-1', null),
    segments: null,
    segmentUsers: [createSegmentUser('segment-1', 'user-1', SegmentUser.State.INCLUDED)],
    expected: true,
  },
  {
    desc: 'single segment ID, user excluded',
    segmentIDs: ['segment-1'],
    user: createUser('user-1', null),
    segments: null,
    segmentUsers: [createSegmentUser('segment-1', 'user-1', SegmentUser.State.EXCLUDED)],
    expected: false,
  },
  {
    desc: 'user included in segments, but not all segments defined',
    segmentIDs: ['segment-1', 'segment-2'],
    user: createUser('user-1', null),
    segments: null,
    segmentUsers: [
      createSegmentUser('segment-1', 'user-1', SegmentUser.State.INCLUDED),
      createSegmentUser('segment-3', 'user-1', SegmentUser.State.INCLUDED), // segment-3 is not in the IDs
    ],
    expected: true,
  },
  {
    desc: 'multiple users with mixed states across segments',
    segmentIDs: ['segment-1', 'segment-2'],
    user: createUser('user-2', null),
    segments: null,
    segmentUsers: [
      createSegmentUser('segment-1', 'user-1', SegmentUser.State.INCLUDED),
      createSegmentUser('segment-1', 'user-2', SegmentUser.State.INCLUDED),
      createSegmentUser('segment-2', 'user-2', SegmentUser.State.INCLUDED),
    ],
    expected: true,
  },
  {
    desc: 'rule-only segment: user matches a segment rule',
    segmentIDs: ['segment-1'],
    user: createUser('user-1', { plan: 'premium' }),
    segments: new Map([
      ['segment-1', createSegment('segment-1', [createSegmentRule('plan', ['premium'])])],
    ]),
    segmentUsers: [],
    expected: true,
  },
  {
    desc: 'rule-only segment: user does not match the segment rule',
    segmentIDs: ['segment-1'],
    user: createUser('user-1', { plan: 'free' }),
    segments: new Map([
      ['segment-1', createSegment('segment-1', [createSegmentRule('plan', ['premium'])])],
    ]),
    segmentUsers: [],
    expected: false,
  },
  {
    desc: 'mixed segment: user in the include list, rule does not match',
    segmentIDs: ['segment-1'],
    user: createUser('user-1', { plan: 'free' }),
    segments: new Map([
      ['segment-1', createSegment('segment-1', [createSegmentRule('plan', ['premium'])])],
    ]),
    segmentUsers: [createSegmentUser('segment-1', 'user-1', SegmentUser.State.INCLUDED)],
    expected: true,
  },
  {
    desc: 'mixed segment: user not in the include list, rule matches',
    segmentIDs: ['segment-1'],
    user: createUser('user-2', { plan: 'premium' }),
    segments: new Map([
      ['segment-1', createSegment('segment-1', [createSegmentRule('plan', ['premium'])])],
    ]),
    segmentUsers: [createSegmentUser('segment-1', 'user-1', SegmentUser.State.INCLUDED)],
    expected: true,
  },
  {
    desc: 'missing attribute: segment rule fails closed',
    segmentIDs: ['segment-1'],
    user: createUser('user-1', null),
    segments: new Map([
      ['segment-1', createSegment('segment-1', [createSegmentRule('plan', ['premium'])])],
    ]),
    segmentUsers: [],
    expected: false,
  },
];

function createSegmentUser(
  segmentId: string,
  userId: string,
  state: SegmentUser.StateMap[keyof SegmentUser.StateMap],
): SegmentUser {
  const user = new SegmentUser();
  user.setSegmentId(segmentId);
  user.setUserId(userId);
  user.setState(state);
  return user;
}

function createSegment(id: string, rules: Rule[]): Segment {
  const segment = new Segment();
  segment.setId(id);
  segment.setRulesList(rules);
  return segment;
}

function createSegmentRule(attribute: string, values: string[]): Rule {
  const clause = new Clause();
  clause.setId(`clause-${attribute}`);
  clause.setAttribute(attribute);
  clause.setOperator(Clause.Operator.EQUALS);
  clause.setValuesList(values);
  const rule = new Rule();
  rule.setId(`rule-${attribute}`);
  rule.setClausesList([clause]);
  return rule;
}

SegmentEvaluatorTestCases.forEach(({ desc, segmentIDs, user, segments, segmentUsers, expected }) => {
  test(desc, (t) => {
    const evaluator = new SegmentEvaluator();
    const actual = evaluator.evaluate(segmentIDs, user, segments, segmentUsers);
    t.is(actual, expected);
  });
});
