import { SegmentUser } from './proto/feature/segment_pb';
//
class SegmentEvaluator {
  evaluate(segmentIDs: string[], userID: string, segmentUsers: SegmentUser[]): boolean {
    return this.findSegmentUser(segmentIDs, userID, SegmentUser.State.INCLUDED, segmentUsers);
  }

  private findSegmentUser(
    segmentIDs: string[],
    userID: string,
    state: SegmentUser.StateMap[keyof SegmentUser.StateMap],
    segmentUsers: SegmentUser[],
  ): boolean {
    for (const segmentID of segmentIDs) {
      if (!this.containsSegmentUser(segmentID, userID, state, segmentUsers)) {
        return false;
      }
    }
    return true;
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
