import test from 'ava';
import { Feature } from '../../proto/feature/feature_pb';
import { EvaluationID, Evaluator } from '../../evaluation';
import { SegmentUser } from '../../proto/feature/segment_pb';
import { Evaluation } from '../../proto/feature/evaluation_pb';
import { createPrerequisite, createVariation, createFeature, createUser } from '../../modelFactory';
import { Strategy } from '../../proto/feature/strategy_pb';
import { Clause } from '../../proto/feature/clause_pb';
import { Reason } from '../../proto/feature/reason_pb';

const ErrVariationNotFound = new Error('evaluator: variation not found');

const TestFeatureIDs = {
  fID0: 'fID0',
  fID1: 'fID1',
  fID2: 'fID2',
};

export const TestVariations = {
  variationA: createVariation('variation-A', 'A', 'Variation A', 'Thing does A'),
  variationB: createVariation('variation-B', 'B', 'Variation B', 'Thing does B'),
  variationC: createVariation('variation-C', 'C', 'Variation C', 'Thing does C'),
};

export function newTestFeature(id: string): Feature {
  // Variations
  const variations = [
    TestVariations.variationA.toObject(),
    TestVariations.variationB.toObject(),
    TestVariations.variationC.toObject(),
  ];

  // Targets
  const targets = [
    { variation: 'variation-A', users: ['user1'] },
    { variation: 'variation-B', users: ['user2'] },
    { variation: 'variation-C', users: ['user3'] },
  ];

  // Rules
  const rules = [
    {
      id: 'rule-1',
      attribute: 'name',
      operator: Clause.Operator.EQUALS,
      values: ['user1', 'user2'],
      fixedVariation: 'variation-A',
    },
    {
      id: 'rule-2',
      attribute: 'name',
      operator: Clause.Operator.EQUALS,
      values: ['user3', 'user4'],
      fixedVariation: 'variation-B',
    },
  ];

  // Default Strategy
  const defaultStrategy = {
    type: Strategy.Type.FIXED,
    variation: 'variation-B',
  };

  // Call the second function to create and return the Feature
  return createFeature({
    id: id,
    name: 'test feature',
    version: 1,
    enabled: true,
    createdAt: Date.now(),
    variationType: Feature.VariationType.STRING,
    variations: variations,
    targets: targets,
    rules: rules,
    defaultStrategy: defaultStrategy,
    prerequisitesList: [],
  });
}

const findEvaluation = (evaluations: Evaluation[], featureId: string): Evaluation | null => {
  return evaluations.find((e) => e.getFeatureId() === featureId) || null;
};

test('EvaluateFeature', async (t) => {
  const f = newTestFeature(TestFeatureIDs.fID0);
  f.getTagsList().push('tag1');

  const f1 = newTestFeature(TestFeatureIDs.fID1);
  f1.getTagsList().push('tag1');
  f1.setEnabled(false);
  f1.setOffVariation(TestVariations.variationA.getId());

  const f2 = newTestFeature(TestFeatureIDs.fID2);
  f2.getTagsList().push('tag1');

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
        variation: { id: 'variation-A', name: 'Variation A', value: 'A', description: '' },
        reason: { ruleId: '', type: Reason.Type.OFF_VARIATION },
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
        variation: { id: 'variation-B', name: 'Variation B', value: 'B', description: '' },
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
        variation: { id: 'variation-B', name: 'Variation B', value: 'B', description: '' },
        reason: { ruleId: '', type: Reason.Type.DEFAULT },
      },
      expectedError: null,
    },
    {
      enabled: true,
      offVariation: 'variation-A',
      userID: 'uID2',
      prerequisite: [createPrerequisite(f1.getId(), TestVariations.variationB.getId())],
      expected: {
        id: EvaluationID(f.getId(), f.getVersion(), 'uID2'),
        featureId: 'fID0',
        featureVersion: 1,
        userId: 'uID2',
        variationId: 'variation-A',
        variationName: 'Variation A',
        variationValue: 'A',
        variation: { id: 'variation-A', name: 'Variation A', value: 'A', description: '' },
        reason: { ruleId: '', type: Reason.Type.PREREQUISITE },
      },
      expectedError: null,
    },
    {
      enabled: true,
      offVariation: '',
      userID: 'uID2',
      prerequisite: [createPrerequisite(f2.getId(), TestVariations.variationA.getId())],
      expected: {
        id: EvaluationID(f.getId(), f.getVersion(), 'uID2'),
        featureId: 'fID0',
        featureVersion: 1,
        userId: 'uID2',
        variationId: 'variation-B',
        variationName: 'Variation B',
        variationValue: 'B',
        variation: { id: 'variation-B', name: 'Variation B', value: 'B', description: '' },
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
    } catch (error) {
      if (error instanceof Error || error === null) {
        t.deepEqual(p.expectedError, error);
      } else {
        t.fail(`Unexpected error type: ${typeof error}: ${error}`);
      }
    }
  }
});
