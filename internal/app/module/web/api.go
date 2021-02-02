package web

import (
	"egoctl/internal/app/module/web/core"
	"encoding/json"
	"errors"
	"github.com/gotomicro/ego/server/egin"
	"github.com/syndtr/goleveldb/leveldb"
)

type Projects struct {
	Name string `json:"name" binding:"required"`
	Path string `json:"path" binding:"required"`
}

func (c *Container) API(component *egin.Component) {
	component.POST("/api/projects", core.Handle(c.apiProjectCreate))
	component.GET("/api/projects", core.Handle(c.apiProjectList))
}

func (c *Container) apiProjectList(ctx *core.Context) {
	value, err := c.leveldb.Get([]byte("projects"), nil)
	if err != nil {
		ctx.JSONE(1, "获取projects失败: err"+err.Error(), make([]struct{}, 0))
		return
	}

	projectsList := make([]Projects, 0)
	err = json.Unmarshal(value, &projectsList)
	if err != nil {
		ctx.JSONE(1, "解析项目json失败: err"+err.Error(), make([]struct{}, 0))
		return
	}
	ctx.JSONOK(projectsList)
}

func (c *Container) apiProjectCreate(ctx *core.Context) {
	req := Projects{}
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSONE(1, "获取参数失败: err"+err.Error(), err)
		return
	}
	value, err := c.leveldb.Get([]byte("projects"), nil)

	if err != nil && !errors.Is(err, leveldb.ErrNotFound) {
		ctx.JSONE(1, "获取projects失败: err"+err.Error(), err)
		return
	}

	projectsList := make([]Projects, 0)
	if !errors.Is(err, leveldb.ErrNotFound) {
		err = json.Unmarshal(value, projectsList)
		if err != nil {
			ctx.JSONE(1, "解析项目json失败: err"+err.Error(), err)
			return
		}
	}

	isExist := false
	for _, project := range projectsList {
		if req.Path == project.Path {
			isExist = true
		}
	}

	if isExist {
		ctx.JSONE(1, "已存在该项目", nil)
		return
	}

	projectsList = append(projectsList, req)

	jsonBytes, err := json.Marshal(projectsList)
	if err != nil {
		ctx.JSONE(1, "JSON编码失败: err"+err.Error(), nil)
		return
	}

	err = c.leveldb.Put([]byte("projects"), jsonBytes, nil)
	if err != nil {
		ctx.JSONE(1, "存入leveldb失败, err: "+err.Error(), nil)
		return
	}

	ctx.JSONOK(projectsList)
}
