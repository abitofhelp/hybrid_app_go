# Application vs Library Architecture Reference

**Version:** 1.0.0
**Date:** November 26, 2025
**SPDX-License-Identifier:** BSD-3-Clause
**Copyright:** (c) 2025 Michael Gardner, A Bit of Help, Inc.
**Status:** Released

Reference document for creating library starters (hybrid_lib_go, hybrid_lib_ada) based on
the application starters (hybrid_app_go, hybrid_app_ada).

---

## Layer Structure Comparison

| Layer | Application | Library | Notes |
|-------|:-----------:|:-------:|-------|
| **bootstrap** | Yes | No | Entry point, DI wiring |
| **presentation** | Yes | No | CLI/UI adapters |
| **api** | No | Yes | Public facade, re-exports Domain |
| **infrastructure** | Yes | Yes | Adapters, I/O (hidden in libs) |
| **application** | Yes | Yes | Use cases, ports, commands |
| **domain** | Yes | Yes | Pure business logic |

---

## Dependency Rules

### Application
```
bootstrap -> presentation -> infrastructure -> application -> domain
                                    |
                                    v
                              application -> domain
```

### Library
```
api -> infrastructure -> application -> domain
              |
              v
        application -> domain
```

**Key Difference**: Libraries expose only `api` + `domain` publicly via `Library_Interface`.
Application and Infrastructure are internal implementation details.

---

## Detection Logic

```python
def detect_project_type(project_root: Path) -> str:
    """Detect whether project is application or library."""
    # Ada projects
    if (project_root / 'src' / 'api').exists():
        return 'library'
    if (project_root / 'src' / 'bootstrap').exists():
        return 'application'

    # Go projects
    if (project_root / 'api').exists():
        return 'library'
    if (project_root / 'bootstrap').exists():
        return 'application'

    # Fallback: check for main/cmd
    if (project_root / 'cmd').exists():
        return 'application'

    return 'unknown'
```

---

## Port Terminology (Standardized)

Use **inbound/outbound** (not inward/outward):

| Port Type | Direction | Purpose | Example |
|-----------|-----------|---------|---------|
| **Inbound** | Outside -> In | Use case interfaces (how clients call us) | `FindById`, `Greet` |
| **Outbound** | Inside -> Out | Dependency interfaces (what we need) | `Writer`, `Repository` |

---

## Ada-Specific: SPARK Support

### Application Projects (hybrid_app_ada)
- **No SPARK_Mode pragmas** required
- Focus: Desktop/web services
- Standard Ada 2022 features

### Library Projects (hybrid_lib_ada, tzif)
- **SPARK_Mode => On** in Domain and Application layers
- **SPARK_Mode => Off** in Infrastructure (I/O, RNG)
- Focus: Reusable across desktop, web, and embedded systems
- Formal package pattern for I/O isolation

```ada
-- Domain layer: SPARK verified
package Domain.Parser is
   pragma SPARK_Mode (On);
   -- Pure parsing logic, no I/O
end Domain.Parser;

-- Infrastructure layer: Non-SPARK
package Infrastructure.IO.Desktop is
   pragma SPARK_Mode (Off);
   -- File I/O operations
end Infrastructure.IO.Desktop;
```

---

## Ada GPR Configuration

### Application (.gpr)
```ada
-- Single executable output
for Main use ("greeter.adb");

-- Layer GPRs for architectural enforcement
with "src/domain/domain.gpr";
with "src/application/application.gpr";
-- etc.
```

### Library (.gpr)
```ada
-- Library output
for Library_Name use "hybrid_lib_ada";
for Library_Standalone use "standard";

-- Public interface: Only API + Domain exposed
for Library_Interface use
  ("Hybrid_Lib_Ada",
   "Hybrid_Lib_Ada.API",
   "Domain",
   "Domain.Error",
   "Domain.Value_Object",
   -- Application.* and Infrastructure.* EXCLUDED
  );
```

---

## Go-Specific Notes

Go doesn't have "libraries" per se - just packages. However, structurally:

### Application (hybrid_app_go)
- Has `bootstrap/` and `presentation/` directories
- Produces executable binary
- `cmd/` or `main.go` entry point

### Library (hybrid_lib_go)
- Has `api/` directory instead of bootstrap/presentation
- No `main.go` or `cmd/`
- Exports packages for other Go modules to import
- `go.mod` module path is the "library name"

---

## API Layer Pattern (Libraries Only)

The API layer acts as a public facade:

```ada
-- api/hybrid_lib_ada-api.ads
package Hybrid_Lib_Ada.API is

   -- Re-export Domain types (public)
   subtype Person_Type is Domain.Value_Object.Person.Person;

   -- Re-export Application ports (public interface)
   function Greet (Cmd : Greet_Command) return Unit_Result.Result;

   -- Infrastructure is HIDDEN - not re-exported

end Hybrid_Lib_Ada.API;
```

