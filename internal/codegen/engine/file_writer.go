package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ktarun.reddy/baas/internal/codegen"
)

// FileWriter handles writing generated files to disk
// It is NOT thread-safe and should be used by a single goroutine or protected by synchronization.
type FileWriter struct {
	outputDir string
}

// NewFileWriter creates a new file writer
func NewFileWriter(outputDir string) *FileWriter {
	return &FileWriter{
		outputDir: outputDir,
	}
}

// validatePath ensures the provided path does not escape outputDir
// Returns the validated absolute path or an error if the path attempts to traverse outside
func (w *FileWriter) validatePath(path string) (string, error) {
	// Normalize the provided path
	cleanPath := filepath.Clean(path)

	// Reject absolute paths
	if filepath.IsAbs(cleanPath) {
		return "", fmt.Errorf("path traversal: absolute paths are not allowed: %s", path)
	}

	// Reject paths that start with ".."
	if strings.HasPrefix(cleanPath, "..") {
		return "", fmt.Errorf("path traversal: path escapes output directory: %s", path)
	}

	// Join and normalize the full path
	fullPath := filepath.Join(w.outputDir, cleanPath)
	fullPath = filepath.Clean(fullPath)

	// Verify the result is within outputDir
	absOutputDir, err := filepath.Abs(w.outputDir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve output directory: %w", err)
	}

	absFullPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve path: %w", err)
	}

	// Ensure absFullPath is within absOutputDir
	if !strings.HasPrefix(absFullPath, absOutputDir+string(filepath.Separator)) && absFullPath != absOutputDir {
		return "", fmt.Errorf("path traversal: path escapes output directory: %s", path)
	}

	return fullPath, nil
}

// WriteFiles writes all generated files to disk
func (w *FileWriter) WriteFiles(files []*codegen.GeneratedFile) error {
	for _, file := range files {
		if err := w.WriteFile(file); err != nil {
			return fmt.Errorf("failed to write %s: %w", file.Path, err)
		}
	}
	return nil
}

// WriteFile writes a single file to disk
func (w *FileWriter) WriteFile(file *codegen.GeneratedFile) error {
	// Validate path to prevent directory traversal
	fullPath, err := w.validatePath(file.Path)
	if err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write file with specified mode
	mode := os.FileMode(file.Mode)
	if mode == 0 {
		mode = 0644 // Default mode
	}

	if err := os.WriteFile(fullPath, []byte(file.Content), mode); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// EnsureDir ensures a directory exists
func (w *FileWriter) EnsureDir(path string) error {
	// Validate path to prevent directory traversal
	fullPath, err := w.validatePath(path)
	if err != nil {
		return fmt.Errorf("invalid directory path: %w", err)
	}

	if err := os.MkdirAll(fullPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return nil
}

// Remove removes a file or directory
func (w *FileWriter) Remove(path string) error {
	// Validate path to prevent directory traversal
	fullPath, err := w.validatePath(path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	if err := os.RemoveAll(fullPath); err != nil {
		return fmt.Errorf("failed to remove: %w", err)
	}

	return nil
}

// Exists checks if a file or directory exists
// Returns false if the path is invalid or the file does not exist
func (w *FileWriter) Exists(path string) bool {
	// Validate path to prevent directory traversal
	fullPath, err := w.validatePath(path)
	if err != nil {
		return false
	}

	_, err = os.Stat(fullPath)
	return err == nil
}
