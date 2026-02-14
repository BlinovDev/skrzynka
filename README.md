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
