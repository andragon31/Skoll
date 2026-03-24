package rsaw

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/yuin/goldmark"
)

type Loader struct {
	root   string
	gm     goldmark.Markdown
	logger *log.Logger
}

func NewLoader(logger *log.Logger, root string) *Loader {
	return &Loader{
		root:   root,
		gm:     goldmark.New(),
		logger: logger,
	}
}

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

type Skill struct {
	Name        string
	Description string
	Path        string
	Metadata    map[string]string
}

type Agent struct {
	Name   string
	Path   string
	Skills []string
	Scope  map[string]interface{}
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

func (l *Loader) LoadRules(category string) ([]Item, error) {
	rulesDir := filepath.Join(l.root, ".skoll", "rules")
	var items []Item

	if _, err := os.Stat(rulesDir); os.IsNotExist(err) {
		return items, nil
	}

	err := filepath.Walk(rulesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			itemType, name := parseHeader(path)
			if itemType == TypeRule {
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

func (l *Loader) LoadSkillsIndex() ([]Skill, error) {
	skillsDir := filepath.Join(l.root, ".skoll", "skills")
	var skills []Skill

	if _, err := os.Stat(skillsDir); os.IsNotExist(err) {
		return skills, nil
	}

	err := filepath.Walk(skillsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != skillsDir {
			skillPath := filepath.Join(path, "SKILL.md")
			if _, err := os.Stat(skillPath); err == nil {
				content, _ := os.ReadFile(skillPath)
				metadata := l.ParseSKILLFrontmatter(string(content))
				name := filepath.Base(path)
				desc := metadata["description"]
				if desc == "" {
					desc = "No description"
				}
				skills = append(skills, Skill{
					Name:        name,
					Description: desc,
					Path:        skillPath,
					Metadata:    metadata,
				})
			}
		}
		return nil
	})

	return skills, err
}

func (l *Loader) ParseSKILLFrontmatter(content string) map[string]string {
	metadata := map[string]string{}

	if !strings.HasPrefix(content, "---") {
		return metadata
	}

	lines := strings.Split(content, "\n")
	var inFrontmatter bool
	var key, value string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "---" && !inFrontmatter {
			inFrontmatter = true
			continue
		}
		if line == "---" && inFrontmatter {
			break
		}
		if inFrontmatter && strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			key = strings.TrimSpace(parts[0])
			value = strings.TrimSpace(parts[1])
			metadata[key] = value
		}
	}

	return metadata
}

func (l *Loader) ListSkillFiles(skillName string) map[string][]string {
	skillDir := filepath.Join(l.root, ".skoll", "skills", skillName)
	files := map[string][]string{
		"scripts":    {},
		"references": {},
		"assets":     {},
	}

	for _, dir := range []string{"scripts", "references", "assets"} {
		path := filepath.Join(skillDir, dir)
		if entries, err := os.ReadDir(path); err == nil {
			for _, e := range entries {
				if !e.IsDir() {
					files[dir] = append(files[dir], e.Name())
				}
			}
		}
	}

	return files
}

func (l *Loader) LoadAgents() ([]Agent, error) {
	agentsDir := filepath.Join(l.root, ".skoll", "agents")
	var agents []Agent

	if _, err := os.Stat(agentsDir); os.IsNotExist(err) {
		return agents, nil
	}

	err := filepath.Walk(agentsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			itemType, name := parseHeader(path)
			if itemType == TypeAgent {
				agents = append(agents, Agent{
					Name: name,
					Path: path,
				})
			}
		}
		return nil
	})

	return agents, err
}

func (l *Loader) CreateSkillDir(name string) error {
	skillDir := filepath.Join(l.root, ".skoll", "skills", name)

	err := os.MkdirAll(filepath.Join(skillDir, "scripts"), 0755)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Join(skillDir, "references"), 0755)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Join(skillDir, "assets"), 0755)
	if err != nil {
		return err
	}

	template := `---
name: ` + name + `
description: |
  Add description here.
license: MIT
metadata:
  author: team
  version: "1.0"
---

## Cuándo aplicar

## Proceso

## Checklist

## Anti-patrones
`

	skillPath := filepath.Join(skillDir, "SKILL.md")
	return os.WriteFile(skillPath, []byte(template), 0644)
}

type Workflow struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Path      string   `json:"path"`
	Steps     []Step   `json:"steps"`
	Standards []string `json:"standards"`
}

type Step struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	DoD         []string `json:"dod"`
}

func (l *Loader) LoadWorkflows() ([]Workflow, error) {
	workflowsDir := filepath.Join(l.root, ".skoll", "workflows")
	var workflows []Workflow

	if _, err := os.Stat(workflowsDir); os.IsNotExist(err) {
		return workflows, nil
	}

	err := filepath.Walk(workflowsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			itemType, name := parseHeader(path)
			if itemType == TypeWorkflow {
				workflows = append(workflows, Workflow{
					ID:   "wf-" + uuid.New().String(),
					Name: name,
					Path: path,
				})
			}
		}
		return nil
	})

	return workflows, err
}

type PendingRule struct {
	ID        string `json:"id"`
	RuleID    string `json:"rule_id"`
	Reason    string `json:"reason"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

func (l *Loader) LoadPendingRules() ([]PendingRule, error) {
	proposedDir := filepath.Join(l.root, ".skoll", "rules", "_proposed")
	var rules []PendingRule

	if _, err := os.Stat(proposedDir); os.IsNotExist(err) {
		return rules, nil
	}

	err := filepath.Walk(proposedDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			rules = append(rules, PendingRule{
				ID:     "pending-" + uuid.New().String(),
				RuleID: filepath.Base(path),
				Status: "pending",
			})
		}
		return nil
	})

	return rules, err
}
