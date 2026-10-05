# Task API

A secure Task Management API built with **Go**, **PostgreSQL**, and **Supabase Auth**.

This project extends a RESTful Task API with authentication, JWT-based access control, reusable authentication middleware, protected routes, logout, and interactive Swagger/OpenAPI documentation.

The project was built incrementally in stages, with each major stage represented by a Git commit.

---

## Features

- RESTful Task API
- User signup with Supabase Auth
- User login with Supabase Auth
- JWT access-token authentication
- Refresh-token support through Supabase Auth
- Bearer-token protected routes
- Reusable authentication middleware
- Protected user profile endpoint
- Protected dashboard endpoint
- Authenticated logout
- Public API endpoint
- Swagger/OpenAPI documentation
- Environment-based configuration
- PostgreSQL persistence
- Docker / Docker Compose support
- Git-based staged development

---

## Architecture

```text
                        Client
                          |
                          | HTTP / JSON
                          v
                  +------------------+
                  |      Go API      |
                  |    net/http      |
                  +--------+---------+
                           |
             +-------------+-------------+
             |                           |
             v                           v
      +-------------+              +-------------+
      | Supabase    |              | PostgreSQL  |
      | Auth        |              |   Database  |
      +------+------+              +-------------+
             |
             | JWT Access Token
             v
      Authentication
          Middleware
             |
             v
       Protected Routes
```

### Authentication Flow

```text
1. Client creates an account
       |
       v
POST /auth/signup
       |
       v
Supabase Auth creates the user

2. Client logs in
       |
       v
POST /auth/login
       |
       v
Supabase Auth validates credentials
       |
       v
Access Token + Refresh Token

3. Client requests a protected resource
       |
       v
Authorization: Bearer <JWT>
       |
       v
Go authentication middleware
       |
       v
Supabase token verification
       |
       v
Protected route
```

Supabase Auth is responsible for account management, password handling, and token issuance. The Go API does not store or hash user passwords itself.

---

## Technology Stack

- **Go**
- **net/http**
- **PostgreSQL**
- **lib/pq**
- **Supabase Auth**
- **github.com/supabase-community/auth-go**
- **Swagger / OpenAPI**
- **swaggo/swag**
- **swaggo/http-swagger**
- **godotenv**
- **Docker / Docker Compose**
- **Git / GitHub**

---

## Project Structure

```text
task-api/
│
├── ai-version/
│
├── database/
│   └── database.go
│
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
│
├── handlers/
│   ├── auth.go
│   ├── middleware.go
│   ├── profile.go
│   └── session.go
│
├── models/
│
├── images/
│
├── config.go
├── main.go
├── compose.yaml
├── Dockerfile
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

---

# Authentication

Supabase is used as the Identity Provider.

The API follows this model:

```text
Client
   |
   | email + password
   v
Supabase Auth
   |
   | JWT
   v
Client
   |
   | Authorization: Bearer <JWT>
   v
Go API
   |
   | Verify token
   v
Supabase Auth
```

The backend does not implement its own password hashing or cryptographic token signing.

---

# API Reference

## Authentication Endpoints

### `POST /auth/signup`

Creates a new user account through Supabase Auth.

#### Request

```json
{
  "email": "test@example.com",
  "password": "Test1234!"
}
```

#### Success

```text
201 Created
```

#### Example

```bash
curl -i -X POST http://localhost:8080/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Test1234!"}'
```

---

### `POST /auth/login`

Authenticates a user and returns a Supabase access token and refresh token.

#### Request

```json
{
  "email": "test@example.com",
  "password": "Test1234!"
}
```

#### Success

```text
200 OK
```

The response includes:

- `access_token`
- `refresh_token`
- `token_type`
- `expires_in`
- `expires_at`
- authenticated user information

#### Example

```bash
curl -i -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Test1234!"}'
```

---

### `POST /auth/logout`

Logs out the authenticated user and revokes the user's refresh-token session through Supabase Auth.

Authentication is required.

#### Header

```http
Authorization: Bearer <ACCESS_TOKEN>
```

#### Success

```text
204 No Content
```

#### Example

```bash
curl -i -X POST http://localhost:8080/auth/logout \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

---

# Public Endpoint

### `GET /public/info`

Returns information available to unauthenticated clients.

Authentication is not required.

#### Success

```text
200 OK
```

#### Response

```json
{
  "message": "Welcome stranger! This info is public."
}
```

#### Example

```bash
curl -i http://localhost:8080/public/info
```

---

# Protected Endpoints

All protected endpoints require:

```http
Authorization: Bearer <ACCESS_TOKEN>
```

---

### `GET /protected/profile`

Returns safe metadata for the authenticated user.

#### Success

```text
200 OK
```

#### Response

```json
{
  "id": "user-id",
  "email": "test@example.com",
  "account_created_at": "timestamp"
}
```

#### Example

```bash
curl -i http://localhost:8080/protected/profile \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

---

### `GET /protected/dashboard`

Demonstrates reuse of the authentication middleware on another protected route.

#### Success

```text
200 OK
```

#### Response

```json
{
  "message": "Welcome to your protected dashboard."
}
```

#### Example

```bash
curl -i http://localhost:8080/protected/dashboard \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

