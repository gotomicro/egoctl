package gen

import (
	"errors"
	"fmt"
	"go/ast"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrStructNotFound      = errors.New("struct not found")
	ErrUnSupportInlineType = errors.New("unsupport inline type")
	interfaceExpr          = `interface{}`
	objectM                = make(map[string]*SpecType)
)

func (c *astParser) getInlineName(tp interface{}) (string, error) {
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

func (sp *astParser) getInlineTypePrefix(tp interface{}) (string, error) {
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

func parseTag(basicLit *ast.BasicLit) SpecTags {
	if basicLit == nil {
		return SpecTags{}
	}
	tags := SpecTags{
		Origin: reflect.StructTag(basicLit.Value),
		Value:  parseLineTag(basicLit.Value),
	}
	return tags
}

func (tag StructTag) Lookup() []SpecTag {
	// When modifying this code, also update the validateStructTag code
	// in cmd/vet/structtag.go.

	output := make([]SpecTag, 0)
	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		value, err := strconv.Unquote(qvalue)
		if err != nil {
			break
		}
		output = append(output, SpecTag{
			Name:   name,
			Origin: value,
			Value:  strings.Split(value, ";"),
		})
	}
	return output
}

// 需要trim `gorm:"not null;PRIMARY_KEY;comment:'用户uid'" json:"uid"`
func parseLineTag(value string) []SpecTag {
	value = strings.TrimSuffix(value, "`")
	value = strings.TrimPrefix(value, "`")
	info := StructTag(value)
	info.Lookup()
	return info.Lookup()
}

// returns
// resp1: type can convert to *SpecPointerType|*SpecBasicType|*SpecMapType|*SpecArrayType|*SpecInterfaceType
// resp2: type's string expression,like int、string、[]int64、map[string]User、*User
// resp3: error
func (sp *astParser) parseType(expr ast.Expr) (interface{}, string, error) {
	if expr == nil {
		return nil, "", errors.New("parse error ")
	}
	exprStr := sp.readContent[expr.Pos():expr.End()]
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
