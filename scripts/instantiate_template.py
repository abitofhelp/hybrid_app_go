#!/usr/bin/env python3
# ==============================================================================
# instantiate_template.py - Ada Project Template Instantiation Script
# ==============================================================================
# Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.
# SPDX-License-Identifier: BSD-3-Clause
# ==============================================================================
"""
Instantiate a new Ada project from the hybrid_app_ada template.

This script performs the following transformations:
1. Renames files containing the old project name
2. Replaces all occurrences of the old name in file contents
3. Handles three case variations: snake_case, Pascal_Case, UPPER_CASE

Usage:
    python3 instantiate_template.py --old-name hybrid_app_ada --new-name my_app
    python3 instantiate_template.py --old-name hybrid_app_ada --new-name my_app --dry-run

Arguments:
    --old-name    : Current project name (snake_case)
    --new-name    : New project name (snake_case)
    --dry-run     : Show what would be changed without making changes
    --project-dir : Project directory (default: current directory)
"""

import argparse
import re
import subprocess
import sys
from pathlib import Path
from typing import List, Tuple


class NameConverter:
    """Convert between different naming conventions."""

    @staticmethod
    def to_pascal_case(snake_case: str) -> str:
        """Convert snake_case to PascalCase (e.g., my_app → My_App)."""
        return '_'.join(word.capitalize() for word in snake_case.split('_'))

    @staticmethod
    def to_upper_case(snake_case: str) -> str:
        """Convert snake_case to UPPER_CASE."""
        return snake_case.upper()


