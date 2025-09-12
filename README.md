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

---

## Submission Flow Overview

```
[Frontend User]
      |
      | POST /submissions
      v
[Backend API] --> insert Cassandra (pending)
      |
      | produce Kafka: submissions topic
      v
[Worker Service] --> execute code, determine status/result
      |
      | produce Kafka: submission-results topic
      v
[Backend Consumer] --> UpdateSubmissionResult in Cassandra
      |
      v
[Frontend User] GET /submissions/:id --> sees result
```

The submission process in LiteCode follows an asynchronous workflow using Kafka and Cassandra:

1. **User Submits Code** (`POST /submissions`)

   - The API writes the submission to Cassandra with `status = "pending"`.
   - The API publishes a Kafka message to the `submissions` topic with:
     ```json
     {
       "submission_id": "...",
       "user_id": "...",
       "problem_id": "...",
       "code": "...",
       "language": "..."
     }
     ```

2. **Worker Consumes Submission**

   - The worker reads the Kafka message from the `submissions` topic.
   - Executes the code in a Docker sandbox.
   - Calculates:
     - **Status**: `"success"`, `"runtime_error"`, `"compilation_error"`, etc.
     - **Result**: `"Accepted"`, `"Wrong Answer"`, `"Time Limit Exceeded"`, etc.
     - **Runtime** and **Memory** usage.

3. **Worker Publishes Result**

   - After processing, the worker publishes a result message to the `submission-results` Kafka topic:
     ```json
     {
       "submission_id": "...",
       "user_id": "...",
       "problem_id": "...",
       "status": "...",
       "result": "...",
       "runtime": 0.123,
       "memory": 2048
     }
     ```

4. **Backend Consumes Results**

   - The backend runs a consumer loop listening to the `submission-results` topic.
   - Calls `UpdateSubmissionResult(res)` to update all denormalized Cassandra tables:
     - `submissions_by_user`
     - `submissions_by_problem`
     - `submissions_by_problem_and_user`

5. **Frontend Fetches Updated Submissions**
   - Users call `GET /submissions/:id` (or a list endpoint).
   - Backend queries Cassandra and returns the latest **status** and **result** for display.
