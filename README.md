# PIM

**PIM (Prompt Instruction Manager)** is a command-line utility for managing prompt instructions and related files from multiple sources.

## Features

- 📦 Fetch files from multiple sources (local directories, Git repositories)
- 🔧 Flexible configuration using YAML
- 🎯 Organize files into different targets
- 🚀 Automatic `working_dir` source for quick setups
- 📝 Simple and intuitive configuration format

## Installation

### From Source

```bash
git clone https://github.com/hubble-works/pim.git
cd pim
make install
```

This will install the `pim` binary to `$GOPATH/bin` (usually `~/go/bin`).

### Build Locally

```bash
make build
```

This creates the `pim` binary in the current directory.

## Quick Start

1. Create a `pim.yaml` configuration file:

```yaml
version: 1

targets:
  - name: prompts
    output: ./output
    include:
      - files:
          - prompts/system.txt
          - prompts/user.txt
```

2. Run PIM:

```bash
pim install
```

## Usage

### Commands

- `pim install [directory]` - Fetch files from sources to targets (defaults to current directory)
- `pim version` - Print version information
- `pim help` - Show help

### Configuration

PIM looks for `pim.yaml` or `.pim.yaml` in the current directory (or the directory specified as an argument).

#### Basic Configuration

```yaml
version: 1

sources:
  - key: local-prompts
    url: /path/to/prompts
  - key: shared-repo
    url: https://github.com/user/prompts-repo.git

targets:
  - name: my-project
    output: ./prompts
    include:
      - source: local-prompts
        files:
          - system.txt
          - user.txt
      - source: shared-repo
        files:
          - templates/common.txt
```

#### Minimal Configuration

The `working_dir` source is automatically available and points to the current directory:

```yaml
version: 1

targets:
  - name: local-files
    output: ./output
    include:
      - files:  # source defaults to working_dir
          - file1.txt
          - file2.txt
```

### Configuration Options

**Sources:**
- `key` - Unique identifier for the source
- `url` - Local directory path or Git repository URL
  - Local: `/absolute/path` or `./relative/path`
  - Git: `https://github.com/user/repo.git`

**Special Sources:**
- `working_dir` - Automatically added, points to current working directory

**Targets:**
- `name` - Target name
- `output` - Directory where files will be copied
- `include` - List of files to include
  - `source` - Source key (optional, defaults to `working_dir`)
  - `files` - List of file paths to include

## Development

### Running Tests

```bash
make test          # Run all tests
make test-verbose  # Run tests with verbose output
```

### Building

```bash
make build  # Build the binary
make clean  # Remove build artifacts
```

## License

See [LICENSE](LICENSE) file for details.

## Documentation

For detailed specification, see [SPEC.md](SPEC.md). 
