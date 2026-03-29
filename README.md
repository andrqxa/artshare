# **ArtShare**

## **Project Goal:**

ArtShare is an application which allows artists to publish artworks, and other users to browse them, like them, and create exchange requests. The focus of the project is a **simple, maintainable monolithic application** with a clear internal structure and straightforward local/cloud deployment.

---

## **Architecture:**

**Monolithic MVC application** with a single REST API backend.

The project is structured as one Go application with clear separation of responsibilities:

- **Models** — domain entities and persistence-related structures
- **Controllers** — HTTP handlers and request orchestration
- **Views** — JSON responses now, with the `web/` directory reserved for templates/static assets
- **Repositories** — database access layer
- **Middleware** — cross-cutting HTTP concerns
- **Router** — route registration and composition

---

## MVP Modules

### 1. Auth

- register user
- login user
- JWT-based authentication

### 2. Users

- get current user profile
- update current user profile
- view public artist profile

### 3. Artworks

- create artwork
- update artwork
- delete artwork
- list artworks
- get artwork details

### 4. Interactions

- like artwork
- remove like from artwork

### 5. Exchanges

- create exchange request
- list current user's exchange requests
- get exchange request details
- accept / reject / cancel exchange request

---
## Technology Stack

- **Language:** Go
- **HTTP router:** Chi
- **Architecture style:** MVC
- **Protocol:** REST
- **Database:** PostgreSQL
- **Authentication:** JWT
- **Containerization:** Docker
- **Task runner:** Task
- **API documentation:** OpenAPI 3

---

## Repository Structure

```text
.
├── api
│   └── openapi.yaml
├── bin
├── ci
│   ├── deployments
│   ├── docker-compose.yml
│   ├── Dockerfile
│   └── scripts
├── cmd
│   └── artshare
│       └── main.go
├── docs
│   ├── component.drawio
│   └── Context.drawio
├── go.mod
├── go.sum
├── internal
│   ├── configs
│   │   └── mainConfig.go
│   ├── controller
│   ├── middleware
│   ├── model
│   ├── pkg
│   ├── repository
│   └── router
├── LICENSE
├── migrations
├── pkg
│   ├── api
│   │   └── v1
│   ├── db
│   └── util
├── README.md
├── Taskfile.yml
├── test
└── web
    ├── app
    ├── static
    └── template
```

---

## MVC Mapping in This Project

### Controllers

- auth controller
- user controller
- artwork controller
- exchange controller
- interaction controller

### Models

- user
- artist profile
- artwork
- like
- exchange request

### Views

- JSON API responses
- validation and error responses
- optional server-rendered templates later via `web/template`

---

## API Contract

The draft OpenAPI specification for the current MVP lives in [api/openapi.yaml](./api/openapi.yaml).

It is intentionally aligned with the minimal feature set described in this README.

---

## Non-functional Goals

- clear and understandable project structure
- minimal but consistent API contract
- support for local development
- easy containerized запуск and deployment
- testable code organization

---

## License

[MIT](./LICENSE)

---
