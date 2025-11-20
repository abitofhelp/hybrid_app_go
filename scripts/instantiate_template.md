# Hybrid App Ada - Application Template

This is a template for creating Ada applications using a hybrid DDD/Clean/Hexagonal architecture.

## Using This Template

### Quick Start

To create a new application from this template:

```bash
# 1. Copy the template to a new directory
cp -r hybrid_app_ada my_new_app
cd my_new_app

# 2. Initialize git repository
git init
git add -A
git commit -m "Initial commit from hybrid_app_ada template"

# 3. Run the instantiation script
python3 scripts/instantiate_template.py \
  --old-name hybrid_app_ada \
  --new-name my_new_app

# 4. Clean and build
make clean
make build

# 5. Run tests
make test

# 6. Commit the changes
git add -A
git commit -m "Instantiate template with project name my_new_app"
```

### Instantiation Script Usage

The `scripts/instantiate_template.py` script automates the renaming process:

```bash
# Basic usage
python3 scripts/instantiate_template.py \
  --old-name hybrid_app_ada \
  --new-name your_app_name

# Dry run (see what would change without applying)
python3 scripts/instantiate_template.py \
  --old-name hybrid_app_ada \
  --new-name your_app_name \
  --dry-run

# Specify project directory
python3 scripts/instantiate_template.py \
  --old-name hybrid_app_ada \
  --new-name your_app_name \
  --project-dir /path/to/project
```

**Important:** Project names must be in `snake_case` format (lowercase with underscores).

### What Gets Renamed

The script performs the following transformations:

1. **File Names**: Renames all files containing the old project name
   - `hybrid_app_ada.gpr` → `your_app_name.gpr`
   - `config/hybrid_app_ada_config.ads` → `config/your_app_name_config.ads`

2. **File Contents**: Replaces all occurrences in three case variations:
   - `hybrid_app_ada` → `your_app_name` (snake_case)
   - `Hybrid_App_Ada` → `Your_App_Name` (Pascal_Case)
   - `HYBRID_APP_ADA` → `YOUR_APP_NAME` (UPPER_CASE)

3. **Processed File Types**:
   - Ada source files (`.ads`, `.adb`)
   - GNAT project files (`.gpr`)
   - Configuration files (`.toml`, `.yaml`)
   - Documentation (`.md`)
   - Build scripts (`.py`, `.sh`, `Makefile`)
   - Diagrams (`.puml`, `.svg`)

### After Instantiation

1. **Update `alire.toml`**:
   - Modify `description` field
   - Update `website` URL
   - Change `authors` and `maintainers`
   - Adjust `tags` as needed

2. **Update Documentation**:
   - Review and update `README.md`
   - Update `CHANGELOG.md` with your project's history
   - Modify `docs/` files to reflect your application

3. **Customize Implementation**:
   - Replace the greeter example with your domain logic
   - Implement your use cases in `src/application/usecase/`
   - Add your domain entities in `src/domain/`
   - Implement infrastructure adapters in `src/infrastructure/adapter/`

4. **Configure GitHub**:
   - Create a new repository on GitHub
   - Update the remote URL:
     ```bash
     git remote add origin https://github.com/yourusername/your_app_name.git
     git push -u origin main
     ```

## Architecture Overview

This template implements a hybrid architecture combining:
- **Domain-Driven Design (DDD)**: Rich domain model in `src/domain/`
- **Clean Architecture**: Dependency rules with domain at the center
- **Hexagonal Architecture**: Ports and adapters pattern in `src/application/port/`

### Directory Structure

```
src/
├── domain/              # Domain entities and business rules (innermost)
│   ├── error/          # Domain error types
│   └── value_object/   # Value objects
├── application/         # Application use cases and interfaces
│   ├── port/
│   │   ├── inbound/    # Input ports (use case interfaces)
│   │   └── outbound/   # Output ports (infrastructure interfaces)
│   ├── usecase/        # Use case implementations
│   ├── error/          # Application error types
│   └── model/          # Application-specific models
├── infrastructure/      # External interface implementations
│   └── adapter/        # Concrete implementations of outbound ports
├── presentation/        # User interface layer
│   └── cli/            # Command-line interface
└── bootstrap/          # Application startup and wiring
```

### Dependency Rules

- **Domain**: No dependencies on other layers
- **Application**: Depends only on Domain
- **Infrastructure**: Depends on Application and Domain
- **Presentation**: Depends on Application (not Infrastructure or Domain)
- **Bootstrap**: Wires everything together

These rules are enforced by `scripts/makefile/arch_guard.py` during build.

## Key Features

- ✅ **Hexagonal Architecture** with enforced dependency rules
- ✅ **Functional Error Handling** using Result<T,E> and Option<T> monads
- ✅ **Type Safety** with strong typing and contracts
- ✅ **Build Modes**: debug, release, optimize
- ✅ **Comprehensive Testing**: unit, integration, e2e test suites
- ✅ **Architecture Validation**: Automated with Python tests
- ✅ **Code Quality**: Style checks, pragma enforcement
- ✅ **Modern Ada**: Ada 2022 features
- ✅ **Alire Integration**: Package management with alire.toml

## Build System

```bash
make help           # Show all available commands
make build          # Build in development mode
make test           # Run all tests
make test-unit      # Run unit tests only
make test-python    # Run architecture validation tests
make clean          # Clean build artifacts
make check-arch     # Validate architecture boundaries
```

## Dependencies

- **functional**: Result/Option/Either types for functional error handling
- **aunit**: Unit testing framework

## License

BSD-3-Clause (same as template)

## Support

For issues with the template itself, see https://github.com/abitofhelp/hybrid_app_ada

For your application-specific questions, update this section with your support information.
