# 1.0.0 (2025-03-17)


### Bug Fixes

* **ci.yml:** added prevent to errors on the prebuild to not wait until the end of the deployment ([4c96fb0](https://github.com/italo-gouveia/myanimeapi/commit/4c96fb0a7891af30cdc8d63c79b850e60b89eba3))
* **ci.yml:** resolve JSON parsing issue in pre-build-validation job ([5c7ec9b](https://github.com/italo-gouveia/myanimeapi/commit/5c7ec9b7174a738cc7650ebe3b6076ab7068269d))
* **ci.yml:** solving pr description workflow problem ([22246e3](https://github.com/italo-gouveia/myanimeapi/commit/22246e3873bd13e89e01da4f0fae538240ed95e6))
* **ci.yml:** to validate all pr changes and not only the head ([e7c6c97](https://github.com/italo-gouveia/myanimeapi/commit/e7c6c97975f91bcd1b1832d10d1e07f67fccc0eb))


### Features

* [GITISSUE-104] Add test coverage report and Semantic Release automation ([af963bd](https://github.com/italo-gouveia/myanimeapi/commit/af963bd92581412647fdd903b666e6e9cbfb2989)), closes [#104](https://github.com/italo-gouveia/myanimeapi/issues/104)
* [GITISSUE-104] using another pload-artifact version ([c4d5cbf](https://github.com/italo-gouveia/myanimeapi/commit/c4d5cbf2db92b924724c79a6e8ee4c66864dd4c8))
* [GITISSUE-104] using latest upload-artifact version ([00318d6](https://github.com/italo-gouveia/myanimeapi/commit/00318d63bc548aef88bb70c85305520e4eda814e))

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).


## [1.0.0] - 2025-03-17

### Added
- Initial project setup with Go, PostgreSQL, and Docker.
- RESTful API endpoints for managing users, anime, and reviews.
- JWT-based authentication and authorization.
- Swagger documentation for API endpoints.
- Unit tests for handlers using mocks.
- Initial release of MyAnimeAPI.

### Changed
- Improved error handling and validation in API endpoints.
- Updated database schema to include relationships between users, anime, and reviews.

### Fixed
- Fixed pagination logic for reviews and users.
- Resolved issues with JWT token generation and validation.

### Removed
- Removed unused dependencies and code.

---
