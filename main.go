package main

import (
    "flag"
    format "fmt"
    "log"
    "os"
    "strings"
)

var packageName string

func init() {
    flag.StringVar(&packageName, "package", "", "Set the package which is responsible for version management")
    flag.Parse()
}

func main() {

    log.Println("Going to target: ", packageName)


    currentDirectory, err := os.Getwd()
    if err != nil {
        log.Fatalf("Unable to get working directory path: %v", err)
    }
    log.Println(getFlags(currentDirectory))
}

func getFlags(directory string) string {
    repositoryInformation := git{ directory: directory }
    if packageName == "" {
        packageName = "main"
    }

    return strings.Join([]string {
        format.Sprintf("-X %s.BuildDateLocal=%s",   packageName, repositoryInformation.BuildDateLocal()),
        format.Sprintf("-X %s.BuildDateUTC=%s",     packageName, repositoryInformation.BuildDateUTC()),
        format.Sprintf("-X %s.GitCommitLong=%s",    packageName, repositoryInformation.LatestCommitFullHash()),
        format.Sprintf("-X %s.GitCommitShort=%s",   packageName, repositoryInformation.LatestCommitShortHash()),
        format.Sprintf("-X %s.GitBranch=%s",        packageName, repositoryInformation.BranchName()),
        format.Sprintf("-X %s.GitState=%s",         packageName, repositoryInformation.IsClean()),
        format.Sprintf("-X %s.GitAuthor=%s",        packageName, repositoryInformation.LatestCommitAuthor()),
        format.Sprintf("-X %s.GitVersion=%s",       packageName, repositoryInformation.VersionAsString()),
        format.Sprintf("-X %s.GitSummary=%s",       packageName, format.Sprintf("%s-%s", repositoryInformation.VersionAsString(), repositoryInformation.IsClean())),
    }, " ")
}