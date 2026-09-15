# 🚀 Task API

A lightweight RESTful Task Management API built with Go, PostgreSQL, and Docker Compose.

This project started as an in-memory CRUD API and was progressively extended through
multiple storage implementations:

```text
A1 → In-memory storage
A2 → SQLite persistence
A3 → PostgreSQL + Docker Compose
```

For A3 of the FlyRank Backend Engineering Track, the task API was migrated from
SQLite to PostgreSQL and the complete application stack was containerized with
Docker Compose.

The API endpoints and request/response behavior remain consistent while the
storage layer and runtime infrastructure changed.

---

## 🏗️ Architecture

### Application Architecture

```text
Client
  │
  ▼
Go Task API
  │
  ▼
PostgreSQL
```

### Docker Compose Architecture

```text
                    Docker Compose
                         │
              ┌──────────┴──────────┐
              │                     │
              ▼                     ▼
       API Container        PostgreSQL Container
       Go application                db
              │                     │
              └────── Compose ──────┘
                     network
                         │
                         ▼
                  taskdata volume
```

The API connects to PostgreSQL using the Docker Compose service name:

```text
db
```

When running inside Docker Compose, the API uses a PostgreSQL connection string
with `db` as the database host:

```text
postgres://postgres:<password>@db:5432/tasks?sslmode=disable
```

The API does not use `localhost` to reach PostgreSQL from inside the API container.

---

## ✨ Features

- Create tasks
- Retrieve all tasks
- Retrieve a task by ID
- Update an existing task
- Delete an existing task
- PostgreSQL persistence
- Automatic database connection
- Automatic `tasks` table creation
- Three seed tasks when the table is empty
- Seed-once behavior
- Parameterized SQL queries
- Input validation
- JSON error responses
- Correct HTTP status codes
- Swagger / OpenAPI documentation
- Dockerized Go API
- Dockerized PostgreSQL
- Docker Compose orchestration
- Persistent PostgreSQL volume

---

## 🛠️ Technologies Used

- Go
- `net/http`
- `encoding/json`
- `database/sql`
- PostgreSQL
- `github.com/lib/pq`
- `github.com/joho/godotenv`
- Docker
- Docker Compose
- Swagger / OpenAPI
- Git
- GitHub
- `curl`
- `psql`

---

## 📂 Project Structure

