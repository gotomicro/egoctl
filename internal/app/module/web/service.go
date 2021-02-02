package web

import "sync"

const (
	LevelDBProjects  = "projects"
	LevelDBProjectId = "projectId"
)

var ProjectSrv *projectSrv

type projectSrv struct {
	l sync.RWMutex
}

func init() {
	ProjectSrv = &projectSrv{}
}

func (p *projectSrv) ProjectCreate() {

}
