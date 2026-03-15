# **Art Exchange Platform**

## **Project Goal:**

To build an online platform where artists can publish their works and users can view, like, and exchange artworks. The platform should be based on a scalable architecture and be easily deployable in the cloud.

---

## **Architecture:**

Microservice-based, with gRPC communication between services.

### Core Microservices:

1. **User Service**

   * User registration and authentication
   * Artist profile (name, avatar, bio, links)
   * Roles: artist, viewer, admin

2. **Art Service**

   * Artwork publishing (title, description, images, tags)
   * Association with the user (creator)
   * Moderation and publication status

3. **Like/Interaction Service**

   * Likes, views, favorites
   * Ability to follow artists

4. **Exchange Service** (optional)

   * Exchange requests between users for artworks
   * Statuses: pending, accepted, rejected
   * Optional message/comment on the exchange

---

## **Technologies:**

* Language: **Go**
* Framework: **Chi**
* Protocol: **REST**
* Database: **PostgreSQL**, **Redis**
* File storage: **AWS S3** (optional)
* Containerization: **Docker**, orchestration — **Kubernetes**
* CI/CD: **GitHub Actions**
* Authentication: **JWT**
* Runner [Task](https://taskfile.dev)

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
---
## License
[MIT](./LICENSE)
