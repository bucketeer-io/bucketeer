# Evaluation module for Node.JS

## Development

### Setup

```sh
export NPM_TOKEN="YOUR_NPM_TOKEN"
make init
make gen_proto
```

### Build

```sh
make build
```

### Unit tests

```sh
make test
```

### Lint

```sh
make lint
```

## Release to NPM

```sh
export NPM_TOKEN="YOUR_NPM_TOKEN"
make init
make gen_proto
make build
make publish
```