---

# API Summary

| Method | Endpoint | Authentication | Description |
|---|---|---|---|
| `POST` | `/auth/signup` | No | Create a new user |
| `POST` | `/auth/login` | No | Authenticate and receive tokens |
| `POST` | `/auth/logout` | Yes | End the authenticated session |
| `GET` | `/public/info` | No | Read public information |
| `GET` | `/protected/profile` | Yes | Return authenticated user metadata |
| `GET` | `/protected/dashboard` | Yes | Example protected route |

---

# Status Codes

| Status Code | Meaning |
|---|---|
| `200` | Request successful |
| `201` | Resource created |
| `204` | Successful logout with no response body |
| `400` | Invalid or missing input |
| `401` | Authentication required or token invalid |
| `405` | HTTP method not allowed |

---

# Authentication Middleware

Protected routes use a reusable authentication middleware.

The middleware:

1. Reads the `Authorization` header.
2. Extracts the `Bearer` token.
3. Rejects missing or malformed tokens.
4. Sends the token to Supabase Auth for verification.
5. Rejects invalid or expired tokens.
6. Attaches authenticated user information to the request context.
7. Allows the protected handler to execute.

The same middleware is reused by:

```text
GET  /protected/profile
GET  /protected/dashboard
POST /auth/logout
```

This avoids duplicating authentication logic across protected endpoints.

---

# Error Handling

## Missing credentials

Example:

```json
{
  "error": "Email and password are required"
}
```

Response:

```text
400 Bad Request
```

---

## Missing access token

Example:

```json
{
  "error": "Access token required"
}
```

Response:

```text
401 Unauthorized
```

---

## Invalid or expired access token

Example:

```json
{
  "error": "Invalid or expired token"
}
```

Response:

```text
401 Unauthorized
```

---

# Environment Configuration

Create a `.env` file in the project root.

Example:

```env
DATABASE_URL=postgres://postgres:your_database_password@localhost:5432/tasks?sslmode=disable
POSTGRES_PASSWORD=your_database_password
SUPABASE_URL=https://your-project-ref.supabase.co
SUPABASE_KEY=your_supabase_anon_key
PORT=8080
```

### Important

Never commit your real `.env` file.

Use `.env.example` as the configuration template.

The `.env` file is excluded through `.gitignore`.

Never place the following in Git:

- Supabase access tokens
- Supabase refresh tokens
- Database passwords
- Supabase service-role keys
- Other private credentials

---

# Supabase Configuration

Create a Supabase project and obtain:

- Project URL
- anon/public API key

The API expects the Supabase project URL in this form:

```text
https://your-project-ref.supabase.co
```

Do not append:

```text
/rest/v1/
```

to `SUPABASE_URL`.

For this practice project, email confirmation should be disabled in the Supabase Email authentication provider so a newly created user can log in immediately.

In a production application, email confirmation may be kept enabled as an additional account-security measure.

---

# PostgreSQL Setup

This project uses PostgreSQL for the Task API data layer.

Create the local database:

```bash
sudo -u postgres createdb tasks
```

Verify it:

```bash
sudo -u postgres psql -tAc "SELECT datname FROM pg_database WHERE datname='tasks';"
```

The expected result is:

```text
tasks
```

The local development database URL should follow this structure:

```env
DATABASE_URL=postgres://postgres:your_database_password@localhost:5432/tasks?sslmode=disable
```

---

# Running the API

## Local Development

Make sure PostgreSQL is running.

Configure your `.env` file.

Then run:

```bash
go run .
```

The API starts at:

```text
http://localhost:8080
```

---

# Swagger / OpenAPI

Swagger UI is available at:

```text
http://localhost:8080/docs/
```

Swagger provides interactive documentation for the API.

The protected routes use the `BearerAuth` security scheme.

Protected routes include:

```text
POST /auth/logout
GET  /protected/profile
GET  /protected/dashboard
```

---

## Authorizing Swagger

1. Open:

```text
http://localhost:8080/docs/
```

2. Click **Authorize**.

3. Select:

```text
BearerAuth
```

4. Enter:

```text
Bearer <ACCESS_TOKEN>
```

5. Click **Authorize**.

6. Execute:

```text
GET /protected/profile
```

or:

```text
GET /protected/dashboard
```

A valid token should return:

```text
200 OK
```

---

# Swagger Screenshot

Add the Stage 5 Swagger screenshot to the repository, for example:

```text
images/swagger-ui.png
```

Then include it here:

![Swagger UI](images/swagger-ui.png)

The screenshot should demonstrate:

- Swagger UI running locally
- Bearer authorization
- Protected endpoint lock icons
- Interactive API documentation

---

# Testing

Run the Go test suite:

```bash
go test ./...
```

Run formatting:

```bash
gofmt -w .
```

Check the Git diff for whitespace errors:

```bash
git diff --check
```

---

# End-to-End Authentication Test

## 1. Signup

```bash
curl -i -X POST http://localhost:8080/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Test1234!"}'
```

Expected:

```text
201 Created
```

---

## 2. Login

```bash
curl -i -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Test1234!"}'
```

