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
	"sort"

	"gonum.org/v1/gonum/stat/distuv"

	"github.com/bucketeer-io/bucketeer/v2/proto/eventcounter"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

const (
	// DefaultSRMThreshold is the p-value threshold below which a Sample Ratio
	// Mismatch is flagged. The 0.001 cutoff is the long-standing default in
	// the experimentation literature (see Fabijan et al., "Diagnosing Sample
	// Ratio Mismatch in Online Controlled Experiments", KDD 2019) and is what
	// Eppo, GrowthBook and Microsoft ExP use as their out-of-the-box default.
	DefaultSRMThreshold = 0.001

	// minSRMSampleSize is the smallest total observed user count for which we
	// trust the chi-square approximation. Below this floor the expected
	// per-cell counts can be small enough (≪ 5) that the chi-square's
	// asymptotic distribution misbehaves; report SKIPPED instead of a
	// potentially misleading p-value.
	minSRMSampleSize = 100
)

var (
	errSRMFeatureMissing     = errors.New("feature definition not available")
	errSRMNoDefaultStrategy  = errors.New("feature has no default strategy")
	errSRMNotRolloutStrategy = errors.New(
		"feature default strategy is not a rollout (no per-variation weights to test against)")
	errSRMNoRolloutVariations = errors.New("feature rollout strategy has no variations")
	errSRMAllWeightsZero      = errors.New("all rollout weights are zero")
	errSRMInsufficientSamples = errors.New(
		"total observed users below the minimum required for a reliable chi-square test")
	errSRMTooFewExpectedCells = errors.New("fewer than 2 variations with positive expected user counts")
)

// computeSRM runs a chi-square goodness-of-fit test comparing each variation's
// observed user count (from VariationResult.evaluation_count.user_count) against
// the intended traffic split defined by the feature's default rollout strategy
// weights. The result is always non-nil and always populated with the
// per-variation observed/expected breakdown (when available), so the UI can
// render a meaningful diagnostic in both the OK and SKIPPED cases.
//
// Status semantics:
//   - OK       — p_value >= threshold; observed split matches intended split.
//   - MISMATCH — p_value < threshold; warn the user, results may be invalid.
//   - SKIPPED  — inputs are unusable (no rollout strategy, weights all zero,
//     total sample below the chi-square's reliable floor, or fewer than two
//     cells have positive expected counts). skip_reason explains which.
//
// Known caveats (documented in docs/mechanism/experiment-calculator-math.md):
//   - When the feature has targeting rules or individual overrides, some users
//     may be assigned by rule rather than by the rollout split. The reported
//     SRM will then include rule-matched users and may flag mismatches that
//     reflect targeting rather than a real bucketing bug. MVP intentionally
//     errs on the side of false positives (recoverable) over false negatives
//     (silently invalid experiments).
func computeSRM(
	variationResults []*eventcounter.VariationResult,
	feature *featureproto.Feature,
	threshold float64,
) *eventcounter.SrmResult {
	if threshold <= 0 {
		threshold = DefaultSRMThreshold
	}
	res := &eventcounter.SrmResult{Threshold: threshold}

	weights, err := extractRolloutWeights(feature)
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

	// Build the variation set as the union of (a) the rollout-strategy weights
	// and (b) the observed variation IDs. Iterating only the weights would
	// silently drop any user assigned to a variation that's not in the
	// rollout — which can happen with experiment schema drift (a variation
	// removed from the rollout but stale assignments still in flight) or with
	// genuinely leaked traffic. Those users still belong in totalObserved, and
	// the per-variation breakdown should still surface them (with
	// expected_weight=0) so the UI can show "unknown variation X received N
	// users". Sorting yields a deterministic per-variation order across runs.
	seen := make(map[string]struct{}, len(weights)+len(observedByID))
	vids := make([]string, 0, len(weights)+len(observedByID))
	for vid := range weights {
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
	var totalWeight int64
	perVariation := make([]*eventcounter.SrmVariation, 0, len(vids))
	for _, vid := range vids {
		w := weights[vid] // 0 for variations observed but absent from rollout
		observed := observedByID[vid]
		totalObserved += observed
		totalWeight += w
		perVariation = append(perVariation, &eventcounter.SrmVariation{
			VariationId:       vid,
			ObservedUserCount: observed,
			// ExpectedWeight stored as raw weight here; normalized below
			// once totalWeight is known.
			ExpectedWeight: float64(w),
		})
	}
	res.Variations = perVariation

	if totalWeight <= 0 {
		res.Status = eventcounter.SrmResult_SKIPPED
		res.SkipReason = errSRMAllWeightsZero.Error()
		return res
	}

	for _, v := range perVariation {
		v.ExpectedWeight = v.ExpectedWeight / float64(totalWeight)
		v.ExpectedUserCount = v.ExpectedWeight * float64(totalObserved)
	}

	if totalObserved < minSRMSampleSize {
		res.Status = eventcounter.SrmResult_SKIPPED
		res.SkipReason = fmt.Sprintf("%s (got %d, need >= %d)",
			errSRMInsufficientSamples.Error(), totalObserved, minSRMSampleSize)
		return res
	}

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
	if df < 1 {
		res.Status = eventcounter.SrmResult_SKIPPED
		res.SkipReason = errSRMTooFewExpectedCells.Error()
		return res
	}

	cs := distuv.ChiSquared{K: float64(df)}
	pValue := 1.0 - cs.CDF(chiSq)

	res.ChiSquare = chiSq
	res.DegreesOfFreedom = df
	res.PValue = pValue
	if pValue < threshold {
		res.Status = eventcounter.SrmResult_MISMATCH
	} else {
		res.Status = eventcounter.SrmResult_OK
	}
	return res
}

// extractRolloutWeights returns the intended per-variation traffic weights
// from the feature's default rollout strategy. Returns an error (mapped 1:1
// to an SrmResult SKIPPED reason by computeSRM) when no usable rollout exists.
func extractRolloutWeights(feature *featureproto.Feature) (map[string]int64, error) {
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
	if len(strat.RolloutStrategy.Variations) == 0 {
		return nil, errSRMNoRolloutVariations
	}
	weights := make(map[string]int64, len(strat.RolloutStrategy.Variations))
	for _, v := range strat.RolloutStrategy.Variations {
		weights[v.Variation] = int64(v.Weight)
	}
	return weights, nil
}
