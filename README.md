# Skrzynka project

This service is a minimal stateless message transport layer.

It does NOT:
- authenticate users
- encrypt/decrypt data
- understand message contents

All cryptography happens on the client side.

The backend simply:
- issues anonymous user IDs
- stores encrypted payloads
- delivers messages via HTTP + WebSocket

---

## Philosophy

The backend is intentionally "dumb".

Security comes from client-side encryption, not server logic.
The server assumes all stored messages are opaque ciphertext blobs.

This allows:
- replacing frontend crypto anytime
- replacing backend implementation anytime
- zero trust towards backend

---

## How to run the application

### Prerequisites

- **Go 1.22+**

### Quick start

1. **Clone and enter the repo** (if not already):

   ```bash
   cd anon-skrzynka
   ```

2. **Install dependencies** (optional; `go run` will fetch them):

   ```bash
   go mod download
   ```

3. **Start the server**:

   ```bash
   go run ./cmd/server
   ```

   You should see: `listening on :8080`. The API is available at **http://localhost:8080**.

### Configuration

The server currently uses defaults:

| Setting   | Default | Description        |
|----------|---------|--------------------|
| HTTP port| `8080`  | API and WebSocket  |
| WebSocket path | `/ws` | Upgrade path for WS |

Configuration is loaded in `app/config`; it can be extended later (e.g. environment variables or a config file) without changing how you run the binary.

### Running as a binary

To build and run a standalone binary:

```bash
go build -o skrzynka ./cmd/server
./skrzynka
```

### Swagger UI (API docs)

The OpenAPI spec lives in **`api/openapi.yaml`**. The server serves it and an interactive Swagger UI:

- **Spec (raw YAML):** [http://localhost:8080/openapi.yaml](http://localhost:8080/openapi.yaml)
- **Swagger UI:** [http://localhost:8080/docs](http://localhost:8080/docs)

With the server running, open **http://localhost:8080/docs** in a browser to explore and try the API. The UI loads the spec from `/openapi.yaml`; edit `api/openapi.yaml` and restart the server to see changes.

---

## Database

**Current behavior:** The application uses **in-memory storage** only. No database setup is required to run the server. Data is lost when the process stops.

**If you add a persistent database later:**

1. **Create the database** (example for PostgreSQL):

   ```bash
   createdb skrzynka
   # or in psql: CREATE DATABASE skrzynka;
   ```

2. **Run migrations** (or initial schema), if your project includes them (e.g. in a `migrations/` directory or via a migration tool).

3. **Point the app at the DB** by setting the connection string (e.g. `DATABASE_DSN` or similar) in the config. The `app/storage` package would then use a DB implementation of `Repository` instead of `storage.NewMemory()`.

Until then, no database creation or migrations are needed.

---

## Testing with Postman

### Base URL

- **Base URL:** `http://localhost:8080`  
  Ensure the server is running (`go run ./cmd/server`) before sending requests.

### Import the API spec (recommended)

1. In Postman: **Import** → **Link** or **File**.
2. Use the OpenAPI definition: **`api/openapi.yaml`** (or the full path to it).  
   Postman will create a collection with the HTTP endpoints and example bodies.

### Manual request reference

#### 1. Create a message — `POST /messages`

- **Method:** POST  
- **URL:** `http://localhost:8080/messages`  
- **Headers:** `Content-Type: application/json`  
- **Body (raw JSON):**

  ```json
  {
    "sender_id": "550e8400-e29b-41d4-a716-446655440000",
    "recipient_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "payload": "<opaque encrypted string>"
  }
  ```

- **Expected:** `201 Created` with the created message (includes `id` and `created_at`).
- **Tips:** Use valid UUIDs for `sender_id` and `recipient_id`. `payload` can be any string (e.g. base64 ciphertext).

#### 2. Get dialog — `GET /messages`

- **Method:** GET  
- **URL:** `http://localhost:8080/messages?sender_id=550e8400-e29b-41d4-a716-446655440000&recipient_id=6ba7b810-9dad-11d1-80b4-00c04fd430c8`  
- **Query params:** `sender_id`, `recipient_id` (both required, UUIDs).  
- **Expected:** `200 OK` with a JSON array of messages between those two participants (both directions).

### Suggested test flow

1. **POST /messages** — Create one or two messages between two UUIDs (e.g. Alice and Bob).
2. **GET /messages** — Same `sender_id` and `recipient_id` (order of the two doesn’t matter for the dialog).  
   You should see the messages created in step 1.
3. **WebSocket (optional)** — In Postman, create a **WebSocket request** to  
   `ws://localhost:8080/ws?user_id=<uuid>`  
   Use one of the participant UUIDs. When you create a new message via POST with that user as sender or recipient, the server can push the new message over the WebSocket (if the client is subscribed).

### Postman tips

- **Environment:** Create a Postman environment with variable `base_url` = `http://localhost:8080` and use `{{base_url}}/messages`, `{{base_url}}/ws` in requests.
- **Collection variables:** Store two UUIDs (e.g. `user_a`, `user_b`) and reuse them in POST body and GET query params.
- **Tests tab:** Add scripts to assert status code (e.g. 201 for POST, 200 for GET) and that the response body contains expected fields (`id`, `created_at` for POST; array for GET).
- **WebSocket:** For `/ws`, use the **WebSocket** request type; after connection, sending a new message via a separate POST request may trigger a broadcast that appears in the WebSocket response panel.
