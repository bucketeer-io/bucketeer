import test from 'ava';
import { SegmentUser } from '../proto/feature/segment_pb';
import { SegmentEvaluator } from '../segmentEvaluator';

// Define the type for the test cases
interface SegmentEvaluatorTestCase {
  desc: string;
  segmentIDs: string[];
  userID: string;
  segmentUsers: SegmentUser[];
  expected: boolean;
}

export const SegmentEvaluatorTestCases: SegmentEvaluatorTestCase[] = [
  {
    desc: 'user is included in all segments',
    segmentIDs: ['segment-1', 'segment-2'],
    userID: 'user-1',
    segmentUsers: [
      createSegmentUser('segment-1', 'user-1', SegmentUser.State.INCLUDED),
      createSegmentUser('segment-2', 'user-1', SegmentUser.State.INCLUDED),
    ],
    expected: true,
  },
  {
    desc: 'user is excluded in one segment',
    segmentIDs: ['segment-1', 'segment-2'],
    userID: 'user-1',
    segmentUsers: [
      createSegmentUser('segment-1', 'user-1', SegmentUser.State.INCLUDED),
      createSegmentUser('segment-2', 'user-1', SegmentUser.State.EXCLUDED),
    ],
    expected: false,
  },
  {
    desc: 'user does not exist in any segments',
    segmentIDs: ['segment-1', 'segment-2'],
    userID: 'user-1',
    segmentUsers: [
      createSegmentUser('segment-1', 'user-2', SegmentUser.State.INCLUDED),
      createSegmentUser('segment-2', 'user-2', SegmentUser.State.INCLUDED),
    ],
    expected: false,
  },
  {
    desc: 'empty segment IDs',
    segmentIDs: [],
    userID: 'user-1',
    segmentUsers: [createSegmentUser('segment-1', 'user-1', SegmentUser.State.INCLUDED)],
    expected: true, // No segments to evaluate means always true
  },
  {
    desc: 'single segment ID, user included',
    segmentIDs: ['segment-1'],
    userID: 'user-1',
    segmentUsers: [createSegmentUser('segment-1', 'user-1', SegmentUser.State.INCLUDED)],
    expected: true,
  },
  {
    desc: 'single segment ID, user excluded',
    segmentIDs: ['segment-1'],
    userID: 'user-1',
    segmentUsers: [createSegmentUser('segment-1', 'user-1', SegmentUser.State.EXCLUDED)],
    expected: false,
  },
  {
    desc: 'multiple segment IDs with mixed states',
    segmentIDs: ['segment-1', 'segment-2'],
    userID: 'user-1',
    segmentUsers: [
      createSegmentUser('segment-1', 'user-1', SegmentUser.State.INCLUDED),
      createSegmentUser('segment-2', 'user-1', SegmentUser.State.INCLUDED),
    ],
    expected: true,
  },
  {
    desc: 'multiple segment IDs, only one excluded',
    segmentIDs: ['segment-1', 'segment-2', 'segment-3'],
    userID: 'user-1',
    segmentUsers: [
      createSegmentUser('segment-1', 'user-1', SegmentUser.State.INCLUDED),
      createSegmentUser('segment-2', 'user-1', SegmentUser.State.EXCLUDED),
      createSegmentUser('segment-3', 'user-1', SegmentUser.State.INCLUDED),
    ],
    expected: false,
  },
  {
    desc: 'user included in segments, but not all segments defined',
    segmentIDs: ['segment-1', 'segment-2'],
    userID: 'user-1',
    segmentUsers: [
      createSegmentUser('segment-1', 'user-1', SegmentUser.State.INCLUDED),
      createSegmentUser('segment-3', 'user-1', SegmentUser.State.INCLUDED), // segment-3 is not in the IDs
    ],
    expected: false,
  },
  {
    desc: 'user included in all segments but one',
    segmentIDs: ['segment-1', 'segment-2'],
    userID: 'user-1',
    segmentUsers: [
      createSegmentUser('segment-1', 'user-1', SegmentUser.State.INCLUDED),
      createSegmentUser('segment-2', 'user-1', SegmentUser.State.EXCLUDED), // Excluded in segment-2
    ],
    expected: false,
  },
  {
    desc: 'multiple users with mixed states across segments',
    segmentIDs: ['segment-1', 'segment-2'],
    userID: 'user-2',
    segmentUsers: [
      createSegmentUser('segment-1', 'user-1', SegmentUser.State.INCLUDED),
      createSegmentUser('segment-1', 'user-2', SegmentUser.State.INCLUDED),
      createSegmentUser('segment-2', 'user-2', SegmentUser.State.INCLUDED),
    ],
    expected: true,
  },
  {
    desc: 'user present in segments but marked as deleted',
    segmentIDs: ['segment-1', 'segment-2'],
    userID: 'user-1',
    segmentUsers: [
      createSegmentUser('segment-1', 'user-1', SegmentUser.State.INCLUDED),
      createSegmentUser('segment-2', 'user-1', SegmentUser.State.EXCLUDED),
    ],
    expected: false, // The user is excluded in segment-2
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

SegmentEvaluatorTestCases.forEach(({ desc, segmentIDs, userID, segmentUsers, expected }) => {
  test(desc, (t) => {
    const evaluator = new SegmentEvaluator();
    const actual = evaluator.evaluate(segmentIDs, userID, segmentUsers);
    t.is(actual, expected);
  });
});
