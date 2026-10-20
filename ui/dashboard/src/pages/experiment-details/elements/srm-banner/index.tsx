import { useMemo, useState } from 'react';
import { useTranslation } from 'i18n';
import { Experiment, SrmResult, SrmStatus } from '@types';
import { cn } from 'utils/style';
import { IconAlert, IconChevronDown, IconInfoFilled } from '@icons';
import Icon from 'components/icon';
import SrmVariationsTable from './srm-variations-table';

// Canonical skip-reason prefixes emitted by the backend SRM pipeline.
// Reasons live in two files:
//   - pkg/experimentcalculator/experimentcalc/srm.go            (errSRM*)
//   - pkg/experimentcalculator/experimentcalc/experiment_calculator.go
//     (no-goal-results / feature-fetch-failed pre-checks)
// Most entries match the full canonical string so the localized text is
// shown verbatim. Only `insufficient-samples` and `small-expected-cell`
// append a genuinely dynamic numeric suffix (e.g. "... (got 37, need >=
// 100)"); for those the prefix matches only the canonical part and the
// suffix is appended verbatim. When extending this list, mirror the
// matching strategy to whichever style the backend emits.
const SKIP_REASON_KEYS: { prefix: string; key: string }[] = [
  {
    prefix: 'feature definition not available',
    key: 'table:results.srm.reason.feature-missing'
  },
  {
    prefix: 'could not fetch feature definition',
    key: 'table:results.srm.reason.feature-fetch-failed'
  },
  {
    prefix: 'no goal results available to derive observed traffic split',
    key: 'table:results.srm.reason.no-goal-results'
  },
  {
    prefix: 'feature has no default strategy',
    key: 'table:results.srm.reason.no-default-strategy'
  },
  {
    // Full canonical string (including the fixed parenthetical) — the
    // backend's errSRM* strings carry these clarifiers as static text,
    // not as dynamic context. Matching only the leading portion would
    // leave the clarifier dangling in English after the localized
    // translation, producing mixed-language output. Only the
    // *insufficient-samples* and *small-expected-cell* skip reasons
    // append genuinely dynamic numeric context — they keep prefix-only
    // matches below.
    prefix:
      'feature default strategy is not a rollout (no per-variation weights to test against)',
    key: 'table:results.srm.reason.not-rollout-strategy'
  },
  {
    prefix: 'feature rollout strategy has no variations',
    key: 'table:results.srm.reason.no-rollout-variations'
  },
  {
    prefix: 'all rollout weights are zero',
    key: 'table:results.srm.reason.all-zero-weights'
  },
  {
    prefix:
      'audience default_variation is not one of the rollout variations (cannot compute expected split)',
    key: 'table:results.srm.reason.audience-default-not-in-rollout'
  },
  {
    prefix:
      'total observed users below the minimum required for a reliable chi-square test',
    key: 'table:results.srm.reason.insufficient-samples'
  },
  {
    prefix:
      'smallest expected per-variation count below the chi-square reliability floor',
    key: 'table:results.srm.reason.small-expected-cell'
  },
  {
    prefix: 'fewer than 2 variations with positive expected user counts',
    key: 'table:results.srm.reason.too-few-cells'
  }
];

const formatPValue = (p: number) => {
  if (!Number.isFinite(p) || p <= 0) return '0';
  if (p < 0.001) return '< 0.001';
  return p.toFixed(3);
};

