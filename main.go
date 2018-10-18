package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"
	"strings"
	"text/template"
)

const (
	helpMsg = "testGen - usage: testGen [-o output] (input regex)"
)

type Param struct {
	Name string
	Type string
}

type Result Param

type Receiver Param

type Func struct {
	Name    string
	Params  []Param
	Results []Result
}

type Method struct {
	Name    string
	Rcvr    Receiver
	Params  []Param
	Results []Result
}

var testFuncTemplate = template.Must(template.New("").
	Parse(`
func Test{{ .Name }}(t *testing.T) {
	type (
		Input struct { {{- range .Params }}
			{{ .Name }} {{ .Type }}{{ end }}
		}

		Output struct { {{- range .Results }}
			{{ .Name }} {{ .Type }}{{ end }}
		}

		Case struct {
			Name   string
			In     Input
			Expect Output
		}
	)

	// TODO create test cases
	cases := []Case{}

	tester := func(c Case) func(*testing.T) {
		return func(t *testing.T) {
			var actual Output
			{{ range $idx, $ret := .Results }}{{ if $idx }}, {{ end }}actual.{{ $ret.Name }}{{ end }} = {{ .Name }}({{ range $idx, $ret := .Params }}{{ if $idx }}, {{ end }}c.In.{{ $ret.Name }}{{ end }})

			if actual != c.Expect {
				t.Errorf("expected %+v, actual: %+v\n", c.Expect, actual)
			}
		}
	}

	for _, c := range cases {
		t.Run(c.Name, tester(c))
	}
}
`))

// TODO method Receiver stuff
var testMethodTemplate = template.Must(template.New("").
	Parse(`
func Test{{ .Name }}(t *testing.T) {
	type (
		Input struct { {{- range .Params }}
			{{ .Name }} {{ .Type }}{{ end }}
		}

		Output struct { {{- range .Results }}
			{{ .Name }} {{ .Type }}{{ end }}
		}

		Case struct {
			Name   string
			In     Input
			Expect Output
		}
	)

	// TODO create test cases
	cases := []Case{}

	tester := func(c Case) func(*testing.T) {
		return func(t *testing.T) {
			var actual Output
			{{ range $idx, $ret := .Results }}{{ if $idx }}, {{ end }}actual.{{ $ret.Name }}{{ end }} = {{ .Name }}({{ range $idx, $ret := .Params }}{{ if $idx }}, {{ end }}c.In.{{ $ret.Name }}{{ end }})

			if actual != c.Expect {
				t.Errorf("expected %+v, actual: %+v\n", c.Expect, actual)
			}
		}
	}

	for _, c := range cases {
		t.Run(c.Name, tester(c))
	}
}
`))

var (
	output      string
	packagePath string
	input       string
)

func init() {
	flag.StringVar(&output, "o", "stdout", "specify where to output. stdout is the default")
	flag.StringVar(&output, "output", "stdout", "specify where to output. stdout is the default")

	flag.StringVar(&packagePath, "p", "pwd", "the absolute or GOPATH relative path to the package to generate tests for")
	flag.StringVar(&packagePath, "package", "pwd", "the absolute or GOPATH relative path to the package to generate tests for")
}

