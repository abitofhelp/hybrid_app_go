#!/usr/bin/env python3
# ==============================================================================
# adapters/base.py - Base adapter for release management
# ==============================================================================
# Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.
# SPDX-License-Identifier: BSD-3-Clause
# See LICENSE file in the project root.
#
# Purpose:
#   Abstract base class for language-specific release adapters.
#   Provides common release operations and defines the interface for
#   language-specific implementations.
#
# Design Notes:
#   Following the same adapter pattern as brand_project and arch_guard.
#   Each language adapter implements version management, build, test,
#   and configuration file handling specific to that language.
#
# ==============================================================================

from abc import ABC, abstractmethod
from datetime import datetime
from pathlib import Path
from typing import List, Optional, Tuple
import re
import subprocess
import sys


class BaseReleaseAdapter(ABC):
    """
    Abstract base class for language-specific release operations.

    Subclasses must implement:
        - name: Human-readable language name
        - detect(): Check if project is this language
        - load_project_info(): Load project metadata from config files
        - update_version(): Update version in config files
        - run_build(): Execute build commands
        - run_tests(): Execute test commands
    """

    @property
    @abstractmethod
    def name(self) -> str:
        """Human-readable language name (e.g., 'Go', 'Ada')."""
        pass

    @staticmethod
    @abstractmethod
    def detect(project_root: Path) -> bool:
        """
        Detect if a directory is a project of this language.

        Args:
            project_root: Path to check

        Returns:
            True if project detected
        """
        pass

    @abstractmethod
    def load_project_info(self, config) -> Tuple[str, str]:
        """
        Load project name and URL from language-specific config.

        Args:
            config: ReleaseConfig instance

        Returns:
            Tuple of (project_name, project_url)
        """
        pass

    @abstractmethod
    def update_version(self, config) -> bool:
        """
        Update version in language-specific configuration files.

        Args:
            config: ReleaseConfig instance

        Returns:
            True if successful
        """
        pass

    @abstractmethod
    def run_build(self, config) -> bool:
        """
        Run language-specific build commands.

        Args:
            config: ReleaseConfig instance

        Returns:
            True if build successful
        """
        pass

    @abstractmethod
    def run_tests(self, config) -> bool:
        """
        Run language-specific test commands.

        Args:
            config: ReleaseConfig instance

        Returns:
            True if tests pass
        """
        pass

    def run_format(self, config) -> bool:
        """
        Run language-specific code formatting (optional).

        Default implementation does nothing.

        Args:
            config: ReleaseConfig instance

        Returns:
            True if successful
        """
        return True

    def sync_versions(self, config) -> bool:
        """
        Synchronize versions across project (optional).

        Default implementation does nothing.

        Args:
            config: ReleaseConfig instance

        Returns:
            True if successful
        """
        return True

    def generate_version_file(self, config) -> bool:
        """
        Generate version source file (optional).

        Default implementation does nothing.

        Args:
            config: ReleaseConfig instance

        Returns:
            True if successful
        """
        return True

    def cleanup_temp_files(self, config) -> bool:
        """
        Clean up temporary and build files.

        Default implementation runs 'make clean' if Makefile exists.

        Args:
            config: ReleaseConfig instance

        Returns:
            True if successful
        """
        makefile = config.project_root / 'Makefile'
        if makefile.exists():
            return self.run_command(['make', 'clean'], config.project_root) is not None
        return True

    def validate_makefile(self, config) -> bool:
        """
        Validate that all key Makefile targets work properly.

        Default implementation does basic validation. Language-specific
        adapters can override with comprehensive target lists.

        Args:
            config: ReleaseConfig instance

        Returns:
            True if all targets work
        """
        makefile = config.project_root / 'Makefile'
        if not makefile.exists():
            print("  No Makefile found, skipping validation")
            return True

        # Default targets to validate (override in subclass for language-specific)
        targets = ['help', 'build', 'clean']

        print("Validating Makefile targets...")
        for target in targets:
            result = self.run_command(
                ['make', target],
                config.project_root,
                capture_output=True,
                check=False
            )
            if result is None:
                print(f"  ✗ 'make {target}' failed")
                return False
            print(f"  ✓ make {target}")

        return True

    # =========================================================================
    # Common Operations (shared by all adapters)
    # =========================================================================

    def run_command(
        self,
        cmd: List[str],
        cwd: Path,
        capture_output: bool = False,
        check: bool = True
    ) -> Optional[str]:
        """
        Run a shell command.

        Args:
            cmd: Command as list of strings
            cwd: Working directory
            capture_output: Capture stdout/stderr
            check: Return None on non-zero exit (don't raise)

        Returns:
            stdout if capture_output, "SUCCESS" if succeeded, None on failure
        """
        try:
            result = subprocess.run(
                cmd,
                cwd=cwd,
                capture_output=capture_output,
                text=True
            )
            if result.returncode != 0:
                if check:
                    print(f"Command failed: {' '.join(cmd)}")
                    if result.stderr:
                        print(f"Error: {result.stderr}")
                return None
            return result.stdout if capture_output else "SUCCESS"
        except Exception as e:
            print(f"Command exception: {' '.join(cmd)}: {e}")
            return None

    def find_markdown_files(self, project_root: Path) -> List[Path]:
        """Find all markdown files with version headers."""
        md_files = []

        # Search in docs and root
        for pattern in ["docs/**/*.md", "*.md"]:
            md_files.extend(project_root.glob(pattern))

        # Filter to only files with version headers
        versioned_files = []
        for md_file in md_files:
            try:
                content = md_file.read_text(encoding='utf-8')
                if re.search(
                    r'Version\s*[:)]|version\s*[:)]|\*\*Version\s+\d+\.\d+|Copyright\s*©\s*\d{4}',
                    content, re.IGNORECASE
                ):
                    versioned_files.append(md_file)
            except Exception:
                pass

        return versioned_files

    def update_markdown_version(self, file_path: Path, config) -> bool:
        """
        Update version and metadata in markdown file headers.

        Handles these patterns:
        - **Version:** 1.0.0
        - **Date:** October 24, 2025
        - **Copyright:** (c) 2025 Michael Gardner, A Bit of Help, Inc.
        - **Status:** Unreleased / Released
        """
        try:
            content = file_path.read_text(encoding='utf-8')
            old_content = content

            # Pattern 1: **Version:** 1.0.0 (bold with colon)
            content = re.sub(
                r'(\*\*Version:\*\*\s+)[^\s\n]+',
                rf'\g<1>{config.version}',
                content,
                flags=re.IGNORECASE
            )

            # Pattern 2: Version: 1.0.0 (plain with colon)
            content = re.sub(
                r'^(Version:\s+)[^\s\n]+',
                rf'\g<1>{config.version}',
                content,
                flags=re.IGNORECASE | re.MULTILINE
            )

            # Pattern 3: **Version 1.0.0** (bold without colon)
            content = re.sub(
                r'(\*\*Version\s+)\d+\.\d+\.\d+(\*\*)',
                rf'\g<1>{config.version}\g<2>',
                content,
                flags=re.IGNORECASE
            )

            # Pattern 4: **Date:** October 24, 2025 (bold)
            content = re.sub(
                r'(\*\*Date:\*\*\s+)[^\n]+',
                rf'\g<1>{config.date_str}',
                content,
                flags=re.IGNORECASE
            )

            # Pattern 5: Date: October 24, 2025 (plain)
            content = re.sub(
                r'^(Date:\s+)[^\n]+',
                rf'\g<1>{config.date_str}',
                content,
                flags=re.IGNORECASE | re.MULTILINE
            )

            # Pattern 6: **Copyright:** (c) 2024 -> (c) 2025 (update year only)
            content = re.sub(
                r'(\*\*Copyright:\*\*\s+©\s*)\d{4}',
                rf'\g<1>{config.year}',
                content,
                flags=re.IGNORECASE
            )

            # Pattern 7: Copyright: (c) 2024 -> (c) 2025 (plain)
            content = re.sub(
                r'^(Copyright:\s+©\s*)\d{4}',
                rf'\g<1>{config.year}',
                content,
                flags=re.IGNORECASE | re.MULTILINE
            )

            # Pattern 8: **Status:** Unreleased -> Released
            if not config.is_prerelease:
                content = re.sub(
                    r'(\*\*Status:\*\*\s+)Unreleased',
                    r'\g<1>Released',
                    content,
                    flags=re.IGNORECASE
                )
                content = re.sub(
                    r'^(Status:\s+)Unreleased',
                    r'\g<1>Released',
                    content,
                    flags=re.IGNORECASE | re.MULTILINE
                )

            # Add trailing spaces for proper GitHub markdown rendering
            lines = content.split('\n')
            new_lines = []
            for line in lines:
                if re.match(r'^\*\*(Version|Date|SPDX|License|Copyright|Status):', line):
                    if not line.endswith('  '):
                        line = line.rstrip() + '  '
                new_lines.append(line)
            content = '\n'.join(new_lines)

            if content != old_content:
                file_path.write_text(content, encoding='utf-8')
                return True

            return False

        except Exception as e:
            print(f"Error updating {file_path}: {e}")
            return False

    def add_markdown_header(self, file_path: Path, config) -> bool:
        """Add metadata header to markdown file if missing."""
        try:
            content = file_path.read_text(encoding='utf-8')
            lines = content.splitlines(keepends=True)

            # Find first # heading
            title_idx = None
            for i, line in enumerate(lines):
                if re.match(r'^#\s+\S', line):
                    title_idx = i
                    break

            if title_idx is None:
                return False

            # Create header
            status = "Unreleased" if config.is_prerelease else "Released"
            header = (
                "\n"
                f"**Version:** {config.version}  \n"
                f"**Date:** {config.date_str}  \n"
                f"**SPDX-License-Identifier:** BSD-3-Clause  \n"
                f"**License File:** See the LICENSE file in the project root.  \n"
                f"**Copyright:** (c) {config.year} Michael Gardner, A Bit of Help, Inc.  \n"
                f"**Status:** {status}  \n"
                "\n"
            )

            lines.insert(title_idx + 1, header)
            file_path.write_text(''.join(lines), encoding='utf-8')
            return True

        except Exception as e:
            print(f"Error adding header to {file_path}: {e}")
            return False

    def update_all_markdown_files(self, config) -> int:
        """Update version in all markdown files. Returns count of updated files."""
        skip_files = {
            "software_requirements_specification.md",
            "software_design_specification.md",
            "software_test_guide.md"
        }

        all_md_files = []
        for pattern in ["docs/**/*.md", "*.md"]:
            all_md_files.extend(config.project_root.glob(pattern))

        updated_count = 0

        for md_file in all_md_files:
            if md_file.name in skip_files:
                continue

            try:
                content = md_file.read_text(encoding='utf-8')
                has_metadata = bool(re.search(
                    r'Version\s*[:)]|version\s*[:)]|\*\*Version\s+\d+\.\d+|Copyright\s*©\s*\d{4}',
                    content, re.IGNORECASE
                ))

                if has_metadata:
                    if self.update_markdown_version(md_file, config):
                        rel_path = md_file.relative_to(config.project_root)
                        print(f"  Updated {rel_path}")
                        updated_count += 1
                else:
                    if self.add_markdown_header(md_file, config):
                        rel_path = md_file.relative_to(config.project_root)
                        print(f"  Added header to {rel_path}")
                        updated_count += 1

            except Exception as e:
                print(f"  Warning: Error processing {md_file}: {e}")

        return updated_count

    def generate_diagrams(self, config) -> bool:
        """Generate PlantUML diagrams."""
        try:
            subprocess.run(
                ["plantuml", "-version"],
                capture_output=True,
                check=True
            )
        except (subprocess.CalledProcessError, FileNotFoundError):
            print("  plantuml not found, skipping diagram generation")
            return True

        diagrams_dir = config.project_root / "docs" / "diagrams"
        if not diagrams_dir.exists():
            print("  No diagrams directory found")
            return True

        puml_files = list(diagrams_dir.glob("*.puml"))
        if not puml_files:
            print("  No PlantUML files found")
            return True

        for puml_file in puml_files:
            self.run_command(
                ["plantuml", "-tsvg", str(puml_file)],
                config.project_root,
                capture_output=True
            )

        print(f"  Generated {len(puml_files)} diagram(s)")
        return True

    def verify_clean_working_tree(self, config) -> bool:
        """Verify git working tree is clean."""
        result = self.run_command(
            ["git", "status", "--porcelain"],
            config.project_root,
            capture_output=True
        )
        if result is None:
            return False
        return len(result.strip()) == 0

    def create_git_tag(self, config) -> bool:
        """Create annotated git tag."""
        tag_name = config.tag_name
        message = f"Release version {config.version}"

        result = self.run_command(
            ["git", "tag", "-a", tag_name, "-m", message],
            config.project_root
        )

        if result:
            print(f"  Created tag {tag_name}")
        return result is not None

    def push_changes(self, config) -> bool:
        """Push changes and tags to origin."""
        commands = [
            (["git", "push", "origin", "main"], "Pushed to main"),
            (["git", "push", "origin", config.tag_name], f"Pushed tag {config.tag_name}")
        ]

        for cmd, success_msg in commands:
            if self.run_command(cmd, config.project_root) is None:
                return False
            print(f"  {success_msg}")

        return True

    def create_github_release(self, config) -> bool:
        """Create GitHub release using gh CLI."""
        changelog_file = config.project_root / "CHANGELOG.md"
        release_notes = f"Release version {config.version}"

        if changelog_file.exists():
            try:
                content = changelog_file.read_text(encoding='utf-8')
                version_pattern = rf'## \[{re.escape(config.version)}\][^\n]*\n(.*?)(?=\n## |\Z)'
                match = re.search(version_pattern, content, re.DOTALL)
                if match:
                    release_notes = match.group(1).strip()
            except Exception as e:
                print(f"Warning: Could not extract release notes: {e}")

        cmd = [
            "gh", "release", "create", config.tag_name,
            "--title", f"Release {config.version}",
            "--notes", release_notes
        ]

        result = self.run_command(cmd, config.project_root)
        if result:
            print(f"  Created GitHub release {config.tag_name}")
        return result is not None
