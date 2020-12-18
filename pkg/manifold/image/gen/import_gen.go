package gen

import (
	"bytes"
	"fmt"
	"path"
	"strings"

	. "github.com/dave/jennifer/jen"
)

func qualObj(id string) string {
	return fmt.Sprintf("worksite/pkg/obj/%s", id)
}

func ObjectPackageIndex(objIDs []string) ([]byte, error) {
	f := NewFile("obj")
	f.Comment("GENERATED INDEX")

	f.ImportName("github.com/progrium/zt100/pkg/manifold/library", "library")
	for _, objId := range objIDs {
		f.ImportName(qualObj(objId), objId)
		f.ImportAlias(qualObj(objId), objId)
	}

	f.Func().Id("relPath").Params(Id("subpath").Id("string")).Id("string").Block(
		List(Id("_"), Id("filename"), Id("_"), Id("_")).Op(":=").Qual("runtime", "Caller").Call(Lit(1)),
		Return(Qual("path", "Join").Call(Qual("path", "Dir").Call(Id("filename")), Id("subpath"))),
	)
	f.Line()
	var registrations []Code
	for _, objId := range objIDs {
		registrations = append(registrations,
			Qual("github.com/progrium/zt100/pkg/manifold/library", "Register").Call(
				Op("&").Qual(qualObj(objId), "Main").Values(),
				Lit(objId),
				Id("relPath").Call(Lit(path.Join(objId, "component.go"))),
				Lit("fas fa-user-robot"),
			),
		)
	}
	f.Func().Id("init").Params().Block(registrations...)
	buf := &bytes.Buffer{}
	err := f.Render(buf)
	return buf.Bytes(), err
}

func UserPackageIndex(pkgNames []string) ([]byte, error) {
	indexTmpl := `package usr

// GENERATED INDEX

import(
	%s
)
`
	importTmpl := `_ "worksite/pkg/usr/%s"`

	// f := NewFile("usr")
	// f.Comment("GENERATED INDEX")

	// f.ImportName("github.com/progrium/zt100/pkg/manifold/library", "library")
	// for _, objId := range objIDs {
	// 	f.ImportName(qualObj(objId), objId)
	// 	f.ImportAlias(qualObj(objId), objId)
	// }

	var imports []string
	for _, name := range pkgNames {
		imports = append(imports, fmt.Sprintf(importTmpl, name))
	}

	// f.Func().Id("relPath").Params(Id("subpath").Id("string")).Id("string").Block(
	// 	List(Id("_"), Id("filename"), Id("_"), Id("_")).Op(":=").Qual("runtime", "Caller").Call(Lit(1)),
	// 	Return(Qual("path", "Join").Call(Qual("path", "Dir").Call(Id("filename")), Id("subpath"))),
	// )
	// f.Line()
	// var registrations []Code
	// for _, objId := range objIDs {
	// 	registrations = append(registrations,
	// 		Qual("github.com/progrium/zt100/pkg/manifold/library", "Register").Call(
	// 			Op("&").Qual(qualObj(objId), "Main").Values(),
	// 			Lit(objId),
	// 			Id("relPath").Call(Lit(path.Join(objId, "component.go"))),
	// 			Lit("fas fa-user-robot"),
	// 		),
	// 	)
	// }
	// f.Func().Id("init").Params().Block(registrations...)
	// buf := &bytes.Buffer{}
	// err := f.Render(buf)
	index := fmt.Sprintf(indexTmpl, strings.Join(imports, "\n"))
	return []byte(index), nil
}
