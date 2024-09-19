import test from 'ava';
import { Clause } from '../proto/feature/clause_pb';
import { ClauseEvaluator } from '../clauseEvaluator';

const clauseEvaluator = new ClauseEvaluator();

test('GreaterFloat', t => {
  const testCases = [
    // Int
    { targetValue: '1', values: ['1'], expected: false },
    { targetValue: '1', values: ['1', '2', '3'], expected: false },
    { targetValue: '1', values: ['a', '1', '2.0'], expected: false },
    { targetValue: '1', values: ['a', 'b', 'c'], expected: false },
    { targetValue: '1', values: ['0a', '1a'], expected: false },
    { targetValue: '1', values: ['0'], expected: true },
    { targetValue: '1', values: ['0.0', '1.0', '2.0'], expected: true },
    { targetValue: '1', values: ['0.9', '1.0', '2.0'], expected: true },
    { targetValue: '1', values: ['0', '1', '2'], expected: true },
    { targetValue: '1', values: ['a', '0', '1.0'], expected: true },
    { targetValue: '1', values: ['a', '0', '1'], expected: true },
    { targetValue: '1', values: ['0a', '0'], expected: true },
    // Float
    { targetValue: '1.0', values: ['1.0', '2.0', '3.0'], expected: false },
    { targetValue: '1.0', values: ['1', '2', '3'], expected: false },
    { targetValue: '1.0', values: ['a', '1', '2.0'], expected: false },
    { targetValue: '1.0', values: ['a', 'b', 'c'], expected: false },
    { targetValue: '1.0', values: ['0', '1.0', '2.0'], expected: true },
    { targetValue: '1.0', values: ['a', '0.0', '1.0'], expected: true },
    { targetValue: '1.2', values: ['a', '1.1', '2.0'], expected: true },
  ];

  testCases.forEach((tc, i) => {
    const clause = new Clause();
    clause.setOperator(Clause.Operator.GREATER);
    clause.setValuesList(tc.values);

    const result = clauseEvaluator.evaluate(tc.targetValue, clause, 'userId', [], {});
    t.is(result, tc.expected, `Test case ${i} failed : targetValue ${tc.targetValue} : value ${tc.values}`);
  });
});

test('GreaterSemver', t => {
  const testCases = [
    { targetValue: '1.0.0', values: ['1.0.0', '0.0', '1.0.1'], expected: false },
    // This case is difference with golang version which expected: false. 
    // Because node.js semver is able to parse v0.0.7 to 0.0.7, 1.0.9-alpha to 1.0.9
    { targetValue: '1.0.0', values: ['1.0.0', '1.0.1', 'v0.0.7'], expected: true },
    { targetValue: '1.0.0', values: ['1.0.0', '1.0.1', '0.0.7a'], expected: false },
    { targetValue: '1.0.0', values: ['1.0.0', '1.0.1', 'a-0.0.7'], expected: false },

    { targetValue: '0.0.8', values: ['1.0.0', '0.0.9', '1.0.1'], expected: false },

    // This case is difference with golang version which expected: false. 
    { targetValue: '1.1.0', values: ['1.1.0', 'v1.0.9', '1.1.1'], expected: true },

    { targetValue: '1.1.0', values: ['1.1.0', '1.0.9-alpha', '1.1.1'], expected: true },
    { targetValue: '1.1.0', values: ['1.1.0', '1.0.9a', '1.1.1'], expected: false },

    // This case is difference with golang version which expected: false.
    { targetValue: '2.1.0', values: ['2.1.0', 'v2.0.9', '2.1.1'], expected: true },

    { targetValue: '1.0.1', values: ['1.0.1', '1.0.0', 'v0.0.7'], expected: true },
    { targetValue: '1.1.1', values: ['1.1.1', 'v1.0.9', '1.1.0'], expected: true },
    { targetValue: '2.1.1', values: ['2.1.1', 'v2.0.9', '2.1.0'], expected: true },
  ];

  testCases.forEach((tc, i) => {
    const clause = new Clause();
    clause.setOperator(Clause.Operator.GREATER);
    clause.setValuesList(tc.values);

    const result = clauseEvaluator.evaluate(tc.targetValue, clause, 'userId', [], {});
    t.is(result, tc.expected, `Test case ${i} failed`);
  });
});

