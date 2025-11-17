import test from 'ava';
import { Feature } from '../../proto/feature/feature_pb';
import { Evaluator } from '../../evaluation';
import { createFeature, createUser } from '../../modelFactory';
import { Strategy } from '../../proto/feature/strategy_pb';

test('convertVariationValue: Non-YAML type returns original value', async (t) => {
  const evaluator = new Evaluator();
  const feature = createFeature({
    id: 'feature-1',
    name: 'String Feature',
    version: 1,
    enabled: true,
    createdAt: Date.now(),
    variationType: Feature.VariationType.STRING,
    variations: [
      {
        id: 'var-1',
        value: 'simple string',
        name: 'String Value',
        description: '',
      },
    ],
    targets: [],
    rules: [],
    defaultStrategy: {
      type: Strategy.Type.FIXED,
      variation: 'var-1',
    },
    prerequisitesList: [],
  });

  const user = createUser('test-user-1', {});
  const result = await evaluator.evaluateFeatures([feature], user, new Map(), '');

  t.is(result.getEvaluationsList().length, 1);
  const evaluation = result.getEvaluationsList()[0];
  t.truthy(evaluation);
  if (!evaluation) return;
  t.is(evaluation.getVariationValue(), 'simple string');
});

test('convertVariationValue: JSON type returns original value', async (t) => {
  const evaluator = new Evaluator();
  const feature = createFeature({
    id: 'feature-1',
    name: 'JSON Feature',
    version: 1,
    enabled: true,
    createdAt: Date.now(),
    variationType: Feature.VariationType.JSON,
    variations: [
      {
        id: 'var-1',
        value: '{"key": "value"}',
        name: 'JSON Value',
        description: '',
      },
    ],
    targets: [],
    rules: [],
    defaultStrategy: {
      type: Strategy.Type.FIXED,
      variation: 'var-1',
    },
    prerequisitesList: [],
  });

  const user = createUser('test-user-1', {});
  const result = await evaluator.evaluateFeatures([feature], user, new Map(), '');

  t.is(result.getEvaluationsList().length, 1);
  const evaluation = result.getEvaluationsList()[0];
  t.truthy(evaluation);
  if (!evaluation) return;
  t.is(evaluation.getVariationValue(), '{"key": "value"}');
});

test('convertVariationValue: YAML type converts to JSON', async (t) => {
  const evaluator = new Evaluator();
  const yamlValue = `name: John Doe
age: 30
active: true`;

  const feature = createFeature({
    id: 'yaml-feature',
    name: 'YAML Feature',
    version: 1,
    enabled: true,
    createdAt: Date.now(),
    variationType: Feature.VariationType.YAML,
    variations: [
      {
        id: 'yaml-var-1',
        value: yamlValue,
        name: 'YAML Config',
        description: '',
      },
    ],
    targets: [],
    rules: [],
    defaultStrategy: {
      type: Strategy.Type.FIXED,
      variation: 'yaml-var-1',
    },
    prerequisitesList: [],
  });

  const user = createUser('test-user-1', {});
  const result = await evaluator.evaluateFeatures([feature], user, new Map(), '');

  t.is(result.getEvaluationsList().length, 1);
  const evaluation = result.getEvaluationsList()[0];
  t.truthy(evaluation);
  if (!evaluation) return;

  // Verify it's valid JSON
  const jsonData = JSON.parse(evaluation.getVariationValue());
  t.is(jsonData.name, 'John Doe');
  t.is(jsonData.age, 30);
  t.is(jsonData.active, true);

  // Verify Variation.Value is also converted
  t.is(evaluation.getVariation()?.getValue(), evaluation.getVariationValue());
});

test('convertVariationValue: YAML with nested objects converts to JSON', async (t) => {
  const evaluator = new Evaluator();
  const yamlValue = `user:
  name: Jane
  email: jane@example.com
settings:
  theme: dark
  notifications: true`;

  const feature = createFeature({
    id: 'yaml-feature',
    name: 'YAML Feature',
    version: 1,
    enabled: true,
    createdAt: Date.now(),
    variationType: Feature.VariationType.YAML,
    variations: [
      {
        id: 'yaml-var-nested',
        value: yamlValue,
        name: 'Nested YAML Config',
        description: '',
      },
    ],
    targets: [],
    rules: [],
    defaultStrategy: {
      type: Strategy.Type.FIXED,
      variation: 'yaml-var-nested',
    },
    prerequisitesList: [],
  });

  const user = createUser('test-user-1', {});
  const result = await evaluator.evaluateFeatures([feature], user, new Map(), '');

  t.is(result.getEvaluationsList().length, 1);
  const evaluation = result.getEvaluationsList()[0];
  t.truthy(evaluation);
  if (!evaluation) return;

  // Verify it's valid JSON
  const jsonData = JSON.parse(evaluation.getVariationValue());
  t.is(jsonData.user.name, 'Jane');
  t.is(jsonData.user.email, 'jane@example.com');
  t.is(jsonData.settings.theme, 'dark');
  t.is(jsonData.settings.notifications, true);
});

