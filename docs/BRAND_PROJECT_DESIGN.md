# brand_project.py - Unified Project Instantiation Script Design

**Status:** Design complete, implementation pending
**Date:** 2025-11-25

## Overview

A unified script (like arch_guard.py) that instantiates new projects from hybrid_app or hybrid_lib templates for Go, Ada, and Rust.

## CLI Interface

```bash
python3 brand_project.py --git-repo https://github.com/abitofhelp/my_awesome_app.git

# Or shorthand (script adds https:// and .git):
python3 brand_project.py --git-repo github.com/abitofhelp/my_awesome_app
```

## Parsing the Git Repo URL

```
https://github.com/abitofhelp/my_awesome_app.git
         └──────┬──────┘ └────┬────┘ └─────┬─────┘
              host         account    project_name
```

**Derived values:**
- **Project name**: `my_awesome_app`
- **Go module**: `github.com/abitofhelp/my_awesome_app`
- **Ada alire.toml**: `name = "my_awesome_app"`, `website = "https://github.com/abitofhelp/my_awesome_app"`

## Unified Process (All Languages)

| Step | Action |
|------|--------|
| 1 | Detect template language (Go/Ada/Rust) from source |
| 2 | Create target directory |
| 3 | Copy template files (skip build artifacts) |
| 4 | Rename files containing template name |
| 5 | Search/replace template name in file contents |
| 6 | Update language-specific config (go.mod, alire.toml, Cargo.toml) |
| 7 | Verify no old references remain |

## Language-Specific Config Updates

| Language | File | Field | Value |
|----------|------|-------|-------|
| **Go** | `go.mod` | `module` | `github.com/xxxxx/yyyyy` |
| **Go** | `go.work` | `use` paths | Updated references |
| **Ada** | `alire.toml` | `name` | `my_awesome_app` |
| **Ada** | `alire.toml` | `website` | `https://github.com/xxxxx/yyyyy` |
| **Ada** | `*.gpr` | project name | `My_Awesome_App` (Pascal_Case) |
| **Rust** | `Cargo.toml` | `name` | `my_awesome_app` |
| **Rust** | `Cargo.toml` | `repository` | `https://github.com/xxxxx/yyyyy` |

## Excluded from Copy

```python
EXCLUDED_DIRS = {
    '.git',           # All languages
    'alire', 'obj', 'bin', 'lib',  # Ada
    'vendor',         # Go
    'target',         # Rust
    '__pycache__', '.pytest_cache',  # Python tooling
}

EXCLUDED_FILES = {
    '*.gz',
    '*.zip',
}
```

## Directory Structure

```
scripts/brand_project/
├── __init__.py
├── __main__.py           # python -m brand_project
├── brand_project.py      # Main engine
├── models.py             # GitRepoUrl, ProjectConfig
└── adapters/
    ├── base.py           # BaseAdapter
    ├── go.py             # GoAdapter
    ├── ada.py            # AdaAdapter
    └── rust.py           # RustAdapter (future)
```

## Case Conventions by Language

| Language | snake_case | PascalCase | UPPER_CASE |
|----------|------------|------------|------------|
| **Go** | `my_app` | `MyApp` | `MY_APP` |
| **Ada** | `my_app` | `My_App` | `MY_APP` |
| **Rust** | `my_app` | `MyApp` | `MY_APP` |
