# Releasing

This guide describes how maintainers release v4 versions from the `4.x` branch: assess the version, prepare a release PR, push module tags, and publish a GitHub Release.

Run commands from the repository root unless a step says otherwise. Use the same shell throughout so that the variables set below remain available. You need Git, a Go toolchain supported by the target branch, and the GitHub CLI (`gh`) authenticated with permission to push and create releases. Make builds the release tools from `internal/tools` as needed.

## 1. Update the main branch and assess the version

**Before running `make gorelease`, update your local `4.x` branch to the latest remote commit.**

### Update the branch

Check the working tree and resolve any outstanding changes before continuing:

```sh
git status --short --branch
```

Fetch branches and tags, then update `4.x`:

```sh
git fetch --prune --tags origin
git switch 4.x
git merge --ff-only origin/4.x
```

### Assess release compatibility

```sh
make gorelease
```

This runs `gorelease` separately for the core and OpenTelemetry modules. It checks release-related issues and compares API compatibility when a suitable release baseline is available.

Use the report together with the changes since the previous release to choose the version:

| Change | Version consideration |
| --- | --- |
| Backward-compatible bug fixes | Patch release, such as `v4.0.1` |
| Backward-compatible features | Minor release, such as `v4.1.0` |
| Incompatible public API changes | Reassess the major version and migration plan |
| Go requirements, dependencies, or build environment changes | Assess against the project's compatibility policy and explain upgrade requirements |

The report is one input to the decision. Also review runtime behavior, dependency changes, and the minimum Go version. For the first v4 release, there may be no previous v4 tag to compare against; assess the migration itself and inspect the reported issues.

> The current Makefile ignores nonzero exit statuses from `gorelease`. Read each module's output: a successful `make` exit status alone does not mean the checks passed.

### Set release variables

The following examples use **`v4.1.0`**. Replace it with the version you have chosen:

```sh
BASE_BRANCH=4.x
VERSION=v4.1.0
MODSET=stable
RELEASE_BRANCH="release/$VERSION"
PRERELEASE_BRANCH="prerelease_${MODSET}_${VERSION}"
```

Check for existing tags:

```sh
git tag --list "*$VERSION"
```

A new release normally uses tags that do not exist yet. If a tag exists, determine whether it belongs to an earlier, incomplete attempt at this release before proceeding.

## 2. Understand the release model

`versions.yaml` defines module sets and their versions. The `multimod` tool synchronizes version references and creates module tags.

The `stable` set contains:

| Module | Module path | Tag format |
| --- | --- | --- |
| Core | `github.com/flc1125/go-gitlab-webhook/v4` | `v4.X.Y` |
| OpenTelemetry middleware | `github.com/flc1125/go-gitlab-webhook/middleware/otel/v4` | `middleware/otel/v4.X.Y` |

Both modules use the same version and are tagged at the same release commit. The internal tools module is listed under `excluded-modules` and is not released with the `stable` set.

The release sequence is:

```text
Update 4.x
  → Run make gorelease and choose the version
  → Create a release branch
  → Update and commit versions.yaml
  → Run prerelease and merge its generated commit
  → Validate and open a release PR
  → Merge the PR into 4.x
  → Update 4.x
  → Tag the merge commit and push both module tags
  → Publish a GitHub Release
  → Verify the release
```

### Makefile commands

| Command | Purpose | Parameters |
| --- | --- | --- |
| `make gorelease` | Assess module release issues and API compatibility | None required |
| `make verify-mods` | Verify module sets, versions, and inter-module dependency constraints | None required |
| `make prerelease` | Generate and commit version updates | `MODSET=stable` |
| `make add-tags` | Create local tags for a module set | `MODSET=stable`; optional `COMMIT` |
| `make push-tags` | Push local tags matching a version | `TAG=v4.X.Y` |

## 3. Create a release branch

Create the branch from the updated `4.x` checkout:

```sh
git switch -c "$RELEASE_BRANCH"
```

Collect the version configuration and generated updates on this branch, then merge them into `4.x` through a PR.

> If the target version's configuration, both `Version()` functions, and the internal module dependency have already been updated and merged, as with the initial v4.0.0 migration preparation, continue with validation. There is no need to create duplicate version commits or an empty release PR. Set `PR_NUMBER` to the PR that merged those updates, skip PR creation, and use its merge commit in step 7.

## 4. Update and commit the version configuration

Change the `stable` version in `versions.yaml`:

```yaml
module-sets:
  stable:
    version: v4.1.0
    modules:
      - github.com/flc1125/go-gitlab-webhook/v4
      - github.com/flc1125/go-gitlab-webhook/middleware/otel/v4

excluded-modules:
  - github.com/flc1125/go-gitlab-webhook/internal/tools/v4
```

Only change the version in this step; retain the existing module paths and lists.

Review and commit it:

```sh
git diff -- versions.yaml
git add versions.yaml
git commit -m "chore(release): prepare $VERSION"
```

