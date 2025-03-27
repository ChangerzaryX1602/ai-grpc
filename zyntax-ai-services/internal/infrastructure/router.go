package infrastructure

import (
	"github.com/Zentrix-Software-Hive/zyntax-ai-services/internal/handlers"
	"github.com/Zentrix-Software-Hive/zyntax-ai-services/pkg/ask"
	"github.com/Zentrix-Software-Hive/zyntax-ai-services/pkg/auth"
	"github.com/Zentrix-Software-Hive/zyntax-ai-services/pkg/models"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes is the Router for GoFiber App
func (s *Server) SetupRoutes(app *fiber.App) {

	// Prepare a static middleware to serve the built React files.
	app.Static("/", "./web/build")

	// API routes group
	groupApiV1 := app.Group("/api/v:version?", handlers.ApiLimiter)
	{
		groupApiV1.Get("/", handlers.Index())
	}
	app.Get("/api/v1/swagger/*", swagger.HandlerDefault)
	s.MainDbConn.AutoMigrate(&models.MainUser{})
	s.MainDbConn.AutoMigrate(&ask.History{}, &ask.HistoryMessage{})
	s.MainDbConn.AutoMigrate(&ask.MapHistoryMessage{}, &ask.MapUserHistory{})
	routerResource := handlers.NewRouterResources(s.JwtResources.JwtKeyfunc)
	authRepository := auth.NewAuthRepository(s.MainDbConn)
	authService := auth.NewAuthService(authRepository)
	ask.NewAskHandler(app.Group("/api/v1/ask"), routerResource, ":50051", s.MainDbConn)
	auth.NewAuthHandler(app.Group("/api/v1/auth"), authService, s.JwtResources)
	// Prepare a fallback route to always serve the 'index.html', had there not be any matching routes.
	app.Static("*", "./web/build/index.html")
}
