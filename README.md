# Books Backend

## Overview
Books Backend is a REST API built using Golang with the Gin framework. It provides functionalities to manage a collection of books, including creating, retrieving, updating, and deleting books. The project uses PostgreSQL as the database, Redis for caching, and Kafka for event streaming.

## Features
- CRUD operations for books
- PostgreSQL as the primary database
- Redis caching for optimized data retrieval
- Kafka for event-driven messaging
- Swagger documentation

## Technologies Used
- **Golang** (Gin framework)
- **PostgreSQL** (Database)
- **Redis** (Caching layer)
- **Kafka** (Event streaming)
- **Docker** (For containerization)

## API Endpoints

### Books API
| Method | Endpoint        | Description |
|--------|---------------|-------------|
| GET    | `/books`       | Get all books with pagination |
| GET    | `/books/:id`   | Get book by ID |
| POST   | `/books`       | Create a new book |
| PUT    | `/books/:id`   | Update an existing book |
| DELETE | `/books/:id`   | Delete a book |

## Prerequisites
Ensure you have the following installed:
- Golang
- PostgreSQL
- Redis
- Kafka (Zookeeper required)
- Docker (optional for containerized setup)

## Environment Variables
Create a `.env` file in the project root with the following configuration:
```env
DB_HOST=localhost
DB_USER=rohan
DB_PASSWORD=password1
DB_NAME=booksdb
DB_PORT=5432
REDIS_ADDR=localhost:6379
KAFKA_BROKER=localhost:9092
```

## Setup and Run Locally

### 1. Clone the repository
```bash
git clone https://github.com/rohans540/books-backend.git
cd books-backend
```

### 2. Install dependencies
```bash
go mod tidy
```

### 3. Start PostgreSQL, Redis, and Kafka
Ensure PostgreSQL, Redis, and Kafka services are running.

To start Redis:
```bash
redis-server
```

To start Kafka (with Docker Compose):
```bash
docker-compose up -d
```

### 4. Run the application
```bash
go run main.go
```

### 5. Access the API
- Swagger Documentation: `http://localhost:8000/swagger/index.html`
- API Base URL: `http://localhost:8000/books`

### 6. Running with Docker
To run the project using Docker:
```bash
docker build -t books-backend .
docker run -p 8000:8000 --env-file .env books-backend
```

## Logs and Debugging
- Check PostgreSQL logs: `sudo journalctl -u postgresql --no-pager`
- Check Redis logs: `redis-cli monitor`
- Check Kafka logs: `docker logs kafka`
- Application logs will be printed in the terminal.

