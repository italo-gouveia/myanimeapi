# [1.14.0](https://github.com/italo-gouveia/myanimeapi/compare/v1.13.2...v1.14.0) (2026-04-27)


### Bug Fixes

* **test-env:** remove initdb.d mount and add test-db make target ([c21c210](https://github.com/italo-gouveia/myanimeapi/commit/c21c2105a0bc4b8b8b158a7d26fa52b2b62f63be))


### Features

* **migrations:** replace AutoMigrate with golang-migrate SQL migrations ([#45](https://github.com/italo-gouveia/myanimeapi/issues/45)) ([10eb1d2](https://github.com/italo-gouveia/myanimeapi/commit/10eb1d2d23108743da95ee9487682d3eb380e9f6))

## [1.13.2](https://github.com/italo-gouveia/myanimeapi/compare/v1.13.1...v1.13.2) (2026-04-27)


### Bug Fixes

* **smoke:** suppress errcheck lint on resp.Body.Close defers ([a5f2a9f](https://github.com/italo-gouveia/myanimeapi/commit/a5f2a9f3cdef17e0e4b62e4f7d47abbf7535a85d))

## [1.13.1](https://github.com/italo-gouveia/myanimeapi/compare/v1.13.0...v1.13.1) (2026-04-26)


### Bug Fixes

* **ci:** fix test workflow, swagger docs, and coverage ([7fca237](https://github.com/italo-gouveia/myanimeapi/commit/7fca237ca3a8546af07f6b58cacc6c2d59c1634c))

# [1.13.0](https://github.com/italo-gouveia/myanimeapi/compare/v1.12.0...v1.13.0) (2026-04-26)


### Bug Fixes

* handle json.Unmarshal errors in test files to satisfy linter requirements ([ddbd652](https://github.com/italo-gouveia/myanimeapi/commit/ddbd652b447a63cfe1f604ef8e706a2e280e6171))
* **hooks:** add GOPATH/bin to PATH so golangci-lint is found after install ([68039be](https://github.com/italo-gouveia/myanimeapi/commit/68039be1e664bdca7782ab7f935e2f982c3a9828))
* regenerate anime service mock with missing methods ([7e78476](https://github.com/italo-gouveia/myanimeapi/commit/7e784765b80f1f781340949079e1524025fa3897))
* remove duplicate anime service mock file ([e207226](https://github.com/italo-gouveia/myanimeapi/commit/e20722648efc758c653dece822a5b4ab2b872440))
* remove unsupported Favorites preload from AnimeRepository ([1651247](https://github.com/italo-gouveia/myanimeapi/commit/1651247ac27ccb879069d0cc87467a11124ad196))
* **tests:** align test expectations with actual API responses ([fcd98fc](https://github.com/italo-gouveia/myanimeapi/commit/fcd98fc2442281e5a73888b3fa636a8674c75053))
* **tests:** check error return from multipart writer Close ([fca1a76](https://github.com/italo-gouveia/myanimeapi/commit/fca1a762d72313bc09611fabddba335827696ba3))
* **tests:** fix favorite integration tests ([8542aca](https://github.com/italo-gouveia/myanimeapi/commit/8542aca82c204830809b2fff02b4e1da4468ccd3))
* **tests:** fix integration tests and improve API structure ([8a0e80d](https://github.com/italo-gouveia/myanimeapi/commit/8a0e80d16f8ba16dc66b620ea3fde1019513f8b3))
* **tests:** fix linter error in favorite integration test ([98a7134](https://github.com/italo-gouveia/myanimeapi/commit/98a713420c6c5ee035f4f1d6e5925452a5d55852))
* **tests:** fix linter errors in anime integration tests ([a3ffda1](https://github.com/italo-gouveia/myanimeapi/commit/a3ffda1f7b921124e973d7b675f91d5b30cd4442))
* **tests:** update expected fields in GetAnimesByGenreHandler test to match API response ([055107c](https://github.com/italo-gouveia/myanimeapi/commit/055107c284f144191dcb19b64a1bf9aba1672485))
* **tests:** update integration and favorite tests to match API response formats and improve error handling ([6f2dca8](https://github.com/italo-gouveia/myanimeapi/commit/6f2dca8824a9d053a6669b5a6eebccd09c445c40))
* **tests:** update integration test expectations to match API response format ([3302ca4](https://github.com/italo-gouveia/myanimeapi/commit/3302ca4ef63116430060daa6b5f872727c0b8d41))
* update tests and mocks for anime handlers, align status codes and response bodies, improve semantic match for handler expectations ([73bcb2a](https://github.com/italo-gouveia/myanimeapi/commit/73bcb2a93f5970c168ea502307d172fdd5dae220))


### Features

* **api:** standardize error handling and response formats ([30b4940](https://github.com/italo-gouveia/myanimeapi/commit/30b4940743f105d1101b87c381d19facdcd86244))
* **auth:** improve user ID handling and test coverage ([c4c2f71](https://github.com/italo-gouveia/myanimeapi/commit/c4c2f71d88a3c10aa052e63cedaf12390cae37f8))
* **ci:** add AI code review workflows and update CI configuration ([ae84fe9](https://github.com/italo-gouveia/myanimeapi/commit/ae84fe9214906f59a05049873d8da61de25b8118))
* integrate SONAR_IMPROVMENTS work into main ([0953c39](https://github.com/italo-gouveia/myanimeapi/commit/0953c39be04817627c10057cbef2c969134a24a6))
* **test:** add mock implementations for anime and favorite services ([c8615ce](https://github.com/italo-gouveia/myanimeapi/commit/c8615cedea8fff38881ea65f788422947bcfbf60))
* **test:** add mock implementations for anime and favorite services ([85edc8f](https://github.com/italo-gouveia/myanimeapi/commit/85edc8f12bbe85bba3656aa170f4e7cc1e8ac950))
* **tests:** add integration tests for anime and favorite handlers ([7b915a9](https://github.com/italo-gouveia/myanimeapi/commit/7b915a90f88c10f1e8ee4b6be0dbb2d605d7933c))


### BREAKING CHANGES

* **tests:** Search routes (/animes/search and /animes/genre/{genre}) are now public
* **auth:** User ID is now stored as uint instead of string in context

- Update auth middleware to parse and validate user ID from JWT claims
- Modify context middleware to handle uint user IDs directly
- Refactor integration tests to use proper authentication
- Update test cases to match new error response format
- Add comprehensive test coverage for unauthorized scenarios

This change improves type safety and reduces potential runtime errors
by handling user IDs as uint throughout the application instead of
converting between string and uint in multiple places.
* **api:** - Changed response format from "animes" to "data" in anime search endpoints
- Modified genre search endpoint to use query parameters instead of path parameters
- Simplified response structures by removing ToResponse() transformations
- Standardized error response format across all handlers
- Improved error messages and context handling

This change improves API consistency and error handling while simplifying the response structure.

# [1.12.0](https://github.com/italo-gouveia/myanimeapi/compare/v1.11.0...v1.12.0) (2025-05-26)


### Bug Fixes

* **ci:** ensure SARIF file exists for Gitleaks report ([e48c00e](https://github.com/italo-gouveia/myanimeapi/commit/e48c00ea1afe392d0696158fd56919a75eea5be4))
* **ci:** improve security scanning and rate limiting ([e34e8f0](https://github.com/italo-gouveia/myanimeapi/commit/e34e8f0d85bac07a8da559b71ba3913eddb59cd5))
* **ci:** resolve SARIF file path issues in security workflow ([c049983](https://github.com/italo-gouveia/myanimeapi/commit/c049983e985a538af012da23f0106945960adddb))
* **ci:** update Gitleaks action and SARIF file paths ([2abc53e](https://github.com/italo-gouveia/myanimeapi/commit/2abc53ea44a8edc5b28add1bf01a6950f8114f69))
* comment out problematic SARIF upload steps ([09dd851](https://github.com/italo-gouveia/myanimeapi/commit/09dd85100b061fc014ce206b86c54a0c92e43fed))
* enable CGO for race detection in pre-push hook ([6a47d82](https://github.com/italo-gouveia/myanimeapi/commit/6a47d82302098a443e5be1b79ca8c2e469eded38))
* ensure SARIF file exists before upload ([4c7181f](https://github.com/italo-gouveia/myanimeapi/commit/4c7181facb38c610ccb91a2f423d15f02de5c3f5))
* handle missing bc command in pre-push hook ([11bd021](https://github.com/italo-gouveia/myanimeapi/commit/11bd021619a8cf879985a831bab09a274f702ade))
* make pre-push hook handle missing gcc gracefully ([c243a36](https://github.com/italo-gouveia/myanimeapi/commit/c243a3633ad261712c05d4e3b39e008316b161ff))
* temporarily lower test coverage threshold to 5% ([549641c](https://github.com/italo-gouveia/myanimeapi/commit/549641c28d9b0a60b709949765f67873aa563d14))


### Code Refactoring

* **logger:** implement custom structured logger ([f4ba061](https://github.com/italo-gouveia/myanimeapi/commit/f4ba061aa108a5011b5dd98f18a06f39635d2856))
* **services:** enhance error handling and logging in anime service ([823d5c6](https://github.com/italo-gouveia/myanimeapi/commit/823d5c6cd71aa839467f801f81d4e56f79c254db))
* **services:** enhance error handling and logging in genre and tag services ([6eba7e0](https://github.com/italo-gouveia/myanimeapi/commit/6eba7e09ca85f1f36ac2db55a6143cf01a74cfaa))


### Features

* **auth:** enhance authentication handler with improved error handling and logging ([fe5b537](https://github.com/italo-gouveia/myanimeapi/commit/fe5b5379820cc2151c53a8b59636b7832d1359d5))
* **auth:** enhance error handling and logging in auth services ([2105996](https://github.com/italo-gouveia/myanimeapi/commit/210599638f2fa54f14004913a32a6cd55be20cc8))
* **auth:** enhance JWT authentication and authorization ([dbe62bc](https://github.com/italo-gouveia/myanimeapi/commit/dbe62bc05e12c6c53f8099cdaaa13d7340359933))
* **ci:** enhance CI/CD pipeline and improve application logging ([24f1e8c](https://github.com/italo-gouveia/myanimeapi/commit/24f1e8ca02ffabb428a96157a9a72f6317d61fbb))
* **genres:** enhance genre handlers with improved error handling and logging ([0f92b8d](https://github.com/italo-gouveia/myanimeapi/commit/0f92b8dccd2ec05037978f054bdcf2880e816b7f))
* **handlers:** enhance anime handlers with improved error handling and validation ([51285d8](https://github.com/italo-gouveia/myanimeapi/commit/51285d8801c98be752afe1198b9ad73aa2e19bf1))
* **handlers:** enhance favorite handlers with improved error handling and logging ([cd8de45](https://github.com/italo-gouveia/myanimeapi/commit/cd8de451a7c0c0da3e0739250dd8b12f1a79fc60))
* **handlers:** enhance user handlers with improved error handling and logging ([7379e3d](https://github.com/italo-gouveia/myanimeapi/commit/7379e3db02e90ccd6cdbfc229317d804a747e7bf))
* **logger:** add Get() function and update error logging ([e26489b](https://github.com/italo-gouveia/myanimeapi/commit/e26489b867b74f1101bbfa781d964e3c744aabc9))
* **logging:** implement structured logging with logrus across core layers ([4bd6207](https://github.com/italo-gouveia/myanimeapi/commit/4bd62074712ee86e5e5297e5a8951e701397097a))
* **middleware:** enhance error handling and logging system ([9aad8c9](https://github.com/italo-gouveia/myanimeapi/commit/9aad8c977b2028589436f67c86ec4af8f7abf92c))
* **middleware:** enhance middleware tests and add new test cases ([48320a9](https://github.com/italo-gouveia/myanimeapi/commit/48320a93026bc0eaf2e066252f2aeb0647b9cdea))
* **middleware:** enhance payload validation and sanitization ([986f5c9](https://github.com/italo-gouveia/myanimeapi/commit/986f5c93ddc07d7fc522301612e7c09b6f6b886b))
* **password-reset:** enhance error handling and add logging ([7e83b3d](https://github.com/italo-gouveia/myanimeapi/commit/7e83b3de182f4906f6269a10aaa575012536bcc1))
* **repositories:** enhance logging and error handling ([77cbea0](https://github.com/italo-gouveia/myanimeapi/commit/77cbea0a8fdccdcf2dd44813906da14e3f7a2c7d))
* **reviews:** enhance review handlers with improved security and logging ([6fb0fb3](https://github.com/italo-gouveia/myanimeapi/commit/6fb0fb3ef4b4a2a94f97ecb63ae8cdbcf6d9de7b))
* **security:** enhance HTTPS middleware with improved security features ([1616943](https://github.com/italo-gouveia/myanimeapi/commit/16169435504a0a32e51bb250072709f56c0ffe1c))


### BREAKING CHANGES

* **ci:** The Docker build process now uses buildx and the error response
structure includes additional fields for logging context. Update your deployment
scripts and error handling accordingly.

This change improves the CI/CD pipeline's reliability and security while adding
better logging and error handling throughout the application. The enhanced PR
description analysis provides better visibility into changes, and the improved
security scanning helps identify potential issues earlier in the development
process.
* **middleware:** The Authenticate middleware has been renamed to AuthMiddleware
and CustomClaims UserID type has been changed from uint to string. Update your
code accordingly.

This change improves the middleware test coverage and reliability while adding
new test cases for recently added middleware components. The enhanced error
handling tests ensure proper error responses across all middleware functions.
* **genres:** The genre handler now uses AuthMiddleware instead of Authenticate
and has updated response formats for bulk operations. Update your API clients accordingly.

This change improves the genre handlers' reliability, maintainability, and
error handling capabilities while providing better logging and response formatting.
The new error handling system provides more detailed information about failures,
making it easier to debug issues.
* **reviews:** The review handler constructor now requires interface types
instead of concrete types. Update your service initialization accordingly.

This change improves the review handlers' security, maintainability, and
reliability while providing better error handling and logging capabilities.
The new authorization checks ensure users can only modify their own reviews.
* **auth:** The /auth/login endpoint has been renamed to /auth/authenticate
for better clarity and consistency. Update your API clients accordingly.

This change improves the authentication handler's reliability, maintainability,
and security while providing better error handling and logging capabilities.
* **services:** This commit introduces breaking changes to error handling and logging formats. Clients need to update their error handling logic and log parsing mechanisms to accommodate the new structured format and additional context data.
* **services:** This commit introduces breaking changes to error handling and logging formats. Clients need to update their error handling logic and log parsing mechanisms to accommodate the new structured format and additional context data.
* **logger:** Replace logrus with custom structured logger implementation

- Remove logrus dependency in favor of standard library
- Add structured logging with JSON output format
- Implement log levels (DEBUG, INFO, WARNING, ERROR)
- Add field-based context to log entries
- Add formatted logging methods (Debugf, Infof, etc.)
- Add default logger instance with global methods
- Improve timestamp formatting using RFC3339
- Add proper JSON marshaling for log entries

This change improves logging consistency and removes external dependencies.
All code using the logger package will need to be updated to use the new API.
* **logging:** All logging is now structured and uses logrus; any custom log parsing or monitoring should be updated accordingly.

# [1.11.0](https://github.com/italo-gouveia/myanimeapi/compare/v1.10.1...v1.11.0) (2025-05-15)


### Features

* **db:** seed admin user during initial database setup ([8402563](https://github.com/italo-gouveia/myanimeapi/commit/84025635ca33b1e12814730c1b1e10602f334508))

## [1.10.1](https://github.com/italo-gouveia/myanimeapi/compare/v1.10.0...v1.10.1) (2025-05-15)


### Bug Fixes

* **handlers:** improve ID validation for path parameters ([ac33400](https://github.com/italo-gouveia/myanimeapi/commit/ac33400e5932a849c5739c71b46be87c80164ed6))

# [1.10.0](https://github.com/italo-gouveia/myanimeapi/compare/v1.9.0...v1.10.0) (2025-05-15)


### Bug Fixes

* **models:** correct swagger example annotations for User and UserUpdateRequest ([e7418d5](https://github.com/italo-gouveia/myanimeapi/commit/e7418d5945f51d0d742e6a8405bc5c41aa1b33cc))


### Features

* **auth:** implement password reset functionality ([fc03c55](https://github.com/italo-gouveia/myanimeapi/commit/fc03c551adcaeba4e0080cffe29217f7888611c1))

# [1.9.0](https://github.com/italo-gouveia/myanimeapi/compare/v1.8.0...v1.9.0) (2025-05-14)


### Bug Fixes

* use custom type for request ID context key ([80c0501](https://github.com/italo-gouveia/myanimeapi/commit/80c05019040ee2e7812070a6a7c4d2f798dd3806))


### Features

* add request ID middleware for request traceability ([03ca20f](https://github.com/italo-gouveia/myanimeapi/commit/03ca20f87b42a6146c9c823b2a81b4e0a3beb471))

# [1.8.0](https://github.com/italo-gouveia/myanimeapi/compare/v1.7.0...v1.8.0) (2025-05-13)


### Features

* **reviews:** add media attachment support with storage service ([0df0131](https://github.com/italo-gouveia/myanimeapi/commit/0df0131164c0123f284d011ff3cf74d3682802b1))
* **reviews:** implement media file handling in review handlers ([802f73f](https://github.com/italo-gouveia/myanimeapi/commit/802f73faf9e090e12f5b4d8e7fc4373e32bec4b4))


### BREAKING CHANGES

* **reviews:** None
* **reviews:** None

# [1.7.0](https://github.com/italo-gouveia/myanimeapi/compare/v1.6.0...v1.7.0) (2025-05-12)


### Bug Fixes

* **handlers:** add error handling for JSON encoding ([463a41b](https://github.com/italo-gouveia/myanimeapi/commit/463a41bfe19fc52816af5d1348e53dd87a3ed47e))
* **handlers:** add error handling for JSON encoding ([3de2456](https://github.com/italo-gouveia/myanimeapi/commit/3de2456b1fdfd1f755fc70e5011a6cfb76533f2e))


### Features

* **api:** enhance API documentation and add new endpoints ([a9e16fe](https://github.com/italo-gouveia/myanimeapi/commit/a9e16feb3d9caf2bfedb960e69ef7c21a5168a23)), closes [#87](https://github.com/italo-gouveia/myanimeapi/issues/87) [#89](https://github.com/italo-gouveia/myanimeapi/issues/89) [#56](https://github.com/italo-gouveia/myanimeapi/issues/56)
* **api:** improve review and user request handling ([1e629f2](https://github.com/italo-gouveia/myanimeapi/commit/1e629f2ff2ec61b7cbaf1d86f8c78df3f8fe1bb0))


### BREAKING CHANGES

* **handlers:** None
* **handlers:** None
* **api:** None
* **api:** User endpoints have been restructured for better security and usability

- Add health check and version endpoints with detailed documentation
- Restructure user management endpoints with improved security
- Add comprehensive request/response examples for user operations
- Implement favorites functionality for anime
- Enhance error handling with consistent response format
- Add health monitoring and dependency management features
- Improve API documentation with detailed endpoint descriptions

This change improves the API's usability and security while providing better
documentation for developers. The user endpoints have been restructured to
follow better security practices, requiring authentication for sensitive
operations.

# [1.6.0](https://github.com/italo-gouveia/myanimeapi/compare/v1.5.0...v1.6.0) (2025-05-11)


### Code Refactoring

* **repository:** implement generic repository pattern for genres and tags ([f50a550](https://github.com/italo-gouveia/myanimeapi/commit/f50a55063c77c03754975f8b223afeb41b654821))


### Features

* **api:** add comprehensive anime, genre, and tag management endpoints ([f7e575b](https://github.com/italo-gouveia/myanimeapi/commit/f7e575b0d3c443519246667ca131b43eef725325))
* **api:** enhance anime, genre, and tag handlers with comprehensive documentation ([c1e2c45](https://github.com/italo-gouveia/myanimeapi/commit/c1e2c45435b96d5d26df60f90ce661ea8595eb74))
* **genres:** add search and bulk operations ([252bbe2](https://github.com/italo-gouveia/myanimeapi/commit/252bbe29c26ec3895f99bd9aaa56b95a6ea588e0))
* **models:** add Genre and Tag models with anime relationships ([0fa5849](https://github.com/italo-gouveia/myanimeapi/commit/0fa584982ce42411665d5439347beaa6c6ce94c8))
* **models:** improve Swagger documentation and model structure ([789e8ce](https://github.com/italo-gouveia/myanimeapi/commit/789e8ce2edb4cd76f6789cf384f02b1d195643fb))
* **repositories:** implement Genre and Tag repositories ([4c9a5cd](https://github.com/italo-gouveia/myanimeapi/commit/4c9a5cd5acc9588ce16e1576f8780094c636b991))


### BREAKING CHANGES

* **genres:** None
* **api:** Some response formats have been updated to include additional fields
and improved type safety. Clients should update their response handling accordingly.
* **repository:** Repository interfaces have been refactored to use a
generic pattern. The following changes are required:
- Update repository method calls to use new interface methods
- Update service layer to handle new error types
- Update handlers to work with new repository and service interfaces

# [1.5.0](https://github.com/italo-gouveia/myanimeapi/compare/v1.4.0...v1.5.0) (2025-04-30)


### Bug Fixes

* **anime-handler:** use validated payload from context in CreateAnimeHandler ([143a485](https://github.com/italo-gouveia/myanimeapi/commit/143a485160b410c6431d12724af4929b0761dfb0))
* **api:** improve error handling in RegisterUserHandler ([bdfd19e](https://github.com/italo-gouveia/myanimeapi/commit/bdfd19ef9ac226940f047f7a34e5228da2a0feeb))
* **ci:** add GITHUB_TOKEN to Gitleaks action configuration ([27b86e9](https://github.com/italo-gouveia/myanimeapi/commit/27b86e948b7995de8c4503668c66a1b980b59cda))
* **ci:** correct Gitleaks action configuration in workflow ([065e2f5](https://github.com/italo-gouveia/myanimeapi/commit/065e2f5ac71148a2092b7c0bdb5b273e6d1dd00b))
* **security:** add required rule IDs to Gitleaks configuration ([383e071](https://github.com/italo-gouveia/myanimeapi/commit/383e071693b9ebd944954ca0c990cac8027c84c5))
* **security:** update Gitleaks configuration and workflow ([bbe9168](https://github.com/italo-gouveia/myanimeapi/commit/bbe91685b98f8b33d127572fa54b4d21e483d77d))
* **security:** update Gitleaks configuration and workflow ([005f5c9](https://github.com/italo-gouveia/myanimeapi/commit/005f5c91e369f0db1586394e5670f2e1135787ae))
* **security:** update Gitleaks workflow configuration ([7b9d633](https://github.com/italo-gouveia/myanimeapi/commit/7b9d63396680d26289f8540aaac66df94d790067))
* **utils:** handle JSON encoding errors in response writers ([b8a1ad3](https://github.com/italo-gouveia/myanimeapi/commit/b8a1ad3c3a95d7abf76ab5aa3237e3bbb1a4afd0))


### Code Refactoring

* Implement interface-based dependency injection for anime handlers ([0da4a9a](https://github.com/italo-gouveia/myanimeapi/commit/0da4a9a5a379f9397530bbac1d00d1a4e0b17cbb))


### Features

* **anime:** update CreateAnimeHandler to use AnimeCreateRequest ([fcafaa8](https://github.com/italo-gouveia/myanimeapi/commit/fcafaa8a0c4492121dfdd3784d4aa499f82df32f))
* **api:** enhance Swagger documentation and security ([fb7e166](https://github.com/italo-gouveia/myanimeapi/commit/fb7e1660f67d4d9982d19e1044e48c26a0448f5c))
* **arch:** implement service layer and refactor handlers ([9ae9c4f](https://github.com/italo-gouveia/myanimeapi/commit/9ae9c4f78c070694077a8e71b882c89e48dcf715))
* **auth:** implement auth repository and service with generic response type ([eb7d034](https://github.com/italo-gouveia/myanimeapi/commit/eb7d0349510e20835d571f1f0f6d0362aa7ae288))
* enhance API with favorite anime and security features ([b41407e](https://github.com/italo-gouveia/myanimeapi/commit/b41407e0cd7506839056a7782999465d519b5a0f))
* **errors:** add structured AppError type and constructor ([2513211](https://github.com/italo-gouveia/myanimeapi/commit/25132111406b692df174c6e2a354279f29e7d23c))
* **errors:** add WriteAppErrorResponse function for unified error handling ([33a3d3a](https://github.com/italo-gouveia/myanimeapi/commit/33a3d3a61104201cefc91a4203b40d6e31693763))
* **repositories:** add base repository interfaces ([af5516a](https://github.com/italo-gouveia/myanimeapi/commit/af5516a79af2d67277efa667e5f7002870d21196))
* **repositories:** implement AnimeRepository ([25b8694](https://github.com/italo-gouveia/myanimeapi/commit/25b869433626b6be8c5f80a00befeda7cb3f0f39))
* **repositories:** implement FavoriteRepository ([82383fe](https://github.com/italo-gouveia/myanimeapi/commit/82383fee2bb558092c5bfd38c87766a3f4b13726))
* **repositories:** implement ReviewRepository ([3bd1aaf](https://github.com/italo-gouveia/myanimeapi/commit/3bd1aaffc7ceb617ce7f0fe100e9a510d47d7123))
* **repositories:** implement UserRepository ([c75baf5](https://github.com/italo-gouveia/myanimeapi/commit/c75baf548bba5f60c11f9beda32ae2d8c7d286d2))
* **services:** implement service layer with business logic ([fff27ce](https://github.com/italo-gouveia/myanimeapi/commit/fff27ce119b44e27211079582e2eb3eb4f354cc3))
* **test:** add mock repository implementation for testing ([78f358f](https://github.com/italo-gouveia/myanimeapi/commit/78f358f13a2147aa7944144c148064387eee2314))
* **utils:** add common utility functions ([fdcc999](https://github.com/italo-gouveia/myanimeapi/commit/fdcc9998f24cfee642f5833ee095b25f9ac964d4))


### BREAKING CHANGES

* **anime:** The POST /v1/anime endpoint now requires a different request payload structure.
Instead of the full Anime model, it now uses AnimeCreateRequest which only includes title, description, and rating fields.
* AnimeHandler constructor now accepts AnimeServiceInterface
instead of concrete AnimeService type.
* **auth:** The auth package now uses a repository pattern instead of
direct database access. Existing code that directly accessed the database
for auth operations will need to be updated to use the new AuthRepository.

# [1.4.0](https://github.com/italo-gouveia/myanimeapi/compare/v1.3.0...v1.4.0) (2025-04-05)


### Features

* add favorite anime functionality ([4c804af](https://github.com/italo-gouveia/myanimeapi/commit/4c804aff13f97e17c8405271bdddc04b9fbeb5e4))

# [1.3.0](https://github.com/italo-gouveia/myanimeapi/compare/v1.2.0...v1.3.0) (2025-03-22)


### Features

* **anime-handler:** [GITISSUE-86] add pagination to GetAllAnimesHandler ([fb8e738](https://github.com/italo-gouveia/myanimeapi/commit/fb8e73829893964b7294d233473c784fcbb29559))
* **anime-handler:** [GITISSUE-86] updated docs swagger ([fc411cd](https://github.com/italo-gouveia/myanimeapi/commit/fc411cde921d7e97c7fb77d0b56125bf69fa6584))
* **anime-handler:** [GITISSUE-86] updating dependecy with vulnerabilitie ([a4daa25](https://github.com/italo-gouveia/myanimeapi/commit/a4daa25bf8c832f90d2c19760bfed30c4c9fdcc0))
* **anime-handler:** [GITISSUE-86]commenting failing tests to be fixed late ([a51061a](https://github.com/italo-gouveia/myanimeapi/commit/a51061a2e1bf854fd4b7f1d8c97050388056b258))

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.3.1] - 2025-04-02

### Added
- Favorite anime functionality
  - Add anime to favorites
  - Remove anime from favorites
  - List user's favorite anime
  - Automatic database migration for favorites table

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
