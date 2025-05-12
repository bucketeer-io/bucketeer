import { PrerequisiteSchema, IndividualRuleItem } from './types';

const initialPrerequisite: PrerequisiteSchema = {
  featureId: '',
  variationId: ''
};

const initialIndividualRule: IndividualRuleItem = {
  variationId: '',
  name: '',
  users: []
};

export { initialPrerequisite, initialIndividualRule };