```text
task-api/
│
├── database/
│   └── database.go
│
├── handlers/
│   └── task.go
│
├── models/
│   └── task.go
│
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
│
├── images/
│   ├── db-browser.png
│   ├── postgres.png
│   └── swagger-ui.png
│
├── main.go
├── Dockerfile
├── compose.yaml
├── .dockerignore
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

> `images/postgres.png` should contain the PostgreSQL database screenshot required
> for the A3 submission.

The `.env` file is intentionally excluded from the repository because it contains
the local database password.

---

## 🔐 Environment Variables

Create a `.env` file in the project root.

Use `.env.example` as the template.

Example:

```env
POSTGRES_PASSWORD=dev
DATABASE_URL=postgres://postgres:dev@localhost:5432/tasks?sslmode=disable
```

### Variables

| Variable | Purpose |
|---|---|
| `POSTGRES_PASSWORD` | Password used by the PostgreSQL container |
| `DATABASE_URL` | PostgreSQL connection string used by the application |

### Important

`.env` is ignored by Git.

`.env.example` is committed with placeholder values so another developer can create
their own local configuration.

Never commit real database credentials to the repository.

---

## ▶️ Running the Application

### Prerequisites

Install:

- Docker
- Docker Compose

Verify Docker:

```bash
docker --version
```

Verify Docker Compose:

```bash
docker compose version
```

---

## 🚀 Start the Whole Stack

The complete application stack can be started with one command:

```bash
docker compose up --build
```

Docker Compose starts:

1. The Go API container
2. The PostgreSQL container
3. The Docker network connecting the services
4. The persistent `taskdata` volume

The API waits for PostgreSQL to become healthy before starting.

Once running, the API is available at:

```text
http://localhost:8080
```

---

## 🛑 Stop the Stack

Stop the running containers with:

```bash
docker compose down
```

The `taskdata` named volume is preserved when the containers are removed, so the
PostgreSQL data remains available when the stack is started again.

---

## 🗄️ Database

The application uses PostgreSQL.

The database name is:

```text
tasks
```

The PostgreSQL Compose service is:

```text
db
```

The application automatically creates the `tasks` table if it does not already exist.

### Tasks Table

The schema is:

```sql
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    done BOOLEAN NOT NULL DEFAULT FALSE
);
```

### Columns

| Column | Description |
|---|---|
| `id` | Auto-generated primary key |
| `title` | Task title |
| `done` | Completion status |

---

## 🌱 Seed Data

When the `tasks` table is empty, the application inserts three example tasks:

```text
Learn Go
Build Task API
Connect API to PostgreSQL
```

The seed operation only runs when the table is empty.

Restarting the application does not create duplicate seed records.

---

## 💾 Database Persistence

PostgreSQL data is stored in the Docker named volume:

```text
taskdata
```

The volume is mounted at:

```text
/var/lib/postgresql/data
```

This allows PostgreSQL data to survive container recreation.

### Persistence Verification

Persistence was verified by:

1. Starting the application with Docker Compose.
2. Creating a new task through the API.
3. Confirming the task was returned by `GET /tasks`.
4. Stopping the Compose stack.
5. Starting the Compose stack again.
6. Calling `GET /tasks`.
7. Confirming the previously created task was still present.

For example, the following task was created successfully:

```json
{
  "id": 5,
  "title": "Docker Compose",
  "completed": false
}
```

After restarting the Compose stack, the task remained available.

---

## 📡 API Endpoints

| Method | Endpoint | Description | Success |
|---|---|---|---|
| GET | `/tasks` | Retrieve all tasks | `200 OK` |
| GET | `/tasks/{id}` | Retrieve one task | `200 OK` |
| POST | `/tasks` | Create a task | `201 Created` |
| PUT | `/tasks/{id}` | Update a task | `200 OK` |
| DELETE | `/tasks/{id}` | Delete a task | `204 No Content` |

Unknown task IDs return:

```text
404 Not Found
```

Invalid request bodies return:

```text
400 Bad Request
```

with a JSON error response.

---

## 🧪 API Testing with curl

### Get all tasks

```bash
curl -i http://localhost:8080/tasks
```

Example:

```text
HTTP/1.1 200 OK
Content-Type: application/json

[{"id":2,"title":"Build Task API","completed":false},{"id":3,"title":"Connect API to PostgreSQL","completed":false},{"id":1,"title":"Learn Go","completed":false}]
```

### Get a task by ID

```bash
curl -i http://localhost:8080/tasks/1
```

### Create a task

```bash
curl -i -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Docker Compose","completed":false}'
```

Example response:

```text
HTTP/1.1 201 Created
Content-Type: application/json

{"id":5,"title":"Docker Compose","completed":false}
```

### Update a task

```bash
curl -i -X PUT http://localhost:8080/tasks/5 \
  -H "Content-Type: application/json" \
  -d '{"title":"Docker Compose deeply","completed":true}'