---

## Build Profiles

### Application
Typically 4 modes:
- `development` - Debug symbols, all checks
- `test` - Coverage instrumentation
- `release` - Optimized, minimal checks
- `benchmark` - Performance testing

### Library (Ada with SPARK)
May have 6+ profiles for embedded targets:
- `standard` - Desktop/server (1+ GB RAM)
- `embedded` - Ravenscar (512KB+ RAM)
- `concurrent` - Multi-threaded
- `baremetal` - Zero footprint (128KB+ RAM)
- `stm32h7s78` - Specific board target
- `stm32mp135_linux` - Linux on ARM

---

## Tool Support Matrix

| Tool | App Support | Lib Support | Changes Needed |
|------|:-----------:|:-----------:|----------------|
| **arch_guard** | Yes | Planned | Add lib layer rules, project type detection |
| **brand_project** | Yes | Yes | Architecture-agnostic (no changes) |
| **release** | Yes | Yes | Architecture-agnostic (no changes) |

---

## arch_guard Layer Rules

### Application Rules
```python
APP_LAYER_RULES = {
    'bootstrap': {'can_import': ['presentation', 'infrastructure', 'application', 'domain']},
    'presentation': {'can_import': ['application']},  # NOT domain directly
    'infrastructure': {'can_import': ['application', 'domain']},
    'application': {'can_import': ['domain']},
    'domain': {'can_import': []},  # No external dependencies
}
```

### Library Rules
```python
LIB_LAYER_RULES = {
    'api': {'can_import': ['infrastructure', 'application', 'domain']},
    'infrastructure': {'can_import': ['application', 'domain']},
    'application': {'can_import': ['domain']},
    'domain': {'can_import': []},
}
```

---

## Migration Checklist: App to Lib Starter

### Go Migration (Simple)

Go library creation is straightforward - no SPARK or embedded concerns:

1. **Remove layers**:
   - [ ] Delete `bootstrap/` directory
   - [ ] Delete `presentation/` directory

2. **Add API layer**:
   - [ ] Create `api/` directory
   - [ ] Create facade package re-exporting Domain types

3. **Update go.mod**:
   - [ ] Remove main package / cmd directory
   - [ ] Ensure clean public package exports

4. **Update tests**:
   - [ ] Test via public API (black-box testing)

5. **Update documentation**:
   - [ ] Update architecture diagrams
   - [ ] Add API usage examples

**That's it for Go** - no build profile changes, no SPARK, no GPR files.

---

### Ada Migration (More Complex - SPARK + Embedded)

Ada libraries require additional work for SPARK formal verification and embedded system support:

1. **Remove layers**:
   - [ ] Delete `bootstrap/` directory
   - [ ] Delete `presentation/` directory

2. **Add API layer**:
   - [ ] Create `api/` directory structure
   - [ ] Create facade package re-exporting Domain types
   - [ ] Create operations package for use case wrappers
   - [ ] Create platform-specific instantiations (e.g., `api/desktop/`)

3. **Update GPR files**:
   - [ ] Change main .gpr from executable to library output
   - [ ] Add `Library_Standalone use "standard"`
   - [ ] Define `Library_Interface` (API + Domain only, exclude Application/Infrastructure)
   - [ ] Update layer .gpr files for library context

4. **Add SPARK support** (required for Ada libs):
   - [ ] Add `pragma SPARK_Mode (On)` to Domain layer packages
   - [ ] Add `pragma SPARK_Mode (On)` to Application layer packages
   - [ ] Add `pragma SPARK_Mode (Off)` to Infrastructure I/O packages
   - [ ] Use formal packages for I/O isolation
   - [ ] Ensure Result monad supports SPARK verification

5. **Add embedded build profiles**:
   - [ ] `standard` - Desktop/server
   - [ ] `embedded` - Ravenscar
   - [ ] `baremetal` - Zero footprint
   - [ ] Board-specific profiles as needed (STM32, etc.)

6. **Update tests**:
   - [ ] Test via public API only (black-box)
   - [ ] Add SPARK proof harnesses for Domain layer
   - [ ] Add GNATprove configuration

7. **Update documentation**:
   - [ ] Change architecture diagrams
   - [ ] Document SPARK verification status
   - [ ] Add API usage examples
   - [ ] Document embedded platform support

---

## References

- **App starter (Go)**: hybrid_app_go
- **App starter (Ada)**: hybrid_app_ada
- **Lib reference (Ada)**: tzif (SPARK-ready timezone library)
- **Real app example**: adafmt (needs updating to match current template)
