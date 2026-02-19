# Experiment Calculator: Mathematical Foundations

This document explains the mathematical theory behind Bucketeer's Experiment Calculator, which uses Bayesian statistical methods to analyze A/B test results.

---

## 1. The Problem: A/B Testing

When running an A/B test, we have:
- **Variation A (Control)**: The original version
- **Variation B (Treatment)**: The new version we're testing

We observe:
- How many users saw each variation (evaluations)
- How many users completed our goal (conversions)

**Question**: Which variation has the higher true conversion rate?

### Example Scenario

```
Variation A (Control):   10,000 users saw it, 500 converted  →  5.0% observed CVR
Variation B (Treatment): 10,000 users saw it, 600 converted  →  6.0% observed CVR
```

B looks better, but is it *really* better, or just random chance?

---

## 2. Bayesian Inference: The Core Idea

### Frequentist vs Bayesian

| Frequentist Approach | Bayesian Approach |
|---------------------|-------------------|
| "Is the difference statistically significant?" | "What's the probability B is better than A?" |
| Binary answer (yes/no at p < 0.05) | Full probability distribution |
| "Reject null hypothesis" | "95% probability B is better" |

### Bayes' Theorem

The foundation of Bayesian inference:

```
                    P(Data | θ) × P(θ)
P(θ | Data) = ─────────────────────────
                       P(Data)
```

In plain English:

```
Posterior = (Likelihood × Prior) / Evidence
```

Where:
- **Prior P(θ)**: What we believe before seeing data
- **Likelihood P(Data|θ)**: How probable is the data given a parameter value
- **Posterior P(θ|Data)**: What we believe after seeing data

---

## 3. The Binomial-Beta Model

For conversion rate estimation, we use the **Binomial-Beta** model.

### The Setup

For each variation:
- **n** = number of users who saw it (trials)
- **x** = number of users who converted (successes)
- **p** = true conversion rate (unknown - what we want to estimate)

### Step 1: The Likelihood (Binomial Distribution)

The probability of observing x conversions out of n users, given conversion rate p:

```
P(x | n, p) = C(n,x) × p^x × (1-p)^(n-x)
```

Where C(n,x) = n! / (x!(n-x)!) is the binomial coefficient.

**Intuition**: If the true rate is 5%, what's the chance of seeing exactly 500 out of 10,000?

### Step 2: The Prior (Beta Distribution)

We use a Beta distribution as our prior belief about p:

```
P(p) = Beta(α, β) ∝ p^(α-1) × (1-p)^(β-1)
```

Special cases:
- **Beta(1, 1)** = Uniform distribution (no prior knowledge)
- **Beta(2, 2)** = Slight belief that p is around 0.5
- **Beta(10, 90)** = Prior belief that p ≈ 10%

**We use Beta(1,1) - uniform prior** (no prior assumptions).

### Step 3: The Posterior (Also Beta!)

The magic of **conjugate priors**: When likelihood is Binomial and prior is Beta, the posterior is also Beta:

```
Prior:      Beta(α, β)
+ Data:     x successes, n-x failures
────────────────────────────────────
Posterior:  Beta(α + x, β + n - x)
```

### Example Calculation

**Variation A**: 500 conversions out of 10,000 users

```
Prior:     Beta(1, 1)
+ Data:    500 successes, 9,500 failures
─────────────────────────────────────────
Posterior: Beta(1 + 500, 1 + 9500) = Beta(501, 9501)
```

**Variation B**: 600 conversions out of 10,000 users

```
Prior:     Beta(1, 1)
+ Data:    600 successes, 9,400 failures
─────────────────────────────────────────
Posterior: Beta(1 + 600, 1 + 9400) = Beta(601, 9401)
```

### Posterior Statistics

For a Beta(α, β) distribution:

```
Mean     = α / (α + β)
Variance = αβ / ((α + β)² × (α + β + 1))
```

**Variation A**: Beta(501, 9501)
```
Mean = 501 / (501 + 9501) = 501 / 10002 = 0.0501 = 5.01%
```

**Variation B**: Beta(601, 9401)
```
Mean = 601 / (601 + 9401) = 601 / 10002 = 0.0601 = 6.01%
```

---

## 4. Comparing Variations: Who's the Winner?

Now we have two posterior distributions:
- **p_A ~ Beta(501, 9501)**
- **p_B ~ Beta(601, 9401)**

### Question: What's P(p_B > p_A)?

This is where **Monte Carlo simulation** comes in.

