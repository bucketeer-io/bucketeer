# Experiment Calculator: Code Implementation

This document explains how Bucketeer implements Bayesian A/B test analysis in code. For the mathematical foundations, see [Experiment Calculator: Mathematical Foundations](./experiment-calculator-math.md).

---

## 1. Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                     ExperimentCalculator                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│   ┌─────────────────┐    ┌─────────────────┐                    │
│   │ Event Counter   │    │ Stan Server     │                    │
│   │ Service         │    │ (httpstan)      │                    │
│   │                 │    │                 │                    │
│   │ - Eval counts   │    │ - Compile model │                    │
│   │ - Goal counts   │    │ - Run MCMC      │                    │
│   │ - Value sums    │    │ - Return samples│                    │
│   └────────┬────────┘    └────────┬────────┘                    │
│            │                      │                             │
│            ▼                      ▼                             │
│   ┌──────────────────────────────────────────────────────┐      │
│   │              createExperimentResult()                │      │
│   │                                                      │      │
│   │  1. Fetch data (eval counts, goal counts)            │      │
│   │  2. Run Bayesian inference (Stan)                    │      │
│   │  3. Calculate statistics                             │      │
│   │  4. Store results                                    │      │
│   └──────────────────────────────────────────────────────┘      │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## 2. The Stan Model

**File**: `pkg/experimentcalculator/stan/experiment.stan`

```stan
data {
    int<lower=0> g;           // Number of variations
    int<lower=0> x[g];        // Conversions per variation
    int<lower=0> n[g];        // Users per variation
}

parameters {
    real<lower=0, upper=1> p[g];  // Conversion rates (what we estimate)
}

model {
    for(i in 1:g){
        x[i] ~ binomial(n[i], p[i]);  // Likelihood
    }
    // Implicit prior: p[i] ~ uniform(0, 1) = Beta(1, 1)
}

generated quantities {
    matrix[g, g] prob_upper;   // Pairwise: is p[i] > p[j]?
    real prob_best[g];         // Is p[i] the best?

    for(i in 1:g){
        real others[g-1];
        others = append_array(p[:i-1], p[i+1:]);
        prob_best[i] = p[i] > max(others) ? 1 : 0;
        for(j in 1:g){
            prob_upper[i, j] = p[i] > p[j] ? 1 : 0;
        }
    }
}
```

**Mapping to Math**:
- `x[i] ~ binomial(n[i], p[i])` → The Binomial likelihood
- `p[i]` bounded to [0,1] → Implicit Beta(1,1) prior
- `prob_best[i]` → Monte Carlo estimate of P(variation i is best)
- `prob_upper[i,j]` → Monte Carlo estimate of P(p[i] > p[j])

## 3. Main Function: `createExperimentResult`

**File**: `pkg/experimentcalculator/experimentcalc/experiment_calculator.go:130`

### Code Walkthrough with Example

Let's trace through with our example (A: 500/10000, B: 600/10000):

```go
func (e ExperimentCalculator) createExperimentResult(
    ctx context.Context,
    envNamespace string,
    experiment *experiment.Experiment,
) (*eventcounter.ExperimentResult, error) {
```

**Step 1: Initialize result**
```go
    experimentResult := &eventcounter.ExperimentResult{
        Id:           experiment.Id,
        ExperimentId: experiment.Id,
        UpdatedAt:    time.Now().Unix(),
    }
```

**Step 2: Collect variation IDs**
```go
    var variationIDs []string
    for _, variation := range experiment.Variations {
        variationIDs = append(variationIDs, variation.Id)
    }
    // variationIDs = ["variation-a-id", "variation-b-id"]
```

**Step 3: Generate timestamps (daily intervals)**
```go
    endAts := listEndAt(experiment.StartAt, experiment.StopAt, time.Now().Unix())
    // endAts = [day1, day2, day3, ..., today]
```

**Step 4: For each goal, fetch data and calculate**
```go
    for _, goalID := range experiment.GoalIds {
        goalResult := &eventcounter.GoalResult{GoalId: goalID}

        for _, timestamp := range endAts {
            // Fetch evaluation counts
            // Math: n_A = 10000, n_B = 10000
            evalVc := e.getEvaluationCount(ctx, ...)

            // Fetch goal counts
            // Math: x_A = 500, x_B = 600
            goalVc := e.getGoalCount(ctx, ...)

            // Run Bayesian calculation
            // Math: Compute posteriors Beta(501,9501) and Beta(601,9401)
            gr := e.calcGoalResult(ctx, evalVc, goalVc, experiment)

            // Append to time series
            e.appendVariationResult(ctx, timestamp, goalResult, gr.VariationResults)
        }
```

