// Copyright 2026 The Bucketeer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package storage_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// layout describes how a storage package lays out its two implementations.
type layout int

const (
	// subPkg: implementations live in mysql/ and postgres/ subpackages, the
	// interfaces and error sentinels in the parent package.
	subPkg layout = iota
	// sameDir: interfaces, errors, and both implementations share one package,
	// with implementations split across *_mysql.go and *_postgres.go files.
	sameDir
)

// pair is a storage package whose MySQL and PostgreSQL implementations must
// stay in sync.
//
// Only register domain storage packages that keep their shared interfaces and
// error sentinels in the parent directory (parsed as `parent` below). Do NOT
// register the low-level driver layer at pkg/storage/v2/{mysql,postgres}: it
// has no parent interface/error package, so the parity assumptions don't hold.
type pair struct {
	dir    string // module-root-relative directory of the storage package
	layout layout
}

var registry = []pair{
	{"pkg/account/storage/v2", subPkg},
	{"pkg/auditlog/storage/v2", subPkg},
	{"pkg/autoops/storage/v2", subPkg},
	{"pkg/coderef/storage", subPkg},
	{"pkg/environment/storage/v2", subPkg},
	{"pkg/eventcounter/storage/v2/dwh_database", subPkg},
	{"pkg/eventcounter/storage/v2/operational_database", subPkg},
	{"pkg/experiment/storage/v2", subPkg},
	{"pkg/experimentcalculator/storage/v2", subPkg},
	{"pkg/feature/storage/v2", subPkg},
	{"pkg/insights/storage/v2", subPkg},
	{"pkg/opsevent/storage/v2", subPkg},
	{"pkg/push/storage/v2", sameDir},
	{"pkg/subscription/storage/v2", subPkg},
	{"pkg/tag/storage", subPkg},
	{"pkg/team/storage", subPkg},
}

// dialectTokens are stripped from identifiers before comparing them across the
// two implementations, so dialect-specific names (e.g. ErrMySQLNoResultsFound
// and ErrPostgresNoResultsFound) compare as the same logical error.
var dialectTokens = []string{"MySQL", "Mysql", "PostgreSQL", "Postgres", "Pg"}

// sources holds the parsed parent / mysql / postgres source for one pair.
type sources struct {
	parent []*ast.File
	mysql  []*ast.File
	pg     []*ast.File
}

func TestStorageMySQLPostgresParity(t *testing.T) {
	t.Parallel()
	root := moduleRoot(t)
	for _, p := range registry {
		p := p
		t.Run(p.dir, func(t *testing.T) {
			t.Parallel()
			src := loadSources(t, root, p)
			checkInterfaceParity(t, src)
			checkErrorParity(t, src)
		})
	}
}

// checkInterfaceParity fails when one dialect implements a shared interface
// method that the other does not. Only interface methods that are actually
// implemented in the implementation packages are considered, so an interface
// not yet migrated to PostgreSQL (its impl still lives in the parent package)
// is not a false failure. The compiler already guarantees a wired-up impl
// satisfies its interface; this names the divergent method explicitly.
func checkInterfaceParity(t *testing.T, src sources) {
	methods := interfaceMethods(src.parent)
	mysqlMethods := receiverMethods(src.mysql)
	pgMethods := receiverMethods(src.pg)
	for _, m := range methods {
		inMySQL := mysqlMethods[m]
		inPg := pgMethods[m]
		if inMySQL != inPg {
			t.Errorf(
				"interface method %q implemented by mysql=%v but postgres=%v; "+
					"both implementations must define it",
				m, inMySQL, inPg,
			)
		}
	}
}

// checkErrorParity fails when an error sentinel declared in the parent package
// is referenced by one implementation but not the other. Names are normalized
// so dialect-specific sentinels are treated as the same logical error.
func checkErrorParity(t *testing.T, src sources) {
	parentErrs := parentErrorNames(src.parent)
	if len(parentErrs) == 0 {
		return
	}
	mysqlRefs := referencedErrorNames(src.mysql, parentErrs)
	pgRefs := referencedErrorNames(src.pg, parentErrs)

	// Group parent errors by normalized name; an error is "used" by an
	// implementation when any member of its group is referenced there.
	groups := map[string][]string{}
	for _, e := range parentErrs {
		n := normalize(e)
		groups[n] = append(groups[n], e)
	}
	for norm, members := range groups {
		usedByMySQL := anyReferenced(members, mysqlRefs)
		usedByPg := anyReferenced(members, pgRefs)
		if usedByMySQL != usedByPg {
			t.Errorf(
				"error %q referenced by mysql=%v but postgres=%v (members %v); "+
					"both implementations must handle it",
				norm, usedByMySQL, usedByPg, members,
			)
		}
	}
}

func loadSources(t *testing.T, root string, p pair) sources {
	t.Helper()
	base := filepath.Join(root, filepath.FromSlash(p.dir))
	switch p.layout {
	case subPkg:
		return sources{
			parent: parseDir(t, base),
			mysql:  parseDir(t, filepath.Join(base, "mysql")),
			pg:     parseDir(t, filepath.Join(base, "postgres")),
		}
	case sameDir:
		var src sources
		for path, f := range parseDirFiles(t, base) {
			name := strings.ToLower(filepath.Base(path))
			switch {
			case strings.Contains(name, "postgres"):
				src.pg = append(src.pg, f)
			case strings.Contains(name, "mysql"):
				src.mysql = append(src.mysql, f)
			default:
				src.parent = append(src.parent, f)
			}
		}
		return src
	default:
		t.Fatalf("unknown layout for %s", p.dir)
		return sources{}
	}
}

