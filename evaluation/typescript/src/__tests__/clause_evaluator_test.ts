import test from 'ava';
import { Clause } from '../proto/feature/clause_pb';
import { ClauseEvaluator } from '../clauseEvaluator';
import { SegmentEvaluatorTestCases } from './segment_evaluator_test';

const clauseEvaluator = new ClauseEvaluator();

test('GreaterFloat', (t) => {
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
    t.is(
      result,
      tc.expected,
      `Test case ${i} failed : targetValue ${tc.targetValue} : value ${tc.values}`,
    );
  });
});

test('GreaterSemver', (t) => {
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

test('TestGreaterString', (t) => {
  const testcases = [
    {
      targetValue: 'b',
      values: ['c', 'd', 'e'],
      expected: false,
    },
    {
      targetValue: 'v1.0.0',
      values: ['v2.0.0', 'v1.0.9', 'v1.0.8'],
      expected: false,
    },
    {
      targetValue: 'b',
      values: ['1', 'a', '2.0'],
      expected: true,
    },
    {
      targetValue: 'b',
      values: ['c', 'd', 'a'],
      expected: true,
    },
    {
      targetValue: 'v1.0.0',
      values: ['v1.0.0', 'v1.0.9', 'v0.0.9'],
      expected: true,
    },
  ];

  const clauseEvaluator = new ClauseEvaluator();

  testcases.forEach((tc, i) => {
    const clause = new Clause();
    clause.setOperator(Clause.Operator.GREATER);
    clause.setValuesList(tc.values);

    const result = clauseEvaluator.evaluate(tc.targetValue, clause, 'userId', [], {});
    t.is(
      result,
      tc.expected,
      `Test case ${i} failed : targetValue ${tc.targetValue} : value ${tc.values}`,
    );
  });
});

test('TestGreaterOrEqualFloat', (t) => {
  const testcases = [
    // Int test cases
    {
      targetValue: '1',
      values: ['2'],
      expected: false,
    },
    {
      targetValue: '1',
      values: ['2', '3', '4'],
      expected: false,
    },
    {
      targetValue: '1',
      values: ['2.0', '3.0', '4.0'],
      expected: false,
    },
    {
      targetValue: '1',
      values: ['a', '2', '3.0'],
      expected: false,
    },
    {
      targetValue: '1',
      values: ['a', 'b', 'c'],
      expected: false,
    },
    {
      targetValue: '1',
      values: ['0a', '1a'],
      expected: false,
    },
    {
      targetValue: '1',
      values: ['1'],
      expected: true,
    },
    {
      targetValue: '1',
      values: ['0', '1', '2'],
      expected: true,
    },
    {
      targetValue: '1',
      values: ['0.0', '1.0', '2.0'],
      expected: true,
    },
    {
      targetValue: '1',
      values: ['1.0', '2.0', '3.0'],
      expected: true,
    },
    {
      targetValue: '1',
      values: ['1', '2', '3'],
      expected: true,
    },
    {
      targetValue: '1',
      values: ['a', '0', '1.0'],
      expected: true,
    },
    {
      targetValue: '1',
      values: ['a', '0', '1'],
      expected: true,
    },
    {
      targetValue: '1',
      values: ['a', '1', '2.0'],
      expected: true,
    },
    {
      targetValue: '1',
      values: ['a', '1.0', '2'],
      expected: true,
    },
    {
      targetValue: '1',
      values: ['0a', '0'],
      expected: true,
    },
    // Float test cases
    {
      targetValue: '1.0',
      values: ['2.0', '3.0', '4.0'],
      expected: false,
    },
    {
      targetValue: '1.0',
      values: ['2', '3', '4'],
      expected: false,
    },
    {
      targetValue: '1.0',
      values: ['a', '1.1', '2.0'],
      expected: false,
    },
    {
      targetValue: '1.0',
      values: ['a', 'b', 'c'],
      expected: false,
    },
    {
      targetValue: '1.0',
      values: ['0.9', '2.0', '3.0'],
      expected: true,
    },
    {
      targetValue: '1.0',
      values: ['a', '0', '2.0'],
      expected: true,
    },
    {
      targetValue: '1.1',
      values: ['1', '2.0', '3.0'],
      expected: true,
    },
    {
      targetValue: '1.1',
      values: ['1.1', '2.0', '3.0'],
      expected: true,
    },
    {
      targetValue: '1.1',
      values: ['a', '1.0', '2.0'],
      expected: true,
    },
  ];

  const clauseEvaluator = new ClauseEvaluator();

  testcases.forEach((tc, i) => {
    const clause = new Clause();
    clause.setOperator(Clause.Operator.GREATER_OR_EQUAL); // Greater or Equal operator
    clause.setValuesList(tc.values);

    const result = clauseEvaluator.evaluate(tc.targetValue, clause, 'userId', [], {});
    t.is(
      result,
      tc.expected,
      `Test case ${i} failed: targetValue ${tc.targetValue} : values ${tc.values}`,
    );
  });
});

test('TestGreaterOrEqualSemver', (t) => {
  const testcases = [
    {
      targetValue: '1.0.0',
      values: ['1.0.1', '0.0', '1.0.2'],
      expected: false,
    },
    // This case is difference with golang version which expected: false.
    // Because node.js semver is able to parse v0.0.7 to 0.0.7, 1.0.9-alpha to 1.0.9
    {
      targetValue: '1.0.0',
      values: ['1.0.1', '1.0.2', 'v0.0.7'],
      expected: true,
    },
    {
      targetValue: '1.0.0',
      values: ['1.0.1', '1.0.2', '0.0.7-alpha'],
      expected: true,
    },
    {
      targetValue: '1.0.0',
      values: ['1.0.1', '1.0.2', '0.0.7a'],
      expected: false,
    },
    {
      targetValue: '0.0.8',
      values: ['1.0.0', '0.0.9', '1.0.1'],
      expected: false,
    },
    // This case is difference with golang version which expected: false.
    {
      targetValue: '1.1.0',
      values: ['1.1.1', 'v1.0.9', '1.1.2'],
      expected: true,
    },
    // This case is difference with golang version which expected: false.
    {
      targetValue: '2.1.0',
      values: ['2.1.1', 'v2.0.9', '2.1.2'],
      expected: true,
    },
    {
      targetValue: '1.0.0',
      values: ['1.0.1', '1.0.0', 'v0.0.7'],
      expected: true,
    },
    {
      targetValue: '1.1.1',
      values: ['1.1.2', 'v1.0.9', '1.1.1'],
      expected: true,
    },
    {
      targetValue: '2.1.1',
      values: ['2.1.2', 'v2.0.9', '2.1.1'],
      expected: true,
    },
    {
      targetValue: '1.0.1',
      values: ['1.0.2', '1.0.1', 'v0.0.7'],
      expected: true,
    },
    {
      targetValue: '1.1.1',
      values: ['1.1.2', 'v1.0.9', '1.1.0'],
      expected: true,
    },
    {
      targetValue: '2.1.1',
      values: ['2.1.2', 'v2.0.9', '2.1.0'],
      expected: true,
    },
  ];

  const clauseEvaluator = new ClauseEvaluator();

  testcases.forEach((tc, i) => {
    const clause = new Clause();
    clause.setOperator(Clause.Operator.GREATER_OR_EQUAL);
    clause.setValuesList(tc.values);

    const result = clauseEvaluator.evaluate(tc.targetValue, clause, 'userId', [], {});
    t.is(
      result,
      tc.expected,
      `Test case ${i} failed: targetValue ${tc.targetValue} : values ${tc.values}`,
    );
  });
});

test('TestGreaterOrEqualString', (t) => {
  const testcases = [
    {
      targetValue: 'b',
      values: ['c', 'd', 'e'],
      expected: false,
    },
    {
      targetValue: 'v1.0.0',
      values: ['v2.0.0', 'v1.0.9', 'v1.0.8'],
      expected: false,
    },
    {
      targetValue: 'b',
      values: ['1', 'a', '2.0'],
      expected: true,
    },
    {
      targetValue: 'b',
      values: ['d', 'c', 'b'],
      expected: true,
    },
    {
      targetValue: 'b',
      values: ['c', 'd', 'a'],
      expected: true,
    },
    {
      targetValue: 'b',
      values: ['d', 'c', 'b'],
      expected: true,
    },
    {
      targetValue: 'v1.0.0',
      values: ['v1.0.8', 'v1.0.9', 'v1.0.0'],
      expected: true,
    },
    {
      targetValue: 'v1.0.0',
      values: ['v1.0.8', 'v1.0.9', 'v0.0.9'],
      expected: true,
    },
  ];

  const clauseEvaluator = new ClauseEvaluator();

  testcases.forEach((tc, i) => {
    const clause = new Clause();
    clause.setOperator(Clause.Operator.GREATER_OR_EQUAL);
    clause.setValuesList(tc.values);

    const result = clauseEvaluator.evaluate(tc.targetValue, clause, 'userId', [], {});
    t.is(
      result,
      tc.expected,
      `Test case ${i} failed: targetValue ${tc.targetValue} : values ${tc.values}`,
    );
  });
});

test('TestLessThanSemver', (t) => {
  const testcases = [
    {
      targetValue: '1.0.0',
      values: ['1.0.0', '0.0', '0.0.9'],
      expected: false,
    },
    {
      targetValue: '1.0.0',
      values: ['1.0.0', 'v0.0.8', '0.0.7'],
      expected: false,
    },
    {
      targetValue: '0.0.8',
      values: ['0.0.8', '0.0.7', 'v0.0.9'],
      expected: true,
    },
    {
      targetValue: '1.1.0',
      values: ['1.1.0', 'v1.0.9', '1.0.8'],
      expected: false,
    },
    {
      targetValue: '2.1.0',
      values: ['2.1.0', 'v2.0.9', '2.0.9'],
      expected: false,
    },
    {
      targetValue: '1.0.1',
      values: ['1.0.1', 'v0.0.7', '1.0.2'],
      expected: true,
    },
    {
      targetValue: '1.1.1',
      values: ['1.1.1', 'v1.0.9', '1.1.2'],
      expected: true,
    },
    {
      targetValue: '2.1.1',
      values: ['2.1.1', 'v2.0.9', '2.1.2'],
      expected: true,
    },
  ];

  const clauseEvaluator = new ClauseEvaluator();

  testcases.forEach((tc, i) => {
    const clause = new Clause();
    clause.setOperator(Clause.Operator.LESS);
    clause.setValuesList(tc.values);

    const result = clauseEvaluator.evaluate(tc.targetValue, clause, 'userId', [], {});
    t.is(
      result,
      tc.expected,
      `Test case ${i} failed: targetValue ${tc.targetValue} : values ${tc.values}`,
    );
  });
});

test('TestLessFloat', (t) => {
  const testcases = [
    // Int cases
    {
      targetValue: '3',
      values: ['3'],
      expected: false,
    },
    {
      targetValue: '3',
      values: ['1', '2', '3'],
      expected: false,
    },
    {
      targetValue: '3',
      values: ['a', '1', '2.0'],
      expected: false,
    },
    {
      targetValue: '3',
      values: ['a', 'b', 'c'],
      expected: false,
    },
    {
      targetValue: '3',
      values: ['0a', '1a'],
      expected: false,
    },
    {
      targetValue: '3',
      values: ['4'],
      expected: true,
    },
    {
      targetValue: '3',
      values: ['2.0', '3.0', '4.0'],
      expected: true,
    },
    {
      targetValue: '3',
      values: ['1.0', '2.0', '3.1'],
      expected: true,
    },
    {
      targetValue: '3',
      values: ['2', '3', '4'],
      expected: true,
    },
    {
      targetValue: '3',
      values: ['d', '3', '3.5'],
      expected: true,
    },
    {
      targetValue: '3',
      values: ['a', '0', '4'],
      expected: true,
    },
    {
      targetValue: '3',
      values: ['4a', '4'],
      expected: true,
    },
    // Float cases
    {
      targetValue: '3.0',
      values: ['1.0', '2.0', '3.0'],
      expected: false,
    },
    {
      targetValue: '3.0',
      values: ['1', '2', '3'],
      expected: false,
    },
    {
      targetValue: '3.0',
      values: ['a', '1', '2.0'],
      expected: false,
    },
    {
      targetValue: '3.0',
      values: ['a', 'b', 'c'],
      expected: false,
    },
    {
      targetValue: '3.0',
      values: ['2', '3.0', '3.1'],
      expected: true,
    },
    {
      targetValue: '3.0',
      values: ['a', '0.0', '3.9'],
      expected: true,
    },
    {
      targetValue: '3.2',
      values: ['a', '1.1', '3.5'],
      expected: true,
    },
  ];

  const clauseEvaluator = new ClauseEvaluator();

  testcases.forEach((tc, i) => {
    const clause = new Clause();
    clause.setOperator(Clause.Operator.LESS);
    clause.setValuesList(tc.values);

    const result = clauseEvaluator.evaluate(tc.targetValue, clause, 'userId', [], {});
    t.is(
      result,
      tc.expected,
      `Test case ${i} failed: targetValue ${tc.targetValue} : values ${tc.values}`,
    );
  });
});

test('TestLessString', (t) => {
  const testcases = [
    {
      targetValue: 'c',
      values: ['c', 'b', 'a'],
      expected: false,
    },
    {
      targetValue: 'c',
      values: ['1', 'a', '2.0'],
      expected: false,
    },
    {
      targetValue: 'v2.0.0',
      values: ['v2.0.0', 'v1.0.9', 'v1.0.8'],
      expected: false,
    },
    {
      targetValue: 'c',
      values: ['b', 'c', 'd'],
      expected: true,
    },
    {
      targetValue: 'c',
      values: ['3', '1.0', 'd'],
      expected: true,
    },
    {
      targetValue: 'v2.0.0',
      values: ['v1.0.0', 'v1.0.9', 'v2.1.0'],
      expected: true,
    },
  ];

  const clauseEvaluator = new ClauseEvaluator();

  testcases.forEach((tc, i) => {
    const clause = new Clause();
    clause.setOperator(Clause.Operator.LESS);
    clause.setValuesList(tc.values);

    const result = clauseEvaluator.evaluate(tc.targetValue, clause, 'userId', [], {});
    t.is(
      result,
      tc.expected,
      `Test case ${i} failed: targetValue ${tc.targetValue} : values ${tc.values}`,
    );
  });
});

test('TestLessOrEqualFloat', (t) => {
  const testcases = [
    // Int
    {
      targetValue: '3',
      values: ['2'],
      expected: false,
    },
    {
      targetValue: '3',
      values: ['0', '1', '2'],
      expected: false,
    },
    {
      targetValue: '3',
      values: ['0', '1.0', '2.0'],
      expected: false,
    },
    {
      targetValue: '3',
      values: ['a', '1', '2.0'],
      expected: false,
    },
    {
      targetValue: '3',
      values: ['a', 'b', 'c'],
      expected: false,
    },
    {
      targetValue: '3',
      values: ['3a', '4a'],
      expected: false,
    },
    {
      targetValue: '3',
      values: ['3'],
      expected: true,
    },
    {
      targetValue: '3',
      values: ['2', '3', '4'],
      expected: true,
    },
    {
      targetValue: '3',
      values: ['1.0', '2.0', '3.0'],
      expected: true,
    },
    {
      targetValue: '3',
      values: ['1.0', '2.0', '3.1'],
      expected: true,
    },
    {
      targetValue: '3',
      values: ['1', '2', '4'],
      expected: true,
    },
    {
      targetValue: '3',
      values: ['a', '0', '3.0'],
      expected: true,
    },
    {
      targetValue: '3',
      values: ['a', '1.0', '4'],
      expected: true,
    },
    {
      targetValue: '3',
      values: ['a', '1', '3.5'],
      expected: true,
    },
    {
      targetValue: '3',
      values: ['3a', '3'],
      expected: true,
    },
    // Float
    {
      targetValue: '3.0',
      values: ['0', '1.0', '2.0'],
      expected: false,
    },
    {
      targetValue: '3.0',
      values: ['0', '1', '2'],
      expected: false,
    },
    {
      targetValue: '3.0',
      values: ['a', '1.1', '2.0'],
      expected: false,
    },
    {
      targetValue: '3.0',
      values: ['a', 'b', 'c'],
      expected: false,
    },
    {
      targetValue: '3.0',
      values: ['0.9', '2.0', '3.0'],
      expected: true,
    },
    {
      targetValue: '3.0',
      values: ['a', '0', '3.1'],
      expected: true,
    },
    {
      targetValue: '3.1',
      values: ['1', '2.0', '3.9'],
      expected: true,
    },
    {
      targetValue: '3.1',
      values: ['1.1', '2.0', '4'],
      expected: true,
    },
    {
      targetValue: '3.1',
      values: ['a', '1.0', '3.1'],
      expected: true,
    },
  ];

  const clauseEvaluator = new ClauseEvaluator();

  testcases.forEach((tc, i) => {
    const clause = new Clause();
    clause.setOperator(Clause.Operator.LESS_OR_EQUAL);
    clause.setValuesList(tc.values);

    const result = clauseEvaluator.evaluate(tc.targetValue, clause, 'userId', [], {});
    t.is(
      result,
      tc.expected,
      `Test case ${i} failed: targetValue ${tc.targetValue} : values ${tc.values}`,
    );
  });
});

test('TestLessThanOrEqualSemver', (t) => {
  const testcases = [
    {
      targetValue: '1.0.1',
      values: ['1.0.0', '0.0', '0.0.9'],
      expected: false,
    },
    {
      targetValue: '1.0.1',
      values: ['1.0.0', 'v0.0.8', '0.0.7'],
      expected: false,
    },
    {
      targetValue: '0.0.9',
      values: ['0.0.8', '0.0.7', 'v0.0.9'],
      expected: true,
    },
    {
      targetValue: '1.1.1',
      values: ['1.1.0', 'v1.0.9', '1.0.8'],
      expected: false,
    },
    {
      targetValue: '2.1.1',
      values: ['2.1.0', 'v2.0.9', '2.0.9'],
      expected: false,
    },
    {
      targetValue: '1.0.1',
      values: ['1.0.1', 'v0.0.7', '1.0.0'],
      expected: true,
    },
    {
      targetValue: '1.1.1',
      values: ['1.1.1', 'v1.0.9', '1.1.0'],
      expected: true,
    },
    {
      targetValue: '2.1.1',
      values: ['2.1.1', 'v2.0.9', '2.1.0'],
      expected: true,
    },
    {
      targetValue: '1.0.1',
      values: ['1.0.0', 'v0.0.7', '1.0.2'],
      expected: true,
    },
    {
      targetValue: '1.1.1',
      values: ['1.1.0', 'v1.0.9', '1.1.2'],
      expected: true,
    },
    {
      targetValue: '2.1.1',
      values: ['2.1.0', 'v2.0.9', '2.1.2'],
      expected: true,
    },
  ];

  const clauseEvaluator = new ClauseEvaluator();

  testcases.forEach((tc, i) => {
    const clause = new Clause();
    clause.setOperator(Clause.Operator.LESS_OR_EQUAL);
    clause.setValuesList(tc.values);

    const result = clauseEvaluator.evaluate(tc.targetValue, clause, 'userId', [], {});
    t.is(
      result,
      tc.expected,
      `Test case ${i} failed: targetValue ${tc.targetValue} : values ${tc.values}`,
    );
  });
});

test('TestLessOrEqualString', (t) => {
  const testcases = [
    {
      targetValue: 'd',
      values: ['a', 'b', 'c'],
      expected: false,
    },
    {
      targetValue: 'c',
      values: ['1', 'a', '2.0'],
      expected: false,
    },
    {
      targetValue: 'v2.0.0',
      values: ['v1.0.0', 'v1.0.9', 'v1.0.8'],
      expected: false,
    },
    {
      targetValue: 'c',
      values: ['3.0', 'c', 'b'],
      expected: true,
    },
    {
      targetValue: 'c',
      values: ['c', 'b', 'a'],
      expected: true,
    },
    {
      targetValue: 'c',
      values: ['a', 'b', 'd'],
      expected: true,
    },
    {
      targetValue: 'v2.0.0',
      values: ['v1.0.0', 'v1.0.9', 'v2.0.0'],
      expected: true,
    },
    {
      targetValue: 'v2.0.0',
      values: ['v1.0.0', 'v1.0.9', 'v2.0.1'],
      expected: true,
    },
  ];

  const clauseEvaluator = new ClauseEvaluator();

  testcases.forEach((tc, i) => {
    const clause = new Clause();
    clause.setOperator(Clause.Operator.LESS_OR_EQUAL); // LESS_OR_EQUAL operator
    clause.setValuesList(tc.values);

    const result = clauseEvaluator.evaluate(tc.targetValue, clause, 'userId', [], {});
    t.is(
      result,
      tc.expected,
      `Test case ${i} failed: targetValue ${tc.targetValue} : values ${tc.values}`,
    );
  });
});

test('TestBeforeInt', (t) => {
  const testcases = [
    // Int
    {
      targetValue: '1519223320',
      values: ['1419223320'],
      expected: false,
    },
    {
      targetValue: '1519223320',
      values: ['1619223320'],
      expected: true,
    },
    {
      targetValue: '1519223320',
      values: ['1519223320', '1519200000'],
      expected: false,
    },
    // Strings
    {
      targetValue: '15192XXX23320',
      values: ['1519223330', '1519223311', '1519223300'],
      expected: false,
    },
    {
      targetValue: '1519223320',
      values: ['1519223320', '1519200000', '15192XXX23300'],
      expected: false,
    },
    // Float
    {
      // TODO: review me again
      // This case is difference with golang version which expected: false.
      // when parse `15192233.30` to int, it shoule be `15192233`
      // So 15192233 is before 1519223330, 1519223311, 1519223300
      targetValue: '15192233.30',
      values: ['1519223330', '1519223311', '1519223300'],
      expected: true,
    },
    {
      targetValue: '1519223320',
      values: ['1519223320', '1519200000', '15192233.00'],
      expected: false,
    },
  ];

  const clauseEvaluator = new ClauseEvaluator();

  testcases.forEach((tc, i) => {
    const clause = new Clause();
    clause.setOperator(Clause.Operator.BEFORE); // BEFORE operator
    clause.setValuesList(tc.values);

    const result = clauseEvaluator.evaluate(tc.targetValue, clause, 'userId', [], {});
    t.is(
      result,
      tc.expected,
      `Test case ${i} failed: targetValue ${tc.targetValue} : values ${tc.values}`,
    );
  });
});

test('TestAfterInt', (t) => {
  const testcases = [
    // Int
    {
      targetValue: '1519223320',
      values: ['1419223320'],
      expected: true,
    },
    {
      targetValue: '1519223320',
      values: ['1619223320'],
      expected: false,
    },
    {
      targetValue: '1519223320',
      values: ['1519223320', '1519223319'],
      expected: true,
    },
    // Strings
    {
      targetValue: '15192XXX23320',
      values: ['1519223330', '1519223311', '1519223300'],
      expected: false,
    },
    {
      targetValue: '1519223320',
      values: ['1519223320', '1519200000', '15192XXX23300'],
      expected: true,
    },
    // Float
    {
      targetValue: '15192233.30',
      values: ['1519223330', '1519223311', '1519223300'],
      expected: false,
    },
    {
      targetValue: '1519223320',
      values: ['1519223320', '1519200000', '15192233.00'],
      expected: true,
    },
  ];

  const clauseEvaluator = new ClauseEvaluator();

  testcases.forEach((tc, i) => {
    const clause = new Clause();
    clause.setOperator(Clause.Operator.AFTER);
    clause.setValuesList(tc.values);

    const result = clauseEvaluator.evaluate(tc.targetValue, clause, 'userId', [], {});
    t.is(
      result,
      tc.expected,
      `Test case ${i} failed: targetValue ${tc.targetValue} : values ${tc.values}`,
    );
  });
});

// TestInOperator does not exist on the Go version testcases
test('TestInOperator', (t) => {
  const testcases = [
    // Exact match
    {
      targetValue: '1519223320',
      values: ['1519223320', '1419223320', '1619223320'],
      expected: true, // The targetValue exists in the values array
    },
    {
      targetValue: '1519223321',
      values: ['1519223320', '1419223320', '1619223320'],
      expected: false, // The targetValue does not exist in the values array
    },
    // Strings
    {
      targetValue: 'apple',
      values: ['banana', 'orange', 'apple'],
      expected: true, // The targetValue exists in the values array
    },
    {
      targetValue: 'grape',
      values: ['banana', 'orange', 'apple'],
      expected: false, // The targetValue does not exist in the values array
    },
    // Numbers as strings
    {
      targetValue: '123',
      values: ['456', '789', '123'],
      expected: true, // The targetValue exists in the values array
    },
    {
      targetValue: '321',
      values: ['456', '789', '123'],
      expected: false, // The targetValue does not exist in the values array
    },
    // Mixed types
    {
      targetValue: '123',
      values: ['123', 'apple', 'banana'],
      expected: true, // The targetValue exists in the values array
    },
    {
      targetValue: 'apple',
      values: ['123', '456', '789'],
      expected: false, // The targetValue does not exist in the values array
    },
    // Empty values list
    {
      targetValue: 'anything',
      values: [],
      expected: false, // The values array is empty, so targetValue cannot exist
    },
    // Case sensitivity
    {
      targetValue: 'Apple',
      values: ['apple', 'banana', 'cherry'],
      expected: false, // Case-sensitive comparison, 'Apple' != 'apple'
    },
    {
      targetValue: 'apple',
      values: ['apple', 'banana', 'cherry'],
      expected: true, // Exact match, case-sensitive
    },
  ];

  const clauseEvaluator = new ClauseEvaluator();

  testcases.forEach((tc, i) => {
    const clause = new Clause();
    clause.setOperator(Clause.Operator.IN); // Change operator to 'IN'
    clause.setValuesList(tc.values);

    const result = clauseEvaluator.evaluate(tc.targetValue, clause, 'userId', [], {});
    t.is(
      result,
      tc.expected,
      `Test case ${i} failed: targetValue ${tc.targetValue} : values ${tc.values}`,
    );
  });
});

// TestStartsWithOperator does not exist on the Go version testcases
test('TestStartsWithOperator', (t) => {
  const testcases = [
    // Exact starts with match
    {
      targetValue: '1519223320',
      values: ['1519'],
      expected: true, // targetValue starts with '1519'
    },
    {
      targetValue: '1519223320',
      values: ['1619'],
      expected: false, // targetValue does not start with '1619'
    },
    // Multiple possible prefixes
    {
      targetValue: '1519223320',
      values: ['1519', '1419', '1319'],
      expected: true, // targetValue starts with '1519'
    },
    {
      targetValue: '1519223320',
      values: ['1619', '1719', '1819'],
      expected: false, // targetValue does not start with any of these values
    },
    // Full match is also a valid "starts with"
    {
      targetValue: 'apple',
      values: ['apple'],
      expected: true, // targetValue fully matches the value, so it starts with 'apple'
    },
    // String comparisons
    {
      targetValue: 'banana',
      values: ['ban'],
      expected: true, // targetValue starts with 'ban'
    },
    {
      targetValue: 'banana',
      values: ['ba', 'ban', 'bana'],
      expected: true, // targetValue starts with 'ban' or 'ba' or 'bana'
    },
    {
      targetValue: 'banana',
      values: ['car', 'dog'],
      expected: false, // targetValue does not start with any of these values
    },
    // Numbers as strings
    {
      targetValue: '123456',
      values: ['123'],
      expected: true, // targetValue starts with '123'
    },
    {
      targetValue: '123456',
      values: ['456'],
      expected: false, // targetValue does not start with '456'
    },
    // Mixed strings and numbers
    {
      targetValue: 'hello123',
      values: ['hello'],
      expected: true, // targetValue starts with 'hello'
    },
    {
      targetValue: 'hello123',
      values: ['123'],
      expected: false, // targetValue does not start with '123'
    },
    // Empty values list
    {
      targetValue: 'anything',
      values: [],
      expected: false, // An empty list means no prefix matches, so it should return false
    },
    // Case sensitivity
    {
      targetValue: 'Apple',
      values: ['apple'],
      expected: false, // Case-sensitive comparison, 'Apple' does not start with 'apple'
    },
    {
      targetValue: 'apple',
      values: ['apple'],
      expected: true, // Exact match, targetValue starts with 'apple'
    },
    // Special characters
    {
      targetValue: 'hello_world',
      values: ['hello_'],
      expected: true, // targetValue starts with 'hello_'
    },
    {
      targetValue: 'hello_world',
      values: ['world'],
      expected: false, // targetValue does not start with 'world'
    },
  ];

  const clauseEvaluator = new ClauseEvaluator();

  testcases.forEach((tc, i) => {
    const clause = new Clause();
    clause.setOperator(Clause.Operator.STARTS_WITH); // Change operator to 'STARTS_WITH'
    clause.setValuesList(tc.values);

    const result = clauseEvaluator.evaluate(tc.targetValue, clause, 'userId', [], {});
    t.is(
      result,
      tc.expected,
      `Test case ${i} failed: targetValue ${tc.targetValue} : values ${tc.values}`,
    );
  });
});

// TestEndsWithOperator does not exist on the Go version testcases
test('TestEndsWithOperator', (t) => {
  const testcases = [
    // Simple ends with match
    {
      targetValue: '1519223320',
      values: ['3320'],
      expected: true, // targetValue ends with '3320'
    },
    {
      targetValue: '1519223320',
      values: ['3321'],
      expected: false, // targetValue does not end with '3321'
    },
    // Multiple possible suffixes
    {
      targetValue: '1519223320',
      values: ['3320', '1234', '5678'],
      expected: true, // targetValue ends with '3320'
    },
    {
      targetValue: '1519223320',
      values: ['1234', '5678'],
      expected: false, // targetValue does not end with any of these values
    },
    // Full match is also a valid "ends with"
    {
      targetValue: 'apple',
      values: ['apple'],
      expected: true, // targetValue fully matches, so it ends with 'apple'
    },
    // String comparisons
    {
      targetValue: 'banana',
      values: ['nana'],
      expected: true, // targetValue ends with 'nana'
    },
    {
      targetValue: 'banana',
      values: ['ana', 'na', 'nana'],
      expected: true, // targetValue ends with 'nana', 'ana', or 'na'
    },
    {
      targetValue: 'banana',
      values: ['car', 'dog'],
      expected: false, // targetValue does not end with any of these values
    },
    // Numbers as strings
    {
      targetValue: '123456',
      values: ['456'],
      expected: true, // targetValue ends with '456'
    },
    {
      targetValue: '123456',
      values: ['123'],
      expected: false, // targetValue does not end with '123'
    },
    // Mixed strings and numbers
    {
      targetValue: 'hello123',
      values: ['123'],
      expected: true, // targetValue ends with '123'
    },
    {
      targetValue: 'hello123',
      values: ['hello'],
      expected: false, // targetValue does not end with 'hello'
    },
    // Empty values list
    {
      targetValue: 'anything',
      values: [],
      expected: false, // An empty list means no suffix matches, so it should return false
    },
    // Case sensitivity
    {
      targetValue: 'Apple',
      values: ['apple'],
      expected: false, // Case-sensitive comparison, 'Apple' does not end with 'apple'
    },
    {
      targetValue: 'apple',
      values: ['apple'],
      expected: true, // Exact match, targetValue ends with 'apple'
    },
    // Special characters
    {
      targetValue: 'hello_world',
      values: ['_world'],
      expected: true, // targetValue ends with '_world'
    },
    {
      targetValue: 'hello_world',
      values: ['hello'],
      expected: false, // targetValue does not end with 'hello'
    },
  ];

  const clauseEvaluator = new ClauseEvaluator();

  testcases.forEach((tc, i) => {
    const clause = new Clause();
    clause.setOperator(Clause.Operator.ENDS_WITH); // Change operator to 'ENDS_WITH'
    clause.setValuesList(tc.values);

    const result = clauseEvaluator.evaluate(tc.targetValue, clause, 'userId', [], {});
    t.is(
      result,
      tc.expected,
      `Test case ${i} failed: targetValue ${tc.targetValue} : values ${tc.values}`,
    );
  });
});

// TestEndsWithOperator does not exist on the Go version testcases
test('TestPartiallyMatchesOperator', (t) => {
  const testcases = [
    // Simple partial match
    {
      targetValue: '1519223320',
      values: ['922'],
      expected: true, // targetValue contains '922'
    },
    {
      targetValue: '1519223320',
      values: ['333'],
      expected: false, // targetValue does not contain '333'
    },
    // Multiple possible matches
    {
      targetValue: '1519223320',
      values: ['333', '223', '920'],
      expected: true, // targetValue contains '223'
    },
    {
      targetValue: '1519223320',
      values: ['333', '222', '920'],
      expected: false, // targetValue does not contain any of these values
    },
    {
      targetValue: '1519223320',
      values: ['555', '666', '777'],
      expected: false, // targetValue does not contain any of these values
    },
    // Full match is valid as a partial match
    {
      targetValue: 'apple',
      values: ['apple'],
      expected: true, // targetValue contains 'apple' fully
    },
    // String comparisons
    {
      targetValue: 'banana',
      values: ['nan'],
      expected: true, // targetValue contains 'nan'
    },
    {
      targetValue: 'banana',
      values: ['ana', 'ban'],
      expected: true, // targetValue contains 'ana' and 'ban'
    },
    {
      targetValue: 'banana',
      values: ['car', 'dog'],
      expected: false, // targetValue does not contain 'car' or 'dog'
    },
    // Numbers as strings
    {
      targetValue: '123456',
      values: ['234'],
      expected: true, // targetValue contains '234'
    },
    {
      targetValue: '123456',
      values: ['789'],
      expected: false, // targetValue does not contain '789'
    },
    // Mixed strings and numbers
    {
      targetValue: 'hello123',
      values: ['123'],
      expected: true, // targetValue contains '123'
    },
    {
      targetValue: 'hello123',
      values: ['hello'],
      expected: true, // targetValue contains 'hello'
    },
    {
      targetValue: 'hello123',
      values: ['world'],
      expected: false, // targetValue does not contain 'world'
    },
    // Empty values list
    {
      targetValue: 'anything',
      values: [],
      expected: false, // An empty list means no partial match, so it should return false
    },
    // Case sensitivity
    {
      targetValue: 'Apple',
      values: ['apple'],
      expected: false, // Case-sensitive comparison, 'Apple' does not contain 'apple'
    },
    {
      targetValue: 'apple',
      values: ['apple'],
      expected: true, // Exact match is valid for partial match
    },
    // Special characters
    {
      targetValue: 'hello_world',
      values: ['_world'],
      expected: true, // targetValue contains '_world'
    },
    {
      targetValue: 'hello_world',
      values: ['hello'],
      expected: true, // targetValue contains 'hello'
    },
    {
      targetValue: 'hello_world',
      values: ['world'],
      expected: true, // targetValue contains 'world'
    },
    {
      targetValue: 'hello_world',
      values: ['!'],
      expected: false, // targetValue does not contain '!'
    },
  ];

  const clauseEvaluator = new ClauseEvaluator();

  testcases.forEach((tc, i) => {
    const clause = new Clause();
    clause.setOperator(Clause.Operator.PARTIALLY_MATCH); // Change operator to 'PARTIALLY_MATCH'
    clause.setValuesList(tc.values);

    const result = clauseEvaluator.evaluate(tc.targetValue, clause, 'userId', [], {});
    t.is(
      result,
      tc.expected,
      `Test case ${i} failed: targetValue ${tc.targetValue} : values ${tc.values}`,
    );
  });
});

// TestSegementMatchesOperator does not exist on the Go version testcases
test('TestSegementMatchesOperator', (t) => {
  const clauseEvaluator = new ClauseEvaluator();

  SegmentEvaluatorTestCases.forEach((tc, i) => {
    const clause = new Clause();
    clause.setOperator(Clause.Operator.SEGMENT);
    clause.setValuesList(tc.segmentIDs);
    const result = clauseEvaluator.evaluate('', clause, tc.userID, tc.segmentUsers, {});
    t.is(
      result,
      tc.expected,
      `Test case ${i} failed: userID ${tc.userID} : segmentUsers ${tc.segmentUsers}`,
    );
  });
});
