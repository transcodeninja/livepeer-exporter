# Changelog

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