class ProjectInstantiator:
    """Instantiate a new Ada project from a template."""

    # File extensions to process for content replacement
    PROCESSABLE_EXTENSIONS = {
        '.ads', '.adb', '.gpr', '.h', '.c',
        '.md', '.toml', '.py', '.sh', '.bash',
        '.puml', '.svg', '.xml', '.yaml', '.yml',
        'Makefile', '.gitignore'
    }

    # Directories to exclude from processing
    EXCLUDED_DIRS = {
        '.git', 'obj', 'bin', 'lib', 'alire',
        'external', '.pytest_cache', '__pycache__'
    }

    def __init__(
        self,
        old_name: str,
        new_name: str,
        project_dir: Path,
        dry_run: bool = False
    ):
        """
        Initialize the instantiator.

        Args:
            old_name: Current project name (snake_case)
            new_name: New project name (snake_case)
            project_dir: Project directory
            dry_run: If True, show changes without applying them
        """
        self.old_name = old_name
        self.new_name = new_name
        self.project_dir = project_dir.resolve()
        self.dry_run = dry_run

        # Generate all case variations
        self.old_pascal = NameConverter.to_pascal_case(old_name)
        self.new_pascal = NameConverter.to_pascal_case(new_name)
        self.old_upper = NameConverter.to_upper_case(old_name)
        self.new_upper = NameConverter.to_upper_case(new_name)

        self.files_renamed: List[Tuple[Path, Path]] = []
        self.files_modified: List[Path] = []

    def log(self, message: str, prefix: str = "ℹ️"):
        """Log a message to stdout."""
        print(f"{prefix} {message}")

    def is_processable_file(self, file_path: Path) -> bool:
        """Check if a file should be processed for content replacement."""
        # Check if file has a processable extension
        if file_path.name == 'Makefile':
            return True
        return file_path.suffix in self.PROCESSABLE_EXTENSIONS

    def is_excluded_path(self, path: Path) -> bool:
        """Check if a path should be excluded from processing."""
        return any(excluded in path.parts for excluded in self.EXCLUDED_DIRS)

    def find_files_to_rename(self) -> List[Tuple[Path, Path]]:
        """Find all files that need to be renamed."""
        files_to_rename = []

        for file_path in self.project_dir.rglob('*'):
            if not file_path.is_file():
                continue
            if self.is_excluded_path(file_path):
                continue

            # Check if filename contains old name
            if self.old_name in file_path.name:
                new_name = file_path.name.replace(self.old_name, self.new_name)
                new_path = file_path.parent / new_name
                files_to_rename.append((file_path, new_path))

        return files_to_rename

    def find_files_for_content_replacement(self) -> List[Path]:
        """Find all files that need content replacement."""
        files_to_process = []

        for file_path in self.project_dir.rglob('*'):
            if not file_path.is_file():
                continue
            if self.is_excluded_path(file_path):
                continue
            if not self.is_processable_file(file_path):
                continue

            files_to_process.append(file_path)

        return files_to_process

    def rename_files(self) -> None:
        """Rename all files containing the old project name."""
        files_to_rename = self.find_files_to_rename()

        if not files_to_rename:
            self.log("No files to rename", "✓")
            return

        self.log(f"Found {len(files_to_rename)} files to rename")

        for old_path, new_path in files_to_rename:
            rel_old = old_path.relative_to(self.project_dir)
            rel_new = new_path.relative_to(self.project_dir)

            if self.dry_run:
                self.log(f"Would rename: {rel_old} → {rel_new}", "📝")
            else:
                try:
                    # Use git mv if in a git repository
                    result = subprocess.run(
                        ['git', 'mv', str(old_path), str(new_path)],
                        cwd=self.project_dir,
                        capture_output=True,
                        text=True
                    )
                    if result.returncode == 0:
                        self.log(f"Renamed (git): {rel_old} → {rel_new}", "✓")
                        self.files_renamed.append((old_path, new_path))
                    else:
                        # Fallback to regular rename
                        old_path.rename(new_path)
                        self.log(f"Renamed: {rel_old} → {rel_new}", "✓")
                        self.files_renamed.append((old_path, new_path))
                except Exception as e:
                    self.log(f"Failed to rename {rel_old}: {e}", "❌")

    def replace_content(self) -> None:
        """Replace all occurrences of old names in file contents."""
        files_to_process = self.find_files_for_content_replacement()

        if not files_to_process:
            self.log("No files to process for content replacement", "✓")
            return

        self.log(f"Found {len(files_to_process)} files to process")

        for file_path in files_to_process:
            try:
                # Read file content
                content = file_path.read_text(encoding='utf-8')
                original_content = content

                # Replace all three case variations
                # Order matters: do most specific (longest) replacements first
                content = content.replace(self.old_upper, self.new_upper)
                content = content.replace(self.old_pascal, self.new_pascal)
                content = content.replace(self.old_name, self.new_name)

                # Check if content changed
                if content != original_content:
                    rel_path = file_path.relative_to(self.project_dir)

                    if self.dry_run:
                        self.log(f"Would modify: {rel_path}", "📝")
                    else:
                        file_path.write_text(content, encoding='utf-8')
                        self.log(f"Modified: {rel_path}", "✓")
                        self.files_modified.append(file_path)

            except UnicodeDecodeError:
                # Skip binary files
                pass
            except Exception as e:
                rel_path = file_path.relative_to(self.project_dir)
                self.log(f"Failed to process {rel_path}: {e}", "❌")

    def verify_changes(self) -> bool:
        """Verify that no old name references remain in source files."""
        self.log("\n" + "=" * 70)
        self.log("Verifying changes...")
        self.log("=" * 70)

        # Search for any remaining old name references
        try:
            result = subprocess.run(
                [
                    'grep', '-r',
                    f'{self.old_name}\\|{self.old_pascal}\\|{self.old_upper}',
                    '.',
                    '--include=*.ads',
                    '--include=*.adb',
                    '--include=*.gpr',
                    '--include=*.md',
                    '--include=*.toml',
                    '--include=*.py',
                    '--exclude-dir=.git',
                    '--exclude-dir=obj',
                    '--exclude-dir=bin',
                    '--exclude-dir=alire',
                    '--exclude-dir=external'
                ],
                cwd=self.project_dir,
                capture_output=True,
                text=True
            )

            if result.returncode == 0 and result.stdout.strip():
                self.log("⚠️  Warning: Found remaining old name references:", "⚠️")
                print(result.stdout)
                return False
            else:
                self.log("No old name references found in source files", "✓")
                return True

        except Exception as e:
            self.log(f"Verification failed: {e}", "❌")
            return False

    def print_summary(self) -> None:
        """Print a summary of changes."""
        self.log("\n" + "=" * 70)
        self.log("INSTANTIATION SUMMARY")
        self.log("=" * 70)

        mode = "DRY RUN" if self.dry_run else "APPLIED"
        self.log(f"Mode: {mode}")
        self.log(f"Project directory: {self.project_dir}")
        self.log(f"Old name: {self.old_name} / {self.old_pascal} / {self.old_upper}")
        self.log(f"New name: {self.new_name} / {self.new_pascal} / {self.new_upper}")
        self.log("")
        self.log(f"Files renamed: {len(self.files_renamed)}")
        self.log(f"Files modified: {len(self.files_modified)}")

        if not self.dry_run:
            self.log("\n" + "=" * 70)
            self.log("NEXT STEPS")
            self.log("=" * 70)
            self.log("1. Review changes: git status")
            self.log("2. Clean build artifacts: make clean")
            self.log("3. Build project: make build")
            self.log("4. Run tests: make test")
            self.log("5. If everything works, commit: git add -A && git commit")

    def run(self) -> int:
        """
        Run the instantiation process.

        Returns:
            0 on success, non-zero on failure
        """
        try:
            self.log("=" * 70)
            self.log("ADA PROJECT TEMPLATE INSTANTIATION")
            self.log("=" * 70)

            if not self.project_dir.exists():
                self.log(f"Project directory does not exist: {self.project_dir}", "❌")
                return 1

            if self.dry_run:
                self.log("Running in DRY RUN mode (no changes will be made)", "ℹ️")

            # Step 1: Rename files
            self.log("\n" + "=" * 70)
            self.log("Step 1: Renaming files")
            self.log("=" * 70)
            self.rename_files()

            # Step 2: Replace content
            self.log("\n" + "=" * 70)
            self.log("Step 2: Replacing content")
            self.log("=" * 70)
            self.replace_content()

            # Step 3: Verify (only if not dry run)
            if not self.dry_run:
                self.verify_changes()

            # Print summary
            self.print_summary()

            return 0

        except Exception as e:
            self.log(f"Instantiation failed: {e}", "❌")
            import traceback
            traceback.print_exc()
            return 1


