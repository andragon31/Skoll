package rsaw

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
)

type Loader struct {
	root string
	gm   goldmark.Markdown
}

func NewLoader(root string) *Loader {
	return &Loader{
		root: root,
		gm:   goldmark.New(),
	}
}

func (l *Loader) LoadRules() ([]string, error) {
	rulesDir := filepath.Join(l.root, ".skoll", "rules")
	files, err := os.ReadDir(rulesDir)
	if err != nil {
		return nil, err
	}

	var rules []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".md") {
			content, err := os.ReadFile(filepath.Join(rulesDir, f.Name()))
			if err == nil {
				rules = append(rules, string(content))
			}
		}
	}
	return rules, nil
}

// Logic for parsing agents, skills, workflows will be added here
// For now, this is the base for the MCP server to start
