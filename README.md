# Go Backend with SQLc

This is a Go backend application using Gin framework with SQLc and pgx for database operations, following clean architecture principles.

## Features

- RESTful API with Gin framework
- SQLc for type-safe database operations with PostgreSQL
- pgx/v5 for high-performance PostgreSQL driver
- Clean architecture with separation of concerns
- User CRUD operations with repository pattern
- Pagination support
- Database connection pooling
- Environment-based configuration
- Password hashing with bcrypt
- Graceful shutdown
- Request ID tracking
- CORS support
- Comprehensive error handling

## Architecture

The application follows a clean architecture pattern with the following structure:

```
go-backend-valos-id/
├── core/           # Core application logic
│   ├── config/     # Configuration management
│   ├── db/         # Database connection and setup
│   ├── handlers/   # HTTP handlers (presentation layer)
│   ├── middleware/ # HTTP middleware
│   ├── models/     # Data models and repositories
│   ├── server/     # Server setup and routing
│   └── utils/      # Utility functions
├── migrations/     # Database migrations
├── docs/          # Documentation
└── main.go        # Application entry point
```

## Database Schema

The application uses PostgreSQL with the following users table:

```sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

## API Endpoints

### Health Checks
- `GET /ping` - Basic ping endpoint
- `GET /health` - Application and database health check
- `GET /ready` - Readiness probe (Kubernetes)
- `GET /live` - Liveness probe (Kubernetes)

### User Management
- `POST /api/v1/users` - Create a new user
- `GET /api/v1/users` - Get all users
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user
- `GET /api/v1/users/paginate?limit=10&offset=0` - Get users with pagination

## Setup

1. Install dependencies:
```bash
go mod download
```

2. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your database credentials
```

3. Run database migration:
```bash
psql -h localhost -U postgres -d valos_db -f migrations/001_create_users_table.sql
```

4. Run the application:
```bash
go run main.go
```

## Environment Variables

- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5432)
- `DB_USER` - Database username (default: postgres)
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name (default: valos_db)
- `DB_SSL_MODE` - SSL mode (default: disable)
- `SERVER_PORT` - Server port (default: 3210)
- `GIN_MODE` - Gin mode (debug/release)

## Usage Examples

### Create a user
```bash
curl -X POST http://localhost:3210/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Get all users
```bash
curl http://localhost:3210/api/v1/users
```

### Get user by ID
```bash
curl http://localhost:3210/api/v1/users/1
```

### Update user
```bash
curl -X PUT http://localhost:3210/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe_updated",
    "email": "john.doe@example.com"
  }'
```

### Delete user
```bash
curl -X DELETE http://localhost:3210/api/v1/users/1
```

### Get users with pagination
```bash
curl "http://localhost:3210/api/v1/users/paginate?limit=5&offset=10"
```

### Health check
```bash
curl http://localhost:3210/health
```

## Key Components

### Configuration Layer (`core/config/`)
- Environment-based database configuration
- Connection string builder with default values

### Database Layer (`core/db/`)
- pgx/v5 connection pool management
- Connection pooling configuration
- Health check functionality

### Handlers Layer (`core/handlers/`)
- HTTP request handling
- Request validation
- Response formatting
- Error handling

### Middleware Layer (`core/middleware/`)
- CORS handling
- Request ID generation
- Error handling
- Logging

### Models Layer (`core/models/`)
- Data structures with JSON tags
- Repository pattern implementation using SQLc generated code
- Database operations abstraction

### Server Layer (`core/server/`)
- Application setup and initialization
- Route configuration
- Graceful shutdown handling

## SQLc + pgx Features Used

- Type-safe SQL queries with SQLc
- Compile-time query validation
- Auto-generated Go code from SQL
- Type-safe parameters and results
- High-performance pgx/v5 connection pooling
- Transaction support with pgx
- PostgreSQL-specific features support
- Optimized query performance

## Security Features

- Password hashing with bcrypt
- Input validation
- SQL injection prevention through parameterized queries
- CORS configuration

## Performance Features

- Connection pooling
- Database indexes
- Pagination support
- Type-safe compiled queries

## Development

### Running Tests
```bash
go test ./...
```

### Building
```bash
go build -o app .
```

### Running in Development
```bash
go run main.go
```

### Running with Air (hot reload)
```bash
air
```

## Deployment

The application is container-ready with:
- Health check endpoints for Kubernetes
- Graceful shutdown handling
- Environment-based configuration
- Structured logging

## Contributing

1. Follow the existing code structure
2. Add tests for new features
3. Update documentation
4. Ensure all tests pass before submitting

## License

This project is licensed under the MIT License.