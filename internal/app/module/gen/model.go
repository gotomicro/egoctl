package gen

import (
	"reflect"
)

type SpecAnnotation struct {
	Name       string
	Properties map[string]string
	Value      string
}

type SpecMember struct {
	Annotations []SpecAnnotation
	Name        string
	// 数据类型字面值，如：string、map[int]string、[]int64、[]*User
	Type string
	// it can be asserted as BasicType: int、bool、
	// PointerType: *string、*User、
	// MapType: map[${BasicType}]interface、
	// ArrayType:[]int、[]User、[]*User
	// InterfaceType: interface{}
	// Type
	Expr interface{}
	Tag  SpecTags
	// 成员尾部注释说明
	Comments []string
	// 成员头顶注释说明
	Docs     []string
	IsInline bool
}

// A StructTag is the tag string in a struct field.
//
// By convention, tag strings are a concatenation of
// optionally space-separated key:"value" pairs.
// Each key is a non-empty string consisting of non-control
// characters other than space (U+0020 ' '), quote (U+0022 '"'),
// and colon (U+003A ':').  Each value is quoted using U+0022 '"'
// characters and Go string literal syntax.
type StructTag string

type SpecTags struct {
	Origin reflect.StructTag // 原始数据，所有tag信息
	Value  []SpecTag
}

type SpecTag struct {
	Name   string   // tag名称
	Origin string   // 原始数据，例如 json的名称
	Value  []string // tag的数组值
}

type SpecType struct {
	Name        string
	Annotations []SpecAnnotation
	Members     []SpecMember
}

func (content SpecType) ToModelInfos() (output []ModelSchema) {
	output = make([]ModelSchema, 0)
	for _, member := range content.Members {
		comment := content.Name
		if len(member.Comments) > 0 {
			comment = member.Comments[0]
		}

		tags := make(map[string]SpecTag, 0)
		for _, value := range member.Tag.Value {
			tags[value.Name] = value
		}

		m := ModelSchema{
			FieldName:    member.Name,
			FieldType:    member.Type,
			FieldTags:    tags,
			FieldComment: comment,
		}
		output = append(output, m)
	}
	return
}

type SpecPointerType struct {
	StringExpr string
	// it can be asserted as BasicType: int、bool、
	// PointerType: *string、*User、
	// MapType: map[${BasicType}]interface、
	// ArrayType:[]int、[]User、[]*User
	// InterfaceType: interface{}
	// Type
	Star interface{}
}
type (
	SpecMapType struct {
		StringExpr string
		// only support the BasicType
		Key string
		// it can be asserted as BasicType: int、bool、
		// PointerType: *string、*User、
		// MapType: map[${BasicType}]interface、
		// ArrayType:[]int、[]User、[]*User
		// InterfaceType: interface{}
		// Type
		Value interface{}
	}
	SpecArrayType struct {
		StringExpr string
		// it can be asserted as BasicType: int、bool、
		// PointerType: *string、*User、
		// MapType: map[${BasicType}]interface、
		// ArrayType:[]int、[]User、[]*User
		// InterfaceType: interface{}
		// Type
		ArrayType interface{}
	}
	SpecInterfaceType struct {
		StringExpr string
		// do nothing,just for assert
	}
	SpecTimeType struct {
		StringExpr string
	}
	SpecStructType struct {
		StringExpr string
	}

	// 系统预设基本数据类型
	SpecBasicType struct {
		StringExpr string
		Name       string
	}
)
