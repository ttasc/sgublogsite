# 📝 SGU Blog Platform

> A full-stack blog & content management system built with Golang, featuring server-side rendering, admin dashboard, and type-safe database access using sqlc.

---

## 📌 Table of Contents
- [About The Project](#-about-the-project)
- [Features](#-features)
- [Tech Stack](#-tech-stack)
- [Architecture Overview](#-architecture-overview)
- [Project Structure](#-project-structure)
- [Getting Started](#-getting-started)
- [Usage](#-usage)
- [Database Design](#-database-design)
- [Technical Highlights](#-technical-highlights)
- [License](#-license)

---

## 📖 About The Project

**SGU Blog Platform** is a full-featured blogging system built with **Golang**, designed to manage content, users, and media in a structured and scalable way.

Unlike typical SPA-based applications, this project uses:
- **Server-Side Rendering (SSR)** for performance and simplicity
- **HTMX-enhanced interactions** for dynamic UI without heavy frontend frameworks
- **sqlc for type-safe SQL** instead of ORM

It is designed as a **lightweight CMS** suitable for:
- Educational environments
- Internal content systems
- Personal blogging platforms

---

## ✨ Features

### 🌐 Public Site
- 📰 Blog posts (news, announcements)
- 🔍 Search functionality
- 🏷️ Categories & tags
- 📄 Static pages (about, contact)
- 👤 User profile

---

### ⚙️ Admin Dashboard
- 📝 Manage posts (CRUD)
- 🏷️ Manage categories & tags
- 👥 Manage users
- 🖼️ Manage uploaded images
- 📊 Dashboard overview

---

### 👤 User System
- Login page
- Profile management
- Avatar upload

---

### ⚡ Frontend Interaction
- Server-side rendered templates
- HTMX for partial updates
- Minimal JavaScript (progressive enhancement)

---

## 🧰 Tech Stack

| Category        | Technology |
|----------------|------------|
| Language       | Go (Golang) |
| Architecture   | MVC-like layered structure |
| Database       | SQL (via sqlc) |
| Query Tool     | sqlc |
| Frontend       | HTML Templates + HTMX + Vanilla JS |
| Styling        | CSS |
| Container      | Docker + Docker Compose |
| Dev Tooling    | Air (live reload), Makefile |

---

## 🏗️ Architecture Overview

The system follows a layered backend design:

### Request Flow

```text
HTTP Request
    ↓
Router (server/routes.go)
    ↓
Controller (internal/controller)
    ↓
Model / Repository (sqlc + model layer)
    ↓
Database
    ↓
Template Rendering (HTML)
    ↓
HTTP Response
````

---

### Core Layers

* **Controller Layer**

  * Handles HTTP requests
  * Maps routes to business logic
  * Example:

    * `admin_posts.go`
    * `login.go`
    * `search.go`

* **Model Layer**

  * Defines domain entities
  * Integrates with sqlc-generated queries

* **Server Layer**

  * Routing and HTTP server setup

* **Utils Layer**

  * Database connection
  * Shared utilities

---

## 📁 Project Structure

```bash
.
├── src/
│   ├── cmd/
│   │   └── main.go          # Entry point
│   └── internal/
│       ├── controller/      # HTTP handlers (admin + public)
│       ├── model/           # Domain models
│       ├── server/          # Router & server setup
│       └── utils/           # DB & helper functions
│
├── queries/                 # SQL queries (for sqlc)
├── templates/               # HTML templates (SSR)
├── assets/                  # Static assets (CSS, JS, images)
├── statics/                 # Static HTML pages
├── uploads/                 # Uploaded files (images, avatars)
│
├── schema.sql               # Database schema
├── sqlc.yml                 # sqlc configuration
├── docker-compose.yml       # Multi-container setup
├── Dockerfile               # App container
├── Makefile                 # Dev commands
└── .env                     # Environment variables
```

---

## 🚀 Getting Started

### ⚙️ Prerequisites

* Go 1.20+
* Docker & Docker Compose

---

### 🐳 Run with Docker (Recommended)

```bash
docker-compose up --build
```

---

### 💻 Run Locally

1. Setup environment:

```bash
cp .env.example .env
```

2. Start database (Docker or local)

3. Run app:

```bash
go run ./src/cmd
```

or:

```bash
make run
```

---

## ▶️ Usage

* Open: `http://localhost:<port>`
* Admin panel: `/admin`
* Login: `/login`

---

## 🗄️ Database Design

Tables include:

* `users`
* `posts`
* `categories`
* `tags`
* `images`
* `siteinfo`

All SQL queries are:

* Written manually in `/queries`
* Compiled using **sqlc**
* Fully type-safe in Go

---

## ⚡ Technical Highlights

### 🧠 Type-Safe SQL (sqlc)

* Eliminates ORM overhead
* Compile-time query validation
* Strong typing between DB and Go

---

### ⚙️ SSR + HTMX Architecture

* No heavy frontend framework
* Faster load time
* Simpler deployment
* Progressive enhancement with HTMX

---

### 🧩 Modular Backend Design

* Clear separation:

  * controller
  * model
  * server
* Easy to extend and maintain

---

### 🐳 Containerized Environment

* Reproducible setup
* Easy deployment

---

### 🖼️ Media Handling

* Image uploads
* Avatar management
* Organized storage structure

---

## 📄 License

No license specified yet.

---

## 🙌 Acknowledgements

* Golang community
* sqlc
* HTMX
* Docker

---

> 💡 This project demonstrates backend engineering skills in Go, including SSR architecture, type-safe database access, and modular system design.
