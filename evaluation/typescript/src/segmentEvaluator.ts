import { Segment, SegmentUser } from './proto/feature/segment_pb';
import { User } from './proto/user/user_pb';
import { RuleEvaluator } from './ruleEvaluator';
//
class SegmentEvaluator {
  // evaluate reports whether the user belongs to ANY of the given segments (OR).
  evaluate(
    segmentIDs: string[],
    user: User,
    segments: Map<string, Segment> | null,
    segmentUsers: SegmentUser[],
  ): boolean {
    for (const segmentID of segmentIDs) {
      const segment = segments?.get(segmentID) ?? null;
      if (this.isUserInSegment(segmentID, user, segment, segmentUsers)) {
        return true;
      }
    }
    return false;
  }

  // isUserInSegment reports whether the user belongs to the segment:
  // the user is in the included-user list OR matches any of the segment rules.
  private isUserInSegment(
    segmentID: string,
    user: User,
    segment: Segment | null,
    segmentUsers: SegmentUser[],
  ): boolean {
    // 1. Explicit include list — existing behavior, unchanged.
    if (
      this.containsSegmentUser(segmentID, user.getId(), SegmentUser.State.INCLUDED, segmentUsers)
    ) {
      return true;
    }
    // 2. Rules. Segment rules use the same Rule message as feature flag rules,
    // so the flag-side rule evaluator (OR across rules, AND across clauses,
    // attribute resolution, all operators) runs on them as is.
    // segmentUsers/segments/flagVariations are empty/null: validation rejects
    // SEGMENT and FEATURE_FLAG operators inside segment rules, and null makes
    // any such clause fail closed.
    if (!segment || segment.getRulesList().length === 0) {
      return false;
    }
    const ruleEvaluator = new RuleEvaluator();
    const matchedRule = ruleEvaluator.evaluate(segment.getRulesList(), user, [], null, null);
    return matchedRule !== null;
  }

  private containsSegmentUser(
    segmentID: string,
    userID: string,
    state: SegmentUser.StateMap[keyof SegmentUser.StateMap],
    segmentUsers: SegmentUser[],
  ): boolean {
    for (const user of segmentUsers) {
      if (user.getSegmentId() !== segmentID) {
        continue;
      }
      if (user.getUserId() !== userID) {
        continue;
      }
      if (user.getState() !== state) {
        continue;
      }
      return true;
    }
    return false;
  }
}

export { SegmentEvaluator };
