# ⚡ Caching Proxy Server

A lightweight and configurable HTTP proxy server with caching functionality powered by Redis. Built using Go and the redis, this proxy caches responses to reduce repeated requests to the origin server and speed up response times.

---

## 📁 Project Structure

```text
akashgupta1909-caching-proxy-server/
├── configure.go       # Command-line flag parsing and config initialization
├── go.mod             # Module dependencies
├── go.sum             # Dependency checksums
├── handleProxy.go     # Proxy handler with cache logic
├── main.go            # Entry point for the application
├── redis.go           # Redis client logic for caching
└── router.go          # HTTP routing setup using chi
```

---

## 🚀 Features

- 🌐 Acts as a reverse proxy to a specified origin server.
- ⚡ Caches GET responses using Redis.
- ♻️ Supports automatic cache flushing via flag.
- 🌍 CORS support with pre-configured policies.
- 🛠 Command-line flags for flexible configuration.
- 🧪 Cache hit/miss detection (`X-Cache` response header).

---

## 🧰 Requirements

- Go 1.21+
- Redis server running on `localhost:6379`

---

## 🧪 Installation & Running

1. **Clone the repository:**

```bash
git clone https://github.com/akashgupta1909/caching-proxy-server.git
cd caching-proxy-server
```

2. **Build the application:**

```bash
go build -o caching-proxy-server main.go
```

3. **Start the Redis server:**

```bash
# If you have Redis installed, start it with:
redis-server
```

4. **Run the proxy server:**

```bash
./caching-proxy-server -origin <origin_server_url> -port <port_number> -cache-clean
```

Replace `<origin_server_url>` with the URL of the server you want to proxy to (e.g., `https://example.com`) and `<port_number>` with the desired port (default is `8080`).

## ⚙️ CLI Flags

| Flag            | Alias | Description                                  | Default    |
| --------------- | ----- | -------------------------------------------- | ---------- |
| `--port`        | `-p`  | Port to run the server on                    | `8080`     |
| `--origin`      | `-o`  | Origin base URL to forward proxy requests to | _Required_ |
| `--clean-cache` | `-c`  | Flush Redis cache before starting the server | `false`    |
| `--help`        | `-h`  | Show help message and exit                   |            |

### Example

```bash
go run main.go -o https://example.com -p 9090 -c
```

---

## 🧠 How It Works

1. When the server starts, it reads flags for configuration like port, origin URL, and cache clean-up.
2. The server listens for `GET` requests.
3. For each request:
   - It checks Redis for a cached response using the request path as the key.
   - If **cache hit**, it returns the cached response with header `X-Cache: HIT`.
   - If **cache miss**:
     - It forwards the request to the configured origin server.
     - It caches the response in Redis with a TTL of 1 hour.
     - It returns the fresh response to the client with header `X-Cache: MISS`.
4. If `--clean-cache` is passed, it flushes the entire Redis cache and exits.

## 🛠 To Do

- [ ] Add support for caching other HTTP methods like POST, PUT.
- [ ] Allow configurable TTL via CLI flags or config file.
- [ ] Create admin endpoints for manual cache invalidation.
- [ ] Add logging and metrics (e.g. Prometheus, Grafana integration).
- [ ] Add support for Redis cluster/failover mode.
- [ ] Add automated tests (unit/integration).
- [ ] Add Dockerfile and `docker-compose.yml` for easier deployment.
