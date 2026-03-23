package generator

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/andragon31/skoll/internal/utils"
)

func InjectOpenCodeMCP() error {
	dir := utils.OpenCodeConfigDir()
	configPath := filepath.Join(dir, "opencode.json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	cleanData := utils.StripJSONC(data)

	var config map[string]interface{}
	if err := json.Unmarshal(cleanData, &config); err != nil {
		return err
	}

	mcpBlock, _ := config["mcpServers"].(map[string]interface{})
	if mcpBlock == nil {
		mcpBlock = make(map[string]interface{})
	}

	mcpBlock["skoll"] = map[string]interface{}{
		"type":    "command",
		"command": utils.ResolveBinaryPath(),
		"args":    []string{"mcp"},
		"enabled": true,
	}

	config["mcpServers"] = mcpBlock
	output, _ := json.MarshalIndent(config, "", "  ")
	return os.WriteFile(configPath, output, 0644)
}
