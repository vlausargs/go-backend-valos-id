package server

import (
	"log"
	"sort"
	"strings"

	"go-backend-valos-id/core/config"
	"go-backend-valos-id/core/db"
	"go-backend-valos-id/core/handlers"
	"go-backend-valos-id/core/middleware"
	user_handler "go-backend-valos-id/core/user/handler"
	user_repository "go-backend-valos-id/core/user/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Pre-computed method order for fastest comparison
var methodOrder = map[string]int{
	"DELETE":  0,
	"GET":     1,
	"POST":    2,
	"PUT":     3,
	"PATCH":   4,
	"HEAD":    5,
	"OPTIONS": 6,
}

type Server struct {
	router        *gin.Engine
	pool          *pgxpool.Pool
	healthHandler *handlers.HealthHandler
	userHandler   *user_handler.UserHandler
	database      *db.Database // Keep reference for cleanup
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Initialize() error {
	// Initialize database configuration
	dbConfig := config.NewDatabaseConfig()

	// Initialize database connection
	database, err := db.NewDatabase(dbConfig)
	if err != nil {
		return err
	}
	s.pool = database.Pool
	s.database = database // Keep reference for cleanup

	// Initialize repositories
	userRepo := user_repository.NewUserRepository(s.pool)

	// Initialize handlers
	s.healthHandler = handlers.NewHealthHandler(s.pool)
	s.userHandler = user_handler.NewUserHandler(userRepo)

	// Setup router
	s.setupRouter()

	return nil
}

func (s *Server) setupRouter() {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create router
	s.router = gin.New()

	// Add middleware - use built-in gin.Logger for route logging
	s.router.Use(middleware.RequestID())
	s.router.Use(gin.Logger())
	s.router.Use(gin.Recovery())
	s.router.Use(middleware.ErrorHandler())
	s.router.Use(middleware.CORS())

	// Setup routes
	s.setupRoutes()

	// Log all registered routes on startup
	logRegisteredRoutes(s.router)
}

func (s *Server) setupRoutes() {
	// Health check routes
	s.router.GET("/ping", s.healthHandler.Ping)
	s.router.GET("/health", s.healthHandler.HealthCheck)
	s.router.GET("/ready", s.healthHandler.Readiness)
	s.router.GET("/live", s.healthHandler.Liveness)

	// API routes v1
	v1 := s.router.Group("/api/v1")
	{
		// User routes
		users := v1.Group("/users")
		{
			users.POST("", s.userHandler.CreateUser)
			users.GET("", s.userHandler.GetAllUsers)
			users.GET("/paginate", s.userHandler.GetUsersWithPagination)
			users.GET("/:id", s.userHandler.GetUserByID)
			users.PUT("/:id", s.userHandler.UpdateUser)
			users.DELETE("/:id", s.userHandler.DeleteUser)
		}
	}
}

func (s *Server) Start(addr string) error {
	log.Printf("Server starting on %s", addr)
	return s.router.Run(addr)
}

func (s *Server) Close() error {
	if s.database != nil {
		return s.database.Close()
	}
	return nil
}

// logRegisteredRoutes logs all registered routes when server starts
func logRegisteredRoutes(engine *gin.Engine) {
	log.Println("🚀 Registered Routes:")
	log.Println(strings.Repeat("─", 30))

	routes := engine.Routes()

	// Fastest approach: pre-computed order map + direct comparisons
	sort.SliceStable(routes, func(i, j int) bool {
		// First compare paths (most common difference)
		if routes[i].Path != routes[j].Path {
			return routes[i].Path < routes[j].Path
		}

		// Then compare methods using pre-computed order
		orderI := methodOrder[routes[i].Method]
		orderJ := methodOrder[routes[j].Method]

		// If both methods are in our order map, use that
		if orderI != orderJ {
			return orderI < orderJ
		}

		// Fallback: if both unknown methods, compare alphabetically
		return routes[i].Method < routes[j].Method
	})

	for _, route := range routes {
		log.Printf("%-8s %s", route.Method, route.Path)
	}

	log.Println(strings.Repeat("─", 30))
	log.Printf("Total routes registered: %d", len(routes))
}
