import {
  IndividualRuleItem,
  PrerequisiteSchema,
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

const initialPrerequisite: PrerequisiteSchema = {
  featureId: '',
  variationId: ''
};

export { initialSegmentCondition, initialIndividualRule, initialPrerequisite };
