import test from 'ava';
import { Feature } from '../../proto/feature/feature_pb';
import { Evaluator } from '../../evaluation';
import { SegmentUser } from '../../proto/feature/segment_pb';
import { UserEvaluations } from '../../proto/feature/evaluation_pb';
import { createEvaluation, createPrerequisite, createReason, createUser } from '../../modelFactory';
import { Reason } from '../../proto/feature/reason_pb';
import { NewUserEvaluations } from '../../userEvaluation';
import { newTestFeature } from './evaluate_feature_test';

interface TestCase {
  desc: string;
  prevUEID: string;
  evaluatedAt: number;
  userAttributesUpdated: boolean;
  tag: string;
  createFeatures: () => Feature[];
  expectedEvals: UserEvaluations;
  expectedEvalFeatureIDs: string[];
  expectedError: Error | null;
}

function TestEvaluateFeaturesByEvaluatedAtCases() {
  const now = new Date();
  function getTimeAgo(seconds: number): number {
    return new Date(now.getTime() - seconds * 1000).getTime() / 1000;
  }

  const thirtyOneDaysAgo = getTimeAgo(31 * 24 * 60 * 60);
  const fiveMinutesAgo = getTimeAgo(5 * 60);
  const tenMinutesAgo = getTimeAgo(10 * 60);
  const tenMinutesAndNineSecondsAgo = getTimeAgo(609);
  const tenMinutesAndElevenSecondsAgo = getTimeAgo(611);
  const oneHourAgo = getTimeAgo(60 * 60);

  const patterns: TestCase[] = [
    {
      desc: 'success: evaluate all features since the previous UserEvaluationsID is empty',
      prevUEID: '',
      evaluatedAt: thirtyOneDaysAgo,
      userAttributesUpdated: false,
      tag: '',
      createFeatures: () => {
        const f1 = newTestFeature('feature1');
        f1.setUpdatedAt(fiveMinutesAgo);

        const f2 = newTestFeature('feature2');
        f2.setUpdatedAt(fiveMinutesAgo);

        const f3 = newTestFeature('feature3');
        f3.setUpdatedAt(fiveMinutesAgo);
        f3.setArchived(true);
        return [f1, f2, f3];
      },
      expectedEvals: NewUserEvaluations(
        'dummy',
        [
          createEvaluation(
            'feature1:1:user1',
            'feature1',
            1,
            'user1',
            'variation-B',
            'B',
            'Variation B',
            createReason('', Reason.Type.DEFAULT),
          ),
          createEvaluation(
            'feature2:1:user1',
            'feature2',
            1,
            'user1',
            'variation-B',
            'B',
            'Variation B',
            createReason('', Reason.Type.DEFAULT),
          ),
        ],
        ['feature3'],
        true,
      ),
      expectedEvalFeatureIDs: ['feature1', 'feature2'],
      expectedError: null,
    },
    {
      desc: 'success: evaluate all features since the previous evaluation was over a month ago',
      prevUEID: 'prevUEID',
      evaluatedAt: thirtyOneDaysAgo,
      userAttributesUpdated: false,
      tag: '',
      createFeatures: () => {
        const f1 = newTestFeature('feature1');
        f1.setUpdatedAt(fiveMinutesAgo);

        const f2 = newTestFeature('feature2');
        f2.setUpdatedAt(fiveMinutesAgo);

        const f3 = newTestFeature('feature3');
        f3.setUpdatedAt(fiveMinutesAgo);
        f3.setArchived(true);
        return [f1, f2, f3];
      },
      expectedEvals: NewUserEvaluations(
        'dummy',
        [
          createEvaluation(
            'feature1:1:user1',
            'feature1',
            1,
            'user1',
            'variation-B',
            'B',
            'Variation B',
            createReason('', Reason.Type.DEFAULT),
          ),
          createEvaluation(
            'feature2:1:user1',
            'feature2',
            1,
            'user1',
            'variation-B',
            'B',
            'Variation B',
            createReason('', Reason.Type.DEFAULT),
          ),
        ],
        ['feature3'],
        true,
      ),
      expectedEvalFeatureIDs: ['feature1', 'feature2'],
      expectedError: null,
    },
    {
      desc: 'success: evaluate all features since both feature flags and user attributes have not been updated (although the UEID has been updated)',
      prevUEID: 'prevUEID',
      evaluatedAt: tenMinutesAgo,
      userAttributesUpdated: false,
      tag: '',
      createFeatures: () => {
        const f1 = newTestFeature('feature-1');
        f1.setUpdatedAt(oneHourAgo);

        const f2 = newTestFeature('feature-2');
        f2.setUpdatedAt(oneHourAgo);

        const f3 = newTestFeature('feature-3');
        f3.setUpdatedAt(oneHourAgo);
        f3.setArchived(true);
        return [f1, f2, f3];
      },
      expectedEvals: NewUserEvaluations(
        'dummy',
        [
          createEvaluation(
            'feature-1:1:user-1',
            'feature-1',
            1,
            'user-1',
            'variation-B',
            'B',
            'Variation B',
            createReason('', Reason.Type.DEFAULT),
          ),
          createEvaluation(
            'feature-2:1:user-1',
            'feature-2',
            1,
            'user-1',
            'variation-B',
            'B',
            'Variation B',
            createReason('', Reason.Type.DEFAULT),
          ),
        ],
        ['feature-3'],
        true,
      ),
      expectedEvalFeatureIDs: ['feature-1', 'feature-2'],
      expectedError: null,
    },
    {
      desc: 'success: evaluate only features updated since the previous evaluations',
      prevUEID: 'prevUEID',
      evaluatedAt: tenMinutesAgo,
      userAttributesUpdated: false,
      tag: '',
      createFeatures: () => {
        const f1 = newTestFeature('feature1');
        f1.setUpdatedAt(fiveMinutesAgo);

        const f2 = newTestFeature('feature2');
        f2.setUpdatedAt(oneHourAgo);

        const f3 = newTestFeature('feature3');
        f3.setUpdatedAt(fiveMinutesAgo);
        f3.setArchived(true);

        const f4 = newTestFeature('feature4');
        f4.setUpdatedAt(oneHourAgo);
        f4.setArchived(true);
        return [f1, f2, f3, f4];
      },
      expectedEvals: NewUserEvaluations(
        'dummy',
        [
          createEvaluation(
            'feature1:1:user1',
            'feature1',
            1,
            'user1',
            'variation-B',
            'B',
            'Variation B',
            createReason('', Reason.Type.DEFAULT),
          ),
        ],
        ['feature3'],
        false,
      ),
      expectedEvalFeatureIDs: ['feature1'],
      expectedError: null,
    },
    {
      desc: 'success: check the adjustment seconds',
      prevUEID: 'prevUEID',
      evaluatedAt: tenMinutesAgo,
      userAttributesUpdated: false,
      tag: '',
      createFeatures: () => {
        const f1 = newTestFeature('feature1');
        f1.setUpdatedAt(tenMinutesAndNineSecondsAgo);

        const f2 = newTestFeature('feature2');
        f2.setUpdatedAt(tenMinutesAndElevenSecondsAgo);
        return [f1, f2];
      },
      expectedEvals: NewUserEvaluations(
        'dummy',
        [
          createEvaluation(
            'feature1:1:user1',
            'feature1',
            1,
            'user1',
            'variation-B',
            'B',
            'Variation B',
            createReason('', Reason.Type.DEFAULT),
          ),
        ],
        [],
        false,
      ),
      expectedEvalFeatureIDs: ['feature1'],
      expectedError: null,
    },
    {
      desc: 'success: evaluate only features has rules when user attributes updated',
      prevUEID: 'prevUEID',
      evaluatedAt: tenMinutesAgo,
      userAttributesUpdated: true,
      tag: '',
      createFeatures: () => {
        const f1 = newTestFeature('feature1');
        f1.setUpdatedAt(thirtyOneDaysAgo);

        const f2 = newTestFeature('feature2');
        f2.setUpdatedAt(thirtyOneDaysAgo);
        f2.setRulesList([]);

        const f3 = newTestFeature('feature3');
        f3.setUpdatedAt(thirtyOneDaysAgo);
        f3.setArchived(true);
        return [f1, f2, f3];
      },
      expectedEvals: NewUserEvaluations(
        'dummy',
        [
          createEvaluation(
            'feature1:1:user1',
            'feature1',
            1,
            'user1',
            'variation-B',
            'B',
            'Variation B',
            createReason('', Reason.Type.RULE),
          ),
        ],
        [],
        false,
      ),
      expectedEvalFeatureIDs: ['feature1'],
      expectedError: null,
    },
    {
      desc: 'success: evaluate only the features that have been updated since the previous evaluation, or the features that have rules when user attributes are updated',
      prevUEID: 'prevUEID',
      evaluatedAt: tenMinutesAgo,
      userAttributesUpdated: true,
      tag: '',
      createFeatures: () => {
        const f1 = newTestFeature('feature1');
        f1.setUpdatedAt(fiveMinutesAgo);
        f1.setRulesList([]);

        const f2 = newTestFeature('feature2');
        f2.setUpdatedAt(thirtyOneDaysAgo);
        f2.setRulesList([]);

        const f3 = newTestFeature('feature3');
        f3.setUpdatedAt(fiveMinutesAgo);
        f3.setArchived(true);

        const f4 = newTestFeature('feature4');
        f4.setUpdatedAt(fiveMinutesAgo);
        f4.setRulesList([]);
        return [f1, f2, f3, f4];
      },
      expectedEvals: NewUserEvaluations(
        'dummy',
        [
          createEvaluation(
            'feature1:1:user1',
            'feature1',
            1,
            'user1',
            'variation-B',
            'B',
            'Variation B',
            createReason('', Reason.Type.RULE),
          ),
          createEvaluation(
            'feature4:1:user1',
            'feature4',
            1,
            'user1',
            'variation-B',
            'B',
            'Variation B',
            createReason('', Reason.Type.RULE),
          ),
        ],
        ['feature3'],
        false,
      ),
      expectedEvalFeatureIDs: ['feature1', 'feature4'],
      expectedError: null,
    },
    {
      desc: 'success: prerequisite',
      prevUEID: 'prevUEID',
      evaluatedAt: tenMinutesAgo,
      userAttributesUpdated: false,
      tag: '',
      createFeatures: () => {
        const f1 = newTestFeature('feature1');
        f1.setUpdatedAt(thirtyOneDaysAgo);
        f1.setPrerequisitesList([createPrerequisite('feature4', 'B')]);

        const f2 = newTestFeature('feature2');
        f2.setUpdatedAt(thirtyOneDaysAgo);

        const f3 = newTestFeature('feature3');
        f3.setUpdatedAt(thirtyOneDaysAgo);

        const f4 = newTestFeature('feature4');
        f4.setUpdatedAt(fiveMinutesAgo);
        return [f1, f2, f3, f4];
      },
      expectedEvals: NewUserEvaluations(
        'dummy',
        [
          createEvaluation(
            'feature1:1:user1',
            'feature1',
            1,
            'user1',
            'variation-B',
            'B',
            'Variation B',
            createReason('', Reason.Type.RULE),
          ),
          createEvaluation(
            'feature4:1:user1',
            'feature4',
            1,
            'user1',
            'variation-B',
            'B',
            'Variation B',
            createReason('', Reason.Type.RULE),
          ),
        ],
        [],
        false,
      ),
      expectedEvalFeatureIDs: ['feature1', 'feature4'],
      expectedError: null,
    },
    {
      desc: "success: When a tag is specified, it excludes the evaluations that don't have that tag. But archived features are not excluded",
      prevUEID: 'prevUEID',
      evaluatedAt: tenMinutesAgo,
      userAttributesUpdated: false,
      tag: 'tag1',
      createFeatures: () => {
        const f1 = newTestFeature('feature1');
        f1.setTagsList(['tag1']);
        f1.setUpdatedAt(fiveMinutesAgo);

        const f2 = newTestFeature('feature2');
        f2.setTagsList(['tag2']);
        f2.setUpdatedAt(fiveMinutesAgo);

        const f3 = newTestFeature('feature3');
        f3.setTagsList(['tag1']);
        f3.setArchived(true);
        f3.setUpdatedAt(fiveMinutesAgo);

        const f4 = newTestFeature('feature4');
        f4.setTagsList(['tag2']);
        f4.setArchived(true);
        f4.setUpdatedAt(fiveMinutesAgo);
        return [f1, f2, f3, f4];
      },
      expectedEvals: NewUserEvaluations(
        'dummy',
        [
          createEvaluation(
            'feature1:1:user1',
            'feature1',
            1,
            'user1',
            'variation-B',
            'B',
            'Variation B',
            createReason('', Reason.Type.DEFAULT),
          ),
        ],
        ['feature3', 'feature4'],
        false,
      ),
      expectedEvalFeatureIDs: ['feature1'],
      expectedError: null,
    },
    {
      desc: 'success: When a tag is not specified, it does not exclude evaluations that have tags.',
      prevUEID: 'prevUEID',
      evaluatedAt: tenMinutesAgo,
      userAttributesUpdated: false,
      tag: '',
      createFeatures: () => {
        const f1 = newTestFeature('feature1');
        f1.setTagsList(['tag1']);
        f1.setUpdatedAt(fiveMinutesAgo);

        const f2 = newTestFeature('feature2');
        f2.setTagsList(['tag2']);
        f2.setUpdatedAt(fiveMinutesAgo);

        const f3 = newTestFeature('feature3');
        f3.setUpdatedAt(fiveMinutesAgo);

        const f4 = newTestFeature('feature4');
        f4.setTagsList(['tag1', 'tag2']);
        f4.setUpdatedAt(fiveMinutesAgo);
        return [f1, f2, f3, f4];
      },
      expectedEvals: NewUserEvaluations(
        'dummy',
        [
          createEvaluation(
            'feature1:1:user1',
            'feature1',
            1,
            'user1',
            'variation-B',
            'B',
            'Variation B',
            createReason('', Reason.Type.DEFAULT),
          ),
          createEvaluation(
            'feature2:1:user1',
            'feature2',
            1,
            'user1',
            'variation-B',
            'B',
            'Variation B',
            createReason('', Reason.Type.DEFAULT),
          ),
          createEvaluation(
            'feature3:1:user1',
            'feature3',
            1,
            'user1',
            'variation-B',
            'B',
            'Variation B',
            createReason('', Reason.Type.DEFAULT),
          ),
          createEvaluation(
            'feature4:1:user1',
            'feature4',
            1,
            'user1',
            'variation-B',
            'B',
            'Variation B',
            createReason('', Reason.Type.DEFAULT),
          ),
        ],
        [],
        false,
      ),
      expectedEvalFeatureIDs: ['feature1', 'feature2', 'feature3', 'feature4'],
      expectedError: null,
    },
    {
      desc: 'success: including up/downwards features of target feature with prerequisite',
      prevUEID: 'prevUEID',
      evaluatedAt: tenMinutesAgo,
      userAttributesUpdated: false,
      tag: '',
      createFeatures: () => {
        const f1 = newTestFeature('feature1');
        f1.setUpdatedAt(oneHourAgo);
        f1.setPrerequisitesList([createPrerequisite('feature2', 'B')]);

        const f2 = newTestFeature('feature2');
        f2.setUpdatedAt(fiveMinutesAgo);
        f2.setPrerequisitesList([createPrerequisite('feature3', 'B')]);

        const f3 = newTestFeature('feature3');
        f3.setUpdatedAt(oneHourAgo);
        return [f1, f2, f3];
      },
      expectedEvals: NewUserEvaluations(
        'dummy',
        [
          createEvaluation(
            'feature1:1:user1',
            'feature1',
            1,
            'user1',
            'variation-B',
            'B',
            'Variation B',
            createReason('', Reason.Type.RULE),
          ),
          createEvaluation(
            'feature2:1:user1',
            'feature2',
            1,
            'user1',
            'variation-B',
            'B',
            'Variation B',
            createReason('', Reason.Type.RULE),
          ),
          createEvaluation(
            'feature3:1:user1',
            'feature3',
            1,
            'user1',
            'variation-B',
            'B',
            'Variation B',
            createReason('', Reason.Type.RULE),
          ),
        ],
        [],
        false,
      ),
      expectedEvalFeatureIDs: ['feature1', 'feature2', 'feature3'],
      expectedError: null,
    },
  ];
  return patterns;
}

TestEvaluateFeaturesByEvaluatedAtCases().forEach((p) => {
  test(p.desc, async (t) => {
    const user = createUser('user-1', {});
    const segmentUser: Map<string, SegmentUser[]> = new Map<string, SegmentUser[]>();
    const evaluator = new Evaluator();
    try {
      const actual = evaluator.evaluateFeaturesByEvaluatedAt(
        p.createFeatures(),
        user,
        segmentUser,
        p.prevUEID,
        p.evaluatedAt,
        p.userAttributesUpdated,
        p.tag,
      );

      t.deepEqual(p.expectedEvals.getForceUpdate(), actual.getForceUpdate());
      t.deepEqual(p.expectedEvals.getArchivedFeatureIdsList(), actual.getArchivedFeatureIdsList());
      t.is(p.expectedEvals.getEvaluationsList().length, actual.getEvaluationsList().length);
      actual.getEvaluationsList().forEach((e) => {
        t.true(p.expectedEvalFeatureIDs.includes(e.getFeatureId()));
      });
      //TODO: Check me - did Golang test is missing this
      //t.deepEqual(p.expectedEvals.toObject(), actual.toObject());
    } catch (error) {
      t.deepEqual(p.expectedError, error);
    }
  });
});
