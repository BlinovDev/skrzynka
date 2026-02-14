# AI Architecture — Go Backend

## Goal

Provide a minimal message transport service for an anonymous encrypted chat system.

The backend must never:
- decrypt messages
- generate encryption keys
- manage user accounts
- enforce access permissions

Messages are assumed to be already encrypted.

---

## Core Principles

1. Zero Knowledge Server
Backend treats payload as opaque JSON.

2. Replaceable Layer
Frontend must be able to switch backend implementations without crypto changes.

3. Stateless Identity
User identity is just UUID generation.

---

## Data Model

messages:
- id
- sender_id (uuid)
- recipient_id (uuid)
- payload (text)
- created_at

Server never interprets payload structure.

---

## Responsibilities

- persist messages
- stream new messages via WebSocket

---

## API

- create new message(all field are required except id and created_at which are iterated automatically)
- read messages for dialog between two(by sender_id and recipient_id) should return all messages from you and to you

---

## Non Goals

- authentication
- rate limiting (can be added later)
- message validation
- key storage
