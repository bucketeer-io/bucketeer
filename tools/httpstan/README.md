# bucketeer-httpstan

Currently, the calculator service is written in Python, using Bayesian Algorithm to calculate the experiment
probabilities.\
Due to the following reasons, we want to rewrite the service in Go and, for Bayesian analysis, use the `HttpStan`.\
This directory is used to build the `HttpStan` Docker image.

## How to build httpstan image

```bash
make docker-build
```

## How to run httpstan image

```bash
docker run -p 8080:8080 -it --rm --name httpstan bucketeer-httpstan:latest
```

## Test HttpStan API

There is a Python script that can be used to test the HttpStan API.\
The script is `httpstan_api_test.py`, the script contains the following functions:

- `compile_model` - test the `POST /v1/models` endpoint, used to compile the model

Just run the following command to test the HttpStan API locally:

```bash
python httpstan_api_test.py
```