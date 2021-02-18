package parser

import (
	"fmt"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/egoctl/internal/app/module/web/constx"
	"github.com/gotomicro/egoctl/internal/pkg/system"
	"github.com/gotomicro/egoctl/internal/pkg/utils"
	"github.com/pelletier/go-toml"
	"sync"
	"time"
)

const (
	MDateFormat = "20060102_150405"
)

func NewParser(option UserOption) *Container {
	obj := &Container{
		GenerateTime:     time.Now().Format(MDateFormat),
		GenerateTimeUnix: time.Now().Unix(),
		TmplOption:       TmplOption{},
		CurPath:          system.CurrentDir,
		EnableModules:    make(map[string]interface{}), // get the user configuration, get the enable module result
		FunctionOnce:     make(map[string]sync.Once),   // get the tmpl configuration, get the function once result
		StoreData: StoreData{
			UserOption: option,
		},
	}
	obj.UserOption = option
	return obj
}

func (c *Container) Run() error {
	c.initUserOption()
	c.initTemplateOption()
	c.initParser()
	c.initRender()
	return c.err
}

// 初始化用户配置
func (c *Container) initUserOption() {
	// 如果是Go语言，那么就需要判断是否有go.mod，因为需要go.mod里的数据
	if c.UserOption.Language == constx.LanguageGo && !utils.IsExist(c.UserOption.ProjectPath+"/go.mod") {
		c.err = fmt.Errorf("请在%s目录下创建go.mod文件", c.UserOption.ProjectPath)
		return
	}
	c.EnableModules["*"] = struct{}{}
	c.StoreData.EnableModules = c.EnableModules
}

// 解析模板配置
func (c *Container) initTemplateOption() {
	if c.err != nil {
		return
	}
	tree, err := toml.LoadFile(c.UserOption.GitLocalPath + "/" + c.UserOption.ProType + "/egoctl.toml")
	if err != nil {
		c.err = fmt.Errorf("egoctl tmpl exec error, err: %w", err)
		return
	}

	err = tree.Unmarshal(&c.TmplOption)
	if err != nil {
		c.err = fmt.Errorf("egoctl tmpl parse error, err: %w", err)
		return
	}

	c.StoreData.TemplateOption = c.TmplOption

	for _, value := range c.TmplOption.Descriptor {
		if value.Once {
			c.FunctionOnce[value.SrcName] = sync.Once{}
		}
	}
}

func (c *Container) initParser() {
	if c.err != nil {
		return
	}
	c.parser = AstParserBuild(c.UserOption, c.TmplOption)
}

func (c *Container) initRender() {
	if c.err != nil {
		return
	}
	for _, desc := range c.TmplOption.Descriptor {
		_, allFlag := c.EnableModules["*"]
		_, moduleFlag := c.EnableModules[desc.Module]
		if !allFlag && !moduleFlag {
			continue
		}

		models := c.parser.GetRenderInfos(desc)
		c.StoreData.ModelData = models
		// model table name, model table schema
		for _, m := range models {
			// some render exec once
			syncOnce, flag := c.FunctionOnce[desc.SrcName]
			if flag {
				syncOnce.Do(func() {
					c.err = c.renderModel(m)
					if c.err != nil {
						return
					}
				})
				continue
			}
			c.err = c.renderModel(m)
			if c.err != nil {
				return
			}
		}
	}
}

func (c *Container) renderModel(m RenderInfo) error {
	// todo optimize
	m.GenerateTime = c.GenerateTime
	render := NewRender(m)
	// 如果只给json数据
	if c.UserOption.Mode == "json" {
		return nil
	}

	err := render.Exec(m.Descriptor.SrcName)
	if err != nil {
		return err
	}
	if render.Descriptor.IsExistScript() {
		err := render.Descriptor.ExecScript(c.CurPath)
		if err != nil {
			elog.Errorf("egoctl exec shell error, err: %s", err)
		}
	}
	return nil
}

func (c *Container) GetRenderData() StoreData {
	return c.StoreData
}
