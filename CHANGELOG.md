## [5.0.5](https://github.com/zenrocklabs/zrchain/compare/v5.0.4...v5.0.5) (2024-11-04)


### Bug Fixes

* regressions (zenBTC e2e flow) ([0db15e3](https://github.com/zenrocklabs/zrchain/commit/0db15e3cf1509dbf0e18372de66e101fc5fdd616))
* regressions (zenBTC e2e flow) ([0608f45](https://github.com/zenrocklabs/zrchain/commit/0608f458adc1fa25847331f868a027cda69e9ebd))

## [5.0.4](https://github.com/zenrocklabs/zrchain/compare/v5.0.3...v5.0.4) (2024-11-04)


### Bug Fixes

* Dockerfile sha256sum ([2742ab1](https://github.com/zenrocklabs/zrchain/commit/2742ab1e39efa3a05ef92219c06c8e8c2f2e920c))
* libwasm version bump ([8acb93c](https://github.com/zenrocklabs/zrchain/commit/8acb93cc67fa8ebf60d8d4e93b9daa145edde58a))

## [5.0.3](https://github.com/zenrocklabs/zrchain/compare/v5.0.2...v5.0.3) (2024-11-04)


### Bug Fixes

* ci pipeline include zip file ([e05990a](https://github.com/zenrocklabs/zrchain/commit/e05990a7e07fd632c07ee686ea1a07336f0f5ef2))

## [5.0.2](https://github.com/zenrocklabs/zrchain/compare/v5.0.1...v5.0.2) (2024-11-04)


### Bug Fixes

* improve tests ([ba57a0b](https://github.com/zenrocklabs/zrchain/commit/ba57a0ba0c672922c33f8561575a3a7857082e93))
* remove gardia-1 related file ([9b1fd5d](https://github.com/zenrocklabs/zrchain/commit/9b1fd5d1556cf4c432bef4942acb52b6e26a1fc2))

## [5.0.1](https://github.com/zenrocklabs/zrchain/compare/v5.0.0...v5.0.1) (2024-11-01)


### Bug Fixes

* v5 upgrade-related issues ([57bf22f](https://github.com/zenrocklabs/zrchain/commit/57bf22f27c7b9cf77cc0f963d4bdd9caac33dfad))
* v5 upgrade-related issues ([8d77806](https://github.com/zenrocklabs/zrchain/commit/8d778061c229507674e98a1633ba299dfd8164ae))

# [5.0.0](https://github.com/zenrocklabs/zrchain/compare/v4.16.1...v5.0.0) (2024-11-01)


### Bug Fixes

* bump `zenrock-avs` ver ([bac81ce](https://github.com/zenrocklabs/zrchain/commit/bac81ce33691ade27307512072dd98745b771283))
* cleanup, fix ethnonce test ([7edafdc](https://github.com/zenrocklabs/zrchain/commit/7edafdcfab92897b8a2f4aad27df3210487a32da))
* Fix for Solana data encoding ([222000b](https://github.com/zenrocklabs/zrchain/commit/222000bf6241eecb1914863aa77842db56b95792))
* migration script ([34ca604](https://github.com/zenrocklabs/zrchain/commit/34ca604d92abb68c83e27e855f5e4e37cb2e1865))
* named upgrade to v5 ([e5980aa](https://github.com/zenrocklabs/zrchain/commit/e5980aa21e5532e6b943ff3cc4345b959cee1385))
* populate `ZenbtcMetadata` field correctly ([6ad7c03](https://github.com/zenrocklabs/zrchain/commit/6ad7c03ab97411f24432daa468b76116b8a184b3))
* temporarily comment out zenbtc keeper `app.go` ([3ecb497](https://github.com/zenrocklabs/zrchain/commit/3ecb4970c88bfc3690079c1ccf7176ba052be148))
* temporarily comment out zenbtc keeper `app.go` ([8ef9e81](https://github.com/zenrocklabs/zrchain/commit/8ef9e815b393df7f630c62233b7b3e45a9ce8375))
* typos in protos ([29c0bf0](https://github.com/zenrocklabs/zrchain/commit/29c0bf03c0c2e2ef7cc51d6778e4465e22c84e16))
* zenBTC work, VE refactor/cleanup + temporarily disable some AVS registration code ([00f1255](https://github.com/zenrocklabs/zrchain/commit/00f1255e464605679bbddd622ae3ebdaf4873960))


### Features

* add `go-client` ([e1a2b11](https://github.com/zenrocklabs/zrchain/commit/e1a2b11a218b2b2a3610d7695c7a408d20506426))
* add `go-client` ([e01c780](https://github.com/zenrocklabs/zrchain/commit/e01c780fa842148aa40301fa2a7748bb8f9ccb5a))
* add `make sidecar`, update `go.mod` ([95dd050](https://github.com/zenrocklabs/zrchain/commit/95dd0504433b3cd97ccf5bc7c445de5be5196f4c))
* adjusting docker-build ([fd7ff3b](https://github.com/zenrocklabs/zrchain/commit/fd7ff3b31233e03e8fffa1076607928c2ae209e6))
* bump zrchain v4 to v5 ([98d8ec0](https://github.com/zenrocklabs/zrchain/commit/98d8ec0d2769ac340e44ac369dc590b62f9bb215))
* trigger CI (BREAKING CHANGE) ([6bb664f](https://github.com/zenrocklabs/zrchain/commit/6bb664f7530a0f67c99d5d03a4eb5d57910ae7ec))
* trigger CI for v5 (BREAKING CHANGE) ([cff3ea4](https://github.com/zenrocklabs/zrchain/commit/cff3ea4255b0d68da7e58d8d0fd5e233a126e697))
* v4 upgrade files ([ac640dc](https://github.com/zenrocklabs/zrchain/commit/ac640dc504c3ed490f39247abad7a326a2cd9f2e))
* working e2e zenBTC deposit & mint flow ([d0af536](https://github.com/zenrocklabs/zrchain/commit/d0af53679c17073743772dc6a239b25c521ceda6))


### BREAKING CHANGES

* breaks keeper and imports

## [4.16.1](https://github.com/zenrocklabs/zrchain/compare/v4.16.0...v4.16.1) (2024-10-29)


### Bug Fixes

* bump buf version and update protos ([6ecea71](https://github.com/zenrocklabs/zrchain/commit/6ecea7116c17880542ff010841b2e757ee4c45f7))

# [4.16.0](https://github.com/zenrocklabs/zrchain/compare/v4.15.0...v4.16.0) (2024-10-28)


### Bug Fixes

* added zenbtc to upgrade ([6e0ff27](https://github.com/zenrocklabs/zrchain/commit/6e0ff27c6c71702e6b7c051559b23b952c17eeba))


### Features

* approvernumber verification ([1da8a59](https://github.com/zenrocklabs/zrchain/commit/1da8a5936dcb687c0670f7efabcb5305999c8fb5))
* remove abbreviations in policies ([9aa7db1](https://github.com/zenrocklabs/zrchain/commit/9aa7db1aba8b06a9e8743eecf1d4efcc07e1f265))

# [4.15.0](https://github.com/zenrocklabs/zrchain/compare/v4.14.1...v4.15.0) (2024-10-23)


### Features

* ci build pipeline ([#4](https://github.com/zenrocklabs/zrchain/issues/4)) ([f770241](https://github.com/zenrocklabs/zrchain/commit/f770241f3cdedbe0ca1b4357c40979f28d691dd7))