func main() {
	flag.CommandLine.Init("testGen", flag.ContinueOnError)
	err := flag.CommandLine.Parse(os.Args[1:])
	if err != nil {
		if err == flag.ErrHelp {
			fmt.Println(helpMsg)
			flag.CommandLine.PrintDefaults()
		} else {
			fmt.Println(err)
		}
		os.Exit(1)
	}

	input = flag.CommandLine.Arg(0)
	if input == "" {
		// match anything
		input = `\.*`
	}

	var outputFile *os.File
	if output == "stdout" {
		outputFile = os.Stdout
	} else {
		outputFile, err = os.OpenFile(output, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			fmt.Printf("could not open '%s': %v\n", output, err)
			os.Exit(1)
		}
		defer outputFile.Close()
	}

	if packagePath == "pwd" {
		packagePath, err = os.Getwd()
		if err != nil {
			fmt.Println("could not get working directory:", err)
			os.Exit(1)
		}
	} else if !path.IsAbs(packagePath) {
		// it's relative to $GOPATH
		packagePath := os.ExpandEnv(path.Join("$GOPATH", packagePath))
		if !path.IsAbs(packagePath) {
			fmt.Println("GOPATH does not seem to be set - cannot use GOPATH relative package paths")
			os.Exit(1)
		}
	}

	filter, err := regexp.Compile(input)
	if err != nil {
		fmt.Println("bad regex:", err)
		os.Exit(1)
	}

	fset := token.NewFileSet()
	packages, err := parser.ParseDir(fset, packagePath, filesToIgnore, 0)
	if err != nil {
		fmt.Println("error parsing package:", err)
		os.Exit(1)
	}

	for _, pkg := range packages {
		for _, file := range pkg.Files {
			for _, obj := range file.Scope.Objects {
				if obj.Kind != ast.Fun {
					continue
				}

				if !obj.Pos().IsValid() {
					fmt.Println("skipping invalid", obj.Name)
					continue
				}

				if !filter.MatchString(obj.Name) {
					continue
				}

				if packageContains(pkg, "Test"+obj.Name) {
					continue
				}

				fn := obj.Decl.(*ast.FuncDecl)

				// no point in testing functions without input (?)
				if fn.Type.Params.NumFields() == 0 {
					continue
				}

				// no point in testing functions without output (?)
				if fn.Type.Results == nil || fn.Type.Results.NumFields() == 0 {
					continue
				}

				var (
					ps []Param
					rs []Result
				)

				for _, p := range fn.Type.Params.List {
					for _, name := range p.Names {
						ps = append(ps, Param{
							Name: name.Name,
							Type: typename(p.Type),
						})
					}
				}

				for _, r := range fn.Type.Results.List {
					typeName := typename(r.Type)

					if r.Names == nil {
						rs = append(rs, Result{
							Name: strings.ToLower(typeName),
							Type: typeName,
						})
						continue
					}

					for _, name := range r.Names {
						rs = append(rs, Result{
							Name: name.Name,
							Type: typeName,
						})
					}
				}

				if fn.Recv == nil {
					f := Func{
						Name:    fn.Name.Name,
						Params:  ps,
						Results: rs,
					}
					err = testFuncTemplate.Execute(outputFile, f)
				} else {
					m := Method{
						Name: fn.Name.Name,
						Rcvr: Receiver{
							Name: fn.Recv.List[0].Names[0].Name,
							Type: typename(fn.Recv.List[0].Type),
						},
						Params:  ps,
						Results: rs,
					}

					err = testMethodTemplate.Execute(outputFile, m)
				}

				if err != nil {
					fmt.Printf("creating test function for %s: %v\n", fn.Name, err)
					os.Exit(1)
				}
			}
		}
	}

	if goimports := findGoImports(); goimports != "" {
		var cmd *exec.Cmd
		if outputFile != os.Stdout {
			cmd = exec.Command(goimports, "-w", outputFile.Name())
		} else {
			cmd = exec.Command(goimports, "-w")
		}

		err = cmd.Run()
		if err != nil {
			fmt.Println("goimports:", err)
		}
	} else {
		fmt.Println("Could not find goimports, you will have to handle imports yourself")
	}
}

func packageContains(pkg *ast.Package, name string) bool {
	for _, file := range pkg.Files {
		for _, obj := range file.Scope.Objects {
			if obj.Name == name {
				return true
			}
		}
	}
	return false
}

func filesToIgnore(info os.FileInfo) bool {
	if strings.HasSuffix(info.Name(), "_test.go") {
		return false
	}
	return true
}

func findGoImports() string {
	binName := "goimports"
tryAgain:
	p, _ := exec.LookPath(binName)
	if p != "" {
		return binName
	}

	p, _ = exec.LookPath(os.ExpandEnv("$GOPATH/bin/" + binName))
	if p != "" {
		return p
	}

	if runtime.GOOS == "windows" {
		binName += ".exe"
		goto tryAgain
	}

	return ""
}

func typename(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return t.Sel.Name
	case *ast.StarExpr:
		return typename(t.X)
	default:
		fmt.Printf("%T\n", t)
		return ""
	}
}
