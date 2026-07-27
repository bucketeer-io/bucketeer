import test from 'ava';
import * as fs from 'fs';
import * as path from 'path';
import { Segment, SegmentUser } from '../proto/feature/segment_pb';
import { Rule } from '../proto/feature/rule_pb';
import { Clause } from '../proto/feature/clause_pb';
import { SegmentEvaluator } from '../segmentEvaluator';
import { createUser } from '../modelFactory';

// The conformance fixtures are shared with evaluation/go so the two engines
// stay in lockstep. See evaluation/testdata/segment_rules_conformance.json.
interface ConformanceClause {
  id: string;
  attribute: string;
  operator: string;
  values: string[];
}

interface ConformanceRule {
  id: string;
  clauses: ConformanceClause[];
}

interface ConformanceSegment {
  id: string;
  rules: ConformanceRule[];
}

interface ConformanceSegmentUser {
  segmentId: string;
  userId: string;
  state: string;
}

interface ConformanceTestCase {
  desc: string;
  segmentIds: string[];
  user: { id: string; data: { [key: string]: string } };
  expected: boolean;
}

interface ConformanceFixture {
  segments: ConformanceSegment[];
  segmentUsers: ConformanceSegmentUser[];
  testCases: ConformanceTestCase[];
}

function loadFixture(): ConformanceFixture {
  // Compiled tests run from __test/__tests__, source runs from src/__tests__.
  const candidates = [
    path.join(__dirname, '../../../testdata/segment_rules_conformance.json'),
    path.join(__dirname, '../../testdata/segment_rules_conformance.json'),
  ];
  for (const candidate of candidates) {
    if (fs.existsSync(candidate)) {
      return JSON.parse(fs.readFileSync(candidate, 'utf-8'));
    }
  }
  throw new Error('segment_rules_conformance.json not found');
}

function toOperator(name: string): Clause.OperatorMap[keyof Clause.OperatorMap] {
  const operator = Clause.Operator[name as keyof Clause.OperatorMap];
  if (operator === undefined) {
    throw new Error(`unknown clause operator: ${name}`);
  }
  return operator;
}

function toState(name: string): SegmentUser.StateMap[keyof SegmentUser.StateMap] {
  const state = SegmentUser.State[name as keyof SegmentUser.StateMap];
  if (state === undefined) {
    throw new Error(`unknown segment user state: ${name}`);
  }
  return state;
}

function toProtoRules(rules: ConformanceRule[]): Rule[] {
  return rules.map((r) => {
    const rule = new Rule();
    rule.setId(r.id);
    rule.setClausesList(
      r.clauses.map((c) => {
        const clause = new Clause();
        clause.setId(c.id);
        clause.setAttribute(c.attribute);
        clause.setOperator(toOperator(c.operator));
        clause.setValuesList(c.values);
        return clause;
      }),
    );
    return rule;
  });
}

const fixture = loadFixture();

const segments = new Map<string, Segment>();
fixture.segments.forEach((s) => {
  const segment = new Segment();
  segment.setId(s.id);
  segment.setRulesList(toProtoRules(s.rules));
  segments.set(s.id, segment);
});

const segmentUsers = fixture.segmentUsers.map((su) => {
  const segmentUser = new SegmentUser();
  segmentUser.setSegmentId(su.segmentId);
  segmentUser.setUserId(su.userId);
  segmentUser.setState(toState(su.state));
  return segmentUser;
});

fixture.testCases.forEach((tc) => {
  test(`conformance: ${tc.desc}`, (t) => {
    const evaluator = new SegmentEvaluator();
    const user = createUser(tc.user.id, tc.user.data);
    const actual = evaluator.evaluate(tc.segmentIds, user, segments, segmentUsers);
    t.is(actual, tc.expected);
  });
});
