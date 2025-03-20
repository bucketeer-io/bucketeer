import {
  IndividualRuleItem,
  PrerequisiteRuleType,
  SegmentConditionType
} from './types';

const initialSegmentCondition: SegmentConditionType = {
  situation: 'compare',
  conditioner: '=',
  firstValue: '1',
  secondValue: '',
  value: '',
  date: '',
  flagId: '',
  variation: ''
};

const initialIndividualRule: IndividualRuleItem = {
  variationId: '',
  users: []
};

const initialPrerequisitesRule: PrerequisiteRuleType = {
  featureFlag: '',
  variation: ''
};

export {
  initialSegmentCondition,
  initialIndividualRule,
  initialPrerequisitesRule
};
