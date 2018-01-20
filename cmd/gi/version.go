package main

import (
	"fmt"
)

var LastGitCommitHash string
var BuildTimeStamp string
var NearestGitTag string
var GitBranch string
var GoVersion string
var LuajitVersion string = "LuaJIT_2.1.0-beta3+jea-int64-additions"

func Version() string {
	return fmt.Sprintf("built: '%s'\nlast-git-commit-hash: '%s'\nnearest-git-tag: '%s'\ngit-branch: '%s'\ngo-version: '%s'\nluajit-version: '%s'", BuildTimeStamp, LastGitCommitHash, NearestGitTag, GitBranch, GoVersion, LuajitVersion)
}
