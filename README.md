# Go Classifieds Watcher 🕵️‍♂️

A lightweight, concurrent real-time watcher for real estate listings, built with **Go** using **Clean Architecture**.

## 🏗 Architecture

- **Core:** Pure business logic.
- **Adapters:** 
  - `asi67`: Reverse-engineered API Client (10x faster than scraping).
  - `fs`: JSON-based full data persistence.
- **Config:** 12-Factor App compliant (Env vars).

## 🚀 Usage

1. Create a `.env` file (see source code for keys).
2. Run with Docker:
   ```bash
   make run