**Commit before running `prerelease`**, which requires a clean working tree. The next step updates the modules' `Version()` functions and internal dependency version automatically.

## 5. Generate and merge version updates

### Run release preparation

```sh
make prerelease MODSET="$MODSET"
```

This first runs `verify-mods`. Then `multimod`:

1. Checks the working tree and target version tags.
2. Updates each module's `version.go`.
3. Updates references to modules in the set in `go.mod` files.
4. Runs dependency cleanup.
5. Commits the generated changes on `prerelease_stable_v4.X.Y`.
6. Returns to the branch from which the command was run.

Here, **`prerelease` means release preparation**. It does not add an RC suffix or create a GitHub prerelease.

### Review the generated changes

After the command finishes, you should be back on `release/v4.X.Y`:

```sh
git status --short --branch
git log --oneline HEAD.."$PRERELEASE_BRANCH"
git diff HEAD.."$PRERELEASE_BRANCH"
```

Expected updates include:

- The core module's `Version()` returns the target version without the `v` prefix.
- The OpenTelemetry module's `Version()` returns the same version.
- The OpenTelemetry module requires the target version of the core module.
- Any necessary `go.mod` and `go.sum` changes from dependency cleanup.

Merge the generated commit into the release branch:

```sh
git merge --ff-only "$PRERELEASE_BRANCH"
```

## 6. Validate and open a release PR

### Validate locally

```sh
make verify-mods
make ci
make test-race
```

| Command | Checks |
| --- | --- |
| `make verify-mods` | Module sets, versions, and inter-module dependency constraints |
| `make ci` | Dependency cleanup, lint, builds, and uncommitted changes to tracked files |
| `make test-race` | Core and middleware tests with race detection |

If dependency cleanup produces changes that should be retained, review and commit them, then confirm validation passes.

Review the complete PR diff and working tree:

```sh
git diff --check "origin/$BASE_BRANCH"...HEAD
git diff --stat "origin/$BASE_BRANCH"...HEAD
git status --short --branch
```

### Push and open the PR

```sh
git push --set-upstream origin "$RELEASE_BRANCH"
```

The PR description should include:

- The target version and version files updated by this PR.
- The main changes since the previous release tag, distinguished from the version-only PR diff.
- Compatibility changes and Go toolchain requirements that consumers need to act on.
- Validation performed and the two planned module tags.

Save the PR description to a temporary file outside the repository. Replace the placeholder path below with that file:

```sh
gh pr create \
  --base "$BASE_BRANCH" \
  --head "$RELEASE_BRANCH" \
  --title "Release $VERSION" \
  --body-file /path/to/pr-body.md
```

Record the PR number and watch its checks:

```sh
PR_NUMBER=$(gh pr view "$RELEASE_BRANCH" --json number --jq '.number')
gh pr checks "$PR_NUMBER" --watch
```

Once checks and review are complete, merge the PR into `4.x`.

## 7. Update the merged branch and create tags

### Identify the release commit

After the PR is merged:

```sh
git fetch --prune --tags origin
git switch "$BASE_BRANCH"
git merge --ff-only "origin/$BASE_BRANCH"
```

Check the PR state and base branch:

```sh
gh pr view "$PR_NUMBER" --json state,baseRefName,mergeCommit
```

Confirm that it is `MERGED` into `4.x`, then capture and inspect the merge commit:

```sh
RELEASE_COMMIT=$(gh pr view "$PR_NUMBER" --json mergeCommit --jq '.mergeCommit.oid')

git show --stat "$RELEASE_COMMIT"
git show "$RELEASE_COMMIT:versions.yaml"
```

The module set and version in your current `versions.yaml` must match the release commit: `add-tags` reads tag information from the current configuration file.

### Create local tags

```sh
make add-tags MODSET="$MODSET" COMMIT="$RELEASE_COMMIT"
```

For v4.1.0, this creates:

```text
v4.1.0
middleware/otel/v4.1.0
```

**Tag the release PR's merge commit.** With squash merging, this differs from the original commits on the release branch.

`COMMIT` defaults to `HEAD`. Pass the merge commit explicitly to avoid tagging a later commit if `4.x` has advanced.

### Verify tag targets

```sh
git rev-parse "$VERSION^{commit}"
git rev-parse "middleware/otel/$VERSION^{commit}"
```

Both results must equal `$RELEASE_COMMIT`. `make add-tags` only creates local tags; pushing is a separate step.

## 8. Push module tags

Inspect the tags to be pushed:

```sh
git tag --list "*$VERSION"
```

Then push them:

```sh
make push-tags TAG="$VERSION"
```

This command pushes local tags ending in the specified version, covering both the core and OpenTelemetry tags.

Verify the remote tags:

```sh
git ls-remote --tags origin \
  "refs/tags/$VERSION" \
  "refs/tags/$VERSION^{}" \
  "refs/tags/middleware/otel/$VERSION" \
  "refs/tags/middleware/otel/$VERSION^{}"
```

