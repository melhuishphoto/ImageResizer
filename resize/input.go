package resize

import (
	"runtime"

	"github.com/skratchdot/open-golang/open"
	"github.com/sqweek/dialog"
)

func ChooseFile() string {
	file, e := dialog.File().Load()
	if e != nil {
		return ""
	}
	return file
}

func ChooseDir() string {
	dir, e := dialog.Directory().Browse()
	if e != nil {
		return ""
	}
	return dir
}

func OpenDir(dir string) error {
	return open.RunWith(dir, fileViewer())
}

func fileViewer() string {
	switch runtime.GOOS {
	case "windows":
		return "Explorer"
	case "darwin":
		return "Finder"
	default:
		return ""
	}
}
