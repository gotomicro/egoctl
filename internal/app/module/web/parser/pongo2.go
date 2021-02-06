package parser

import (
	"github.com/flosch/pongo2"
	"github.com/gotomicro/egoctl/utils"
	"strings"
	"unicode/utf8"
)

func init() {
	_ = pongo2.RegisterFilter("lowerFirst", pongo2LowerFirst)
	_ = pongo2.RegisterFilter("upperFirst", pongo2UpperFirst)
	_ = pongo2.RegisterFilter("snakeString", pongo2SnakeString)
	_ = pongo2.RegisterFilter("camelString", pongo2CamelString)
	_ = pongo2.RegisterFilter("fieldsGetPrimaryKey", pongo2ModelFieldsGetPrimaryKey) // 根据字段数组获取主键
	_ = pongo2.RegisterFilter("fieldsExist", pongo2ModelFieldsExist)
	_ = pongo2.RegisterFilter("fieldGetTag", pongo2ModelFieldGetTag)
}

func pongo2ModelFieldsExist(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	arr, flag := in.Interface().(ModelSchemas)
	if !flag {
		return pongo2.AsSafeValue(""), nil
	}

	for _, info := range arr {
		if info.FieldName == param.String() {
			return pongo2.AsSafeValue(true), nil
		}

	}
	return pongo2.AsSafeValue(false), nil
}

func pongo2ModelFieldGetTag(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	info, flag := in.Interface().(ModelSchema)
	if !flag {
		return pongo2.AsSafeValue(""), nil
	}

	tag, flag := info.FieldTags[param.String()]
	if !flag {
		return pongo2.AsSafeValue(""), nil
	}

	return pongo2.AsSafeValue(tag.Origin), nil
}

func pongo2ModelFieldsGetPrimaryKey(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	arr, flag := in.Interface().(ModelSchemas)
	if !flag {
		return pongo2.AsSafeValue("Id"), nil
	}

	primaryKey := "Id"

	for _, info := range arr {
		tags, flag := info.FieldTags["ego"]
		if !flag {
			continue
		}
		for _, tag := range tags.Value {
			if tag == "primary_key" {
				primaryKey = info.FieldName
			}
		}
	}
	return pongo2.AsSafeValue(primaryKey), nil
}

func pongo2LowerFirst(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if in.Len() <= 0 {
		return pongo2.AsSafeValue(""), nil
	}
	t := in.String()
	r, size := utf8.DecodeRuneInString(t)
	return pongo2.AsSafeValue(strings.ToLower(string(r)) + t[size:]), nil
}

func pongo2UpperFirst(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if in.Len() <= 0 {
		return pongo2.AsSafeValue(""), nil
	}
	t := in.String()
	return pongo2.AsSafeValue(strings.Replace(t, string(t[0]), strings.ToUpper(string(t[0])), 1)), nil
}

// snake string, XxYy to xx_yy
func pongo2SnakeString(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if in.Len() <= 0 {
		return pongo2.AsSafeValue(""), nil
	}
	t := in.String()
	return pongo2.AsSafeValue(utils.SnakeString(t)), nil
}

// snake string, XxYy to xx_yy
func pongo2CamelString(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if in.Len() <= 0 {
		return pongo2.AsSafeValue(""), nil
	}
	t := in.String()
	return pongo2.AsSafeValue(utils.CamelString(t)), nil
}

//func upperFirst(str string) string {
//	return strings.Replace(str, string(str[0]), strings.ToUpper(string(str[0])), 1)
//}

func lowerFirst(str string) string {
	return strings.Replace(str, string(str[0]), strings.ToLower(string(str[0])), 1)
}
