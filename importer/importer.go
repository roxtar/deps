package importer

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

// Returns imports for a go file
func GetImportsFile(path string) ([]string, error) {

	var result []string
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
	if err != nil {
		return result, err
	}

	for _, s := range f.Imports {
		val, err := strconv.Unquote(s.Path.Value)
		if err != nil {
			return result, err
		}
		result = append(result, val)
	}
	return result, nil
}

// Returns imports for all the files present in a go package
func GetImportsPackage(packageName string) ([]string, error) {
	var result []string

	packagePath, err := GetPackagePath(packageName)
	if err != nil {
		return result, err
	}

	importSet := make(map[string]bool)

	gofiles := GetGoFiles(packagePath)
	for _, gofile := range gofiles {
		imports, err := GetImportsFile(gofile)
		if err != nil {
			return result, err
		}

		for _, pkg := range imports {
			importSet[pkg] = true
		}
	}

	for pkg, _ := range importSet {
		result = append(result, pkg)
	}

	return result, nil
}

// Returns the package path of the given package name
func GetPackagePath(packageName string) (string, error) {
	env := os.Getenv("GOPATH")
	if env == "" {
		return "", fmt.Errorf("GOPATH not set")
	}

	gopaths := filepath.SplitList(env)
	for _, gopath := range gopaths {
		packagePath := path.Join(gopath, "src", packageName)
		_, err := os.Stat(packagePath)
		if err == nil {
			return packagePath, nil
		}
	}

	return "", fmt.Errorf("package not found: %s", packageName)
}

// Recursively gets all the go files under the given path
func GetGoFiles(dirPath string) []string {
	var gofiles []string
	filepath.Walk(dirPath, func(subPath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(subPath) != ".go" {
			return nil
		}

		gofiles = append(gofiles, subPath)
		return nil
	})
	return gofiles
}