### Monte Carlo Method

1. Draw a sample from p_A's posterior
2. Draw a sample from p_B's posterior
3. Check if p_B > p_A
4. Repeat 100,000 times
5. Count how often p_B > p_A

```
Sample 1: p_A = 0.0498, p_B = 0.0612 → B wins ✓
Sample 2: p_A = 0.0512, p_B = 0.0589 → B wins ✓
Sample 3: p_A = 0.0507, p_B = 0.0595 → B wins ✓
...
Sample 100,000: p_A = 0.0495, p_B = 0.0608 → B wins ✓

Result: B won 99,800 out of 100,000 times
P(B > A) = 99.8%
```

### Visual Representation

```
Posterior Distributions:

        Variation A              Variation B
        Beta(501, 9501)          Beta(601, 9401)

              ▲                        ▲
             /|\                      /|\
            / | \                    / | \
           /  |  \                  /  |  \
          /   |   \                /   |   \
    ─────/────|────\──────────────/────|────\─────
        4%   5%   6%            5%   6%   7%

        ← Almost no overlap →

    P(B > A) ≈ 99.8%
```

---

## 5. MCMC Sampling

For simple Beta distributions, we can sample directly. But for more complex models, we need **Markov Chain Monte Carlo (MCMC)**.

### Why MCMC?

When we have multiple variations and want to compute things like:
- P(variation i is THE BEST among all)
- Pairwise comparisons between all variations

Direct sampling becomes complex. MCMC explores the posterior systematically.

### How MCMC Works (Conceptually)

```
1. Start at some random parameter values
2. Propose a move to new values
3. Accept or reject based on posterior probability
4. Repeat thousands of times
5. The samples approximate the posterior distribution
```

### HMC-NUTS (What Bucketeer Uses)

**Hamiltonian Monte Carlo with No-U-Turn Sampler** is an efficient MCMC algorithm:

- Uses gradient information to make smarter proposals
- Automatically tunes step sizes
- Much faster than basic MCMC

### Multiple Chains

We run **5 parallel chains** to:
- Check convergence (chains should agree)
- Get more samples efficiently

```
Chain 1: ●●●●●●●●●● (20,000 samples)
Chain 2: ●●●●●●●●●● (20,000 samples)
Chain 3: ●●●●●●●●●● (20,000 samples)
Chain 4: ●●●●●●●●●● (20,000 samples)
Chain 5: ●●●●●●●●●● (20,000 samples)
────────────────────────────────────
Total:   100,000 samples
```

---

## 6. R-hat: Convergence Diagnostic

How do we know MCMC has converged to the true posterior?

### The Idea

If all chains have converged, they should be sampling from the same distribution. So:
- **Within-chain variance** ≈ **Between-chain variance**

### The Formula

```
R̂ = sqrt( (Between-chain variance + Within-chain variance) / Within-chain variance )
```

Simplified:
```
R̂ ≈ sqrt( (B/W + n-1) / n )
```

### Interpretation

| R̂ Value | Meaning |
|---------|---------|
| R̂ ≈ 1.0 | Perfect convergence |
| R̂ < 1.01 | Excellent |
| R̂ < 1.1 | Acceptable |
| R̂ > 1.1 | **Not converged** - need more samples |

### Example

```
Chain 1 mean: 0.0502
Chain 2 mean: 0.0499
Chain 3 mean: 0.0501
Chain 4 mean: 0.0500
Chain 5 mean: 0.0503

Between-chain variance: very small (all ≈ 0.050)
Within-chain variance: similar to between

R̂ ≈ 1.001 ✓ Converged!
```

---

## 7. Expected Loss (Regret)

Expected loss answers: **"How much do I lose if I pick the wrong variation?"**

### Definition

For variation i:
```
Expected Loss(i) = E[ max(p_1, p_2, ..., p_k) - p_i ]
```

The average difference between the best possible and this variation.

### Monte Carlo Calculation

```
For each of 100,000 samples:
    best = max(p_A, p_B)
    loss_A = best - p_A
    loss_B = best - p_B

Expected Loss(A) = average(loss_A) × 100%
Expected Loss(B) = average(loss_B) × 100%
```

### Example

