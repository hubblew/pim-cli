package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      string
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid config",
			config: `version: 1
sources:
  - key: source1
    url: /path/to/source1
  - key: source2
    url: https://github.com/user/repo.git
targets:
  - name: target1
    output: ./output
    include:
      - source: source1
        files:
          - file1.txt
      - source: source2
        files:
          - file2.txt
`,
			expectError: false,
		},
		{
			name: "duplicate source keys",
			config: `version: 1
sources:
  - key: duplicate
    url: /path/one
  - key: duplicate
    url: /path/two
`,
			expectError: true,
			errorMsg:    "duplicate source key: duplicate",
		},
		{
			name: "empty source key",
			config: `version: 1
sources:
  - key: ""
    url: /path/to/source
`,
			expectError: true,
			errorMsg:    "source key cannot be empty",
		},
		{
			name: "reference to non-existent source",
			config: `version: 1
sources:
  - key: source1
    url: /path/to/source
targets:
  - name: target1
    output: ./output
    include:
      - source: nonexistent
        files:
          - file.txt
`,
			expectError: true,
			errorMsg:    "target 'target1' references unknown source: nonexistent",
		},
		{
			name: "empty config with defaults",
			config: `version: 1
`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "test-config.yaml")

			if err := os.WriteFile(configPath, []byte(tt.config), 0644); err != nil {
				t.Fatalf("failed to write test config: %v", err)
			}

			cfg, err := LoadConfig(configPath)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("expected error %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if cfg == nil {
					t.Error("expected config to be non-nil")
				}
				if cfg != nil && cfg.Version != 1 {
					t.Errorf("expected version 1, got %d", cfg.Version)
				}
			}
		})
	}
}

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()

	if cfg.Version != 1 {
		t.Errorf("expected version 1, got %d", cfg.Version)
	}

	if cfg.Sources == nil {
		t.Error("expected Sources to be initialized")
	}

	if cfg.Targets == nil {
		t.Error("expected Targets to be initialized")
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid config",
			config: &Config{
				Version: 1,
				Sources: []Source{
					{Key: "s1", URL: "/path1"},
					{Key: "s2", URL: "/path2"},
				},
				Targets: []Target{
					{
						Name:   "t1",
						Output: "/output",
						Include: []Include{
							{Source: "s1", Files: []string{"file.txt"}},
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "duplicate keys",
			config: &Config{
				Version: 1,
				Sources: []Source{
					{Key: "same", URL: "/path1"},
					{Key: "same", URL: "/path2"},
				},
			},
			expectError: true,
			errorMsg:    "duplicate source key: same",
		},
		{
			name: "empty key",
			config: &Config{
				Version: 1,
				Sources: []Source{
					{Key: "", URL: "/path1"},
				},
			},
			expectError: true,
			errorMsg:    "source key cannot be empty",
		},
		{
			name: "unknown source reference",
			config: &Config{
				Version: 1,
				Sources: []Source{
					{Key: "s1", URL: "/path1"},
				},
				Targets: []Target{
					{
						Name:   "t1",
						Output: "/output",
						Include: []Include{
							{Source: "unknown", Files: []string{"file.txt"}},
						},
					},
				},
			},
			expectError: true,
			errorMsg:    "target 't1' references unknown source: unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				} else if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("expected error %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
