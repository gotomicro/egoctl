package web

import (
	"egoctl/internal/app/module/web/core"
	"egoctl/internal/app/module/web/project"
	"egoctl/internal/app/module/web/template"
	"github.com/gotomicro/ego/server/egin"
	"github.com/gotomicro/gotoant"
)

func (c *Container) API(component *egin.Component) {
	component.GET("/api/projects", core.Handle(c.apiProjectList))
	component.GET("/api/projects/gen", core.Handle(c.apiProjectGen)) // 生成代码
	component.POST("/api/projects", core.Handle(c.apiProjectCreate))
	component.PUT("/api/projects", core.Handle(c.apiProjectUpdate))
	component.PUT("/api/projects/dsl", core.Handle(c.apiProjectDSL))
	component.DELETE("/api/projects", core.Handle(c.apiProjectDelete))
	component.GET("/api/templates", core.Handle(c.apiTemplateList))
	component.GET("/api/templates/select", core.Handle(c.apiTemplateSelect))
	component.POST("/api/templates", core.Handle(c.apiTemplateCreate))
	component.PUT("/api/templates", core.Handle(c.apiTemplateUpdate))
	component.PUT("/api/templates/sync", core.Handle(c.apiTemplateSync)) // 同步模板代码
	component.DELETE("/api/templates", core.Handle(c.apiTemplateDelete))
}

func (c *Container) apiProjectList(ctx *core.Context) {
	list, err := project.Srv.ProjectList()
	if err != nil {
		ctx.JSONE(1, "获取信息失败: err"+err.Error(), make([]struct{}, 0))
		return
	}
	ctx.JSONOK(list)
}

func (c *Container) apiProjectGen(ctx *core.Context) {
	req := project.InfoUniqId{}
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSONE(1, "获取参数失败: err"+err.Error(), err)
		return
	}
	err = project.Srv.ProjectGen(req)
	if err != nil {
		ctx.JSONE(1, "生成代码失败: err"+err.Error(), make([]struct{}, 0))
		return
	}
	ctx.JSONOK()
}

func (c *Container) apiProjectCreate(ctx *core.Context) {
	req := project.Info{}
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSONE(1, "获取参数失败: err"+err.Error(), err)
		return
	}
	err = project.Srv.ProjectCreate(req)
	if err != nil {
		ctx.JSONE(1, "创建项目失败: err"+err.Error(), err)
		return
	}
	ctx.JSONOK()
}

func (c *Container) apiProjectUpdate(ctx *core.Context) {
	req := project.Info{}
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSONE(1, "获取参数失败: err"+err.Error(), err)
		return
	}
	err = project.Srv.ProjectUpdate(req)
	if err != nil {
		ctx.JSONE(1, "更新项目失败: err"+err.Error(), err)
		return
	}
	ctx.JSONOK()
}

func (c *Container) apiProjectDSL(ctx *core.Context) {
	req := project.InfoDSL{}
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSONE(1, "获取参数失败: err"+err.Error(), err)
		return
	}
	err = project.Srv.ProjectDSL(req)
	if err != nil {
		ctx.JSONE(1, "更新项目失败: err"+err.Error(), err)
		return
	}
	ctx.JSONOK()
}

func (c *Container) apiProjectDelete(ctx *core.Context) {
	req := project.InfoUniqId{}
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSONE(1, "获取参数失败: err"+err.Error(), err)
		return
	}
	err = project.Srv.ProjectDelete(req)
	if err != nil {
		ctx.JSONE(1, "删除项目失败: err"+err.Error(), err)
		return
	}
	ctx.JSONOK()
}

func (c *Container) apiTemplateList(ctx *core.Context) {
	list, err := template.Srv.TemplateList()
	if err != nil {
		ctx.JSONE(1, "获取模板列表失败: err"+err.Error(), make([]struct{}, 0))
		return
	}
	ctx.JSONOK(list)
}

func (c *Container) apiTemplateSelect(ctx *core.Context) {
	list, err := template.Srv.TemplateList()
	if err != nil {
		ctx.JSONE(1, "获取模板列表失败: err"+err.Error(), make([]struct{}, 0))
		return
	}

	antselect := gotoant.NewSelect()
	for _, value := range list {
		antselect.SetOption(value.Name, value.GitRemotePath)
	}
	ctx.JSONOK(antselect.GetOptions())
}

// 创建模板
func (c *Container) apiTemplateCreate(ctx *core.Context) {
	req := template.Info{}
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSONE(1, "获取参数失败: err"+err.Error(), err)
		return
	}
	err = template.Srv.TemplateCreate(req)
	if err != nil {
		ctx.JSONE(1, "创建模板失败: err"+err.Error(), nil)
		return
	}
	ctx.JSONOK()
}

func (c *Container) apiTemplateUpdate(ctx *core.Context) {
	req := template.Info{}
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSONE(1, "获取参数失败: err"+err.Error(), err)
		return
	}
	err = template.Srv.TemplateUpdate(req)
	if err != nil {
		ctx.JSONE(1, "更新模板失败: err"+err.Error(), nil)
		return
	}
	ctx.JSONOK()
}

func (c *Container) apiTemplateSync(ctx *core.Context) {
	req := template.InfoUniqId{}
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSONE(1, "获取参数失败: err"+err.Error(), err)
		return
	}
	err = template.Srv.TemplateSync(req)
	if err != nil {
		ctx.JSONE(1, "同步模板失败: err"+err.Error(), nil)
		return
	}
	ctx.JSONOK()
}

func (c *Container) apiTemplateDelete(ctx *core.Context) {
	req := template.InfoUniqId{}
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSONE(1, "获取参数失败: err"+err.Error(), err)
		return
	}
	err = template.Srv.TemplateDelete(template.Info{
		GitRemotePath: req.GitRemotePath,
	})
	if err != nil {
		ctx.JSONE(1, "删除模板失败: err"+err.Error(), nil)
		return
	}
	ctx.JSONOK()
}
