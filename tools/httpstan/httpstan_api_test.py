import requests


def compile_model():
    with open("model.stan", "r") as f:
        model_code = f.read()

    data = {
        "program_code": model_code,
    }
    resp = requests.post("http://localhost:8080/v1/models", json=data)
    print(f"compiler_output: {resp.json()['compiler_output']}")
    print(f"stanc_warnings: {resp.json()['stanc_warnings']}")
    print(f"model_name: {resp.json()['name']}")


if __name__ == '__main__':
    compile_model()
