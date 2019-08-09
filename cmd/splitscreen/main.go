package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/tools/imports"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Handler struct {
	File      string
	Package   string
	Aggregate string
	Events    []string
	Commands  []string
}

type Type struct {
	Pkg  string
	Name string
}

type Import struct {
	Alias string
	Path string
}

type View struct {
	Imports []Import
	Package   string
	View      string
	Listeners []Type
}

var action = flag.String("generate", "handler", "handler/view")

func init(){
	log.SetFlags(log.LstdFlags|log.Lshortfile)
}

func main() {
	flag.Parse()
	defFile, err := filepath.Abs(os.Getenv("GOFILE"))
	log.Printf("generating code from %s", defFile)
	if err != nil {
		log.Fatal(err)
	}

	lookupPaths := []string{os.Getenv("SSPATH"), filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "netbrain", "splitscreen", "cmd", "splitscreen")}
	for _, p := range lookupPaths {
		if _, err := os.Stat(p); !os.IsNotExist(err) {
			if err := os.Chdir(p); err != nil {
				log.Fatal(err)
			}
			break
		}
	}

	switch *action {
	case "view":
		err := writeView(defFile)
		if err != nil {
			log.Fatal(err)
		}
	case "handler":
		meta, err := writeBoilerplate(defFile)
		if err != nil {
			log.Fatal(err)
		}

		err = writeHandlers(meta)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("unknown action")
	}

}

func writeView(defFile string) error {
	fset := token.NewFileSet()
	src, err := ioutil.ReadFile(defFile)
	if err != nil {
		return err
	}
	f, err := parser.ParseFile(fset, defFile, src, 0)
	if err != nil {
		return err
	}

	meta := View{}
	meta.Package = f.Name.Name

	ast.Inspect(f, func(n ast.Node) bool {
		imp,ok := n.(*ast.ImportSpec)
		if ok {
			var alias string
			path := strings.Trim(imp.Path.Value,`"`)
			if imp.Name != nil && imp.Name.Name != "." {
				alias = imp.Name.Name
			}else {
				parts := strings.Split(path,"/")
				alias = parts[len(parts)-1]
			}
			meta.Imports = append(meta.Imports,Import{
				Alias: alias,
				Path:  path,
			})
		}

		typ, ok := n.(*ast.TypeSpec)
		if ok {
			if strings.HasSuffix(typ.Name.Name, "View") {
				meta.View = typ.Name.Name
			}
		}

		fn, ok := n.(*ast.FuncDecl)
		if ok {
			if strings.HasPrefix(fn.Name.Name, "On") {
				for _, f := range fn.Type.Params.List {
					ptr, ok := f.Type.(*ast.StarExpr)
					if !ok {
						continue
					}
					selc, ok := ptr.X.(*ast.SelectorExpr)
					if !ok {
						continue
					}
					pkgIdent,ok := selc.X.(*ast.Ident)
					if !ok {
						continue
					}

					event := selc.Sel.Name
					pkg := pkgIdent.Name

					meta.Listeners = append(meta.Listeners, Type{
						Pkg:  pkg,
						Name: event,
					})
				}
			}
		}
		return true
	})

	tmpl := template.Must(template.ParseGlob("./tmpl/*"))
	buffer := &bytes.Buffer{}
	err = tmpl.ExecuteTemplate(buffer, "view", &meta)
	if err != nil {
		return err
	}
	output := path.Join(path.Dir(defFile), strings.TrimSuffix(path.Base(defFile), path.Ext(defFile))+"_gen.go")
	buf, err := imports.Process(output, buffer.Bytes(), nil)
	if err != nil {
		log.Println(buffer.String())
		return err
	}
	return ioutil.WriteFile(output, buf, 0644)
}