**Step 5: Post-processing**
```go
        // Calculate simple conversion rate
        // Math: CVR_A = 500/10000 = 5%, CVR_B = 600/10000 = 6%
        for _, vr := range goalResult.VariationResults {
            vr.ConversionRate = float64(vr.ExperimentCount.UserCount) /
                                float64(vr.EvaluationCount.UserCount) * 100
        }

        // Calculate expected loss
        // Math: E[Loss_A] = 1.0%, E[Loss_B] = 0.02%
        e.calculateExpectedLoss(goalResult.VariationResults)

        // Find best variations (>95% prob of beating baseline)
        e.calculateSummary(ctx, goalResult)
    }
```

## 4. MCMC Sampling: `binomialModelSample`

**File**: `pkg/experimentcalculator/experimentcalc/experiment_calculator.go:536`

```go
func (e ExperimentCalculator) binomialModelSample(
    ctx context.Context,
    vids []string,        // ["variation-a", "variation-b"]
    goalUc, evalUc []int64,  // [500, 600], [10000, 10000]
    baseLineIdx int,
    experiment *experiment.Experiment,
) (map[string]*eventcounter.VariationResult, error) {

    // Run 5 parallel MCMC chains
    numOfChains := 5

    for i := 1; i <= numOfChains; i++ {
        go func(chain int) {
            // Prepare data for Stan
            // Math: g=2, x=[500,600], n=[10000,10000]
            req := stan.CreateFitReq{
                Chain: chain,
                Data: map[string]interface{}{
                    "g": len(goalUc),    // 2 variations
                    "x": goalUc,          // [500, 600] conversions
                    "n": evalUc,          // [10000, 10000] users
                },
                Function:   stan.HmcNUTSFunction,  // HMC-NUTS algorithm
                NumSamples: 21000,                  // 21000 samples
                NumWarmup:  1000,                   // 1000 warmup (discarded)
                RandomSeed: 1234,
            }

            // Submit to Stan server
            fitResp := e.httpStan.CreateFit(ctx, e.stanModelID, req)

            // Wait for MCMC to complete
            for {
                details := e.httpStan.GetOperationDetails(ctx, fitId)
                if details.Done {
                    break
                }
                time.Sleep(50 * time.Millisecond)
            }

            // Get posterior samples
            result := e.httpStan.GetFitResult(ctx, e.stanModelID, fitId)
            samplesChan <- result
        }(i)
    }

    // Combine all chains: 5 × 20000 = 100,000 samples
    samples := collectAll(samplesChan)

    // Convert to statistics
    return e.convertFitSamples(ctx, samples, vids, baseLineIdx)
}
```

**What Stan Returns** (100,000 rows):
```
| p.1    | p.2    | prob_best.1 | prob_best.2 | prob_upper.1.2 | prob_upper.2.1 |
|--------|--------|-------------|-------------|----------------|----------------|
| 0.0498 | 0.0612 | 0           | 1           | 0              | 1              |
| 0.0512 | 0.0589 | 0           | 1           | 0              | 1              |
| 0.0507 | 0.0595 | 0           | 1           | 0              | 1              |
| ...    | ...    | ...         | ...         | ...            | ...            |
```

## 5. Computing Statistics: `statistics.go`

**File**: `pkg/experimentcalculator/experimentcalc/statistics.go`

### Conversion Rate Statistics

```go
func createCvrProb(df dataframe.DataFrame, samples []dataframe.DataFrame, index int) *DistributionSummary {
    col := fmt.Sprintf("p.%d", index)  // "p.1" for variation A
    p := df.Col(col)                    // All 100,000 samples of p_A

    ordered := p.Subset(p.Order(false)).Float()  // Sort for percentiles

    return &eventcounter.DistributionSummary{
        Mean:          p.Mean(),      // Math: E[p_A] = 5.01%
        Sd:            p.StdDev(),    // Math: σ[p_A]
        Median:        p.Median(),    // Math: median(p_A) = 5.01%
        Percentile025: stat.Quantile(0.025, ...),  // Math: 2.5th percentile = 4.58%
        Percentile975: stat.Quantile(0.975, ...),  // Math: 97.5th percentile = 5.44%
        Rhat:          rHat(paramSamples),         // Convergence check
        Histogram:     histogram,
    }
}
```

