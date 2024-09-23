import test from 'ava';
import { Feature } from '../proto/feature/feature_pb';
import { EvaluationID, Evaluator } from '../evaluation';
import { SegmentUser } from '../proto/feature/segment_pb';
import { Evaluation, UserEvaluations } from '../proto/feature/evaluation_pb';
import { createUser } from './rule_evaluator_test';
import { createEvaluation, createPrerequisite, createReason, createVariation, creatFeature } from './utils/test_data';
import { Strategy } from '../proto/feature/strategy_pb';
import { Clause } from '../proto/feature/clause_pb';
import { Reason } from '../proto/feature/reason_pb';
import { NewUserEvaluations } from '../userEvaluation';

const ErrVariationNotFound = new Error('evaluator: variation not found');

const TestFeatureIDs = {
  fID0: 'fID0',
  fID1: 'fID1',
  fID2: 'fID2',
}

const TestVariations = {
  variationA: createVariation('variation-A', 'A', 'Variation A', 'Thing does A'),
  variationB: createVariation('variation-B', 'B', 'Variation B', 'Thing does B')
}

function newTestFeature(id: string): Feature {
  // Variations
  const variations = [
    { id: 'variation-A', value: 'A', name: 'Variation A', description: 'Thing does A' },
    { id: 'variation-B', value: 'B', name: 'Variation B', description: 'Thing does B' },
    { id: 'variation-C', value: 'C', name: 'Variation C', description: 'Thing does C' }
  ];

  // Targets
  const targets = [
    { variation: 'variation-A', users: ['user1'] },
    { variation: 'variation-B', users: ['user2'] },
    { variation: 'variation-C', users: ['user3'] }
  ];

  // Rules
  const rules = [
    {
      id: 'rule-1',
      attribute: 'name',
      operator: Clause.Operator.EQUALS,
      values: ['user1', 'user2'],
      fixedVariation: 'variation-A'
    },
    {
      id: 'rule-2',
      attribute: 'name',
      operator: Clause.Operator.EQUALS,
      values: ['user3', 'user4'],
      fixedVariation: 'variation-B'
    }
  ];

  // Default Strategy
  const defaultStrategy = {
    type: Strategy.Type.FIXED,
    variation: 'variation-B'
  };

  // Call the second function to create and return the Feature
  return creatFeature(
    id,
    'test feature',
    1,
    true,
    Date.now(),
    Feature.VariationType.STRING,
    variations,
    targets,
    rules,
    defaultStrategy
  );
}

const findEvaluation = (evaluations: Evaluation[], featureId: string): Evaluation | null => {
  return evaluations.find(e => e.getFeatureId() === featureId) || null;
};

