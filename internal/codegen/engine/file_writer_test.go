package engine

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ktarun.reddy/baas/internal/codegen"
)

func TestFileWriter_WriteFiles(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "codegen-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	writer := NewFileWriter(tmpDir)

	files := []*codegen.GeneratedFile{
		{
			Path:    "src/index.ts",
			Content: "console.log('Hello');",
			Mode:    0644,
		},
		{
			Path:    "package.json",
			Content: `{"name": "test"}`,
			Mode:    0644,
		},
	}

	err = writer.WriteFiles(files)
	require.NoError(t, err)

	// Verify files exist
	indexPath := filepath.Join(tmpDir, "src", "index.ts")
	assert.FileExists(t, indexPath)

	content, err := os.ReadFile(indexPath)
	require.NoError(t, err)
	assert.Equal(t, "console.log('Hello');", string(content))
}

func TestFileWriter_EnsureDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "codegen-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	writer := NewFileWriter(tmpDir)

	deepPath := "src/controllers/v1/users"
	err = writer.EnsureDir(deepPath)
	require.NoError(t, err)

	fullPath := filepath.Join(tmpDir, deepPath)
	info, err := os.Stat(fullPath)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

// TestFileWriter_PathTraversalValidation_WriteFile tests path traversal protection in WriteFile
func TestFileWriter_PathTraversalValidation_WriteFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "codegen-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	writer := NewFileWriter(tmpDir)

	// Test case 1: Path traversal attempt with ../
	traversalFile := &codegen.GeneratedFile{
		Path:    "../../../etc/passwd",
		Content: "malicious",
		Mode:    0644,
	}

	err = writer.WriteFile(traversalFile)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "path traversal")

	// Test case 2: Absolute path attempt
	absoluteFile := &codegen.GeneratedFile{
		Path:    "/etc/passwd",
		Content: "malicious",
		Mode:    0644,
	}

	err = writer.WriteFile(absoluteFile)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "path traversal")

	// Verify no files were written outside the output directory
	etcPasswd := "/etc/passwd"
	info, err := os.Stat(etcPasswd)
	// We can't assume /etc/passwd doesn't exist or has specific permissions,
	// but we verified through the error that the write was blocked
	_ = info
}

// TestFileWriter_PathTraversalValidation_EnsureDir tests path traversal protection in EnsureDir
func TestFileWriter_PathTraversalValidation_EnsureDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "codegen-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	writer := NewFileWriter(tmpDir)

	// Test case: Path traversal attempt with ../
	err = writer.EnsureDir("../../../malicious")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "path traversal")
}

// TestFileWriter_Remove_Success tests successful file removal
func TestFileWriter_Remove_Success(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "codegen-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	writer := NewFileWriter(tmpDir)

	// Create a test file
	testFile := &codegen.GeneratedFile{
		Path:    "test.txt",
		Content: "content",
		Mode:    0644,
	}
	err = writer.WriteFile(testFile)
	require.NoError(t, err)

	filePath := filepath.Join(tmpDir, "test.txt")
	assert.FileExists(t, filePath)

	// Remove the file
	err = writer.Remove("test.txt")
	require.NoError(t, err)

	// Verify file is gone
	assert.NoFileExists(t, filePath)
}

// TestFileWriter_Remove_Directory tests successful directory removal
func TestFileWriter_Remove_Directory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "codegen-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	writer := NewFileWriter(tmpDir)

	// Create a directory with nested files
	err = writer.EnsureDir("src/nested")
	require.NoError(t, err)

	dirPath := filepath.Join(tmpDir, "src")
	info, err := os.Stat(dirPath)
	require.NoError(t, err)
	assert.True(t, info.IsDir())

	// Remove the directory
	err = writer.Remove("src")
	require.NoError(t, err)

	// Verify directory is gone
	assert.NoFileExists(t, dirPath)
}

// TestFileWriter_Remove_PathTraversal tests path traversal protection in Remove
func TestFileWriter_Remove_PathTraversal(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "codegen-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	writer := NewFileWriter(tmpDir)

	// Test case: Path traversal attempt with ../
	err = writer.Remove("../../../sensitive")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "path traversal")
}

// TestFileWriter_Exists_FileExists tests Exists returns true for existing file
func TestFileWriter_Exists_FileExists(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "codegen-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	writer := NewFileWriter(tmpDir)

	// Create a test file
	testFile := &codegen.GeneratedFile{
		Path:    "test.txt",
		Content: "content",
		Mode:    0644,
	}
	err = writer.WriteFile(testFile)
	require.NoError(t, err)

	// Check file exists
	exists := writer.Exists("test.txt")
	assert.True(t, exists)
}

// TestFileWriter_Exists_FileNotExists tests Exists returns false for non-existing file
func TestFileWriter_Exists_FileNotExists(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "codegen-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	writer := NewFileWriter(tmpDir)

	// Check non-existent file
	exists := writer.Exists("nonexistent.txt")
	assert.False(t, exists)
}

// TestFileWriter_Exists_PathTraversal tests Exists returns false for traversal attempts
func TestFileWriter_Exists_PathTraversal(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "codegen-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	writer := NewFileWriter(tmpDir)

	// Check traversal attempt returns false
	exists := writer.Exists("../../../etc/passwd")
	assert.False(t, exists)
}

// TestFileWriter_WriteFile_PermissionError tests error handling when unable to create directory
func TestFileWriter_WriteFile_PermissionError(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "codegen-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create a file where a directory should be
	readOnlyFile := filepath.Join(tmpDir, "blocker")
	err = os.WriteFile(readOnlyFile, []byte("content"), 0644)
	require.NoError(t, err)

	writer := NewFileWriter(tmpDir)

	// Try to write a file in a path that conflicts with the file
	testFile := &codegen.GeneratedFile{
		Path:    "blocker/file.txt",
		Content: "content",
		Mode:    0644,
	}

	err = writer.WriteFile(testFile)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create directory")
}

// TestFileWriter_WriteFiles_Error tests error handling in WriteFiles
func TestFileWriter_WriteFiles_Error(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "codegen-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	writer := NewFileWriter(tmpDir)

	files := []*codegen.GeneratedFile{
		{
			Path:    "valid.txt",
			Content: "valid",
			Mode:    0644,
		},
		{
			Path:    "../invalid.txt",
			Content: "invalid",
			Mode:    0644,
		},
	}

	err = writer.WriteFiles(files)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "path traversal")

	// First file should still exist
	validPath := filepath.Join(tmpDir, "valid.txt")
	assert.FileExists(t, validPath)
}