def validate_name(name: str) -> bool:
    """Validate that a name is in snake_case format."""
    return bool(re.match(r'^[a-z][a-z0-9_]*$', name))


def main() -> int:
    """Main entry point."""
    parser = argparse.ArgumentParser(
        description='Instantiate a new Ada project from the hybrid_app_ada template',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog=__doc__
    )

    parser.add_argument(
        '--old-name',
        required=True,
        help='Current project name (snake_case, e.g., hybrid_app_ada)'
    )

    parser.add_argument(
        '--new-name',
        required=True,
        help='New project name (snake_case, e.g., my_awesome_app)'
    )

    parser.add_argument(
        '--project-dir',
        type=Path,
        default=Path.cwd(),
        help='Project directory (default: current directory)'
    )

    parser.add_argument(
        '--dry-run',
        action='store_true',
        help='Show what would be changed without making changes'
    )

    args = parser.parse_args()

    # Validate names
    if not validate_name(args.old_name):
        print(f"❌ Error: old-name must be in snake_case: {args.old_name}")
        return 1

    if not validate_name(args.new_name):
        print(f"❌ Error: new-name must be in snake_case: {args.new_name}")
        return 1

    # Run instantiation
    instantiator = ProjectInstantiator(
        old_name=args.old_name,
        new_name=args.new_name,
        project_dir=args.project_dir,
        dry_run=args.dry_run
    )

    return instantiator.run()


if __name__ == '__main__':
    sys.exit(main())
