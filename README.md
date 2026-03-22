# 🚀 Distributed Task Scheduler (Go Microservices)

> A high-performance distributed system built with **Golang**, using **gRPC**, designed with **Clean Architecture**, and powered by **Go concurrency (channels & context)**.

---

## ⚡ Overview

This system allows clients to submit tasks which are intelligently distributed and executed across worker nodes in a scalable and fault-tolerant manner.

It demonstrates:

* ⚡ High concurrency using Goroutines & Channels
* 🔗 Efficient service-to-service communication via gRPC
* 🧠 Smart task orchestration
* 🧱 Clean Architecture for maintainability

---

## 🏗️ Architecture

```id="4b6y2p"
        Client
          |
          | HTTP
          ↓
+----------------------+
| Task Scheduler       |
+----------+-----------+
           |
           | gRPC
           ↓
+----------------------+
| Coordinator Service  |
+----------+-----------+
           |
           | Assign Task
           ↓
+----------------------+
| Worker Service(s)    |
+----------------------+
```

---

## 🔄 Working Flow

1. 📩 **Client → Task Scheduler**

   * User submits a task request

2. 🧠 **Task Scheduler**

   * Validates and processes the task
   * Forwards task for execution

3. 🎯 **Coordinator Service**

   * Decides *which worker* should handle the task
   * Handles load balancing / assignment logic

4. ⚙️ **Worker Service**

   * Executes the task
   * Returns result/status

---

## 🧬 Tech Stack

* 🟦 **Golang**
* 🔗 **gRPC (Protocol Buffers)**
* 🧱 **Clean Architecture**
* 🐳 **Docker**
* ⚡ **Channels & Goroutines**
* ⏱️ **Context API (timeouts, cancellation)**

---

## 🧠 Core Engineering Concepts

### 🔗 gRPC Communication

* Strong contracts using `.proto` files
* Low-latency service communication
* Scalable microservice interaction

---

### 🧱 Clean Architecture

Each service follows strict separation:

```id="y0kp63"
cmd/                → Application entry
internal/
    ├── domain/     → Core business logic
    ├── service/    → Use cases
    ├── transport/  → gRPC handlers
    └── repository/ → Data access
```

✔ Decoupled
✔ Testable
✔ Maintainable

---

### ⚡ Concurrency Design

#### 🧵 Channels for Task Flow

* Tasks are streamed internally using channels
* Enables async, non-blocking processing


#### ⏱️ Context for Control

* Timeout control for gRPC calls
* Graceful shutdown of services
* Cancellation propagation

```go id="g5tx7y"
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()
```

---

## 🐳 Running the Project

```bash id="v6j9bc"
git clone https://github.com/<your-username>/distributed-task-scheduler.git
cd distributed-task-scheduler

# Setup env files
cp worker_service/.env.example worker_service/.env
cp coordinator_service/.env.example coordinator_service/.env
cp task_scheduler_service/.env.example task_scheduler_service/.env

# Run services
docker-compose up --build
```

---

## 📁 Project Structure

```id="p3t9xk"
distributed-task-scheduler/
│
├── coordinator_service/
├── task_scheduler_service/
├── worker_service/
│
├── proto/                  → gRPC contracts
├── docker-compose.yml
└── README.md
```

---

## 🔮 Future Enhancements

* 🔁 Retry & failure handling
* 📊 Observability (logs, metrics, tracing)
* 🧠 Intelligent scheduling strategies
* 🌐 Horizontal scaling (multiple workers)
* 🔐 Secure inter-service communication

---

## 💥 Why This Project is Strong

✔ Real distributed system design
✔ Deep use of Go concurrency primitives
✔ Clean and scalable architecture
✔ Practical gRPC microservices implementation
✔ Production-ready foundation

---

## 👨‍💻 Author

Built with ⚡ by **Gaganpreet**

---

## ⭐ Support

If you like this project, consider giving it a ⭐ on GitHub!
