// Copyright 2013 bee authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"text/template"
	"time"
	"unicode"

	"egoctl/config"
	"egoctl/internal/pkg/system"
	"egoctl/logger"
	"egoctl/logger/colors"
)

type tagName struct {
	Name string `json:"name"`
}

func GetBeeWorkPath() string {
	curpath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return curpath
}

// Go is a basic promise implementation: it wraps calls a function in a goroutine
// and returns a channel which will later return the function's return value.
func Go(f func() error) chan error {
	ch := make(chan error)
	go func() {
		ch <- f()
	}()
	return ch
}

// GetGOPATHs returns all paths in GOPATH variable.
func GetGOPATHs() []string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" && strings.Compare(runtime.Version(), "go1.8") >= 0 {
		gopath = defaultGOPATH()
	}
	return filepath.SplitList(gopath)
}

// CloseFile attempts to close the passed file
// or panics with the actual error
func CloseFile(f *os.File) {
	err := f.Close()
	MustCheck(err)
}

// MustCheck panics when the error is not nil
func MustCheck(err error) {
	if err != nil {
		panic(err)
	}
}

// WriteToFile creates a file and writes content to it
func WriteToFile(filename, content string) {
	f, err := os.Create(filename)
	MustCheck(err)
	defer CloseFile(f)
	_, err = f.WriteString(content)
	MustCheck(err)
}

// __FILE__ returns the file name in which the function was invoked
func FILE() string {
	_, file, _, _ := runtime.Caller(1)
	return file
}

// __LINE__ returns the line number at which the function was invoked
func LINE() int {
	_, _, line, _ := runtime.Caller(1)
	return line
}

// BeeFuncMap returns a FuncMap of functions used in different templates.
func BeeFuncMap() template.FuncMap {
	return template.FuncMap{
		"trim":       strings.TrimSpace,
		"bold":       colors.Bold,
		"headline":   colors.MagentaBold,
		"foldername": colors.RedBold,
		"endline":    EndLine,
		"tmpltostr":  TmplToString,
	}
}

// TmplToString parses a text template and return the result as a string.
func TmplToString(tmpl string, data interface{}) string {
	t := template.New("tmpl").Funcs(BeeFuncMap())
	template.Must(t.Parse(tmpl))

	var doc bytes.Buffer
	err := t.Execute(&doc, data)
	MustCheck(err)

	return doc.String()
}

// EndLine returns the a newline escape character
func EndLine() string {
	return "\n"
}

func Tmpl(text string, data interface{}) {
	output := colors.NewColorWriter(os.Stderr)

	t := template.New("Usage").Funcs(BeeFuncMap())
	template.Must(t.Parse(text))

	err := t.Execute(output, data)
	if err != nil {
		logger.Log.Error(err.Error())
	}
}

func CheckEnv(appname string) (apppath, packpath string, err error) {
	gps := GetGOPATHs()
	if len(gps) == 0 {
		logger.Log.Error("if you want new a go module project,please add param `-gopath=false`.")
		logger.Log.Fatal("GOPATH environment variable is not set or empty")
	}
	currpath, _ := os.Getwd()
	currpath = filepath.Join(currpath, appname)
	for _, gpath := range gps {
		gsrcpath := filepath.Join(gpath, "src")
		if strings.HasPrefix(strings.ToLower(currpath), strings.ToLower(gsrcpath)) {
			packpath = strings.Replace(currpath[len(gsrcpath)+1:], string(filepath.Separator), "/", -1)
			return currpath, packpath, nil
		}
	}

	// In case of multiple paths in the GOPATH, by default
	// we use the first path
	gopath := gps[0]

	logger.Log.Warn("You current workdir is not inside $GOPATH/src.")
	logger.Log.Debugf("GOPATH: %s", FILE(), LINE(), gopath)

	gosrcpath := filepath.Join(gopath, "src")
	apppath = filepath.Join(gosrcpath, appname)

	if _, e := os.Stat(apppath); !os.IsNotExist(e) {
		err = fmt.Errorf("cannot create application without removing '%s' first", apppath)
		logger.Log.Errorf("Path '%s' already exists", apppath)
		return
	}
	packpath = strings.Join(strings.Split(apppath[len(gosrcpath)+1:], string(filepath.Separator)), "/")
	return
}

func PrintErrorAndExit(message, errorTemplate string) {
	Tmpl(fmt.Sprintf(errorTemplate, message), nil)
	os.Exit(2)
}

// GoCommand executes the passed command using Go tool
func GoCommand(command string, args ...string) error {
	allargs := []string{command}
	allargs = append(allargs, args...)
	goBuild := exec.Command("go", allargs...)
	goBuild.Stderr = os.Stderr
	return goBuild.Run()
}

