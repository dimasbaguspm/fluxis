# Product Requirements Document (PRD): Fluxis

| Project Name | Status      | Version | Author              | Date         |
| :----------- | :---------- | :------ | :------------------ | :----------- |
| **Fluxis**   | Draft / MVP | 1.0     | Gemini (AI Partner) | Jan 11, 2026 |

---

## 1. Executive Summary

**Fluxis** is a minimalist, self-hosted project management tool designed for personal use (Single Admin). It combines a structured **Todo List** with a **Task Logger** (Activity Feed) to help users track progress across multiple isolated projects. Fluxis offers two primary interfaces: a modern Web Dashboard and a high-performance CLI/TUI.

---

## 2. Goals & Vision

- **Dual Interface Access:** Provide a seamless experience between a React-based Web UI and a Bubble Tea-based Terminal UI.
- **Project Isolation:** Ensure tasks and logs are strictly categorized by project to maintain focus.
- **Self-Hosted Sovereignty:** Built for users who want to host their own data using Docker and Postgres.
- **Auditability:** Every project has a dedicated "Log" to track history, thoughts, and progress over time.

---

## 3. User Persona

- **The Power User / Developer:** An individual who spends time in both the browser and the terminal, needing a unified tool to log work progress and manage "to-do" items without the overhead of enterprise software.

---

## 4. Functional Requirements

### 4.1 Authentication & Access

| ID        | Feature           | Description                                                                     |
| :-------- | :---------------- | :------------------------------------------------------------------------------ |
| **FR-01** | Single-Admin Auth | Access is restricted to one admin account configured via environment variables. |
| **FR-02** | JWT-based Auth    | Secure communication between the Go backend and both the Web and CLI frontends. |

### 4.2 Project Management

| ID        | Feature        | Description                                               |
| :-------- | :------------- | :-------------------------------------------------------- |
| **FR-03** | Project CRUD   | Create, view, update, and delete distinct projects.       |
| **FR-04** | Project Status | Categorize projects as `Active`, `Paused`, or `Archived`. |

### 4.3 Task Management (Todo List)

| ID        | Feature         | Description                                      |
| :-------- | :-------------- | :----------------------------------------------- |
| **FR-05** | Task CRUD       | Manage tasks within a specific project.          |
| **FR-06** | Task States     | Basic states: `Todo`, `In Progress`, `Done`.     |
| **FR-07** | Priority Levels | Assign priority labels: `Low`, `Medium`, `High`. |

### 4.4 Task & Activity Logger

| ID        | Feature          | Description                                                                         |
| :-------- | :--------------- | :---------------------------------------------------------------------------------- |
| **FR-08** | Activity Logging | Add text entries to a project's timeline (e.g., "Finished the database migration"). |
| **FR-09** | Task Correlation | (Optional) Link a specific log entry to a specific task for better context.         |

---

## 5. Technical Stack

- **Backend:** Golang Standard Library.
- **Frontend (Web):** Vue.
- **Frontend (CLI):** Bubble Tea
- **Database:** PostgreSQL
- **Infrastructure:**
  - **Containerization:** Docker & Docker Compose.
  - **Registry:** GitHub Container Registry (GHCR) for image hosting.