test('convertVariationValue: YAML with arrays converts to JSON', async (t) => {
  const evaluator = new Evaluator();
  const yamlValue = `items:
  - id: 1
    name: Item 1
  - id: 2
    name: Item 2`;

  const feature = createFeature({
    id: 'yaml-feature',
    name: 'YAML Feature',
    version: 1,
    enabled: true,
    createdAt: Date.now(),
    variationType: Feature.VariationType.YAML,
    variations: [
      {
        id: 'yaml-var-array',
        value: yamlValue,
        name: 'Array YAML Config',
        description: '',
      },
    ],
    targets: [],
    rules: [],
    defaultStrategy: {
      type: Strategy.Type.FIXED,
      variation: 'yaml-var-array',
    },
    prerequisitesList: [],
  });

  const user = createUser('test-user-1', {});
  const result = await evaluator.evaluateFeatures([feature], user, new Map(), '');

  t.is(result.getEvaluationsList().length, 1);
  const evaluation = result.getEvaluationsList()[0];
  t.truthy(evaluation);
  if (!evaluation) return;

  // Verify it's valid JSON
  const jsonData = JSON.parse(evaluation.getVariationValue());
  t.is(jsonData.items.length, 2);
  t.is(jsonData.items[0].id, 1);
  t.is(jsonData.items[0].name, 'Item 1');
  t.is(jsonData.items[1].id, 2);
  t.is(jsonData.items[1].name, 'Item 2');
});

test('convertVariationValue: YAML with comments converts to JSON', async (t) => {
  const evaluator = new Evaluator();
  const yamlValue = `# This is a configuration
name: John Doe
# Age in years
age: 30
active: true # User is active`;

  const feature = createFeature({
    id: 'yaml-feature',
    name: 'YAML Feature',
    version: 1,
    enabled: true,
    createdAt: Date.now(),
    variationType: Feature.VariationType.YAML,
    variations: [
      {
        id: 'yaml-var-comments',
        value: yamlValue,
        name: 'YAML with Comments',
        description: '',
      },
    ],
    targets: [],
    rules: [],
    defaultStrategy: {
      type: Strategy.Type.FIXED,
      variation: 'yaml-var-comments',
    },
    prerequisitesList: [],
  });

  const user = createUser('test-user-1', {});
  const result = await evaluator.evaluateFeatures([feature], user, new Map(), '');

  t.is(result.getEvaluationsList().length, 1);
  const evaluation = result.getEvaluationsList()[0];
  t.truthy(evaluation);
  if (!evaluation) return;

  // Verify it's valid JSON and comments are stripped
  const jsonData = JSON.parse(evaluation.getVariationValue());
  t.is(jsonData.name, 'John Doe');
  t.is(jsonData.age, 30);
  t.is(jsonData.active, true);

  // Ensure no comments in the JSON string
  t.false(evaluation.getVariationValue().includes('#'));
});

test('convertVariationValue: Invalid YAML returns original value as fallback', async (t) => {
  const evaluator = new Evaluator();
  const invalidYamlValue = 'invalid: yaml: [unclosed';

  const feature = createFeature({
    id: 'yaml-feature',
    name: 'YAML Feature',
    version: 1,
    enabled: true,
    createdAt: Date.now(),
    variationType: Feature.VariationType.YAML,
    variations: [
      {
        id: 'yaml-var-invalid',
        value: invalidYamlValue,
        name: 'Invalid YAML',
        description: '',
      },
    ],
    targets: [],
    rules: [],
    defaultStrategy: {
      type: Strategy.Type.FIXED,
      variation: 'yaml-var-invalid',
    },
    prerequisitesList: [],
  });

  const user = createUser('test-user-1', {});
  const result = await evaluator.evaluateFeatures([feature], user, new Map(), '');

  t.is(result.getEvaluationsList().length, 1);
  const evaluation = result.getEvaluationsList()[0];
  t.truthy(evaluation);
  if (!evaluation) return;

  // Should return original value as fallback
  t.is(evaluation.getVariationValue(), invalidYamlValue);
});

