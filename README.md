# 🕵️‍♂️ Go Classifieds Watcher

![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go)
![Architecture](https://img.shields.io/badge/Architecture-Hexagonal-orange?style=for-the-badge)
![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

> **From a dusty single-threaded script to a high-performance concurrent system.**

## 📖 The Story

This project is the spiritual successor to an old legacy automation script. The goal was not just to bring it back to life, but to use it as a **playground to demonstrate modern Go capabilities** regarding concurrency, architecture, and performance optimization.

What started as a heavy browser automation task (Headless Chrome) was refactored into a **lightweight, reverse-engineered API client**, drastically reducing resource consumption and execution time.

## ⚡ Key Technical Highlights

### 1. From "Brute Force" to "Smart Engineering"
Initial versions used `chromedp` to render JavaScript and scrape the DOM. While functional, it was resource-heavy.
**The Pivot:** By reverse-engineering the target's private API, I transitioned to a pure HTTP implementation.

| Metric | Headless Chrome Strategy | Reverse-Engineered API Strategy | Improvement |
| :--- | :--- | :--- | :--- |
| **Execution Time** | ~45s | **~1.5s** | **30x Faster** |
| **Memory Usage** | ~600MB | **~15MB** | **40x Lighter** |
| **Docker Image** | ~500MB | **~20MB** | **Alpine Native** |

### 2. Concurrency Pattern (Fan-Out / Fan-In)
To fetch paginated results efficiently, the application uses a **Fan-Out** pattern:
1.  **Discovery:** Fetches Page 1 synchronously to determine total item count.
2.  **Fan-Out:** Spawns Goroutines to fetch all remaining pages (2..N) in parallel.
3.  **Fan-In:** Uses a `Mutex` protected slice to aggregate results safely.
4.  **Resilience:** Uses a `sync.WaitGroup` to ensure all routines complete before processing.

### 3. Clean Architecture (Hexagonal)
The project is strictly structured to decouple business logic from infrastructure:
* **Core (Domain):** Contains `Item` entities and the `WatcherService`. Zero dependencies.
* **Ports:** Interfaces defining `Provider`, `Repository`, and `Notifier`.
* **Adapters:**
    * `asi67`: The HTTP Client adapter (API).
    * `fs`: A JSON-based file persistence adapter.
    * `rememberme`: The HTTP Client adapter (API).
    * `std`: Structured logging adapter (`slog`).

## 🏗 Project Structure

```bash
.
├── cmd/
│   └── watcher/       # Application Entrypoint (Wiring & Config)
├── internal/
│   ├── core/          # PURE BUSINESS LOGIC (Ports & Domain)
│   ├── config/        # 12-Factor App Configuration
│   └── adapters/      # INFRASTRUCTURE LAYERS
│       ├── asi67/     # API Client (HTTP, JSON Unmarshal)
│       ├── fs/        # File System Persistence
│       └── std/       # Logging wrapper
├── .docker/           # Docker build context
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

**Using Make (Standard):**
```bash
make run
```

**Using Docker Compose manually:**
```bash
docker-compose up --build
```

### 3. Output
The application will:
1.  **Fetch** items from the API in parallel.
2.  **Filter** duplicates.
3.  **Compare** with the local history (`data/seen.json`).
4.  **Log** new items to stdout.
5.  **Exit** cleanly.

To verify persistence, run the command a second time. It should report `0 new items`.

### 4. Development
To clean artifacts (binary and local database) and start fresh:
```bash
make clean
```

## 🔮 Future Roadmap

* [ ] **Notification Adapter:** Replace Logger with Discord/Slack Webhooks.
* [ ] **Metrics:** Add Prometheus adapter for monitoring execution duration.