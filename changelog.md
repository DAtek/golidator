# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.1.2] - 2022-11-27

### Changed
- `GetValidatorsForObject` accepts context
- `GetValidatorsForList` accepts context

## [v0.1.1] - 2022-11-19

### Fixed
- `Validate()` returned `ValidationError` without `source` field being set
- `GetValidators()` was private method, therefore `Validate()` couldn't be used externally


## [v0.1.0] - 2022-11-19
