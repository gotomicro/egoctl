package gen

import (
	"testing"
)

func Test_astParser_parserStruct(t *testing.T) {
	ast := AstParserBuild(UserOption{
		ScaffoldDSLFile: "testdata/user/ego.go",
	}, TmplOption{})
	if len(ast.modelArr) != 1 {
		t.Fatalf("got %d model arr, want 1", len(ast.modelArr))
	}
	t.Log(ast.modelArr)
}
