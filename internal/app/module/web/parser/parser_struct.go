package parser

import (
	"go/ast"
)

func (c *astParser) parseObject(structName string, obj *ast.Object) (*SpecType, error) {
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

func (c *astParser) parseFields(fields []*ast.Field) ([]SpecMember, error) {
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
