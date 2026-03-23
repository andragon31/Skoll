package generator

import (
	"fmt"
	"os"
	"path/filepath"
)

func InitProject(root string) error {
	skollDir := filepath.Join(root, ".skoll")
	dirs := []string{
		filepath.Join(skollDir, "rules"),
		filepath.Join(skollDir, "skills"),
		filepath.Join(skollDir, "agents"),
		filepath.Join(skollDir, "workflows"),
	}

	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return err
		}
	}

	// Create initial files
	files := map[string]string{
		filepath.Join(skollDir, "rules", "global.md"):      GlobalRulesTemplate,
		filepath.Join(skollDir, "skills", "git-workflow.md"): GitSkillTemplate,
		filepath.Join(skollDir, "agents", "backend.md"):     BackendAgentTemplate,
		filepath.Join(skollDir, "workflows", "feature.md"):  FeatureWorkflowTemplate,
	}

	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return err
		}
	}

	fmt.Printf("Skoll initialized in %s\n", skollDir)
	return nil
}
