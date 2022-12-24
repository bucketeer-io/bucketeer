# Web Console

This console is built using [Tailwind CSS](https://tailwindcss.com).

## Installation

First, ensure that the node and yarn are installed in your local environment, using the same version configured [here](https://github.com/bucketeer-io/bucketeer/blob/master/WORKSPACE).

### Install dependencies

```sh
yarn install
```

## Local Development

### Set the API and Web endpoint

```sh
export NX_DEV_WEB_API_ENDPOINT=https://example.com
export NX_DEV_AUTH_REDIRECT_ENDPOINT=http://localhost:8000
```

### Serve locally

```sh
yarn start
```

### Build

```sh
yarn build
```

### Lint codes

```sh
yarn lint
```

### Internationalization

You need to run the following command when [the message file](https://github.com/bucketeer-io/bucketeer/blob/master/ui/web-v2/apps/admin/src/lang/messages.ts) is modified.

```sh
yarn translate
```

It will generate [en.json](https://github.com/bucketeer-io/bucketeer/blob/master/ui/web-v2/apps/admin/src/assets/lang/en.json). Then, you need manually do the same modifications on the other languages files, including the translation under the same directory.