test('convertVariationValue: Cache is used across multiple evaluations', async (t) => {
  const evaluator = new Evaluator();
  const yamlValue = `settings:
  theme: dark
  language: en
  notifications:
    email: true
    push: false`;

  const feature = createFeature({
    id: 'cached-yaml-feature',
    name: 'Cached YAML Feature',
    version: 1,
    enabled: true,
    createdAt: Date.now(),
    variationType: Feature.VariationType.YAML,
    variations: [
      {
        id: 'cached-yaml-var',
        value: yamlValue,
        name: 'Cached Config',
        description: '',
      },
    ],
    targets: [],
    rules: [],
    defaultStrategy: {
      type: Strategy.Type.FIXED,
      variation: 'cached-yaml-var',
    },
    prerequisitesList: [],
  });

  // First evaluation
  const user1 = createUser('test-user-1', {});
  const result1 = await evaluator.evaluateFeatures([feature], user1, new Map(), '');
  t.truthy(result1.getEvaluationsList()[0]);
  const value1 = result1.getEvaluationsList()[0]?.getVariationValue();
  t.truthy(value1);
  if (!value1) return;

  // Second evaluation with different user (should use cache)
  const user2 = createUser('test-user-2', {});
  const result2 = await evaluator.evaluateFeatures([feature], user2, new Map(), '');
  t.truthy(result2.getEvaluationsList()[0]);
  const value2 = result2.getEvaluationsList()[0]?.getVariationValue();
  t.truthy(value2);
  if (!value2) return;

  // Both should have the same converted JSON value
  t.is(value1, value2);

  // Verify structure
  const jsonData = JSON.parse(value1);
  t.is(jsonData.settings.theme, 'dark');
  t.is(jsonData.settings.language, 'en');
  t.is(jsonData.settings.notifications.email, true);
  t.is(jsonData.settings.notifications.push, false);
});

test('convertVariationValue: Mixed variation types in single evaluation', async (t) => {
  const evaluator = new Evaluator();

  const yamlFeature = createFeature({
    id: 'yaml-feature',
    name: 'YAML Feature',
    version: 1,
    enabled: true,
    createdAt: Date.now(),
    variationType: Feature.VariationType.YAML,
    variations: [
      {
        id: 'yaml-var',
        value: 'enabled: true\ntimeout: 30',
        name: 'YAML Config',
        description: '',
      },
    ],
    targets: [],
    rules: [],
    defaultStrategy: {
      type: Strategy.Type.FIXED,
      variation: 'yaml-var',
    },
    prerequisitesList: [],
  });

  const stringFeature = createFeature({
    id: 'string-feature',
    name: 'String Feature',
    version: 1,
    enabled: true,
    createdAt: Date.now(),
    variationType: Feature.VariationType.STRING,
    variations: [
      {
        id: 'string-var',
        value: 'simple-string',
        name: 'String Value',
        description: '',
      },
    ],
    targets: [],
    rules: [],
    defaultStrategy: {
      type: Strategy.Type.FIXED,
      variation: 'string-var',
    },
    prerequisitesList: [],
  });

  const jsonFeature = createFeature({
    id: 'json-feature',
    name: 'JSON Feature',
    version: 1,
    enabled: true,
    createdAt: Date.now(),
    variationType: Feature.VariationType.JSON,
    variations: [
      {
        id: 'json-var',
        value: '{"key":"value"}',
        name: 'JSON Value',
        description: '',
      },
    ],
    targets: [],
    rules: [],
    defaultStrategy: {
      type: Strategy.Type.FIXED,
      variation: 'json-var',
    },
    prerequisitesList: [],
  });

  const user = createUser('test-user', {});
  const result = await evaluator.evaluateFeatures(
    [yamlFeature, stringFeature, jsonFeature],
    user,
    new Map(),
    '',
  );

  t.is(result.getEvaluationsList().length, 3);

  // Find each evaluation
  const yamlEval = result.getEvaluationsList().find((e) => e.getFeatureId() === 'yaml-feature');
  const stringEval = result.getEvaluationsList().find((e) => e.getFeatureId() === 'string-feature');
  const jsonEval = result.getEvaluationsList().find((e) => e.getFeatureId() === 'json-feature');

  t.truthy(yamlEval);
  t.truthy(stringEval);
  t.truthy(jsonEval);
  if (!yamlEval || !stringEval || !jsonEval) return;

  // YAML should be converted to JSON
  const yamlData = JSON.parse(yamlEval.getVariationValue());
  t.is(yamlData.enabled, true);
  t.is(yamlData.timeout, 30);

  // String should remain as-is
  t.is(stringEval.getVariationValue(), 'simple-string');

  // JSON should remain as-is
  t.is(jsonEval.getVariationValue(), '{"key":"value"}');
});
