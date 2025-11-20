#!/usr/bin/env python3
"""
Hexagonal/Clean Architecture Guard for Go Projects

Validates layer dependencies in a hybrid DDD/Clean/Hexagonal architecture.

Architecture Model (Concentric Spheres - Center-Seeking Dependencies):
- Domain: ZERO dependencies (innermost core)
- Application: Depends ONLY on Domain (middle sphere)
- Infrastructure: Depends on Application + Domain (outer half-sphere)
- Presentation: Depends ONLY on Application (NOT Domain directly!)
- Bootstrap: Depends on all layers (outermost)

Exit Codes:
  0: All architecture rules satisfied
  1: Architecture violations detected
  2: Script error
"""

import os
import sys
import re
from pathlib import Path
from typing import Set, Dict, List, Tuple
from dataclasses import dataclass


@dataclass
class ArchitectureViolation:
    """Represents a single architecture rule violation"""
    file_path: str
    line_number: int
    violation_type: str
    details: str


class ArchitectureGuard:
    """Validates hexagonal architecture layer dependencies for Go projects"""

    # Layer dependency rules (layer -> allowed dependencies)
    LAYER_RULES = {
        'domain': set(),
        'application': {'domain'},
        'infrastructure': {'application', 'domain'},
        'presentation': {'application'},
        'bootstrap': {'domain', 'application', 'infrastructure', 'presentation'}
    }

    FORBIDDEN_TEST_IMPORTS = ['testing', 'testify', 'assert', 'require', 'mock']
    FORBIDDEN_DEPENDENCY = 'bootstrap'

    def __init__(self, project_root: Path):
        self.project_root = project_root
        self.violations: List[ArchitectureViolation] = []
        self.layers_present = self._detect_layers()
        self.gomod_config_valid = False
        self.module_path = self._get_module_path()

    def _get_module_path(self) -> str:
        """Extract the module path from the root go.mod file"""
        root_gomod = self.project_root / 'go.mod'
        if root_gomod.exists():
            try:
                with open(root_gomod, 'r', encoding='utf-8') as f:
                    for line in f:
                        match = re.match(r'^\s*module\s+(.+)\s*$', line)
                        if match:
                            return match.group(1).strip()
            except Exception as e:
                print(f"Warning: Could not read root go.mod: {e}")
        return "github.com/abitofhelp/hybrid_app_go"

    def _detect_layers(self) -> Set[str]:
        """Detect which layers exist in the project"""
        layers = set()
        if not self.project_root.exists():
            print(f"Warning: Project root {self.project_root} does not exist")
            return layers

        print("Layer Detection:")
        for layer in sorted(self.LAYER_RULES.keys()):
            layer_dir = self.project_root / layer
            if layer_dir.exists() and layer_dir.is_dir():
                layers.add(layer)
                print(f"  ✓ {layer:15} - present")
            else:
                print(f"  ○ {layer:15} - not present (skipped)")

        return layers

    def _extract_imports(self, file_path: Path) -> List[Tuple[int, str]]:
        """Extract all import statements from a Go file"""
        imports = []

        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                in_import_block = False

                for line_num, line in enumerate(f, start=1):
                    stripped = line.strip()

                    if stripped.startswith('import ('):
                        in_import_block = True
                        continue

                    if in_import_block:
                        if stripped == ')':
                            in_import_block = False
                            continue

                        match = re.match(r'^\s*(?:\w+\s+)?"([^"]+)"\s*$', stripped)
                        if match:
                            imports.append((line_num, match.group(1)))
                    else:
                        match = re.match(r'^\s*import\s+(?:\w+\s+)?"([^"]+)"\s*$', stripped)
                        if match:
                            imports.append((line_num, match.group(1)))

        except Exception as e:
            print(f"Warning: Could not read {file_path}: {e}")

        return imports

    def _get_layer_from_import(self, import_path: str) -> str | None:
        """Determine which layer an import belongs to"""
        if not import_path.startswith(self.module_path):
            return None

        relative = import_path[len(self.module_path):].lstrip('/')
        parts = relative.split('/')
        if parts and parts[0] in self.LAYER_RULES:
            return parts[0]

        return None

    def _get_file_layer(self, file_path: Path) -> str | None:
        """Determine which layer a file belongs to"""
        try:
            relative_path = file_path.relative_to(self.project_root)
            parts = relative_path.parts
            if parts and parts[0] in self.LAYER_RULES:
                return parts[0]
        except ValueError:
            pass
        return None

    def _validate_gomod_config(self) -> bool:
        """Validate go.mod files are configured correctly"""
        valid = True

        print("\nValidating go.mod configurations:")

        for layer in sorted(self.layers_present):
            gomod_path = self.project_root / layer / 'go.mod'

            if not gomod_path.exists():
                print(f"  ⚠ {layer:15} - no go.mod file")
                continue

            layer_deps = set()
            try:
                with open(gomod_path, 'r', encoding='utf-8') as f:
                    content = f.read()

                    for match in re.finditer(r'require\s+([^\s]+)', content):
                        dep = match.group(1)
                        dep_layer = self._get_layer_from_import(dep)
                        if dep_layer:
                            layer_deps.add(dep_layer)

            except Exception as e:
                print(f"  ❌ {layer:15} - error reading go.mod: {e}")
                valid = False
                continue

            allowed_deps = self.LAYER_RULES[layer]
            forbidden = layer_deps - allowed_deps
            if forbidden:
                print(f"  ❌ {layer:15} - forbidden dependencies: {', '.join(sorted(forbidden))}")
                valid = False
            else:
                deps_str = ', '.join(sorted(layer_deps)) if layer_deps else 'none'
                print(f"  ✓ {layer:15} - dependencies: {deps_str}")

        return valid

    def _validate_no_test_imports(self, file_path: Path) -> None:
        """Ensure production code doesn't import test frameworks"""
        if file_path.name.endswith('_test.go'):
            return

        imports = self._extract_imports(file_path)

        for line_num, import_path in imports:
            if import_path == 'testing':
                self.violations.append(ArchitectureViolation(
                    file_path=str(file_path),
                    line_number=line_num,
                    violation_type='TEST_CODE_IN_PRODUCTION',
                    details=f"Production code cannot import testing package: {import_path}"
                ))

    def validate_file(self, file_path: Path) -> None:
        """Validate a single Go file against architecture rules"""
        self._validate_no_test_imports(file_path)

        current_layer = self._get_file_layer(file_path)
        if not current_layer or current_layer not in self.layers_present:
            return

        allowed_deps = self.LAYER_RULES[current_layer]
        imports = self._extract_imports(file_path)

        for line_num, import_path in imports:
            dependency_layer = self._get_layer_from_import(import_path)

            if not dependency_layer or dependency_layer not in self.layers_present:
                continue

            if dependency_layer == current_layer:
                continue

            # Check for bootstrap dependency violation
            if dependency_layer == self.FORBIDDEN_DEPENDENCY:
                self.violations.append(ArchitectureViolation(
                    file_path=str(file_path),
                    line_number=line_num,
                    violation_type='FORBIDDEN_BOOTSTRAP_DEPENDENCY',
                    details=f"Layer '{current_layer}' cannot depend on '{self.FORBIDDEN_DEPENDENCY}' (import: {import_path})"
                ))
                continue

            # Check for forbidden lateral dependencies
            if current_layer == 'presentation' and dependency_layer == 'infrastructure':
                self.violations.append(ArchitectureViolation(
                    file_path=str(file_path),
                    line_number=line_num,
                    violation_type='FORBIDDEN_LATERAL_DEPENDENCY',
                    details=f"Presentation cannot depend on Infrastructure (import: {import_path})"
                ))
                continue

            if current_layer == 'infrastructure' and dependency_layer == 'presentation':
                self.violations.append(ArchitectureViolation(
                    file_path=str(file_path),
                    line_number=line_num,
                    violation_type='FORBIDDEN_LATERAL_DEPENDENCY',
                    details=f"Infrastructure cannot depend on Presentation (import: {import_path})"
                ))
                continue

            # CRITICAL: Presentation cannot import Domain directly
            if current_layer == 'presentation' and dependency_layer == 'domain':
                self.violations.append(ArchitectureViolation(
                    file_path=str(file_path),
                    line_number=line_num,
                    violation_type='PRESENTATION_IMPORTS_DOMAIN',
                    details=f"Presentation MUST NOT import Domain directly (import: {import_path})\n" +
                            f"      → Use application/error re-exports instead (e.g., apperr.Result[T])"
                ))
                continue

            # Check if inter-layer dependency is allowed
            if dependency_layer not in allowed_deps:
                self.violations.append(ArchitectureViolation(
                    file_path=str(file_path),
                    line_number=line_num,
                    violation_type='ILLEGAL_LAYER_DEPENDENCY',
                    details=f"Layer '{current_layer}' cannot depend on '{dependency_layer}' (import: {import_path})"
                ))

    def validate_all(self) -> bool:
        """Validate all Go files in the project"""
        if not self.layers_present:
            print("⚠ No architecture layers detected - skipping validation")
            return True

        print(f"\nValidating architecture rules for layers: {', '.join(sorted(self.layers_present))}\n")

        # Step 1: Validate go.mod configuration
        print("=" * 70)
        print("Step 1: Validate go.mod Configuration")
        print("=" * 70)
        self.gomod_config_valid = self._validate_gomod_config()
        print()

        # Step 2: Validate Go source file dependencies
        print("=" * 70)
        print("Step 2: Validate Go Source File Dependencies")
        print("=" * 70)

        go_files = []
        for layer in self.layers_present:
            layer_path = self.project_root / layer
            for go_file in layer_path.rglob('*.go'):
                if 'vendor' in go_file.parts:
                    continue
                go_files.append(go_file)

        print(f"Scanning {len(go_files)} Go files...\n")

        for go_file in go_files:
            self.validate_file(go_file)

        return self.gomod_config_valid and len(self.violations) == 0

    def report_violations(self) -> None:
        """Print violation report"""
        print("\n" + "=" * 70)
        print("FINAL RESULTS")
        print("=" * 70)

        if self.gomod_config_valid:
            print("✅ go.mod Configuration: VALID")
        else:
            print("❌ go.mod Configuration: INVALID")

        if not self.violations:
            print("✅ Source File Dependencies: VALID")
        else:
            print(f"❌ Source File Dependencies: {len(self.violations)} violation(s)")

        if self.gomod_config_valid and not self.violations:
            print("\n✅ Architecture validation PASSED - All rules satisfied!")
            return

        print(f"\n❌ Architecture validation FAILED")

        if self.violations:
            print(f"\nSource Dependency Violations ({len(self.violations)}):\n")

        by_type: Dict[str, List[ArchitectureViolation]] = {}
        for v in self.violations:
            by_type.setdefault(v.violation_type, []).append(v)

        for violation_type, violations in sorted(by_type.items()):
            print(f"  [{violation_type}] ({len(violations)} violations)")
            for v in violations:
                print(f"    {v.file_path}:{v.line_number}")
                print(f"      → {v.details}")
            print()


def main():
    """Main entry point"""
    script_dir = Path(__file__).parent
    project_root = script_dir.parent.parent

    print("=" * 70)
    print("Hexagonal Architecture Guard for Go")
    print("=" * 70)
    print(f"Project root: {project_root}")
    print()

    if not project_root.exists():
        print(f"ERROR: Project directory not found: {project_root}")
        return 2

    guard = ArchitectureGuard(project_root)

    if not guard.layers_present:
        print("No architecture layers to validate - exiting")
        return 2

    is_valid = guard.validate_all()
    guard.report_violations()

    return 0 if is_valid else 1


if __name__ == '__main__':
    sys.exit(main())
