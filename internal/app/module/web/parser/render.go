package parser

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"github.com/flosch/pongo2"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/egoctl/internal/app/module/web/constx"
	"github.com/gotomicro/egoctl/internal/logger"
	"github.com/gotomicro/egoctl/internal/system"
	"github.com/gotomicro/egoctl/internal/utils"
	"github.com/smartwalle/pongo2render"
	"go.uber.org/zap"
)

// render
type RenderFile struct {
	*pongo2render.Render
	Context      pongo2.Context
	GenerateTime string
	Option       UserOption
	ModelName    string // 单个model name
	PackageName  string
	FlushFile    string
	PkgPath      string
	TmplPath     string
	Descriptor   Descriptor
}

func NewRender(m RenderInfo) *RenderFile {
	var (
		pathCtx       pongo2.Context
		newDescriptor Descriptor
	)

	// parse descriptor, get flush file path, beego path, etc...
	newDescriptor, pathCtx = m.Descriptor.Parse(m.Option, m.ModelName, m.ModelNames, m.Option.Path)

	obj := &RenderFile{
		Context:      make(pongo2.Context),
		Option:       m.Option,
		ModelName:    m.ModelName,
		GenerateTime: m.GenerateTime,
		Descriptor:   newDescriptor,
	}

	obj.FlushFile = newDescriptor.DstPath

	// new render
	obj.Render = pongo2render.NewRender(path.Join(obj.Option.GitLocalPath, obj.Option.ProType, m.TmplPath))

	filePath := path.Dir(obj.FlushFile)
	err := createPath(filePath)
	if err != nil {
		logger.Log.Fatalf("Could not create the controllers directory: %s", err)
	}
	// get go package path
	if m.Option.Language == constx.LanguageGo {
		obj.PkgPath = getPackagePath(m.Option.ProjectPath)
	}
	obj.SetContext("packagePath", obj.PkgPath)

	relativePath, err := filepath.Rel(system.CurrentDir, obj.FlushFile)
	if err != nil {
		logger.Log.Fatalf("Could not get the relative path: %s", err)
	}

	modelSchemas := m.Content

	importMaps := make(map[string]struct{})

	obj.PackageName = filepath.Base(filepath.Dir(relativePath))

	elog.Info("render", zap.String("modelName", obj.ModelName), zap.String("packageName", obj.PackageName))

	// package
	obj.SetContext("packageName", obj.PackageName)
	obj.SetContext("packageImports", importMaps)

	// if pathCtx["pathRelEgo"] == "." {
	//	obj.SetContext("packagePath", obj.PkgPath)
	// } else {
	//	obj.SetContext("packagePath", obj.PkgPath+"/"+pathCtx["pathRelEgo"].(string))
	// }

	obj.SetContext("packageMod", obj.PkgPath)

	obj.SetContext("modelSchemas", modelSchemas)

	for key, value := range pathCtx {
		obj.SetContext(key, value)
	}

	obj.SetContext("apiPrefix", obj.Option.ApiPrefix)
	obj.SetContext("generateTime", obj.GenerateTime)

	if obj.Option.ContextDebug {
		utils.DumpWrapper("TEMPLATE-CONTEXT-DUMP", func() { spew.Dump(obj.Context) })
	}
	return obj
}

func (r *RenderFile) SetContext(key string, value interface{}) {
	r.Context[key] = value
}

func (r *RenderFile) Exec(name string) error {
	var (
		buf string
		err error
	)
	buf, err = r.Render.Template(name).Execute(r.Context)
	if err != nil {
		return fmt.Errorf("Could not create the %s render tmpl , err: %w", name, err)
	}
	_, err = os.Stat(r.Descriptor.DstPath)
	var orgContent []byte
	if err == nil {
		if org, err := os.OpenFile(r.Descriptor.DstPath, os.O_RDONLY, 0666); err == nil {
			orgContent, _ = ioutil.ReadAll(org)
			org.Close()
		} else {
			logger.Log.Infof("file err %s", err)
		}
	}
	// Replace or create when content changes
	output := []byte(buf)
	ext := filepath.Ext(r.FlushFile)
	if r.Option.EnableFormat && ext == ".go" {
		// format code
		var bts []byte
		bts, err = format.Source([]byte(buf))
		if err != nil {
			logger.Log.Warnf("format buf error %s", err.Error())
		}
		output = bts
	}

	if FileContentChange(orgContent, output, GetSeg(ext)) {
		err = r.write(r.FlushFile, output)
		if err != nil {
			return fmt.Errorf("创建文件失败, err: %w", err)
		}
		elog.Info("create file", elog.String("packageName", r.PackageName), elog.String("flushFile", r.FlushFile))
	}
	return nil
}
