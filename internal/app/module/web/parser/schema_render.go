package parser

type RenderInfo struct {
	ModelNames   []string     `json:"modelNames"` // 所有model names
	ModelName    string       `json:"modelName"`  // 当前model names
	Module       string       `json:"-"`
	TmplPath     string       `json:"tmplPath"`
	GenerateTime string       `json:"generateTime"`
	Option       UserOption   `json:"-"`
	Content      ModelSchemas `json:"content"`
	Descriptor   Descriptor   `json:"-"`
}
