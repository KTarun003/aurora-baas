package handlers

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ktarun.reddy/baas/internal/domain"
	"github.com/ktarun.reddy/baas/internal/service"
)

// SchemaHandler handles schema-related requests
type SchemaHandler struct {
	service *service.SchemaService
}

// NewSchemaHandler creates a new schema handler
func NewSchemaHandler(svc *service.SchemaService) *SchemaHandler {
	if svc == nil {
		panic("schema service cannot be nil")
	}
	return &SchemaHandler{
		service: svc,
	}
}

// Apply handles applying a new schema
// @Summary Apply a new schema version
// @Description Apply a new YAML schema definition to a project. This will increment the schema version.
// @Tags schemas
// @Accept plain
// @Produce json
// @Param id path string true "Project ID"
// @Param schema body string true "YAML schema content"
// @Success 201 {object} domain.Schema
// @Failure 400 {object} map[string]string
// @Router /projects/{id}/schemas [post]
func (h *SchemaHandler) Apply(c *gin.Context) {
	projectID := c.Param("id")

	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "project id is required",
		})
		return
	}

	// Limit request body size to 10MB to prevent memory exhaustion
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20)

	// Read raw YAML body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read request body",
		})
		return
	}

	if len(body) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "request body cannot be empty",
		})
		return
	}

	yamlContent := string(body)

	schema, err := h.service.ApplySchema(projectID, yamlContent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, schema)
}

// GetLatest retrieves the latest schema version for a project
// @Summary Get latest schema version
// @Description Retrieve the most recent schema version for a specific project
// @Tags schemas
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} domain.Schema
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /projects/{id}/schemas/latest [get]
func (h *SchemaHandler) GetLatest(c *gin.Context) {
	projectID := c.Param("id")

	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "project id is required",
		})
		return
	}

	schema, err := h.service.GetLatestSchema(projectID)
	if err != nil || schema == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "no schema found for project",
		})
		return
	}

	c.JSON(http.StatusOK, schema)
}

// List retrieves all schema versions for a project
// @Summary List all schema versions
// @Description Get a list of all schema versions associated with a specific project
// @Tags schemas
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} map[string][]domain.Schema
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /projects/{id}/schemas [get]
func (h *SchemaHandler) List(c *gin.Context) {
	projectID := c.Param("id")

	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "project id is required",
		})
		return
	}

	schemas, err := h.service.ListSchemas(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to list schemas: " + err.Error(),
		})
		return
	}

	if schemas == nil {
		schemas = []*domain.Schema{}
	}

	schemasData := make([]gin.H, 0, len(schemas))
	for _, s := range schemas {
		schemasData = append(schemasData, gin.H{
			"id":         s.ID,
			"project_id": s.ProjectID,
			"version":    s.Version,
			"content":    s.Content,
			"created_at": s.CreatedAt,
			"updated_at": s.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"schemas": schemasData,
	})
}
