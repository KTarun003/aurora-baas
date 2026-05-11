package api

import (
	"github.com/gin-gonic/gin"
	_ "github.com/ktarun.reddy/baas/docs"
	"github.com/ktarun.reddy/baas/internal/api/handlers"
	"github.com/ktarun.reddy/baas/internal/api/middleware"
	"github.com/ktarun.reddy/baas/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter creates and configures the Gin router
func NewRouter(
	projectService *service.ProjectService,
	schemaService *service.SchemaService,
	codegenService *service.CodegenService,
) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	// Add middleware
	router.Use(middleware.Logger())

	// Swagger documentation route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Create handlers
	healthHandler := handlers.NewHealthHandler()
	projectHandler := handlers.NewProjectHandler(projectService)
	schemaHandler := handlers.NewSchemaHandler(schemaService)
	codegenHandler := handlers.NewCodegenHandler(codegenService)

	// Health check endpoint
	router.GET("/health", healthHandler.Check)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Project routes
		projectsGroup := v1.Group("/projects")
		{
			projectsGroup.POST("", projectHandler.Create)
			projectsGroup.GET("", projectHandler.List)
			projectsGroup.GET("/:id", projectHandler.Get)
			projectsGroup.PUT("/:id", projectHandler.Update)
			projectsGroup.DELETE("/:id", projectHandler.Delete)

			// Code generation route
			projectsGroup.POST("/:id/generate", codegenHandler.GenerateProject)
		}

		// Schema routes
		schemasGroup := v1.Group("/projects/:id/schemas")
		{
			schemasGroup.POST("", schemaHandler.Apply)
			schemasGroup.GET("/latest", schemaHandler.GetLatest)
			schemasGroup.GET("", schemaHandler.List)
		}
	}

	return router
}
