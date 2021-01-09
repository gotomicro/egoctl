package version

import (
	"io"
	"io/ioutil"
	"text/template"
	"time"

	"egoctl/logger"
)

var (
	buildVersion    string
	buildGitVersion string
	buildTag        string
	buildStatus     string
	buildUser       string
	buildHost       string
	buildTime       string
)

// InitBanner loads the banner and prints it to output
// All errors are ignored, the application will not
// print the banner in case of error.
func InitBanner(out io.Writer, in io.Reader) {
	if in == nil {
		logger.Log.Fatal("The input is nil")
	}

	banner, err := ioutil.ReadAll(in)
	if err != nil {
		logger.Log.Fatalf("Error while trying to read the banner: %s", err)
	}

	show(out, string(banner))
}

func show(out io.Writer, content string) {
	t, err := template.New("banner").
		Funcs(template.FuncMap{"Now": Now}).
		Parse(content)

	if err != nil {
		logger.Log.Fatalf("Cannot parse the banner template: %s", err)
	}

	data := map[string]interface{}{
		"Version":         buildVersion,
		"BuildGitVersion": buildGitVersion,
		"BuildTag":        buildTag,
		"BuildStatus":     buildStatus,
		"BuildUser":       buildUser,
		"BuildHost":       buildHost,
		"BuildTime":       buildTime,
	}
	err = t.Execute(out, data)
	if err != nil {
		logger.Log.Error(err.Error())
	}
}

// Now returns the current local time in the specified layout
func Now(layout string) string {
	return time.Now().Format(layout)
}
