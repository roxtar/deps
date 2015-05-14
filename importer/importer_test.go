package importer_test

import (
	"github.com/roxtar/deps/importer"
	"os"
	"path"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Importer", func() {

	Context("GetImportsFile", func() {
		It("gets the right imports for a file", func() {
			imports, err := importer.GetImportsFile("fixtures/test_good1.go")
			Expect(err).ToNot(HaveOccurred())
			Expect(imports).To(ConsistOf("github.com/test/test1", "os/exec", "json"))
		})

		It("errors out if import statement is bad", func() {
			_, err := importer.GetImportsFile("fixtures/test_bad1.go")
			Expect(err).To(HaveOccurred())
		})

	})

	Context("GetPackagePath", func() {

		It("gets the package path", func() {
			pwd, _ := os.Getwd()
			gopath := path.Join(pwd, "fixtures")
			os.Setenv("GOPATH", gopath)

			packagePath, err := importer.GetPackagePath("github.com/test/test1")
			Expect(err).ToNot(HaveOccurred())
			Expect(packagePath).To(Equal(path.Join(gopath, "src", "github.com/test/test1")))
		})

		It("gets package path with multiple GOPATHs", func() {
			pwd, _ := os.Getwd()
			gopath1 := "/tmp/gopath"
			gopath2 := path.Join(pwd, "fixtures")
			os.Setenv("GOPATH", strings.Join([]string{gopath1, gopath2}, string(os.PathListSeparator)))

			packagePath, err := importer.GetPackagePath("github.com/test/test1")
			Expect(err).ToNot(HaveOccurred())
			Expect(packagePath).To(Equal(path.Join(gopath2, "src", "github.com/test/test1")))

		})

		It("errors out if GOPATH does not have package", func() {
			gopath := "/tmp/gopath"
			os.Setenv("GOPATH", gopath)

			_, err := importer.GetPackagePath("github.com/test/test1")
			Expect(err).To(HaveOccurred())
		})

		It("errors out if GOPATH is not set", func() {
			os.Setenv("GOPATH", "")
			_, err := importer.GetPackagePath("github.com/test/test1")
			Expect(err).To(HaveOccurred())
		})
	})

	Context("GetGoFiles", func() {
		It("gets all go files under given path", func() {
			gofiles := importer.GetGoFiles("fixtures")
			path1 := path.Join("fixtures", "test_good1.go")
			path2 := path.Join("fixtures", "test_bad1.go")
			path3 := path.Join("fixtures", "src", "github.com", "test", "test1", "test1.go")
			path4 := path.Join("fixtures", "src", "github.com", "test", "test1", "sub", "sub.go")
			Expect(gofiles).To(ConsistOf(path1, path2, path3, path4))
		})
	})

	Context("GetImportsPackage", func() {
		It("gets all imports in given package", func() {
			pwd, _ := os.Getwd()
			gopath := path.Join(pwd, "fixtures")
			os.Setenv("GOPATH", gopath)

			imports, err := importer.GetImportsPackage("github.com/test/test1")
			Expect(err).ToNot(HaveOccurred())
			Expect(imports).To(ConsistOf("os", "path", "github.com/test/test2", "filepath", "ioutils", "json"))
		})
	})

})