Expected:

```text
200 OK
```

The response provides an access token and refresh token.

---

## 3. Protected profile

```bash
curl -i http://localhost:8080/protected/profile \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

Expected:

```text
200 OK
```

---

## 4. Tampered token

Changing a character in the access token should cause verification to fail.

Expected:

```text
401 Unauthorized
```

Response:

```json
{
  "error": "Invalid or expired token"
}
```

---

## 5. Protected dashboard

```bash
curl -i http://localhost:8080/protected/dashboard \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

Expected:

```text
200 OK
```

---

## 6. Logout

```bash
curl -i -X POST http://localhost:8080/auth/logout \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

Expected:

```text
204 No Content
```

---

# Docker

The repository also contains:

```text
Dockerfile
compose.yaml
```

Docker Compose can be used to run the API and PostgreSQL together.

Start the stack with:

```bash
docker compose up --build
```

The API is available at:

```text
http://localhost:8080
```

Swagger is available at:

```text
http://localhost:8080/docs/
```

For the Docker database, the host port is:

```text
5433
```

while the PostgreSQL container continues to use:

```text
5432
```

internally.

Stop the stack with:

```bash
docker compose down
```

Do not use:

```bash
docker compose down -v
```

unless you intentionally want to remove the PostgreSQL volume and its stored data.

---

# Security

This project follows several important authentication practices:

- Passwords are not stored by the Go API.
- Password hashing is delegated to Supabase Auth.
- JWT issuance is delegated to Supabase Auth.
- Protected routes require bearer authentication.
- Access tokens are verified through Supabase.
- Invalid or expired tokens return `401 Unauthorized`.
- Supabase credentials are stored in environment variables.
- `.env` is excluded from Git.
- `.env.example` contains placeholders rather than secrets.
- The Supabase `service_role` key is not used.
- Access tokens should never be logged or committed to source control.

---

# Git Workflow

The project was developed incrementally using Git.

Each major implementation stage was committed separately.

Example:

```bash
git log --oneline
```

Development history:

```text
Stage 1 - Signup and Login
Stage 2 - Public and Protected Routes
Stage 3 - Token Verification
Stage 4 - Authentication Middleware and Logout
Stage 5 - Swagger UI with Bearer Authentication
Stage 6 - Publication and Documentation
```

---

# Development Stages

## Stage 1 — Signup and Login

Implemented:

```text
POST /auth/signup
POST /auth/login
```

The API validates credentials and delegates authentication to Supabase Auth.

---

## Stage 2 — Public and Protected Routes

Implemented:

```text
GET /public/info
GET /protected/profile
```

The protected endpoint initially checked only whether a bearer token was presented.

---

## Stage 3 — Token Verification

Implemented actual token verification through Supabase Auth.

The protected profile endpoint now:

- extracts the bearer token
- sends the token to Supabase
- verifies the authenticated user
- returns safe user metadata
- rejects invalid or expired tokens

---

## Stage 4 — Authentication Middleware and Logout

Introduced reusable authentication middleware.

Added:

```text
GET  /protected/dashboard
POST /auth/logout
```

The same authentication guard is reused across multiple protected routes.

---

## Stage 5 — Swagger Bearer Authentication

Implemented Swagger/OpenAPI documentation with bearer authentication.

Swagger documents:

- authentication routes
- public routes
- protected routes
- bearer authorization
- request and response models

---

## Stage 6 — Publication and Documentation

The final required stage focuses on:

- README documentation
- `.env` security
- `.env.example`
- Git history
- GitHub publication
- final end-to-end verification

---

# AI Rematch

The `ai-version/` directory is reserved for the optional AI rematch stage.

The purpose of the AI rematch is to:

1. Ask an AI to build the same secured API.
2. Generate the AI version separately.
3. Run the generated version.
4. Compare it with the manually built implementation.
5. Review security decisions.
6. Document differences.
7. Improve the prompt.
8. Regenerate the AI version.

The AI-generated implementation should remain separate from the main submission implementation.

---

# Lessons Demonstrated

## Authentication vs Authorization

Authentication answers:

```text
Who is this user?
```

Authorization answers:

```text
What is this authenticated user allowed to access?
```

---

## Bearer Tokens

Authenticated requests use:

```http
Authorization: Bearer <JWT>
```

---

## Middleware

Authentication logic is centralized in reusable middleware instead of being copied into every protected handler.

---

## Identity Provider

Supabase acts as the Identity Provider responsible for:

- user accounts
- password handling
- authentication
- token issuance

---

## Environment Variables

Configuration and secrets are separated from source code through environment variables.

---

## API Status Codes

The project intentionally uses HTTP status codes to communicate request results:

```text
200 OK
201 Created
204 No Content
400 Bad Request
401 Unauthorized
405 Method Not Allowed
```

---

# Project Author

## Joseph Amos Ekpe

Backend / Go Developer

Focused on:

- Go
- Backend Engineering
- REST APIs
- PostgreSQL
- Authentication
- AI-powered applications

GitHub:

```text
https://github.com/Nezzy-joe
```

---

#License

This project was created as part of a backend engineering learning and internship workflow.