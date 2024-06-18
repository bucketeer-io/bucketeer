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

If you want to connect to a real API server, additional settings on the `.env` file are needed.

First, create your own `.env` file based on the `.env.example` file.

```bash
cp .env.example .env
```

Add your API endpoint to the `.env` file like this:

```
DEV_WEB_API_ENDPOINT=https://example.com
DEV_AUTH_REDIRECT_ENDPOINT=http://localhost:8000
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

You need to run the following command when [the message file](https://github.com/bucketeer-io/bucketeer/blob/master/ui/web-v2/src/lang/messages.ts) is modified.

```sh
yarn translate
```

It will generate [en.json](https://github.com/bucketeer-io/bucketeer/blob/master/ui/web-v2/src/assets/lang/en.json). Then, you need manually do the same modifications on the other languages files, including the translation under the same directory.
