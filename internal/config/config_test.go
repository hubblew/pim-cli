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
		{
			name: "omitted source defaults to working_dir",
			config: `version: 1
targets:
  - name: target1
    output: ./output
    include:
      - files:
          - file1.txt
`,
			expectError: false,
		},
		{
			name: "working_dir source is always added",
			config: `version: 1
sources:
  - key: custom
    url: /path/to/custom
targets:
  - name: target1
    output: ./output
    include:
      - files:
          - file1.txt
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
				if cfg != nil {
					hasWorkingDir := false
					for _, source := range cfg.Sources {
						if source.Name == "working_dir" {
							hasWorkingDir = true
							break
						}
					}
					if !hasWorkingDir {
						t.Error("expected working_dir source to be present")
					}
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
					{Name: "s1", URL: "/path1"},
					{Name: "s2", URL: "/path2"},
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
					{Name: "same", URL: "/path1"},
					{Name: "same", URL: "/path2"},
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
					{Name: "", URL: "/path1"},
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
					{Name: "s1", URL: "/path1"},
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

func TestWorkingDirSource(t *testing.T) {
	cfg := &Config{
		Version: 1,
		Sources: []Source{
			{Name: "custom", URL: "/path/to/custom"},
		},
	}

	if err := cfg.addWorkingDirSource(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cfg.Sources) != 2 {
		t.Fatalf("expected 2 sources, got %d", len(cfg.Sources))
	}

	if cfg.Sources[0].Name != "working_dir" {
		t.Errorf("expected first source to be working_dir, got %s", cfg.Sources[0].Name)
	}

	if cfg.Sources[0].URL == "" {
		t.Error("expected working_dir URL to be set")
	}
}

func TestDefaultSourceForIncludes(t *testing.T) {
	cfg := &Config{
		Version: 1,
		Targets: []Target{
			{
				Name:   "target1",
				Output: "/output",
				Include: []Include{
					{Files: []string{"file1.txt"}},
					{Source: "custom", Files: []string{"file2.txt"}},
					{Files: []string{"file3.txt"}},
				},
			},
		},
	}

	cfg.setDefaultSourceForIncludes()

	if cfg.Targets[0].Include[0].Source != "working_dir" {
		t.Errorf("expected first include source to be working_dir, got %s", cfg.Targets[0].Include[0].Source)
	}

	if cfg.Targets[0].Include[1].Source != "custom" {
		t.Errorf("expected second include source to remain custom, got %s", cfg.Targets[0].Include[1].Source)
	}

	if cfg.Targets[0].Include[2].Source != "working_dir" {
		t.Errorf("expected third include source to be working_dir, got %s", cfg.Targets[0].Include[2].Source)
	}
}