// SplitQuotedFields is like strings.Fields but ignores spaces
// inside areas surrounded by single quotes.
// To specify a single quote use backslash to escape it: '\''
func SplitQuotedFields(in string) []string {
	type stateEnum int
	const (
		inSpace stateEnum = iota
		inField
		inQuote
		inQuoteEscaped
	)
	state := inSpace
	r := []string{}
	var buf bytes.Buffer

	for _, ch := range in {
		switch state {
		case inSpace:
			if ch == '\'' {
				state = inQuote
			} else if !unicode.IsSpace(ch) {
				buf.WriteRune(ch)
				state = inField
			}

		case inField:
			if ch == '\'' {
				state = inQuote
			} else if unicode.IsSpace(ch) {
				r = append(r, buf.String())
				buf.Reset()
			} else {
				buf.WriteRune(ch)
			}

		case inQuote:
			if ch == '\'' {
				state = inField
			} else if ch == '\\' {
				state = inQuoteEscaped
			} else {
				buf.WriteRune(ch)
			}

		case inQuoteEscaped:
			buf.WriteRune(ch)
			state = inQuote
		}
	}

	if buf.Len() != 0 {
		r = append(r, buf.String())
	}

	return r
}

// GetFileModTime returns unix timestamp of `os.File.ModTime` for the given path.
func GetFileModTime(path string) int64 {
	path = strings.Replace(path, "\\", "/", -1)
	f, err := os.Open(path)
	if err != nil {
		logger.Log.Errorf("Failed to open file on '%s': %s", path, err)
		return time.Now().Unix()
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		logger.Log.Errorf("Failed to get file stats: %s", err)
		return time.Now().Unix()
	}

	return fi.ModTime().Unix()
}

func defaultGOPATH() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	if home := os.Getenv(env); home != "" {
		return filepath.Join(home, "go")
	}
	return ""
}

func GetGoVersionSkipMinor() string {
	strArray := strings.Split(runtime.Version()[2:], `.`)
	return strArray[0] + `.` + strArray[1]
}

func IsGOMODULE() bool {
	if combinedOutput, e := exec.Command(`go`, `env`).CombinedOutput(); e != nil {
		logger.Log.Errorf("i cann't find go.")
	} else {
		regex := regexp.MustCompile(`GOMOD="?(.+go.mod)"?`)
		stringSubmatch := regex.FindStringSubmatch(string(combinedOutput))
		return len(stringSubmatch) == 2
	}
	return false
}

func NoticeUpdateEgoctl() {
	cmd := exec.Command("go", "version")
	cmd.Output()
	if cmd.Process == nil || cmd.Process.Pid <= 0 {
		logger.Log.Warn("There is no go environment")
		return
	}
	egoctlHome := system.EgoctlHome
	if !IsExist(egoctlHome) {
		err := os.Mkdir(egoctlHome, os.ModePerm)
		if err != nil {
			logger.Log.Warnf("Create egoctlHome file err: %s", err)
			return
		}
	}

	fp := egoctlHome + "/.noticeUpdateEgoctl"
	timeNow := time.Now().Unix()
	var timeOld int64
	if !IsExist(fp) {
		f, err := os.Create(fp)
		if err != nil {
			logger.Log.Warnf("Create noticeUpdateEgoctl file err: %s", err)
			return
		}
		defer f.Close()
	}
	oldContent, err := ioutil.ReadFile(fp)
	if err != nil {
		logger.Log.Warnf("Read noticeUpdateEgoctl file err: %s", err)
		return
	}
	timeOld, _ = strconv.ParseInt(string(oldContent), 10, 64)
	if timeNow-timeOld < 24*60*60 {
		return
	}
	w, err := os.OpenFile(fp, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		logger.Log.Warnf("Open noticeUpdateEgoctl file err: %s", err)
		return
	}
	defer w.Close()
	timeNowStr := strconv.FormatInt(timeNow, 10)
	if _, err := w.WriteString(timeNowStr); err != nil {
		logger.Log.Warnf("Update noticeUpdateEgoctl file err: %s", err)
		return
	}
	logger.Log.Info("Getting egoctl latest version...")
	versionLast := EgoctlLastVersion()
	versionNow := config.Version
	if versionLast == "" {
		logger.Log.Warn("Get latest version err")
		return
	}
	if versionNow != versionLast {
		logger.Log.Warnf("Update available %s ==> %s", versionNow, versionLast)
		logger.Log.Warn("Run `egoctl update` to update")
	}
	logger.Log.Info("Your egoctl are up to date")
}

func EgoctlLastVersion() (version string) {
	var url = ""
	resp, err := http.Get(url)
	if err != nil {
		logger.Log.Warnf("Get egoctl tags from github error: %s", err)
		return
	}
	defer resp.Body.Close()
	bodyContent, _ := ioutil.ReadAll(resp.Body)
	var tags []tagName
	if err = json.Unmarshal(bodyContent, &tags); err != nil {
		logger.Log.Warnf("Unmarshal tags body error: %s", err)
		return
	}
	if len(tags) < 1 {
		logger.Log.Warn("There is no tags！")
		return
	}
	last := tags[0]
	re, _ := regexp.Compile(`[0-9.]+`)
	versionList := re.FindStringSubmatch(last.Name)
	if len(versionList) > 0 {
		return versionList[0]
	}
	logger.Log.Warn("There is no tags！")
	return
}

func DumpWrapper(msg string, dump func()) {
	fmt.Println()
	logger.Log.Infof("------------------------%s--------------------------", msg)
	dump()
	logger.Log.Infof("------------------------%s--------------------------\n\n", msg)
}
