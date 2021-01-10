package gen

type RenderInfo struct {
	Module       string
	ModelName    string
	Option       UserOption
	Content      ModelSchemas
	Descriptor   Descriptor
	TmplPath     string
	GenerateTime string
}
