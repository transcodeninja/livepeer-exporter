# Changelog

## [2.6.0](https://github.com/transcodeninja/livepeer-exporter/compare/v2.5.0...v2.6.0) (2023-12-09)


### Features

* add ticket and rewards gas cost ([#83](https://github.com/transcodeninja/livepeer-exporter/issues/83)) ([bdf28bf](https://github.com/transcodeninja/livepeer-exporter/commit/bdf28bfb9f958ec78bc23e9e2b03f1972ad96dcc))

## [2.5.0](https://github.com/transcodeninja/livepeer-exporter/compare/v2.4.0...v2.5.0) (2023-12-09)


### Features

* add gas price metric ([#81](https://github.com/transcodeninja/livepeer-exporter/issues/81)) ([27237aa](https://github.com/transcodeninja/livepeer-exporter/commit/27237aac867e5e36c7ba2b90105a0050ada92d61))

## [2.4.0](https://github.com/transcodeninja/livepeer-exporter/compare/v2.3.0...v2.4.0) (2023-12-08)


### Features

* add Livepeer X-Device-ID header ([#78](https://github.com/transcodeninja/livepeer-exporter/issues/78)) ([8e72745](https://github.com/transcodeninja/livepeer-exporter/commit/8e72745b337c5b64b842579f698bdc4f6342e509))

## [2.3.0](https://github.com/transcodeninja/livepeer-exporter/compare/v2.2.0...v2.3.0) (2023-12-06)


### Features

* add orch address validation ([#75](https://github.com/transcodeninja/livepeer-exporter/issues/75)) ([286e03d](https://github.com/transcodeninja/livepeer-exporter/commit/286e03da3a04b1e25b80c6655c9fa8b84e13d8c1))
* **info_exporter:** migrate to Livepeer subgraph ([#73](https://github.com/transcodeninja/livepeer-exporter/issues/73)) ([7e94f65](https://github.com/transcodeninja/livepeer-exporter/commit/7e94f65160717f1d90c5a7cd592f3eee57598e37))

## [2.2.0](https://github.com/transcodeninja/livepeer-exporter/compare/v2.1.5...v2.2.0) (2023-12-06)


### Features

* **delegators_exporter:** migrate to Livepeer subgraph ([#72](https://github.com/transcodeninja/livepeer-exporter/issues/72)) ([d25dda8](https://github.com/transcodeninja/livepeer-exporter/commit/d25dda8bc9c9d8568d9565aab36d987734aa7b2d))
* **rewards_exporter:** migrate to Livepeer subgraph ([#69](https://github.com/transcodeninja/livepeer-exporter/issues/69)) ([4eda8af](https://github.com/transcodeninja/livepeer-exporter/commit/4eda8af7f1179507a8193edce5caf37f078aeb54))
* **tickets_exporter:** migrate to Livepeer subgraph ([#67](https://github.com/transcodeninja/livepeer-exporter/issues/67)) ([1d4d50e](https://github.com/transcodeninja/livepeer-exporter/commit/1d4d50ef1403477efd70fde92de0c9a67d8ca9a2))


### Bug Fixes

* **tickets_exporter:** ensures that the graphql query is dynamic ([#70](https://github.com/transcodeninja/livepeer-exporter/issues/70)) ([e0e36ac](https://github.com/transcodeninja/livepeer-exporter/commit/e0e36ac41d7bbaab840028884e7c433ab0fbfe78))

## [2.1.5](https://github.com/transcodeninja/livepeer-exporter/compare/v2.1.4...v2.1.5) (2023-12-05)


### Bug Fixes

* apply orch info endpoint hotfix ([#65](https://github.com/transcodeninja/livepeer-exporter/issues/65)) ([c27c3da](https://github.com/transcodeninja/livepeer-exporter/commit/c27c3da536d478da50090d31a391694798546cd9))

## [2.1.3](https://github.com/rickstaa/livepeer-exporter/compare/v2.1.2...v2.1.3) (2023-11-18)


### Bug Fixes

* remove test stream update delay ([#52](https://github.com/rickstaa/livepeer-exporter/issues/52)) ([5d8b2d7](https://github.com/rickstaa/livepeer-exporter/commit/5d8b2d7b426b41244641afcc42faaeb2b8fc8bce))

## [2.1.1](https://github.com/rickstaa/livepeer-exporter/compare/v2.1.0...v2.1.1) (2023-11-17)


### Bug Fixes

* restore exposed port ([fb5b6a1](https://github.com/rickstaa/livepeer-exporter/commit/fb5b6a13f562f81020669c8fda1aec4455e0ceb1))

## [2.1.0](https://github.com/rickstaa/livepeer-exporter/compare/v2.0.0...v2.1.0) (2023-11-17)


### Features

* ctach http serve errors ([#46](https://github.com/rickstaa/livepeer-exporter/issues/46)) ([300c7b7](https://github.com/rickstaa/livepeer-exporter/commit/300c7b715481fc59c4063cf252fa6a4295a80016))
* increase default update intervals ([#47](https://github.com/rickstaa/livepeer-exporter/issues/47)) ([ff6c721](https://github.com/rickstaa/livepeer-exporter/commit/ff6c7217ff7f926eb5022501c0644ffa18a3f462))
* mitigate racing conditions ([#44](https://github.com/rickstaa/livepeer-exporter/issues/44)) ([9336d1c](https://github.com/rickstaa/livepeer-exporter/commit/9336d1cebc189c3d17b9b43a34eddb860c5347d0))

## [2.0.0](https://github.com/rickstaa/livepeer-exporter/compare/v1.1.1...v2.0.0) (2023-11-16)


### ⚠ BREAKING CHANGES

* improve environment variables ([#41](https://github.com/rickstaa/livepeer-exporter/issues/41))
* fix total claimed rewards and reward call ratio ([#40](https://github.com/rickstaa/livepeer-exporter/issues/40))
* change orchestrator score range ([#37](https://github.com/rickstaa/livepeer-exporter/issues/37))

### Features

* add crypto prices exporter ([#42](https://github.com/rickstaa/livepeer-exporter/issues/42)) ([08635d5](https://github.com/rickstaa/livepeer-exporter/commit/08635d570022341afbd893c42a64f860b06898c6))


### Bug Fixes

* fix total claimed rewards and reward call ratio ([#40](https://github.com/rickstaa/livepeer-exporter/issues/40)) ([2de5e5c](https://github.com/rickstaa/livepeer-exporter/commit/2de5e5c4ceb13382756688a0b8de8176efb2af05))


### Code Refactoring

* change orchestrator score range ([#37](https://github.com/rickstaa/livepeer-exporter/issues/37)) ([49b1594](https://github.com/rickstaa/livepeer-exporter/commit/49b159442bd0634dec418d031e60f9e3a1adb895))
* improve environment variables ([#41](https://github.com/rickstaa/livepeer-exporter/issues/41)) ([9c71e16](https://github.com/rickstaa/livepeer-exporter/commit/9c71e16af89058cb7952993c43c01bf41dc1be1b))

## [1.1.1](https://github.com/rickstaa/livepeer-exporter/compare/v1.1.0...v1.1.1) (2023-11-14)


### Bug Fixes

* fix empty string parse bug ([#34](https://github.com/rickstaa/livepeer-exporter/issues/34)) ([1f3e819](https://github.com/rickstaa/livepeer-exporter/commit/1f3e819acaa57d0598ce001c8a281400eebfbe11))

## [1.1.1](https://github.com/rickstaa/livepeer-exporter/compare/v1.1.0...v1.1.1) (2023-11-14)


### Bug Fixes

* fix empty string parse bug ([#34](https://github.com/rickstaa/livepeer-exporter/issues/34)) ([1f3e819](https://github.com/rickstaa/livepeer-exporter/commit/1f3e819acaa57d0598ce001c8a281400eebfbe11))

## [1.1.0](https://github.com/rickstaa/livepeer-exporter/compare/v1.0.1...v1.1.0) (2023-11-14)


### Features

* add 'LIVEPEER_EXPORTER_TICKETS_FETCH_INTERVAL' env variable ([f399f68](https://github.com/rickstaa/livepeer-exporter/commit/f399f684c819228f0e0816ba19cc9706d8e0c348))
* add 'LIVEPEER_EXPORTER_TICKETS_FETCH_INTERVAL' env variable ([#29](https://github.com/rickstaa/livepeer-exporter/issues/29)) ([33e5778](https://github.com/rickstaa/livepeer-exporter/commit/33e577845a45aef69e28eba75387d42a4fbdd998))
* add orchestrator rewards sub-exporter ([#30](https://github.com/rickstaa/livepeer-exporter/issues/30)) ([8156e81](https://github.com/rickstaa/livepeer-exporter/commit/8156e817f817876ae959f9f0be9e8139aeecbd9b))
* add orchestrator tickets sub-exporter ([980ad7d](https://github.com/rickstaa/livepeer-exporter/commit/980ad7dbc0502f357814e63465f401fa04328441))
* add orchestrator tickets sub-exporter ([#26](https://github.com/rickstaa/livepeer-exporter/issues/26)) ([646c5e9](https://github.com/rickstaa/livepeer-exporter/commit/646c5e95405f126afbbb10acaed51b4f3f433e4b))
* start sub-exporters in goroutines ([#28](https://github.com/rickstaa/livepeer-exporter/issues/28)) ([92880d8](https://github.com/rickstaa/livepeer-exporter/commit/92880d8d4945cbfb23973774c24ee517cb23aab2))


### Bug Fixes

* fix incorrect tickets exporter folder name ([fd32192](https://github.com/rickstaa/livepeer-exporter/commit/fd32192a2ce791ee490bb5c47f6e7dd6a0c8c70d))

## [1.1.0](https://github.com/rickstaa/livepeer-exporter/compare/v1.0.1...v1.1.0) (2023-11-14)


### Features

* add 'LIVEPEER_EXPORTER_TICKETS_FETCH_INTERVAL' env variable ([f399f68](https://github.com/rickstaa/livepeer-exporter/commit/f399f684c819228f0e0816ba19cc9706d8e0c348))
* add 'LIVEPEER_EXPORTER_TICKETS_FETCH_INTERVAL' env variable ([#29](https://github.com/rickstaa/livepeer-exporter/issues/29)) ([33e5778](https://github.com/rickstaa/livepeer-exporter/commit/33e577845a45aef69e28eba75387d42a4fbdd998))
* add orchestrator rewards sub-exporter ([#30](https://github.com/rickstaa/livepeer-exporter/issues/30)) ([8156e81](https://github.com/rickstaa/livepeer-exporter/commit/8156e817f817876ae959f9f0be9e8139aeecbd9b))
* add orchestrator tickets sub-exporter ([980ad7d](https://github.com/rickstaa/livepeer-exporter/commit/980ad7dbc0502f357814e63465f401fa04328441))
* add orchestrator tickets sub-exporter ([#26](https://github.com/rickstaa/livepeer-exporter/issues/26)) ([646c5e9](https://github.com/rickstaa/livepeer-exporter/commit/646c5e95405f126afbbb10acaed51b4f3f433e4b))
* start sub-exporters in goroutines ([#28](https://github.com/rickstaa/livepeer-exporter/issues/28)) ([92880d8](https://github.com/rickstaa/livepeer-exporter/commit/92880d8d4945cbfb23973774c24ee517cb23aab2))


### Bug Fixes

* fix incorrect tickets exporter folder name ([fd32192](https://github.com/rickstaa/livepeer-exporter/commit/fd32192a2ce791ee490bb5c47f6e7dd6a0c8c70d))

## [1.0.1](https://github.com/rickstaa/livepeer-exporter/compare/v1.0.0...v1.0.1) (2023-11-14)


### Bug Fixes

* fix reward and fee cut format ([#24](https://github.com/rickstaa/livepeer-exporter/issues/24)) ([c8b6a97](https://github.com/rickstaa/livepeer-exporter/commit/c8b6a9741fe8134c6f25ccc6a672125424b33896))

## 1.0.0 (2023-11-14)


### ⚠ BREAKING CHANGES

* add 'LIVEPEER_EXPORTER' prefix to environment variables ([#22](https://github.com/rickstaa/livepeer-exporter/issues/22))

### Features

* add exporter log statements ([db24ee1](https://github.com/rickstaa/livepeer-exporter/commit/db24ee1945bffb1698fb4440c81c5d32e431a33a))
* add livepeer-exporter ([92309b4](https://github.com/rickstaa/livepeer-exporter/commit/92309b4240d7114ec44c6c30ba36fbc0fc50b50a))
* add reward call ratio ([#17](https://github.com/rickstaa/livepeer-exporter/issues/17)) ([ff4d85d](https://github.com/rickstaa/livepeer-exporter/commit/ff4d85d3b0bc49bdeea23a79f07af646bf6ef648))


### Bug Fixes

* fix exporter port ([c3fbb74](https://github.com/rickstaa/livepeer-exporter/commit/c3fbb74ce2169abdbc17541c3aa1393c62f328bc))
* fix reward and fee cut format ([#14](https://github.com/rickstaa/livepeer-exporter/issues/14)) ([695683a](https://github.com/rickstaa/livepeer-exporter/commit/695683a7d0af3d1c72dc9cc21c441e066a56e667))
* fix reward call ratio calculation ([#19](https://github.com/rickstaa/livepeer-exporter/issues/19)) ([c8e1bed](https://github.com/rickstaa/livepeer-exporter/commit/c8e1bedce9e447ba5b036abdbd8b60b3a49ec289))


### Code Refactoring

* add 'LIVEPEER_EXPORTER' prefix to environment variables ([#22](https://github.com/rickstaa/livepeer-exporter/issues/22)) ([f92ad12](https://github.com/rickstaa/livepeer-exporter/commit/f92ad126153adee3ae0fae5ffab7f808988dc83b))