For annotated tags, the entry without `^{}` identifies the tag object; the entry with `^{}` identifies the commit it points to. Both tags must resolve to `$RELEASE_COMMIT`.

## 9. Publish a GitHub Release

Create one GitHub Release for the core tag, following the project's existing convention. Include both module tags and upgrade commands in its notes.

### Write release notes

Include:

- Highlights and the minimum Go version.
- Compatibility changes and migration steps.
- Major dependency updates.
- Upgrade commands for the core and OpenTelemetry modules.
- A full changelog link comparing the previous and current releases.

Example upgrade commands:

```sh
go get github.com/flc1125/go-gitlab-webhook/v4@v4.1.0
go get github.com/flc1125/go-gitlab-webhook/middleware/otel/v4@v4.1.0
```

Example changelog link:

```text
https://github.com/flc1125/go-gitlab-webhook/compare/v4.0.0...v4.1.0
```

For the first v4.0.0 release, include migration guidance and choose a comparison baseline appropriate to the major-version transition.

### Publish a stable release

Save the notes outside the repository and replace the placeholder path:

```sh
gh release create "$VERSION" \
  --verify-tag \
  --target "$BASE_BRANCH" \
  --title "$VERSION" \
  --notes-file /path/to/release-notes.md \
  --latest
```

`--verify-tag` requires the remote tag to exist, preventing release creation from implicitly creating a tag.

### Choose Latest and prerelease flags

| Release intent | Flags |
| --- | --- |
| Stable release that should become the recommended version | `--latest` |
| Stable release that should not replace the current Latest | `--latest=false` |
| RC, beta, or other prerelease | `--prerelease --latest=false` |

For an RC such as `v4.1.0-rc.1`, use that full version from the start of release preparation. Publish it with:

```sh
gh release create "$VERSION" \
  --verify-tag \
  --target "$BASE_BRANCH" \
  --title "$VERSION" \
  --notes-file /path/to/release-notes.md \
  --prerelease \
  --latest=false
```

## 10. Verify the published release

### Check GitHub

```sh
gh release view "$VERSION" \
  --json url,tagName,isDraft,isPrerelease,publishedAt
```

Confirm the tag, publication status, and prerelease flag are correct. If this release should be Latest, also check:

```sh
gh api repos/flc1125/go-gitlab-webhook/releases/latest --jq .tag_name
```

The result should be the version just released.

### Check external module resolution and builds

Create a temporary Go project outside the repository:

```sh
CHECK_DIR=$(mktemp -d)

(
  set -e
  cd "$CHECK_DIR" || exit 1
  export GOWORK=off

  go mod init example.com/release-check

  go get github.com/flc1125/go-gitlab-webhook/v4@"$VERSION"
  go get github.com/flc1125/go-gitlab-webhook/middleware/otel/v4@"$VERSION"

  go build \
    github.com/flc1125/go-gitlab-webhook/v4 \
    github.com/flc1125/go-gitlab-webhook/middleware/otel/v4
)
```

This resolves the published versions without local `replace` directives and verifies that an external consumer can fetch and build both modules. If a module proxy cannot resolve the version immediately after tags are pushed, first verify the remote tags and their targets, then retry later.

## 11. Troubleshooting

### Git reports a clean tree, but prerelease rejects it

The tool may report:

```text
VerifyWorkingTreeClean failed: working tree not clean
```

During a release, the Git library used by `multimod` did not honor local ignore rules in `.git/info/exclude` and treated those files as untracked.

To resolve this:

1. Identify the local files and ignore rules involved.
2. Temporarily move the affected files outside the repository.
3. Retry `make prerelease MODSET=stable`.
4. Restore the files to their original locations whether the command succeeds or fails.

These local files do not belong in the release commit.

### Prerelease succeeds, but the current branch has no version updates

The tool commits updates on `prerelease_<module set>_<version>`, such as `prerelease_stable_v4.1.0`, then returns to the original branch. Review and merge the generated branch as described in step 5.

### Make ci reports an unclean working tree

`make ci` includes dependency cleanup. If that updates `go.mod` or `go.sum`, the final working-tree check fails. Review the changes, commit those that should be retained, and validate again. Include the required updates in the release PR.

### Make gorelease exits successfully despite errors in its output

The current Makefile ignores failures from each module's `gorelease` invocation. Inspect the output to distinguish API compatibility reports, unavailable release baselines, module resolution or build failures, and other release issues. Do not rely on the exit status alone.

### Only one module tag was pushed

`make push-tags` pushes tags separately, not atomically. Inspect both remote tags, confirm that any successful tag points to the correct commit, and push the missing tag:

```sh
git push origin "refs/tags/$VERSION"
```

Or, for the middleware tag:

```sh
git push origin "refs/tags/middleware/otel/$VERSION"
```

Confirm both tags before publishing the GitHub Release.

---

A release is complete when its version updates are merged into `4.x`, both module tags point to the intended merge commit and are pushed, the GitHub Release is published, and an external project can fetch and build the released modules.
