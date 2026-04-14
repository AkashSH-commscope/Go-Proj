# Edge telemetry bridge — milestone checklist (you implement)

Hands-on project: **simulated devices → MQTT (Mosquitto) → Go ingest → gRPC API** (unary + server streaming). Concepts and VRIOT references: [python-concurrency-and-mqtt-patterns.md](python-concurrency-and-mqtt-patterns.md).

**Suggested layout (your repo or folder, not committed by this guide):**

- `proto/` — `.proto` definitions  
- `cmd/ingest/` — MQTT subscriber + gRPC server binary  
- `cmd/sim/` — publisher simulator (Go or Python)  
- `docker-compose.yml` — Mosquitto + optional build of ingest  

---

## Phase 0 — Environment (1–2 h)

**Do**

- Install Go toolchain, `buf` or `protoc` + `protoc-gen-go` + `protoc-gen-go-grpc`, `mosquitto` client tools, optional `grpcurl`.
- Confirm `docker compose` works.

**Acceptance**

- `go version` and `docker compose version` succeed.
- You can run Mosquitto locally or via Compose and `mosquitto_sub -t '#' -v` sees test publishes.

---

## Phase 1 — Broker and topic contract (1–2 h)

**Do**

- Choose a topic pattern, e.g. `telemetry/{gateway_id}/readings`.
- Decide JSON payload fields (`gateway_id`, `timestamp`, `sensor_id`, `value`, optional `unit`).

**Acceptance**

- Written topic + JSON schema in your README (even if informal).
- Manual `mosquitto_pub` publishes; `mosquitto_sub` receives.

---

## Phase 2 — Device simulator (2–4 h)

**Do**

- Loop: publish N messages/sec per simulated gateway with realistic jitter.
- Use QoS 0 or 1 consciously; note trade-offs in README.

**Acceptance**

- Run simulator against Mosquitto; subscriber CLI shows stable flow for several minutes.

---

## Phase 3 — Go MQTT consumer skeleton (3–5 h)

**Do**

- Connect with a maintained Go MQTT client (e.g. Paho/Eclipse MQTT Go).
- Subscribe to your topic pattern; log parsed JSON.

**Acceptance**

- Consumer runs until Ctrl+C; no panic on malformed JSON (log/skip).
- Reconnect behavior documented (what happens if broker restarts).

---

## Phase 4 — In-memory store + worker pool (4–6 h)

**Do**

- MQTT message callback **does minimal work**; push to **buffered channel**.
- **Worker goroutines** parse + validate + append to a **ring buffer** or slice with **mutex** (cap e.g. last 10k readings).
- Use **`context`** for shutdown: cancel → close channel → workers exit (`sync.WaitGroup`).

**Acceptance**

- Under moderate publish rate, no sustained growth of goroutines or channel depth without backpressure strategy (document if you drop or block).

---

## Phase 5 — gRPC: `.proto` + codegen (2–4 h)

**Do**

- Define messages (`Reading`, `ListReadingsRequest`, etc.).
- **Unary:** e.g. `GetLatestReading(device_id)` or `ListRecentReadings(limit)`.
- **Server streaming:** e.g. `SubscribeReadings(filter)` streaming new readings as they arrive (push from ingest path via channel/safe map).

**Acceptance**

- `go build` succeeds with generated code in repo (or documented generation step).
- **Reflection** registered for `grpcurl` debugging.

---

## Phase 6 — gRPC server in same process (4–6 h)

**Do**

- Start **gRPC server** in the same binary as MQTT ingest (or separate binary + shared store—your choice; same binary is simpler for learning).
- Map **MQTT ingest** updates to **stream** subscribers (fan-out: careful with mutexes / broadcast channel patterns).

**Acceptance**

- `grpcurl` unary call returns data you ingested from MQTT.
- `grpcurl` streaming call shows new messages as simulator publishes.

---

## Phase 7 — Docker Compose + polish (3–5 h)

**Do**

- `docker-compose.yml`: Mosquitto service; optional multi-stage build for ingest binary.
- README: how to start stack, run simulator, run `grpcurl` examples, ports.

**Acceptance**

- Fresh clone: `docker compose up` + documented commands reproduce the demo.

---

## Phase 8 — Stretch (optional)

- Unary **ingest** RPC from a second tool (client sends readings over gRPC instead of MQTT) — compares transport roles.
- **grpc-gateway** slice exposing one JSON endpoint (study [`../riot-config-services/pluginmgr/main.go`](../riot-config-services/pluginmgr/main.go)).
- TLS for gRPC or MQTT (document only if you skip implementation).

---

## Optional code review (how to use help)

When you want feedback:

1. Complete a phase and its **acceptance** checks yourself.
2. Share a **branch link**, **patch**, or **paste** of the critical paths (MQTT handler, shutdown, gRPC stream fan-out).
3. Ask specific questions (e.g. “Is this channel close pattern safe?”).

Review is optional and does not require committing learning code to this workspace.

---

## Self-check summary

| Area | Question you can answer after finishing |
|------|----------------------------------------|
| MQTT | Why QoS 0 vs 1 for your simulator? |
| Go | How does your process exit cleanly? |
| gRPC | Why unary vs stream for “live” readings? |
| Ops | What happens if broker dies mid-run? |
