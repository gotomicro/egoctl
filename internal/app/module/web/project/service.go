package project

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gotomicro/egoctl/internal/app/module/web/constx"
	"github.com/gotomicro/egoctl/internal/app/module/web/parser"
	"github.com/gotomicro/egoctl/internal/app/module/web/template"
	"github.com/syndtr/goleveldb/leveldb"
)

type Info struct {
	Name          string   `json:"name" binding:"required"`
	Path          string   `json:"path" binding:"required"`
	GitRemotePath string   `json:"gitRemotePath" binding:"required"`
	ProType       string   `json:"proType"`      // 默认类型
	Language      string   `json:"language"`     // Go React Vue 其他
	ApiPrefix     string   `json:"apiPrefix"`    // API 前缀
	DSL           string   `json:"dsl"`          // dsl 描述
	EnableModule  []string `json:"enableModule"` // 开启模块
	Ctime         int64    `json:"ctime"`
	Utime         int64    `json:"utime"`
}

type InfoDSL struct {
	Path string `json:"path" binding:"required"`
	DSL  string `json:"dsl"` // dsl 描述
}

// 用户看到的列表数据
type InfoDto struct {
	Name          string `json:"name" binding:"required"`          // 名称
	GitRemotePath string `json:"gitRemotePath" binding:"required"` // 远程地址
	Path          string `json:"path"`                             // 存储路径
	TemplateName  string `json:"templateName"`                     // 模板名称
	ProType       string `json:"proType"`                          // 默认类型
	Language      string `json:"language"`                         // Go React Vue 其他
	ApiPrefix     string `json:"apiPrefix"`                        // API 前缀
	DSL           string `json:"dsl"`                              // dsl 描述
	Ctime         int64  `json:"ctime"`
	Utime         int64  `json:"utime"`
}

type InfoUniqId struct {
	Path string `json:"path" form:"path"`
}

type Infos []Info

func (i Infos) ToInfoDtos() []InfoDto {
	output := make([]InfoDto, 0)
	for _, value := range i {
		tmplInfo, _ := template.Srv.TemplateInfo(template.InfoUniqId{GitRemotePath: template.GitURL(value.GitRemotePath)})
		output = append(output, InfoDto{
			Name:          value.Name,
			GitRemotePath: value.GitRemotePath,
			Path:          value.Path,
			TemplateName:  tmplInfo.Name,
			Ctime:         value.Ctime,
			Utime:         value.Utime,
			ProType:       value.ProType,
			ApiPrefix:     value.ApiPrefix,
			DSL:           value.DSL,
			Language:      value.Language,
		})
	}
	return output
}

var Srv *projectSrv

type projectSrv struct {
	l       sync.RWMutex
	leveldb *leveldb.DB
}

func InitProjectSrv(leveldb *leveldb.DB) {
	Srv = &projectSrv{
		leveldb: leveldb,
	}
}

func (p *projectSrv) ProjectList() (list []InfoDto, err error) {
	value, err := p.leveldb.Get([]byte("projects"), nil)
	if err != nil {
		err = fmt.Errorf("获取projects失败: %w", err)
		return
	}

	projectsList := make(Infos, 0)
	err = json.Unmarshal(value, &projectsList)
	if err != nil {
		err = fmt.Errorf("解析项目json失败: %w", err)
		return
	}

	list = projectsList.ToInfoDtos()
	return
}

