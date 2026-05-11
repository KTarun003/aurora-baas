package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ktarun.reddy/baas/internal/service"
)

// CodegenHandler handles code generation requests
type CodegenHandler struct {
	service *service.CodegenService
}

// NewCodegenHandler creates a new codegen handler
func NewCodegenHandler(svc *service.CodegenService) *CodegenHandler {
	if svc == nil {
		panic("codegen service cannot be nil")
	}
	return &CodegenHandler{
		service: svc,
	}
}

// GenerateProjectRequest represents the request body for generating code
type GenerateProjectRequest struct {
	OutputDir string `json:"output_dir" binding:"required"`
}

// GenerateProject handles code generation for a project
// @Summary Generate backend code
// @Description Trigger code generation for a specific project based on its latest schema
// @Tags codegen
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Param request body GenerateProjectRequest true "Code generation request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /projects/{id}/generate [post]
func (h *CodegenHandler) GenerateProject(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "project_id is required",
		})
		return
	}

	var req GenerateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body: " + err.Error(),
		})
		return
	}

	// Generate the project
	generatedProject, err := h.service.GenerateProject(projectID, req.OutputDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate project: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"project_id":  generatedProject.ProjectID,
		"language":    generatedProject.Language,
		"files_count": len(generatedProject.Files),
		"metadata":    generatedProject.Metadata,
	})
}
