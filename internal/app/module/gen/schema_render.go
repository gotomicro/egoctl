package gen

type RenderInfo struct {
	ModelNames   []string // 所有model names
	ModelName    string   // 当前model names
	Module       string
	Option       UserOption
	Content      ModelSchemas
	Descriptor   Descriptor
	TmplPath     string
	GenerateTime string
}
