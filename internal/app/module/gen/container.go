package gen

import (
	"egoctl/internal/pkg/git"
	"egoctl/internal/pkg/system"
	"egoctl/internal/pkg/utils"
	"egoctl/logger"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/pelletier/go-toml"
	"github.com/spf13/viper"
	"io/ioutil"
	"path/filepath"
	"sync"
	"time"
)

const (
	defaultBinName = "egoctl"
	MDateFormat    = "20060102_150405"
)

var DefaultEgoctlPro = &Container{
	ScaffoldBinName:  defaultBinName,
	ScaffoldConfFile: filepath.Join(system.CurrentDir, fmt.Sprintf("%s.toml", defaultBinName)),

	TimestampFile: filepath.Join(system.CurrentDir, fmt.Sprintf(".%s.timestamp", defaultBinName)),
	GoModFile:     system.CurrentDir + "/go.mod",
	UserOption: UserOption{
		Debug:           false,
		ScaffoldDSLFile: filepath.Join(system.CurrentDir, fmt.Sprintf("%s.go", defaultBinName)),
		ContextDebug:    false,
		ProType:         "default",
		ApiPrefix:       "/api",
		EnableModule:    nil,
		GitRemotePath:   "",
		Branch:          "master",
		GitLocalPath:    system.EgoctlHome + "/egoctl",
		EnableFormat:    true,
		EnableGitPull:   true,
		Path: map[string]string{
			"ego": ".",
		},
		EnableGomod:    true,
		RefreshGitTime: 24 * 3600,
	},
	GenerateTime:     time.Now().Format(MDateFormat),
	GenerateTimeUnix: time.Now().Unix(),
	TmplOption:       TmplOption{},
	CurPath:          system.CurrentDir,
	EnableModules:    make(map[string]interface{}), // get the user configuration, get the enable module result
	FunctionOnce:     make(map[string]sync.Once),   // get the tmpl configuration, get the function once result
}

func (c *Container) Run() {
	// 项目第一次使用egoctl时间
	c.initTimestamp()
	c.initUserOption()
	c.initTemplateOption()
	c.initParser()
	c.initRender()
	c.flushTimestamp()
}

// initTimestamp 项目第一次使用egoctl时间
func (c *Container) initTimestamp() {
	// 如果存在该时间，读取该文件
	if utils.IsExist(c.TimestampFile) {
		tomlFile, err := toml.LoadFile(c.TimestampFile)
		if err != nil {
			logger.Log.Fatalf("egoctl timestamp tmpl load error, err: %s", err)
			return
		}
		err = tomlFile.Unmarshal(&c.Timestamp)
		if err != nil {
			logger.Log.Fatalf("egoctl timestamp tmpl parse error, err: %s", err)
			return
		}
	}
	// 将代码生成时间存入时间戳文件里
	c.Timestamp.Generate = c.GenerateTimeUnix
}

// 初始化用户配置
func (c *Container) initUserOption() {
	if !utils.IsExist(c.ScaffoldConfFile) {
		logger.Log.Fatalf("egoctl config is not exist, path: %s", c.ScaffoldConfFile)
		return
	}
	viper.SetConfigFile(c.ScaffoldConfFile)
	err := viper.ReadInConfig()
	if err != nil {
		logger.Log.Fatalf("read egoctl config content, err: %s", err.Error())
		return
	}

	err = viper.Unmarshal(&c.UserOption)
	if err != nil {
		logger.Log.Fatalf("egoctl config unmarshal error, err: %s", err.Error())
		return
	}
	if c.UserOption.Debug {
		utils.DumpWrapper("VIPER-DUMP", func() { viper.Debug() })
	}

	if c.UserOption.EnableGomod {
		if !utils.IsExist(c.GoModFile) {
			logger.Log.Fatalf("go mod not exist, please create go mod file")
			return
		}
	}

	for _, value := range c.UserOption.EnableModule {
		c.EnableModules[value] = struct{}{}
	}

	if len(c.EnableModules) == 0 {
		c.EnableModules["*"] = struct{}{}
	}

	if c.UserOption.Debug {
		logger.Log.Infof("c.modules: %+v", c.EnableModules)
	}
}

func (c *Container) initTemplateOption() {
	if c.UserOption.EnableGitPull && (c.GenerateTimeUnix-c.Timestamp.GitCacheLastRefresh > c.UserOption.RefreshGitTime) {
		err := git.CloneORPullRepo(c.UserOption.GitRemotePath, c.UserOption.GitLocalPath)
		if err != nil {
			logger.Log.Fatalf("egoctl pro git clone or pull repo error, err: %s", err)
			return
		}
		c.Timestamp.GitCacheLastRefresh = c.GenerateTimeUnix
	}

	tree, err := toml.LoadFile(c.UserOption.GitLocalPath + "/" + c.UserOption.ProType + "/egoctl.toml")

	if err != nil {
		logger.Log.Fatalf("egoctl tmpl exec error, err: %s", err)
		return
	}
	err = tree.Unmarshal(&c.TmplOption)
	if err != nil {
		logger.Log.Fatalf("egoctl tmpl parse error, err: %s", err)
		return
	}

	if c.UserOption.Debug {
		utils.DumpWrapper("TEMPLATE-DUMP", func() { spew.Dump(c.TmplOption) })
	}

	for _, value := range c.TmplOption.Descriptor {
		if value.Once {
			c.FunctionOnce[value.SrcName] = sync.Once{}
		}
	}
}

func (c *Container) initParser() {
	c.parser = AstParserBuild(c.UserOption, c.TmplOption)
}

func (c *Container) initRender() {
	for _, desc := range c.TmplOption.Descriptor {
		_, allFlag := c.EnableModules["*"]
		_, moduleFlag := c.EnableModules[desc.Module]
		if !allFlag && !moduleFlag {
			continue
		}

		models := c.parser.GetRenderInfos(desc)
		// model table name, model table schema
		for _, m := range models {
			// some render exec once
			syncOnce, flag := c.FunctionOnce[desc.SrcName]
			if flag {
				syncOnce.Do(func() {
					c.renderModel(m)
				})
				continue
			}
			c.renderModel(m)
		}
	}
}

func (c *Container) renderModel(m RenderInfo) {
	// todo optimize
	m.GenerateTime = c.GenerateTime
	render := NewRender(m)
	render.Exec(m.Descriptor.SrcName)
	if render.Descriptor.IsExistScript() {
		err := render.Descriptor.ExecScript(c.CurPath)
		if err != nil {
			logger.Log.Fatalf("egoctl exec shell error, err: %s", err)
		}
	}
}

func (c *Container) flushTimestamp() {
	tomlByte, err := toml.Marshal(c.Timestamp)
	if err != nil {
		logger.Log.Fatalf("marshal timestamp tmpl parse error, err: %s", err)
	}
	err = ioutil.WriteFile(c.TimestampFile, tomlByte, 0644)
	if err != nil {
		logger.Log.Fatalf("flush timestamp tmpl parse error, err: %s", err)
	}
}
