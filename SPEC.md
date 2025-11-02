# PIM Specification

## Overview
PIM (Prompt Instruction Manager) is a command-line utility for managing prompt instructions and related files.

## Configuration File

The tool uses a YAML configuration file to manage settings and package definitions.

### Configuration Format

```yaml
version: 1  # Configuration schema version (default: 1)

sources:
  - key: local-dir        # Unique identifier for this source
    url: /path/to/dir     # Local directory path or git repository URL
  - key: git-repo
    url: https://github.com/username/repo.git

targets:
  - name: my-target       # Target name
    output: ./output/dir  # Output directory for downloaded files
    include:
      - source: local-dir # Reference to source key
        files:            # List of file paths to include
          - file1.txt
          - folder/file2.txt
      - source: git-repo
        files:
          - README.md
```

#### Minimal Configuration Example

Since the `working_dir` source is automatically available and `source` defaults to `working_dir`, you can write a minimal configuration:

```yaml
version: 1

targets:
  - name: my-target
    output: ./output
    include:
      - files:              # source defaults to working_dir
          - prompts/system.txt
          - prompts/user.txt
```

This is equivalent to:

```yaml
version: 1

sources:
  - key: working_dir      # Automatically added
    url: /current/working/directory

targets:
  - name: my-target
    output: ./output
    include:
      - source: working_dir
        files:
          - prompts/system.txt
          - prompts/user.txt
```

#### Sources
- `key`: Unique identifier for the source
- `url`: Either a local directory path or a git repository URL
  - Local directories: `/path/to/directory` or `./relative/path`
  - Git repositories: `https://github.com/user/repo.git` or `git@github.com:user/repo.git`

**Special Sources:**
- `working_dir`: Automatically added to all configurations, pointing to the current working directory. This source is always available even if not explicitly defined in the YAML.

#### Targets
- `name`: Name of the target
- `output`: Directory where files will be downloaded/copied
- `include`: List of includes from sources
  - `source`: Reference to a source key (optional, defaults to `working_dir`)
  - `files`: List of file paths to include from that source

### Configuration Location
- Default: `pim.yaml` or `.pim.yaml` in the current directory
- Can be overridden with `--config` flag

## Features (To Be Defined)
- Package management
- Version control
- Dependency resolution
- Configuration management

## Future Work
- Define package structure
- Define repository format
- Add authentication mechanisms
- Add caching strategies
