# App vs Server Architecture

## Overview

In this Go backend application, we've separated concerns into two main structures: `App` and `Server`. This follows the principle of **separation of concerns** and makes the code more maintainable and testable.

## Server

The `Server` is responsible for the **technical infrastructure** of the application.

### Responsibilities
- **HTTP Server Management**: Creating and configuring the Gin router
- **Dependency Injection**: Setting up database connections, repositories, and handlers
- **Route Configuration**: Defining API endpoints and middleware
- **Database Setup**: Establishing database connections
- **Handler Initialization**: Creating and wiring up HTTP handlers

### What it Contains
```go
type Server struct {
    router        *gin.Engine        // HTTP router
    db            *db.Database      // Database connection
    healthHandler *handlers.HealthHandler
    userHandler   *handlers.UserHandler
}
```

### Key Methods
- `Initialize()` - Sets up all dependencies (database, handlers, router)
- `setupRouter()` - Configures middleware and routes
- `Start(addr)` - Starts the HTTP server
- `Close()` - Closes database connections

### Analogy
Think of the `Server` as the **"engine room"** - it handles all the technical machinery, wiring, and infrastructure.

## App

The `App` is responsible for the **application lifecycle** and runtime behavior.

### Responsibilities
- **Application Bootstrap**: Creating and initializing the server
- **Graceful Shutdown**: Handling SIGINT/SIGTERM signals
- **Process Management**: Managing the application's runtime
- **Error Handling**: Top-level error handling and logging
- **Signal Handling**: Responding to system signals

### What it Contains
```go
type App struct {
    server *Server    // Reference to the server
}
```

### Key Methods
- `NewApp()` - Factory function to create a new app instance
- `Initialize()` - Initializes the underlying server
- `Run(addr)` - Starts the server and waits for shutdown signals
- **Signal Handling** - Listens for shutdown signals and handles graceful shutdown

### Analogy
Think of the `App` as the **"captain's bridge"** - it manages the overall operation and lifecycle of the application.

## Why Separate Them?

### 1. **Single Responsibility Principle**
- `Server`: Manages technical infrastructure
- `App`: Manages application lifecycle

### 2. **Testability**
```go
// You can test the server without running the full app
func TestServerRoutes(t *testing.T) {
    server := NewServer()
    server.Initialize()
    // Test routes, handlers, etc.
}

// You can test app lifecycle separately
func TestAppShutdown(t *testing.T) {
    app := NewApp()
    // Test graceful shutdown behavior
}
```

### 3. **Flexibility**
- You can use the `Server` in different contexts (testing, CLI tools, etc.)
- You can have multiple `App` instances with different configurations
- Easier to extend for different deployment scenarios

### 4. **Clean Architecture**
```
main.go (Entry Point)
    ↓
App (Application Layer)
    ↓
Server (Infrastructure Layer)
    ↓
Handlers → Models → Database
```

## Flow Diagram

```
┌─────────────┐
│    main.go  │
│ (Entry Point)│
└─────┬───────┘
      │ creates
      ▼
┌─────────────┐     ┌─────────────────┐
│     App     │────▶│     Server      │
│ (Lifecycle) │     │ (Infrastructure)│
└─────┬───────┘     └─────────┬───────┘
      │ initializes               │ setup
      ▼                           ▼
┌─────────────┐           ┌─────────────┐
│   Signals   │           │   Router    │
│ (SIGINT/SIG│           │ (Gin Engine)│
│    TERM)    │           └─────────────┘
└─────┬───────┘                 │
      │ handles                 │
      ▼                         ▼
┌─────────────┐           ┌─────────────┐
│   Shutdown  │           │   Handlers  │
│  (Graceful) │           │ (Business   │
└─────────────┘           │    Logic)   │
                          └─────────────┘
```

## Usage Examples

### Current Usage
```go
func main() {
    app := core/server.NewApp()
    
    if err := app.Initialize(); err != nil {
        log.Fatalf("Failed to initialize application: %v", err)
        os.Exit(1)
    }
    
    addr := getServerAddr()
    if err := app.Run(addr); err != nil {
        log.Fatalf("Failed to run application: %v", err)
        os.Exit(1)
    }
}
```

### Alternative: Using Server Directly (for testing)
```go
func main() {
    server := core/server.NewServer()
    
    if err := server.Initialize(); err != nil {
        log.Fatal(err)
    }
    
    // No graceful shutdown - just run until process ends
    log.Fatal(server.Start(":8080"))
}
```

### Alternative: Multiple Apps (advanced scenarios)
```go
func main() {
    // HTTP API app
    apiApp := core/server.NewApp()
    apiApp.Initialize()
    go apiApp.Run(":8080")
    
    // Admin app with different configuration
    adminApp := core/server.NewApp()
    adminApp.Initialize()
    go adminApp.Run(":8081")
    
    // Wait for signals
    select {}
}
```

## Benefits Summary

| Aspect | Server | App |
|--------|---------|-----|
| **Focus** | Technical setup | Application lifecycle |
| **Scope** | Infrastructure | Process management |
| **Testing** | Unit test routes/handlers | Test lifecycle/shutdown |
| **Reusability** | Can be reused in tools | Specific to this application |
| **Complexity** | Lower (single responsibility) | Higher (orchestration) |

This separation makes the code more maintainable, testable, and follows clean architecture principles. Each component has a clear, focused responsibility.