### Probability of Being Best

```go
func createCvrProbBest(df dataframe.DataFrame, samples []dataframe.DataFrame, index int) *DistributionSummary {
    col := fmt.Sprintf("prob_best.%d", index)  // "prob_best.1"
    probBest := df.Col(col)  // 100,000 values of 0 or 1

    return &eventcounter.DistributionSummary{
        Mean: probBest.Mean(),  // Math: P(A is best) = 0.002 = 0.2%
    }
}
```

### R-hat Calculation

```go
func rHat(samples [][]float64) float64 {
    chains := len(samples)      // 5 chains
    n := len(samples[0]) / 2    // Split each chain in half

    // Calculate mean of each half-chain
    splitChainMean := make([]float64, 2*chains)  // 10 means
    splitChainVar := make([]float64, 2*chains)   // 10 variances

    for chain := 0; chain < chains; chain++ {
        splitChainMean[2*chain] = stat.Mean(samples[chain][:n], nil)
        splitChainMean[2*chain+1] = stat.Mean(samples[chain][n:], nil)
        splitChainVar[2*chain] = stat.Variance(samples[chain][:n], nil)
        splitChainVar[2*chain+1] = stat.Variance(samples[chain][n:], nil)
    }

    // Between-chain variance
    varBetween := float64(n) * stat.Variance(splitChainMean, nil)

    // Within-chain variance (average)
    varWithin := stat.Mean(splitChainVar, nil)

    // R-hat formula
    // Math: R̂ = sqrt((B/W + n-1) / n)
    rhat := math.Sqrt((varBetween/varWithin + float64(n-1)) / float64(n))

    return rhat  // Should be < 1.1
}
```

## 6. Expected Loss: `calculateExpectedLoss`

**File**: `pkg/experimentcalculator/experimentcalc/experiment_calculator.go:843`

```go
func (e ExperimentCalculator) calculateExpectedLoss(variationResults []*VariationResult) {
    numDraws := len(variationResults[0].CvrSamples)  // 100,000 samples
    regretSum := make(map[string]float64)

    // For each posterior sample
    for t := 0; t < numDraws; t++ {
        // Find the best CVR in this sample
        // Math: best_t = max(p_A^t, p_B^t)
        best := variationResults[0].CvrSamples[t]
        for _, vr := range variationResults[1:] {
            if vr.CvrSamples[t] > best {
                best = vr.CvrSamples[t]
            }
        }

        // Accumulate regret
        // Math: loss_i^t = best_t - p_i^t
        for _, vr := range variationResults {
            regretSum[vr.VariationId] += best - vr.CvrSamples[t]
        }
    }

    // Average and convert to percentage
    // Math: E[Loss_i] = (1/N) × Σ loss_i^t × 100%
    for _, vr := range variationResults {
        vr.ExpectedLoss = (regretSum[vr.VariationId] / float64(numDraws)) * 100
    }
}
```

## 7. Normal-Inverse-Gamma: `normal_inverse_gamma.go`

**File**: `pkg/experimentcalculator/experimentcalc/normal_inverse_gamma.go`

For value-based metrics (revenue per user):

```go
// Prior parameters
const (
    priorMean  = 30     // μ₀: prior mean
    priorVar   = 2      // Prior variance
    priorSize  = 20     // κ₀: prior precision weight
    priorAlpha = 10     // α₀: Inverse-Gamma shape
    priorBeta  = 1000   // β₀: Inverse-Gamma rate
)

func calcPosterior(thisN int64, thisMu, thisSigma float64, ...) distr {
    // Effective sample size (log-scaled for stability)
    n2 := math.Log(float64(thisN)) / math.Log(1.1)

    // Posterior mean: weighted average of prior and data
    // Math: μ_n = (κ₀×μ₀ + n×x̄) / (κ₀ + n)
    postMu := (priorNu*priorMu + n2*thisMu) / (priorNu + n2)

    // Posterior precision
    postNu := priorNu + n2

    // Posterior shape
    postAlpha := priorAlpha + (n2 / 2)

    // Posterior rate
    postBeta := priorBeta + 0.5*thisSigma*thisSigma*n2 + ...

    return distr{mu: postMu, nu: postNu, alpha: postAlpha, beta: postBeta}
}

func generateNormalGamma(n int, mu, lambda, alpha, beta float64) []float64 {
    samples := make([]float64, n)

    for i := 0; i < n; i++ {
        // Step 1: Sample precision from Gamma
        // Math: τ ~ Gamma(α, β)
        tau := distuv.Gamma{Alpha: alpha, Beta: beta}.Rand()

        // Step 2: Sample mean from Normal
        // Math: x ~ Normal(μ, 1/(τ×λ))
        samples[i] = mu + math.Sqrt(1/(tau*lambda)) * distuv.Normal{Mu: 0, Sigma: 1}.Rand()
    }

    return samples
}
```

