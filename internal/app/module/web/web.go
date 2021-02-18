package web

import (
	"embed"
	"errors"
	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server/egin"
	"github.com/gotomicro/egoctl/internal/app/module/web/project"
	"github.com/gotomicro/egoctl/internal/app/module/web/template"
	"github.com/gotomicro/egoctl/internal/pkg/system"
	webui2 "github.com/gotomicro/egoctl/webui"
	"github.com/syndtr/goleveldb/leveldb"
	"io/fs"
	"net/http"
	"path"
	"path/filepath"
	"strings"
)

type Container struct {
	leveldb  *leveldb.DB
	DataPath string
}

var DefaultWebContainer = &Container{
	leveldb:  nil,
	DataPath: system.EgoctlHome + "/egoctl/data",
}

var config = `[logger.default]
debug=true
enableAsync=false
[trace.jaeger]
[logger.ego]
debug=true
enableAsync=false
[server.http]
host="0.0.0.0"
port=9999`

func (c *Container) Run() {
	var err error
	c.leveldb, err = leveldb.OpenFile(c.DataPath, nil)
	if err != nil {
		elog.Panic("level db open file error", elog.FieldErr(err), elog.FieldName(c.DataPath))
	}
	defer c.leveldb.Close()
	project.InitProjectSrv(c.leveldb)
	template.InitTemplateSrv(c.leveldb)

	webuiObj := &webui{
		webuiEmbed: webui2.WebUI,
		path:       "dist",
	}

	// 设置Ant Design前端访问，try file到index.html
	webuiAntIndexObj := &webuiIndex{
		webui: webuiObj,
	}
	econf.LoadFromReader(strings.NewReader(config), toml.Unmarshal)
	if err := ego.New().Serve(func() *egin.Component {
		server := egin.Load("server.http").Build()
		server.GET("/", func(ctx *gin.Context) {
			ctx.Redirect(302, "/projects")
			return
		})

		server.GET("/projects", func(context *gin.Context) {
			context.FileFromFS("/projects", http.FS(webuiAntIndexObj))
		})
		server.GET("/templates", func(context *gin.Context) {
			context.FileFromFS("/templates", http.FS(webuiAntIndexObj))
		})

		server.StaticFS("/webui/", http.FS(webuiObj))
		c.API(server)
		server.GET("/hello", func(ctx *gin.Context) {
			ctx.JSON(200, "Hello EGO")
			return
		})
		return server
	}()).Run(); err != nil {
		elog.Panic("startup", elog.FieldErr(err))
	}
}

// 嵌入普通的静态资源
type webui struct {
	webuiEmbed embed.FS // 静态资源
	path       string   // 设置embed文件到静态资源的相对路径，也就是embed注释里的路径
}

// 静态资源被访问的核心逻辑
func (w *webui) Open(name string) (fs.File, error) {
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return nil, errors.New("http: invalid character in file path")
	}
	fullName := filepath.Join(w.path, filepath.FromSlash(path.Clean("/"+name)))
	file, err := w.webuiEmbed.Open(fullName)
	return file, err
}

// Ant Design前端页面，需要该方式，实现刷新，访问到前端index.html
type webuiIndex struct {
	webui *webui
}

func (w *webuiIndex) Open(name string) (fs.File, error) {
	return w.webui.Open("index.html")
}
