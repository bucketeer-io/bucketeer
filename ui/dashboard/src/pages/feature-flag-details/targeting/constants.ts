import { i18n } from 'i18n';
import { FeatureRuleClauseOperator } from '@types';
import { PrerequisiteSchema, RuleClauseType } from './types';

const translation = i18n.t;

const initialPrerequisite: PrerequisiteSchema = {
  featureId: '',
  variationId: ''
};

const situationOptions = [
  {
    label: translation('form:feature-flags.compare'),
    value: RuleClauseType.COMPARE
  },
  {
    label: translation('form:feature-flags.user-segment'),
    value: RuleClauseType.SEGMENT
  },
  {
    label: translation('form:feature-flags.date'),
    value: RuleClauseType.DATE
  },
  {
    label: translation('form:feature-flags.feature-flag'),
    value: RuleClauseType.FEATURE_FLAG
  }
];

const conditionerCompareOptions = [
  {
    label: '=',
    value: FeatureRuleClauseOperator.EQUALS
  },
  {
    label: '>=',
    value: FeatureRuleClauseOperator.GREATER_OR_EQUAL
  },
  {
    label: '>',
    value: FeatureRuleClauseOperator.GREATER
  },
  {
    label: '<=',
    value: FeatureRuleClauseOperator.LESS_OR_EQUAL
  },
  {
    label: '<',
    value: FeatureRuleClauseOperator.LESS
  },
  {
    label: translation('form:contains'),
    value: FeatureRuleClauseOperator.IN
  },
  {
    label: translation('form:partially-matches'),
    value: FeatureRuleClauseOperator.PARTIALLY_MATCH
  },
  {
    label: translation('form:starts-with'),
    value: FeatureRuleClauseOperator.STARTS_WITH
  },
  {
    label: translation('form:ends-with'),
    value: FeatureRuleClauseOperator.ENDS_WITH
  }
];

const conditionerDateOptions = [
  {
    label: translation('form:before'),
    value: FeatureRuleClauseOperator.BEFORE
  },
  {
    label: translation('form:after'),
    value: FeatureRuleClauseOperator.AFTER
  }
];

export {
  initialPrerequisite,
  situationOptions,
  conditionerCompareOptions,
  conditionerDateOptions
};
