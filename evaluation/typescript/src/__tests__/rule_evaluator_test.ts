import test from 'ava';
import { Rule } from '../proto/feature/rule_pb';
import { Feature } from '../proto/feature/feature_pb';
import { FixedStrategy, Strategy } from '../proto/feature/strategy_pb';
import { Clause } from '../proto/feature/clause_pb';
import { RuleEvaluator } from '../ruleEvaluator';
import { User } from '../proto/user/user_pb';
import { SegmentUser } from '../proto/feature/segment_pb';
import { createRule, createSegmentUser, createUser, createVariation } from '../modelFactory';

// Define the type for the test cases
interface TestCase {
  user: User;
  expected: Rule | null;
  description: string;
}

const TestRuleList = {
  EQUALS: createRule(
    'rule-id-1',
    'full-name',
    Clause.Operator.EQUALS,
    ['bucketeer project'],
    'variation-A',
  ),
  STARTS_WITH: createRule(
    'rule-id-2',
    'first-name',
    Clause.Operator.STARTS_WITH,
    ['buck'],
    'variation-A',
  ),
  ENDS_WITH: createRule(
    'rule-id-3',
    'last-name',
    Clause.Operator.ENDS_WITH,
    ['ject'],
    'variation-A',
  ),
  SEGMENT: createRule(
    'rule-id-4',
    '',
    Clause.Operator.SEGMENT,
    ['segment-id-1', 'segment-id-2'],
    'variation-B',
  ),
  IN: createRule('rule-id-5', 'email', Clause.Operator.IN, ['bucketeer@gmail.com'], 'variation-B'),
};

function newTestFeature(): Feature {
  const feature = new Feature();
  feature.setId('feature-id');
  feature.setName('test feature');
  feature.setVersion(1);
  feature.setCreatedAt(Math.floor(Date.now() / 1000)); // Unix timestamp in seconds

  const variations = [
    createVariation('variation-A', 'A', 'Variation A', 'Thing does A'),
    createVariation('variation-B', 'B', 'Variation B', 'Thing does B'),
  ];
  feature.setVariationsList(variations);
  const rules = [
    TestRuleList.EQUALS,
    TestRuleList.STARTS_WITH,
    TestRuleList.ENDS_WITH,
    TestRuleList.SEGMENT,
    TestRuleList.IN,
  ];

  feature.setRulesList(rules);

  const fixedStrategy = new FixedStrategy();
  fixedStrategy.setVariation('variation-B');
  const strategy = new Strategy();
  strategy.setType(Strategy.Type.FIXED);
  strategy.setFixedStrategy(fixedStrategy);

  feature.setDefaultStrategy(strategy);

  return feature;
}

function newSegmentUserIDs(): SegmentUser[] {
  return [
    createSegmentUser('user-id-1', 'segment-id-1', SegmentUser.State.INCLUDED),
    createSegmentUser('user-id-1', 'segment-id-2', SegmentUser.State.INCLUDED),
    createSegmentUser('user-id-2', 'segment-id-1', SegmentUser.State.INCLUDED),
    createSegmentUser('user-id-2', 'segment-id-2', SegmentUser.State.INCLUDED),
    createSegmentUser('user-id-3', 'segment-id-1', SegmentUser.State.INCLUDED),
    createSegmentUser('user-id-4', 'segment-id-2', SegmentUser.State.INCLUDED),
  ];
}

// Define test cases
const testcases: TestCase[] = [
  {
    user: createUser('user-id-1', { 'full-name': 'bucketeer project' }),
    expected: TestRuleList.EQUALS,
    description: 'TestRuleList.EQUALS',
  },
  {
    user: createUser('user-id-1', { 'first-name': 'bucketeer' }),
    expected: TestRuleList.STARTS_WITH,
    description: 'TestRuleList.STARTS_WITH',
  },
  {
    user: createUser('user-id-1', { 'last-name': 'project' }),
    expected: TestRuleList.ENDS_WITH,
    description: 'TestRuleList.ENDS_WITH',
  },
  {
    user: createUser('user-id-3', { email: 'bucketeer@gmail.com' }),
    expected: TestRuleList.IN,
    description: 'TestRuleList.IN',
  },
  {
    user: createUser('user-id-1', null),
    expected: TestRuleList.SEGMENT,
    description: 'TestRuleList.SEGMENT - user-id-1',
  },
  {
    user: createUser('user-id-2', null),
    expected: TestRuleList.SEGMENT,
    description: 'TestRuleList.SEGMENT - user-id-2',
  },
  {
    user: createUser('user-id-3', null),
    expected: null,
    description: 'No rule matched - user-id-3',
  },
  {
    user: createUser('user-id-4', null),
    expected: null,
    description: 'No rule matched - user-id-4',
  },
];

// Test cases
testcases.forEach(({ user, expected, description }) => {
  test(description, (t) => {
    const ruleEvaluator = new RuleEvaluator();
    const values = newSegmentUserIDs();
    const actual = ruleEvaluator.evaluate(newTestFeature().getRulesList(), user, values, {});
    t.deepEqual(actual, expected);
  });
});
