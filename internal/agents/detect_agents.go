package agents

import (
	"os/exec"
)

// DetectAgents checks for known LLM CLI tools in the system
func DetectAgents() []string {
	var tools []string

	if isCommandAvailable("copilot") {
		// Check if gh copilot extension is available
		cmd := exec.Command("copilot", "--version")
		if err := cmd.Run(); err == nil {
			tools = append(tools, "GitHub Copilot")
		}
	}

	return tools
}

// isCommandAvailable checks if a command is available in PATH
func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
