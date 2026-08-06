import { describe, expect, it } from 'vitest';
import { FeatureRule, FeatureRuleClauseOperator } from '@types';
import { RuleClauseType } from 'pages/feature-flag-details/targeting/types';
import {
  createSegmentRule,
  getSegmentClauseType,
  hasSegmentRulesChanged,
  SegmentRuleFormValue,
  toSegmentFormRules,
  toSegmentRulesPayload
} from './utils';

const makeOriginalRule = (
  ruleId: string,
  clauseId: string,
  overrides?: Partial<{ attribute: string; operator: string; values: string[] }>
): FeatureRule =>
  ({
    id: ruleId,
    clauses: [
      {
        id: clauseId,
        attribute: overrides?.attribute ?? 'email',
        operator: overrides?.operator ?? FeatureRuleClauseOperator.ENDS_WITH,
        values: overrides?.values ?? ['@example.com']
      }
    ]
  }) as FeatureRule;

describe('getSegmentClauseType', () => {
  it('maps BEFORE/AFTER to date and everything else to compare', () => {
    expect(getSegmentClauseType(FeatureRuleClauseOperator.BEFORE)).toBe(
      RuleClauseType.DATE
    );
    expect(getSegmentClauseType(FeatureRuleClauseOperator.AFTER)).toBe(
      RuleClauseType.DATE
    );
    expect(getSegmentClauseType(FeatureRuleClauseOperator.EQUALS)).toBe(
      RuleClauseType.COMPARE
    );
    expect(getSegmentClauseType(FeatureRuleClauseOperator.ENDS_WITH)).toBe(
      RuleClauseType.COMPARE
    );
  });
});

describe('createSegmentRule', () => {
  it('creates a rule with one empty compare clause and local ids', () => {
    const rule = createSegmentRule();
    expect(rule.id).toBeTruthy();
    expect(rule.clauses).toHaveLength(1);
    expect(rule.clauses[0]).toMatchObject({
      type: RuleClauseType.COMPARE,
      attribute: '',
      operator: FeatureRuleClauseOperator.EQUALS,
      values: []
    });
    expect(rule.clauses[0].id).toBeTruthy();
  });
});

describe('toSegmentFormRules', () => {
  it('maps API rules to form values with a UI clause type', () => {
    const rules = [
      makeOriginalRule('rule-1', 'clause-1'),
      makeOriginalRule('rule-2', 'clause-2', {
        operator: FeatureRuleClauseOperator.BEFORE,
        attribute: 'signup_date',
        values: ['1700000000']
      })
    ];
    expect(toSegmentFormRules(rules)).toEqual([
      {
        id: 'rule-1',
        clauses: [
          {
            id: 'clause-1',
            type: RuleClauseType.COMPARE,
            attribute: 'email',
            operator: FeatureRuleClauseOperator.ENDS_WITH,
            values: ['@example.com']
          }
        ]
      },
      {
        id: 'rule-2',
        clauses: [
          {
            id: 'clause-2',
            type: RuleClauseType.DATE,
            attribute: 'signup_date',
            operator: FeatureRuleClauseOperator.BEFORE,
            values: ['1700000000']
          }
        ]
      }
    ]);
  });

  it('returns an empty array for missing rules', () => {
    expect(toSegmentFormRules(undefined)).toEqual([]);
  });
});