func (p *projectSrv) ProjectCreate(req Info) (err error) {
	// 防止并发请求
	p.l.Lock()
	defer p.l.Unlock()
	value, err := p.leveldb.Get([]byte(constx.LevelDBProjects), nil)

	if err != nil && !errors.Is(err, leveldb.ErrNotFound) {
		err = fmt.Errorf("获取projects失败: %w", err)
		return
	}

	projectsList := make([]Info, 0)
	if !errors.Is(err, leveldb.ErrNotFound) {
		err = json.Unmarshal(value, &projectsList)
		if err != nil {
			err = fmt.Errorf("解析项目json失败: %w", err)
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
		err = fmt.Errorf("已存在该项目")
		return
	}

	req.Ctime = time.Now().Unix()
	req.Utime = time.Now().Unix()
	projectsList = append(projectsList, req)

	jsonBytes, err := json.Marshal(projectsList)
	if err != nil {
		err = fmt.Errorf("JSON编码失败: %w", err)
		return
	}

	err = p.leveldb.Put([]byte(constx.LevelDBProjects), jsonBytes, nil)
	if err != nil {
		err = fmt.Errorf("写入leveldb失败: %w", err)
		return
	}
	return
}

func (p *projectSrv) ProjectUpdate(req Info) (err error) {
	// 防止并发请求
	p.l.Lock()
	defer p.l.Unlock()
	value, err := p.leveldb.Get([]byte(constx.LevelDBProjects), nil)

	if err != nil && !errors.Is(err, leveldb.ErrNotFound) {
		err = fmt.Errorf("获取projects失败: %w", err)
		return
	}

	projectsList := make([]Info, 0)
	if !errors.Is(err, leveldb.ErrNotFound) {
		err = json.Unmarshal(value, &projectsList)
		if err != nil {
			err = fmt.Errorf("解析项目json失败: %w", err)
			return
		}
	}

	isExist := false
	for _, project := range projectsList {
		if req.Path == project.Path {
			isExist = true
		}
	}

	if !isExist {
		err = fmt.Errorf("不存在该项目数据")
		return
	}

	listNew := make([]Info, 0)
	for _, value := range projectsList {
		if value.Path == req.Path {
			value.Name = req.Name
			value.GitRemotePath = req.GitRemotePath
			value.Utime = time.Now().Unix()
			value.ApiPrefix = req.ApiPrefix
			value.ProType = req.ProType
			value.Language = req.Language
		}
		listNew = append(listNew, value)
	}

	jsonBytes, err := json.Marshal(listNew)
	if err != nil {
		err = fmt.Errorf("JSON编码失败: %w", err)
		return
	}

	err = p.leveldb.Put([]byte(constx.LevelDBProjects), jsonBytes, nil)
	if err != nil {
		err = fmt.Errorf("写入leveldb失败: %w", err)
		return
	}
	return
}

func (p *projectSrv) ProjectDSL(req InfoDSL) (err error) {
	// 防止并发请求
	p.l.Lock()
	defer p.l.Unlock()
	value, err := p.leveldb.Get([]byte(constx.LevelDBProjects), nil)

	if err != nil && !errors.Is(err, leveldb.ErrNotFound) {
		err = fmt.Errorf("获取projects失败: %w", err)
		return
	}

	projectsList := make([]Info, 0)
	if !errors.Is(err, leveldb.ErrNotFound) {
		err = json.Unmarshal(value, &projectsList)
		if err != nil {
			err = fmt.Errorf("解析项目json失败: %w", err)
			return
		}
	}

	isExist := false
	for _, project := range projectsList {
		if req.Path == project.Path {
			isExist = true
		}
	}

	if !isExist {
		err = fmt.Errorf("不存在该项目数据")
		return
	}

	listNew := make([]Info, 0)
	for _, value := range projectsList {
		if value.Path == req.Path {
			value.DSL = req.DSL
			value.Utime = time.Now().Unix()
		}
		listNew = append(listNew, value)
	}

	jsonBytes, err := json.Marshal(listNew)
	if err != nil {
		err = fmt.Errorf("JSON编码失败: %w", err)
		return
	}

	err = p.leveldb.Put([]byte(constx.LevelDBProjects), jsonBytes, nil)
	if err != nil {
		err = fmt.Errorf("写入leveldb失败: %w", err)
		return
	}
	return
}

func (t *projectSrv) ProjectInfo(info InfoUniqId) (resp Info, err error) {
	// 防止并发请求
	t.l.RLock()
	defer t.l.RUnlock()

	value, err := t.leveldb.Get([]byte(constx.LevelDBProjects), nil)
	if err != nil {
		err = fmt.Errorf("获取LevelDB模板列表数据失败, err: %w", err)
		return
	}

	list := make([]Info, 0)
	// 如果已经存在数据，那么进行解析
	err = json.Unmarshal(value, &list)
	if err != nil {
		err = fmt.Errorf("解析LevelDB模板列表数据失败, err: %w", err)
		return
	}

	isExist := false
	for _, value := range list {
		// 该模板地址已存在，不允许插入
		if value.Path == info.Path {
			isExist = true
			resp = value
		}
	}

	if !isExist {
		err = fmt.Errorf("不存在该git模板数据")
		return
	}

	return
}

func (p *projectSrv) ProjectGen(req InfoUniqId) (err error) {
	// 防止并发请求
	info, err := p.ProjectInfo(req)
	if err != nil {
		return fmt.Errorf("获取projects失败: %w", err)
	}

	templateInfo, err := template.Srv.TemplateInfo(template.InfoUniqId{GitRemotePath: template.GitURL(info.GitRemotePath)})
	if err != nil {
		return fmt.Errorf("获取模板信息失败: %w", err)
	}
	parserObj := parser.NewParser(parser.UserOption{
		Language:           info.Language,
		ScaffoldDSLContent: info.DSL,
		ProType:            info.ProType,
		ApiPrefix:          info.ApiPrefix,
		EnableModule:       make([]string, 0),
		ProjectPath:        info.Path,
		GitLocalPath:       templateInfo.Path,
		EnableFormat:       false,
		Path: map[string]string{
			"backend": ".",
		},
	})

	err = parserObj.Run()
	if err != nil {
		return fmt.Errorf("生成代码失败: %w", err)
	}
	return
}

func (p *projectSrv) ProjectRender(req InfoUniqId) (resp parser.StoreData, err error) {
	// 防止并发请求
	info, err := p.ProjectInfo(req)
	if err != nil {
		return resp, fmt.Errorf("获取projects失败: %w", err)
	}

	templateInfo, err := template.Srv.TemplateInfo(template.InfoUniqId{GitRemotePath: template.GitURL(info.GitRemotePath)})
	if err != nil {
		return resp, fmt.Errorf("获取模板信息失败: %w", err)
	}
	parserObj := parser.NewParser(parser.UserOption{
		Mode:               "json",
		Language:           info.Language,
		ScaffoldDSLContent: info.DSL,
		ProType:            info.ProType,
		ApiPrefix:          info.ApiPrefix,
		EnableModule:       make([]string, 0),
		ProjectPath:        info.Path,
		GitLocalPath:       templateInfo.Path,
		EnableFormat:       false,
		Path: map[string]string{
			"backend": ".",
		},
	})

	err = parserObj.Run()
	if err != nil {
		return resp, fmt.Errorf("生成代码失败: %w", err)
	}

	return parserObj.GetRenderData(), nil
}

func (t *projectSrv) ProjectDelete(info InfoUniqId) (err error) {
	// 防止并发请求
	t.l.Lock()
	defer t.l.Unlock()
	value, err := t.leveldb.Get([]byte(constx.LevelDBProjects), nil)
	if err != nil {
		err = fmt.Errorf("获取LevelDB项目列表数据失败, err: %w", err)
		return
	}

	list := make([]Info, 0)
	// 如果已经存在数据，那么进行解析
	err = json.Unmarshal(value, &list)
	if err != nil {
		return fmt.Errorf("解析LevelDB项目列表数据失败, err: %w", err)
	}

	isExist := false
	for _, value := range list {
		// 该模板地址已存在，不允许插入
		if value.Path == info.Path {
			isExist = true
		}
	}

	if !isExist {
		err = fmt.Errorf("不存在该项目数据")
		return
	}

	listNew := make([]Info, 0)

	for _, value := range list {
		if value.Path == info.Path {
			continue
		}
		listNew = append(listNew, value)
	}

	jsonBytes, err := json.Marshal(listNew)
	if err != nil {
		return fmt.Errorf("编码LevelDB项目列表数据失败, err: %w", err)
	}
	err = t.leveldb.Put([]byte(constx.LevelDBProjects), jsonBytes, nil)
	if err != nil {
		return fmt.Errorf("存入LevelDB项目列表数据失败, err: %w", err)
	}
	return
}
