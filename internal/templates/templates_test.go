package templates

import "testing"

func TestRenderGenerateInstructionsPrompt(t *testing.T) {
	_, err := RenderGenerateInstructionsPrompt("/path/to/instructions")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
