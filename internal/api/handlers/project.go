package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ktarun.reddy/baas/internal/domain"
	"github.com/ktarun.reddy/baas/internal/service"
)

// ProjectHandler handles project-related requests
type ProjectHandler struct {
	service *service.ProjectService
}

// NewProjectHandler creates a new project handler
func NewProjectHandler(svc *service.ProjectService) *ProjectHandler {
	if svc == nil {
		panic("project service cannot be nil")
	}
	return &ProjectHandler{
		service: svc,
	}
}

// CreateProjectRequest represents the request body for creating a project
type CreateProjectRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	Language     string `json:"language" binding:"required"`
	DatabaseType string `json:"database_type" binding:"required"`
	APIStyle     string `json:"api_style" binding:"required"`
}

// Create handles creating a new project
// @Summary Create a new project
// @Description Create a new Backend as a Service project with specified configuration
// @Tags projects
// @Accept json
// @Produce json
// @Param project body CreateProjectRequest true "Project creation request"
// @Success 201 {object} domain.Project
// @Failure 400 {object} map[string]string
// @Router /projects [post]
func (h *ProjectHandler) Create(c *gin.Context) {
	var req CreateProjectRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body: " + err.Error(),
		})
		return
	}

	project := &domain.Project{
		Name:         req.Name,
		Description:  req.Description,
		Language:     req.Language,
		DatabaseType: req.DatabaseType,
		APIStyle:     req.APIStyle,
	}

	if err := h.service.CreateProject(project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, project)
}

// Get retrieves a project by ID
// @Summary Get project by ID
// @Description Retrieve detailed information about a specific project
// @Tags projects
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} domain.Project
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /projects/{id} [get]
func (h *ProjectHandler) Get(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "project id is required",
		})
		return
	}

	project, err := h.service.GetProject(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "project not found",
		})
		return
	}

	c.JSON(http.StatusOK, project)
}

// List retrieves all projects
// @Summary List all projects
// @Description Get a list of all Backend as a Service projects
// @Tags projects
// @Produce json
// @Success 200 {object} map[string][]domain.Project
// @Failure 500 {object} map[string]string
// @Router /projects [get]
func (h *ProjectHandler) List(c *gin.Context) {
	projects, err := h.service.ListProjects()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to list projects: " + err.Error(),
		})
		return
	}

	if projects == nil {
		projects = []*domain.Project{}
	}

	projectsData := make([]gin.H, 0, len(projects))
	for _, p := range projects {
		projectsData = append(projectsData, gin.H{
			"id":             p.ID,
			"name":           p.Name,
			"description":    p.Description,
			"language":       p.Language,
			"database_type":  p.DatabaseType,
			"api_style":      p.APIStyle,
			"created_at":     p.CreatedAt,
			"updated_at":     p.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"projects": projectsData,
	})
}

// Update updates an existing project
// @Summary Update project by ID
// @Description Update the configuration of an existing project
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Param project body CreateProjectRequest true "Project update request"
// @Success 200 {object} domain.Project
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /projects/{id} [put]
func (h *ProjectHandler) Update(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "project id is required",
		})
		return
	}

	// Check if project exists
	project, err := h.service.GetProject(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "project not found",
		})
		return
	}

	var req CreateProjectRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body: " + err.Error(),
		})
		return
	}

	// Update project fields
	project.Name = req.Name
	project.Description = req.Description
	project.Language = req.Language
	project.DatabaseType = req.DatabaseType
	project.APIStyle = req.APIStyle

	if err := h.service.UpdateProject(project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, project)
}

// Delete deletes a project by ID
// @Summary Delete project by ID
// @Description Remove a project and all its associated schemas
// @Tags projects
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /projects/{id} [delete]
func (h *ProjectHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "project id is required",
		})
		return
	}

	// Check if project exists
	_, err := h.service.GetProject(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "project not found",
		})
		return
	}

	if err := h.service.DeleteProject(id); err != nil {
		if err.Error() == "project not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "project not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to delete project: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "project deleted",
	})
}
