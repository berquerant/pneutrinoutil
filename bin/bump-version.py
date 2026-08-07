#!/usr/bin/env python3
"""Bump version script for updating project configuration files and optionally creating a PR."""

import argparse
import logging
import re
import subprocess
import sys
from dataclasses import dataclass
from pathlib import Path

SEMVER_REGEX = r"^v?(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-[0-9A-Za-z.-]+)?(\+[0-9A-Za-z.-]+)?$"


@dataclass
class CLIArgs:
    """Command line arguments for the version bump script."""

    version: str
    create_pr: bool = False
    dry_run: bool = False


@dataclass
class VersionInfo:
    """Parsed semantic version details."""

    clean_version: str
    tag_version: str


@dataclass
class ReplacementRule:
    """Regex pattern and replacement text for modifying configuration files."""

    pattern: str
    replacement: str
    flags: int = 0


@dataclass
class FileUpdateConfig:
    """Target file path and corresponding replacement rules."""

    filepath: Path
    rules: list[ReplacementRule]


class BumpVersionError(Exception):
    """Custom exception raised when version bumping encounters an error."""

    pass


def parse_args() -> CLIArgs:
    """Parse command line arguments.

    Returns:
        CLIArgs containing the parsed input version and flags.
    """
    parser = argparse.ArgumentParser(
        description="Bump version in configuration files and optionally create a Pull Request."
    )
    parser.add_argument(
        "version",
        help="Target version number to bump (e.g., 0.10.0 or v0.10.0)",
    )
    parser.add_argument(
        "--create-pr",
        action="store_true",
        help="Create a git branch, commit changes, push to origin, and open a Pull Request using gh CLI.",
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Dry run mode. Print actions/commands instead of executing PR creation.",
    )
    args = parser.parse_args()
    return CLIArgs(version=args.version, create_pr=args.create_pr, dry_run=args.dry_run)


def run_command(
    args: list[str], capture_output: bool = False
) -> subprocess.CompletedProcess:
    """Execute a system command with standard error output redirection.

    Args:
        args: List of command arguments.
        capture_output: Whether to capture standard output as text.

    Returns:
        CompletedProcess instance resulting from command execution.
    """
    if capture_output:
        return subprocess.run(args, capture_output=True, text=True)
    return subprocess.run(args, stdout=sys.stderr, stderr=sys.stderr)


def validate_semver(input_version: str) -> VersionInfo:
    """Validate that the input string matches SemVer format.

    Args:
        input_version: Version string provided by user.

    Returns:
        VersionInfo object with clean and tagged version strings.

    Raises:
        BumpVersionError: If the input version string is invalid SemVer.
    """
    if not re.match(SEMVER_REGEX, input_version):
        raise BumpVersionError(f"'{input_version}' is not a valid SemVer.")

    clean_version = input_version.lstrip("v")
    tag_version = f"v{clean_version}"
    logging.info(f"Valid SemVer format: {clean_version} (Git tag: {tag_version})")
    return VersionInfo(clean_version=clean_version, tag_version=tag_version)


def check_git_tag_availability(version_info: VersionInfo) -> None:
    """Check if the target Git tag already exists locally or on remote.

    Args:
        version_info: VersionInfo containing target version strings.

    Raises:
        BumpVersionError: If the tag exists locally or on remote 'origin'.
    """
    run_command(["git", "fetch", "--tags"])

    for tag in (version_info.tag_version, version_info.clean_version):
        res = run_command(["git", "rev-parse", "-q", "--verify", f"refs/tags/{tag}"])
        if res.returncode == 0:
            raise BumpVersionError(f"Git tag '{tag}' already exists locally!")

    ls_remote = run_command(
        ["git", "ls-remote", "--tags", "origin"], capture_output=True
    )
    if ls_remote.returncode == 0:
        for line in ls_remote.stdout.splitlines():
            if any(
                line.endswith(f"refs/tags/{t}")
                for t in (version_info.tag_version, version_info.clean_version)
            ):
                raise BumpVersionError(
                    f"Git tag '{version_info.tag_version}' or '{version_info.clean_version}' already exists on remote 'origin'!"
                )

    logging.info("Tag availability check passed.")


def replace_in_file(config: FileUpdateConfig) -> None:
    """Apply replacement rules to a target file.

    Args:
        config: FileUpdateConfig specifying the file path and replacement rules.
    """
    content = config.filepath.read_text(encoding="utf-8")
    for rule in config.rules:
        content = re.sub(rule.pattern, rule.replacement, content, flags=rule.flags)
    config.filepath.write_text(content, encoding="utf-8")
    logging.info(f"Updated {config.filepath}")


