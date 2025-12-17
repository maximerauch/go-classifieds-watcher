# 🕵️‍♂️ Go Classifieds Watcher

![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8?style=for-the-badge&logo=go)
![Architecture](https://img.shields.io/badge/Architecture-Hexagonal-orange?style=for-the-badge)
![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)
![CI Status](https://img.shields.io/github/actions/workflow/status/maximerauch/go-classifieds-watcher/unit-tests.yml?branch=master&label=Tests&logo=github&style=for-the-badge)

> **From a dusty single-threaded script to a high-performance, cloud-native concurrent system.**

## 📖 The Story

This project is the spiritual successor to an old legacy automation script. The goal was not just to bring it back to life, but to use it as a **playground to demonstrate modern Go capabilities** regarding concurrency, architecture, and reliability.

What started as a heavy browser automation task was refactored into a **lightweight, intelligent agent** capable of handling both JSON APIs and HTML scraping with industrial-grade reliability, designed to run continuously as a cloud daemon.

## ⚡ Key Technical Highlights

### 1. Hybrid Data Fetching Strategies
The system adapts to the source target using distinct strategies, all behind a unified Port:
* **Reverse-Engineered API (ASI67):** Bypasses the UI to hit private JSON endpoints directly (~30x faster than Headless Chrome).
* **Polite HTML Scraping (RememberMe):** Uses `goquery` to parse DOM efficiently.
  * *Feature:* Implements **Rate Limiting** via a weighted semaphore pattern to respect the target server's load and avoid IP bans.

### 2. Advanced Concurrency Patterns
To fetch data efficiently without overwhelming resources:
* **Fan-Out/Fan-In:** Spawns Goroutines to fetch pages in parallel.
* **Semaphore Pattern:** Limits concurrent HTTP requests (e.g., max 5 active workers) during scraping.
* **Mutex Synchronization:** Aggregates results safely into shared slices.

### 3. Clean Architecture (Hexagonal)
The project is strictly structured to decouple business logic from infrastructure.
* **Core (Domain):** Pure business logic. Zero external dependencies.
* **Ports:** Interfaces defining `Provider`, `Repository`, and `Notifier`.
* **Adapters:**
  * `postgres`: SQL persistence using `lib/pq`, optimized for PaaS (Heroku/Scalingo) with SSL support.
  * `fs`: Fallback JSON file persistence for local development.
  * `email`: SMTP adapter using `gomail` for rich HTML alerts.
  * `composite`: **Composite Pattern** implementation to broadcast notifications to multiple channels (Logs + Email) simultaneously.

### 4. Cloud-Native & Daemon Mode
Unlike simple scripts, this application is designed to run as a long-living **Daemon**:
* **Ticker Loop:** Executes scans on a defined schedule (e.g., every 10m).
* **Graceful Shutdown:** Handles `SIGTERM` signals to finish current jobs before exiting.
* **Stateless:** Uses PostgreSQL to persist state, allowing the application to be deployed on ephemeral file systems.

## 🏗 Project Structure

```bash
.
├── cmd/
│   └── watcher/       # Application Entrypoint (Wiring & Config)
├── internal/
│   ├── core/          # PURE BUSINESS LOGIC (Ports & Domain)
│   ├── config/        # 12-Factor App Configuration
│   └── adapters/      # INFRASTRUCTURE LAYERS
│       ├── asi67/     # API Client (HTTP JSON)
│       ├── rememberme/# HTML Scraper (goquery + Rate Limiting)
│       ├── postgres/  # SQL Repository (lib/pq)
│       ├── email/     # SMTP Client
│       ├── composite/ # Composite Pattern for Notifiers
│       └── std/       # Structured Logging (slog)
├── .docker/           # Docker build context
└── Procfile           # PaaS Deployment Config (Worker definition)
└── data/              # Local persistence volume (GitIgnored)
```

## 🚀 Usage

### 1. Configuration
The application follows the **12-Factor App** methodology. Configuration is managed via environment variables.

Copy the example configuration:
```bash
cp .env.example .env
```

### 2. Running with Docker
This approach requires no local Go installation. It runs the application as a one-off job inside a lightweight Alpine container.

```bash
docker-compose up --build -d
```

### 3. PaaS Deployment (Scalingo / Heroku)
The project includes a Procfile for seamless cloud deployment.
1. Provision a PostgreSQL addon on your PaaS.
2. Set the environment variables. 
3. **Important:** Scale the worker process to 1 and web to 0 (as this is a background daemon).

## 🔮 Future Roadmap

* [x] PostgreSQL Support: Enable stateless cloud deployment.
* [x] Email Notifications: Rich HTML alerts.
* [x] Daemon Mode: Continuous monitoring with Graceful Shutdown.
* [ ] **Metrics:** Add Prometheus adapter for monitoring execution duration.