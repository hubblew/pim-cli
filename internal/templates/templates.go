package templates

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
)

//go:embed generate_instructions.gohtml
var templatesFS embed.FS

func RenderGenerateInstructionsPrompt(instructionsDir string) (string, error) {
	content, err := templatesFS.ReadFile("generate_instructions.gohtml")
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %w", err)
	}

	tmpl, err := template.New("generate_instructions").Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	data := map[string]string{
		"InstructionsDir": instructionsDir,
	}
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
