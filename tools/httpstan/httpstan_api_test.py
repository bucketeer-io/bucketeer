import requests


def compile_model():
    model_code = """
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
        """
    data = {
        "program_code": model_code,
    }
    resp = requests.post("http://localhost:8080/v1/models", json=data)
    print(f"compiler_output: {resp.json()['compiler_output']}")
    print(f"stanc_warnings: {resp.json()['stanc_warnings']}")
    print(f"model_name: {resp.json()['name']}")


if __name__ == '__main__':
    compile_model()