```

Expected:

```text
HTTP/1.1 200 OK
```

### Delete a task

```bash
curl -i -X DELETE http://localhost:8080/tasks/5
```

Expected:

```text
HTTP/1.1 204 No Content
```

### Request an unknown task

```bash
curl -i http://localhost:8080/tasks/999
```

Expected:

```text
HTTP/1.1 404 Not Found
```

with a JSON error message.

---

## 🔒 Parameterized Queries

The PostgreSQL implementation uses parameterized SQL queries.

For example:

```sql
SELECT id, title, done
FROM tasks
WHERE id = $1;
```

The task ID is supplied separately to the database driver instead of being
concatenated directly into the SQL statement.

This keeps user-controlled values separate from the SQL statement and helps prevent
SQL injection.

Parameterized queries are also used for inserts, updates, and deletes.

---

## 🔄 PostgreSQL CRUD Operations

### Read all tasks

```sql
SELECT id, title, done
FROM tasks;
```

### Read a task by ID

```sql
SELECT id, title, done
FROM tasks
WHERE id = $1;
```

### Create a task

```sql
INSERT INTO tasks (title, done)
VALUES ($1, $2)
RETURNING id;
```

The `RETURNING` clause provides the newly generated task ID.

### Update a task

```sql
UPDATE tasks
SET title = $1, done = $2
WHERE id = $3;
```

### Delete a task

```sql
DELETE FROM tasks
WHERE id = $1;
```

---

## 🐳 Docker

The project includes a Dockerfile for building the Go API image.

### Dockerfile responsibilities

The Dockerfile:

1. Uses a Go Alpine base image.
2. Sets `/app` as the working directory.
3. Copies Go module files.
4. Downloads dependencies.
5. Copies the application source.
6. Builds the Go API binary.
7. Exposes port `8080`.
8. Starts the compiled API binary.

Build the API image manually with:

```bash
docker build -t task-api .
```

The normal workflow is to use Docker Compose:

```bash
docker compose up --build
```

---

## 🧩 Docker Compose

The `compose.yaml` file defines two services:

```text
api
db
```

### API service

The API service:

- Builds from the project Dockerfile
- Exposes port `8080`
- Receives `DATABASE_URL` from the environment
- Depends on PostgreSQL health status

### Database service

The PostgreSQL service:

- Uses the PostgreSQL image
- Creates the `tasks` database
- Uses `POSTGRES_PASSWORD`
- Exposes PostgreSQL on port `5432`
- Stores data in the `taskdata` volume
- Uses `pg_isready` as a health check

---

## 🌐 Container Networking

Inside Docker Compose, the API does not connect to PostgreSQL using:

```text
localhost
```

Instead it connects to:

```text
db
```

because `db` is the PostgreSQL service name on the Compose network.

The API therefore uses a connection string shaped like:

```text
postgres://postgres:<password>@db:5432/tasks?sslmode=disable
```

This allows the API and PostgreSQL containers to communicate through Docker's
internal network.

---

## 🩺 PostgreSQL Health Check

The PostgreSQL service uses:

```bash
pg_isready -U postgres -d tasks
```

Docker Compose waits for the database service to become healthy before starting
the API service.

This reduces startup failures caused by the API attempting to connect before
PostgreSQL is ready.

---

## 🖥️ Inspecting PostgreSQL

The database can be inspected directly from the PostgreSQL container.

Open a PostgreSQL shell:

```bash
docker compose exec db psql -U postgres -d tasks
```

List tables:

```sql
\dt
```

View task records:

```sql
SELECT * FROM tasks;
```

Exit PostgreSQL:

```sql
\q
```

---

## 📸 Database Evidence

The A3 assignment requires a screenshot showing the data in the PostgreSQL
database.

The screenshot is stored in:

```text
images/postgres.png
```

The screenshot should show PostgreSQL being inspected with commands such as:

```sql
\dt
```

and:

```sql
SELECT * FROM tasks;
```

Example Markdown reference:

```markdown
![PostgreSQL Database](images/postgres.png)
```

![PostgreSQL Database](images/postgres.png)

---

## 📖 Swagger Documentation

The project includes Swagger / OpenAPI documentation.

Start the application:

```bash
docker compose up
```

Then open:

```text
http://localhost:8080/swagger/index.html
```

Swagger UI provides interactive documentation for the available API endpoints.

![Swagger UI](images/swagger-ui.png)

---

## 🧠 Storage Evolution

This project demonstrates that the API contract can remain stable while the storage
implementation changes.

### A1 — In-memory

```text
Client
  ↓
Go API
  ↓
