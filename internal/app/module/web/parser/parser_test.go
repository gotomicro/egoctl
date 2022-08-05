package parser

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func Test_astParser_parserStruct(t *testing.T) {
	ast, _ := AstParserBuild(UserOption{
		ScaffoldDSLContent: "testdata/user/ego.go",
	}, TmplOption{})
	if len(ast.modelArr) != 1 {
		t.Fatalf("got %d model arr, want 1", len(ast.modelArr))
	}
	t.Log(ast.modelArr)
}

func Test_astParser_parserStructTag(t *testing.T) {
	ast, _ := AstParserBuild(UserOption{
		ScaffoldDSLContent: "testdata/user/ego.go",
	}, TmplOption{})
	if len(ast.modelArr) != 1 {
		t.Fatalf("got %d model arr, want 1", len(ast.modelArr))
	}
	spew.Dump(ast.modelArr)
}