## 8. Complete Code Example

Let's trace through the complete flow with our example:

### Input
```go
experiment := &Experiment{
    Id: "exp-123",
    Variations: [
        {Id: "var-a", Value: "blue"},   // Control
        {Id: "var-b", Value: "green"},  // Treatment
    ],
    BaseVariationId: "var-a",
    GoalIds: ["purchase"],
}

// Data from Event Counter:
evalCounts = {"var-a": 10000, "var-b": 10000}
goalCounts = {"var-a": 500, "var-b": 600}
```

### Stan Input
```json
{
    "g": 2,
    "x": [500, 600],
    "n": [10000, 10000]
}
```

### Stan Output (100,000 samples)
```
p.1 samples:     [0.0498, 0.0512, 0.0507, 0.0495, ...]  → Mean: 0.0501
p.2 samples:     [0.0612, 0.0589, 0.0595, 0.0608, ...]  → Mean: 0.0601
prob_best.1:     [0, 0, 0, 0, ...]                       → Mean: 0.002
prob_best.2:     [1, 1, 1, 1, ...]                       → Mean: 0.998
prob_upper.2.1:  [1, 1, 1, 1, ...]                       → Mean: 0.998
```

### Final Output
```go
ExperimentResult{
    Id: "exp-123",
    GoalResults: [{
        GoalId: "purchase",
        VariationResults: [
            {
                VariationId: "var-a",
                ConversionRate: 5.0,
                CvrProb: {
                    Mean: 0.0501,
                    Median: 0.0501,
                    Percentile025: 0.0458,
                    Percentile975: 0.0544,
                    Rhat: 1.001,
                },
                CvrProbBest: {Mean: 0.002},
                CvrProbBeatBaseline: {Mean: 0.0},  // This IS the baseline
                ExpectedLoss: 1.0,
            },
            {
                VariationId: "var-b",
                ConversionRate: 6.0,
                CvrProb: {
                    Mean: 0.0601,
                    Median: 0.0601,
                    Percentile025: 0.0555,
                    Percentile975: 0.0648,
                    Rhat: 1.001,
                },
                CvrProbBest: {Mean: 0.998},
                CvrProbBeatBaseline: {Mean: 0.998},
                ExpectedLoss: 0.02,
            },
        ],
        Summary: {
            BestVariations: [{Id: "var-b", Probability: 0.998, IsBest: true}],
            GoalUserCount: 1100,
        },
    }],
    TotalEvaluationUserCount: 20000,
}
```

---

## Summary: Math to Code Mapping

| Concept | Math | Code Location |
|---------|------|---------------|
| Bayesian inference | Posterior ∝ Likelihood × Prior | Stan model |
| Binomial likelihood | x ~ Binomial(n, p) | `experiment.stan:13` |
| Beta prior | p ~ Beta(1, 1) | Implicit in Stan |
| MCMC sampling | HMC-NUTS algorithm | `binomialModelSample()` |
| R-hat convergence | R̂ = sqrt((B/W + n-1)/n) | `statistics.go:100` |
| Expected loss | E[max(p) - pᵢ] | `calculateExpectedLoss()` |
| Normal-Inverse-Gamma | For value metrics | `normal_inverse_gamma.go` |
| Credible intervals | 2.5th and 97.5th percentiles | `statistics.go:60-61` |

---

## Related Documents

- [Experiment Calculator: Mathematical Foundations](./experiment-calculator-math.md) - The mathematical theory behind these implementations
