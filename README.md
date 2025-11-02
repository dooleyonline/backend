# DooleyOnline Backend

A RESTful API backend service built with Go for the DooleyOnline platform. This service provides comprehensive functionality for user management, item listings, categories, and real-time chat capabilities.

## 🏗️ Architecture

This project follows a clean architecture pattern with the following structure:

- **API Layer** (`internal/api/`) - HTTP handlers and middleware
- **Service Layer** (`internal/service/`) - Business logic
- **Data Layer** (`internal/db/`) - Database operations using sqlc
- **Models** (`internal/model/`) - Domain models
- **Storage** (`internal/storage/`) - File storage management

## ✨ Features

- **Authentication & Authorization**

  - JWT-based authentication
  - Secure password hashing
  - Token management

- **Item Management**

  - CRUD operations for items
  - Category-based organization
  - Pagination support
  - Image storage integration

- **User Management**

  - User registration and profiles
  - User data management

- **Real-time Chat**

  - WebSocket-based messaging
  - Chat rooms and participants
  - Message history

- **Storage Integration**
  - AWS S3-compatible storage
  - Image processing and optimization
  - Secure file uploads

## 🚀 Tech Stack

- **Language**: Go 1.25
- **Web Framework**: Echo v4
- **Database**: PostgreSQL (via pgx/v5)
- **Database Management**: sqlc for type-safe SQL
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **Storage**: AWS SDK v2
- **WebSocket**: Gorilla WebSocket
- **Documentation**: Swagger/OpenAPI
- **Image Processing**: disintegration/imaging
- **Task Runner**: Task (taskfile)
- **AI Integration**: Google Gen AI

## 📋 Prerequisites

- Go 1.25 or higher
- PostgreSQL database
- AWS S3 or S3-compatible storage (for file uploads)
- [Task](https://taskfile.dev/) (optional, for task automation)

## 🛠️ Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/dooleyonline/backend.git
   cd backend
   ```

2. **Install dependencies**

   ```bash
   go mod download
   ```

3. **Set up environment variables**

   Create a `.env` file in the root directory:

   ```env
   # Environment
   ENV=development
   PORT=8080

   # Database
   DATABASE_URL=postgres://username:password@localhost:5432/dooleyonline?sslmode=disable

   # Storage (AWS S3 or compatible)
   STORAGE_BUCKET=your-bucket-name
   STORAGE_URL=https://your-storage-url
   STORAGE_S3_URL=https://s3.your-region.amazonaws.com
   STORAGE_REGION=your-region
   STORAGE_ACCESS_ID=your-access-key-id
   STORAGE_ACCESS_SECRET=your-secret-access-key

   # Authentication
   AUTH_TOKEN_SECRET=your-secret-key-for-jwt
   ```

4. **Generate database code (if needed)**

   ```bash
   task sql
   # or
   sqlc generate
   ```

5. **Generate API documentation**
   ```bash
   task docs
   # or
   swag init --dir ./internal/api,./internal/db/item,./internal/db/chat,./internal/db/user --generalInfo api.go --parseInternal --parseDependency
   ```

## 🏃‍♂️ Running the Application

### Using Task

```bash
# Run the server
task run

# Run tests
task test

# Generate documentation
task docs

# Generate database code
task sql

# Run data generator
task datagen
```

### Using Go directly

```bash
# Run the server
go run ./cmd/server

# Run tests
go test ./... -v
```

### Using Docker

```bash
# Build the image
docker build -t dooleyonline-backend .

# Run the container
docker run -p 8080:8080 --env-file .env dooleyonline-backend
```

## 📚 API Documentation

Once the server is running, you can access the Swagger documentation at:

```
http://localhost:8080/swagger/index.html
```

## 🗂️ Project Structure

```
.
├── cmd/
│   ├── datagen/          # Data generation utility
│   └── server/           # Main application entry point
├── internal/
│   ├── api/              # HTTP handlers
│   │   ├── auth/         # Authentication endpoints
│   │   ├── category/     # Category endpoints
│   │   ├── chat/         # Chat & WebSocket endpoints
│   │   ├── item/         # Item endpoints
│   │   ├── shared/       # Shared API utilities
│   │   └── user/         # User endpoints
│   ├── config/           # Configuration management
│   ├── db/               # Database layer (sqlc generated)
│   │   ├── chat/         # Chat-related queries
│   │   ├── item/         # Item-related queries
│   │   └── user/         # User-related queries
│   ├── model/            # Domain models
│   ├── service/          # Business logic layer
│   │   ├── auth/         # Authentication service
│   │   ├── category/     # Category service
│   │   ├── chat/         # Chat service
│   │   ├── item/         # Item service
│   │   └── user/         # User service
│   └── storage/          # File storage client
├── sqlc/                 # SQL schemas and queries
├── docs/                 # Generated API documentation
├── Dockerfile            # Container definition
├── go.mod                # Go dependencies
├── sqlc.yaml             # sqlc configuration
└── taskfile.yaml         # Task automation
```

## 🧪 Testing

Run all tests:

```bash
go test ./... -v
```

Run tests with coverage:

```bash
go test ./... -cover
```

## 🔧 Development

### Database Migrations

This project uses raw SQL managed by sqlc. Schema definitions are in `sqlc/*/schema.sql` and queries are in `sqlc/*/query.sql`.

After modifying SQL files, regenerate the Go code:

```bash
task sql
```

### Adding New API Endpoints

1. Define the SQL schema in `sqlc/{domain}/{entity}/schema.sql`
2. Define queries in `sqlc/{domain}/{entity}/query.sql`
3. Generate database code: `task sql`
4. Implement service logic in `internal/service/{domain}/`
5. Create handler in `internal/api/{domain}/`
6. Register routes in `internal/api/api.go`
7. Update documentation: `task docs`

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the Apache 2.0 License - see the LICENSE file for details.

## 📧 Contact

- **API Support**: support@swagger.io
- **Project Link**: https://github.com/dooleyonline/backend

## 🙏 Acknowledgments

- Echo Framework for the excellent web framework
- sqlc for type-safe SQL code generation
- The Go community for amazing tools and libraries
