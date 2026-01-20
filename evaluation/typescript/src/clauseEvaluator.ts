import { Clause } from './proto/feature/clause_pb';
import { SegmentUser } from './proto/feature/segment_pb';
import { SegmentEvaluator } from './segmentEvaluator';
import { DependencyEvaluator } from './dependencyEvaluator';
import * as semver from 'semver';
//
class ClauseEvaluator {
  private segmentEvaluator: SegmentEvaluator;
  private dependencyEvaluator: DependencyEvaluator;

  constructor() {
    this.segmentEvaluator = new SegmentEvaluator();
    this.dependencyEvaluator = new DependencyEvaluator();
  }

  evaluate(
    targetValue: string,
    clause: Clause,
    userID: string,
    segmentUsers: SegmentUser[],
    flagVariations: { [key: string]: string },
  ): boolean {
    try {
      switch (clause.getOperator()) {
        case Clause.Operator.EQUALS:
          return this.equals(targetValue, clause.getValuesList());
        case Clause.Operator.IN:
          return this.in(targetValue, clause.getValuesList());
        case Clause.Operator.STARTS_WITH:
          return this.startsWith(targetValue, clause.getValuesList());
        case Clause.Operator.ENDS_WITH:
          return this.endsWith(targetValue, clause.getValuesList());
        case Clause.Operator.SEGMENT:
          return this.segmentEvaluator.evaluate(clause.getValuesList(), userID, segmentUsers);
        case Clause.Operator.GREATER:
          return this.greater(targetValue, clause.getValuesList());
        case Clause.Operator.GREATER_OR_EQUAL:
          return this.greaterOrEqual(targetValue, clause.getValuesList());
        case Clause.Operator.LESS:
          return this.less(targetValue, clause.getValuesList());
        case Clause.Operator.LESS_OR_EQUAL:
          return this.lessOrEqual(targetValue, clause.getValuesList());
        case Clause.Operator.BEFORE:
          return this.before(targetValue, clause.getValuesList());
        case Clause.Operator.AFTER:
          return this.after(targetValue, clause.getValuesList());
        case Clause.Operator.FEATURE_FLAG:
          return this.dependencyEvaluator.evaluate(
            clause.getAttribute(),
            clause.getValuesList(),
            flagVariations,
          );
        case Clause.Operator.PARTIALLY_MATCH:
          return this.partiallyMatches(targetValue, clause.getValuesList());
        case Clause.Operator.NOT_EQUALS:
          return !this.equals(targetValue, clause.getValuesList());
        default:
          return false;
      }
    } catch (error) {
      console.error('Error evaluating clause:', error);
      throw error;
    }
  }

  private equals(targetValue: string, values: string[]): boolean {
    // First try semver comparison (with v-prefix normalization via semver.valid)
    const semverTarget = semver.valid(targetValue);
    if (semverTarget) {
      const semverValues = values
        .map((v) => semver.valid(v))
        .filter((v): v is string => v !== null);
      // Target parsed as semver; if no values parsed or none matched, return false
      // (consistent with greaterOrEqual, greater, less, lessOrEqual behavior)
      return semverValues.some((value) => semver.eq(semverTarget, value));
    }
    // Fall back to exact string comparison only when target fails to parse as semver
    return values.includes(targetValue);
  }

  private partiallyMatches(targetValue: string, values: string[]): boolean {
    return values.some((value) => targetValue.includes(value));
  }

  private in(targetValue: string, values: string[]): boolean {
    // First try semver comparison (with v-prefix normalization via semver.valid)
    const semverTarget = semver.valid(targetValue);
    if (semverTarget) {
      const semverValues = values
        .map((v) => semver.valid(v))
        .filter((v): v is string => v !== null);
      // Target parsed as semver; if no values parsed or none matched, return false
      // (consistent with greaterOrEqual, greater, less, lessOrEqual behavior)
      return semverValues.some((value) => semver.eq(semverTarget, value));
    }
    // Fall back to exact string comparison only when target fails to parse as semver
    return values.includes(targetValue);
  }

  private startsWith(targetValue: string, values: string[]): boolean {
    return values.some((value) => targetValue.startsWith(value));
  }