func writeBoilerplate(defFile string) (meta Handler, err error) {
	meta.File = defFile
	fset := token.NewFileSet()
	src, err := ioutil.ReadFile(defFile)
	if err != nil {
		return
	}
	f, err := parser.ParseFile(fset, defFile, src, 0)
	if err != nil {
		return
	}
	meta.Package = f.Name.Name
	ast.Inspect(f, func(n ast.Node) bool {
		typ, ok := n.(*ast.TypeSpec)
		if ok {
			if strings.HasSuffix(typ.Name.Name, "Event") {
				meta.Events = append(meta.Events, typ.Name.Name)
			}
			if strings.HasSuffix(typ.Name.Name, "Command") {
				meta.Commands = append(meta.Commands, typ.Name.Name)
			}
			if strings.HasSuffix(typ.Name.Name, "Aggregate") {
				meta.Aggregate = typ.Name.Name
			}
		}
		return true
	})

	if meta.Aggregate == "" || len(meta.Commands) == 0 && len(meta.Events) == 0 {
		err = fmt.Errorf("missing aggregate or event/command definition(s)")
		return
	}

	tmpl := template.Must(template.ParseGlob("./tmpl/*"))

	buffer := &bytes.Buffer{}
	err = tmpl.ExecuteTemplate(buffer, "boilerplate", meta)
	if err != nil {
		return
	}

	output := path.Join(path.Dir(defFile), strings.TrimSuffix(path.Base(defFile), path.Ext(defFile))+"_gen.go")
	buf, err := imports.Process(output, buffer.Bytes(), nil)
	if err != nil {
		log.Println(buffer.String())
		return
	}
	err = ioutil.WriteFile(output, buf, 0644)
	return
}

func writeHandlers(meta Handler) error {
	output := path.Join(path.Dir(meta.File), strings.TrimSuffix(path.Base(meta.File), path.Ext(meta.File))+"_handler.go")
	if _, err := os.Stat(output); os.IsNotExist(err) {
		buffer := bytes.NewBufferString(fmt.Sprintf(`package %s`, meta.Package))
		tmpl := template.Must(template.ParseGlob("./tmpl/*"))

		if err := tmpl.ExecuteTemplate(buffer, "handler", meta); err != nil {
			return err
		}
		buf, err := imports.Process(output, buffer.Bytes(), nil)
		if err != nil {
			return err
		}
		return ioutil.WriteFile(output, buf, 0644)
	}

	fset := token.NewFileSet()
	src, err := ioutil.ReadFile(output)
	if err != nil {
		return err
	}
	f, err := parser.ParseFile(fset, meta.File, src, 0)
	ast.Inspect(f, func(n ast.Node) (ret bool) {
		ret = true
		typ, ok := n.(*ast.FuncDecl)
		if !ok {
			return
		}
		if typ.Recv == nil {
			return
		}
		fn, ok := n.(*ast.FuncDecl)

		if strings.HasPrefix(fn.Name.Name, "Apply") && strings.HasSuffix(fn.Name.Name, "Event") {
			for i := 0; i < len(meta.Events); i++ {
				if meta.Events[i] == strings.TrimPrefix(fn.Name.Name, "Apply") {
					meta.Events = append(meta.Events[:i], meta.Events[i+1:]...)
					break
				}
			}
		}
		if strings.HasPrefix(fn.Name.Name, "Handle") && strings.HasSuffix(fn.Name.Name, "Command") {
			for i := 0; i < len(meta.Commands); i++ {
				if meta.Commands[i] == strings.TrimPrefix(fn.Name.Name, "Handle") {
					meta.Commands = append(meta.Commands[:i], meta.Commands[i+1:]...)
					break
				}
			}
		}
		return
	})
	inBuf, err := ioutil.ReadFile(output)
	if err != nil {
		return err
	}
	buffer := bytes.NewBuffer(inBuf)
	tmpl := template.Must(template.ParseGlob("./tmpl/*"))

	if err := tmpl.ExecuteTemplate(buffer, "handler_partial", meta); err != nil {
		return err
	}

	buf, err := imports.Process(output, buffer.Bytes(), nil)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(output, buf, 0644)
}
