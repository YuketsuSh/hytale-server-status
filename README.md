# Hytale Server Status

> **Open-source daemon for querying Hytale servers using QUIC + TLS 1.3 and exposing live server status through a unified REST API.**

**Hytale Server Status** is a standalone, high-performance daemon written in **Go** that connects to real Hytale servers, retrieves live status information (online state, MOTD, players, latency, version), and exposes it through a simple **JSON REST API**.

It is designed to be language-agnostic thanks to lightweight **wrappers** for PHP, Java, Node.js, and more.

---

## ✨ Features

- ✅ Native **QUIC + TLS 1.3** connection
- 📡 Real Hytale server status (online/offline)
- 📝 MOTD & server version
- 👥 Players online / max players
- ⏱️ Accurate latency measurement (ms)
- 🚀 High-performance **Go daemon**
- 🧠 In-memory cache with configurable TTL
- 🌍 REST JSON API
- 🔌 Multi-language wrappers (PHP, Java, Node.js)
- 🐳 Docker-ready & cross-platform
- 📊 Designed for monitoring & integrations

---

## 🧱 Repository Structure

```text
hytale-server-status/
│
├── daemon/                 # Core Go daemon (QUIC + API)
│   ├── src/
│   ├── Dockerfile
│   ├── go.mod
│   └── README.md
│
├── wrappers/               # Language wrappers consuming the REST API
│   ├── php/
│   │   ├── src/
│   │   └── README.md
│   │
│   ├── java/
│   │   ├── src/
│   │   └── README.md
│   │
│   └── node/
│       ├── src/
│       └── README.md
│
├── docs/                   # Documentation & specs
│   ├── api.md
│   ├── protocol.md
│   └── architecture.md
│
├── .github/
│   └── workflows/
│
├── LICENSE
├── README.md
└── .gitignore
````

---

## 🚀 Quick Start (Daemon)

### Requirements

* Go **1.21+**
* UDP outbound access (QUIC)

### Run locally

```bash
cd daemon
go run .
```

The REST API will be available at:

```
http://localhost:8080
```

---

## 🔌 REST API Example

### Status endpoint

```http
GET /status?host=example.com&port=5520
```

### Example response

```json
{
  "address": "example.com:5520",
  "online": true,
  "motd": "Welcome to Hytale!",
  "server_version": "1.0.0",
  "players": {
    "online": 12,
    "max": 100
  },
  "latency_ms": 87,
  "packet_type": "status"
}
```

---

## 🔧 Wrappers

Wrappers are lightweight clients that consume the daemon REST API.

| Language | Path            |
| -------- | --------------- |
| PHP      | `wrappers/php`  |
| Java     | `wrappers/java` |
| Node.js  | `wrappers/node` |

Each wrapper has its own README with usage examples.

---

## 📚 Documentation

* 📘 **API Reference** → `docs/api.md`
* 📐 **Architecture** → `docs/architecture.md`
* 🔌 **Hytale Protocol Notes** → `docs/protocol.md`

---

## 🛡️ Security Notes

* The daemon should **not be exposed directly to the public Internet**
* Use a reverse proxy or firewall if deployed publicly
* Optional API token support (planned)
* No sensitive data is stored

---

## 🐳 Docker (Planned)

A production-ready Docker image will be provided for easy deployment.

---

## 🛡️ License

This project is licensed under the **MIT License**.
You are free to use, modify, and distribute it.

---

## ❤️ Contributing

Contributions are welcome!

* Open an issue for bugs or feature requests
* Submit pull requests for improvements
* Keep code clean and documented

---

## ⚠️ Disclaimer

This project is **not affiliated with Riot Games or Hypixel Studios**.

Hytale is a trademark of **Hypixel Studios**.
This project is a community-made, open-source tool.

---

## ⭐ Support the project

If you find this project useful, consider giving it a ⭐ on GitHub!