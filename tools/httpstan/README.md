# bucketeer-httpstan

Currently, the calculator service is written in Python, using Bayesian Algorithm to calculate the experiment
probabilities.\
Due to the following reasons, we want to rewrite the service in Go and, for Bayesian analysis, use the `HttpStan`.\
This directory is used to build the `HttpStan` Docker image.

## How to build httpstan image

If you have compiling errors during the image build process, you might need to increase the docker VM memory.

```bash
make docker-build
```

## How to run httpstan image

Use the following command to run the `HttpStan` Docker image on local machine:

```bash
make run-httpstan-container
```

Then you will see a container called `httpstan` running in the background and listening on port: `8080`.

## Test HttpStan API

There is a Python script that can be used to test the HttpStan API.\
The script is `httpstan_api_test.py`, the script contains the following functions:

- `compile_model` - test the `POST /v1/models` endpoint, used to compile the model

Install dependencies:

* `requests` - used to send HTTP requests

```bash
pip install requests
```

Just run the following command to test the HttpStan API locally:

```bash
python httpstan_api_test.py
```

Here are the example results:

* `compile_model`:

```
compiler_output: /root/.cache/httpstan/4.10.1/models/h7wjpoel/model_h7wjpoel.cpp: ... 

stanc_warnings: Warning in '/tmp/httpstan_zdtggaqa/model_h7wjpoel.stan', line 4, column 12: Declaration
    of arrays by placing brackets after a variable name is deprecated and
    will be removed in Stan 2.33.0. Instead use the array keyword before the
    type. This can be changed automatically using the auto-format flag to
    stanc...
    
model_name: models/h7wjpoel
```

As you can see, the `compiler_output` contains the compiler output, the `stanc_warnings` contains the warnings,
and the `model_name` contains the compiled model name.
