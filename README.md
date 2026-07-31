# Chirpy

Chirpy is a small social-network backend written in Go. It provides a JSON API for user accounts, JWT authentication, refresh tokens, and short posts called chirps, backed by PostgreSQL.

## Features

- Register and update user accounts
- Log in with Argon2id password hashing
- Authenticate with one-hour JWT access tokens
- Refresh and revoke long-lived refresh tokens
- Create, list, filter, fetch, and delete chirps
- Limit chirps to 140 characters and replace blocked words
- Handle Chirpy Red upgrade webhooks
- Serve a basic web page and request metrics

## Requirements

- [Go 1.26.5](https://go.dev/dl/) or a compatible newer version
- [PostgreSQL](https://www.postgresql.org/)
- [Goose](https://github.com/pressly/goose) for database migrations
- `make` (optional)

Install Goose if needed:

```sh
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Makefile:
```makefile
server:
	go build -o out && ./out

migrate-up:
	goose postgres -dir sql/schema postgres://user:password@localhost:5432/chirpy up

migrate-down:
	goose postgres -dir sql/schema postgres://user:password@localhost:5432/chirpy down
```

## Getting started

1. Clone the repository and enter it:

   ```sh
   git clone https://github.com/bePramudya/chirpy.git
   cd chirpy
   ```

2. Create the PostgreSQL database:

   postgre#
   ```sql
   CREATE DATABASE chirpy;
   ```
   ```sh
   \c chirpy;
   ```

   chirpy#
   ```sql
   ALTER USER postgres WITH PASSWORD 'postgres';
   ```
   
3. Create a `.env` file:

   ```dotenv
   DB_URL=postgres://user:password@localhost:5432/chirpy?sslmode=disable
   PLATFORM=dev
   JWT_SECRET=replace-with-a-random-secret
   POLKA_KEY=f271c81ff7084ee5b99a5091b42d486e
   ```

   `PLATFORM=dev` enables the database reset endpoint. Use another value outside local development.

4. Apply the migrations:

   ```sh
   make migrate-up
   ```

5. Start the server:

   ```sh
   make server
   ```

Chirpy listens on `http://localhost:8080`. Check it with:

```sh
curl http://localhost:8080/api/healthz
```

## API

JSON errors use the form `{"error":"message"}`.

### Users and authentication

| Method | Path | Authentication | Description |
| --- | --- | --- | --- |
| `POST` | `/api/users` | None | Register with `email` and `password` |
| `POST` | `/api/login` | None | Log in and receive `token` and `refresh_token` |
| `PUT` | `/api/users` | Access token | Update the current user's `email` and `password` |
| `POST` | `/api/refresh` | Refresh token | Create a new access token |
| `POST` | `/api/revoke` | Refresh token | Revoke a refresh token |

Pass access and refresh tokens as bearer tokens:

```text
Authorization: Bearer <token>
```

### Chirps

| Method | Path | Authentication | Description |
| --- | --- | --- | --- |
| `GET` | `/api/chirps` | None | List chirps |
| `GET` | `/api/chirps/{chirpID}` | None | Get one chirp |
| `POST` | `/api/chirps` | Access token | Create a chirp from `{"body":"..."}` |
| `DELETE` | `/api/chirps/{chirpID}` | Access token | Delete a chirp owned by the current user |

`GET /api/chirps` accepts these optional query parameters:

- `author_id=<uuid>` filters by author.
- `sort=asc|desc` controls creation-time order; the default is `asc`.

Chirp bodies may contain at most 140 characters. The words `kerfuffle`, `sharbert`, and `fornax` are replaced with `****`, case-insensitively.

### Operations

| Method | Path | Authentication | Description |
| --- | --- | --- | --- |
| `GET` | `/api/healthz` | None | Readiness check |
| `GET` | `/admin/metrics` | None | View static-file request metrics |
| `POST` | `/admin/reset` | None | Reset metrics and users in the `dev` platform |
| `POST` | `/api/polka/webhooks` | Polka API key | Process a Chirpy Red upgrade |
| `GET` | `/app/` | None | Serve the basic web page |

Polka webhooks use a separate authorization scheme:

```text
Authorization: ApiKey <POLKA_KEY>
```

The supported payload is:

```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "<uuid>"
  }
}
```

## Example

Register a user:

```sh
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"bird@example.com","password":"secret"}'
```

Log in:

```sh
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"bird@example.com","password":"secret"}'
```

Use the returned access token to create a chirp:

```sh
curl -X POST http://localhost:8080/api/chirps \
  -H "Authorization: Bearer <access-token>" \
  -H "Content-Type: application/json" \
  -d '{"body":"Hello from Chirpy!"}'
```

## Development

Run the tests:

```sh
go test ./...
```

Database queries live in `sql/queries`, migrations in `sql/schema`, and generated [sqlc](https://sqlc.dev/) code in `internal/database`.
