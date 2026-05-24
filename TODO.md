# TODO

* [x] Create models for other entities
* [x] Refactor controllers to depend on repository interfaces
* [x] Move dependency wiring to constructors and the application entry point
* [ ] Replace hardcoded current user values with middleware/context
* [ ] Update roles and authorization with Casbin
* [x] Create Postgres connection and configuration
* [x] Implement the artist Postgres repository
* [ ] Implement the remaining repository functions for Postgres
* [ ] Add tests for controllers and repositories
* [ ] Create docker compose
* [ ] Configure both service and DB in docker compose 
* [ ] Change Taskfile to run both service and DB from docker compose
* [ ] Write integration tests for artist
* [ ] Create santinel errors for controller
* [ ] Create error global specific like router/artist/error
* [ ] Create separate container in docke-compose for integration tests
* [ ] Add run integration tests from Taskfile
