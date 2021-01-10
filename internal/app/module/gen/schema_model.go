package gen

// parse get the model info
type ModelSchema struct {
	FieldName    string             `json:"name"`   // 字段名称
	FieldType    string             `json:"goType"` // go type
	FieldTags    map[string]SpecTag // map[gorm]{name:"gorm",origin:"not null;comment:"名称""}
	FieldComment string             `json:"comment"` // mysql comment
}

type ModelSchemas []ModelSchema