def update_version_files(clean_version: str) -> None:
    """Update all project configuration files with the new version string.

    Args:
        clean_version: The new version string without 'v' prefix.
    """
    logging.info(f"Updating version to {clean_version} in configuration files...")

    configs = [
        FileUpdateConfig(
            filepath=Path("pyproject.toml"),
            rules=[
                ReplacementRule(
                    pattern=r'^(version\s*=\s*")[^"]+(")',
                    replacement=rf"\g<1>{clean_version}\g<2>",
                    flags=re.MULTILINE,
                )
            ],
        ),
        FileUpdateConfig(
            filepath=Path("charts/pneutrinoutil/Chart.yaml"),
            rules=[
                ReplacementRule(
                    pattern=r"^(version:\s*)[0-9A-Za-z.-]+",
                    replacement=rf"\g<1>{clean_version}",
                    flags=re.MULTILINE,
                ),
                ReplacementRule(
                    pattern=r'^(appVersion:\s*")[^"]+(")',
                    replacement=rf"\g<1>{clean_version}\g<2>",
                    flags=re.MULTILINE,
                ),
            ],
        ),
        FileUpdateConfig(
            filepath=Path("server/main.go"),
            rules=[
                ReplacementRule(
                    pattern=r"(//\s*@version\s+)[0-9A-Za-z.-]+",
                    replacement=rf"\g<1>{clean_version}",
                )
            ],
        ),
    ]

    for config in configs:
        replace_in_file(config)

    Path("VERSION").write_text(f"{clean_version}\n", encoding="utf-8")
    logging.info("Updated VERSION")

    logging.info("Version files updated successfully.")


def create_bump_pr(version_info: VersionInfo, dryrun: bool = False) -> None:
    """Create a new git branch, commit version bump changes, push, and open a Pull Request via gh CLI.

    Args:
        version_info: VersionInfo containing target version strings.
        dryrun: If True, log planned execution commands without running them.

    Raises:
        BumpVersionError: If git or gh commands encounter a failure.
    """
    logging.info("Starting Pull Request creation process...")

    def execute(
        cmd: list[str],
        capture_output: bool = False,
        dryrun_stdout: str = "",
        dryrun_stderr: str = "",
        dryrun_returncode: int = 0,
    ) -> subprocess.CompletedProcess:
        if dryrun:
            logging.info(f"[DRY-RUN] Would execute: {' '.join(cmd)}")
            return subprocess.CompletedProcess(
                args=cmd,
                returncode=dryrun_returncode,
                stdout=dryrun_stdout,
                stderr=dryrun_stderr,
            )
        return run_command(cmd, capture_output=capture_output)

    sha_res = execute(
        ["git", "rev-parse", "--short", "HEAD"],
        capture_output=True,
        dryrun_stdout="dryrunsha",
    )
    if sha_res.returncode != 0:
        raise BumpVersionError("Failed to retrieve short commit SHA.")
    short_sha = sha_res.stdout.strip()

    branch_name = f"bump-version-{version_info.tag_version}-{short_sha}"

    execute(["git", "config", "user.name", "github-actions[bot]"])
    execute(
        ["git", "config", "user.email", "github-actions[bot]@users.noreply.github.com"]
    )

    execute(["git", "checkout", "-b", branch_name])
    execute(["git", "add", "-A"])

    diff_res = execute(
        ["git", "diff", "--staged", "--quiet"],
        dryrun_returncode=1,  # In dryrun, simulate having staged changes
    )
    if diff_res.returncode == 0:
        logging.info("No changes to commit. Skipping PR creation.")
        return

    commit_msg = f"chore: bump version to {version_info.tag_version}"
    execute(["git", "commit", "-m", commit_msg])

    logging.info(f"Pushing branch '{branch_name}' to remote origin...")
    execute(["git", "push", "origin", branch_name])

    title = f"chore: bump version to {version_info.tag_version}"
    body = f"""\
## Automated Version Bump

This PR updates the project version to `{version_info.tag_version}` across all required locations:

- `VERSION`
- `pyproject.toml`
- `charts/pneutrinoutil/Chart.yaml` (`appVersion` & `version`)
- `server/main.go` (`@version` annotation)
- Regenerated Swagger API specs (`server/docs/`) and UI client (`ui/app/api/client/`) via `./task gen:swag`

*Validation for SemVer and non-existence of Git tag `{version_info.tag_version}` passed successfully.*"""

    logging.info(f"Creating new PR for branch {branch_name}...")
    gh_res = execute(
        [
            "gh",
            "pr",
            "create",
            "--title",
            title,
            "--body",
            body,
            "--base",
            "main",
            "--head",
            branch_name,
        ]
    )
    if gh_res.returncode != 0:
        raise BumpVersionError("Failed to create Pull Request using gh CLI.")

    logging.info("Pull Request creation step processed successfully.")


def main() -> None:
    """Entry point for the version bump script."""
    logging.basicConfig(level=logging.INFO, format="%(levelname)s: %(message)s")
    cli_args = parse_args()
    try:
        version_info = validate_semver(cli_args.version)
        check_git_tag_availability(version_info)
        update_version_files(version_info.clean_version)
        if cli_args.create_pr:
            create_bump_pr(version_info, dryrun=cli_args.dry_run)
    except BumpVersionError as e:
        logging.error(str(e))
        sys.exit(1)


if __name__ == "__main__":
    main()
