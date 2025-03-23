README.md
markdown
Copy
# Go Template Service

This is a simple template service written in Go with a SQLite database. It uses [Echo](https://echo.labstack.com/) as the web framework and includes [Swagger](https://swagger.io/) for API documentation.

## Features

- **RESTful API** for managing documents
- **SQLite** database with auto-migration
- **Swagger UI** for interactive API documentation (available at `/swagger/index.html#/`)
- Structured with a clean architecture separating controllers, services, and repositories

## Getting Started

### Prerequisites

- [Go 1.20+](https://golang.org/dl/)
- [Docker](https://www.docker.com/get-started) (if you prefer containerized deployment)

### Running Locally

1. Clone the repository:

```bash
   git clone https://github.com/yourusername/gorag.git
   cd gorag
```

Download dependencies:

```bash
go mod download
```

Run the application:

```bash
go run ./cmd/main/main.go
```

Open your browser and navigate to http://localhost:8711/swagger/index.html#/ to view the API documentation.

## Using Docker
### Build and Run with Docker
Build the Docker image:

```bash
docker build -t gorag .
```

Run the Docker container:

```bash
docker run -p 8711:8711 -v $(pwd)/data:/app gorag
```

Using Docker Compose
Alternatively, use Docker Compose to build and run the service:

```bash
docker compose up --build
```

This will start the service on port 8711 and persist the SQLite database file in the ./data directory.

## Environment Variables
DB_PATH - Path to the SQLite database file (default: ./gorag.db)
PORT - Port for the service to listen on (default: 8711)

## API Endpoints
POST /api/documents - Create a new document
GET /api/documents - Retrieve all documents
GET /api/documents/{id} - Retrieve a document by ID
PUT /api/documents/{id} - Update a document
DELETE /api/documents/{id} - Delete a document

## License
This project is licensed under the MIT License.