# LiteCode

A LeetCode-style platform for mobile.  
This side project is also an exercise in designing and experimenting with scalable architecture.

## Planned Tech Stack

**Backend**

- **Go** – core backend services
- **GORM** – ORM for PostgreSQL integration
- **Postman** – API testing and development
- **Neon (PostgreSQL)** – for structured, lower-throughput data (problems, test cases, users, etc.)
- **Cassandra** – for high-throughput write-heavy workloads (e.g., submissions)
- **Redis** – caching and session management
- **Kafka** – distributed submission queue and event streaming
- **Docker** – sandboxed code execution
- **Kubernetes** – orchestration of containerized execution environments

**Frontend**

- **React Native** – cross-platform mobile app

---

## Why this stack

The goal is to combine practicality with scalability:

- **Go** offers high performance, concurrency, and simplicity for backend services.
- **PostgreSQL (via Neon)** is reliable for relational data that benefits from strong consistency and structured queries.
- **Cassandra** excels at high-throughput, write-heavy workloads, making it well-suited for large volumes of submissions.
- **Redis** provides fast in-memory caching and session management to reduce database load and speed up requests.
- **Kafka** enables reliable, scalable event-driven processing for handling submission queues.
- **Docker** isolates code execution securely, ensuring user submissions run in sandboxed environments.
- **Kubernetes** adds the ability to scale and orchestrate these containers efficiently in production.
- **React Native** allows a single codebase for iOS and Android, accelerating mobile development.

This mix ensures the platform can scale horizontally as usage grows while still being approachable for a side project.
