package analyzer

import (
	"os"
	"path/filepath"
	"strings"
)

type ProjectMeta struct {
	Stack      []string
	Modules    []string
	Language   string
	BuildTools []string
}

func Analyze(root string) ProjectMeta {
	meta := ProjectMeta{
		Stack:    []string{},
		Modules:  []string{},
		Language: "Unknown",
	}

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		// Skip invisible dirs like .git
		if strings.Contains(path, "/.") || strings.Contains(path, "\\.") {
			return nil
		}

		name := filepath.Base(path)
		switch name {
		case "go.mod":
			meta.Language = "Go"
			meta.BuildTools = append(meta.BuildTools, "go")
		case "package.json":
			meta.Language = "Node/JS/TS"
			meta.BuildTools = append(meta.BuildTools, "npm/yarn/pnpm")
		case "pom.xml", "build.gradle":
			meta.Language = "Java/Kotlin"
			meta.BuildTools = append(meta.BuildTools, "maven/gradle")
		case "pyproject.toml", "requirements.txt":
			meta.Language = "Python"
			meta.BuildTools = append(meta.BuildTools, "pip/poetry")
		case "Cargo.toml":
			meta.Language = "Rust"
			meta.BuildTools = append(meta.BuildTools, "cargo")
		}

		// Simple heuristics for modules
		if info.IsDir() && !strings.HasPrefix(info.Name(), ".") {
			meta.Modules = append(meta.Modules, info.Name())
		}

		return nil
	})

	return meta
}