```
Sample 1: p_A=0.050, p_B=0.060 → best=0.060 → loss_A=0.010, loss_B=0.000
Sample 2: p_A=0.052, p_B=0.058 → best=0.058 → loss_A=0.006, loss_B=0.000
Sample 3: p_A=0.054, p_B=0.055 → best=0.055 → loss_A=0.001, loss_B=0.000
Sample 4: p_A=0.051, p_B=0.059 → best=0.059 → loss_A=0.008, loss_B=0.000
...

Average loss_A = 0.010 → Expected Loss(A) = 1.0%
Average loss_B = 0.0002 → Expected Loss(B) = 0.02%
```

### Interpretation

- **Expected Loss(A) = 1.0%**: Choosing A costs you ~1 percentage point on average
- **Expected Loss(B) = 0.02%**: Choosing B costs you almost nothing

**Decision**: Pick B (lowest expected loss)

---

## 8. Normal-Inverse-Gamma: For Value Metrics

When tracking **revenue per user** (not just conversion), we need a different model.

### The Setup

- We observe values: $10, $25, $15, $30, ...
- We want to estimate the **mean value per user**
- We also need to estimate **variance** (uncertainty in individual values)

### The Model

```
Data:    x_i ~ Normal(μ, σ²)    (values are normally distributed)
Unknown: μ (mean) and σ² (variance)
```

### Prior: Normal-Inverse-Gamma

```
σ² ~ Inverse-Gamma(α, β)
μ | σ² ~ Normal(μ₀, σ²/κ)
```

This is a **conjugate prior** for the Normal likelihood.

### Posterior Updates

After observing data with sample mean x̄ and sample variance s²:

```
κ_n = κ₀ + n
μ_n = (κ₀×μ₀ + n×x̄) / κ_n
α_n = α₀ + n/2
β_n = β₀ + ½×Σ(x_i - x̄)² + (κ₀×n×(x̄ - μ₀)²) / (2×κ_n)
```

### Sampling

```
1. Sample precision: τ ~ Gamma(α_n, β_n)
2. Convert to variance: σ² = 1/τ
3. Sample mean: μ ~ Normal(μ_n, σ²/κ_n)
```

---

## 9. Credible Intervals

A **95% credible interval** means:

```
"There is a 95% probability that the true parameter is in this interval."
```

This is different from frequentist confidence intervals!

### Calculation

From posterior samples, take the 2.5th and 97.5th percentiles:

```
Samples: [0.048, 0.049, 0.050, 0.051, 0.052, 0.053, ...]
         (sorted)

2.5th percentile:  0.046
97.5th percentile: 0.054

95% Credible Interval: [4.6%, 5.4%]
```

---

## 10. Complete Mathematical Example

### Scenario

A/B test for a "Buy Now" button:

| Variation | Users | Conversions |
|-----------|-------|-------------|
| A (Blue button) | 10,000 | 500 |
| B (Green button) | 10,000 | 600 |

### Step 1: Define Priors

Using non-informative priors:
```
p_A ~ Beta(1, 1)
p_B ~ Beta(1, 1)
```

### Step 2: Calculate Posteriors

```
p_A | data ~ Beta(1 + 500, 1 + 9500) = Beta(501, 9501)
p_B | data ~ Beta(1 + 600, 1 + 9400) = Beta(601, 9401)
```

### Step 3: Compute Statistics

**Variation A** - Beta(501, 9501):
```
Mean   = 501/10002 = 5.01%
Median ≈ 5.01%
95% CI = [4.58%, 5.44%]
```

**Variation B** - Beta(601, 9401):
```
Mean   = 601/10002 = 6.01%
Median ≈ 6.01%
95% CI = [5.55%, 6.48%]
```

### Step 4: Compare (Monte Carlo)

Draw 100,000 samples from each posterior:

```
Comparisons:
- P(B > A) = 99.8%
- P(B is best) = 99.8%
- P(A is best) = 0.2%
```

### Step 5: Calculate Expected Loss

```
Expected Loss(A) = 1.00%  (choosing A costs ~1 percentage point)
Expected Loss(B) = 0.02%  (choosing B costs almost nothing)
```

### Step 6: Final Results

| Metric | Variation A | Variation B |
|--------|-------------|-------------|
| Observed CVR | 5.00% | 6.00% |
| Posterior Mean | 5.01% | 6.01% |
| 95% Credible Interval | [4.58%, 5.44%] | [5.55%, 6.48%] |
| Probability of Being Best | 0.2% | 99.8% |
| Expected Loss | 1.00% | 0.02% |

**Conclusion**: Variation B (green button) wins with 99.8% probability.

---

## Related Documents

- [Experiment Calculator: Code Implementation](./experiment-calculator-code.md) - How these concepts are implemented in Bucketeer's codebase
