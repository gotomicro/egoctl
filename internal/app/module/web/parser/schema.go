package parser

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gotomicro/egoctl/internal/app/module/web/parser/pongo2"
	"github.com/gotomicro/egoctl/internal/app/module/web/parser/pongo2render"
	"github.com/gotomicro/egoctl/internal/command"
	"github.com/gotomicro/egoctl/internal/logger"
	"github.com/gotomicro/egoctl/internal/system"
	"github.com/gotomicro/egoctl/internal/utils"
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
	StoreData        StoreData
}

// user option
type UserOption struct {
	Mode               string            `json:"mode"` // mode: tmpl 模板，json json数据
	ContextDebug       bool              `json:"contextDebug"`
	ScaffoldDSLContent string            `json:"scaffoldDslContent"`
	Language           string            `json:"language"`
	ProType            string            `json:"proType"`
	ApiPrefix          string            `json:"apiPrefix"`
	EnableModule       []string          `json:"enableModule"`
	ProjectPath        string            `json:"projectPath"`
	GitLocalPath       string            `json:"gitLocalPath"`
	EnableFormat       bool              `json:"enableFormat"`
	Path               map[string]string `json:"path"`
}

type StoreData struct {
	EnableModules  map[string]interface{} `json:"enableModules"` // 开启的模块
	UserOption     UserOption             `json:"userOption"`
	TemplateOption TmplOption             `json:"templateOption"`
	ModelData      []RenderInfo           `json:"modelData"`
}

// tmpl option
type TmplOption struct {
	RenderPath string       `toml:"renderPath" json:"renderPath"`
	Descriptor []Descriptor `json:"descriptor"`
}

type Descriptor struct {
	Module  string `toml:"module" json:"module"`
	SrcName string `toml:"srcName" json:"srcName"`
	DstPath string `toml:"dstPath" json:"dstPath"`
	Once    bool   `toml:"once" json:"once"`
	Script  string `toml:"script" json:"script"`
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
