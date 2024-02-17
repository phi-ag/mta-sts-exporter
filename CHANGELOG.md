# Changelog

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

## Changelog