In-memory task list
```

Data disappeared when the program stopped.

### A2 — SQLite

```text
Client
  ↓
Go API
  ↓
SQLite
  ↓
tasks.db
```

Data persisted in a local database file.

### A3 — PostgreSQL + Docker

```text
Client
  ↓
Go API container
  ↓
PostgreSQL container
  ↓
taskdata volume
```

The database now runs as its own server inside a container and persists through
a Docker volume.

The API routes remain the same while the storage engine changes.

---

## 🔁 Storage as an Implementation Detail

The API contract remains unchanged even though the storage implementation changed
from memory to SQLite and then PostgreSQL.

The client only interacts with the API endpoints.

The underlying storage engine is therefore an implementation detail of the
application rather than something the client needs to know about.

---

## ✅ A3 Verification

The following parts of the containerized stack were tested:

- Docker installation
- PostgreSQL container
- PostgreSQL health check
- Database connection
- Automatic table creation
- Seed data
- `GET /tasks`
- `GET /tasks/{id}`
- `POST /tasks`
- `PUT /tasks/{id}`
- `DELETE /tasks/{id}`
- Unknown task `404`
- PostgreSQL persistence
- Docker Compose startup
- API-to-database communication through service name `db`
- Persistent Docker volume

A successful containerized startup is:

```bash
docker compose up --build
```

The API is then available at:

```text
http://localhost:8080
```

---

## 🧪 Clean Clone

A fresh clone should be able to start the application without manual database
setup.

Clone the repository:

```bash
git clone https://github.com/Nezzy-joe/task-api.git
```

Enter the project:

```bash
cd task-api
```

Create the local environment file:

```bash
cp .env.example .env
```

Edit `.env` and set a real local PostgreSQL password.

Then start the complete stack:

```bash
docker compose up
```

The database and API should start automatically.

The seeded tasks should then be available through:

```bash
curl -i http://localhost:8080/tasks
```

No manual PostgreSQL table creation is required.

---

## 🧰 Useful Docker Commands

### Start the stack

```bash
docker compose up --build
```

### Start in detached mode

```bash
docker compose up -d --build
```

### View Compose services

```bash
docker compose ps
```

### View all logs

```bash
docker compose logs
```

### View API logs

```bash
docker compose logs api
```

### View database logs

```bash
docker compose logs db
```

### Stop the stack

```bash
docker compose down
```

### Enter PostgreSQL

```bash
docker compose exec db psql -U postgres -d tasks
```

---

## 🔐 Security Notes

The project keeps the database password outside the application source code.

The real `.env` file is ignored by Git.

Only `.env.example` with placeholder values is committed.

Database queries use parameterized values rather than directly concatenating
user input into SQL strings.

For a production deployment, additional controls would still be required,
including stronger secret management, authentication, authorization, TLS,
database access restrictions, logging, monitoring, and other environment-specific
security controls.

---

## 📚 Learning Outcomes

This project demonstrates practical experience with:

- REST API development in Go
- HTTP routing
- JSON request and response handling
- CRUD operations
- PostgreSQL
- SQL
- Parameterized queries
- Database initialization
- Database seeding
- Environment-based configuration
- Docker images
- Docker containers
- Docker Compose
- Container networking
- Docker volumes
- Database persistence
- Health checks
- API testing with `curl`
- PostgreSQL inspection with `psql`
- Git and GitHub workflow
- Swagger / OpenAPI documentation

---

## 🔮 Future Improvements

Potential future improvements include:

- Automated integration tests
- Authentication and authorization
- JWT-based sessions
- User accounts
- Pagination
- Filtering and search
- Request logging
- Structured logging
- Database migrations
- Multi-stage Docker builds
- Production secret management
- Dedicated health endpoint with database connectivity checks

---

## 👨‍💻 Author

**Joseph Amos**

GitHub:

https://github.com/Nezzy-joe

---

## 📄 License

This project was developed as part of the FlyRank Backend Engineering Internship.