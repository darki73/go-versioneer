package main

import (
    "bytes"
    format "fmt"
    "os/exec"
    "strings"
    "time"
)

type git struct {
    directory       string
}

type version struct {
    Major           string
    Minor           string
    Patch           string
    Build           string
}

func (git git) exec(arguments ...string) (string, error) {
    var errorOut bytes.Buffer
    command := exec.Command("git", arguments...)
    command.Dir = git.directory
    command.Stderr = &errorOut
    commandOutput, err := command.Output()
    outputToString := strings.TrimSpace(string(commandOutput))
    if err != nil {
        err = format.Errorf("git: error=%q stderr=%s", err, string(errorOut.Bytes()))
    }
    return outputToString, err
}

func (git git) LatestCommitFullHash() string {
    result, _ := git.latestCommitFullHash()
    return result
}

func (git git) LatestCommitShortHash() string {
    result, _ := git.latestCommitShortHash()
    return result
}

func (git git) LatestCommitTag() string {
    result, _ := git.latestCommitTag()
    return result
}

func (git git) LatestCommitTagClean() string {
    commitTag := git.LatestCommitTag()
    if strings.HasPrefix(commitTag, "v") {
        commitTag = strings.Replace(commitTag, "v", "", -1)
    }
    return commitTag
}

func (git git) Version() version {
    return git.extractVersionFromTag()
}

func (git git) VersionAsString() string {
    var versionParts []string
    extractedVersion := git.extractVersionFromTag()
    if extractedVersion.Major != "" {
        versionParts = append(versionParts, extractedVersion.Major)
    }
    if extractedVersion.Minor != "" {
        versionParts = append(versionParts, extractedVersion.Minor)
    }
    if extractedVersion.Patch != "" {
        versionParts = append(versionParts, extractedVersion.Patch)
    }
    if extractedVersion.Build != "" {
        versionParts = append(versionParts, extractedVersion.Build)
    }
    return strings.Join(versionParts, ".")
}

func (git git) BuildDateLocal() string {
    return git.artificialBuildDateLocal()
}

func (git git) BuildDateUTC() string {
    return git.artificialBuildDateUTC()
}

func (git git) BuildDate(atLocal bool) string {
    if atLocal {
        return git.BuildDateLocal()
    }
    return git.BuildDateUTC()
}

func (git git) BranchName() string {
    result, _ := git.extractBranchName()
    return result
}

func (git git) LatestCommitAuthor() string {
    result, _ := git.extractCommitAuthor()
    return result
}

func (git git) RepositoryStateWithBoolean(ignoreNotTracked bool) string {
    result, _ := git.extractRepositoryState(ignoreNotTracked)
    return result
}

func (git git) IsCleanOnlyTracked() string {
    return git.RepositoryStateWithBoolean(true)
}

func (git git) IsClean() string {
    return git.RepositoryStateWithBoolean(false)
}


// [INTERNAL] Get full hash of the latest commit
func (git git) latestCommitFullHash() (string, error) {
    commitHash, err := git.exec("rev-list",  "--tags", "--max-count=1")
    if err != nil {
        return "error", err
    }
    return commitHash, nil
}

// [INTERNAL] Get short hash of the latest commit
func (git git) latestCommitShortHash() (string, error) {
    commitHash, err := git.exec("rev-parse", "--short", "HEAD")
    if err != nil {
        return "error", err
    }
    return commitHash, nil
}

// [INTERNAL] Get tag of the latest commit
func (git git) latestCommitTag() (string, error) {
    commitHash, err := git.latestCommitFullHash()
    if err != nil {
        return "error", err
    }
    commitTag, err := git.exec("describe", "--tags", commitHash)
    if err != nil {
        return "", err
    }
    return strings.TrimSpace(commitTag), nil
}

// [INTERNAL] Get artificial build date in local time
func (git git) artificialBuildDateLocal() string {
    return time.Now().Format(time.RFC3339)
}

// [INTERNAL] Get artificial build date in UTC time
func (git git) artificialBuildDateUTC() string {
    return time.Now().UTC().Format(time.RFC3339)
}

// [INTERNAL] Extract version information from tag
func (git git) extractVersionFromTag() version {
    commitTag := git.LatestCommitTagClean()
    commitParts := strings.Split(commitTag, ".")
    var version version
    switch len(commitParts) {
        case 1:
            version.Major = commitParts[0]
            break
        case 2:
            version.Major = commitParts[0]
            version.Minor = commitParts[1]
            break
        case 3:
            version.Major = commitParts[0]
            version.Minor = commitParts[1]
            version.Patch = commitParts[2]
            break
        default:
            version.Major = commitParts[0]
            version.Minor = commitParts[1]
            version.Patch = commitParts[2]
            version.Build = commitParts[3]
    }
    return version
}

// [INTERNAL] Extract current branch name
func (git git) extractBranchName() (string, error) {
    result, err := git.exec("rev-parse", "--abbrev-ref", "HEAD")
    if err != nil {
        return "undefined", err
    }
    return result, nil
}

// [INTERNAL] Extract author of the latest commit
func (git git) extractCommitAuthor() (string, error) {
    result, err := git.exec("log", "--format=%ae", git.LatestCommitFullHash(), "|", "tail", "-n 1")
    if err != nil {
        return "undefined", err
    }
    return result, nil
}

// [INTERNAL] Check if repository has any changes from HEAD
func (git git) extractRepositoryState(ignoreNotTracked bool) (string, error) {
    var result string
    var err error
    if ignoreNotTracked {
        result, err = git.exec("status", "--untracked-files=no", "--porcelain")
    } else {
        result, err = git.exec("status", "--porcelain")
    }
    if err != nil {
        return "", nil
    }
    if len(result) > 0 {
        return "dirty", nil
    }
    return "clean", nil
}