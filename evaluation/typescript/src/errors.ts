const EVALUATOR_ERRORS = {
  StrategyNotFound: new Error('evaluator: strategy not found'),
  DefaultStrategyNotFound: new Error('evaluator: default strategy not found'),
  FeatureNotFound: new Error('evaluator: feature not found'),
  PrerequisiteVariationNotFound: new Error('evaluator: prerequisite variation not found'),
  VariationNotFound: new Error('evaluator: variation not found'),
  UnsupportedStrategy: new Error('evaluator: unsupported strategy'),
};

export { EVALUATOR_ERRORS };
