package rsaw

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type RSAWType string

const (
	TypeRule     RSAWType = "Rule"
	TypeSkill    RSAWType = "Skill"
	TypeAgent    RSAWType = "Agent"
	TypeWorkflow RSAWType = "Workflow"
	TypeUnknown  RSAWType = "Unknown"
)

type Item struct {
	Type RSAWType
	Name string
	Path string
}

func (l *Loader) ScanProject() ([]Item, error) {
	skollDir := filepath.Join(l.root, ".skoll")
	var items []Item

	err := filepath.Walk(skollDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			itemType, name := parseHeader(path)
			if itemType != TypeUnknown {
				items = append(items, Item{
					Type: itemType,
					Name: name,
					Path: path,
				})
			}
		}
		return nil
	})

	return items, err
}

func parseHeader(path string) (RSAWType, string) {
	file, err := os.Open(path)
	if err != nil {
		return TypeUnknown, ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "# ") {
			parts := strings.SplitN(strings.TrimPrefix(line, "# "), ":", 2)
			if len(parts) == 2 {
				typeName := strings.TrimSpace(parts[0])
				name := strings.TrimSpace(parts[1])
				
				switch typeName {
				case "Rules", "Rule":
					return TypeRule, name
				case "Skill":
					return TypeSkill, name
				case "Agent":
					return TypeAgent, name
				case "Workflow":
					return TypeWorkflow, name
				}
			}
		}
	}
	return TypeUnknown, ""
}
