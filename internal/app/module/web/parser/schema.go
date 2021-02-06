package parser

import (
	"fmt"
	"github.com/flosch/pongo2"
	"github.com/gotomicro/egoctl/internal/pkg/command"
	"github.com/gotomicro/egoctl/internal/pkg/system"
	"github.com/gotomicro/egoctl/logger"
	"github.com/gotomicro/egoctl/utils"
	"github.com/smartwalle/pongo2render"
	"path/filepath"
	"strings"
	"sync"
)

// store all data
type Container struct {
	UserOption       UserOption             // user option
	TmplOption       TmplOption             // tmpl option
	CurPath          string                 // user current path
	EnableModules    map[string]interface{} // beego pro provider a collection of module
	FunctionOnce     map[string]sync.Once   // exec function once
	GenerateTime     string
	GenerateTimeUnix int64
	Timestamp        Timestamp
	parser           *astParser
	err              error
}

// user option
type UserOption struct {
	Debug              bool     `json:"debug"`
	ContextDebug       bool     `json:"contextDebug"`
	ScaffoldDSLContent string   // ego pro dsl
	Language           string   `json:"language"`
	ProType            string   `json:"proType"`
	ApiPrefix          string   `json:"apiPrefix"`
	EnableModule       []string `json:"enableModule"`
	ProjectPath        string
	GitLocalPath       string            `json:"gitLocalPath"`
	EnableFormat       bool              `json:"enableFormat"`
	Path               map[string]string `json:"path"`
}

// tmpl option
type TmplOption struct {
	RenderPath string `toml:"renderPath"`
	Descriptor []Descriptor
}

type Descriptor struct {
	Module  string `toml:"module"`
	SrcName string `toml:"srcName"`
	DstPath string `toml:"dstPath"`
	Once    bool   `toml:"once"`
	Script  string `toml:"script"`
}

func (descriptor Descriptor) Parse(option UserOption, modelName string, modelNames []string, paths map[string]string) (newDescriptor Descriptor, ctx pongo2.Context) {
	var (
		err             error
		relativeDstPath string
		absFile         string
		relPath         string
	)

	newDescriptor = descriptor
	render := pongo2render.NewRender("")
	ctx = make(pongo2.Context)
	for key, value := range paths {
		absFile, err = filepath.Abs(value)
		if err != nil {
			logger.Log.Fatalf("absolute path error %s from key %s and value %s", err, key, value)
		}
		relPath, err = filepath.Rel(system.CurrentDir, absFile)
		if err != nil {
			logger.Log.Fatalf("Could not get the relative path: %s", err)
		}
		// user input path
		ctx["path"+utils.CamelCase(key)] = option.ProjectPath + "/" + value
		// relativePath
		ctx["pathRel"+utils.CamelCase(key)] = relPath
	}
	ctx["modelName"] = lowerFirst(utils.CamelString(modelName))
	ctx["modelNames"] = modelNames
	ctx["modelNameSnake"] = utils.SnakeString(modelName)
	relativeDstPath, err = render.TemplateFromString(descriptor.DstPath).Execute(ctx)
	if err != nil {
		logger.Log.Fatalf("egoctl tmpl exec error, err: %s", err)
		return
	}
	newDescriptor.DstPath, err = filepath.Abs(relativeDstPath)
	if err != nil {
		logger.Log.Fatalf("absolute path error %s from flush file %s", err, relativeDstPath)
	}

	newDescriptor.Script, err = render.TemplateFromString(descriptor.Script).Execute(ctx)
	if err != nil {
		logger.Log.Fatalf("parse script %s, error %s", descriptor.Script, err)
	}
	return
}

func (descriptor Descriptor) IsExistScript() bool {
	return descriptor.Script != ""
}

func (d Descriptor) ExecScript(path string) (err error) {
	arr := strings.Split(d.Script, " ")
	if len(arr) == 0 {
		return
	}

	stdout, stderr, err := command.ExecCmdDir(path, arr[0], arr[1:]...)
	if err != nil {
		return concatenateError(err, stderr)
	}

	logger.Log.Info(stdout)
	return nil
}

type Timestamp struct {
	GitCacheLastRefresh int64 `toml:"gitCacheLastRefresh"`
	Generate            int64 `toml:"generate"`
}

func concatenateError(err error, stderr string) error {
	if len(stderr) == 0 {
		return err
	}
	return fmt.Errorf("%v: %s", err, stderr)
}
