# AI Dev Rules — Go Backend (for Cursor)

## 1. Public Method Comments
- Every exported (public) type, function, and method MUST have a GoDoc comment.
- Comment must start with the identifier name.
- Describe WHAT it does, not HOW.

## 2. No Extra Comments
- Do NOT add inline/explanatory comments inside code.
- Only GoDoc comments for exported identifiers are allowed.
- No TODO/FIXME unless explicitly requested.

## 3. Minimal Changes
- Implement tasks with the smallest possible diff.
- Reuse existing code and structure.
- Avoid refactors, renames, or large rewrites unless strictly necessary.

## 4. Modular Architecture
Each responsibility MUST live in its own package and encapsulate its logic:

- app/http       → handlers + router
- app/ws         → websocket hub & clients
- app/storage    → database + repositories
- app/config     → configuration loading

Rules:
- Do not mix responsibilities between modules.
- Handlers must not contain SQL.
- cmd/server only wires dependencies.

## 5. Swagger Required
- All HTTP endpoints MUST be described in `api/openapi.yaml`.
- Update Swagger when adding or changing endpoints, params, or responses.
- Include example payloads.
- WebSocket endpoint must also be documented (URL + payload format).

## Cursor Execution Rules
1. Modify only the relevant module.
2. Keep changes minimal.
3. Add GoDoc comments for any new public API.
4. Update Swagger if HTTP API changes.
5. Do not introduce new abstractions unless required.
