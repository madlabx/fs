package buildcontext

import "github.com/madlabx/pkgx/utils"

type BuildInfo struct {
	Module    string
	Version   string
	Branch    string
	Commit    string
	BuildDate string
}

var (
	Version   = "(untracked)"
	Commit    = "(unknown)"
	BuildDate = "(unknown)"
	Module    = "(unknown)"
	Branch    = "(unknown)"
)

func Get() *BuildInfo {
	return &BuildInfo{
		Module:    Module,
		Version:   Version,
		Commit:    Commit,
		Branch:    Branch,
		BuildDate: BuildDate,
	}
}

func (bi *BuildInfo) String() string {
	return utils.ToString(bi)
}
