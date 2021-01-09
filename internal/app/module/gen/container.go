package gen

import (
	"bytes"
	"errors"
	"github.com/davecgh/go-spew/spew"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"sort"
)

type Container struct {
	objectM map[string]*SpecType
	Src     string
}

var DefaultEgoctlPro = &Container{
	objectM :make(map[string]*SpecType),
}

func (c *Container) Run() error {
	fSet := token.NewFileSet()
	readinfo , err := ioutil.ReadFile("egoctl.go")
	if err != nil {
		panic(err)
	}
	c.Src = string(readinfo)
	// strings.NewReader
	f, err := parser.ParseFile(fSet,"" , bytes.NewReader(readinfo), parser.ParseComments)
	if err != nil {
		panic(err)
	}

	commentMap := ast.NewCommentMap(fSet, f, f.Comments)
	f.Comments = commentMap.Filter(f).Comments()

	scope := f.Scope
	if scope == nil {
		return errors.New("struct nil")
	}
	objects := scope.Objects
	structs := make([]*SpecType, 0)
	for structName, obj := range objects {
		st, err := c.parseObject(structName, obj)
		if err != nil {
			return err
		}
		structs = append(structs, st)
	}
	sort.Slice(structs, func(i, j int) bool {
		return structs[i].Name < structs[j].Name
	})

	resp := make([]SpecType, 0)
	for _, item := range structs {
		resp = append(resp, *item)
	}

	spew.Dump("f", resp)
	return nil
}

func (c *Container) parseObject(structName string, obj *ast.Object) (*SpecType, error) {
	if data, ok := c.objectM[structName]; ok {
		return data, nil
	}
	var st SpecType
	st.Name = structName
	if obj.Decl == nil {
		c.objectM[structName] = &st
		return &st, nil
	}
	decl, ok := obj.Decl.(*ast.TypeSpec)
	if !ok {
		c.objectM[structName] = &st
		return &st, nil
	}
	if decl.Type == nil {
		c.objectM[structName] = &st
		return &st, nil
	}
	tp, ok := decl.Type.(*ast.StructType)
	if !ok {
		c.objectM[structName] = &st
		return &st, nil
	}
	fields := tp.Fields
	if fields == nil {
		c.objectM[structName] = &st
		return &st, nil
	}
	fieldList := fields.List
	members, err := c.parseFields(fieldList)
	if err != nil {
		return nil, err
	}
	st.Members = members
	c.objectM[structName] = &st
	return &st, nil
}

func (c *Container) parseFields(fields []*ast.Field) ([]SpecMember, error) {
	members := make([]SpecMember, 0)
	for _, field := range fields {
		docs := parseCommentOrDoc(field.Doc)
		comments := parseCommentOrDoc(field.Comment)
		name := parseName(field.Names)
		tp, stringExpr, err := c.parseType(field.Type)
		if err != nil {
			return nil, err
		}
		tag := parseTag(field.Tag)
		isInline := name == ""
		if isInline {
			var err error
			name, err = c.getInlineName(tp)
			if err != nil {
				return nil, err
			}
		}
		members = append(members, SpecMember{
			Name:     name,
			Type:     stringExpr,
			Expr:     tp,
			Tag:      tag,
			Comments: comments,
			Docs:     docs,
			IsInline: isInline,
		})

	}
	return members, nil
}
