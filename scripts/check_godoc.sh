#!/usr/bin/env bash
# check_godoc.sh — verify that every exported symbol has a doc comment.
# Uses the AST-based scanner (same as /tmp/opencode/godocscan.go) to avoid
# false positives from go/doc's grouping of const/var blocks.
# Used by CI and `make check-godoc`.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$PROJECT_DIR"

TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

cat > "$TMPDIR/scan.go" <<'EOF'
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

func hasDoc(cg *ast.CommentGroup) bool {
	return cg != nil && strings.TrimSpace(cg.Text()) != ""
}

func main() {
	root := os.Args[1]
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, root, func(fi os.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, parser.ParseComments)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fail := false
	for _, pkg := range pkgs {
		for name, f := range pkg.Files {
			if strings.Contains(name, "_test.go") {
				continue
			}
			for _, d := range f.Decls {
				switch dd := d.(type) {
				case *ast.FuncDecl:
					if dd.Name.IsExported() && !hasDoc(dd.Doc) {
						// Skip methods on unexported receiver types.
						if dd.Recv != nil && len(dd.Recv.List) > 0 {
							rt := ""
							switch t := dd.Recv.List[0].Type.(type) {
							case *ast.StarExpr:
								if id, ok := t.X.(*ast.Ident); ok {
									rt = id.Name
								}
							case *ast.Ident:
								rt = t.Name
							}
							if rt == "" || !ast.IsExported(rt) {
								continue
							}
						}
						fmt.Printf("MISSING DOC: %s:%d func %s\n", name, fset.Position(dd.Pos()).Line, dd.Name.Name)
						fail = true
					}
				case *ast.GenDecl:
					if dd.Tok == token.TYPE || dd.Tok == token.VAR || dd.Tok == token.CONST {
						for _, s := range dd.Specs {
							switch sp := s.(type) {
							case *ast.TypeSpec:
								if sp.Name.IsExported() && !hasDoc(sp.Doc) && !hasDoc(dd.Doc) {
									fmt.Printf("MISSING DOC: %s:%d type %s\n", name, fset.Position(sp.Pos()).Line, sp.Name.Name)
									fail = true
								}
							case *ast.ValueSpec:
								for _, id := range sp.Names {
									if id.IsExported() && !hasDoc(sp.Doc) && !hasDoc(dd.Doc) {
										fmt.Printf("MISSING DOC: %s:%d %s %s\n", name, fset.Position(id.Pos()).Line, dd.Tok, id.Name)
										fail = true
									}
								}
							}
						}
					}
				}
			}
		}
	}
	if fail {
		os.Exit(1)
	}
}
EOF

# Run from each package directory
MISSING=0
for d in $(go list -f '{{.Dir}}' ./...); do
	if ! (cd "$PROJECT_DIR" && go run "$TMPDIR/scan.go" "$d") >/dev/null 2>&1; then
		go run "$TMPDIR/scan.go" "$d" 2>&1 || true
		MISSING=1
	fi
done

if [ "$MISSING" -ne 0 ]; then
	echo "ERROR: Some exported symbols lack doc comments." >&2
	exit 1
fi
echo "PASS: All exported symbols have doc comments."