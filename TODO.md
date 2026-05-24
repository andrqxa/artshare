# TODO

* [x] Create models for other entities
* [x] Refactor controllers to depend on repository interfaces
* [x] Move dependency wiring to constructors and the application entry point
* [x] Replace hardcoded current user values with middleware/context
* [ ] Update roles and authorization with Casbin
* [x] Create Postgres connection and configuration
* [x] Implement the artist Postgres repository
* [x] Implement the remaining repository functions for Postgres
* [x] Add tests for controllers and repositories
* [x] Create docker compose
* [x] Configure both service and DB in docker compose
* [x] Change Taskfile to run both service and DB from docker compose
* [x] Write integration tests for artist
* [x] Create sentinel errors for controller
* [x] Create error global specific like router/artist/error
* [x] Create separate container in docke-compose for integration tests
* [x] Add run integration tests from Taskfile
