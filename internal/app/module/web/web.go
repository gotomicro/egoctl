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
	"io"
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

type webui struct {
	webuiembed embed.FS
	path       string
}

func (w *webui) Open(name string) (http.File, error) {
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return nil, errors.New("http: invalid character in file path")
	}
	fullName := filepath.Join(w.path, filepath.FromSlash(path.Clean("/"+name)))
	file, err := w.webuiembed.Open(fullName)
	wf := &WebuiFile{
		File: file,
	}
	return wf, err
}

type WebuiFile struct {
	io.Seeker
	fs.File
}

func (*WebuiFile) Readdir(count int) ([]fs.FileInfo, error) {
	return nil, nil
}

type webuiIndex struct {
	webuiembed embed.FS
	path       string
}

func (w *webuiIndex) Open(name string) (http.File, error) {
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return nil, errors.New("http: invalid character in file path")
	}
	fullName := filepath.Join(w.path, filepath.FromSlash(path.Clean("/index.html")))
	file, err := w.webuiembed.Open(fullName)
	wf := &WebuiFile{
		File: file,
	}
	return wf, err
}

func (c *Container) Run() {
	var err error
	c.leveldb, err = leveldb.OpenFile(c.DataPath, nil)
	if err != nil {
		elog.Panic("level db open file error", elog.FieldErr(err))
	}
	defer c.leveldb.Close()
	project.InitProjectSrv(c.leveldb)
	template.InitTemplateSrv(c.leveldb)

	webuiObj := &webui{
		webuiembed: webui2.WebUI,
		path:       "dist",
	}

	webuiIndexObj := &webuiIndex{
		webuiembed: webui2.WebUI,
		path:       "dist",
	}
	econf.LoadFromReader(strings.NewReader(config), toml.Unmarshal)
	if err := ego.New().Serve(func() *egin.Component {
		server := egin.Load("server.http").Build()
		server.GET("/", func(ctx *gin.Context) {
			ctx.Redirect(302, "/projects")
			return
		})

		server.GET("/projects", func(context *gin.Context) {
			context.FileFromFS("/projects", webuiIndexObj)
		})
		server.GET("/templates", func(context *gin.Context) {
			context.FileFromFS("/templates", webuiIndexObj)
		})

		server.StaticFS("/webui/", webuiObj)
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
