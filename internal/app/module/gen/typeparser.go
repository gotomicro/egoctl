package gen

import (
	"errors"
	"fmt"
	"go/ast"
	"strings"
)

var (
	ErrStructNotFound      = errors.New("struct not found")
	ErrUnSupportInlineType = errors.New("unsupport inline type")
	interfaceExpr          = `interface{}`
	objectM                = make(map[string]*SpecType)
)

func (c *Container) getInlineName(tp interface{}) (string, error) {
	switch v := tp.(type) {
	case *SpecType:
		return v.Name, nil
	case *SpecPointerType:
		return c.getInlineName(v.Star)
	case *SpecStructType:
		return v.StringExpr, nil
	default:
		return "", ErrUnSupportInlineType
	}
}

func (sp *Container) getInlineTypePrefix(tp interface{}) (string, error) {
	if tp == nil {
		return "", nil
	}
	switch tp.(type) {
	case *ast.Ident:
		return "", nil
	case *ast.StarExpr:
		return "*", nil
	case *ast.TypeSpec:
		return "", nil
	default:
		return "", ErrUnSupportInlineType
	}
}

func parseTag(basicLit *ast.BasicLit) string {
	if basicLit == nil {
		return ""
	}
	fmt.Printf("basicLit.Value--------------->"+"%+v\n", basicLit.Value)
	return basicLit.Value
}

// returns
// resp1: type can convert to *SpecPointerType|*SpecBasicType|*SpecMapType|*SpecArrayType|*SpecInterfaceType
// resp2: type's string expression,like int、string、[]int64、map[string]User、*User
// resp3: error
func (sp *Container) parseType(expr ast.Expr) (interface{}, string, error) {
	if expr == nil {
		return nil, "", errors.New("parse error ")
	}
	exprStr := sp.Src[expr.Pos():expr.End()]
	switch v := expr.(type) {
	case *ast.StarExpr:
		star, stringExpr, err := sp.parseType(v.X)
		if err != nil {
			return nil, "", err
		}
		e := fmt.Sprintf("*%s", stringExpr)
		return &SpecPointerType{Star: star, StringExpr: e}, e, nil
	case *ast.Ident:
		if isBasicType(v.Name) {
			return &SpecBasicType{Name: v.Name, StringExpr: v.Name}, v.Name, nil
		} else if v.Obj != nil {
			obj := v.Obj
			if obj.Name != v.Name { // 防止引用自己而无限递归
				specType, err := sp.parseObject(v.Name, v.Obj)
				if err != nil {
					return nil, "", err
				} else {
					return specType, v.Obj.Name, nil
				}
			} else {
				inlineType, err := sp.getInlineTypePrefix(obj.Decl)
				if err != nil {
					return nil, "", err
				}
				return &SpecStructType{
					StringExpr: fmt.Sprintf("%s%s", inlineType, v.Name),
				}, v.Name, nil
			}
		} else {
			return nil, "", fmt.Errorf(" [%s] - member is not exist, expr is %s", v.Name, exprStr)
		}
	case *ast.MapType:
		key, keyStringExpr, err := sp.parseType(v.Key)
		if err != nil {
			return nil, "", err
		}

		value, valueStringExpr, err := sp.parseType(v.Value)
		if err != nil {
			return nil, "", err
		}

		keyType, ok := key.(*SpecBasicType)
		if !ok {
			return nil, "", fmt.Errorf("[%+v] - unsupported type of map key, expr is  %s", v.Key, exprStr)
		}

		e := fmt.Sprintf("map[%s]%s", keyStringExpr, valueStringExpr)
		return &SpecMapType{
			Key:        keyType.Name,
			Value:      value,
			StringExpr: e,
		}, e, nil
	case *ast.ArrayType:
		arrayType, stringExpr, err := sp.parseType(v.Elt)
		if err != nil {
			return nil, "", err
		}

		e := fmt.Sprintf("[]%s", stringExpr)
		return &SpecArrayType{ArrayType: arrayType, StringExpr: e}, e, nil
	case *ast.InterfaceType:
		return &SpecInterfaceType{StringExpr: interfaceExpr}, interfaceExpr, nil
	case *ast.ChanType:
		return nil, "", errors.New("[chan] - unsupported type, expr is " + exprStr)
	case *ast.FuncType:
		return nil, "", errors.New("[func] - unsupported type, expr is " + exprStr)
	case *ast.StructType: // todo can optimize
		return nil, "", errors.New("[struct] - unsupported inline struct type, expr is " + exprStr)
	case *ast.SelectorExpr:
		x := v.X
		sel := v.Sel
		xIdent, ok := x.(*ast.Ident)
		if ok {
			name := xIdent.Name
			if name != "time" && sel.Name != "Time" {
				return nil, "", fmt.Errorf("[outter package] - package: %s, unsupport type", exprStr)
			}

			tm := fmt.Sprintf("time.Time")
			return &SpecTimeType{
				StringExpr: tm,
			}, tm, nil
		}
		return nil, "", errors.New("parse error " + exprStr)
	default:
		return nil, "", errors.New("parse error " + exprStr)
	}
}

func isBasicType(tp string) bool {
	switch tp {
	case
		"bool",
		"uint8",
		"uint16",
		"uint32",
		"uint64",
		"int8",
		"int16",
		"int32",
		"int64",
		"float32",
		"float64",
		"complex64",
		"complex128",
		"string",
		"int",
		"uint",
		"uintptr",
		"byte",
		"rune",
		"Type",
		"Type1",
		"IntegerType",
		"FloatType",
		"ComplexType":
		return true
	default:
		return false
	}
}
func parseName(names []*ast.Ident) string {
	if len(names) == 0 {
		return ""
	}
	name := names[0]
	return parseIdent(name)
}

func parseIdent(ident *ast.Ident) string {
	if ident == nil {
		return ""
	}
	return ident.Name
}

func parseCommentOrDoc(cg *ast.CommentGroup) []string {
	if cg == nil {
		return nil
	}
	comments := make([]string, 0)
	for _, comment := range cg.List {
		if comment == nil {
			continue
		}
		text := strings.TrimSpace(comment.Text)
		if text == "" {
			continue
		}
		comments = append(comments, text)
	}
	return comments
}