  private endsWith(targetValue: string, values: string[]): boolean {
    return values.some((value) => targetValue.endsWith(value));
  }

  private greater(targetValue: string, values: string[]): boolean {
    const floatTarget = validateAndParseFloat(targetValue);
    const floatValues = values.map(validateAndParseFloat).filter((value) => !isNaN(value));

    if (!isNaN(floatTarget)) {
      if (floatValues.length == 0) {
        return false;
      }
      return floatValues.some((value) => floatTarget > value);
    }

    const semverTarget = semver.valid(targetValue);
    if (semverTarget) {
      return values.some((value) => {
        if (semver.valid(value)) {
          return semver.gt(semverTarget, value);
        }
        return false;
      });
    }

    return values.some((value) => targetValue > value);
  }

  private greaterOrEqual(targetValue: string, values: string[]): boolean {
    const floatTarget = validateAndParseFloat(targetValue);
    const floatValues = values.map(validateAndParseFloat).filter((value) => !isNaN(value));

    if (!isNaN(floatTarget)) {
      if (floatValues.length == 0) {
        return false;
      }
      return floatValues.some((value) => floatTarget >= value);
    }

    const semverTarget = semver.valid(targetValue);
    if (semverTarget) {
      return values.some((value) => {
        if (semver.valid(value)) {
          return semver.gte(semverTarget, value);
        }
        return false;
      });
    }

    return values.some((value) => targetValue >= value);
  }

  private less(targetValue: string, values: string[]): boolean {
    const floatTarget = validateAndParseFloat(targetValue);
    const floatValues = values.map(validateAndParseFloat).filter((value) => !isNaN(value));

    if (!isNaN(floatTarget)) {
      if (floatValues.length == 0) {
        return false;
      }
      return floatValues.some((value) => floatTarget < value);
    }

    const semverTarget = semver.valid(targetValue);
    if (semverTarget) {
      return values.some((value) => {
        if (semver.valid(value)) {
          return semver.lt(semverTarget, value);
        }
        return false;
      });
    }

    return values.some((value) => targetValue < value);
  }

  private lessOrEqual(targetValue: string, values: string[]): boolean {
    const floatTarget = validateAndParseFloat(targetValue);
    const floatValues = values.map(validateAndParseFloat).filter((value) => !isNaN(value));

    if (!isNaN(floatTarget)) {
      if (floatValues.length == 0) {
        return false;
      }
      return floatValues.some((value) => floatTarget <= value);
    }

    const semverTarget = semver.valid(targetValue);
    if (semverTarget) {
      return values.some((value) => {
        if (semver.valid(value)) {
          return semver.lte(semverTarget, value);
        }
        return false;
      });
    }

    return values.some((value) => targetValue <= value);
  }

  private before(targetValue: string, values: string[]): boolean {
    const intTarget = validateAndParseInt(targetValue);
    if (isNaN(intTarget)) {
      return false;
    }
    const intValues = values.map((value) => parseInt(value, 10)).filter((value) => !isNaN(value));
    if (intValues.length == 0) {
      return false;
    }

    return intValues.some((value) => intTarget < value);
  }

  private after(targetValue: string, values: string[]): boolean {
    const intTarget = parseInt(targetValue, 10);
    const intValues = values.map((value) => parseInt(value, 10)).filter((value) => !isNaN(value));

    return intValues.some((value) => intTarget > value);
  }
}

function isNumericString(input: string): boolean {
  return /^-?\d*\.?\d+$/.test(input);
}

function validateAndParseInt(input: string): number {
  // Check if the string is a valid number using regex
  if (!isNumericString(input)) {
    // Will not allow a string like `1a`, `0b`
    return NaN;
  }

  // Parse the string to float
  const result = parseInt(input, 10);

  // If parsing fails, return NaN
  return isNaN(result) ? NaN : result;
}

function validateAndParseFloat(input: string): number {
  // Check if the string is a valid number using regex
  if (!isNumericString(input)) {
    // Will not allow a string like `1a`, `0b`
    return NaN;
  }

  // Parse the string to float
  const result = parseFloat(input);

  // If parsing fails, return NaN
  return isNaN(result) ? NaN : result;
}

export { ClauseEvaluator };
