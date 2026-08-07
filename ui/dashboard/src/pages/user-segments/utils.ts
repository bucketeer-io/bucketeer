import isEqual from 'lodash/isEqual';
import { v4 as uuid } from 'uuid';
import {
  FeatureRule,
  FeatureRuleClauseOperator,
  SegmentRulePayload
} from '@types';
import { RuleClauseType } from 'pages/feature-flag-details/targeting/types';

export type SegmentClauseType = RuleClauseType.COMPARE | RuleClauseType.DATE;

// Local mirror of the targeting tab's getClauseType, restricted to the clause
// types segment rules support (the server rejects SEGMENT and FEATURE_FLAG
// operators inside segment rules).
export const getSegmentClauseType = (
  operator: FeatureRuleClauseOperator
): SegmentClauseType => {
  const { BEFORE, AFTER } = FeatureRuleClauseOperator;
  if ([BEFORE, AFTER].includes(operator)) return RuleClauseType.DATE;
  return RuleClauseType.COMPARE;
};

export interface SegmentRuleClauseFormValue {
  id: string;
  type: SegmentClauseType;
  attribute: string;
  operator: string;
  values: (string | undefined)[];
}

export interface SegmentRuleFormValue {
  id: string;
  clauses: SegmentRuleClauseFormValue[];
}

export const createSegmentRule = (): SegmentRuleFormValue => ({
  id: uuid(),
  clauses: [
    {
      id: uuid(),
      type: RuleClauseType.COMPARE,
      attribute: '',
      operator: FeatureRuleClauseOperator.EQUALS,
      values: []
    }
  ]
});

export const toSegmentFormRules = (
  rules?: FeatureRule[]
): SegmentRuleFormValue[] =>
  (rules || [])
    .filter(rule => rule !== undefined && rule !== null)
    .map(({ id, clauses }) => ({
      id,
      clauses: (clauses || [])
        .filter(clause => clause !== undefined && clause !== null)
        .map(clause => ({
          id: clause.id,
          type: getSegmentClauseType(clause.operator),
          attribute: clause.attribute || '',
          operator: clause.operator as string,
          values: clause.values || []
        }))
    }));

/**
 * Builds the rules payload sent to the segment API.
 * Rule/clause IDs are generated server-side: IDs that exist on the loaded
 * segment are preserved, everything else (locally generated form keys) is
 * sent empty.
 */
export const toSegmentRulesPayload = (
  formRules: SegmentRuleFormValue[] | undefined,
  originalRules: FeatureRule[]
): SegmentRulePayload[] => {
  const originalClauseIds = new Map(
    (originalRules || []).map(rule => [
      rule.id,
      new Set((rule.clauses || []).map(clause => clause.id))
    ])
  );
  return (formRules || []).map(rule => {
    const existingClauseIds = originalClauseIds.get(rule.id);
    return {
      id: existingClauseIds ? rule.id : '',
      clauses: (rule.clauses || []).map(clause => ({
        id: existingClauseIds?.has(clause.id) ? clause.id : '',
        attribute: clause.attribute || '',
        operator: clause.operator,
        values: (clause.values || []).filter(
          (value): value is string => value !== undefined
        )
      }))
    };
  });
};

/**
 * Whether the form's rules differ from the segment's stored rules.
 * The update API treats an absent rules field as "unchanged", so the caller
 * must only send rules when this returns true.
 */
export const hasSegmentRulesChanged = (
  formRules: SegmentRuleFormValue[] | undefined,
  originalRules: FeatureRule[]
): boolean => {
  const current = toSegmentRulesPayload(formRules, originalRules);
  const original = (originalRules || []).map(rule => ({
    id: rule.id,
    clauses: (rule.clauses || []).map(clause => ({
      id: clause.id,
      attribute: clause.attribute || '',
      operator: clause.operator as string,
      values: clause.values || []
    }))
  }));
  return !isEqual(current, original);
};
