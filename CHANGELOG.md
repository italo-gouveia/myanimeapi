# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.3.0] - 2025-03-20

### Added
- **Structured Error Handling:** Introduced a new `ErrorResponse` type for consistent error responses across the application.
- **Swagger Documentation Updates:** Enhanced Swagger annotations to reflect the new error response structure and added detailed error response examples for all endpoints.
- **CI/CD Pipeline Improvements:** Added new steps for Swagger validation, security analysis (Semgrep, Gitleaks, Gosec), and Docker image scanning with Trivy.

### Changed
- **Error Handling Mechanism:** Updated all handlers to use the new error handling mechanism, replacing direct `http.Error` calls.
- **Testing Framework:** Added new unit tests for the error handling middleware and updated existing tests to align with the new error response format.

### Fixed
- **Swagger Consistency:** Improved the readability and consistency of the Swagger documentation.

### Impact
- **Developers:** The new error handling mechanism provides a consistent way to handle and return errors, making it easier to debug and maintain the codebase.
- **API Consumers:** The updated Swagger documentation offers clearer insights into the expected error responses, improving the overall API experience.
- **Security:** The enhanced CI/CD pipeline ensures that security vulnerabilities are caught early in the development process.

## [1.2.0] - 2025-03-18

### Added
- **Project Restructuring**:
  - Moved `pkg/` to `api/` to better reflect the purpose of the package as API-related code.
  - Renamed `resources/` to `assets/` for consistency and clarity.
  - Updated references to diagram images in `README.md` to point to the new `assets/` directory.

### Changed
- **Swagger Documentation**:
  - Updated the `swag init` command in the Dockerfile to reflect the new project structure (`./api/handlers` and `./api/models` instead of `./pkg/handlers` and `./pkg/models`).
  - Ensured Swagger documentation is generated correctly during the Docker image build process.
  - Updated the Swagger docs generation path in the CI workflow to match the new structure.

- **Validation and Mocks**:
  - Consolidated validation functions from `pkg/validation` to `pkg/utils` to centralize utility functions.
  - Updated all handlers to use `utils.ValidateID` and `utils.ValidatePagination` instead of the deprecated `validation` package.
  - Relocated `mocks` package from `internal/mocks` to `pkg/mocks` for better organization and accessibility across the project.
  - Renamed `pkg/validation/validation.go` to `pkg/utils/validation.go` to align with the new structure.

### Fixed
- **CI Workflow**:
  - Disabled the `pre-build-validation` job in the CI workflow as it is no longer needed.
  - Ensured Swagger docs are generated correctly and checked for changes in the CI pipeline.

### Refactor
- **Project Layout**:
  - Updated all import paths to reflect the new directory structure.
  - Adjusted import paths in `cmd/main.go` and `tests/smoke/smoke_test.go` to match the new structure.
  - No functional changes were made; this is purely a structural refactor.

### Removed
- **Deprecated Code**:
  - Removed the `pkg/validation` package as its functionality has been consolidated into `pkg/utils`.
  - Removed the `internal/mocks` package as it has been relocated to `pkg/mocks`.

### Issues Closed
- **[GITISSUE-112]**: Restructured project layout and updated paths.
- **[GITISSUE-112]**: Consolidated validation and mocks into `pkg/utils`.


## [1.0.1](https://github.com/italo-gouveia/myanimeapi/compare/v1.0.0...v1.0.1) (2025-03-18)


### Bug Fixes

* GITISSUE-50 Resolve linting and CI workflow issues ([3ae6db2](https://github.com/italo-gouveia/myanimeapi/commit/3ae6db2001418e52b12584e260488909b9c48938))

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
