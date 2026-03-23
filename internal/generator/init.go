package generator

import (
	"fmt"
	"os"
	"path/filepath"
)

func InitProjectWithMeta(root string, meta interface{}) error {
	skollDir := filepath.Join(root, ".skoll")
	dirs := []string{"rules", "skills", "agents", "workflows"}
	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(skollDir, d), 0755); err != nil {
			return err
		}
	}

	// Just copy basic templates for now
	if err := CopyTemplate("agent", filepath.Join(skollDir, "agents", "default.md")); err != nil {
		return err
	}
	if err := CopyTemplate("skill", filepath.Join(skollDir, "skills", "base.md")); err != nil {
		return err
	}
	if err := CopyTemplate("rule", filepath.Join(skollDir, "rules", "global.md")); err != nil {
		return err
	}

	fmt.Printf("Skoll .skoll/ structure generated from analysis.\n")
	return nil
}

func InitProject(root string) error {
	return InitProjectWithMeta(root, nil)
}
