import { Rule } from './proto/feature/rule_pb';
import { Clause } from './proto/feature/clause_pb';
import { User } from './proto/user/user_pb';
import { Segment, SegmentUser } from './proto/feature/segment_pb';
import { ClauseEvaluator } from './clauseEvaluator';
//
class RuleEvaluator {
  private clauseEvaluator: ClauseEvaluator;

  constructor() {
    this.clauseEvaluator = new ClauseEvaluator();
  }

  evaluate(
    rules: Rule[],
    user: User,
    segmentUsers: SegmentUser[],
    segments: Map<string, Segment> | null,
    flagVariations: { [key: string]: string } | null,
  ): Rule | null {
    for (const rule of rules) {
      const matched = this.evaluateRule(rule, user, segmentUsers, segments, flagVariations);
      if (matched) {
        return rule;
      }
    }
    return null;
  }

  private evaluateRule(
    rule: Rule,
    user: User,
    segmentUsers: SegmentUser[],
    segments: Map<string, Segment> | null,
    flagVariations: { [key: string]: string } | null,
  ): boolean {
    for (const clause of rule.getClausesList()) {
      const matched = this.evaluateClause(clause, user, segmentUsers, segments, flagVariations);
      if (!matched) {
        return false;
      }
    }
    return true;
  }

  private evaluateClause(
    clause: Clause,
    user: User,
    segmentUsers: SegmentUser[],
    segments: Map<string, Segment> | null,
    flagVariations: { [key: string]: string } | null,
  ): boolean {
    let targetAttr: string | undefined;
    if (clause.getAttribute() === 'id') {
      targetAttr = user.getId();
    } else {
      targetAttr = user.getDataMap().get(clause.getAttribute());
    }

    return this.clauseEvaluator.evaluate(
      targetAttr || '', // Handling the case where targetAttr is undefined
      clause,
      user,
      segmentUsers,
      segments,
      flagVariations,
    );
  }
}

export { RuleEvaluator };
