github_repo(
    name = "pleasings2",
    repo = "sagikazarmark/mypleasings",
    revision = "ec5cf0342df6f6a015668edc353a964bf71ee28d",
)

genrule(
    name = "go.mod",
    srcs = ["go.mod"],
    outs = ["go.mod"],
    cmd = [
        """export GO=$(cat "$SRC" | grep '^go .*$') && cat > $OUT << EOF
module plz-out

$GO
EOF
      """,
    ],
    labels = [
        "link:plz-out/",
    ],
)

CONFIG.setdefault("RELEASE_TAG_PREFIX", "v")

sh_cmd(
    name = "release",
    shell = "/usr/bin/env bash",
    deps = ["///pleasings2//tools/misc:semver"],
    labels = ["manual"],
    cmd = f"\nTAG_PREFIX={CONFIG.RELEASE_TAG_PREFIX}\n" + """
set -o errexit -o pipefail

# Global regex variables
NAT="0|[1-9][0-9]*"
VERSION_REGEX="^[vV]?(\\\$NAT)\.(\\\$NAT)\.(\\\$NAT)$"

PROG=release

USAGE="\\
Usage:
  \\\$PROG (major|minor|patch)
  \\\$PROG --help

Options:
  -t, --tag              Tag the current HEAD with the new version.
  -p, --push             Push the new tag to the "origin" remote.
  -n, --dry-run          Print the commands to be executed.
  -h, --help             Print this help message.
"

function error {
  echo -e "\x1B[31m\\\$1\x1B[0m" >&2

  exit 1
}

# Default values of arguments
CREATE_TAG=false
PUSH=false
DRY_RUN=false
HELP=false
ARGS=()

# Loop through arguments and process them
for arg in "\\\$@"
do
    case \\\$arg in
        -t|--tag)
        CREATE_TAG=true
        shift # Remove --tag from processing
        ;;
        -p|--push)
        PUSH=true
        shift # Remove --push from processing
        ;;
        -n|--dry-run)
        DRY_RUN=true
        shift # Remove --dry-run from processing
        ;;
        -h|--help)
        HELP=true
        shift # Remove --help from processing
        ;;
        *)
        ARGS+=("\\\$1")
        shift # Remove generic argument from processing
        ;;
    esac
done

# Show usage
if [[ \\\$HELP == true ]]; then
    echo -e "\\\$USAGE"

    exit 0
fi

# Prepare dry-run command
CMD=""
if [[ \\\$DRY_RUN == true ]]; then
    CMD="echo"
else
    CMD=""
fi

CURRENT_BRANCH="\\\$(git branch --show-current)"
if [[ "\\\$CURRENT_BRANCH" != "master" ]]; then
    error "cannot release from branch \\"\\\$CURRENT_BRANCH\\": please switch to master"
fi

LATEST_VERSION="\\\$(git tag | grep -E "\\\$VERSION_REGEX" | sort -r --version-sort | head -1)"
if [[ "\\\$LATEST_VERSION" == "" ]]; then
    error "failed to determine latest version"
fi

VERSION_TO_BUMP="\\\${ARGS[0]}"
if [[ ! "\\\$VERSION_TO_BUMP" =~ ^(major|minor|patch)$ ]]; then
    error "invalid version bump command \\"\\\$VERSION_TO_BUMP\\""
fi

TAG="\\\$($(out_location ///pleasings2//tools/misc:semver) bump \\\$VERSION_TO_BUMP \\\$LATEST_VERSION)"

# Validate tag
if [[ "\\\$TAG" == "" ]]; then
    error "missing tag"
fi

sed -e "s/^## \[Unreleased\]$/## [Unreleased]\\\\\\"$'\\n'"\\\\\\"$'\\n'"\\\\\\"$'\\n'"## [\\\$TAG] - $(date +%Y-%m-%d)/g; s|^\[Unreleased\]: \(.*\/compare\/\)\(.*\)...HEAD$|[Unreleased]: \1\\\$TAG_PREFIX\\\$TAG...HEAD\\\\\\"$'\\n'"[\\\$TAG]: \1\2...\\\$TAG_PREFIX\\\$TAG|g" CHANGELOG.md > CHANGELOG.md.new
\\\$CMD mv CHANGELOG.md.new CHANGELOG.md

if [[ \\\$CREATE_TAG == true ]]; then
    \\\$CMD git add CHANGELOG.md
    \\\$CMD git commit -m "Prepare release \\\$TAG"
    \\\$CMD git tag -m "Release \\\$TAG" \\\$TAG_PREFIX\\\$TAG

    if [[ \\\$PUSH == true ]]; then
        \\\$CMD git push
        \\\$CMD git push origin \\\$TAG_PREFIX\\\$TAG
    fi
fi

echo "Version updated to \\\$TAG!"

if [[ \\\$PUSH == false ]]; then
	echo
	echo "Review the changes made by this script then execute the following:"

    if [[ \\\$CREATE_TAG == false ]]; then
        echo
        echo "git add CHANGELOG.md && git commit -m 'Prepare release \\\$TAG' && git tag -m 'Release \\\$TAG' \\\$TAG_PREFIX\\\$TAG"
        echo
        echo "Finally, push the changes:"
    fi

	echo
	echo "git push; git push origin \\\$TAG_PREFIX\\\$TAG"
fi
""",
)