// parseDir parses the non-test Go files directly in dir (not its subdirs).
func parseDir(t *testing.T, dir string) []*ast.File {
	t.Helper()
	files := parseDirFiles(t, dir)
	out := make([]*ast.File, 0, len(files))
	for _, f := range files {
		out = append(out, f)
	}
	return out
}

func parseDirFiles(t *testing.T, dir string) map[string]*ast.File {
	t.Helper()
	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		t.Fatalf("registered storage dir missing: %s", dir)
	}
	notTest := func(fi fs.FileInfo) bool { return !strings.HasSuffix(fi.Name(), "_test.go") }
	pkgs, err := parser.ParseDir(token.NewFileSet(), dir, notTest, 0)
	if err != nil {
		t.Fatalf("failed to parse %s: %v", dir, err)
	}
	out := map[string]*ast.File{}
	for name, pkg := range pkgs {
		if strings.HasSuffix(name, "_test") {
			continue
		}
		for path, f := range pkg.Files {
			out[path] = f
		}
	}
	return out
}

// interfaceMethods returns every method name declared on any interface type.
func interfaceMethods(files []*ast.File) []string {
	var methods []string
	for _, f := range files {
		ast.Inspect(f, func(n ast.Node) bool {
			it, ok := n.(*ast.InterfaceType)
			if !ok {
				return true
			}
			for _, m := range it.Methods.List {
				for _, name := range m.Names {
					methods = append(methods, name.Name)
				}
			}
			return true
		})
	}
	return methods
}

// receiverMethods returns the set of method names defined on any type.
func receiverMethods(files []*ast.File) map[string]bool {
	out := map[string]bool{}
	for _, f := range files {
		for _, decl := range f.Decls {
			if fd, ok := decl.(*ast.FuncDecl); ok && fd.Recv != nil {
				out[fd.Name.Name] = true
			}
		}
	}
	return out
}

// parentErrorNames returns top-level var names that look like error sentinels.
func parentErrorNames(files []*ast.File) []string {
	var names []string
	for _, f := range files {
		for _, decl := range f.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.VAR {
				continue
			}
			for _, spec := range gd.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				for _, name := range vs.Names {
					if strings.HasPrefix(name.Name, "Err") {
						names = append(names, name.Name)
					}
				}
			}
		}
	}
	return names
}

// referencedErrorNames returns the subset of candidate error names that the
// files genuinely reference as the shared sentinel — not merely an identifier
// that happens to share the name.
//
// A candidate counts as referenced when it appears either:
//   - as a qualified selector target (e.g. v2as.ErrFoo), which can only be a
//     reference to another package's exported symbol; or
//   - as a bare identifier that is NOT declared locally in these files.
//
// This excludes the case the parity check could otherwise be fooled by: a local
// variable, parameter, or field that coincidentally shares a sentinel's name.
// Such a local declaration is recorded in declaredNames and therefore ignored,
// so the bare-identifier path only matches same-package sentinels (the sameDir
// layout) that are declared in the parent file rather than the impl files.
func referencedErrorNames(files []*ast.File, candidates []string) map[string]bool {
	declared := declaredNames(files)
	allIdents := map[string]bool{}
	selectorTargets := map[string]bool{}
	for _, f := range files {
		ast.Inspect(f, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.SelectorExpr:
				selectorTargets[x.Sel.Name] = true
			case *ast.Ident:
				allIdents[x.Name] = true
			}
			return true
		})
	}
	out := map[string]bool{}
	for _, c := range candidates {
		if selectorTargets[c] || (allIdents[c] && !declared[c]) {
			out[c] = true
		}
	}
	return out
}

// declaredNames returns identifier names introduced as declarations in the
// given files: short var definitions, var/const specs, type names, function
// names, and field/parameter/result names.
func declaredNames(files []*ast.File) map[string]bool {
	out := map[string]bool{}
	for _, f := range files {
		ast.Inspect(f, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.AssignStmt:
				if x.Tok == token.DEFINE {
					for _, lhs := range x.Lhs {
						if id, ok := lhs.(*ast.Ident); ok {
							out[id.Name] = true
						}
					}
				}
			case *ast.ValueSpec:
				for _, id := range x.Names {
					out[id.Name] = true
				}
			case *ast.TypeSpec:
				out[x.Name.Name] = true
			case *ast.FuncDecl:
				out[x.Name.Name] = true
			case *ast.Field:
				for _, id := range x.Names {
					out[id.Name] = true
				}
			}
			return true
		})
	}
	return out
}

// ---- helpers ------------------------------------------------------------

func anyReferenced(names []string, refs map[string]bool) bool {
	for _, n := range names {
		if refs[n] {
			return true
		}
	}
	return false
}

func normalize(s string) string {
	for _, tok := range dialectTokens {
		s = strings.ReplaceAll(s, tok, "")
	}
	return s
}

func moduleRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not locate go.mod (module root)")
		}
		dir = parent
	}
}