describe('toSegmentRulesPayload', () => {
  it('preserves ids of existing rules/clauses and empties new ones', () => {
    const original = [makeOriginalRule('rule-1', 'clause-1')];
    const formRules: SegmentRuleFormValue[] = [
      {
        // Existing rule: id preserved, new clause id emptied.
        id: 'rule-1',
        clauses: [
          {
            id: 'clause-1',
            type: RuleClauseType.COMPARE,
            attribute: 'email',
            operator: FeatureRuleClauseOperator.ENDS_WITH as string,
            values: ['@example.com']
          },
          {
            id: 'local-uuid-clause',
            type: RuleClauseType.COMPARE,
            attribute: 'plan',
            operator: FeatureRuleClauseOperator.EQUALS as string,
            values: ['pro']
          }
        ]
      },
      {
        // New rule: rule and clause ids emptied.
        id: 'local-uuid-rule',
        clauses: [
          {
            id: 'local-uuid-clause-2',
            type: RuleClauseType.COMPARE,
            attribute: 'country',
            operator: FeatureRuleClauseOperator.IN as string,
            values: ['JP', 'US']
          }
        ]
      }
    ];

    expect(toSegmentRulesPayload(formRules, original)).toEqual([
      {
        id: 'rule-1',
        clauses: [
          {
            id: 'clause-1',
            attribute: 'email',
            operator: FeatureRuleClauseOperator.ENDS_WITH,
            values: ['@example.com']
          },
          {
            id: '',
            attribute: 'plan',
            operator: FeatureRuleClauseOperator.EQUALS,
            values: ['pro']
          }
        ]
      },
      {
        id: '',
        clauses: [
          {
            id: '',
            attribute: 'country',
            operator: FeatureRuleClauseOperator.IN,
            values: ['JP', 'US']
          }
        ]
      }
    ]);
  });

  it('does not preserve a clause id moved to a different rule', () => {
    const original = [
      makeOriginalRule('rule-1', 'clause-1'),
      makeOriginalRule('rule-2', 'clause-2')
    ];
    const formRules: SegmentRuleFormValue[] = [
      {
        id: 'rule-1',
        clauses: [
          {
            id: 'clause-2',
            type: RuleClauseType.COMPARE,
            attribute: 'email',
            operator: FeatureRuleClauseOperator.EQUALS as string,
            values: ['a']
          }
        ]
      }
    ];
    expect(toSegmentRulesPayload(formRules, original)[0].clauses[0].id).toBe(
      ''
    );
  });

  it('strips the UI-only clause type and undefined values', () => {
    const payload = toSegmentRulesPayload(
      [
        {
          id: 'new',
          clauses: [
            {
              id: 'new-clause',
              type: RuleClauseType.COMPARE,
              attribute: 'a',
              operator: FeatureRuleClauseOperator.EQUALS as string,
              values: ['1', undefined, '2']
            }
          ]
        }
      ],
      []
    );
    expect(payload[0].clauses[0]).not.toHaveProperty('type');
    expect(payload[0].clauses[0].values).toEqual(['1', '2']);
  });

  it('returns an empty list when the form has no rules', () => {
    expect(toSegmentRulesPayload(undefined, [])).toEqual([]);
    expect(toSegmentRulesPayload([], [])).toEqual([]);
  });
});

describe('hasSegmentRulesChanged', () => {
  const original = [makeOriginalRule('rule-1', 'clause-1')];

  it('returns false when the form mirrors the stored rules', () => {
    expect(hasSegmentRulesChanged(toSegmentFormRules(original), original)).toBe(
      false
    );
  });

  it('returns false for empty form rules and empty stored rules', () => {
    expect(hasSegmentRulesChanged([], [])).toBe(false);
    expect(hasSegmentRulesChanged(undefined, [])).toBe(false);
  });

  it('detects an edited clause value', () => {
    const formRules = toSegmentFormRules(original);
    formRules[0].clauses[0].values = ['@other.com'];
    expect(hasSegmentRulesChanged(formRules, original)).toBe(true);
  });

  it('detects an added rule', () => {
    const formRules = [...toSegmentFormRules(original), createSegmentRule()];
    expect(hasSegmentRulesChanged(formRules, original)).toBe(true);
  });

  it('detects a deleted rule (empty replacement list)', () => {
    expect(hasSegmentRulesChanged([], original)).toBe(true);
  });

  it('detects reordered rules (full replacement is ordered)', () => {
    const twoRules = [
      makeOriginalRule('rule-1', 'clause-1'),
      makeOriginalRule('rule-2', 'clause-2', { attribute: 'plan' })
    ];
    const formRules = toSegmentFormRules(twoRules).reverse();
    expect(hasSegmentRulesChanged(formRules, twoRules)).toBe(true);
  });
});
