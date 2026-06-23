// Copyright 2026 The Bucketeer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package experimentcalc

import (
	"errors"
	"fmt"
	"math"
	"sort"

	"gonum.org/v1/gonum/stat/distuv"

	"github.com/bucketeer-io/bucketeer/v2/proto/eventcounter"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

const (
	// DefaultSRMThreshold is the p-value threshold below which a Sample Ratio
	// Mismatch is flagged. The 0.001 cutoff is the long-standing default in
	// the experimentation literature (see Fabijan et al., "Diagnosing Sample
	// Ratio Mismatch in Online Controlled Experiments", KDD 2019).
	DefaultSRMThreshold = 0.001

	// minSRMSampleSize is the smallest total observed user count for which we
	// trust the chi-square approximation. Below this floor the expected
	// per-cell counts can be small enough (≪ 5) that the chi-square's
	// asymptotic distribution misbehaves; report SKIPPED instead of a
	// potentially misleading p-value.
	minSRMSampleSize = 100

	// minExpectedCellCount is the textbook reliability floor for the
	// chi-square goodness-of-fit approximation: every expected cell count
	// should be at least 5 (Cochran, 1954). minSRMSampleSize alone is not
	// sufficient — a very skewed rollout (e.g. 99/1 with total=100 gives an
	// expected count of 1 in the small cell) can still violate this floor.
	// When violated we report SKIPPED rather than a p-value the user
	// shouldn't trust.
	minExpectedCellCount = 5.0

	// leakNoiseFloor and leakRateFloor define the zero-expected-cell leak
	// detector that runs alongside the chi-square test. Any variation with
	// expected_weight == 0 (a variation explicitly weighted 0 in the rollout,
	// OR a "leaked" variation that's not in the rollout at all) receiving
	// non-zero observed traffic is, by configuration, traffic going somewhere
	// the SDK should never have routed it. Chi-square can't model this
	// (its (O - E)² / E term is undefined when E = 0), so it would silently
	// pass — a real bucketing-bug class that SRM exists to catch.
	//
	// We trigger MISMATCH when the total zero-expected observed user count
	// strictly exceeds max(leakNoiseFloor, leakRateFloor * totalObserved):
	//   - leakNoiseFloor (= 5) prevents false-positives at small n where a
	//     single misrouted user would otherwise trip the rate floor.
	//   - leakRateFloor (= 0.0001, i.e. 0.01%) scales the threshold for
	//     large experiments. It's intentionally aligned in spirit with the
	//     chi-square p < 0.001 threshold — both say "we're confidently sure
	//     this isn't noise". Real production bucketing bugs (broken hash
	//     function, stale config) typically leak >> 1% of traffic, so this
	//     floor leaves comfortable headroom above realistic noise rates.
	leakNoiseFloor = int64(5)
	leakRateFloor  = 0.0001
)

var (
	errSRMFeatureMissing     = errors.New("feature definition not available")
	errSRMNoDefaultStrategy  = errors.New("feature has no default strategy")
	errSRMNotRolloutStrategy = errors.New(
		"feature default strategy is not a rollout (no per-variation weights to test against)")
	errSRMNoRolloutVariations         = errors.New("feature rollout strategy has no variations")
	errSRMAllWeightsZero              = errors.New("all rollout weights are zero")
	errSRMAudienceDefaultNotInRollout = errors.New(
		"audience default_variation is not one of the rollout variations (cannot compute expected split)")
	errSRMInsufficientSamples = errors.New(
		"total observed users below the minimum required for a reliable chi-square test")
	errSRMTooFewExpectedCells = errors.New("fewer than 2 variations with positive expected user counts")
	errSRMSmallExpectedCell   = errors.New(
		"smallest expected per-variation count below the chi-square reliability floor")
)

// computeSRM compares each variation's observed user count (from
// VariationResult.evaluation_count.user_count) against the intended traffic
// split defined by the feature's default rollout strategy, using two
// detection mechanisms in parallel:
//
//  1. A chi-square goodness-of-fit test over the variations with positive
//     expected counts (the standard SRM test).
//  2. A zero-expected-cell leak detector over the variations with
//     expected_weight == 0 — these are variations the rollout says should
//     receive 0% of traffic, either because they're explicitly weighted 0
//     in the rollout or because they're not in the rollout at all
//     ("leaked"). Chi-square can't include these cells (its denominator
//     would be 0), but any non-trivial observed traffic on a zero-expected
//     variation is by definition a bucketing bug. See leakNoiseFloor /
//     leakRateFloor for the noise thresholds.
//
// Status is MISMATCH when either mechanism fires. The chi-square statistic
// is still reported (when df >= 1) even if the leak detector triggered, so
// the UI can show both signals.
//
// The expected split is audience-aware (see extractExpectedFractions): the
// Audience Traffic Allocation's out-of-audience users are correctly
// attributed to audience.default_variation rather than being treated as a
// bucketing bug.
//
// The result is always non-nil and always populated with the per-variation
// observed/expected breakdown (when available), so the UI can render a
// meaningful diagnostic in both the OK and SKIPPED cases.
//
// Status semantics:
//   - OK       — both chi-square (when applicable) and leak detector pass.
//   - MISMATCH — chi-square's p_value < threshold OR a zero-expected
//     variation received observed traffic above the leak floor.
//   - SKIPPED  — inputs are unusable (no rollout strategy, weights all zero,
//     audience default_variation not in the rollout, total sample below the
//     chi-square's reliable floor, smallest expected cell below the per-cell
//     reliability floor, or fewer than two cells have positive expected
//     counts AND no leak was detected). skip_reason explains which.
//
// Known caveat (documented in docs/mechanism/experiment-calculator-math.md):
//
//	When the feature has rule-based targeting or individual overrides, some
//	users may be assigned by rule rather than by the rollout. The reported
//	SRM will then include rule-matched users and may flag mismatches that
//	reflect targeting rather than a real bucketing bug. The MVP intentionally
//	errs on the side of false positives (recoverable) over false negatives
//	(silently invalid experiments). The Audience-Traffic-Allocation case is
//	handled correctly above — only per-user / per-segment targeting remains
//	a residual caveat.
func computeSRM(
	variationResults []*eventcounter.VariationResult,
	feature *featureproto.Feature,
	threshold float64,
) *eventcounter.SrmResult {
	if threshold <= 0 {
		threshold = DefaultSRMThreshold
	}
	res := &eventcounter.SrmResult{Threshold: threshold}

	fractions, err := extractExpectedFractions(feature)
	if err != nil {
		res.Status = eventcounter.SrmResult_SKIPPED
		res.SkipReason = err.Error()
		return res
	}

	observedByID := make(map[string]int64, len(variationResults))
	for _, vr := range variationResults {
		if vr == nil || vr.EvaluationCount == nil {
			continue
		}
		observedByID[vr.VariationId] = vr.EvaluationCount.UserCount
	}

	// Build the variation set as the union of (a) the rollout-strategy
	// expected fractions and (b) the observed variation IDs. Iterating only
	// the fractions would silently drop any user assigned to a variation
	// that's not in the rollout — which can happen with experiment schema
	// drift (a variation removed from the rollout but stale assignments
	// still in flight) or with genuinely leaked traffic. Those users still
	// belong in totalObserved, and the per-variation breakdown should still
	// surface them (with expected_weight=0) so the UI can show "unknown
	// variation X received N users". Sorting yields a deterministic
	// per-variation order across runs.
	seen := make(map[string]struct{}, len(fractions)+len(observedByID))
	vids := make([]string, 0, len(fractions)+len(observedByID))
	for vid := range fractions {
		if _, ok := seen[vid]; ok {
			continue
		}
		seen[vid] = struct{}{}
		vids = append(vids, vid)
	}
	for vid := range observedByID {
		if _, ok := seen[vid]; ok {
			continue
		}
		seen[vid] = struct{}{}
		vids = append(vids, vid)
	}
	sort.Strings(vids)

	var totalObserved int64
	perVariation := make([]*eventcounter.SrmVariation, 0, len(vids))
	for _, vid := range vids {
		f := fractions[vid] // 0 for variations observed but absent from rollout
		observed := observedByID[vid]
		totalObserved += observed
		perVariation = append(perVariation, &eventcounter.SrmVariation{
			VariationId:       vid,
			ObservedUserCount: observed,
			ExpectedWeight:    f,
		})
	}
	res.Variations = perVariation

	for _, v := range perVariation {
		v.ExpectedUserCount = v.ExpectedWeight * float64(totalObserved)
	}

	if totalObserved < minSRMSampleSize {
		res.Status = eventcounter.SrmResult_SKIPPED
		res.SkipReason = fmt.Sprintf("%s (got %d, need >= %d)",
			errSRMInsufficientSamples.Error(), totalObserved, minSRMSampleSize)
		return res
	}

	// Per-cell reliability floor. We only require expected >= 5 on cells
	// the chi-square sum will actually use — cells with expected == 0
	// (unknown/leaked variations) are excluded from the sum below, so they
	// don't constrain reliability here.
	minExpected := math.Inf(1)
	for _, v := range perVariation {
		if v.ExpectedUserCount > 0 && v.ExpectedUserCount < minExpected {
			minExpected = v.ExpectedUserCount
		}
	}
	if !math.IsInf(minExpected, 1) && minExpected < minExpectedCellCount {
		res.Status = eventcounter.SrmResult_SKIPPED
		res.SkipReason = fmt.Sprintf("%s (smallest expected = %.2f, need >= %.0f)",
			errSRMSmallExpectedCell.Error(), minExpected, minExpectedCellCount)
		return res
	}

	// Zero-expected-cell leak detection. Any variation with expected = 0
	// receiving non-zero observed traffic is, by configuration, traffic
	// going somewhere the SDK should never have routed it. Chi-square
	// can't model this (its (O - E)² / E term is undefined for E = 0),
	// so this check runs independently. We trigger MISMATCH when the
	// total leaked observed count strictly exceeds the noise floor
	// (max of the absolute floor and the rate floor scaled by
	// totalObserved). Strict-greater rather than >= so the threshold
	// itself is treated as still-noise, not as a leak.
	var leakedObserved int64
	for _, v := range perVariation {
		if v.ExpectedUserCount == 0 && v.ObservedUserCount > 0 {
			leakedObserved += v.ObservedUserCount
		}
	}
	leakThreshold := leakNoiseFloor
	if scaled := int64(leakRateFloor * float64(totalObserved)); scaled > leakThreshold {
		leakThreshold = scaled
	}
	leakDetected := leakedObserved > leakThreshold

	// Chi-square goodness-of-fit: Σ (O - E)² / E, summed only over cells with
	// positive expected counts. Each contributing cell adds 1 to df; the
	// final df is K_pos - 1 (one parameter consumed by the totalObserved
	// constraint).
	var chiSq float64
	var posCells int64
	for _, v := range perVariation {
		if v.ExpectedUserCount <= 0 {
			continue
		}
		diff := float64(v.ObservedUserCount) - v.ExpectedUserCount
		chiSq += diff * diff / v.ExpectedUserCount
		posCells++
	}
	df := posCells - 1

	// If chi-square has enough cells to run, report its statistic
	// regardless of which check fires for the final status. This keeps the
	// chi-square / p-value / df fields populated so the UI can show the
	// main-rollout test result alongside any leak signal.
	chiSquareApplicable := df >= 1
	var pValue float64
	if chiSquareApplicable {
		cs := distuv.ChiSquared{K: float64(df)}
		// Use Survival (= P(X > chiSq) = 1 - CDF) directly rather than
		// computing 1 - CDF in user code. For large chi-square values the
		// CDF rounds to exactly 1.0 and the subtraction underflows to 0.0,
		// even though the true tail probability is a tiny non-zero number.
		// MISMATCH would still fire (any p < threshold triggers it), but
		// the p-value surfaced to the API / UI would be a misleading 0.0
		// instead of the real underflow value.
		pValue = cs.Survival(chiSq)
		res.ChiSquare = chiSq
		res.DegreesOfFreedom = df
		res.PValue = pValue
	}

	// Final status. MISMATCH if either check fires. SKIPPED only if neither
	// the chi-square (because there aren't enough cells) nor the leak
	// detector (because the leak didn't exceed the noise floor) has
	// anything to say.
	chiSquareMismatch := chiSquareApplicable && pValue < threshold
	switch {
	case chiSquareMismatch || leakDetected:
		res.Status = eventcounter.SrmResult_MISMATCH
	case chiSquareApplicable:
		res.Status = eventcounter.SrmResult_OK
	default:
		res.Status = eventcounter.SrmResult_SKIPPED
		res.SkipReason = errSRMTooFewExpectedCells.Error()
	}
	return res
}

// extractExpectedFractions returns the audience-adjusted expected fraction
// of total observed traffic for each variation in the feature's default
// rollout strategy. The returned fractions sum to exactly 1.0 (modulo
// floating-point rounding) over the rollout's variations.
//
// Bucketeer's rollout strategy has two independent layers:
//
//  1. Audience Traffic Allocation (audience.percentage, in 1-99): a fraction
//     of users is excluded from the experiment and served
//     audience.default_variation. Per pkg/.../strategy_evaluator.go's
//     rollout() function, those excluded users still emit EvaluationEvents
//     (with the default variation's id), so they count toward the SRM
//     observed user counts.
//  2. Variation Allocation (per-variation weight): the in-audience traffic
//     is split between variations according to these weights.
//
// For each variation V_i with rollout weight w_i (in-audience fraction
// p_i = w_i / Σw) and audience fraction a (= audience.percentage / 100,
// or 1.0 when no audience config or audience.percentage in {0, 100}; see
// strategy_evaluator.go which only filters when 0 < pct < 100):
//
//	expected_fraction(V_i) = a · p_i              if V_i != D
//	expected_fraction(D)   = a · p_D + (1 - a)
//
// where D = audience.default_variation. The sum across all V_i is
// a · Σp + (1 - a) = a + (1 - a) = 1.
//
// When D is empty (audience.default_variation unset), excluded users emit
// no EvaluationEvent — observed counts only contain in-audience users — and
// the raw per-variation weights are the correct expected fractions. We
// model this by treating a as 1.0 for the SRM calculation: it is
// mathematically equivalent.
//
// Returns an error (mapped 1:1 to an SrmResult SKIPPED reason by computeSRM)
// when no usable rollout exists, when all weights are zero, or when
// audience.default_variation is set but is not one of the rollout's
// variations (shouldn't happen per UI validation, but defensive — we
// refuse to compute SRM rather than silently mis-attributing the excluded
// traffic).
func extractExpectedFractions(feature *featureproto.Feature) (map[string]float64, error) {
	if feature == nil {
		return nil, errSRMFeatureMissing
	}
	strat := feature.DefaultStrategy
	if strat == nil {
		return nil, errSRMNoDefaultStrategy
	}
	if strat.Type != featureproto.Strategy_ROLLOUT || strat.RolloutStrategy == nil {
		return nil, errSRMNotRolloutStrategy
	}
	rs := strat.RolloutStrategy
	if len(rs.Variations) == 0 {
		return nil, errSRMNoRolloutVariations
	}

	var totalWeight int64
	for _, v := range rs.Variations {
		totalWeight += int64(v.Weight)
	}
	if totalWeight <= 0 {
		return nil, errSRMAllWeightsZero
	}

	// Determine the audience adjustment. The audience filter is only
	// observable to SRM when ALL three conditions hold:
	//   - audience != nil
	//   - 0 < audience.percentage < 100 (strategy_evaluator.go only
	//     filters in this range; outside it, every user goes through
	//     the variation weights as if no audience were configured)
	//   - audience.default_variation != "" (when empty, rollout()
	//     returns ErrVariationNotFound for excluded users so the SDK
	//     never fires an EvaluationEvent for them; observed counts then
	//     reflect only in-audience users and the raw weights are
	//     correct, equivalent to treating a as 1.0 here)
	// Outside those conditions the adjustment is a no-op (a = 1.0,
	// defaultVariation = ""), which is exactly the pre-audience-aware
	// behavior.
	audienceFraction := 1.0
	defaultVariation := ""
	if a := rs.Audience; a != nil {
		if a.Percentage > 0 && a.Percentage < 100 && a.DefaultVariation != "" {
			audienceFraction = float64(a.Percentage) / 100.0
			defaultVariation = a.DefaultVariation
		}
	}

	// Defensive: if audience.default_variation is set but is not one of the
	// rollout's variations, we can't correctly attribute the out-of-audience
	// users. UI validation should prevent this, but if a misconfigured flag
	// reaches the calculator we want to surface a SKIP rather than silently
	// distort the expected counts.
	if defaultVariation != "" {
		found := false
		for _, v := range rs.Variations {
			if v.Variation == defaultVariation {
				found = true
				break
			}
		}
		if !found {
			return nil, errSRMAudienceDefaultNotInRollout
		}
	}

	fractions := make(map[string]float64, len(rs.Variations))
	for _, v := range rs.Variations {
		pi := float64(v.Weight) / float64(totalWeight)
		ef := audienceFraction * pi
		if v.Variation == defaultVariation {
			ef += 1.0 - audienceFraction
		}
		fractions[v.Variation] = ef
	}
	return fractions, nil
}
