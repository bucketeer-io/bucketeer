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
  flag: ''
};

const initialIndividualRule: IndividualRuleItem = {
  on: [],
  off: []
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
