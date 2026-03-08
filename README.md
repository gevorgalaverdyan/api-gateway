# This gets forwarded to a "users" service
curl http://localhost:8080/users

# This gets forwarded to a "products" service  
curl http://localhost:8080/products

# This hits an unknown route
curl http://localhost:8080/unknown
→ 404 Not Found
```

And in Grafana you'd see:
```
Total requests:        1,204
Requests per route:    /users → 800 | /products → 404
Error rate:            2%
Average response time: 120ms
```

---

## The architecture
```
[ Client ]
    │
    ▼
[ Your Go Gateway ]  ← the thing you build
    ├── /users    → forwards to → [ Mock Users Service ]
    ├── /products → forwards to → [ Mock Products Service ]
    └── /metrics  → scraped by → [ Prometheus ] → [ Grafana ]
```

The "Mock Services" are just tiny Go HTTP servers that return fake JSON — nothing fancy.

---

## What you'll learn

| Concept | Where |
|---|---|
| HTTP server & routing | The gateway itself |
| Reverse proxying | Forwarding requests to other services |
| Goroutines & concurrency | Handling many requests at once |
| Prometheus metrics | Counting requests, measuring latency |
| Docker | Packaging your service |
| Kubernetes | Running everything as pods |
| Grafana | Building the dashboard |

---

## Build phases
```
Phase 1 → Basic Go HTTP server with routing
Phase 2 → Forward requests to mock services
Phase 3 → Add Prometheus metrics
Phase 4 → Dockerize everything
Phase 5 → Deploy to Kubernetes
Phase 6 → Grafana dashboard
