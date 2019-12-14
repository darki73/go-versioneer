package main

type Handler interface {
    exec(arguments ...string) (string, error)
    LatestCommitFullHash() string
    LatestCommitShortHash() string
    LatestCommitTag() string
    LatestCommitTagClean() string
    Version() version
    VersionAsString() string
    BuildDateLocal() string
    BuildDateUTC() string
    BuildDate(atLocal bool) string
    BranchName() string
    LatestCommitAuthor() string
    RepositoryStateWithBoolean(ignoreNotTracked bool) string
    IsCleanOnlyTracked() string
    IsClean() string
    latestCommitFullHash() (string, error)
    latestCommitShortHash() (string, error)
    latestCommitTag() (string, error)
    artificialBuildDateLocal() string
    artificialBuildDateUTC() string
    extractVersionFromTag() version
    extractBranchName() (string, error)
    extractCommitAuthor() (string, error)
    extractRepositoryState(ignoreNotTracked bool) (string, error)
}