const SrmBanner = ({
  experiment,
  srmResult
}: {
  experiment: Experiment;
  srmResult?: SrmResult;
}) => {
  const { t } = useTranslation(['table']);
  const [isDetailsOpen, setIsDetailsOpen] = useState(false);
  const [isAdvancedOpen, setIsAdvancedOpen] = useState(false);

  const localizedSkipReason = useMemo(() => {
    if (!srmResult?.skipReason) return '';
    const match = SKIP_REASON_KEYS.find(({ prefix }) =>
      srmResult.skipReason.startsWith(prefix)
    );
    if (!match) return srmResult.skipReason;
    const suffix = srmResult.skipReason.slice(match.prefix.length);
    return `${t(match.key)}${suffix}`;
  }, [srmResult?.skipReason, t]);

  if (!srmResult || srmResult.status === SrmStatus.Unknown) {
    return null;
  }

  // Distinguish the two MISMATCH paths so we can show the right description:
  //   - chi-square fired: degreesOfFreedom > 0 AND pValue < threshold
  //     → show the chi-square-flavored copy with p-value/threshold inline.
  //   - leak-only:        chi-square either wasn't applicable (df = 0,
  //     fields are unpopulated zeros) OR passed (pValue >= threshold),
  //     so MISMATCH must have come from the zero-expected-cell leak
  //     detector → show the leak-flavored copy and let the per-variation
  //     table reveal which variations leaked.
  const chiSquareFired =
    Number(srmResult.degreesOfFreedom) > 0 &&
    srmResult.pValue < srmResult.threshold;

  if (srmResult.status === SrmStatus.Mismatch) {
    return (
      <div className="flex flex-col w-full overflow-hidden rounded-lg border border-accent-red-100 dark:border-dark-black-700 border-l-4 border-l-accent-red-500 bg-white dark:bg-dark-black-800 shadow-card dark:shadow-dark-card">
        <div className="flex flex-col w-full gap-y-2 bg-accent-red-50 dark:bg-accent-red-900/30 px-4 py-3">
          <div className="flex items-center w-full gap-x-2">
            <Icon icon={IconAlert} size="xs" color="accent-red-500" />
            <p className="typo-head-bold-small text-accent-red-500">
              {t('table:results.srm.mismatch-title')}
            </p>
          </div>
          <p className="typo-para-small text-gray-700 dark:text-dark-gray-300 pl-7">
            {chiSquareFired
              ? t('table:results.srm.mismatch-desc', {
                  pValue: formatPValue(srmResult.pValue),
                  threshold: formatPValue(srmResult.threshold)
                })
              : t('table:results.srm.mismatch-desc-leak')}
          </p>
          <button
            type="button"
            onClick={() => setIsDetailsOpen(prev => !prev)}
            className="flex w-fit items-center gap-x-2 pl-7"
          >
            <p className="typo-para-small font-medium text-accent-red-600">
              {t(
                isDetailsOpen
                  ? 'table:results.srm.hide-details'
                  : 'table:results.srm.view-details'
              )}
            </p>
            <Icon
              icon={IconChevronDown}
              size="xxs"
              color="accent-red-600"
              className={cn(
                'transition-transform duration-200',
                isDetailsOpen && 'rotate-180'
              )}
            />
          </button>
        </div>
        {isDetailsOpen && (
          <div className="bg-white dark:bg-dark-black-800 px-4 pt-3 border-t border-accent-red-100 dark:border-dark-black-700 overflow-x-auto overflow-y-hidden">
            <SrmVariationsTable
              variations={srmResult.variations}
              experimentVariations={experiment.variations}
            />
          </div>
        )}
      </div>
    );
  }

  if (srmResult.status === SrmStatus.Skipped) {
    return (
      <div className="flex items-center w-full gap-x-2 rounded border-l-4 border-accent-blue-500 bg-accent-blue-50 dark:bg-accent-blue-900/30 p-4">
        <Icon icon={IconInfoFilled} size="xxs" color="accent-blue-500" />
        <p className="typo-para-small text-accent-blue-500">
          {t('table:results.srm.skipped-prefix', {
            reason: localizedSkipReason
          })}
        </p>
      </div>
    );
  }

  // SrmStatus.Ok — collapsed-by-default advanced disclosure.
  return (
    <div className="flex flex-col w-full gap-y-2 rounded border border-gray-200 dark:border-dark-black-700 bg-white dark:bg-dark-black-800 p-3">
      <button
        type="button"
        onClick={() => setIsAdvancedOpen(prev => !prev)}
        className="flex items-center w-full justify-between gap-x-2"
      >
        <p className="typo-para-small text-gray-700 dark:text-dark-gray-300">
          {t('table:results.srm.advanced')}
        </p>
        <Icon
          icon={IconChevronDown}
          size="xxs"
          color="gray-600"
          className={cn(
            'transition-transform duration-200',
            isAdvancedOpen && 'rotate-180'
          )}
        />
      </button>
      {isAdvancedOpen && (
        <div className="flex flex-col gap-y-3">
          <p className="typo-para-small text-gray-700 dark:text-dark-gray-300">
            {t('table:results.srm.ok-summary', {
              pValue: formatPValue(srmResult.pValue)
            })}
          </p>
          <div className="overflow-x-auto overflow-y-hidden">
            <SrmVariationsTable
              variations={srmResult.variations}
              experimentVariations={experiment.variations}
            />
          </div>
        </div>
      )}
    </div>
  );
};

export default SrmBanner;
