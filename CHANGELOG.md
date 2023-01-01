# Changelog

## [0.2.0](https://github.com/bucketeer-io/bucketeer/compare/v0.1.0...v0.2.0) (2023-01-01)


### Features

* add track api to grpc server ([#45](https://github.com/bucketeer-io/bucketeer/issues/45)) ([f7cb80b](https://github.com/bucketeer-io/bucketeer/commit/f7cb80ba506bcb6e507d0bebed1ffd21c9e98450))
* **ui:** add the menu to change the language in the dashboard ([#135](https://github.com/bucketeer-io/bucketeer/issues/135)) ([36be7b7](https://github.com/bucketeer-io/bucketeer/commit/36be7b753d4ac14bf97bcbebed935327d1d60b09))


### Bug Fixes

* error handling condition in the kafka writer ([#116](https://github.com/bucketeer-io/bucketeer/issues/116)) ([207f46a](https://github.com/bucketeer-io/bucketeer/commit/207f46a0351ff06474c747ddaf594af7bcd5d5b6))
* error handling while trying to send messages to kafka ([#115](https://github.com/bucketeer-io/bucketeer/issues/115)) ([34f616b](https://github.com/bucketeer-io/bucketeer/commit/34f616b05f042fb4157717f134a9eea71e12bcc4))
* evaluation unique count is always zero ([#127](https://github.com/bucketeer-io/bucketeer/issues/127)) ([7e1a455](https://github.com/bucketeer-io/bucketeer/commit/7e1a45529500661c1b3e0fd4ab8ee557850a260c))
* event user persister should not ack message when DB returns an error ([#84](https://github.com/bucketeer-io/bucketeer/issues/84)) ([39eb579](https://github.com/bucketeer-io/bucketeer/commit/39eb5797dfeecf50025e3c235891a499cc0f8a81))
* failed to find feature while trying to update an archived feature ([#53](https://github.com/bucketeer-io/bucketeer/issues/53)) ([e4b9e0d](https://github.com/bucketeer-io/bucketeer/commit/e4b9e0dd55aa015d0be0238a630f9f9793cfa7c9))
* flush size setting being ignored in the event persister ([#117](https://github.com/bucketeer-io/bucketeer/issues/117)) ([a19af75](https://github.com/bucketeer-io/bucketeer/commit/a19af751a2d66e942e29dca77e09b5957e9e15fb))
* internal error when there is no value in the event count ([#137](https://github.com/bucketeer-io/bucketeer/issues/137)) ([b3b5b15](https://github.com/bucketeer-io/bucketeer/commit/b3b5b15b1cb87a3713d817683c6f6ff1337be20f))
* json cannot unmarshal array while trying to create a feature flag ([#15](https://github.com/bucketeer-io/bucketeer/issues/15)) ([142b117](https://github.com/bucketeer-io/bucketeer/commit/142b11742e177f064ae9c7ab98fad5de62faf851))
* redis key expiration is not being set correctly ([#118](https://github.com/bucketeer-io/bucketeer/issues/118)) ([2cf3bfa](https://github.com/bucketeer-io/bucketeer/commit/2cf3bfad5f01e8d328b5a428b474e6937b0fce66))
* table has no partition for value from column_list error ([#144](https://github.com/bucketeer-io/bucketeer/issues/144)) ([d70dcb8](https://github.com/bucketeer-io/bucketeer/commit/d70dcb851e408340c9a296cb4524a23ee258f4eb))
* the event may duplicate in the kafka if the redis request fails ([#128](https://github.com/bucketeer-io/bucketeer/issues/128)) ([393afe7](https://github.com/bucketeer-io/bucketeer/commit/393afe730e96462ef1bf58f319b529e081c818f8))
* **ui/web-v2:** feature flag name not being validated ([#16](https://github.com/bucketeer-io/bucketeer/issues/16)) ([6118f38](https://github.com/bucketeer-io/bucketeer/commit/6118f38507dd3911e5e8493bcd5e463821fbfdfb))


### Build System

* add include_imports flag to protoc ([#57](https://github.com/bucketeer-io/bucketeer/issues/57)) ([6ab4f11](https://github.com/bucketeer-io/bucketeer/commit/6ab4f11b37ae066ab83afde35d7817592c159a94))
* add rules to build and test go files using go command ([#23](https://github.com/bucketeer-io/bucketeer/issues/23)) ([399fd6d](https://github.com/bucketeer-io/bucketeer/commit/399fd6d51432b14fb0e90aee16df91182f0a560b))
* change protoc to use external dependency files ([#25](https://github.com/bucketeer-io/bucketeer/issues/25)) ([56894fe](https://github.com/bucketeer-io/bucketeer/commit/56894fe54bc26efd93f800eb63f27e80cc4ffecb))
* fix not found error while generating proto web files ([#27](https://github.com/bucketeer-io/bucketeer/issues/27)) ([0471f3b](https://github.com/bucketeer-io/bucketeer/commit/0471f3ba2c5382919669d85017a524a2770d9b6a))
* generate proto_descriptor using protoc ([#46](https://github.com/bucketeer-io/bucketeer/issues/46)) ([a8a5f1d](https://github.com/bucketeer-io/bucketeer/commit/a8a5f1dee8526b5454c0e190e3212551ea898357))
* remove bazelisk command from Makefile ([#59](https://github.com/bucketeer-io/bucketeer/issues/59)) ([3d3495b](https://github.com/bucketeer-io/bucketeer/commit/3d3495b1dc314f434e2898b5cdbbf96e56c3d2f3))
* remove go_embed_data rule ([#19](https://github.com/bucketeer-io/bucketeer/issues/19)) ([b0efa8f](https://github.com/bucketeer-io/bucketeer/commit/b0efa8f0fcc84289d10fc094b0e1699500637fb7))
* set ubuntu-20.04 for python build runner ([#64](https://github.com/bucketeer-io/bucketeer/issues/64)) ([a1c9ea0](https://github.com/bucketeer-io/bucketeer/commit/a1c9ea047bf1b320d78618ca24c93f3055b7bf64))
* setup buildifier via bazel ([#8](https://github.com/bucketeer-io/bucketeer/issues/8)) ([ab61d14](https://github.com/bucketeer-io/bucketeer/commit/ab61d149caaa05dec992c136a239024d790ac1e2))
* update renovate settings to check once a month ([#5](https://github.com/bucketeer-io/bucketeer/issues/5)) ([01ddce4](https://github.com/bucketeer-io/bucketeer/commit/01ddce40f9f7384549cde0b0ee8abc328bc6f71e))


### Miscellaneous

* add auto ops service in the event persister envoy settings ([#103](https://github.com/bucketeer-io/bucketeer/issues/103)) ([0753c7d](https://github.com/bucketeer-io/bucketeer/commit/0753c7d343567fb7c3b2295c8902ce432c6bac37))
* add default variation id for evaluation count key ([#111](https://github.com/bucketeer-io/bucketeer/issues/111)) ([f7900d1](https://github.com/bucketeer-io/bucketeer/commit/f7900d114f005cb1b65cce40d7b6c4e26c90f752))
* add env variables for postgres ([#65](https://github.com/bucketeer-io/bucketeer/issues/65)) ([d9f403c](https://github.com/bucketeer-io/bucketeer/commit/d9f403c5f4263679f5947a79f888a6da03d9b18a))
* add experiment client in the event-persister service ([#89](https://github.com/bucketeer-io/bucketeer/issues/89)) ([e231c03](https://github.com/bucketeer-io/bucketeer/commit/e231c036b7f96bc0baa1e6b16a630b48069c49d8))
* add experiment service's envoy config in the event persister ([#80](https://github.com/bucketeer-io/bucketeer/issues/80)) ([54d8d34](https://github.com/bucketeer-io/bucketeer/commit/54d8d3491b962ec2815c887b6abb1781dc090124))
* add expiration to keys for evaluation counts ([#100](https://github.com/bucketeer-io/bucketeer/issues/100)) ([6c9920f](https://github.com/bucketeer-io/bucketeer/commit/6c9920ffeb558ca0b79a0443ff82ead74bb16155))
* add get user evaluation interface implementation ([#90](https://github.com/bucketeer-io/bucketeer/issues/90)) ([e68dcc2](https://github.com/bucketeer-io/bucketeer/commit/e68dcc28dfe98117183e66577b10568603c1ab69))
* add handling for new metrics events in the persister service ([#31](https://github.com/bucketeer-io/bucketeer/issues/31)) ([6fc3419](https://github.com/bucketeer-io/bucketeer/commit/6fc3419adcfa116d483728656d9e59c5bdc99c61))
* add json transcoding to track api for testing ([#37](https://github.com/bucketeer-io/bucketeer/issues/37)) ([876fdaf](https://github.com/bucketeer-io/bucketeer/commit/876fdaf60a4e30bdf70d35cdf0311d3b46811d07))
* add metadata property to proto client events ([#34](https://github.com/bucketeer-io/bucketeer/issues/34)) ([9cabd9f](https://github.com/bucketeer-io/bucketeer/commit/9cabd9f8d90c14af9f0bee18508ae431ce225781))
* add new sdk metrics event types to proto ([#29](https://github.com/bucketeer-io/bucketeer/issues/29)) ([5d173d5](https://github.com/bucketeer-io/bucketeer/commit/5d173d58b1a30cd735ea51f07f2389e12405b014))
* add postgresClient to persister ([#73](https://github.com/bucketeer-io/bucketeer/issues/73)) ([ff105f2](https://github.com/bucketeer-io/bucketeer/commit/ff105f28185883f9440be90888d0161a344a0787))
* add redis pipeline metrics ([#120](https://github.com/bucketeer-io/bucketeer/issues/120)) ([610e07d](https://github.com/bucketeer-io/bucketeer/commit/610e07db41e88d6e071223c45441ac1defddd89d))
* add redis settings in the event persister service ([#92](https://github.com/bucketeer-io/bucketeer/issues/92)) ([bbf62ef](https://github.com/bucketeer-io/bucketeer/commit/bbf62ef02d769a5db3601264f4b9712878d14c9b))
* add sdk version property to proto metrics events ([#32](https://github.com/bucketeer-io/bucketeer/issues/32)) ([2548847](https://github.com/bucketeer-io/bucketeer/commit/2548847f3980d2f8bde908a4ab4cb18e0af91c0c))
* add the ability to handle new metrics event in gateway ([#33](https://github.com/bucketeer-io/bucketeer/issues/33)) ([f29ab67](https://github.com/bucketeer-io/bucketeer/commit/f29ab676033d84cbaabc71de90b09e7bcd140700))
* add the metadata property to metrics event proto message ([#132](https://github.com/bucketeer-io/bucketeer/issues/132)) ([af575de](https://github.com/bucketeer-io/bucketeer/commit/af575deb21ece93700d93fd09552e05178cadfe6))
* add upsert evaluation event in the persister event service ([#88](https://github.com/bucketeer-io/bucketeer/issues/88)) ([ce1f4e0](https://github.com/bucketeer-io/bucketeer/commit/ce1f4e038fa8c6970a9c9a90a5edd15fc234711f))
* change event-counter timeout to 3 hours temporarily ([#86](https://github.com/bucketeer-io/bucketeer/issues/86)) ([c676599](https://github.com/bucketeer-io/bucketeer/commit/c676599ca08c487ff825e8bdafe70cc881018364))
* change experiment batch cronjob ([#126](https://github.com/bucketeer-io/bucketeer/issues/126)) ([9cd856a](https://github.com/bucketeer-io/bucketeer/commit/9cd856af6572b982d0eb8adfff854e3b3e8a72f7))
* change experiment updater cronjob settings ([#55](https://github.com/bucketeer-io/bucketeer/issues/55)) ([9a8a4d0](https://github.com/bucketeer-io/bucketeer/commit/9a8a4d0b3adb5ad50110d0a9a29ff6980e499588))
* change grpc track api name temporarily for testing ([#39](https://github.com/bucketeer-io/bucketeer/issues/39)) ([03c626d](https://github.com/bucketeer-io/bucketeer/commit/03c626dd960b9231bcaf7109a147ad1910f77d21))
* change mau count api in the notification sender ([#136](https://github.com/bucketeer-io/bucketeer/issues/136)) ([5ca4170](https://github.com/bucketeer-io/bucketeer/commit/5ca4170ff9fdbc4419b2ca238ee44e53a616c7d4))
* change the timezone to JP location when saving the count in redis ([#130](https://github.com/bucketeer-io/bucketeer/issues/130)) ([fd8bd2e](https://github.com/bucketeer-io/bucketeer/commit/fd8bd2ebdc966abe8ac94cac8e4caa2a9482099f))
* check for unsent events in the event persister before shutting down the service ([#123](https://github.com/bucketeer-io/bucketeer/issues/123)) ([1c0cd2c](https://github.com/bucketeer-io/bucketeer/commit/1c0cd2c3bebffd24c65bfef3cf1820e0d233c01f))
* configure delete-e2e-data's Makefile to build docker image ([#51](https://github.com/bucketeer-io/bucketeer/issues/51)) ([0aba9f3](https://github.com/bucketeer-io/bucketeer/commit/0aba9f31f4f52d4ba9aa762d6762c9649451a1eb))
* configure renovate for automerge and assigning reviewers ([#71](https://github.com/bucketeer-io/bucketeer/issues/71)) ([7a3429b](https://github.com/bucketeer-io/bucketeer/commit/7a3429b83c13d6f48482982cf87b858f34aa5cf3))
* display variation name in the variation label ([#68](https://github.com/bucketeer-io/bucketeer/issues/68)) ([8c56897](https://github.com/bucketeer-io/bucketeer/commit/8c568977d6b22c20ca916f8e61b5bb144e1c2da9))
* link goal event to the auto ops before sending it to kafka ([#104](https://github.com/bucketeer-io/bucketeer/issues/104)) ([bebc795](https://github.com/bucketeer-io/bucketeer/commit/bebc795a4ba227539f0e04a0babe94e77d27bd11))
* link the goal event to the experiment before sending it to kafka ([#101](https://github.com/bucketeer-io/bucketeer/issues/101)) ([015f326](https://github.com/bucketeer-io/bucketeer/commit/015f3268cd3f3170832eba89a39143bbae103fb5))
* remove bazel config files ([#67](https://github.com/bucketeer-io/bucketeer/issues/67)) ([efb48ea](https://github.com/bucketeer-io/bucketeer/commit/efb48ea66bd2a142a66d6735f0b78b737a44f170))
* remove duplicated environment vars ([#114](https://github.com/bucketeer-io/bucketeer/issues/114)) ([6dd9801](https://github.com/bucketeer-io/bucketeer/commit/6dd980143ef6d1d4001b76e0caffa67779584613))
* remove expiration setting from the event and user count ([#121](https://github.com/bucketeer-io/bucketeer/issues/121)) ([b18232f](https://github.com/bucketeer-io/bucketeer/commit/b18232f2fceca822bcf990a1d86dc5f1f8f26e87))
* remove postgreSQL test implementation from event persister ([#96](https://github.com/bucketeer-io/bucketeer/issues/96)) ([e9e51cd](https://github.com/bucketeer-io/bucketeer/commit/e9e51cd6f7bd92c049747330262b2cd0740ee288))
* remove unnecessary health check in the envoy egress settings ([#85](https://github.com/bucketeer-io/bucketeer/issues/85)) ([5140c11](https://github.com/bucketeer-io/bucketeer/commit/5140c11e5e8395bb27a9be0f08749897533a216a))
* remove unnecessary health checks in the envoy settings ([#87](https://github.com/bucketeer-io/bucketeer/issues/87)) ([77474c2](https://github.com/bucketeer-io/bucketeer/commit/77474c2339075db4acfca2491a36b6534885c9e5))
* remove whitespaces from the tag before upserting it ([#35](https://github.com/bucketeer-io/bucketeer/issues/35)) ([7871b97](https://github.com/bucketeer-io/bucketeer/commit/7871b972a29a707c85d295e8f9305bd4fb3834df))
* set circuit break for api-gateway ([#79](https://github.com/bucketeer-io/bucketeer/issues/79)) ([45d1363](https://github.com/bucketeer-io/bucketeer/commit/45d136301f70e7feb15d1ed1e43d637b673cfedd))
* set version using ldflags ([#48](https://github.com/bucketeer-io/bucketeer/issues/48)) ([1c9cefb](https://github.com/bucketeer-io/bucketeer/commit/1c9cefbebb02c5f96e870b9223299b7e7c06e18a))
* stop inserting events into postgres ([#76](https://github.com/bucketeer-io/bucketeer/issues/76)) ([6eea130](https://github.com/bucketeer-io/bucketeer/commit/6eea1302f1d1f0c00aae23c81725d637ccfed5e4))
* store evaluation count in redis ([#91](https://github.com/bucketeer-io/bucketeer/issues/91)) ([105b4da](https://github.com/bucketeer-io/bucketeer/commit/105b4da5b96a5baea79d2496c7d1a35a75b7f266))
* store evaluation events to postgresql ([#63](https://github.com/bucketeer-io/bucketeer/issues/63)) ([4c82b31](https://github.com/bucketeer-io/bucketeer/commit/4c82b319fef6ac219e59067ba010458a95f8bb13))
* support aws kms ([#62](https://github.com/bucketeer-io/bucketeer/issues/62)) ([77b1ae6](https://github.com/bucketeer-io/bucketeer/commit/77b1ae6515246d58852532ee3e9259b6075a8fd7))
* support sdk version in the metrics ([#54](https://github.com/bucketeer-io/bucketeer/issues/54)) ([16cb007](https://github.com/bucketeer-io/bucketeer/commit/16cb0072b9afa428ebb23d2e9b6421f487c55d41))
* **ui/web-v2:** change the local development server endpoint ([#131](https://github.com/bucketeer-io/bucketeer/issues/131)) ([d273656](https://github.com/bucketeer-io/bucketeer/commit/d273656e240a6fc6de34c6a311461eca5a97d8e3))
* update api-gateway proto descriptor value ([#43](https://github.com/bucketeer-io/bucketeer/issues/43)) ([0c2d619](https://github.com/bucketeer-io/bucketeer/commit/0c2d619e71ef522f104c2d561148d00ab025a580))
* update envoy gateway descriptor ([#42](https://github.com/bucketeer-io/bucketeer/issues/42)) ([0f95e97](https://github.com/bucketeer-io/bucketeer/commit/0f95e97dd4c74dda49e9f6d113fa49fc1c07bc58))
* update eventpersister to store mau to mysql ([#81](https://github.com/bucketeer-io/bucketeer/issues/81)) ([571cf44](https://github.com/bucketeer-io/bucketeer/commit/571cf448dfc794ca5bc2308ef1b45d1567493fc9))
* update ingress api version ([#99](https://github.com/bucketeer-io/bucketeer/issues/99)) ([a33aa7c](https://github.com/bucketeer-io/bucketeer/commit/a33aa7c34cd2a9971adeed356f504786bd4b7be2))
* update redis default settings ([#125](https://github.com/bucketeer-io/bucketeer/issues/125)) ([405e495](https://github.com/bucketeer-io/bucketeer/commit/405e495ee00700d374ff488a94035a39cd5ba442))
* update test runner image ([#38](https://github.com/bucketeer-io/bucketeer/issues/38)) ([8cd8db0](https://github.com/bucketeer-io/bucketeer/commit/8cd8db0166aa0fd2d846dbc0871cf4131e48589d))
* use redis instead of druid in GetEvaluationTimeseriesCount ([#122](https://github.com/bucketeer-io/bucketeer/issues/122)) ([517065e](https://github.com/bucketeer-io/bucketeer/commit/517065e6e6a2a799e631f687b3c585f728965ef7))

## [0.1.0](https://github.com/bucketeer-io/bucketeer/compare/v0.0.0...v0.1.0) (2022-09-28)


### Features

* add the initial implementation ([#1](https://github.com/bucketeer-io/bucketeer/issues/1)) ([038601c](https://github.com/bucketeer-io/bucketeer/commit/038601cc714b1fe66d7bf8b3763b344c89749a35))

## [0.1.0](https://github.com/bucketeer-io/bucketeer/compare/v0.0.0...v0.1.0) (2022-09-25)


### Features

* add initial implementation ([#1](https://github.com/bucketeer-io/bucketeer/issues/1)) ([2ddbb2c](https://github.com/bucketeer-io/bucketeer/commit/2ddbb2c455a99cbce30a6e6da0b3859fdcc4b919))

## [0.1.1](https://github.com/bucketeer-io/bucketeer/compare/v0.1.0...v0.1.1) (2022-09-25)


### Bug Fixes

* publish chart workflow not triggering ([5f73004](https://github.com/bucketeer-io/bucketeer/commit/5f7300484cb20ac5084960185a18b4ffe7160e1f))

## [0.1.0](https://github.com/bucketeer-io/bucketeer/compare/v0.0.0...v0.1.0) (2022-09-25)


### Features

* add initial implementation ([#1](https://github.com/bucketeer-io/bucketeer/issues/1)) ([47bdcec](https://github.com/bucketeer-io/bucketeer/commit/47bdcec22d4237fcc2b16b42198b9f1290e48ad0))
