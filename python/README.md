# Bucketer/python

## Requirements

```
brew install poetry
```

## Dependencies

### Calculator

- pystan: depends on gcc(GNU C Compiler), g++(GNU C++ Compiler)

If you use Mac and have already installed XCode, you should have them.

## Development

To initialize venv;

```sh
make init-py
```

To activate venv;

```sh
source $(poetry env info --path)/bin/activate
```

To deactivate venv;

```sh
deactivate
```

To add dependencies;

See [poetry documentation](https://python-poetry.org/docs/).

To export dependencies to requirements.txt;

```sh
make update-deps
```

To test;

```sh
make test
```

To format;

```sh
make fmt
```

To generate codes from proto files;

```sh
make gen
```