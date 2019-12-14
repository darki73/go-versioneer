package main

import (
    "flag"
    format "fmt"
    "os"
    "strings"
)

var cvsName     string
var packageName string

func init() {
    flag.StringVar(&packageName, "cvs", "git", "Set the cvs type for the ")
    flag.StringVar(&packageName, "package", "", "Set the package which is responsible for version management")
    flag.Parse()
}

func main() {
    currentDirectory, err := os.Getwd()
    if err != nil {
        //log.Fatalf("Unable to get working directory path: %v", err)
        return
    }
    format.Println(getFlags(currentDirectory))
}

func getFlags(directory string) string {
    var repositoryInformation Handler
    switch strings.ToLower(cvsName) {
        case "svn":
            //repositoryInformation = svn { directory: directory }
            panic("SVN Handler is not implemented")
        case "mercurial":
            //repositoryInformation = mercurial { directory: directory }
            panic("Mercurial Handler is not implemented")
        case "git":
        default:
            repositoryInformation = git { directory: directory }
            break
    }
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