package gen

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
	Tag  string
	// Deprecated
	Comment string // 换成标准struct中将废弃
	// 成员尾部注释说明
	Comments []string
	// 成员头顶注释说明
	Docs     []string
	IsInline bool
}

type SpecType struct {
	Name        string
	Annotations []SpecAnnotation
	Members     []SpecMember
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

