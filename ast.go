package codescore

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
)

// ASTSerialize takes a path to a Go source file and returns a
// serialized AST string using the go/ast printer representation.
// It returns only the serialized AST output and a possible parse error.
func ASTSerialize(path string) (string, error) {
	fset := token.NewFileSet()

	snippet, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	ast.Fprint(&buf, fset, snippet, nil)

	return buf.String(), nil
}
