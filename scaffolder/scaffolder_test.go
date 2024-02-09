package scaffolder_test

import (
	"os"
	"path"
	"testing"
	"twistingmercury/forge/scaffolder"

	"github.com/stretchr/testify/assert"
)

const (
	testDir       = "unitTestDir"
	validTemplate = "test_templates/test_template.zip"
	noDeps        = "test_templates/no_deps.zip"
)

var fullProjPath string

func init() {
	wd, _ := os.Getwd()
	fullProjPath = path.Join(wd, testDir)
}

func TestTemplateFuncs(t *testing.T) {
	os.RemoveAll(fullProjPath)
	defer func() {
		err := os.RemoveAll(fullProjPath)
		assert.NoErrorf(t, err, "remove all failed: %v", err)
	}()

	wd, _ := os.Getwd()
	t.Run("test ExtractTemplate", func(t *testing.T) {
		scaffolder.ProjectDir(testDir)
		scaffolder.WorkDir(wd)
		zipFile := path.Join(wd, validTemplate)
		err := scaffolder.ExtractTemplate(zipFile)
		if !assert.NoErrorf(t, err, "write validTemplate zip failed: %v", err) {
			t.FailNow()
		}
		projPath := path.Join(wd, testDir)
		if !assert.DirExistsf(t, projPath, "expected %s to be a directory, but it isn't", projPath) {
			t.FailNow()
		}
	})

	t.Run("test ReplaceTokens", func(t *testing.T) {
		fullProjPath := path.Join(wd, testDir)
		err := scaffolder.ReplaceTokens(fullProjPath, "github.com/username/projectname", "projectname")
		if !assert.NoErrorf(t, err, "replace tokens failed: %v", err) {
			t.FailNow()
		}
	})

	t.Run("test GoModInit", func(t *testing.T) {
		err := scaffolder.GoModInit("github.com/username/projectname", fullProjPath)
		if !assert.NoErrorf(t, err, "go mod init failed: %v", err) {
			t.FailNow()
		}
		if !assert.FileExistsf(t, path.Join(fullProjPath, "go.mod"), "expected %s to be a file, but it isn't", path.Join(fullProjPath, "go.mod")) {
			t.FailNow()
		}
	})

	t.Run("test GoModTidy", func(t *testing.T) {
		err := scaffolder.GoModTidy(fullProjPath)
		if !assert.NoErrorf(t, err, "go mod tidy failed: %v", err) {
			t.FailNow()
		}
	})

	t.Run("test AddDependencies", func(t *testing.T) {
		err := scaffolder.AddDependencies(fullProjPath)
		if !assert.NoErrorf(t, err, "add dependencies failed: %v", err) {
			t.FailNow()
		}
	})

	t.Run("test GitInit", func(t *testing.T) {
		err := scaffolder.GitInit(fullProjPath)
		if !assert.NoErrorf(t, err, "git init failed: %v", err) {
			t.FailNow()
		}
		if !assert.DirExistsf(t, path.Join(fullProjPath, ".git"), "expected %s to be a directory, but it isn't", path.Join(fullProjPath, ".git")) {
			t.FailNow()
		}
	})

	t.Run("test GitAdd", func(t *testing.T) {
		err := scaffolder.GitAdd(fullProjPath)
		if !assert.NoErrorf(t, err, "git add failed: %v", err) {
			t.FailNow()
		}
	})

	t.Run("test GitCommit", func(t *testing.T) {
		err := scaffolder.GitCommit(fullProjPath)
		if !assert.NoErrorf(t, err, "git commit failed: %v", err) {
			t.FailNow()
		}
	})

	t.Run("test Rollback", func(t *testing.T) {
		err := scaffolder.Rollback(fullProjPath)
		if !assert.NoErrorf(t, err, "rollback failed: %v", err) {
			t.FailNow()
		}
	})
}

func TestNewCreateProject(t *testing.T) {
	os.RemoveAll(fullProjPath)
	defer func() {
		err := os.RemoveAll(fullProjPath)
		assert.NoErrorf(t, err, "remove all failed: %v", err)
	}()
	wd, _ := os.Getwd()
	zipFile := path.Join(wd, validTemplate)
	err := scaffolder.CreateProject(
		zipFile,
		testDir,
		"github.com/username/projectname")
	assert.NoErrorf(t, err, "create project failed: %v", err)
}

func TestNewCreateProject_no_deps_file(t *testing.T) {
	os.RemoveAll(fullProjPath)
	defer func() {
		err := os.RemoveAll(fullProjPath)
		assert.NoErrorf(t, err, "remove all failed: %v", err)
	}()
	wd, _ := os.Getwd()
	zipFile := path.Join(wd, noDeps)
	err := scaffolder.CreateProject(
		zipFile,
		testDir,
		"github.com/username/projectname")
	assert.Errorf(t, err, "create project failed: %v", err)
}

func TestTemplateFuncErrors(t *testing.T) {
	invalidPath := "/some/invalid/path"

	t.Run("test Rollback", func(t *testing.T) {
		err := scaffolder.CreateProject(invalidPath, testDir, "github.com/username/projectname")
		assert.Errorf(t, err, "create project should have failed: %v", err)
	})
	wd, _ := os.Getwd()
	t.Run("test ExtractTemplate error", func(t *testing.T) {
		scaffolder.ProjectDir(testDir)
		scaffolder.WorkDir(wd)
		zipFile := path.Join(wd, invalidPath)
		err := scaffolder.ExtractTemplate(zipFile)
		assert.Error(t, err)
	})

	t.Run("test ReplaceTokens error", func(t *testing.T) {
		err := scaffolder.ReplaceTokens(invalidPath, "github.com/username/projectname", "projectname")
		assert.Error(t, err)
	})

	t.Run("test ReplaceTokensInFile error", func(t *testing.T) {
		err := scaffolder.ReplaceTokenInFile(invalidPath, "github.com/username/projectname", "projectname")
		assert.Error(t, err)
	})

	t.Run("test GoModInit error", func(t *testing.T) {
		err := scaffolder.GoModInit("github.com/username/projectname", invalidPath)
		assert.Error(t, err)
	})

	t.Run("test GoModTidy error", func(t *testing.T) {
		err := scaffolder.GoModTidy(invalidPath)
		assert.Error(t, err)
	})

	t.Run("test AddDependencies error", func(t *testing.T) {
		err := scaffolder.AddDependencies(invalidPath)
		assert.Error(t, err)
	})

	t.Run("test GitInit error", func(t *testing.T) {
		err := scaffolder.GitInit(invalidPath)
		assert.Error(t, err)
	})

	t.Run("test GitAdd error", func(t *testing.T) {
		err := scaffolder.GitAdd(invalidPath)
		assert.Error(t, err)
	})

	t.Run("test GitCommit error", func(t *testing.T) {
		err := scaffolder.GitCommit(invalidPath)
		assert.Error(t, err)
	})
}
