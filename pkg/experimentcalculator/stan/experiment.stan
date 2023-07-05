data {
    int<lower=0> g;
    int<lower=0> x[g];
    int<lower=0> n[g];
}

parameters {
    real<lower=0, upper=1> p[g];
}

model {
    for(i in 1:g){
        x[i] ~ binomial(n[i], p[i]);
    }
}

generated quantities {
    matrix[g, g] prob_upper;
    real prob_best[g];

    for(i in 1:g){
        real others[g-1];
        others = append_array(p[:i-1], p[i+1:]);
        prob_best[i] = p[i] > max(others) ? 1 : 0;
        for(j in 1:g){
            prob_upper[i, j] = p[i] > p[j] ? 1 : 0;
        }
    }
}