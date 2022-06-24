package system

import (
	"os"
	"os/user"
	"path/filepath"
)

// Egoctl System Params ...
var (
	Usr, _     = user.Current()
	EgoctlHome = filepath.Join(Usr.HomeDir, "/.egoctl")
	CurrentDir = getCurrentDirectory()
	GoPath     = os.Getenv("GOPATH")
)

func getCurrentDirectory() string {
	if dir, err := os.Getwd(); err == nil {
		return dir
	}
	return ""
}
