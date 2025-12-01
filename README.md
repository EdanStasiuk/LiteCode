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
- **RabbitMQ** – asynchronous submission queue
- **Docker** – sandboxed code execution
- **Kubernetes** – orchestration of containerized execution environments

**Frontend**

- Not yet implemented - planned React Native mobile app

---

## Why this stack

The goal is to combine practicality with scalability:

- **Go** offers high performance, concurrency, and simplicity for backend services.
- **PostgreSQL (via Neon)** is reliable for relational data that benefits from strong consistency and structured queries.
- **Cassandra** excels at high-throughput, write-heavy workloads, making it well-suited for large volumes of submissions.
- **Redis** provides fast in-memory caching and session management to reduce database load and speed up requests.
- **RabbitMQ** provides a reliable, asynchronous queue for submissions and results.
- **Docker** isolates code execution securely, ensuring user submissions run in sandboxed environments.
- **Kubernetes** adds the ability to scale and orchestrate these containers efficiently in production.
- **React Native** allows a single codebase for iOS and Android, accelerating mobile development.

This stack allows horizontal scaling as usage grows while remaining manageable for a side project.

---

## Submission Flow Overview

```
[Frontend User]
    |
    | POST /submissions
    v
[Backend API] --> insert into Cassandra (pending)
    |
    | produce RabbitMQ: submissions queue
    v
[Worker Service] --> execute code in Docker sandbox, determine status/result
    |
    | produce RabbitMQ: submission-results queue
    v
[Backend Consumer] --> UpdateSubmissionResult in Cassandra
    |
    v
[Frontend User] GET /submissions/:id --> sees result
```

Or simply

```
Frontend POST /submissions -> RabbitMQ submissions queue -> Worker -> RabbitMQ submission-results queue -> Backend consumer updates Cassandra -> Frontend GET /submissions/:id
```

---

## Submission Process

1. **User Submits Code (`POST /submissions`)**

   - The API writes the submission to Cassandra with `status = "pending"`.
   - The API publishes a RabbitMQ message to the `submissions` queue:

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

   - The worker reads the RabbitMQ message from the `submissions` queue.
   - Executes the code in a Docker sandbox.
   - Determines:
     - **Status**: `"success"`, `"runtime_error"`, `"compilation_error"`, etc.
     - **Result**: `"Accepted"`, `"Wrong Answer"`, `"Time Limit Exceeded"`, etc.
     - **Runtime** and **Memory** usage.

3. **Worker Publishes Result**

   - After execution, the worker publishes a message to the `submission-results` queue:

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

   - A backend consumer listens to the `submission-results` queue.
   - Calls `UpdateSubmissionResult(res)` to update all denormalized Cassandra tables:
     - `submissions_by_user`
     - `submissions_by_problem`
     - `submissions_by_problem_and_user`

5. **Frontend Fetches Updated Submissions**

   - Users call `GET /submissions/:id` or list endpoints.
   - The backend returns the latest **status** and **result** from Cassandra.

---

### Notes

- LiteCode is currently **backend-only**; no frontend or mobile app exists yet.
- Code execution is **not implemented yet**
- Docker/Kubernetes-based sandboxing is planned for real code execution in the future.
