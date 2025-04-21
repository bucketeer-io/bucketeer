import { Evaluation, Feature } from '@types';

export interface EvaluationFeature extends Evaluation {
  feature: Feature;
}
