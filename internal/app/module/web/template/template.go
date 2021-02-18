package template

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gotomicro/egoctl/internal/app/module/web/constx"
	"github.com/gotomicro/egoctl/internal/pkg/git"
	"github.com/gotomicro/egoctl/internal/pkg/system"
	"github.com/gotomicro/egoctl/internal/pkg/utils"
	"github.com/syndtr/goleveldb/leveldb"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

// 模板服务
var Srv *templateSrv

type templateSrv struct {
	leveldb *leveldb.DB
	l       sync.RWMutex
}

type Info struct {
	Name          string `json:"name" binding:"required"`          // 名称
	GitRemotePath GitURL `json:"gitRemotePath" binding:"required"` // 远程地址
	Path          string `json:"path"`                             // 存储路径
}

// 用户看到的列表数据
type InfoDto struct {
	Name          string `json:"name" binding:"required"`          // 名称
	GitRemotePath GitURL `json:"gitRemotePath" binding:"required"` // 远程地址
	Path          string `json:"path"`                             // 存储路径
	StatusText    string `json:"statusText"`
}

type Infos []Info

func (i Infos) ToInfoDtos() []InfoDto {
	output := make([]InfoDto, 0)
	for _, value := range i {
		output = append(output, InfoDto{
			Name:          value.Name,
			GitRemotePath: value.GitRemotePath,
			Path:          value.Path,
			StatusText:    value.StatusText(),
		})
	}
	return output
}

type InfoUniqId struct {
	GitRemotePath GitURL `json:"gitRemotePath"  binding:"required"` // 远程地址
}

type GitURL string

func (u GitURL) Parse() (TmplURL, error) {
	// 如果不是http和https协议，按git协议解析
	if !strings.HasPrefix(string(u), "http") && !strings.HasPrefix(string(u), "https") {
		return parseGit(string(u))
	}

	urlInfo, err := url.Parse(string(u))
	if err != nil {
		return TmplURL{}, err
	}

	if !strings.HasSuffix(urlInfo.Path, ".git") {
		return TmplURL{}, fmt.Errorf("git url不规范，没有.git后缀")
	}

	return TmplURL{
		Path: "/" + urlInfo.Host + strings.TrimSuffix(urlInfo.Path, ".git"),
	}, nil
}

// git@github.com:gotomicro/egoctl-tmpls.git
func parseGit(url string) (TmplURL, error) {
	reg := regexp.MustCompile(`(\w*)@([\w\.]*):([\w\/-]+)\.git`)
	arr := reg.FindStringSubmatch(url)
	if len(arr) != 4 {
		return TmplURL{}, fmt.Errorf("长度不正确")
	}

	return TmplURL{
		// host: github.com + path: gotomicro/egoctl-tmpls
		Path: "/" + arr[2] + "/" + arr[3],
	}, nil
}

type TmplURL struct {
	Path string // 存储路径
}

var DefaultTemplateInfo = Info{
	Name:          "EGO官方模板",
	GitRemotePath: "https://github.com/gotomicro/egoctl-tmpls.git",
	Path:          system.EgoctlHome + "/egoctl/git/gotomicro/egoctl-tmpls",
}

func InitTemplateSrv(leveldb *leveldb.DB) {
	Srv = &templateSrv{
		leveldb: leveldb,
	}
}

func (t *templateSrv) TemplateList() ([]InfoDto, error) {
	// 防止并发请求
	t.l.Lock()
	defer t.l.Unlock()
	value, err := t.leveldb.Get([]byte(constx.LevelDBTemplates), nil)
	if err != nil && !errors.Is(err, leveldb.ErrNotFound) {
		return nil, fmt.Errorf("获取LevelDB模板列表数据失败, err: %w", err)
	}

	list := make(Infos, 0)
	// 如果没找到，将默认模板放入
	if errors.Is(err, leveldb.ErrNotFound) {
		list = append(list, DefaultTemplateInfo)
		jsonBytes, err := json.Marshal(list)
		if err != nil {
			return list.ToInfoDtos(), fmt.Errorf("编码LevelDB模板列表数据失败, err: %w", err)
		}
		err = t.leveldb.Put([]byte(constx.LevelDBTemplates), jsonBytes, nil)
		if err != nil {
			return list.ToInfoDtos(), fmt.Errorf("存入LevelDB模板列表数据失败, err: %w", err)
		}
	} else {
		err = json.Unmarshal(value, &list)
		if err != nil {
			return list.ToInfoDtos(), fmt.Errorf("解析LevelDB模板列表数据失败, err: %w", err)
		}
	}
	return list.ToInfoDtos(), nil
}

