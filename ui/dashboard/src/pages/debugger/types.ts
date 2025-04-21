import { Account, Feature, FeatureVariation } from '@types';

export interface Evaluation {
  id: string;
  feature_id: string;
  feature_version: number;
  user_id: string;
  variation_id: string;
  variation: FeatureVariation;
  reason: string;
  variation_value: string;
  variation_name: string;
}

export interface EvaluationFeatureAccount extends Evaluation {
  feature: Feature;
  account: Account;
}
