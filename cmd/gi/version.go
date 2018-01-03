package main

import (
	"fmt"
)

var LastGitCommitHash string
var BuildTimeStamp string
var NearestGitTag string
var GitBranch string
var GoVersion string

func Version() string {
	return fmt.Sprintf("built: '%s'\nlast-git-commit-hash: '%s'\nnearest-git-tag: '%s'\ngit-branch: '%s'\ngo-version: '%s'\n", BuildTimeStamp, LastGitCommitHash, NearestGitTag, GitBranch, GoVersion)
}
