# Changelog

## 0.85.0 (2026-07-29)

Full Changelog: [v0.84.0...v0.85.0](https://github.com/kernel/kernel-go-sdk/compare/v0.84.0...v0.85.0)

### Features

* Add encrypted per-proxy CA bundle for BYO MITM proxies ([b0a77ec](https://github.com/kernel/kernel-go-sdk/commit/b0a77ecf220c82d05b3e1a31d7dc7cba5a01e21f))
* Stabilize project lifecycle error codes ([778c389](https://github.com/kernel/kernel-go-sdk/commit/778c3890f53b2e1e5b8d31cc31a367fbc387ff49))

## 0.84.0 (2026-07-27)

Full Changelog: [v0.83.0...v0.84.0](https://github.com/kernel/kernel-go-sdk/compare/v0.83.0...v0.84.0)

### Features

* Expose telemetry state on managed auth timeline events ([9e6eace](https://github.com/kernel/kernel-go-sdk/commit/9e6eace40c6f3da89e0f446e5ed6de93323ebea5))


### Bug Fixes

* **ci:** unbreak update-cli-coverage workflow and switch to Claude Code ([#142](https://github.com/kernel/kernel-go-sdk/issues/142)) ([12b3ec6](https://github.com/kernel/kernel-go-sdk/commit/12b3ec62f63ee66d3cd8d83d9c2b9584696dd15f))

## 0.83.0 (2026-07-24)

Full Changelog: [v0.82.0...v0.83.0](https://github.com/kernel/kernel-go-sdk/compare/v0.82.0...v0.83.0)

### Features

* Add browser telemetry to managed auth connections ([fa63a73](https://github.com/kernel/kernel-go-sdk/commit/fa63a7377dc2c5b78886e54e465a2387c7abbf70))

## 0.82.0 (2026-07-23)

Full Changelog: [v0.81.0...v0.82.0](https://github.com/kernel/kernel-go-sdk/compare/v0.81.0...v0.82.0)

### Features

* Count project names by Unicode code point ([fa9dcd4](https://github.com/kernel/kernel-go-sdk/commit/fa9dcd430baef098caec1057a767ed718af907eb))
* Stream profile downloads as tar archives ([2a2ffe3](https://github.com/kernel/kernel-go-sdk/commit/2a2ffe3a6547e97a39e4b7a677a1f92c492455b9))
* Use format-neutral profile download Accept header ([d6464d5](https://github.com/kernel/kernel-go-sdk/commit/d6464d5539667645bf79f5e10389b04c7a073db2))

## 0.81.0 (2026-07-21)

Full Changelog: [v0.80.0...v0.81.0](https://github.com/kernel/kernel-go-sdk/compare/v0.80.0...v0.81.0)

### Features

* add complete audit log downloads ([2f992e4](https://github.com/kernel/kernel-go-sdk/commit/2f992e48aac664dcce5c7ce22d2879cd0c7dea2e))

## 0.80.0 (2026-07-19)

Full Changelog: [v0.79.0...v0.80.0](https://github.com/kernel/kernel-go-sdk/compare/v0.79.0...v0.80.0)

### Features

* Add example to ProxyCheckRequest.url so API reference sample shows it ([b32ed11](https://github.com/kernel/kernel-go-sdk/commit/b32ed11503054394e20e6d0b81302a4e1e1ded08))
* Expose browser sessions on auth timeline events ([00728f0](https://github.com/kernel/kernel-go-sdk/commit/00728f053c7a478222312c4baf73c39c82fc0156))
* **managed-auth:** add TS CUA worker contract (KERNEL-1456) ([5951dfa](https://github.com/kernel/kernel-go-sdk/commit/5951dfae60dce7e60c420c758b5217fe2f96bc6a))
* **stlc:** configurable CI runner and private-production-repo support in workflow templates ([77f3c17](https://github.com/kernel/kernel-go-sdk/commit/77f3c174783e6c4e8dc227d30c0505aee46f4cfd))

## 0.79.0 (2026-07-15)

Full Changelog: [v0.78.1...v0.79.0](https://github.com/kernel/kernel-go-sdk/compare/v0.78.1...v0.79.0)

### Features

* Add telemetry support for browser pools with BAA enforcement ([042f389](https://github.com/kernel/kernel-go-sdk/commit/042f3899a86034a68bdb3493f91cde4f254c75fc))

## 0.78.1 (2026-07-15)

Full Changelog: [v0.78.0...v0.78.1](https://github.com/kernel/kernel-go-sdk/compare/v0.78.0...v0.78.1)

### Features

* Default dashboard browser pools to 25% fill rate ([d00ac8a](https://github.com/kernel/kernel-go-sdk/commit/d00ac8a78bec4aa2c271ab5917ac065e6db7cbc5))

## 0.78.0 (2026-07-13)

Full Changelog: [v0.77.0...v0.78.0](https://github.com/kernel/kernel-go-sdk/compare/v0.77.0...v0.78.0)

### Features

* Expose telemetry exception message in API/SDK ([7426310](https://github.com/kernel/kernel-go-sdk/commit/74263107550f6d13d74a06f0df0c8d603920d62f))

## 0.77.0 (2026-07-13)

Full Changelog: [v0.76.0...v0.77.0](https://github.com/kernel/kernel-go-sdk/compare/v0.76.0...v0.77.0)

### Features

* Add exact-match name filter to list endpoints ([6ffae2f](https://github.com/kernel/kernel-go-sdk/commit/6ffae2ffa5a381af12ff4776ac3d590ed1a724ef))
* Add name-only rename for profiles and proxies ([c2b2b5e](https://github.com/kernel/kernel-go-sdk/commit/c2b2b5e5cf318b2d6c9e964b65274f8de4a57a35))
* Auto-default refresh_on_profile_update when browser pool profile changes ([b74cae1](https://github.com/kernel/kernel-go-sdk/commit/b74cae1d226293704baefb943c213c77b83b409b))
* Document name uniqueness and query match semantics ([fd4fcff](https://github.com/kernel/kernel-go-sdk/commit/fd4fcffa2ad8d35878a189a1a66e9e8aac2df872))
* Make the browser pool OpenAPI contract truthful ([b1baff2](https://github.com/kernel/kernel-go-sdk/commit/b1baff2db6ad583e8998e4ef6d3d711c73817492))
* Persist and echo deployment source identity ([a09849a](https://github.com/kernel/kernel-go-sdk/commit/a09849ac32a842675babeefdf32d136e18a9d812))
* Support multiple audit log method exclusions ([de83d75](https://github.com/kernel/kernel-go-sdk/commit/de83d75e14f7f819f6c6982a7c9833d7a3f45721))


### Documentation

* **openapi:** describe unified concurrency limit, deprecate max_pooled_sessions (CUS-275) ([edcc0e5](https://github.com/kernel/kernel-go-sdk/commit/edcc0e5e25cbce705d36a316b146a6eda30dc4fb))

## 0.76.0 (2026-07-09)

Full Changelog: [v0.75.0...v0.76.0](https://github.com/kernel/kernel-go-sdk/compare/v0.75.0...v0.76.0)

### Features

* Layer telemetry request config onto the default set ([7ff6801](https://github.com/kernel/kernel-go-sdk/commit/7ff68018d09d2b63d649f373ccc262e8e1e6f4b7))
* Return credential value keys and support removing values on update ([750c15f](https://github.com/kernel/kernel-go-sdk/commit/750c15f460a90377acba471935e5f310c8c70a10))

## 0.75.0 (2026-07-08)

Full Changelog: [v0.74.0...v0.75.0](https://github.com/kernel/kernel-go-sdk/compare/v0.74.0...v0.75.0)

### Features

* Auto-flush pools when a managed auth profile re-authenticates ([9b65633](https://github.com/kernel/kernel-go-sdk/commit/9b656333bbab210f4deaf9e8247b2540571693ca))
* Document env var redaction on deployment and app reads ([1a17cf5](https://github.com/kernel/kernel-go-sdk/commit/1a17cf51617d671b676ec713409a05ee002b086b))
* Expose resolved profile_id and extension_ids on browser pool reads ([3931759](https://github.com/kernel/kernel-go-sdk/commit/393175979cfbd763a610d7d87601a8386feba2ad))
* Reject API key self-deletion ([0fe2a82](https://github.com/kernel/kernel-go-sdk/commit/0fe2a8206644361bcdb439f0903b8658a85c5691))
* Revert "Store and return a sha256 checksum for uploaded extensions (#… ([763d2cd](https://github.com/kernel/kernel-go-sdk/commit/763d2cd76858214f0f4e34cabc2c744a0d44a983))
* Store and return a sha256 checksum for uploaded extensions ([92806ce](https://github.com/kernel/kernel-go-sdk/commit/92806ce19504a6c6bbb8a4301bb2a667843a3b24))
* Store and return a sha256 checksum for uploaded extensions (reland) ([24fea5b](https://github.com/kernel/kernel-go-sdk/commit/24fea5b4a661c7970f7f773261a2d7834042c942))


### Documentation

* **api:** clarify reuse/discard_all_idle pool config staleness ([c2a231d](https://github.com/kernel/kernel-go-sdk/commit/c2a231d2184874248ce5554c61505376e644d6d0))

## 0.74.0 (2026-07-06)

Full Changelog: [v0.73.0...v0.74.0](https://github.com/kernel/kernel-go-sdk/compare/v0.73.0...v0.74.0)

### Features

* Add order=desc pagination to telemetry event reads ([f86b440](https://github.com/kernel/kernel-go-sdk/commit/f86b44009ac649334fe1b6c85b0a2c54326e4650))
* Add tablet and mobile viewport presets to pool dashboard ([eec3e1e](https://github.com/kernel/kernel-go-sdk/commit/eec3e1ee4b4a4386d25c260840d279d3afcd9895))

## 0.73.0 (2026-07-01)

Full Changelog: [v0.72.0...v0.73.0](https://github.com/kernel/kernel-go-sdk/compare/v0.72.0...v0.73.0)

### Features

* Add hidden audit-logs export endpoint ([2afb724](https://github.com/kernel/kernel-go-sdk/commit/2afb72456a09c06411d78a2fd53a27ac8f1ebdc5))

## 0.72.0 (2026-06-26)

Full Changelog: [v0.71.0...v0.72.0](https://github.com/kernel/kernel-go-sdk/compare/v0.71.0...v0.72.0)

### Bug Fixes

* **api:** browser pool profile omits save_changes (BrowserPoolProfile) ([5e31861](https://github.com/kernel/kernel-go-sdk/commit/5e31861edf251908db94a60e32f5f15507a6f272))

## 0.71.0 (2026-06-26)

Full Changelog: [v0.70.0...v0.71.0](https://github.com/kernel/kernel-go-sdk/compare/v0.70.0...v0.71.0)

### Features

* Add auth connection event timeline endpoint ([d14cbe4](https://github.com/kernel/kernel-go-sdk/commit/d14cbe46c35e9c5e58398551bb6a903f0bcb02ff))
* Expose audit logs in public SDK ([8c92e59](https://github.com/kernel/kernel-go-sdk/commit/8c92e59f00802a0fe4393ba7bdea1070554ab8b9))

## 0.70.0 (2026-06-24)

Full Changelog: [v0.69.0...v0.70.0](https://github.com/kernel/kernel-go-sdk/compare/v0.69.0...v0.70.0)

### Features

* Add GET /browsers/{id}/telemetry/events (read from S2) ([2c3e61e](https://github.com/kernel/kernel-go-sdk/commit/2c3e61ed66a4b8810a0ef33a8abb09555c030a8d))
* Align browser-pool timeout/viewport/fill-rate contract with implementation; reject save_changes on update ([a986c45](https://github.com/kernel/kernel-go-sdk/commit/a986c453ecee67ebf6b5cf39ecb8e9accbd2c1dc))
* api: support per-acquire start_url override on browser pool acquire ([1c675a3](https://github.com/kernel/kernel-go-sdk/commit/1c675a387cb712036cea70ac75dac2026276704b))
* **api:** add GET /extensions/{id_or_name}/metadata ([373f33b](https://github.com/kernel/kernel-go-sdk/commit/373f33b293159fad981ff78770762e07ac80b5ca))
* **api:** resolve GET /org/projects/{id} by ID or name ([395e0a7](https://github.com/kernel/kernel-go-sdk/commit/395e0a7d43e289a1b4dc187ed2af5a347f0dfc2b))
* Forward replay param through telemetry stream passthrough ([ae02eeb](https://github.com/kernel/kernel-go-sdk/commit/ae02eeb4d15611a0221ed06a2600f55460e78184))


### Bug Fixes

* don't misroute telemetry/events to the browser VM ([95487e0](https://github.com/kernel/kernel-go-sdk/commit/95487e0182fb969bbaa4d44425a688084020d943))

## 0.69.0 (2026-06-18)

Full Changelog: [v0.68.0...v0.69.0](https://github.com/kernel/kernel-go-sdk/compare/v0.68.0...v0.69.0)

### Features

* Add free-text search to remaining paginated list endpoints ([3d07531](https://github.com/kernel/kernel-go-sdk/commit/3d0753183c1bfc0f251273cc0a34116e4642f848))

## 0.68.0 (2026-06-15)

Full Changelog: [v0.67.0...v0.68.0](https://github.com/kernel/kernel-go-sdk/compare/v0.67.0...v0.68.0)

### Features

* Add API key rotate endpoint ([04d4a12](https://github.com/kernel/kernel-go-sdk/commit/04d4a12f1f98a5745b7ca706a22cab71cb625866))
* **api:** surface deleted/expired API keys for audit trail (KERNEL-1350) ([e7f9fd8](https://github.com/kernel/kernel-go-sdk/commit/e7f9fd8aa890b660d675b22ad48f869504817b2d))


### Refactors

* **api:** align API key audit surface with browser sibling (KERNEL-1350) ([9e23506](https://github.com/kernel/kernel-go-sdk/commit/9e2350674aafa75a45c4b843d8a43b31c4656bd6))

## 0.67.0 (2026-06-11)

Full Changelog: [v0.66.0...v0.67.0](https://github.com/kernel/kernel-go-sdk/compare/v0.66.0...v0.67.0)

### Features

* Add project_id SDK client option mapped to X-Kernel-Project-Id ([fc37447](https://github.com/kernel/kernel-go-sdk/commit/fc37447ec9fc09e21c04d1ecee2009dc44010363))


### Documentation

* **api:** correct project-scoping descriptions in OpenAPI spec ([d10587c](https://github.com/kernel/kernel-go-sdk/commit/d10587c1131648662d969b25ed921264063e36a4))

## 0.66.0 (2026-06-10)

Full Changelog: [v0.65.0...v0.66.0](https://github.com/kernel/kernel-go-sdk/compare/v0.65.0...v0.66.0)

### Features

* Add org-level default per-project concurrency cap ([da8ea41](https://github.com/kernel/kernel-go-sdk/commit/da8ea416e8e2d00f38c7a811cbdd515053b37f89))
* Support updating browser session name and tags via PATCH ([f530bfb](https://github.com/kernel/kernel-go-sdk/commit/f530bfb8b6822f3774eee1fb7bbe81fc4ce9fa41))

## 0.65.0 (2026-06-08)

Full Changelog: [v0.64.0...v0.65.0](https://github.com/kernel/kernel-go-sdk/compare/v0.64.0...v0.65.0)

### Features

* **api:** allow setting a name and tags on a pool-acquired browser session ([ef64773](https://github.com/kernel/kernel-go-sdk/commit/ef647732782e8bc382a06d02377f46edd8ebc2b2))
* **api:** support id-or-name lookup on browser session get/patch/delete ([193727b](https://github.com/kernel/kernel-go-sdk/commit/193727b6134db7c2f97d5170ef06c2b0692b45c8))

## 0.64.0 (2026-06-05)

Full Changelog: [v0.63.0...v0.64.0](https://github.com/kernel/kernel-go-sdk/compare/v0.63.0...v0.64.0)

### Features

* Telemetry: expose opt-in categories + full event taxonomy (public API) ([c965b72](https://github.com/kernel/kernel-go-sdk/commit/c965b72cfe8a381469c43ad1b3089442bbfbc8a3))

## 0.63.0 (2026-06-05)

Full Changelog: [v0.62.0...v0.63.0](https://github.com/kernel/kernel-go-sdk/compare/v0.62.0...v0.63.0)

### Features

* api: label proxy config oneOf options by type ([abbee19](https://github.com/kernel/kernel-go-sdk/commit/abbee1950d51eda4edfed2c0a13162d55881dcf6))
* **api:** allow setting a custom name on a browser session at create time ([22794dc](https://github.com/kernel/kernel-go-sdk/commit/22794dc349e48a6e9552ad80489b268b9e8aeb02))
* **api:** allow setting key-value tags on a browser session at create time ([12d335f](https://github.com/kernel/kernel-go-sdk/commit/12d335fb3e7079f614e7b8e9c174cf20f4579820))


### Documentation

* **api:** use neutral example for browser session name field ([2319e34](https://github.com/kernel/kernel-go-sdk/commit/2319e34f9b63ba11104b3ee7a4add0db9f7494ae))

## 0.62.0 (2026-06-04)

Full Changelog: [v0.61.0...v0.62.0](https://github.com/kernel/kernel-go-sdk/compare/v0.61.0...v0.62.0)

### Features

* api: paginate GET /browser_pools ([ed9dc35](https://github.com/kernel/kernel-go-sdk/commit/ed9dc35f3071d978df1e53e77f333f4582f5752a))
* api: paginate GET /extensions ([f48c9cb](https://github.com/kernel/kernel-go-sdk/commit/f48c9cb83ffd3479eb17c32cef049d703106763c))
* api: paginate GET /org/credential_providers ([d72262b](https://github.com/kernel/kernel-go-sdk/commit/d72262bc8ef36d40e4a4f524a0276b5aad167160))
* api: paginate GET /proxies ([cd37790](https://github.com/kernel/kernel-go-sdk/commit/cd37790a3ed1cb029adeed2493b1a0a22ca65944))

## 0.61.0 (2026-06-03)

Full Changelog: [v0.60.0...v0.61.0](https://github.com/kernel/kernel-go-sdk/compare/v0.60.0...v0.61.0)

### Features

* Add record_audio option to browser replay recording API ([c04e9d8](https://github.com/kernel/kernel-go-sdk/commit/c04e9d8896bfad648ea70fd7b9ecbf37293b1c82))

## 0.60.0 (2026-06-03)

Full Changelog: [v0.59.0...v0.60.0](https://github.com/kernel/kernel-go-sdk/compare/v0.59.0...v0.60.0)

### Features

* Add API-backed API key management endpoints ([6ea145b](https://github.com/kernel/kernel-go-sdk/commit/6ea145bdd708bdcf2f264cf1c379d9d9ccaf3779))
* Fix browser pool update schema ([370464e](https://github.com/kernel/kernel-go-sdk/commit/370464e57cee08cebdfa9a920681d68e0a6f4e14))
* route browser telemetry directly to the VM by default ([105313c](https://github.com/kernel/kernel-go-sdk/commit/105313c08767599494267e665346b0ca0b561d66))


### Bug Fixes

* **ssestream:** skip empty SSE keepalive comment frames ([59bd952](https://github.com/kernel/kernel-go-sdk/commit/59bd952042fdf647481e4c71f86f71077b31fe86))


### Refactors

* **examples:** rename telemetry example to browser-telemetry ([3908ff5](https://github.com/kernel/kernel-go-sdk/commit/3908ff56aab71f925ce8d2a7e62259db6ce9a0ff))

## 0.59.0 (2026-06-03)

Full Changelog: [v0.58.0...v0.59.0](https://github.com/kernel/kernel-go-sdk/compare/v0.58.0...v0.59.0)

### Features

* api: surface category field on browser telemetry events ([803734d](https://github.com/kernel/kernel-go-sdk/commit/803734d1179cb14da37173537fa25cb10995fd78))
* **api:** move browser telemetry SSE stream to /browsers/{id}/telemetry/stream ([2c542d9](https://github.com/kernel/kernel-go-sdk/commit/2c542d922d4200170d630198cfac01173279c17e))
* Support Byteful mobile proxies ([d7162b5](https://github.com/kernel/kernel-go-sdk/commit/d7162b5cbb1a39b46cb25ff9a0124a8584d3efaf))


### Bug Fixes

* **api:** move batch + get_mouse_position into Browser Computer Controls tag ([743624f](https://github.com/kernel/kernel-go-sdk/commit/743624f14a8d5fd86a7690d8f2525e38f671b355))

## 0.58.0 (2026-05-27)

Full Changelog: [v0.57.0...v0.58.0](https://github.com/kernel/kernel-go-sdk/compare/v0.57.0...v0.58.0)

### Features

* [codex] Expose API keys in SDK config ([bfed2b8](https://github.com/kernel/kernel-go-sdk/commit/bfed2b81450681bb0a053d233f2cc1ecf65762c8))
* Fix API key request model SDK metadata ([ebc9f45](https://github.com/kernel/kernel-go-sdk/commit/ebc9f45e57aef38964be31942fcf0e8ad554b37b))
* Support telemetry enabled request config and fix SDK metadata ([53bff5e](https://github.com/kernel/kernel-go-sdk/commit/53bff5ecd62bda11e2fac6db87e0ffaec3edad44))


### Chores

* refresh Go SDK release PR ([fe23a23](https://github.com/kernel/kernel-go-sdk/commit/fe23a23c3918c8ee086ecca4e9f9915d57b439b9))

## 0.57.0 (2026-05-26)

Full Changelog: [v0.56.0...v0.57.0](https://github.com/kernel/kernel-go-sdk/compare/v0.56.0...v0.57.0)

### Features

* [kernel-1116] browser events api integration ([e0a367b](https://github.com/kernel/kernel-go-sdk/commit/e0a367b4fd596424dc913c9e999831ae1fdc068e))
* api: dual-route /projects under /org/projects, deprecate /projects ([166c834](https://github.com/kernel/kernel-go-sdk/commit/166c83456a6c1eaf5a28b821383576eacbec3a3a))
* **api:** type can_reauth_reason as an enum on ManagedAuth ([61f5842](https://github.com/kernel/kernel-go-sdk/commit/61f5842a55a7e846c6cf04f6d2c2dd7df5b40862))
* browsers: accept chrome_policy on POST /browsers (KERNEL-1216) ([6f29adc](https://github.com/kernel/kernel-go-sdk/commit/6f29adc513752623a3ea761195005f5a50d46e5b))
* EOL persistent browsers: openapi limits ([b79f2fd](https://github.com/kernel/kernel-go-sdk/commit/b79f2fd415108f84cf99377bd3b4c60e64af1ba3))
* Support telemetry enabled request config ([ba2ff70](https://github.com/kernel/kernel-go-sdk/commit/ba2ff705279a32f8a3d91ca85696e398923ebe43))

## 0.56.0 (2026-05-18)

Full Changelog: [v0.55.0...v0.56.0](https://github.com/kernel/kernel-go-sdk/compare/v0.55.0...v0.56.0)

### Features

* Expose POST /projects in public API ([6f165cd](https://github.com/kernel/kernel-go-sdk/commit/6f165cd1c8bec73f9948c1ee881ca981a32c0ba4))

## 0.55.0 (2026-05-15)

Full Changelog: [v0.53.0...v0.55.0](https://github.com/kernel/kernel-go-sdk/compare/v0.53.0...v0.55.0)

### Features

* Add health check and auto-reauth controls for managed auth connections ([e381d8e](https://github.com/kernel/kernel-go-sdk/commit/e381d8ed28c0d4fe978c0d27d724c0bc9b43d3d5))
* **client:** optimize json encoder for internal types ([e7f4365](https://github.com/kernel/kernel-go-sdk/commit/e7f436591debbdedd4150a0a4b2fcbdce2320cd1))
* Polish start URL OpenAPI descriptions ([2d82565](https://github.com/kernel/kernel-go-sdk/commit/2d82565e34904922fa4a6b87e076bb1790af43bd))

## 0.53.0 (2026-05-12)

Full Changelog: [v0.52.0...v0.53.0](https://github.com/kernel/kernel-go-sdk/compare/v0.52.0...v0.53.0)

### Features

* Add 'switch' MFA option type for generic method-switcher links ([b6c0d17](https://github.com/kernel/kernel-go-sdk/commit/b6c0d174ffaf8f5a0788f8c16a2fded6b56add76))
* Add opt-in record_session flag to managed auth ([5673b41](https://github.com/kernel/kernel-go-sdk/commit/5673b41f0cc5d6233dbaf9875bc20b84b7cc7762))
* **api:** server-side search on GET /projects ([81a050d](https://github.com/kernel/kernel-go-sdk/commit/81a050d64651caf4fefc58ce25a13c65ffe05bea))
* browser_pools: add start_url config (KERNEL-1217 PR 2) ([b2c8f95](https://github.com/kernel/kernel-go-sdk/commit/b2c8f951ddde27daf1285b78a1ce6a7fc01d74ef))
* managed-auth: surface awaiting_external_action even when fallback actions exist ([c1a9ba4](https://github.com/kernel/kernel-go-sdk/commit/c1a9ba4038b690875a5c8c5198d20671e27acec4))
* Scope name uniqueness to project for profiles, session_pools, extensions, credentials ([be3c6bd](https://github.com/kernel/kernel-go-sdk/commit/be3c6bd46f8826cf5d17523bb1aa79c77d9a8d39))


### Bug Fixes

* **go:** avoid panic when http.DefaultTransport is wrapped ([15e7a88](https://github.com/kernel/kernel-go-sdk/commit/15e7a8860af570cb968ed0153742c86af85c9909))


### Chores

* avoid embedding reflect.Type for dead code elimination ([c267709](https://github.com/kernel/kernel-go-sdk/commit/c2677092b3f611e3470f5b6730f32fe896e2d688))
* redact api-key headers in debug logs ([405d242](https://github.com/kernel/kernel-go-sdk/commit/405d242dd4564cb133cc036b7693169df2537e3b))


### Documentation

* clarify record_session description in OpenAPI spec ([c61abfa](https://github.com/kernel/kernel-go-sdk/commit/c61abfa3c5aa1d355bce51c10d08ad962e08424e))

## 0.52.0 (2026-04-29)

Full Changelog: [v0.51.0...v0.52.0](https://github.com/kernel/kernel-go-sdk/compare/v0.51.0...v0.52.0)

### Features

* **go:** add default http client with timeout ([a999228](https://github.com/kernel/kernel-go-sdk/commit/a999228cc4426709eb1561f6cf46d51833509014))
* profile download: 409 for empty profile + surface API errors in dashboard ([c14aef2](https://github.com/kernel/kernel-go-sdk/commit/c14aef289c5418906fb003893a728ec11cf549cc))
* support setting headers via env ([728fe39](https://github.com/kernel/kernel-go-sdk/commit/728fe3931a9a8b9f7de3ef6f338c3872fb67c4b1))

## 0.51.0 (2026-04-25)

Full Changelog: [v0.50.0...v0.51.0](https://github.com/kernel/kernel-go-sdk/compare/v0.50.0...v0.51.0)

### Features

* add browser-scoped session client ([ef994b4](https://github.com/kernel/kernel-go-sdk/commit/ef994b4d6189ee0cdb3ad6e002816e0d77c7db48))
* Expire stuck IN_PROGRESS managed auth sessions via background worker ([f441a6f](https://github.com/kernel/kernel-go-sdk/commit/f441a6fa6f51f5ceb9fab1248a568025079752cf))
* Expose browser_session_id on managed auth connection ([e23defc](https://github.com/kernel/kernel-go-sdk/commit/e23defc4abad6e49ae854ead1d1a015a2ecad097))
* generate browser-scoped service bindings ([b6a77bc](https://github.com/kernel/kernel-go-sdk/commit/b6a77bc656e6480e1bee0932fcc3b5beb40bcb10))


### Bug Fixes

* align browser-scoped routing with base_url ([64f7f81](https://github.com/kernel/kernel-go-sdk/commit/64f7f811a2f3d7c657ad56a1f34a8bea0206e585))
* clean up go browser routing follow-ups ([77bda33](https://github.com/kernel/kernel-go-sdk/commit/77bda33323c48a60ab652bf42a39cc2bc1f909a0))
* keep browser pool routing in sync with cache ([4f754d1](https://github.com/kernel/kernel-go-sdk/commit/4f754d13f38e5b545de8d84cb26de28f6f2f892a))
* make browser route deletion win over sniffing ([0e44ff3](https://github.com/kernel/kernel-go-sdk/commit/0e44ff363c53b1225d48313be6376b2cda7b057f))
* propagate browser HTTP client config errors ([1ec1358](https://github.com/kernel/kernel-go-sdk/commit/1ec1358150d20982fc582c634e68f42322d1fa4c))
* remove old go browser scope package ([d594f39](https://github.com/kernel/kernel-go-sdk/commit/d594f39a352e2ad5639807f222bb16535cce72c6))
* tighten browser route metadata parsing ([b293866](https://github.com/kernel/kernel-go-sdk/commit/b293866537bddb183ad7e8d014c0e2c298d715a9))
* tighten browser-scoped helper surface ([3e3e33f](https://github.com/kernel/kernel-go-sdk/commit/3e3e33f031d8460662da8a368467a3dcc6a0ec3f))


### Chores

* **internal:** more robust bootstrap script ([c0430e3](https://github.com/kernel/kernel-go-sdk/commit/c0430e3e2952f10648ed0181edf85a69230e9279))


### Documentation

* add browser-scoped raw http example ([92dc96e](https://github.com/kernel/kernel-go-sdk/commit/92dc96e99e14eeb9eb6881115f330369bb3c7542))


### Refactors

* move browser route cache sync into middleware ([62078d3](https://github.com/kernel/kernel-go-sdk/commit/62078d3213dbe4ede9600968b48d5c8e7aa118e4))
* move go browser routing rollout to env ([681e57f](https://github.com/kernel/kernel-go-sdk/commit/681e57f34b3199a4d720fb517da08b3bbd6bc0fb))
* narrow browser-scoped helper exports ([0ac61ef](https://github.com/kernel/kernel-go-sdk/commit/0ac61ef0532e6682eed740765975775767e28a10))
* remove browser session wrapper layer ([3452e53](https://github.com/kernel/kernel-go-sdk/commit/3452e537737576fdb645a4e5de0d9147e50cb33b))
* rename browser routing subresources config ([909c377](https://github.com/kernel/kernel-go-sdk/commit/909c377e21adddbaede526ba2a356352652e947b))
* simplify direct-to-vm route caching ([6bdf25f](https://github.com/kernel/kernel-go-sdk/commit/6bdf25f058d2e0bb552bd47ba4478222d97aa254))

## 0.50.0 (2026-04-20)

Full Changelog: [v0.49.0...v0.50.0](https://github.com/kernel/kernel-go-sdk/compare/v0.49.0...v0.50.0)

### Features

* add POST /browsers/{id}/curl and /curl/raw endpoints ([94b52ae](https://github.com/kernel/kernel-go-sdk/commit/94b52aeac38b20c8ea60e3c016443e0ccc680451))
* remove paid plan gating from project endpoints ([f3cbeb3](https://github.com/kernel/kernel-go-sdk/commit/f3cbeb32ca629fefbadaf5449acd0d7aa438e4c5))


### Bug Fixes

* include MFA and sign-in options in CUA SSO-only step response ([45db2eb](https://github.com/kernel/kernel-go-sdk/commit/45db2ebc96e4d6fbeec34cc43ed13ad4d0bd7de2))

## 0.49.0 (2026-04-10)

Full Changelog: [v0.48.0...v0.49.0](https://github.com/kernel/kernel-go-sdk/compare/v0.48.0...v0.49.0)

### Features

* Neil/kernel 1180 fuzzy matching for browser pools ([9b784cb](https://github.com/kernel/kernel-go-sdk/commit/9b784cb70c17efb0540fb28420c1f73d1c57942d))
* Raise replay framerate limit from 20 to 60 fps ([f095149](https://github.com/kernel/kernel-go-sdk/commit/f095149c6a25040b1c6b5f2b12d5a43fab736eb8))

## 0.48.0 (2026-04-10)

Full Changelog: [v0.47.0...v0.48.0](https://github.com/kernel/kernel-go-sdk/compare/v0.47.0...v0.48.0)

### Features

* [kernel-1116] add base_url field to browser session response ([6e52972](https://github.com/kernel/kernel-go-sdk/commit/6e529727c5fc0a122df257f58a2fa3dd285f66a7))


### Chores

* retrigger Stainless codegen for projects resource ([3e1b78f](https://github.com/kernel/kernel-go-sdk/commit/3e1b78f44879ed57afb128c103276a9bcf84a209))

## 0.47.0 (2026-04-07)

Full Changelog: [v0.46.0...v0.47.0](https://github.com/kernel/kernel-go-sdk/compare/v0.46.0...v0.47.0)

### Features

* Include login_url in managed auth connection response ([795989c](https://github.com/kernel/kernel-go-sdk/commit/795989c703df9fc5bd740eeb69538a4a4cec7cee))

## 0.46.0 (2026-04-06)

Full Changelog: [v0.45.0...v0.46.0](https://github.com/kernel/kernel-go-sdk/compare/v0.45.0...v0.46.0)

### Features

* Add optional url parameter to proxy check endpoint ([084975a](https://github.com/kernel/kernel-go-sdk/commit/084975a54df1144592e3f61655ee423df3d38393))


### Bug Fixes

* fix issue with unmarshaling in some cases ([6f721ac](https://github.com/kernel/kernel-go-sdk/commit/6f721ac3da13678eb476f8231617510e8f0b9e34))

## 0.45.0 (2026-03-30)

Full Changelog: [v0.44.0...v0.45.0](https://github.com/kernel/kernel-go-sdk/compare/v0.44.0...v0.45.0)

### Features

* [kernel-1008] browser pools add custom policy ([e740db7](https://github.com/kernel/kernel-go-sdk/commit/e740db79ec07414b5b0b5ff3be98b9e089a9bb45))
* Add disable_default_proxy for stealth browsers ([f584af5](https://github.com/kernel/kernel-go-sdk/commit/f584af5155e3fe8bb6dd87f48e1c5ee6cfe47111))
* **internal:** support comma format in multipart form encoding ([74f5c6a](https://github.com/kernel/kernel-go-sdk/commit/74f5c6a22be0f953416968dfc92b0354f66d5cd9))


### Bug Fixes

* prevent duplicate ? in query params ([66d97d3](https://github.com/kernel/kernel-go-sdk/commit/66d97d350f30d933b1cd4596da6e3fb4cac8c4a1))


### Chores

* **ci:** skip lint on metadata-only changes ([8296238](https://github.com/kernel/kernel-go-sdk/commit/8296238ce0b45903e3d40420100aa29a3ab105fd))
* **ci:** support opting out of skipping builds on metadata-only commits ([708da0f](https://github.com/kernel/kernel-go-sdk/commit/708da0fffff51a8f20905f9983e5beb1181a6a48))
* **client:** fix multipart serialisation of Default() fields ([cabda49](https://github.com/kernel/kernel-go-sdk/commit/cabda49a25334cd3fa3fc747f69b8a515e2b77d9))
* **internal:** support default value struct tag ([dd77af4](https://github.com/kernel/kernel-go-sdk/commit/dd77af45f77a98e96f5d140bb737377f1d72456a))
* **internal:** update gitignore ([48baab7](https://github.com/kernel/kernel-go-sdk/commit/48baab7aa0e23ed519a08b52142927340a8d451b))
* remove unnecessary error check for url parsing ([05d2048](https://github.com/kernel/kernel-go-sdk/commit/05d2048b0cc9027c63a4036b9f722044423f4e25))
* update docs for api:"required" ([abfd988](https://github.com/kernel/kernel-go-sdk/commit/abfd9883774e74c2d4c2cc482377364432318e27))

## 0.44.0 (2026-03-20)

Full Changelog: [v0.43.0...v0.44.0](https://github.com/kernel/kernel-go-sdk/compare/v0.43.0...v0.44.0)

### Features

* Add GPU viewport presets and GPU encoder defaults ([5096099](https://github.com/kernel/kernel-go-sdk/commit/5096099db1afd3d90b384423c8daf18576e167c7))
* Adds description to OAS spec for docs about delta_x, delta_y ([5d60a03](https://github.com/kernel/kernel-go-sdk/commit/5d60a03c51ab0a14ae115552909c09a43f730631))
* Drop headless GPU support and disable pooling ([188a32b](https://github.com/kernel/kernel-go-sdk/commit/188a32bf2e2bfb3e0aa33a595b2477557db168ea))
* Enhance managed authentication with CUA support and new features ([b1c79e6](https://github.com/kernel/kernel-go-sdk/commit/b1c79e61fba350ba9324c8ebaa2205f7ca96332a))
* expose smooth drag mouse movement via public API ([1bcd6f5](https://github.com/kernel/kernel-go-sdk/commit/1bcd6f5635cbd6bf4073efdbf9b3e97471cff826))
* Rename hardware acceleration UI/docs wording to GPU acceleration ([d225ff6](https://github.com/kernel/kernel-go-sdk/commit/d225ff61b6d63308babf283816ac60ad2902f91e))


### Chores

* **internal:** minor cleanup ([2e50750](https://github.com/kernel/kernel-go-sdk/commit/2e507500fb10787ef597819c04bef36baecef0bf))
* **internal:** tweak CI branches ([3a66a84](https://github.com/kernel/kernel-go-sdk/commit/3a66a8427998fe16353cbff5b14f4822e2f278ff))
* **internal:** use explicit returns ([4023de1](https://github.com/kernel/kernel-go-sdk/commit/4023de1ca30ac2bd31a4003941b35ad0e75b8a86))
* **internal:** use explicit returns in more places ([503a87c](https://github.com/kernel/kernel-go-sdk/commit/503a87cf7e6c0ffc76c1c063ff90f419ad0c78db))

## 0.43.0 (2026-03-10)

Full Changelog: [v0.42.1...v0.43.0](https://github.com/kernel/kernel-go-sdk/compare/v0.42.1...v0.43.0)

### Features

* Add webdriver_ws_url and metro webdriver session proxy ([4e7c5d5](https://github.com/kernel/kernel-go-sdk/commit/4e7c5d567427e28035b489fcbf6afcef0f1775c7))


### Chores

* **ci:** skip uploading artifacts on stainless-internal branches ([0586ec1](https://github.com/kernel/kernel-go-sdk/commit/0586ec11323a030954b76b8dde236db6817f5141))
* **internal:** codegen related update ([63ecb44](https://github.com/kernel/kernel-go-sdk/commit/63ecb4436a223ba865230b94353a4f654205a4e2))
* update placeholder string ([82510c4](https://github.com/kernel/kernel-go-sdk/commit/82510c497a99e3a671b695667cc59aa66b8605ac))

## 0.42.1 (2026-03-05)

Full Changelog: [v0.42.0...v0.42.1](https://github.com/kernel/kernel-go-sdk/compare/v0.42.0...v0.42.1)

### Features

* [kernel-1028] add api clipboard support ([12e68a1](https://github.com/kernel/kernel-go-sdk/commit/12e68a109a129985f55d0660b858f9840c6c8c0d))
* add force flag to viewport resize to bypass live view/recording check ([d387238](https://github.com/kernel/kernel-go-sdk/commit/d387238e5da1f24f8b54a4be84d484da90fe1226))
* expose smooth mouse movement via public API ([0e0f2ac](https://github.com/kernel/kernel-go-sdk/commit/0e0f2ac548b6155a2e7af9825e16b6a5662ec8cf))


### Bug Fixes

* use indices:dots arrayFmt in MarshalRoot for correct multipart array encoding ([11a7720](https://github.com/kernel/kernel-go-sdk/commit/11a7720da7982b30d072d88f86d1c77db1d4902f))


### Chores

* **internal:** codegen related update ([3c706a4](https://github.com/kernel/kernel-go-sdk/commit/3c706a4e55213c5da4fdf256c77d716729c7f658))

## 0.42.0 (2026-03-02)

Full Changelog: [v0.41.0...v0.42.0](https://github.com/kernel/kernel-go-sdk/compare/v0.41.0...v0.42.0)

### Features

* Neil/kernel 1052 deployments list endpoint ([8062ea7](https://github.com/kernel/kernel-go-sdk/commit/8062ea7a30d6670f950e04045bd3e52cdd7422c0))

## 0.41.0 (2026-02-27)

Full Changelog: [v0.40.0...v0.41.0](https://github.com/kernel/kernel-go-sdk/compare/v0.40.0...v0.41.0)

### Features

* Return uptime_ms for deleted browser sessions ([a22343e](https://github.com/kernel/kernel-go-sdk/commit/a22343ec31502653057ce247be6177155ad4d651))

## 0.40.0 (2026-02-26)

Full Changelog: [v0.39.0...v0.40.0](https://github.com/kernel/kernel-go-sdk/compare/v0.39.0...v0.40.0)

### Features

* show pool browsers in dashboard and API ([587b290](https://github.com/kernel/kernel-go-sdk/commit/587b2902ed1d5e618c050b4a05a255e08858ff1c))

## 0.39.0 (2026-02-25)

Full Changelog: [v0.38.0...v0.39.0](https://github.com/kernel/kernel-go-sdk/compare/v0.38.0...v0.39.0)

### Features

* Add proxy hostname bypass hosts ([2b066bf](https://github.com/kernel/kernel-go-sdk/commit/2b066bfd2e42070cbf2551dc6579854de5760807))

## 0.38.0 (2026-02-25)

Full Changelog: [v0.37.0...v0.38.0](https://github.com/kernel/kernel-go-sdk/compare/v0.37.0...v0.38.0)

### Features

* Neil/kernel 1029 past session search ([04750c0](https://github.com/kernel/kernel-go-sdk/commit/04750c05fd4b0491f34d65cf8e003a7ee38537b5))


### Chores

* **internal:** move custom custom `json` tags to `api` ([a35c88c](https://github.com/kernel/kernel-go-sdk/commit/a35c88ce383ecf7a2a98b303837948cd3edd3d39))

## 0.37.0 (2026-02-23)

Full Changelog: [v0.36.1...v0.37.0](https://github.com/kernel/kernel-go-sdk/compare/v0.36.1...v0.37.0)

### Features

* Neil/kernel 1017 profile pagination query parameter ([d6166c8](https://github.com/kernel/kernel-go-sdk/commit/d6166c8793e9ad1e28d97b211af802d0abc87461))


### Bug Fixes

* improve CLI coverage workflow reliability and accuracy ([4dac142](https://github.com/kernel/kernel-go-sdk/commit/4dac142d082cec5ac906426e9edb00bc6b0c1693))
* remove racy branch existence check in CLI coverage workflow ([32639df](https://github.com/kernel/kernel-go-sdk/commit/32639df54540df7e4d15c32092a565d20b32718e))

## 0.36.1 (2026-02-21)

Full Changelog: [v0.36.0...v0.36.1](https://github.com/kernel/kernel-go-sdk/compare/v0.36.0...v0.36.1)

### Features

* Add version filter to GET /deployments endpoint ([cebd474](https://github.com/kernel/kernel-go-sdk/commit/cebd474525fbafe8b5897f92cc6d8b258f360626))

## 0.36.0 (2026-02-21)

Full Changelog: [v0.35.0...v0.36.0](https://github.com/kernel/kernel-go-sdk/compare/v0.35.0...v0.36.0)

### Features

* Add DELETE /deployments/{id} API endpoint ([759add0](https://github.com/kernel/kernel-go-sdk/commit/759add0fb28e9c823d6df9947f47f7266a346062))
* add live smoke testing to CLI coverage workflow ([fd4f331](https://github.com/kernel/kernel-go-sdk/commit/fd4f3311b14a895782821708acc592019b5f3bfb))


### Bug Fixes

* allow canceling a request while it is waiting to retry ([f04d85f](https://github.com/kernel/kernel-go-sdk/commit/f04d85f0dbf49c9e2629e3f26b4b5b9b8dcb1586))


### Chores

* **internal:** remove mock server code ([0676f28](https://github.com/kernel/kernel-go-sdk/commit/0676f28561cbbec3cfa6ad4209710e2210c1c292))
* update mock server docs ([217f713](https://github.com/kernel/kernel-go-sdk/commit/217f713a30fd2350b3245cf55605b6a5f375a697))

## 0.35.0 (2026-02-18)

Full Changelog: [v0.34.0...v0.35.0](https://github.com/kernel/kernel-go-sdk/compare/v0.34.0...v0.35.0)

### Features

* GPU pools ([d370bcb](https://github.com/kernel/kernel-go-sdk/commit/d370bcbdafb94975cdcc8a9e0d22c75ba4d6d6c6))

## 0.34.0 (2026-02-18)

Full Changelog: [v0.33.0...v0.34.0](https://github.com/kernel/kernel-go-sdk/compare/v0.33.0...v0.34.0)

### Features

* Add error_code field to ManagedAuthSession and related components ([09fcbf2](https://github.com/kernel/kernel-go-sdk/commit/09fcbf20f9dfc15b1b279089a9c7e1ff7f6967fc))
* Allow arbitrary viewport dimensions ([7460146](https://github.com/kernel/kernel-go-sdk/commit/746014666e64ce99d97e5eb44986e6360b93b5e1))
* Neil/kernel 873 templates v4 ([dc79b42](https://github.com/kernel/kernel-go-sdk/commit/dc79b427e161c35f4087cdac62efe507f75acf1c))

## 0.33.0 (2026-02-11)

Full Changelog: [v0.32.0...v0.33.0](https://github.com/kernel/kernel-go-sdk/compare/v0.32.0...v0.33.0)

### Features

* **auth:** add save_credentials support ([dc044eb](https://github.com/kernel/kernel-go-sdk/commit/dc044eb18795520ad5969e2337dac13f20adfbde))
* **auth:** plan-based min health check intervals ([b2f5bc2](https://github.com/kernel/kernel-go-sdk/commit/b2f5bc203e14397ab9b4f2641f42e1286935d0fc))
* Browser API endpoint grouping ([b772828](https://github.com/kernel/kernel-go-sdk/commit/b7728287d7e63d612f6fedd2c4e199bda93440e5))


### Bug Fixes

* **encoder:** correctly serialize NullStruct ([85ac49d](https://github.com/kernel/kernel-go-sdk/commit/85ac49d67563554db23e8d39694d27d381d0b1c3))
* use kernel-internal app token in update-cli-coverage workflow ([789d161](https://github.com/kernel/kernel-go-sdk/commit/789d1616fac976f15e44eae364442fc47cad381e))


### Refactors

* **api:** remove deprecated agent-auth endpoints from stainless.… ([cb2010d](https://github.com/kernel/kernel-go-sdk/commit/cb2010d43ca2a93f21442218d978524329e9a57e))
* **auth:** simplify proxy configuration in OpenAPI schema ([7302183](https://github.com/kernel/kernel-go-sdk/commit/73021835cd291ba1e86172670daa599d49e76172))

## 0.32.0 (2026-02-07)

Full Changelog: [v0.31.1...v0.32.0](https://github.com/kernel/kernel-go-sdk/compare/v0.31.1...v0.32.0)

### Features

* **auth:** add reauth circuit breaker logic ([9cfbd02](https://github.com/kernel/kernel-go-sdk/commit/9cfbd021a7ecf10f4a037ce73e775c5dab2d8fa0))

## 0.31.1 (2026-02-06)

Full Changelog: [v0.31.0...v0.31.1](https://github.com/kernel/kernel-go-sdk/compare/v0.31.0...v0.31.1)

### Chores

* add Managed Auth API planning doc ([fe6c74f](https://github.com/kernel/kernel-go-sdk/commit/fe6c74f86576f1939a61e5ea78a0480837629c81))

## 0.31.0 (2026-02-06)

Full Changelog: [v0.30.0...v0.31.0](https://github.com/kernel/kernel-go-sdk/compare/v0.30.0...v0.31.0)

### Features

* add batch computer action proxy endpoint ([19e11cc](https://github.com/kernel/kernel-go-sdk/commit/19e11cc10d4b9331a5342152ca745d06535809c7))

## 0.30.0 (2026-02-03)

Full Changelog: [v0.29.0...v0.30.0](https://github.com/kernel/kernel-go-sdk/compare/v0.29.0...v0.30.0)

### Features

* Neil/kernel 872 templates v3 ([9a6a1bf](https://github.com/kernel/kernel-go-sdk/commit/9a6a1bf18929ef26c9fe7259dbd59682532d4bdc))

## 0.29.0 (2026-01-29)

Full Changelog: [v0.28.0...v0.29.0](https://github.com/kernel/kernel-go-sdk/compare/v0.28.0...v0.29.0)

### Features

* add support for 1280x800@60 viewport ([0e0add6](https://github.com/kernel/kernel-go-sdk/commit/0e0add6efe25cbea44be8876e5b24c3b6fd2db02))
* **client:** add a convenient param.SetJSON helper ([00e86b2](https://github.com/kernel/kernel-go-sdk/commit/00e86b280372cea0186578a4e73bf857310233e8))

## 0.28.0 (2026-01-22)

Full Changelog: [v0.27.0...v0.28.0](https://github.com/kernel/kernel-go-sdk/compare/v0.27.0...v0.28.0)

### Features

* Allow hot loading profiles into sessions ([436ac8c](https://github.com/kernel/kernel-go-sdk/commit/436ac8cfa797b3283ccb0a801a59fb1fa499e32f))


### Bug Fixes

* preserve existing CLI branch work in coverage workflow ([5c2eb99](https://github.com/kernel/kernel-go-sdk/commit/5c2eb9944cf21b4c9b4870666be6332c52f2a4d9))

## 0.27.0 (2026-01-21)

Full Changelog: [v0.26.0...v0.27.0](https://github.com/kernel/kernel-go-sdk/compare/v0.26.0...v0.27.0)

### Features

* **agent-auth:** add 1Password integration for credential providers ([690962f](https://github.com/kernel/kernel-go-sdk/commit/690962f9f276c1e917c7462ec57400afccf3509f))
* **dashboard:** add browser replays support for past browsers ([33f2a9c](https://github.com/kernel/kernel-go-sdk/commit/33f2a9cd0b5391bd058f6c74900c8d3731b3f990))
* Update browser pool org limits ([6848a1c](https://github.com/kernel/kernel-go-sdk/commit/6848a1c11b3b648b163db8cf1e47347d6b6acbb8))


### Refactors

* **agentauth:** enhance discover and submit modules with improve… ([71801e7](https://github.com/kernel/kernel-go-sdk/commit/71801e7d3375037a68f751099bc6333e592fb654))

## 0.26.0 (2026-01-17)

Full Changelog: [v0.25.0...v0.26.0](https://github.com/kernel/kernel-go-sdk/compare/v0.25.0...v0.26.0)

### Features

* Auth agents auth check URL ([b7ecbfe](https://github.com/kernel/kernel-go-sdk/commit/b7ecbfe68333179633c445aa1eed551f8b0a11d8))


### Bug Fixes

* **docs:** add missing pointer prefix to api.md return types ([ec3f3d7](https://github.com/kernel/kernel-go-sdk/commit/ec3f3d7baa03b861afb53e3c12b1d32aba02b7e0))
* **stainless:** use @onkernel/sdk package name for TypeScript SDK ([3a4991d](https://github.com/kernel/kernel-go-sdk/commit/3a4991d1ef8e51bc1d3909cd6f4725448f765530))
* use setup-go@v6 (not checkout@v6) ([f1e6ffc](https://github.com/kernel/kernel-go-sdk/commit/f1e6ffc8352c405647fb9d4d07299a465a241bad))


### Chores

* **internal:** update `actions/checkout` version ([5a2a7d4](https://github.com/kernel/kernel-go-sdk/commit/5a2a7d471453081cc84a4fb987c6f43c36bb16ac))

## 0.25.0 (2026-01-16)

Full Changelog: [v0.24.0...v0.25.0](https://github.com/kernel/kernel-go-sdk/compare/v0.24.0...v0.25.0)

### Features

* add MFA options to agent authentication workflow ([0f5b637](https://github.com/kernel/kernel-go-sdk/commit/0f5b63778ed3fce2f8662053eb213ba55a77cde9))
* add WebSocket process attach and PTY support ([690d473](https://github.com/kernel/kernel-go-sdk/commit/690d4737fe6b3d13c64acc5de08da76deb74a597))
* **api:** add IP address logging for residential and custom proxies ([35a09ad](https://github.com/kernel/kernel-go-sdk/commit/35a09ad5ece2cdaf09a8d3b6420b4744377f5d6f))
* **api:** manual updates ([0980329](https://github.com/kernel/kernel-go-sdk/commit/0980329671baab3951d42b472f69fcd5ef29aab7))
* **api:** update production repos ([0cdfbfd](https://github.com/kernel/kernel-go-sdk/commit/0cdfbfdfccab1abe18180bce8500623be8bcad66))
* Support hot swap proxy on a session ([5c479d9](https://github.com/kernel/kernel-go-sdk/commit/5c479d9abb0f42edbb4bdf8d0cef681f5aeeca3a))


### Chores

* sync repo ([1728134](https://github.com/kernel/kernel-go-sdk/commit/1728134103bf4a4bcb1c74ad335408cc3f7cc65a))
* update module path from onkernel to kernel ([52c5887](https://github.com/kernel/kernel-go-sdk/commit/52c588700307b752e5d239718a00aafddc86ab96))

## 0.24.0 (2025-12-17)

Full Changelog: [v0.23.0...v0.24.0](https://github.com/onkernel/kernel-go-sdk/compare/v0.23.0...v0.24.0)

### Features

* Enhance AuthAgentInvocationCreateResponse to include already_authenti… ([fd4bfbd](https://github.com/onkernel/kernel-go-sdk/commit/fd4bfbd58128160ef68fd7e8fb0e6e33b92e5c35))
* Fix browser pool sdk types ([68de230](https://github.com/onkernel/kernel-go-sdk/commit/68de23056cdb320c39f54a1b77baacd723dd2dde))

## 0.23.0 (2025-12-16)

Full Changelog: [v0.22.0...v0.23.0](https://github.com/onkernel/kernel-go-sdk/compare/v0.22.0...v0.23.0)

### Features

* **encoder:** support bracket encoding form-data object members ([7d11b85](https://github.com/onkernel/kernel-go-sdk/commit/7d11b85f517bfa5875440c516de735009fbd05a0))
* enhance agent authentication API with new endpoints and request… ([7f2d67a](https://github.com/onkernel/kernel-go-sdk/commit/7f2d67aeba880e8d35b084ffa6f0c2364f88378a))
* Enhance AuthAgent model with last_auth_check_at field ([a3cb1e1](https://github.com/onkernel/kernel-go-sdk/commit/a3cb1e13b7700f1ccce3393411a08f822e216d1f))


### Bug Fixes

* **client:** copy over change to params names to ExecuteNewRequeest ([92fd6c1](https://github.com/onkernel/kernel-go-sdk/commit/92fd6c10cb2c451f9381d3969c2aa80b121addfe))
* **mcp:** correct code tool API endpoint ([563016a](https://github.com/onkernel/kernel-go-sdk/commit/563016aabd8010af7bf2d1db0b7436f31b0fd23e))
* rename param to avoid collision ([a20c158](https://github.com/onkernel/kernel-go-sdk/commit/a20c1588a19c48cfa6503af80c0ba6ad4add76f7))


### Chores

* elide duplicate aliases ([9eb4ec3](https://github.com/onkernel/kernel-go-sdk/commit/9eb4ec3048ab55e534e84de255ca0f2d22d5f233))
* **internal:** codegen related update ([839598d](https://github.com/onkernel/kernel-go-sdk/commit/839598d5a188b64afe16e1df16915c94f822b309))

## 0.22.0 (2025-12-06)

Full Changelog: [v0.21.0...v0.22.0](https://github.com/onkernel/kernel-go-sdk/compare/v0.21.0...v0.22.0)

### Features

* [wip] Browser pools polish pass ([628410a](https://github.com/onkernel/kernel-go-sdk/commit/628410a6fc946e6adf5d7a9593b6ebb01385f918))
* Add `async_timeout_seconds` to PostInvocations ([4ee4d7c](https://github.com/onkernel/kernel-go-sdk/commit/4ee4d7c60817d8d4a9e6460e7d04c19f309095d3))
* Enhance agent authentication with optional login page URL and auth ch… ([a2fc4ac](https://github.com/onkernel/kernel-go-sdk/commit/a2fc4acfb2bf88a5cc209d5430fcfb0eee4f751b))


### Refactors

* **browser:** remove persistence option UI ([7ea9ccf](https://github.com/onkernel/kernel-go-sdk/commit/7ea9ccfb9a6c8acb3fd9d8d04b7bf8b9cfb4fef1))

## 0.21.0 (2025-12-02)

Full Changelog: [v0.20.0...v0.21.0](https://github.com/onkernel/kernel-go-sdk/compare/v0.20.0...v0.21.0)

### Features

* Browser pools sdk release ([69dcf34](https://github.com/onkernel/kernel-go-sdk/commit/69dcf3471d1b8b7b6fb8210da6d9c0b42ab1f1e7))
* Mason/agent auth api ([ee63aba](https://github.com/onkernel/kernel-go-sdk/commit/ee63aba391e6a3024e93eeca28d008ea93b18b58))

## 0.20.0 (2025-11-19)

Full Changelog: [v0.19.2...v0.20.0](https://github.com/onkernel/kernel-go-sdk/compare/v0.19.2...v0.20.0)

### Features

* Add pagination to list browsers method and allow it to include deleted browsers when `include_deleted = true` ([2bebf78](https://github.com/onkernel/kernel-go-sdk/commit/2bebf78af2a6a385743700fd8342f2c4f6efd20f))


### Bug Fixes

* **client:** correctly specify Accept header with */* instead of empty ([2a95394](https://github.com/onkernel/kernel-go-sdk/commit/2a95394194865cffa351c650924f41dc6262fce8))

## 0.19.2 (2025-11-17)

Full Changelog: [v0.19.1...v0.19.2](https://github.com/onkernel/kernel-go-sdk/compare/v0.19.1...v0.19.2)

### Features

* Feat increase max timeout to 72h ([f824eae](https://github.com/onkernel/kernel-go-sdk/commit/f824eae83fecc41468817b110b7de0869f1e40be))

## 0.19.1 (2025-11-13)

Full Changelog: [v0.19.0...v0.19.1](https://github.com/onkernel/kernel-go-sdk/compare/v0.19.0...v0.19.1)

### Features

* works locally ([1305424](https://github.com/onkernel/kernel-go-sdk/commit/13054249676464251daeb156138ffb9efd895292))

## 0.19.0 (2025-11-12)

Full Changelog: [v0.18.0...v0.19.0](https://github.com/onkernel/kernel-go-sdk/compare/v0.18.0...v0.19.0)

### Features

* feat hide cursor v2 ([83ccf6e](https://github.com/onkernel/kernel-go-sdk/commit/83ccf6eeaf25d7257cf6992af5b856b1c82fa8c6))
* Remove price gating on computer endpoints ([131fda4](https://github.com/onkernel/kernel-go-sdk/commit/131fda473e619b41bf8df8a0c461e5594142b474))


### Chores

* bump gjson version ([666f492](https://github.com/onkernel/kernel-go-sdk/commit/666f4924d4ec1b8ebe5e083a8df701c19ee414be))
* **internal:** grammar fix (it's -&gt; its) ([be376e8](https://github.com/onkernel/kernel-go-sdk/commit/be376e843484cf6aee031b8ace289f166cd181e1))

## 0.18.0 (2025-10-30)

Full Changelog: [v0.17.0...v0.18.0](https://github.com/onkernel/kernel-go-sdk/compare/v0.17.0...v0.18.0)

### Features

* apps: add offset pagination + headers ([bb1bb22](https://github.com/onkernel/kernel-go-sdk/commit/bb1bb22bab1970f1c2ef0455d3002abbfc46466b))

## 0.17.0 (2025-10-27)

Full Changelog: [v0.16.0...v0.17.0](https://github.com/onkernel/kernel-go-sdk/compare/v0.16.0...v0.17.0)

### Features

* Make country flag optional for DC and ISP proxies ([cd2a114](https://github.com/onkernel/kernel-go-sdk/commit/cd2a1144f56f7dbe8a8dce6e918eb9a311e7ec3d))

## 0.16.0 (2025-10-27)

Full Changelog: [v0.15.0...v0.16.0](https://github.com/onkernel/kernel-go-sdk/compare/v0.15.0...v0.16.0)

### Features

* ad hoc playwright code exec AP| ([8589fb2](https://github.com/onkernel/kernel-go-sdk/commit/8589fb22265ec8ad4c5a02b7a485c78e8a5f45da))

## 0.15.0 (2025-10-17)

Full Changelog: [v0.14.2...v0.15.0](https://github.com/onkernel/kernel-go-sdk/compare/v0.14.2...v0.15.0)

### Features

* click mouse, move mouse, screenshot ([b1dcded](https://github.com/onkernel/kernel-go-sdk/commit/b1dcdedf9662f112048d27622a787b2bacb5ac84))
* Phani/deploy with GitHub url ([5e04228](https://github.com/onkernel/kernel-go-sdk/commit/5e04228a5fdd20381b24b230c64e65761a444300))

## 0.14.2 (2025-10-16)

Full Changelog: [v0.14.1...v0.14.2](https://github.com/onkernel/kernel-go-sdk/compare/v0.14.1...v0.14.2)

### Features

* Kiosk mode ([cd6120a](https://github.com/onkernel/kernel-go-sdk/commit/cd6120a71b58e0c539d5ba1361fca0ecd71baf1c))

## 0.14.1 (2025-10-13)

Full Changelog: [v0.14.1...v0.14.1](https://github.com/onkernel/kernel-go-sdk/compare/v0.14.1...v0.14.1)

### Features

* Hide and deprecate mobile proxy type ([0775dfb](https://github.com/onkernel/kernel-go-sdk/commit/0775dfbadceda532968a0bfc24b40e13eb991b41))
* WIP: Configurable Viewport ([b8a3e85](https://github.com/onkernel/kernel-go-sdk/commit/b8a3e85900d690a5a439d7c46fde7180805d6f2c))

## 0.14.1 (2025-10-07)

Full Changelog: [v0.14.0...v0.14.1](https://github.com/onkernel/kernel-go-sdk/compare/v0.14.0...v0.14.1)

### Features

* WIP browser extensions ([22323cd](https://github.com/onkernel/kernel-go-sdk/commit/22323cdf2376bf2016cf1ebd1384f4c1b5ca6752))

## 0.14.0 (2025-10-03)

Full Changelog: [v0.13.0...v0.14.0](https://github.com/onkernel/kernel-go-sdk/compare/v0.13.0...v0.14.0)

### Features

* Http proxy ([0ebe5f1](https://github.com/onkernel/kernel-go-sdk/commit/0ebe5f1f007fd53267ac096e044f9fd513c28d7b))

## 0.13.0 (2025-10-01)

Full Changelog: [v0.12.0...v0.13.0](https://github.com/onkernel/kernel-go-sdk/compare/v0.12.0...v0.13.0)

### Features

* Update oAPI and data model for proxy status ([87b8365](https://github.com/onkernel/kernel-go-sdk/commit/87b8365a941544eaf38b09ca66d6f9a97e3b050b))

## 0.12.0 (2025-09-30)

Full Changelog: [v0.11.5...v0.12.0](https://github.com/onkernel/kernel-go-sdk/compare/v0.11.5...v0.12.0)

### Features

* Return proxy ID in browsers response ([711c52f](https://github.com/onkernel/kernel-go-sdk/commit/711c52facea0f2170938a0a68040c6781f2c19ee))

## 0.11.5 (2025-09-29)

Full Changelog: [v0.11.4...v0.11.5](https://github.com/onkernel/kernel-go-sdk/compare/v0.11.4...v0.11.5)

### Features

* Add App Version to Invocation and add filtering on App Version ([81f87e2](https://github.com/onkernel/kernel-go-sdk/commit/81f87e27ae242be9f7cfc75a6147e9eac669d4c4))
* Fix my incorrect grammer ([f04186f](https://github.com/onkernel/kernel-go-sdk/commit/f04186fd786d47e49c42028c0237e02fd08d03b1))


### Bug Fixes

* bugfix for setting JSON keys with special characters ([32304ba](https://github.com/onkernel/kernel-go-sdk/commit/32304baca3b01de677fd705da2b03787b56fdf35))

## 0.11.4 (2025-09-25)

Full Changelog: [v0.11.3...v0.11.4](https://github.com/onkernel/kernel-go-sdk/compare/v0.11.3...v0.11.4)

### Features

* getInvocations endpoint ([b9a983c](https://github.com/onkernel/kernel-go-sdk/commit/b9a983c1cf96d2aa6d22f57aaf98ded129e8d4d4))

## 0.11.3 (2025-09-24)

Full Changelog: [v0.11.2...v0.11.3](https://github.com/onkernel/kernel-go-sdk/compare/v0.11.2...v0.11.3)

### Features

* Per Invocation Logs ([f1241f2](https://github.com/onkernel/kernel-go-sdk/commit/f1241f2d66a097bfb11a2cbd9fda9495b54b0690))

## 0.11.2 (2025-09-24)

Full Changelog: [v0.11.1...v0.11.2](https://github.com/onkernel/kernel-go-sdk/compare/v0.11.1...v0.11.2)

### Features

* Add stainless CI ([ac377ec](https://github.com/onkernel/kernel-go-sdk/commit/ac377ecde22bf08aa41dbddf77cbfd686354c1ea))


### Bug Fixes

* use slices.Concat instead of sometimes modifying r.Options ([e893e50](https://github.com/onkernel/kernel-go-sdk/commit/e893e50418d4dcbd11b910f12fad79f53ca5eb85))


### Chores

* bump minimum go version to 1.22 ([8f9a2ac](https://github.com/onkernel/kernel-go-sdk/commit/8f9a2ac0e091826c5ed578e568a250b4f11cc05a))
* do not install brew dependencies in ./scripts/bootstrap by default ([856b62a](https://github.com/onkernel/kernel-go-sdk/commit/856b62a529e4bed1ae8205065dc6c0e6e50e5f6d))
* update more docs for 1.22 ([f38538f](https://github.com/onkernel/kernel-go-sdk/commit/f38538f4031a5d8045b949d6e9b06f5c76cecde9))

## 0.11.1 (2025-09-06)

Full Changelog: [v0.11.0...v0.11.1](https://github.com/onkernel/kernel-go-sdk/compare/v0.11.0...v0.11.1)

### Features

* **api:** add pagination to the deployments endpoint ([9130f8c](https://github.com/onkernel/kernel-go-sdk/commit/9130f8c40783560bc02d5b494c60549dbe092dcd))
* **api:** pagination properties added to response (has_more, next_offset) ([65c5342](https://github.com/onkernel/kernel-go-sdk/commit/65c53429690169c3ca5068b94991af6eba3832ff))
* **api:** update API spec with pagination headers ([ad37eb2](https://github.com/onkernel/kernel-go-sdk/commit/ad37eb2bdc92a4145e5dba9d8ef8b4ad26c9b64a))


### Bug Fixes

* **client:** correctly convert header pagination value to int ([d42bd8a](https://github.com/onkernel/kernel-go-sdk/commit/d42bd8a79dc3f64ea6aa8c3f4239b8cfa8545d82))
* **internal:** unmarshal correctly when there are multiple discriminators ([649203d](https://github.com/onkernel/kernel-go-sdk/commit/649203d6621f6850ac4cb414b04b78c9693d1923))

## 0.11.0 (2025-09-04)

Full Changelog: [v0.10.0...v0.11.0](https://github.com/onkernel/kernel-go-sdk/compare/v0.10.0...v0.11.0)

### Features

* **api:** adding support for browser profiles ([481cdb3](https://github.com/onkernel/kernel-go-sdk/commit/481cdb3500744c9e4ec050e340a920302d8fea19))


### Bug Fixes

* close body before retrying ([a6a2e40](https://github.com/onkernel/kernel-go-sdk/commit/a6a2e4054c629d6ee85997ed81a1b14e70e594dc))


### Chores

* **internal:** codegen related update ([a7030ab](https://github.com/onkernel/kernel-go-sdk/commit/a7030abb99c06c675f60a4f2afde43d376d9981f))

## 0.10.0 (2025-08-27)

Full Changelog: [v0.9.1...v0.10.0](https://github.com/onkernel/kernel-go-sdk/compare/v0.9.1...v0.10.0)

### Features

* **api:** new process, fs, and log endpoints ([fa85f19](https://github.com/onkernel/kernel-go-sdk/commit/fa85f19a8f4c664696d8839539a6ce6172a78e98))

## 0.9.1 (2025-08-15)

Full Changelog: [v0.9.0...v0.9.1](https://github.com/onkernel/kernel-go-sdk/compare/v0.9.0...v0.9.1)

### Features

* **api:** add browser timeouts ([6c7b47f](https://github.com/onkernel/kernel-go-sdk/commit/6c7b47f69ccc4e12d9e21340c543d57b7fc6d314))

### Chores

* **internal:** codegen related update ([e4ca558](https://github.com/onkernel/kernel-go-sdk/commit/e4ca55843e4dbb9b7e71821ca58080a5bf25f025))
* **internal:** update comment in script ([9542333](https://github.com/onkernel/kernel-go-sdk/commit/9542333108abb522bae00b266c89cc3917884b35))
* update @stainless-api/prism-cli to v5.15.0 ([625476c](https://github.com/onkernel/kernel-go-sdk/commit/625476c16f8298f7fa1c18318f67231906d89a56))

## 0.9.0 (2025-08-07)

Full Changelog: [v0.8.2...v0.9.0](https://github.com/onkernel/kernel-go-sdk/compare/v0.8.2...v0.9.0)

### Features

* **api:** browser instance file i/o ([7308fb8](https://github.com/onkernel/kernel-go-sdk/commit/7308fb8acdcb2d10fb9d09612854a18f18083bcc))
* **client:** support optional json html escaping ([7fb83d1](https://github.com/onkernel/kernel-go-sdk/commit/7fb83d1f154d9502833b4f0a17c9912aa9643c65))

## 0.8.2 (2025-07-23)

Full Changelog: [v0.8.1...v0.8.2](https://github.com/onkernel/kernel-go-sdk/compare/v0.8.1...v0.8.2)

### Features

* **api:** add action name to the response to invoke ([e1bfb1d](https://github.com/onkernel/kernel-go-sdk/commit/e1bfb1d1874a0ad1dbb61c035e3f728b9406df5a))


### Bug Fixes

* **client:** process custom base url ahead of time ([4bef9af](https://github.com/onkernel/kernel-go-sdk/commit/4bef9afe3097d7d07596799efae41ca2acbff038))

## 0.8.1 (2025-07-21)

Full Changelog: [v0.8.0...v0.8.1](https://github.com/onkernel/kernel-go-sdk/compare/v0.8.0...v0.8.1)

### Chores

* **api:** remove deprecated endpoints ([7a659dc](https://github.com/onkernel/kernel-go-sdk/commit/7a659dcfeba2f3881b84934c2ee4780c4fbbb852))

## 0.8.0 (2025-07-16)

Full Changelog: [v0.7.1...v0.8.0](https://github.com/onkernel/kernel-go-sdk/compare/v0.7.1...v0.8.0)

### Features

* **api:** manual updates ([ba0a473](https://github.com/onkernel/kernel-go-sdk/commit/ba0a473799c2c9a00865112a22ae356c57d2c114))
* **client:** expand max streaming buffer size ([41bdefe](https://github.com/onkernel/kernel-go-sdk/commit/41bdefed3d43cb52b431e8255a4a4faa0dfc46f0))


### Chores

* **internal:** fix lint script for tests ([5118fdf](https://github.com/onkernel/kernel-go-sdk/commit/5118fdf7e31a237e8def04c5a41cd0767900deef))
* lint tests in subpackages ([71d9002](https://github.com/onkernel/kernel-go-sdk/commit/71d900220761ff811187d5c03236a0fc5d92fcc6))

## 0.7.1 (2025-07-08)

Full Changelog: [v0.7.0...v0.7.1](https://github.com/onkernel/kernel-go-sdk/compare/v0.7.0...v0.7.1)

### Features

* **api:** manual updates ([2162b0a](https://github.com/onkernel/kernel-go-sdk/commit/2162b0a2396b520acf176b8622b8728e51f8787d))


### Chores

* lint tests ([29d98ac](https://github.com/onkernel/kernel-go-sdk/commit/29d98ac8ec7d9fe704d07ac0bce4b60d0f2ed121))

## 0.7.0 (2025-07-02)

Full Changelog: [v0.6.5...v0.7.0](https://github.com/onkernel/kernel-go-sdk/compare/v0.6.5...v0.7.0)

### Features

* **api:** headless browsers ([e7f85f9](https://github.com/onkernel/kernel-go-sdk/commit/e7f85f9432b053b513b1581cb2ece0dd78ad3613))

## 0.6.5 (2025-07-02)

Full Changelog: [v0.6.4...v0.6.5](https://github.com/onkernel/kernel-go-sdk/compare/v0.6.4...v0.6.5)

### Bug Fixes

* don't try to deserialize as json when ResponseBodyInto is []byte ([7865ba8](https://github.com/onkernel/kernel-go-sdk/commit/7865ba807c2766ccdd8cb3d59bb8855721b39769))


### Chores

* **ci:** only run for pushes and fork pull requests ([a6b4ae7](https://github.com/onkernel/kernel-go-sdk/commit/a6b4ae749a0e11f8683dd3c354c1e3b3003d738e))

## 0.6.4 (2025-06-27)

Full Changelog: [v0.6.3...v0.6.4](https://github.com/onkernel/kernel-go-sdk/compare/v0.6.3...v0.6.4)

### Features

* **api:** add GET deployments endpoint ([e2c8f14](https://github.com/onkernel/kernel-go-sdk/commit/e2c8f14da32a692adbf18162ceab4b89e1bb2c4f))
* **api:** deployments ([e0dc9c0](https://github.com/onkernel/kernel-go-sdk/commit/e0dc9c0fc621dabb383975656e58a19bbf8c0714))
* **api:** manual updates ([f2ddf2c](https://github.com/onkernel/kernel-go-sdk/commit/f2ddf2c318eb1254e8c15fb5723a92b0ebd9cbb4))

## 0.6.3 (2025-06-25)

Full Changelog: [v0.6.2...v0.6.3](https://github.com/onkernel/kernel-go-sdk/compare/v0.6.2...v0.6.3)

### Features

* **api:** /browsers no longer requires invocation id ([179da50](https://github.com/onkernel/kernel-go-sdk/commit/179da501250b0950198690bf54fabf046264c3f5))

## 0.6.2 (2025-06-24)

Full Changelog: [v0.6.1...v0.6.2](https://github.com/onkernel/kernel-go-sdk/compare/v0.6.1...v0.6.2)

### Features

* **api:** add `since` parameter to deployment logs endpoint ([dc72c81](https://github.com/onkernel/kernel-go-sdk/commit/dc72c81a2e45918d595c6c00843b1a1d0efffdd0))
* **client:** add escape hatch for null slice & maps ([d5e1ad9](https://github.com/onkernel/kernel-go-sdk/commit/d5e1ad9087aecd6b67369b9ebbeb633ad808c129))


### Chores

* fix documentation of null map ([a62b964](https://github.com/onkernel/kernel-go-sdk/commit/a62b9647386501f43e70ad876dc5c0271c4c4709))

## 0.6.1 (2025-06-18)

Full Changelog: [v0.6.0...v0.6.1](https://github.com/onkernel/kernel-go-sdk/compare/v0.6.0...v0.6.1)

### Features

* **api:** add delete_browsers endpoint ([57fb070](https://github.com/onkernel/kernel-go-sdk/commit/57fb070819a948f64af4882595f99fa37f462c70))

## 0.6.0 (2025-06-17)

Full Changelog: [v0.5.0...v0.6.0](https://github.com/onkernel/kernel-go-sdk/compare/v0.5.0...v0.6.0)

### Features

* **api:** update via SDK Studio ([ee8e77a](https://github.com/onkernel/kernel-go-sdk/commit/ee8e77a653cc83084eec067e1060b80e8e99fc27))
* **api:** update via SDK Studio ([1b2efd8](https://github.com/onkernel/kernel-go-sdk/commit/1b2efd8f083758e39ed2bd9bec8d5f6da0642ece))
* **api:** update via SDK Studio ([3094578](https://github.com/onkernel/kernel-go-sdk/commit/3094578598220a1b837274b20e0dd4cb4f36fd2c))
* **api:** update via SDK Studio ([0fec1e5](https://github.com/onkernel/kernel-go-sdk/commit/0fec1e5e8d6b426ade8030ec6836142ddd18cca4))
* **api:** update via SDK Studio ([b44c90e](https://github.com/onkernel/kernel-go-sdk/commit/b44c90ed4698ca9c849bb20aa85579d0fcd36736))
* **api:** update via SDK Studio ([d6e1cd6](https://github.com/onkernel/kernel-go-sdk/commit/d6e1cd618e324beb797495636a79016aa9cfe7b1))
* **api:** update via SDK Studio ([0bc5b00](https://github.com/onkernel/kernel-go-sdk/commit/0bc5b008aad753e6b2a5bf2647b30b389e36c190))
* **api:** update via SDK Studio ([db224f1](https://github.com/onkernel/kernel-go-sdk/commit/db224f133e8935217833037397839b771e901885))
* **api:** update via SDK Studio ([3be7afb](https://github.com/onkernel/kernel-go-sdk/commit/3be7afb4985902b399ed45192a139b8f2a67bcf2))
* **api:** update via SDK Studio ([562f248](https://github.com/onkernel/kernel-go-sdk/commit/562f248c183bc6bd908433c2920b63b596572711))
* **api:** update via SDK Studio ([984cd7f](https://github.com/onkernel/kernel-go-sdk/commit/984cd7f8a9cbbb153e11895e5b8dd5ba79cc5548))
* **api:** update via SDK Studio ([48d0b3d](https://github.com/onkernel/kernel-go-sdk/commit/48d0b3d6938a8e4e9db9148b7dc35d7440e43ae2))
* **api:** update via SDK Studio ([8c6285b](https://github.com/onkernel/kernel-go-sdk/commit/8c6285b9dfff61a9f6b58053026ac18fbde1b66e))
* **client:** add debug log helper ([5c92a71](https://github.com/onkernel/kernel-go-sdk/commit/5c92a71246d10274fa60eb84a0f42121cbca6e7b))


### Chores

* **ci:** enable for pull requests ([72fd885](https://github.com/onkernel/kernel-go-sdk/commit/72fd885ebc2a949148fe7fbc5bc809b716a783fe))

## 0.5.0 (2025-06-04)

Full Changelog: [v0.4.0...v0.5.0](https://github.com/onkernel/kernel-go-sdk/compare/v0.4.0...v0.5.0)

### Features

* **api:** update via SDK Studio ([f0b66f8](https://github.com/onkernel/kernel-go-sdk/commit/f0b66f87541f4bc109f479906c3554bffb38c26f))
* **api:** update via SDK Studio ([d60d333](https://github.com/onkernel/kernel-go-sdk/commit/d60d3337306d90c8c51a54c6abeece2731dd4834))
* **client:** allow overriding unions ([51d877d](https://github.com/onkernel/kernel-go-sdk/commit/51d877d1ae7584b237fdb9238e97949b832c3cc8))


### Bug Fixes

* **client:** cast to raw message when converting to params ([de14358](https://github.com/onkernel/kernel-go-sdk/commit/de14358a78248e642c03bd669f9c361bf6b3c8ba))
* fix error ([7d27985](https://github.com/onkernel/kernel-go-sdk/commit/7d2798511671e285d41943c51c3ac8ebc6a7d6d6))


### Chores

* make go mod tidy continue on error ([91b5021](https://github.com/onkernel/kernel-go-sdk/commit/91b50217a8a08dd06ba64e05d025386d91d586c8))

## 0.4.0 (2025-05-28)

Full Changelog: [v0.3.0...v0.4.0](https://github.com/onkernel/kernel-go-sdk/compare/v0.3.0...v0.4.0)

### Features

* **api:** update via SDK Studio ([ac5cf50](https://github.com/onkernel/kernel-go-sdk/commit/ac5cf50867c042d1da5329f2441855ab89efd686))

## 0.3.0 (2025-05-22)

Full Changelog: [v0.2.0...v0.3.0](https://github.com/onkernel/kernel-go-sdk/compare/v0.2.0...v0.3.0)

### Features

* **api:** update via SDK Studio ([7c22f9e](https://github.com/onkernel/kernel-go-sdk/commit/7c22f9efecdb2548ce58f81454817a33fc249128))
* **api:** update via SDK Studio ([7941c72](https://github.com/onkernel/kernel-go-sdk/commit/7941c7280d3cc6870947d59c96d5e1c55e0ade46))
* **api:** update via SDK Studio ([470093a](https://github.com/onkernel/kernel-go-sdk/commit/470093adf03fc0b077d11f11f13ee802ae4557a6))
* **api:** update via SDK Studio ([5c567c9](https://github.com/onkernel/kernel-go-sdk/commit/5c567c90794bfd2d997685e5102835d65de029e6))


### Chores

* **docs:** grammar improvements ([cb762f8](https://github.com/onkernel/kernel-go-sdk/commit/cb762f85ee2aef349829477128d84e63e41d4449))
* improve devcontainer setup ([8b46076](https://github.com/onkernel/kernel-go-sdk/commit/8b46076b6e1145d9f4ff14e089aa57588f8ef613))

## 0.2.0 (2025-05-21)

Full Changelog: [v0.1.0...v0.2.0](https://github.com/onkernel/kernel-go-sdk/compare/v0.1.0...v0.2.0)

### Features

* **api:** update via SDK Studio ([788354c](https://github.com/onkernel/kernel-go-sdk/commit/788354c6debdf1f4608c31acaa6bf422af3f6306))

## 0.1.0 (2025-05-21)

Full Changelog: [v0.1.0-alpha.6...v0.1.0](https://github.com/onkernel/kernel-go-sdk/compare/v0.1.0-alpha.6...v0.1.0)

### Features

* **api:** update via SDK Studio ([e05f5ec](https://github.com/onkernel/kernel-go-sdk/commit/e05f5ec48a0139acb5f5789124a15916324ceebb))

## 0.1.0-alpha.6 (2025-05-20)

Full Changelog: [v0.1.0-alpha.5...v0.1.0-alpha.6](https://github.com/onkernel/kernel-go-sdk/compare/v0.1.0-alpha.5...v0.1.0-alpha.6)

### Features

* **api:** update via SDK Studio ([8447b94](https://github.com/onkernel/kernel-go-sdk/commit/8447b94117d58db98c50e1e8248121d38c7afd2e))
* **api:** update via SDK Studio ([da5cfff](https://github.com/onkernel/kernel-go-sdk/commit/da5cfff4d62c3dbbac4c2fded65505368097f5a2))

## 0.1.0-alpha.5 (2025-05-20)

Full Changelog: [v0.1.0-alpha.4...v0.1.0-alpha.5](https://github.com/onkernel/kernel-go-sdk/compare/v0.1.0-alpha.4...v0.1.0-alpha.5)

### Features

* **api:** update via SDK Studio ([2f5051e](https://github.com/onkernel/kernel-go-sdk/commit/2f5051e328b157c2836d93fcd1e957c315f9d4e7))
* **api:** update via SDK Studio ([368ec0a](https://github.com/onkernel/kernel-go-sdk/commit/368ec0a4e3eb54fd5634327f6234d9b4f0a83a95))

## 0.1.0-alpha.4 (2025-05-20)

Full Changelog: [v0.1.0-alpha.3...v0.1.0-alpha.4](https://github.com/onkernel/kernel-go-sdk/compare/v0.1.0-alpha.3...v0.1.0-alpha.4)

### Features

* **api:** update via SDK Studio ([34dfea0](https://github.com/onkernel/kernel-go-sdk/commit/34dfea07fd8538b3a73f8ee7ea35f3a70e679b3b))


### Bug Fixes

* **client:** correctly set stream key for multipart ([2336290](https://github.com/onkernel/kernel-go-sdk/commit/23362905d1acb3616333407f70bdee259dc66ebe))
* **client:** don't panic on marshal with extra null field ([4208ce7](https://github.com/onkernel/kernel-go-sdk/commit/4208ce7bb382ac9bfb2cefeb406be6302df6cd7a))
* **client:** increase max stream buffer size ([eca7429](https://github.com/onkernel/kernel-go-sdk/commit/eca74291efed0e090b1a7437d832ee12a819f833))
* **client:** use scanner for streaming ([6fad3f9](https://github.com/onkernel/kernel-go-sdk/commit/6fad3f9d43fbb49c718a53ae6736ac273e43b7b3))

## 0.1.0-alpha.3 (2025-05-19)

Full Changelog: [v0.1.0-alpha.2...v0.1.0-alpha.3](https://github.com/onkernel/kernel-go-sdk/compare/v0.1.0-alpha.2...v0.1.0-alpha.3)

### Features

* **api:** update via SDK Studio ([4a12204](https://github.com/onkernel/kernel-go-sdk/commit/4a122041257d78d3b583408e912d4d8a5b10cdc0))

## 0.1.0-alpha.2 (2025-05-19)

Full Changelog: [v0.1.0-alpha.1...v0.1.0-alpha.2](https://github.com/onkernel/kernel-go-sdk/compare/v0.1.0-alpha.1...v0.1.0-alpha.2)

### Features

* **api:** update via SDK Studio ([72da2c4](https://github.com/onkernel/kernel-go-sdk/commit/72da2c4b53d8f47d6cc0dab3cfbe8707edb989d2))
* **api:** update via SDK Studio ([58da3af](https://github.com/onkernel/kernel-go-sdk/commit/58da3af97eb8b295a54e22e53d8035f00fe09215))
* **api:** update via SDK Studio ([ebab506](https://github.com/onkernel/kernel-go-sdk/commit/ebab506e3b5f8dbe8d30a47b3734c64024341d87))

## 0.1.0-alpha.1 (2025-05-14)

Full Changelog: [v0.0.1-alpha.0...v0.1.0-alpha.1](https://github.com/onkernel/kernel-go-sdk/compare/v0.0.1-alpha.0...v0.1.0-alpha.1)

### Features

* **api:** update via SDK Studio ([bb5cfe4](https://github.com/onkernel/kernel-go-sdk/commit/bb5cfe49047afce492f2a21ade2acbe1612e9f3c))


### Chores

* configure new SDK language ([dd0120b](https://github.com/onkernel/kernel-go-sdk/commit/dd0120b0ce673e0fe7c842d39e91f01b8ee8106a))
* update SDK settings ([b1b8645](https://github.com/onkernel/kernel-go-sdk/commit/b1b8645621de70f28e3cf0b3622c1a93159b561e))
