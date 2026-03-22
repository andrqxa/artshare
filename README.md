# **Art Exchange Platform**

## **Project Goal:**

To build an online platform where artists can publish their works and users can view, like, and exchange artworks. The platform should be based on a simple and maintainable architecture and be easily deployable in the cloud.

---

## **Architecture:**

**Monolithic MVC application** with a single REST API backend.

The system is built as one application with clear internal separation into layers:

* **Models** — domain entities and database access
* **Views** — API responses / optional server-rendered pages
* **Controllers** — request handling and routing
* **Services** — business logic

This approach keeps the project simpler to build, test, and deploy than a microservice-based solution.

### Core Modules:

1. **User Module**

   * User registration and authentication
   * Artist profile management (name, avatar, bio, links)
   * Roles: artist, viewer, admin

2. **Art Module**

   * Artwork publishing (title, description, images, tags)
   * Association with the user (creator)
   * Moderation and publication status

3. **Interaction Module**

   * Likes, views, favorites
   * Ability to follow artists

4. **Exchange Module** (optional)

   * Exchange requests between users for artworks
   * Statuses: pending, accepted, rejected
   * Optional message/comment on the exchange

---

## **Technologies:**

* Language: **Go**
* Framework: **Chi**
* Architecture style: **MVC**
* Protocol: **REST**
* Database: **PostgreSQL**
* Cache / session / counters: **Redis** (optional)
* File storage: **AWS S3** (optional)
* Containerization: **Docker**
* Deployment: **Kubernetes** (optional for future scaling)
* CI/CD: **GitHub Actions**
* Authentication: **JWT**
* Runner: [Task](https://taskfile.dev)

---

## **Application Structure Example:**

```text
api/
bin/
cmd/
  artshare/
deployments/
internal/
  app/
  auth/
  configs/
  controller/
  middleware/
  model/
  pkg/
  repository/
  service/
  storage/
/migrations/
pkg/
scripts/
test/
web/
  app/
  static/
  template/
```

### Example MVC mapping:

* **Controllers**

  * auth controller
  * user controller
  * artwork controller
  * interaction controller
  * exchange controller

* **Models**

  * user
  * artist_profile
  * artwork
  * like
  * favorite
  * follow
  * exchange_request

* **Views**

  * JSON API responses
  * validation / error responses
  * optional frontend templates if needed later

---

## **User Features:**

* Sign up / log in
* Profile editing
* Upload and publish artwork
* Browse other artists’ galleries
* Like, follow, favorite
* Exchange artworks

---

## **Non-functional Requirements:**

* API documentation
* Testing (unit + integration)
* Clear modular structure inside a monolith
* Easy local development and deployment

---

## License

[MIT](./LICENSE)

---