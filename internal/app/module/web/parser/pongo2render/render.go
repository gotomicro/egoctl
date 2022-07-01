package pongo2render

import (
	"net/http"
	"path"

	"github.com/gotomicro/egoctl/internal/app/module/web/parser/pongo2"
)

//	var render = pongo2render.NewRender("./templates")
//
//	http.HandleFunc("/m", func(w http.ResponseWriter, req *http.Request) {
//		render.HTML(w, 200, "index.html", pongo2.Context{"aa": "eeeeeee"})
//	})
//	http.ListenAndServe(":9005", nil)

// --------------------------------------------------------------------------------
var htmlContentType = []string{"text/html; charset=utf-8"}

type Render struct {
	TemplateDir string
	Cache       bool
}

func NewRender(templateDir string) *Render {
	var r = &Render{}
	r.TemplateDir = templateDir
	return r
}

func (this *Render) Template(name string) *Template {
	var template *pongo2.Template
	var filename string
	if len(this.TemplateDir) > 0 {
		filename = path.Join(this.TemplateDir, name)
	} else {
		filename = name
	}

	if this.Cache {
		template = pongo2.Must(pongo2.DefaultSet.FromCache(filename))
	} else {
		template = pongo2.Must(pongo2.DefaultSet.FromFile(filename))
	}

	if template == nil {
		panic("template " + name + " not exists")
		return nil
	}

	var r = &Template{}
	r.template = template
	return r
}

func (this *Render) TemplateFromString(tpl string) *Template {
	var template = pongo2.Must(pongo2.DefaultSet.FromString(tpl))
	var r = &Template{}
	r.template = template
	return r
}

func (this *Render) HTML(w http.ResponseWriter, status int, name string, data any) {
	w.WriteHeader(status)
	this.Template(name).ExecuteWriter(w, data)
}

// --------------------------------------------------------------------------------
type Template struct {
	template *pongo2.Template
	context  pongo2.Context
}

func (this *Template) ExecuteWriter(w http.ResponseWriter, data any) (err error) {
	WriteContentType(w, htmlContentType)
	this.context = DataToContext(data)
	err = this.template.ExecuteWriter(this.context, w)
	return err
}

func (this *Template) Execute(data any) (string, error) {
	this.context = DataToContext(data)
	return this.template.Execute(this.context)
}

// --------------------------------------------------------------------------------
func WriteContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}

func DataToContext(data interface{}) pongo2.Context {
	var ctx pongo2.Context
	if data != nil {
		switch data.(type) {
		case pongo2.Context:
			ctx = data.(pongo2.Context)
		case map[string]any:
			ctx = data.(map[string]any)
		}
	}
	return ctx
}