test('EvaluateFeature', async t => {
  const f = newTestFeature(TestFeatureIDs.fID0);
  f.getTagsList().push('tag1')

  const f1 = newTestFeature(TestFeatureIDs.fID1);
  f1.getTagsList().push('tag1');
  f1.setEnabled(false);
  f1.setOffVariation(TestVariations.variationA.getId());

  const f2 = newTestFeature(TestFeatureIDs.fID2);
  f2.getTagsList().push('tag1')

  const patterns = [
    {
      enabled: false,
      offVariation: 'notfound',
      userID: 'uID0',
      prerequisite: [],
      expected: null,
      expectedError: ErrVariationNotFound,
    },
    {
      enabled: false,
      offVariation: 'variation-A',
      userID: 'uID0',
      prerequisite: [],
      expected: {
        featureId: 'fID0',
        featureVersion: 1,
        id: EvaluationID(f.getId(), f.getVersion(), 'uID0'),
        userId: 'uID0',
        variationId: 'variation-A',
        variationName: 'Variation A',
        variationValue: 'A',
        variation: { id: 'variation-A', name: 'Variation A', value: 'A', description: '', },
        reason: { ruleId: '', type: Reason.Type.OFF_VARIATION, },
      },
      expectedError: null,
    },
    {
      enabled: false,
      offVariation: '',
      userID: 'uID0',
      prerequisite: [],
      expected: {
        id: EvaluationID(f.getId(), f.getVersion(), 'uID0'),
        featureId: 'fID0',
        featureVersion: 1,
        userId: 'uID0',
        variationId: 'variation-B',
        variationName: 'Variation B',
        variationValue: 'B',
        variation: { id: 'variation-B', name: 'Variation B', value: 'B', description: '', },
        reason: { ruleId: '', type: Reason.Type.DEFAULT },
      },
      expectedError: null,
    },
    {
      enabled: true,
      offVariation: '',
      userID: 'uID2',
      prerequisite: [],
      expected: {
        id: EvaluationID(f.getId(), f.getVersion(), 'uID2'),
        featureId: 'fID0',
        featureVersion: 1,
        userId: 'uID2',
        variationId: 'variation-B',
        variationName: 'Variation B',
        variationValue: 'B',
        variation: { id: 'variation-B', name: 'Variation B', value: 'B', description: '', },
        reason: { ruleId: '', type: Reason.Type.DEFAULT },
      },
      expectedError: null,
    },
    {
      enabled: true,
      offVariation: 'variation-A',
      userID: 'uID2',
      prerequisite: [
        createPrerequisite(f1.getId(), TestVariations.variationB.getId())
      ],
      expected: {
        id: EvaluationID(f.getId(), f.getVersion(), 'uID2'),
        featureId: 'fID0',
        featureVersion: 1,
        userId: 'uID2',
        variationId: 'variation-A',
        variationName: 'Variation A',
        variationValue: 'A',
        variation: { id: 'variation-A', name: 'Variation A', value: 'A', description: '', },
        reason: { ruleId: '', type: Reason.Type.PREREQUISITE },
      },
      expectedError: null,
    },
    {
      enabled: true,
      offVariation: '',
      userID: 'uID2',
      prerequisite: [
        createPrerequisite(f2.getId(), TestVariations.variationA.getId())
      ],
      expected: {
        id: EvaluationID(f.getId(), f.getVersion(), 'uID2'),
        featureId: 'fID0',
        featureVersion: 1,
        userId: 'uID2',
        variationId: 'variation-B',
        variationName: 'Variation B',
        variationValue: 'B',
        variation: { id: 'variation-B', name: 'Variation B', value: 'B',  description: '', },
        reason: { ruleId: '', type: Reason.Type.DEFAULT },
      },
      expectedError: null,
    },
  ];

  for (const p of patterns) {
    const evaluator = new Evaluator();
    const user = createUser(p.userID, {});
    f.setEnabled(p.enabled);
    f.setOffVariation(p.offVariation);
    f.setPrerequisitesList(p.prerequisite);
    
    const segmentUser: Map<string, SegmentUser[]> = new Map<string, SegmentUser[]>();
    try {
      const evaluation = await evaluator.evaluateFeatures([f, f1, f2], user, segmentUser, 'tag1');
      if (evaluation.getEvaluationsList()) {
        const actual = findEvaluation(evaluation.getEvaluationsList(), f.getId());
        t.deepEqual(p.expected, actual?.toObject());
      }
    } catch(error) {
      if (error instanceof Error || error === null) {
        t.deepEqual(p.expectedError, error);
      } else {
        t.fail(`Unexpected error type: ${typeof error}: ${error}`);
      }
    }
  }
});

test('TestEvaluateFeaturesByEvaluatedAt', async t => {
  const now = new Date();
  function getTimeAgo(seconds: number): number {
    return (new Date(now.getTime() - seconds * 1000)).getTime() / 1000;
  }

  const thirtyOneDaysAgo = getTimeAgo(31 * 24 * 60 * 60);
  const fiveMinutesAgo = getTimeAgo(5 * 60);
  // const tenMinutesAgo = getTimeAgo(10 * 60);
  // const tenMinutesAndNineSecondsAgo = getTimeAgo(609);
  // const tenMinutesAndElevenSecondsAgo = getTimeAgo(611);
  // const oneHourAgo = getTimeAgo(60 * 60);
  const user = createUser ('user-1', {});
  const segmentUser: Map<string, SegmentUser[]> = new Map<string, SegmentUser[]>();
  
  interface Pattern {
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
  
  const patterns: Pattern[] = [
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
            createReason( '', Reason.Type.DEFAULT),
          ),
          createEvaluation(
            'feature2:1:user1',
            'feature2',
            1,
            'user1',
            'variation-B',
            'B',
            'Variation B',
            createReason( '', Reason.Type.DEFAULT),
          ),
        ],
        ['feature3'],
        true
      ),
      expectedEvalFeatureIDs: ['feature1', 'feature2'],
      expectedError: null,
    },
  ];

  patterns.forEach(p => {
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
        actual.getEvaluationsList().forEach(e => {
          t.true(p.expectedEvalFeatureIDs.includes(e.getFeatureId()));
        });
        //TODO: Check me - did Golang test is missing this
        //t.deepEqual(p.expectedEvals.toObject(), actual.toObject());
      } catch (error) {
        t.deepEqual(p.expectedError, error);
      }
  });
});