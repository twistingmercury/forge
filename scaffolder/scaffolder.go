// Package scaffolder provides the functionality to create a new golang project.
package scaffolder

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

var (
	projToken  = []byte(`{{project_name}}`)
	descToken  = []byte(`{{project_description}}`)
	modToken   = []byte(`{{module_path}}`)
	workDir    string
	projectDir string
)

func projectPath() string {
	return path.Join(workDir, projectDir)
}

// CreateProject creates a new software project using the given validTemplate.
func CreateProject(templatePath, projectDirName, moduleName string) (err error) {
	slog.Info("creating project", "projectDirName", projectDirName, "moduleName", moduleName, "templatePath", templatePath)
	_, err = os.Stat(templatePath)
	if err != nil {
		return fmt.Errorf("could not find the validTemplate `%s`: %w", templatePath, err)
	}

	workDir, _ = os.Getwd()
	projectDir = projectDirName

	if err = ExtractTemplate(templatePath); err != nil {
		return
	}

	if err = ReplaceTokens(projectPath(), moduleName, projectDirName); err != nil {
		return fmt.Errorf("failed replacing tokens: %w", err)
	}

	if err = GoModInit(moduleName, projectPath()); err != nil {
		return fmt.Errorf("could not initialize go module: %w", err)
	}

	if err = AddDependencies(projectPath()); err != nil {
		return fmt.Errorf("failed adding dependencies: %w", err)
	}

	if err = GoModTidy(projectPath()); err != nil {
		return fmt.Errorf("could not tidy go module: %w", err)
	}

	if err = GitInit(projectPath()); err != nil {
		return
	}

	if err = GitAdd(projectPath()); err != nil {
		return
	}

	if err = GitCommit(projectPath()); err != nil {
		return
	}

	return
}

// ExtractTemplate writes the validTemplate to the given destination.
func ExtractTemplate(templatePath string) (err error) {
	_, err = os.Stat(templatePath)
	if err != nil {
		return fmt.Errorf("could not find the validTemplate `%s`: %w", templatePath, err)
	}

	zReader, err := zip.OpenReader(templatePath)
	if err != nil {
		return fmt.Errorf("could not open validTemplate `%s`: %w", templatePath, err)
	}
	defer zReader.Close()

	tmpPath, err := os.MkdirTemp("", "_forge_*")

	for _, f := range zReader.File {
		err := unzip(f, tmpPath)
		if err != nil {
			return fmt.Errorf("could not unzip file `%s`: %w", f.Name, err)
		}
	}

	slog.Info("renaming temporary directory", "newPath", projectPath())
	if err := os.Rename(tmpPath, projectPath()); err != nil {
		return fmt.Errorf("could not rename directory: %w", err)
	}

	if err := os.RemoveAll(tmpPath); err != nil {
		return fmt.Errorf("could not remove temporary directory: %w", err)
	}

	return
}

func unzip(f *zip.File, tmpPath string) (err error) {
	fpath := filepath.Join(tmpPath, f.Name)
	if f.FileInfo().IsDir() {
		slog.Info("creating directory", "dir", f.Name)
		return os.MkdirAll(fpath, os.ModePerm)
	}

	slog.Info("creating file", "fpath", fpath)
	if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
		return fmt.Errorf("could not create directory `%s`: %w", filepath.Dir(fpath), err)
	}

	outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not open file `%s`: %w", fpath, err)
	}

	rc, err := f.Open()
	if err != nil {
		return fmt.Errorf("could not open file `%s`: %w", f.Name, err)
	}

	if _, err = io.Copy(outFile, rc); err != nil {
		return fmt.Errorf("could not copy file `%s`: %w", outFile.Name(), err)
	}

	if err := outFile.Close(); err != nil {
		return fmt.Errorf("could not close file `%s`: %w", outFile.Name(), err)
	}
	if err := rc.Close(); err != nil {
		return fmt.Errorf("could not close io.ReaderCloser: %w", err)
	}
	return
}

func AddDependencies(projPath string) error {
	slog.Info("adding dependencies")
	deps := path.Join(projectPath(), "_deps.sh")
	if _, err := os.Stat(deps); err != nil && os.IsNotExist(err) {
		slog.Info("no dependencies file found")
	}
	slog.Info("adding dependencies", "absPath", projPath)
	cmd := exec.Command("sh", deps)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	cmd.Dir = projPath
	if err := cmd.Run(); err != nil {
		slog.Error("failed adding dependencies", "error", err)
		return fmt.Errorf("failed adding dependencies: %s; %w", errb.String(), err)
	}
	return nil
}

// ReplaceTokens replaces the tokens in the given directory.
func ReplaceTokens(projPath, modPath, projName string) error {
	slog.Info("replacing tokens", "rootPath", projPath, "modPath", modPath, "projName", projName)
	return filepath.Walk(projPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && path != projPath {
			slog.Info("path != rootPath", path, projPath)
		}

		if !info.IsDir() {
			if err := ReplaceTokenInFile(path, modPath, projName); err != nil {
				return fmt.Errorf("failed replacing tokens in file `%s`: %w", path, err)
			}
		}
		return nil
	})
}

// ReplaceTokenInFile replaces the tokens in the given file.
func ReplaceTokenInFile(filePath, modPath, projName string) error {
	slog.Info("replacing tokens in file", "file path", filePath)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if bytes.ContainsAny(data, string(projToken)) &&
		bytes.ContainsAny(data, string(modToken)) &&
		bytes.ContainsAny(data, string(descToken)) {
		slog.Info("no tokens found in file", "filePath", filePath)
	}

	updatedData := bytes.ReplaceAll(data, projToken, []byte(projName))
	updatedData = bytes.ReplaceAll(updatedData, modToken, []byte(modPath))
	return os.WriteFile(filePath, updatedData, 0777)
}

// GoModInit initializes the go module.
func GoModInit(modPath string, projPath string) error {
	slog.Info("running go mod init")
	cmd := exec.Command("go", "mod", "init", modPath)
	return ExecCmd(cmd, projPath)
}

// GoModTidy runs go mod tidy.
func GoModTidy(projPath string) error {
	slog.Info("running go mod tidy")
	cmd := exec.Command("go", "mod", "tidy")
	return ExecCmd(cmd, projPath)
}

// GitInit initializes the git repo in the current directory.
func GitInit(projPath string) (err error) {
	slog.Info("initializing git repo: working branch = main")
	cmd := exec.Command("git", "init", "-b", "main")
	return ExecCmd(cmd, projPath)
}

// GitAdd adds all files that have been copied to the main trunk.
func GitAdd(projPath string) (err error) {
	slog.Info("adding files to working branch: main")
	cmd := exec.Command("git", "add", ".")
	return ExecCmd(cmd, projPath)
}

// GitCommit commits all files added to the main trunk.
func GitCommit(projPath string) (err error) {
	slog.Info("initial commit created")
	cmd := exec.Command("git", "commit", "-m", "created by forge: initial commit")
	return ExecCmd(cmd, projPath)
}

func ExecCmd(cmd *exec.Cmd, projPath string) (err error) {
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	cmd.Dir = projPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s; %w", errb.String(), err)
	}
	slog.Info(outb.String())
	return
}

// Rollback removes the project directory.
func Rollback(projPath string) error {
	slog.Info("rolling back", "absPath", projPath)
	return os.RemoveAll(projPath)
}
