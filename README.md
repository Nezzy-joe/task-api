# 🚀 Task API

A simple RESTful Task Management API built with **Go (Golang)**.

This project demonstrates how to build a CRUD (Create, Read, Update, Delete) API using Go's standard `net/http` package, JSON handling, and Swagger documentation without relying on a web framework.

---

## 📌 Features

- ✅ Create a task
- ✅ Retrieve all tasks
- ✅ Retrieve a task by ID
- ✅ Update an existing task
- ✅ Delete a task
- ✅ JSON request & response handling
- ✅ Input validation
- ✅ Proper HTTP status codes
- ✅ Interactive Swagger API documentation

---

## 🛠 Technologies Used

- Go (Golang)
- net/http
- encoding/json
- Swaggo / Swagger UI
- Git & GitHub

---

## 📂 Project Structure

```text
task-api/
│
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
│
├── handlers/
│   └── task.go
│
├── models/
│   └── task.go
│
├── main.go
├── go.mod
├── go.sum
├── README.md
└── .gitignore
```

---

## 📡 API Endpoints

| Method | Endpoint | Description |
|---------|----------|-------------|
| GET | `/tasks` | Retrieve all tasks |
| GET | `/tasks/{id}` | Retrieve a task by ID |
| POST | `/tasks` | Create a new task |
| PUT | `/tasks/{id}` | Update an existing task |
| DELETE | `/tasks/{id}` | Delete a task |

---

## ▶️ Getting Started

### Clone the repository

```bash
git clone https://github.com/Nezzy-joe/task-api.git
```

### Navigate into the project

```bash
cd task-api
```

### Install dependencies

```bash
go mod tidy
```

### Start the server

```bash
go run main.go
```

The API will run on:

```
http://localhost:8080
```

---

## 📖 Swagger Documentation

After starting the server, open:

```
http://localhost:8080/swagger/index.html
```

Swagger UI provides interactive API documentation where you can test every endpoint directly from your browser.

---

## 🧪 Example Request

```bash
curl http://localhost:8080/tasks
```

Example Response

```json
[
  {
    "id": 1,
    "title": "Learn Go",
    "completed": false
  }
]
```

---

## 🎯 Learning Outcomes

This project demonstrates:

- RESTful API development
- CRUD operations
- HTTP routing with Go
- JSON encoding & decoding
- URL path parameters
- Request validation
- Error handling
- Swagger/OpenAPI documentation
- Git version control
- GitHub workflow

---

## 🔮 Future Improvements

- Persistent database (PostgreSQL)
- JWT Authentication
- User accounts
- Pagination
- Filtering & search
- Docker support
- Automated tests

---

## 👨‍💻 Author

**Joseph Amos**

GitHub: https://github.com/Nezzy-joe

---

## 📄 License

This project was developed as part of the **FlyRank Backend Engineering Internship**.