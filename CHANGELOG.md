# Changelog

## [1.5.0](https://github.com/phi-ag/mta-sts-exporter/compare/v1.4.1...v1.5.0) (2024-02-18)


### Features

* added enable config for reports and metrics ([8f2b64f](https://github.com/phi-ag/mta-sts-exporter/commit/8f2b64f0aad5c9842e8c8f16fca51c7032393669))
* added some basic prometheus counter ([67d93f6](https://github.com/phi-ag/mta-sts-exporter/commit/67d93f683ffd0134b756d3a8fcd324a0e5a692a9))


### Documentation

* added A and AAAA dns entries to readme and formatted examples ([2d9fdc3](https://github.com/phi-ag/mta-sts-exporter/commit/2d9fdc3290a0c3d491b34e20df67f771b7ef8a0b))
* added traefik example to compose.yaml ([e3a6dcf](https://github.com/phi-ag/mta-sts-exporter/commit/e3a6dcf863b9d63eb2a82f87d1dbf782929d6462))


### Miscellaneous Chores

* **deps:** pin phiag/mta-sts-exporter docker tag to 874a452 ([7dd2d14](https://github.com/phi-ag/mta-sts-exporter/commit/7dd2d14420e8276ff8fcdf67f1c423c0f1ff80bf))


### Code Refactoring

* disabled save report by default ([e7348f9](https://github.com/phi-ag/mta-sts-exporter/commit/e7348f9d77030a14fc1708f59fd8ea3b66c92f7d))
* moved max body/json config ([e2905cd](https://github.com/phi-ag/mta-sts-exporter/commit/e2905cdae3e56fc5b09b09b4ac41d342eb76a4e9))
* moved save config ([aa9596d](https://github.com/phi-ag/mta-sts-exporter/commit/aa9596d59a44cb14b14cec18b021f53e5a1558b7))
* reordered functions in main ([555dba4](https://github.com/phi-ag/mta-sts-exporter/commit/555dba413de79406cbbdfd598e870a905f5bf999))

## [1.4.1](https://github.com/phi-ag/mta-sts-exporter/compare/v1.4.0...v1.4.1) (2024-02-18)


### Performance Improvements

* generate policy response on startup ([5782223](https://github.com/phi-ag/mta-sts-exporter/commit/5782223ac80cad0510e38ac85b95b72b8f474bb5))


### Documentation

* added dns entries to readme ([15fd78d](https://github.com/phi-ag/mta-sts-exporter/commit/15fd78da636d3d282eca06d84740a4c54e1e0aeb))
* added environment based config example to compose.yaml ([e7c0ef5](https://github.com/phi-ag/mta-sts-exporter/commit/e7c0ef55030fb65e79c6ddb16318e920e44c6b48))
* configure policy in compose.yaml ([5d50c4a](https://github.com/phi-ag/mta-sts-exporter/commit/5d50c4a3a84b87d9ecf2fa0d83c0d584c5f6a575))


### Miscellaneous Chores

* **deps:** pin phiag/mta-sts-exporter docker tag to bc1b5c3 ([5557dc8](https://github.com/phi-ag/mta-sts-exporter/commit/5557dc8fa87157442b81ebce12064b58cf5ab562))


### Code Refactoring

* changed default report path to /report ([858fc87](https://github.com/phi-ag/mta-sts-exporter/commit/858fc87f146024fa58327c74fa1e86a04f7a73e7))
* **config:** replaced policy content config with fields ([4dc9995](https://github.com/phi-ag/mta-sts-exporter/commit/4dc9995ae6f729c9fc68ebbdc94bdcb21dafaeb1))

## [1.4.0](https://github.com/phi-ag/mta-sts-exporter/compare/v1.3.0...v1.4.0) (2024-02-18)


### Features

* added option to serve mta-sts policy ([2c59c63](https://github.com/phi-ag/mta-sts-exporter/commit/2c59c632d3297b609b809460cfd989dbb41a17ab))


### Documentation

* fixed docker compose healthcheck ([cb2facc](https://github.com/phi-ag/mta-sts-exporter/commit/cb2facc4031039c1bdece7fd3bff24a76aed1416))


### Miscellaneous Chores

* **deps:** pin phiag/mta-sts-exporter docker tag to adc373c ([2c38da4](https://github.com/phi-ag/mta-sts-exporter/commit/2c38da4092dd22a0c85a724c4450ebfe2a2be3a1))

## [1.3.0](https://github.com/phi-ag/mta-sts-exporter/compare/v1.2.1...v1.3.0) (2024-02-18)


### Features

* added healthcheck ([b37a9ac](https://github.com/phi-ag/mta-sts-exporter/commit/b37a9ac16e7c623e9f2bca63d27b9b944ac1025a))


### Bug Fixes

* panic when unmarshal config fails ([9e87e83](https://github.com/phi-ag/mta-sts-exporter/commit/9e87e83fc5be5853c29592e42fd3d84535d6e380))


### Documentation

* added links to RFCs and comment for mx-host array errata ([90e4774](https://github.com/phi-ag/mta-sts-exporter/commit/90e4774fdfaec92fd379a2ae40894b036b3b4362))


### Code Refactoring

* renamed body and json reader ([9f0ef77](https://github.com/phi-ag/mta-sts-exporter/commit/9f0ef770a16ecf6cb483ba0f09121b76b10e2d6f))

## [1.2.1](https://github.com/phi-ag/mta-sts-exporter/compare/v1.2.0...v1.2.1) (2024-02-18)


### Bug Fixes

* parse mx-host field (string or string array) ([7342162](https://github.com/phi-ag/mta-sts-exporter/commit/73421625bcdff4387d05660b86b65117d9129443))


### Documentation

* fixed compose.yaml config ([37db97e](https://github.com/phi-ag/mta-sts-exporter/commit/37db97e784c2a250f5d87050f56ad07fbf2dffe6))
* fixed uid/gid in readme ([9e08f39](https://github.com/phi-ag/mta-sts-exporter/commit/9e08f39d81dcacb7f921c0468378c0a61a08fd15))


### Miscellaneous Chores

* **deps:** pin phiag/mta-sts-exporter docker tag to ac65f14 ([e797067](https://github.com/phi-ag/mta-sts-exporter/commit/e797067f49a4354b3a0511af715bac32f0f82088))


### Code Refactoring

* split tests ([4236cba](https://github.com/phi-ag/mta-sts-exporter/commit/4236cbac16a70fa3f1582c6d9020795a389b1bab))


### Tests

* added google and microsoft report examples ([61458ad](https://github.com/phi-ag/mta-sts-exporter/commit/61458ad1fd2add41b8cd3b93f8ba5000046ee3b0))


### Continuous Integration

* added note about missing platforms ([160a55e](https://github.com/phi-ag/mta-sts-exporter/commit/160a55e859c0112571b4a4cb7a9443b15ac66340))

## [1.2.0](https://github.com/phi-ag/mta-sts-exporter/compare/v1.1.2...v1.2.0) (2024-02-17)


### Features

* enabled platforms riscv64 and 386, dropped ppc64le and s390x ([4ad7769](https://github.com/phi-ag/mta-sts-exporter/commit/4ad7769105db2ba62228dab633ccd47dd951155a))


### Miscellaneous Chores

* **deps:** update golangci/golangci-lint-action action to v4 ([7402705](https://github.com/phi-ag/mta-sts-exporter/commit/740270580b30a9b2ca65336472c2118f470fe811))
* **deps:** update phiag/mta-sts-exporter docker tag to v1.1.2 ([d73dd5c](https://github.com/phi-ag/mta-sts-exporter/commit/d73dd5c23472fc1e919afa7219ba3f8574b48ce3))

## [1.1.2](https://github.com/phi-ag/mta-sts-exporter/compare/v1.1.1...v1.1.2) (2024-02-17)


### Documentation

* added example docker compose.yaml and infer config file type ([bee678a](https://github.com/phi-ag/mta-sts-exporter/commit/bee678a4e1df67c3e3f3c85dba3ba861c47c0c04))
* added readme badges ([12d5801](https://github.com/phi-ag/mta-sts-exporter/commit/12d5801d694462accc8732cf0bc8d7808f7a3551))
* added save report instructions to readme ([3f993ef](https://github.com/phi-ag/mta-sts-exporter/commit/3f993efb58a0e4e78e0266e9e5a9470570663a1e))


### Miscellaneous Chores

* **deps:** pin golangci/golangci-lint-action action to 3a91952 ([b1c7ebf](https://github.com/phi-ag/mta-sts-exporter/commit/b1c7ebf7818c112bb60098300f10c1a767f7af44))
* **deps:** pin phiag/mta-sts-exporter docker tag to fc97b23 ([92e07ee](https://github.com/phi-ag/mta-sts-exporter/commit/92e07eef19742812c53315d94d3e72eca5648570))


### Code Refactoring

* split into multiple files ([feb21b9](https://github.com/phi-ag/mta-sts-exporter/commit/feb21b91b7f2702fcb2b000acb761132e3de06e4))


### Continuous Integration

* added lint ([005add3](https://github.com/phi-ag/mta-sts-exporter/commit/005add3800b873496888e8f9b6cfc94762a23217))

## [1.1.1](https://github.com/phi-ag/mta-sts-exporter/compare/v1.1.0...v1.1.1) (2024-02-17)


### Bug Fixes

* fixed save report and use viper for config ([8a043cb](https://github.com/phi-ag/mta-sts-exporter/commit/8a043cb0af5db2ccd8c5396b62e01a57b8646ada))
* use uint16 for port config ([ea55a19](https://github.com/phi-ag/mta-sts-exporter/commit/ea55a19ee96f494555c592574143c8c76b3e4017))


### Styles

* removed some whitespace ([999cdcf](https://github.com/phi-ag/mta-sts-exporter/commit/999cdcfea94ca8d172432e827a5e47d01bb9bf48))

## [1.1.0](https://github.com/phi-ag/mta-sts-exporter/compare/v1.0.0...v1.1.0) (2024-02-16)


### Features

* configurable json logging ([dcdcb6e](https://github.com/phi-ag/mta-sts-exporter/commit/dcdcb6e53855fe893c862b0eeb432cb88a604e27))
* configurable report path ([1a59e21](https://github.com/phi-ag/mta-sts-exporter/commit/1a59e21d199d6eb804919e85a2962e7e68716c91))


### Bug Fixes

* added image labels ([f17d5e0](https://github.com/phi-ag/mta-sts-exporter/commit/f17d5e086a02f5edd76d7dbae8d64697417ec0ab))
* fix release action ([ee174b1](https://github.com/phi-ag/mta-sts-exporter/commit/ee174b140b0b5112109324b8ea00018eb2251b12))
* fix release action ([4796118](https://github.com/phi-ag/mta-sts-exporter/commit/479611813c7cfe696442e6582d0bd6948e157bc3))
* login to dockerhub ([78692de](https://github.com/phi-ag/mta-sts-exporter/commit/78692deb28facdc196dd691f78013b7f17322b64))
* shared setup go action ([4f3dcc3](https://github.com/phi-ag/mta-sts-exporter/commit/4f3dcc34daf335a602c9e39addeb5e005aacfd5b))


### Miscellaneous Chores

* **deps:** pin docker/login-action action to 343f7c4 ([5a26d3f](https://github.com/phi-ag/mta-sts-exporter/commit/5a26d3fb0eb9c009ee312c509cd1890f4bf5702d))
* **main:** release 1.0.1 ([9196b54](https://github.com/phi-ag/mta-sts-exporter/commit/9196b54b524869301aa0088bfd8848d168259db4))
* **main:** release 1.0.2 ([43c840f](https://github.com/phi-ag/mta-sts-exporter/commit/43c840fb13a9b1ebb9ff21b337639bd01ca8eebc))
* **main:** release 1.0.3 ([b1fa859](https://github.com/phi-ag/mta-sts-exporter/commit/b1fa8593b139b9d2080aee21a280e2752f7d1739))
* **main:** release 1.0.4 ([2b7849e](https://github.com/phi-ag/mta-sts-exporter/commit/2b7849e9f9011cdbdb6c9cf6bfe3a73dca944e74))
* **main:** release 1.0.5 ([58db23f](https://github.com/phi-ag/mta-sts-exporter/commit/58db23f2145523825e320ae75a813653167f8add))


### Code Refactoring

* **config:** refactored config ([1a52b73](https://github.com/phi-ag/mta-sts-exporter/commit/1a52b73002b8dc93e1cc732a06f554bf1e176bc3))
* save to file ([9b66f01](https://github.com/phi-ag/mta-sts-exporter/commit/9b66f0111031300c000b487984967ccd85bb51ed))


### Continuous Integration

* added release-please config and fixed release type ([d4391a0](https://github.com/phi-ag/mta-sts-exporter/commit/d4391a0542d1c2411477ceb4e695a4bbb2063507))

## [1.0.5](https://github.com/phi-ag/mta-sts-exporter/compare/v1.0.4...v1.0.5) (2024-02-16)


### Bug Fixes

* added image labels ([f17d5e0](https://github.com/phi-ag/mta-sts-exporter/commit/f17d5e086a02f5edd76d7dbae8d64697417ec0ab))

## [1.0.4](https://github.com/phi-ag/mta-sts-exporter/compare/v1.0.3...v1.0.4) (2024-02-16)


### Bug Fixes

* shared setup go action ([4f3dcc3](https://github.com/phi-ag/mta-sts-exporter/commit/4f3dcc34daf335a602c9e39addeb5e005aacfd5b))

## [1.0.3](https://github.com/phi-ag/mta-sts-exporter/compare/v1.0.2...v1.0.3) (2024-02-16)


### Bug Fixes

* login to dockerhub ([78692de](https://github.com/phi-ag/mta-sts-exporter/commit/78692deb28facdc196dd691f78013b7f17322b64))

## [1.0.2](https://github.com/phi-ag/mta-sts-exporter/compare/v1.0.1...v1.0.2) (2024-02-16)


### Bug Fixes

* fix release action ([ee174b1](https://github.com/phi-ag/mta-sts-exporter/commit/ee174b140b0b5112109324b8ea00018eb2251b12))

## [1.0.1](https://github.com/phi-ag/mta-sts-exporter/compare/v1.0.0...v1.0.1) (2024-02-16)


### Bug Fixes

* fix release action ([4796118](https://github.com/phi-ag/mta-sts-exporter/commit/479611813c7cfe696442e6582d0bd6948e157bc3))

## 1.0.0 (2024-02-16)


### Features

* init ([b1a3bd7](https://github.com/phi-ag/mta-sts-exporter/commit/b1a3bd7586dc537073acfe69a50e0c556ce36db2))