func (t *templateSrv) TemplateCreate(info Info) (err error) {
	var urlInfo TmplURL
	urlInfo, err = info.GitRemotePath.Parse()
	if err != nil {
		err = fmt.Errorf("URL解析失败, err: %w", err)
		return
	}

	// 防止并发请求
	t.l.Lock()
	defer t.l.Unlock()
	value, err := t.leveldb.Get([]byte(constx.LevelDBTemplates), nil)
	if err != nil && !errors.Is(err, leveldb.ErrNotFound) {
		err = fmt.Errorf("获取LevelDB模板列表数据失败, err: %w", err)
		return
	}

	list := make([]Info, 0)
	// 如果已经存在数据，那么进行解析
	if !errors.Is(err, leveldb.ErrNotFound) {
		err = json.Unmarshal(value, &list)
		if err != nil {
			return fmt.Errorf("解析LevelDB模板列表数据失败, err: %w", err)
		}

		for _, value := range list {
			// 该模板地址已存在，不允许插入
			if value.GitRemotePath == info.GitRemotePath {
				return fmt.Errorf("该模板地址已存在，git: %s", value.GitRemotePath)
			}
		}

	}
	list = append(list, Info{
		Name:          info.Name,
		GitRemotePath: info.GitRemotePath,
		Path:          system.EgoctlHome + "/egoctl/git" + urlInfo.Path,
	})
	jsonBytes, err := json.Marshal(list)
	if err != nil {
		return fmt.Errorf("编码LevelDB模板列表数据失败, err: %w", err)
	}
	err = t.leveldb.Put([]byte(constx.LevelDBTemplates), jsonBytes, nil)
	if err != nil {
		return fmt.Errorf("存入LevelDB模板列表数据失败, err: %w", err)
	}

	return
}

func (t *templateSrv) TemplateUpdate(info Info) (err error) {
	// 防止并发请求
	t.l.Lock()
	defer t.l.Unlock()
	value, err := t.leveldb.Get([]byte(constx.LevelDBTemplates), nil)
	if err != nil {
		err = fmt.Errorf("获取LevelDB模板列表数据失败, err: %w", err)
		return
	}

	list := make([]Info, 0)
	// 如果已经存在数据，那么进行解析
	err = json.Unmarshal(value, &list)
	if err != nil {
		return fmt.Errorf("解析LevelDB模板列表数据失败, err: %w", err)
	}

	isExist := false
	for _, value := range list {
		// 该模板地址已存在，不允许插入
		if value.GitRemotePath == info.GitRemotePath {
			isExist = true
		}
	}

	if !isExist {
		err = fmt.Errorf("不存在该git模板数据")
		return
	}

	listNew := make([]Info, 0)

	for _, value := range list {
		if value.GitRemotePath == info.GitRemotePath {
			value.Name = info.Name
			value.Path = info.Path
		}
		listNew = append(listNew, value)
	}

	jsonBytes, err := json.Marshal(listNew)
	if err != nil {
		return fmt.Errorf("编码LevelDB模板列表数据失败, err: %w", err)
	}
	err = t.leveldb.Put([]byte(constx.LevelDBTemplates), jsonBytes, nil)
	if err != nil {
		return fmt.Errorf("存入LevelDB模板列表数据失败, err: %w", err)
	}
	return
}

func (t *templateSrv) TemplateDelete(info Info) (err error) {
	// 防止并发请求
	t.l.Lock()
	defer t.l.Unlock()
	value, err := t.leveldb.Get([]byte(constx.LevelDBTemplates), nil)
	if err != nil {
		err = fmt.Errorf("获取LevelDB模板列表数据失败, err: %w", err)
		return
	}

	list := make([]Info, 0)
	// 如果已经存在数据，那么进行解析
	err = json.Unmarshal(value, &list)
	if err != nil {
		return fmt.Errorf("解析LevelDB模板列表数据失败, err: %w", err)
	}

	isExist := false
	for _, value := range list {
		// 该模板地址已存在，不允许插入
		if value.GitRemotePath == info.GitRemotePath {
			isExist = true
		}
	}

	if !isExist {
		err = fmt.Errorf("不存在该git模板数据")
		return
	}

	listNew := make([]Info, 0)

	for _, value := range list {
		if value.GitRemotePath == info.GitRemotePath {
			continue
		}
		listNew = append(listNew, value)
	}

	jsonBytes, err := json.Marshal(listNew)
	if err != nil {
		return fmt.Errorf("编码LevelDB模板列表数据失败, err: %w", err)
	}
	err = t.leveldb.Put([]byte(constx.LevelDBTemplates), jsonBytes, nil)
	if err != nil {
		return fmt.Errorf("存入LevelDB模板列表数据失败, err: %w", err)
	}
	return
}

func (t *templateSrv) TemplateSync(info InfoUniqId) (err error) {
	tInfo, err := t.TemplateInfo(info)
	if err != nil {
		return err
	}
	err = git.CloneORPullRepo(string(tInfo.GitRemotePath), tInfo.Path)
	if err != nil {
		return err
	}
	return nil
}

func (info Info) StatusText() (statusText string) {
	if !utils.IsDir(info.Path) {
		return "模板未下载"
	}

	rep, err := git.OpenRepository(info.Path)
	if err != nil {
		return fmt.Sprintf("打开文件失败: %v", err)
	}

	version, err := rep.GetVersion()
	if err != nil {
		return fmt.Sprintf("获取版本信息失败: %v", err)
	}
	return version

}

func (t *templateSrv) TemplateInfo(info InfoUniqId) (resp Info, err error) {
	// 防止并发请求
	t.l.Lock()
	defer t.l.Unlock()

	value, err := t.leveldb.Get([]byte(constx.LevelDBTemplates), nil)
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
		if value.GitRemotePath == info.GitRemotePath {
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
