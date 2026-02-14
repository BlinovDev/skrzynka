# Application Logic — Modules

## Overview

Anonymous message transport backend: persist and stream opaque (client-encrypted) messages. The server is zero-knowledge: it never decrypts payloads, manages keys, or authenticates users. Identity is stateless (UUIDs only).

---

## Module Responsibilities

### app/model

**Role:** Shared domain types used by HTTP, WebSocket, and storage.

- **Message** — Canonical message entity: `id`, `sender_id`, `recipient_id`, `payload` (text), `created_at`. Server treats `payload` as opaque.

No business logic; only data structures.

---

### app/config

**Role:** Load and expose application configuration.

- Load from environment or config file (e.g. server port, database DSN, WebSocket path).
- Expose a single config struct consumed by `cmd/server` when wiring dependencies.

No HTTP, no storage; only configuration.

---

### app/storage

**Role:** Persistence and retrieval of messages.

- **Repository** interface: create message, list messages for a dialog (by `sender_id` and `recipient_id`; returns all messages in both directions).
- Implementation: database or in-memory store. All SQL/queries live here.
- Uses `app/model` for the Message type.

Handlers and WebSocket code must not contain SQL or storage details; they call the repository only.

---

### app/http

**Role:** HTTP API and routing.

- **Router** — Registers routes and delegates to handlers.
- **Handlers:**
  - Create message: accept JSON (sender_id, recipient_id, payload), generate id and created_at, call storage, return created message or error.
  - Read dialog: query params or path for two participant IDs, return all messages between them (both directions).
- Handlers receive a storage interface (e.g. repository) via constructor or options; no direct DB access.
- Request/response types can live here or in `app/model`; no SQL.

---

### app/ws

**Role:** WebSocket hub and real-time delivery.

- **Hub** — Manages connected clients (e.g. keyed by user/dialog or connection ID).
- **Client** — Per-connection state; reads frames, optionally parses subscribe/unsubscribe for dialogs.
- When a new message is persisted (e.g. via storage callback or channel from HTTP), hub broadcasts to relevant clients (e.g. sender and recipient).
- Uses `app/model` for Message when pushing to clients.

No HTTP routing and no SQL; only WebSocket lifecycle and broadcasting.

---

### cmd/server

**Role:** Composition and process lifecycle.

- Load **config**.
- Create **storage** (e.g. repository implementation).
- Create **ws** hub and optionally a way for HTTP layer to notify “new message” (e.g. channel or callback).
- Create **http** router with handlers, injecting storage and (if needed) the “new message” notifier.
- Start HTTP server; serve WebSocket on configured path (e.g. `/ws`).
- No business logic; only wiring and `ListenAndServe`.

---

### api/openapi.yaml

**Role:** Contract for all HTTP and WebSocket endpoints.

- **POST /messages** — Create message (body: sender_id, recipient_id, payload); response: full message with id and created_at.
- **GET /messages** — List messages for a dialog (query: sender_id, recipient_id; semantics: “messages between these two”); response: array of messages.
- **WebSocket** — Document URL (e.g. `/ws`), optional query/path params, and message format (e.g. JSON envelope for server-sent new messages).

Updated whenever endpoints, parameters, or responses change.

---

## Data Flow

1. **Create message:** Client → HTTP handler → storage.Create → (optional) notify WS hub → response.
2. **Read dialog:** Client → HTTP handler → storage.GetDialog(sender_id, recipient_id) → response.
3. **Stream new messages:** Client connects to WS → hub registers client → on new message from storage/HTTP, hub broadcasts to relevant clients.

---

## Boundaries

- **Handlers** do not contain SQL; they use the storage interface.
- **Config** is the only place that knows about env/files for settings.
- **Storage** is the only place that knows about DB or in-memory implementation.
- **WS** does not perform persistence; it only subscribes to “new message” events and pushes to clients.
- **Replaceable backend:** HTTP and WS depend on interfaces (e.g. repository, notifier), so the persistence or transport can be swapped without changing crypto or client contract.
